package soap

import (
	"encoding/xml"
	"time"
)

const (
	dateLayout = "2006-01-02"
)

// XsdDateTime is a type for representing xsd:datetime
type XsdDateTime struct {
	time.Time
}

// MarshalXML implementation on DateTimeg to skip "zero" time values
func (xdt XsdDateTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if !xdt.IsZero() {
		e.EncodeElement(xdt.Time.Format(time.RFC3339), start)
	}
	return nil
}

// UnmarshalXML implementation on DateTimeg to use dateTimeLayout
func (xdt *XsdDateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var err error
	xdt.Time, err = unmarshalTime(d, start, time.RFC3339)
	return err
}

func unmarshalTime(d *xml.Decoder, start xml.StartElement, format string) (time.Time, error) {
	var t time.Time
	var content string
	err := d.DecodeElement(&content, &start)
	if err != nil {
		return t, err
	}
	if content == "" {
		return t, nil
	}
	if content == "0001-01-01T00:00:00Z" {
		return t, nil
	}
	t, err = time.Parse(format, content)
	if err != nil {
		return t, err
	}
	return t, nil
}

// XsdDate is a type for representing xsd:date
type XsdDate struct {
	time.Time
}

// MarshalXML implementation on DateTimeg to skip "zero" time values
func (xd XsdDate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if !xd.IsZero() {
		e.EncodeElement(xd.Time.Format(dateLayout), start)
	}
	return nil
}

// UnmarshalXML implementation on DateTimeg to use dateTimeLayout
func (xd *XsdDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var err error
	xd.Time, err = unmarshalTime(d, start, dateLayout)
	return err
}
