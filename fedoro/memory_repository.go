package fedoro

import (
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type MemoryRepository struct {
	items map[string]*MemDigitalObject
}

func NewMemRepo() Repository {
	return MemoryRepository{items: make(map[string]*MemDigitalObject)}
}

func (r MemoryRepository) FindPid(pid string) (DigitalObject, error) {
	for i := range r.items {
		if r.items[i].Pid == pid {
			return r.items[i], nil
		}
	}
	return nil, &memError{"Non-existant PID " + pid}
}

func (r MemoryRepository) NewObject(obj ObjectInfo) (DigitalObject, error) {
	result := MemDigitalObject{ObjectInfo: obj}
	now := time.Now()
	result.CreatedDate = now
	result.ModifiedDate = now
	r.items[obj.Pid] = &result
	return &result, nil
}

func (r MemoryRepository) RemoveObject(pid string) error {
	// This will still work if the pid does not exist
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

func (mdo *MemDigitalObject) Info() *ObjectInfo {
	return &mdo.ObjectInfo
}

func (mdo *MemDigitalObject) DsNames() []string {
	result := make([]string, 0, 5)
	for _, dsInfo := range mdo.ds {
		result = addIfNew(result, dsInfo.Name)
	}
	return result
}

func (mdo *MemDigitalObject) DsInfo(dsid string, version int) *DatastreamInfo {
	i := findVersion(mdo, dsid, version)
	if i == -1 {
		return nil
	}
	return &mdo.ds[i].DatastreamInfo
}

func (mdo *MemDigitalObject) DsContent(dsid string, version int) (io.ReadCloser, error) {
	i := findVersion(mdo, dsid, version)
	if i == -1 {
		return nil, &memError{"Cannot find " + mdo.Pid + "/" + dsid}
	}
	return ioutil.NopCloser(strings.NewReader(mdo.ds[i].content)), nil
}

func (mdo *MemDigitalObject) UpdateInfo(obj *ObjectInfo) error {
	mdo.ObjectInfo = *obj
	return nil
}

func (mdo *MemDigitalObject) UpdateDatastream(dsinfo *DatastreamInfo) error {
	ver := findVersion(mdo, dsinfo.Name, -1)
	if ver > -1 {
		mdo.ds[ver].DatastreamInfo = *dsinfo
	} else {
		mdo.ds = append(mdo.ds, MemDatastream{DatastreamInfo: *dsinfo})
	}
	mdo.ModifiedDate = time.Now()
	return nil
}

func (mdo *MemDigitalObject) ReplaceContent(dsid string, r io.Reader) error {
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
	mdo.ModifiedDate = newDs.Created
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

func findVersion(mdo *MemDigitalObject, dsid string, version int) int {
	var bestIndex int = -1
	var bestVersion int = -1
	for i, ds := range mdo.ds {
		if dsid != ds.Name {
			continue
		}
		v := decodeVersion(ds.Id)
		if v == -1 {
			continue
		} else if v == version {
			return i
		} else if v >= bestVersion {
			bestIndex = i
			bestVersion = v
		}
	}
	if version != -1 {
		return -1
	}
	return bestIndex
}

func decodeVersion(s string) int {
	suffix := strings.LastIndex(s, ".")
	i, err := strconv.Atoi(s[suffix+1:])
	if err != nil {
		i = 0
	}
	return i
}

type memError struct {
	s string
}

func (e *memError) Error() string {
	return e.s
}
