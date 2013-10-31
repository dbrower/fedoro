package fedoro

import (
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type MemoryRepository struct {
	items map[string]MemDigitalObject
}

func NewMemRepo() Repository {
	return MemoryRepository{items: make(map[string]MemDigitalObject)}
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

func (r MemoryRepository) NewObject(obj ObjectInfo) (DigitalObject, error) {
	result := MemDigitalObject{
		ObjectInfo: obj,
	}
	r.items[obj.Pid] = result
	return result, nil
}

func (r MemoryRepository) RemoveObject(pid string) error {
	delete(r.items, pid)
	return nil
}

type MemDigitalObject struct {
	ObjectInfo
	ds []MemDatastream
}

type MemDatastream struct {
	DatastreamInfo
	content string
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
	return ioutil.NopCloser(strings.NewReader(mdo.ds[i].content)), nil
}

func (mdo MemDigitalObject) UpdateInfo(obj *ObjectInfo) error {
	mdo.ObjectInfo = *obj
	return nil
}

func (mdo MemDigitalObject) UpdateDatastream(dsinfo *DatastreamInfo) error {
	ver := findVersion(mdo, dsinfo.Name, -1)
	if ver == -1 {
		// TODO: return error
		return nil
	}
	mdo.ds[ver].DatastreamInfo = *dsinfo
	return nil
}

func (mdo MemDigitalObject) ReplaceContent(dsid string, r io.Reader) error {
	var newDs MemDatastream

	ver := findVersion(mdo, dsid, -1)
	if ver == -1 {
		newDs.Name = dsid
		newDs.State = 'A'
		newDs.ControlGroup = 'X'
		newDs.Versionable = true
	} else {
		newDs = mdo.ds[ver]
	}

	newDs.NumVersions += 1
	newDs.Id = newDs.Name + "." + strconv.Itoa(newDs.NumVersions-1)
	newDs.Created = time.Now()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	newDs.content = string(data)
	newDs.Size = len(data)

	mdo.ds = append(mdo.ds, newDs)
	return nil
}

func addIfNew(list []string, text string) []string {
	for _, s := range list {
		if s == text {
			return list
		}
	}
	return append(list, text)
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
			} else if best == -1 || v >= bestVersion {
				best = i
				bestVersion = v
			}
		}
	}
	return best
}

func decodeVersion(s string) int {
	suffix := strings.Scanl(s, ".")
	return strconv.Atoi(suffix)
}
