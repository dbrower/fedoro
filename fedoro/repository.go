
package fedoro

import (
    "github.com/dbrower/fedoro/akubra"
    "github.com/dbrower/fedoro/foxml"
)

type DatastreamVersion struct {
    Id strin
    Label strin
    Created time.Time
    Mimetype string
    Format_uri string
    Size uint
    XmlContent XmlContent
    ContentLocation ContentLocation
}

type Datastream struct {
    Id string           `xml:"ID,attr"`
    State string        `xml:"STATE,attr"`
    ControlGroup string `xml:"CONTROL_GROUP,attr"`
    Versionable bool    `xml:"VERSIONABLE,attr"`
    Versions []DatastreamVersion `xml:"datastreamVersion"`
}

type DigitalObject struct {
	Pid string
  State string
    Label string
    OwnerId string
    CreatedDate time.Time
    ModifiedDate time.Time
    Ds []Datastream     `xml:"datastream"`
}

type Repository interface {
    FindPid(string) (DigitalObject, error)
}

