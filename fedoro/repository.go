package fedoro

import (
	"io"
	"time"
)

var (
	MainRepo Repository
)

type Repository interface {
	// Find the object with the given identifier.
	// errors: could not find pid, or possibly other errors
	FindPid(pid string) (DigitalObject, error)
	// Create a new object and return a DigitalObject which represents it
	// It is the callers responsibility to choose the PID of the new object
	NewObject(obj ObjectInfo) (DigitalObject, error)
}

type DigitalObject interface {
	Info() *ObjectInfo
	DsNames() []string

	// Get information about the datastream `dsid` at version `version`. Some of the info is general to the datastream,
	// other is specific to the version asked for.
	// Set version to -1 to get information about the newest version
	DsInfo(dsid string, version int) *DatastreamInfo
	// Get the content for datastream dsid at the given version (0 is first version, 1 is next, etc.)
	// Set version to -1 to always get the newest version
	DsContent(dsid string, version int) (io.ReadCloser, error)

	/*
	   // Update the object's info to what is given
	   UpdateInfo(obj *ObjectInfo) error
	   // Create a new datastream
	   AddDatastream(dsinfo *DatastreamInfo) error
	   // create a new version of a datastream
	   ReplaceContent(dsid string, r io.Reader) error
	*/
}

// The general information about a digital object
type ObjectInfo struct {
	Pid          string // The object's identifier
	Version      string // ? look up
	State        string // The state of the object. One of "A", "I", or "D"
	Label        string // The human readable label for the object
	OwnerId      string
	CreatedDate  time.Time
	ModifiedDate time.Time
}

type DatastreamInfo struct {
	Name         string // The name of the datastream, e.g. "RELS-EXT" (as opposed to a versioned name e.g. "RELS-EXT.0")
	State        rune   // 'A'ctive, 'I'nactive, or 'D'eleted
	ControlGroup rune   // 'X', 'M', 'E'
	Versionable  bool   // Does this ds keep previous versions of its content? true == yes.
	NumVersions  int    // >= 1
	// All the following only refer to the current version of the datastream
	Id         string // The full identity of this version of the datastream, e.g. "RELS-EXT.0"
	Label      string // The supplied human readable name for the content
	Created    time.Time
	Mimetype   string
	Format_uri string // Not sure what this is
	Size       uint
}
