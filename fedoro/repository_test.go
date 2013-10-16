
package fedoro

import (
    "testing"

    "github.com/dbrower/fedoro/akubra"
)

func newRepo() Repository {
    var obj = akubra.Pool{Root: ".", Format: "##" }
    var ds = akubra.Pool{Root: ".", Format: "##" }
    return NewRepository(obj, ds)
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
