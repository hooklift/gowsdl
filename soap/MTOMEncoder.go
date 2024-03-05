package soap

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime"
	"mime/multipart"
	"net/textproto"
	"reflect"
	"strings"
)

type mtomEncoder struct {
	writer *multipart.Writer
}

// Binary enables binary data to be enchanged in MTOM mode with XOP encoding
// When MTOM is not used, the field is encoded in Base64
type Binary struct {
	content     *[]byte
	contentType string
	packageID   string
	useMTOM     bool
}

type xopPlaceholder struct {
	Href string `xml:"href,attr"`
}

// NewBinary allocate a new Binary backed by the given byte slice, an auto-generated packageID and no MTOM-usage
func NewBinary(v []byte) *Binary {
	return &Binary{&v, "application/octet-stream", "", false}
}

// Bytes returns a slice backed by the content of the field
func (b *Binary) Bytes() []byte {
	return *b.content
}

// SetUseMTOM activates the XOP transformation of binaries in MTOM requests
func (b *Binary) SetUseMTOM(useMTOM bool) *Binary {
	b.useMTOM = useMTOM
	return b
}

// SetPackageID sets and overrides the default auto-generated package ID to be used for the multipart binary
func (b *Binary) SetPackageID(packageID string) *Binary {
	b.packageID = packageID
	return b
}

// SetContentType sets the content type the content will be transmitted as multipart
func (b *Binary) SetContentType(contentType string) *Binary {
	b.contentType = contentType
	return b
}

// ContentType returns the content type
func (b *Binary) ContentType() string {
	return b.contentType
}

// MarshalXML implements the xml.Marshaler interface to encode a Binary to XML
func (b *Binary) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	if b.useMTOM {
		if b.packageID == "" {
			b.packageID = fmt.Sprintf("%d", rand.Int())
		}
		return enc.EncodeElement(struct {
			Include *xopPlaceholder `xml:"http://www.w3.org/2004/08/xop/include Include"`
		}{
			Include: &xopPlaceholder{
				Href: fmt.Sprintf("cid:%s", b.packageID),
			},
		}, start)
	}
	return enc.EncodeElement(b.Bytes, start)
}

// UnmarshalXML implements the xml.Unmarshaler interface to decode a Binary form XML
func (b *Binary) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	ref := struct {
		Content []byte          `xml:",innerxml"`
		Include *xopPlaceholder `xml:"http://www.w3.org/2004/08/xop/include Include"`
	}{}

	if err := d.DecodeElement(&ref, &start); err != nil {
		return err
	}

	b.content = &ref.Content

	if ref.Include != nil {
		b.packageID = strings.TrimPrefix(ref.Include.Href, "cid:")
		b.useMTOM = true
	}
	return nil
}

func newMtomEncoder(w io.Writer) *mtomEncoder {
	return &mtomEncoder{
		writer: multipart.NewWriter(w),
	}
}

func (e *mtomEncoder) Encode(v interface{}) error {
	binaryFields := make([]reflect.Value, 0)
	getBinaryFields(v, &binaryFields)
	enableMTOMMode(binaryFields)

	var partWriter io.Writer
	var err error

	h := make(textproto.MIMEHeader)
	h.Set("Content-Type", `application/xop+xml`)
	h.Set("Content-Transfer-Encoding", "8bit")

	if partWriter, err = e.writer.CreatePart(h); err != nil {
		return err
	}
	xmlEncoder := xml.NewEncoder(partWriter)
	if err := xmlEncoder.Encode(v); err != nil {
		return err
	}

	for _, fld := range binaryFields {
		pkg := fld.Interface().(*Binary)
		h := make(textproto.MIMEHeader)
		if pkg.contentType == "" {
			pkg.contentType = "application/octet-stream"
		}
		h.Set("Content-Type", pkg.contentType)
		h.Set("Content-ID", fmt.Sprintf("<%s>", pkg.packageID))
		h.Set("Content-Transfer-Encoding", "binary")
		if partWriter, err = e.writer.CreatePart(h); err != nil {
			return err
		}
		partWriter.Write(*pkg.content)
	}

	return nil
}

func (e *mtomEncoder) Flush() error {
	return e.writer.Close()
}

func getBinaryFields(data interface{}, fields *[]reflect.Value) {
	v := reflect.Indirect(reflect.ValueOf(data))

	switch v.Kind() {
	case reflect.Ptr:
		getBinaryFields(v.Elem(), fields)
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !v.Field(i).CanInterface() {
				continue
			}
			f := v.Field(i)
			if _, ok := f.Interface().(*Binary); ok {
				*fields = append(*fields, f)
			} else {
				getBinaryFields(f.Interface(), fields)
			}
		}
		break
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			e := v.Index(i)
			if !e.CanInterface() {
				continue
			}
			if _, ok := e.Interface().(*Binary); ok {
				*fields = append(*fields, v)
			} else {
				getBinaryFields(e.Interface(), fields)
			}
		}
	}
}

func enableMTOMMode(fields []reflect.Value) {
	for _, f := range fields {
		b := f.Interface().(*Binary)
		b.useMTOM = true
	}
}

func (e *mtomEncoder) Boundary() string {
	return e.writer.Boundary()
}

type mtomDecoder struct {
	reader *multipart.Reader
}

func getMtomHeader(contentType string) (string, error) {
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		boundary, ok := params["boundary"]
		if !ok || boundary == "" {
			return "", fmt.Errorf("Invalid multipart boundary: %s", boundary)
		}

		cType, ok := params["type"]
		if !ok || cType != "application/xop+xml" {
			// Process as normal xml (Don't resolve XOP parts)
			return "", nil
		}

		startInfo, ok := params["start-info"]
		if !ok || (!strings.Contains(startInfo, "application/soap+xml") && !strings.Contains(startInfo, "text/xml")) {
			return "", fmt.Errorf(`Expected param start-info to contain "application/soap+xml" or "text/xml", got %s`, startInfo)
		}
		return boundary, nil
	}

	return "", nil
}

func newMtomDecoder(r io.Reader, boundary string) *mtomDecoder {
	return &mtomDecoder{
		reader: multipart.NewReader(r, boundary),
	}
}

func (d *mtomDecoder) Decode(v interface{}) error {
	packages := make(map[string]*Binary, 0)
	for {
		p, err := d.reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		contentType := p.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "application/xop+xml") {
			err := xml.NewDecoder(p).Decode(v)
			if err != nil {
				return err
			}
		} else {
			contentID := p.Header.Get("Content-Id")
			if contentID == "" {
				return errors.New("Invalid multipart content ID")
			}
			content, err := io.ReadAll(p)
			if err != nil {
				return err
			}

			contentID = strings.Trim(contentID, "<>")
			packages[contentID] = &Binary{
				content:     &content,
				contentType: contentType,
			}
		}
	}

	// Get binary fields after reading the structure, so I have all fields mapped
	fields := make([]reflect.Value, 0)
	getBinaryFields(v, &fields)

	// Set binary fields with correct content
	for _, f := range fields {
		b := f.Interface().(*Binary)
		b.content = packages[b.packageID].content
		b.contentType = packages[b.packageID].contentType
	}
	return nil
}
