package fedoro

import (
	"strings"
	"testing"
)

func setupRepo(t *testing.T) Repository {
	r := NewMemRepo()
	obj := ObjectInfo{
		Pid: "vecnet:x920fw85p",
	}
	do, err := r.NewObject(obj)
	if err != nil {
		t.Fatal(err)
	}
	dsinfo := DatastreamInfo{
		Name:         "test",
		ControlGroup: 'M',
		Versionable:  true,
	}
	err = do.UpdateDatastream(&dsinfo)
	if err != nil {
		t.Fatal(err)
	}
	err = do.ReplaceContent("test", strings.NewReader("This is a test datastream"))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestListDatastreams(t *testing.T) {
	r := setupRepo(t)
	ListDatastreams(r, "vecnet:x920fw85p")
}
