package fedoro

import (
	"testing"
)

func newRepo() Repository {
	return akubra.NewRepository("test-repo", "test-repo")
}

func TestNewRepository(t *testing.T) {
	_ = newRepo()
}

func TestFindPid(t *testing.T) {
	r := newRepo()
	do, err := r.FindPid("vecnet:x920fw85p")
	if err != nil {
		t.Fatalf("Got %s", err)
	}
	if do.Pid != "vecnet:x920fw85p" {
		t.Fatalf("Wrong pid: %s", do.Pid)
	}
}
