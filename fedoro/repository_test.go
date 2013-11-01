package fedoro

import (
	"testing"
)

func TestFindPid(t *testing.T) {
	r := setupRepo(t)
	_, err := r.FindPid("dummy:1")
	if err != nil {
		t.Fatalf("Got %s", err)
	}
}
