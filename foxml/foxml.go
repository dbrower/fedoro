/*
Provides utility routines for loading digital objects from
a foxml file.
*/
package foxml

import (
	"encoding/xml"
	"errors"
	"io"
	"time"
)

var (
	ErrNoFoxml = errors.New("foxml: No digitalObject in Foxml")
)

type FoxmlProperty struct {
	Name  string `xml:"NAME,attr"`
	Value string `xml:"VALUE,attr"`
}

type ContentLocation struct {
	Kind string `xml:"TYPE,attr"`
	Ref  string `xml:"REF,attr"`
}

type XmlContent struct {
	Contents string `xml:",innerxml"`
}

type DatastreamVersion struct {
	Id              string          `xml:"ID,attr"`
	Label           string          `xml:"LABEL,attr"`
	Created         time.Time       `xml:"CREATED,attr"`
	Mimetype        string          `xml:"MIMETYPE,attr"`
	Format_uri      string          `xml:"FORMAT_URI,attr"`
	Size            int             `xml:"SIZE,attr"`
	XmlContent      XmlContent      `xml:"xmlContent"`
	ContentLocation ContentLocation `xml:"contentLocation"`
}

type Datastream struct {
	Id           string              `xml:"ID,attr"`
	State        rune                `xml:"STATE,attr"`
	ControlGroup rune                `xml:"CONTROL_GROUP,attr"`
	Versionable  bool                `xml:"VERSIONABLE,attr"`
	Versions     []DatastreamVersion `xml:"datastreamVersion"`
}

type DigitalObject struct {
	Version      string `xml:"VERSION,attr"`
	Pid          string `xml:"PID,attr"`
	State        string
	Label        string
	OwnerId      string
	CreatedDate  time.Time
	ModifiedDate time.Time
	Properties   []FoxmlProperty `xml:"objectProperties>property"`
	Ds           []Datastream    `xml:"datastream"`
}

// fill in the structure fields from the list of properties
// Modifies the digital object
func (do *DigitalObject) fixupDigitalObject() error {
	for _, p := range do.Properties {
		switch p.Name {
		case "info:fedora/fedora-system:def/model#state":
			do.State = p.Value
		case "info:fedora/fedora-system:def/model#label":
			do.Label = p.Value
		case "info:fedora/fedora-system:def/model#ownerId":
			do.OwnerId = p.Value
		case "info:fedora/fedora-system:def/model#createdDate":
			var err error
			do.CreatedDate, err = time.Parse("2006-01-02T15:04:05Z", p.Value)
			if err != nil {
				return err
			}
		case "info:fedora/fedora-system:def/view#lastModifiedDate":
			var err error
			do.ModifiedDate, err = time.Parse("2006-01-02T15:04:05Z", p.Value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// decode the foxml data provided by f. Returns a DigitalObject
// or an error. Uses the encoding/xml routines to parse the xml.
func DecodeFoxml(f io.Reader) (DigitalObject, error) {
	var do DigitalObject
	d := xml.NewDecoder(f)
	for {
		t, err := d.Token()
		if err != nil {
			return do, err
		}
		if v, ok := t.(xml.StartElement); ok {
			if v.Name.Local == "digitalObject" {
				var do DigitalObject
				err := d.DecodeElement(&do, &v)
				if err != nil {
					return do, err
				}
				err = do.fixupDigitalObject()
				if err != nil {
					return do, err
				}
				return do, nil
			}
		}
	}
	return do, ErrNoFoxml
}
