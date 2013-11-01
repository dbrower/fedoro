package fedoro

import (
	"testing"
)

func TestListDatastreams(t *testing.T) {
	r := setupRepo(t)
	result, _ := ListDatastreams(r, "dummy:1")
	t.Log(result)
}
