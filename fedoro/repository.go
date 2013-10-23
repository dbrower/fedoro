package fedoro

import (
	//	"time"
	//    "github.com/dbrower/fedoro/akubra"
	"github.com/dbrower/fedoro/foxml"
)

var (
	MainRepo Repository
)

type DigitalObject interface {
    Info() *ObjectInfo
    DsNames() []string
    GetDsInfo(dsid string, version int) *DatastreamInfo
    GetDsContent(dsid string, version int) (io.Reader, error)
}

// The general information about a digital object
type ObjectInfo struct {
    // The object's identifier
    Pid          string
    // ? look up
    Version      string
    // The state of the object. One of "A", "I", or "D"
    State        string
    // The human readable label for the object
	Label        string
	OwnerId      string
	CreatedDate  time.Time
	ModifiedDate time.Time
}

type DatastreamInfo struct {
    // The name of the datastream, e.g. "RELS-EXT"
	Name            string
    // ?? lookup. one of 'A', 'D' 
	State           rune
	ControlGroup    rune
    // Does this ds keep previous versions of its content. true == yes.
	Versionable     bool
    NumVersions     int
    // All the following only refer to the current version of the datastream
    // The full identity of this version of the datastream, e.g. "RELS-EXT.0"
	Id              string
	Label           string
	Created         time.Time
	Mimetype        string
	Format_uri      string
	Size            uint
}

type Repository interface {
	FindPid(string) (*DigitalObject, error)
}
