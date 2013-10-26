package fedoro

import (
	"io"
	"io/ioutil"
	"strings"
)

type MemoryRepository struct {
	items map[string]MemDigitalObject
}

type MemDigitalObject struct {
	ObjectInfo
	ds []MemDatastream
}

type MemDatastream struct {
	DatastreamInfo
	content string
}

func NewMemRepo() Repository {
	return MemoryRepository{}
}

func (r MemoryRepository) FindPid(pid string) (DigitalObject, error) {
	for _, d := range r.items {
		if d.Pid == pid {
			return d, nil
		}
	}
	// should be an error instead of nil
	return MemDigitalObject{}, nil
}

func (mdo MemDigitalObject) Info() *ObjectInfo {
	return &mdo.ObjectInfo
}

func (mdo MemDigitalObject) DsNames() []string {
	result := make([]string, 0, 5)
	for _, dsInfo := range mdo.ds {
		result = addIfNew(result, dsInfo.Name)
	}
	return result
}

func addIfNew(list []string, text string) []string {
	for _, s := range list {
		if s == text {
			return list
		}
	}
	return append(list, text)
}

func (mdo MemDigitalObject) DsInfo(dsid string, version int) *DatastreamInfo {
	i := findVersion(mdo, dsid, version)
	if i == -1 {
		return nil
	}
	return &mdo.ds[i].DatastreamInfo
}

func (mdo MemDigitalObject) DsContent(dsid string, version int) (io.ReadCloser, error) {
	i := findVersion(mdo, dsid, version)
	if i == -1 {
		// TODO: return error
		return nil, nil
	}
	return ioutil.NopCloser(strings.Reader(mdo.ds[i].content)), nil
}

func findVersion(mdo MemDigitalObject, dsid string, version int) int {
	var best int = -1
	var bestVersion int
	var targetVersion int = version
	for i, ds := range mdo.ds {
		if dsid == ds.Name {
			v := decodeVersion(ds.Id)
			if v == targetVersion {
				return i
			} else if best == nil || v >= bestVersion {
				best = i
				bestVersion = v
			}
		}
	}
	return best
}

func decodeVersion(s string) int {
	suffix := strings.Scanl(s, ".")
	return strings.atoi(suffix)
}
