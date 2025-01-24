package soap

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/textproto"
	"strings"
)

const mmaContentType string = `multipart/related; start="<soaprequest@gowsdl.lib>"; type="text/xml"; boundary="%s"`

type mmaEncoder struct {
	writer      *multipart.Writer
	attachments []MIMEMultipartAttachment
}

type mmaDecoder struct {
	reader *multipart.Reader
}

func newMmaEncoder(w io.Writer, attachments []MIMEMultipartAttachment) *mmaEncoder {
	return &mmaEncoder{
		writer:      multipart.NewWriter(w),
		attachments: attachments,
	}
}

func newMmaDecoder(r io.Reader, boundary string) *mmaDecoder {
	return &mmaDecoder{
		reader: multipart.NewReader(r, boundary),
	}
}

func (e *mmaEncoder) Encode(v interface{}) error {
	var err error
	var soapPartWriter io.Writer

	// 1. write SOAP envelope part
	headers := make(textproto.MIMEHeader)
	headers.Set("Content-Type", `text/xml;charset=UTF-8`)
	headers.Set("Content-Transfer-Encoding", "8bit")
	headers.Set("Content-ID", "<soaprequest@gowsdl.lib>")
	if soapPartWriter, err = e.writer.CreatePart(headers); err != nil {
		return err
	}
	xmlEncoder := xml.NewEncoder(soapPartWriter)
	if err := xmlEncoder.Encode(v); err != nil {
		return err
	}

	// 2. write attachments parts
	for _, attachment := range e.attachments {
		attHeader := make(textproto.MIMEHeader)
		attHeader.Set("Content-Type", fmt.Sprintf("application/octet-stream; name=%s", attachment.Name))
		attHeader.Set("Content-Transfer-Encoding", "binary")
		attHeader.Set("Content-ID", fmt.Sprintf("<%s>", attachment.Name))
		attHeader.Set("Content-Disposition",
			fmt.Sprintf("attachment; name=\"%s\"; filename=\"%s\"", attachment.Name, attachment.Name))
		var attachmentPartWriter io.Writer
		attachmentPartWriter, err := e.writer.CreatePart(attHeader)
		if err != nil {
			return err
		}
		_, err = io.Copy(attachmentPartWriter, bytes.NewReader(attachment.Data))
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *mmaEncoder) Flush() error {
	return e.writer.Close()
}

func (e *mmaEncoder) Boundary() string {
	return e.writer.Boundary()
}

func getMmaHeader(contentType string) (string, error) {
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		boundary, ok := params["boundary"]
		if !ok || boundary == "" {
			return "", fmt.Errorf("invalid multipart boundary: %s", boundary)
		}

		//	startInfo, ok := params["start"]
		//	if !ok || startInfo != "<soaprequest@gowsdl.lib>" {
		//		return "", fmt.Errorf(`expected param start="<soaprequest@gowsdl.lib>", got %s`, startInfo)
		//	}
		return boundary, nil
	}

	return "", nil
}

func (d *mmaDecoder) Decode(v interface{}) error {
	soapEnvResp := v.(*SOAPEnvelopeResponse)
	attachments := make([]MIMEMultipartAttachment, 0)
	for {
		p, err := d.reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		contentType := p.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "text/xml") {
			// decode SOAP part
			err := xml.NewDecoder(p).Decode(v)
			if err != nil {
				return err
			}
		} else {
			// decode attachment parts
			var r io.Reader = p

			// Handle part encoding
			if p.Header.Get("Content-Transfer-Encoding") == "base64" {
				r = base64.NewDecoder(base64.StdEncoding, p)
			}

			content, err := io.ReadAll(r)
			if err != nil {
				return err
			}

			attachments = append(attachments, MIMEMultipartAttachment{
				Name: strings.Trim(p.Header.Get("Content-Id"), "<>"),
				Data: content,
			})
		}
	}
	if len(attachments) > 0 {
		soapEnvResp.Attachments = attachments
	}

	return nil
}
