package soap

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	dateLayout     = "2006-01-02"
	timeLayout     = "12:13:14"
	timeFracLayout = ".999999999"
	timeZoneLayout = "Z07:00"
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

// XsdTime is a type for representing xsd:time
type XsdTime struct {
	Hour     int
	Minute   int
	Second   int
	Fraction int
	Tz       *time.Location
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// MarshalXML implementation on DateTimeg to skip "zero" time values
func (xt XsdTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	timeString := fmt.Sprintf("%v:%v:%v", xt.Hour, xt.Minute, xt.Second)
	if xt.Fraction != 0 {
		timeString = fmt.Sprintf("%v.%v", timeString, xt.Fraction)
	}
	if xt.Tz != nil {
		_, offset := time.Now().In(xt.Tz).Zone()
		hrOffset := offset / 3600
		minOffset := abs((offset - hrOffset*3600) / 60)
		if hrOffset == 0 && minOffset == 0 {
			timeString = timeString + "Z"
		} else {
			timeString = fmt.Sprintf("%v%+03d:%02d", timeString, hrOffset, minOffset)
		}
	}
	e.EncodeElement(timeString, start)
	return nil
}

// UnmarshalXML implementation on DateTimeg to use dateTimeLayout
func (xt *XsdTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var err error
	var content string
	err = d.DecodeElement(&content, &start)
	if err != nil {
		return err
	}
	tok1 := strings.SplitN(content, ":", 3)
	if len(tok1) != 3 {
		return fmt.Errorf("Failed to parse time %v", content)
	}
	xt.Hour, err = strconv.Atoi(tok1[0])
	if err != nil {
		return err
	}
	xt.Minute, err = strconv.Atoi(tok1[1])
	if err != nil {
		return err
	}
	xt.Second, err = strconv.Atoi(tok1[2][:2])
	if err != nil {
		return err
	}
	reststr := tok1[2][2:]
	if len(reststr) > 0 && reststr[0] == '.' {
		// there is a fraction
		frac := ""
		for i, c := range reststr[1:] {
			if unicode.IsDigit(c) {
				frac = frac + string(c)
			} else {
				xt.Fraction, err = strconv.Atoi(frac)
				if err != nil {
					return err
				}
				reststr = reststr[i:]
				break
			}
		}
	}
	if len(reststr) > 0 && reststr[0] == 'Z' {
		xt.Tz = time.UTC
	} else if len(reststr) > 0 {
		sign := 1
		if reststr[0] == '+' {
			sign = 1
		} else if reststr[0] == '-' {
			sign = -1
		} else {
			return fmt.Errorf("timezone format is incorrect %v", content)
		}
		reststr = reststr[1:]
		tok2 := strings.Split(reststr, ":")
		if len(tok2) > 1 {
			hrOffset, err := strconv.Atoi(tok2[0])
			if err != nil {
				return err
			}
			minOffset, err := strconv.Atoi(tok2[1])
			if err != nil {
				return err
			}
			xt.Tz = time.FixedZone(reststr, sign*hrOffset*3600+sign*minOffset*60)
		}
	}
	return nil
}
