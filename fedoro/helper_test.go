package fedoro

import (
	"strings"
	"testing"
)

func setupRepo(t *testing.T) Repository {
	r := NewMemRepo()
	do := makeObject(t, r, "dummy:1")
	addContent(t, do, "test", "This is a test datastream")
	addContent(t, do, "test", "This is another test datastream")
	return r
}

func makeObject(t *testing.T, r Repository, pid string) DigitalObject {
	do, err := r.NewObject(ObjectInfo{Pid: pid})
	if err != nil {
		t.Fatal(err)
	}
	return do
}

func addContent(t *testing.T, do DigitalObject, dsid string, content string) {
	err := do.ReplaceContent(dsid, strings.NewReader(content))
	if err != nil {
		t.Fatal(err)
	}
}
