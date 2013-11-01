package fedoro

import (
	"testing"

	"strings"
)

func TestReplaceContent(t *testing.T) {
	table := []struct{ dsid, content string }{
		{"test", "abcdefghijklmnop"},
		{"test", "another test string"},
		{"xxx", "test xxx datastream"},
		{"test", "third test string"},
	}

	mdo := &MemDigitalObject{}
	buffer := make([]byte, 100)

	for i := range table {
		dsid := table[i].dsid
		content := table[i].content
		err := mdo.ReplaceContent(dsid, strings.NewReader(content))
		if err != nil {
			t.Fatal(err)
		}
		r, err := mdo.DsContent(dsid, -1)
		if err != nil {
			t.Fatal(err)
		}
		n, err := r.Read(buffer)
		if err != nil {
			t.Fatal(err)
		}
		data := string(buffer[:n])
		if data != content {
			t.Errorf("Expected %v, got %v", table[i], data)
		}
	}
}

func TestAddIfNew(t *testing.T) {
	table := []struct {
		list   []string
		text   string
		result []string
	}{
		{[]string{}, "a", []string{"a"}},
		{[]string{"a"}, "b", []string{"a", "b"}},
		{[]string{"a"}, "a", []string{"a"}},
	}
	for i := range table {
		result := addIfNew(table[i].list, table[i].text)
		if len(result) != len(table[i].result) {
			t.Errorf("%v, got %v", table[i], result)
		}
	}
}

func TestStringRepository(t *testing.T) {
	r := setupRepo(t)
	do, err := r.FindPid("xxx")
	if err == nil {
		t.Error("Found missing PID")
	}
	do, err = r.FindPid("dummy:1")
	if err != nil {
		t.Fatal("Did not find dummy:1")
	}
	oi := do.Info()
	t.Log(oi)
	if oi.Pid != "dummy:1" {
		t.Error("Pid does not match")
	}
}

func TestFindVersion(t *testing.T) {
	mdo := MemDigitalObject{}
	mdo.ds = append(mdo.ds,
		MemDatastream{DatastreamInfo: DatastreamInfo{Name: "test", Id: "test.0"}},
		MemDatastream{DatastreamInfo: DatastreamInfo{Name: "test", Id: "test.1"}},
		MemDatastream{DatastreamInfo: DatastreamInfo{Name: "test", Id: "test.2"}},
		MemDatastream{DatastreamInfo: DatastreamInfo{Name: "xxx", Id: "xxx.0"}},
		MemDatastream{DatastreamInfo: DatastreamInfo{Name: "xxx", Id: "xxx.1"}},
		MemDatastream{DatastreamInfo: DatastreamInfo{Name: "yyy", Id: "yyy.0"}},
	)
	table := []struct {
		s              string
		version, index int
	}{
		{"test", 0, 0},
		{"test", -1, 2},
		{"test", 1, 1},
		{"test", 2, 2},
		{"xxx", -1, 4},
		{"xxx", 0, 3},
		{"xxx", 1, 4},
		{"xxx", 2, -1},
		{"yyy", -1, 5},
		{"yyy", 0, 5},
	}

	for _, z := range table {
		result := findVersion(&mdo, z.s, z.version)
		if result != z.index {
			t.Errorf("(%v, %v) expected %v, got %v", z.s, z.version, z.index, result)
		}
	}
}

func TestDecodeVersion(t *testing.T) {
	table := []struct {
		s string
		v int
	}{
		{"dummy.0", 0},
		{"dummy.10", 10},
		{"dummy.5.10", 10},
		{"dummy", 0},
		{"dummy.", 0},
	}

	for i := range table {
		result := decodeVersion(table[i].s)
		if result != table[i].v {
			t.Errorf("decodeVersion(%v) = %v, expected %v", table[i].s, result, table[i].v)
		}
	}
}
