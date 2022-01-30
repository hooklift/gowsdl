package soap

import (
	"encoding/xml"
	"time"
)

type (
	// GDay is a type for representing xsd:gDay in Golang
	GDay time.Time
	// GMonth is a type for representing xsd:gMonth in Golang
	GMonth time.Time
	// GMonthDay is a type for representing xsd:gMonthDay in Golang
	GMonthDay time.Time
	// GYear is a type for representing xsd:gYear in Golang
	GYear time.Time
)

// MarshalXML implements xml.Marshaler on GDay
func (gDay GDay) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	str := gDay.String()
	if str != "" {
		return e.EncodeElement(str, start)
	}
	return nil
}

// MarshalXML implements xml.Marshaler on GMonth
func (gMonth GMonth) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	str := gMonth.String()
	if str != "" {
		return e.EncodeElement(str, start)
	}
	return nil
}

// MarshalXML implements xml.Marshaler on GMonthDay
func (gMonthDay GMonthDay) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	str := gMonthDay.String()
	if str != "" {
		return e.EncodeElement(str, start)
	}
	return nil
}

// MarshalXML implements xml.Marshaler on GYear
func (gYear GYear) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	str := gYear.String()
	if str != "" {
		return e.EncodeElement(str, start)
	}
	return nil
}

func (gDay GDay) String() string {
	return time.Time(gDay).Format(gDay.timeFormat())
}

func (gDay GDay) timeFormat() string {
	return "02-07:00"
}

func (gMonth GMonth) String() string {
	return time.Time(gMonth).Format(gMonth.timeFormat())
}

func (gMonth GMonth) timeFormat() string {
	return "--01-07:00"
}

func (gMonthDay GMonthDay) String() string {
	return time.Time(gMonthDay).Format(gMonthDay.timeFormat())
}

func (gMonthDay GMonthDay) timeFormat() string {
	return "--02-01-07:00"
}

func (gYear GYear) String() string {
	return time.Time(gYear).Format(gYear.timeFormat())
}

func (gYear GYear) timeFormat() string {
	return "2006-07:00"
}
