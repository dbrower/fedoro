/*
Provide an adaptor for accessing an akubra filesystem datastore.

It does not implement the generic Akubra BlobStore interface.

There are some configuration parameters:

    - The location of the `objectStore` directory
    - The directory layout of the object store directory (defaults to `##/`
    - The location of the `datastreamStore`
    - The directory layout of the datastream store directory.

*/
package akubra

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/dbrower/fedoro/fedoro"
	"github.com/dbrower/fedoro/foxml"
)

var _ = fmt.Println

func isAllLowerHex(name string) bool {
	for _, r := range name {
		switch {
		case r >= '0' && r <= '9', r >= 'a' && r <= 'f':
			continue
		}
		return false
	}
	return true
}

// Try to guess the format string for the blob store at root.
// returns the empty string, "", if the format cannot be determined.
//
// Expects to find a directory containing directories all of which have names
// that are the same number of lowercase hex characters. It then recurses to
// build the format string.
func guessFormat(root string) string {
	dir, err := os.Open(root)
	if err != nil {
		return ""
	}
	defer dir.Close()
	fiList, err := dir.Readdir(0)
	if err != nil {
		return ""
	}
	if len(fiList) == 0 {
		return ""
	}
	length := len(fiList[0].Name())
	for _, fileinfo := range fiList {
		if !fileinfo.IsDir() {
			return ""
		}
		fname := fileinfo.Name()
		if length != len(fname) {
			return ""
		}
		if !isAllLowerHex(fname) {
			return ""
		}
	}

	segement := strings.Repeat("#", length)
	subformat := guessFormat(root + "/" + fiList[0].Name())
	if subformat != "" {
		return segement + "/" + subformat
	}
	return segement
}

// The basic type which defines an akubra storage pool.
type pool struct {
	// The absolute file path to the root directory of this pool
	Root string
	// The format of the sub-directories. Consists of a sequence of hash-marks
	// with optional embedded forward slashes. E.g. "##" or "#/##".
	// Defaults to ""
	Format string
}

// Given an id, return the relative path name of the object
// from the pool's root.
// Associated to a pool because it needs the format string
//
// The algorithm is:
//
// 1) prepend 'info:fedora/' to the id
// 2) hash the string using md5
// 3) URL encode the id
// 4) concatenate the format string and the URL-encoded id
// 5) replace '#' symbols with letters from the hex representation of the hash as needed
//
// For example the id 'fedora-system:FedoraObject-3.0' with a pool having
// format string '##' resolves to 'e5/info%3Afedora%2Ffedora-system%3AFedoraObject-3.0'
func (p pool) resolveName(id string) string {
	s1 := "info:fedora/" + id
	h := md5.New()
	io.WriteString(h, s1)
	hashchars := hex.EncodeToString(h.Sum(nil))
	var i int
	f := func(r rune) rune {
		if r != '#' {
			return r
		}
		c := hashchars[i]
		i++
		return rune(c)
	}
	prefix := strings.Map(f, p.Format)
	return prefix + "/" + url.QueryEscape(s1)
}

func (p pool) getReader(id string) (io.ReadCloser, error) {
	path := p.Root + "/" + p.resolveName(id)
	return os.Open(path)
}

func getDatastream(do *foxml.DigitalObject, name string) *foxml.Datastream {
	for i, ds := range do.Ds {
		if ds.Id == name {
			return &do.Ds[i]
		}
	}
	return nil
}

func getDatastreamVersion(ds *foxml.Datastream, version int) *foxml.DatastreamVersion {
	// This is a nieve way of doing this...might break if datastreams
	// are not stored from oldest to newest in the foxml
	if version < -1 {
		return nil
	}
	if version == -1 {
		version = len(ds.Versions) - 1
	}
	return &ds.Versions[version]
}

func getDatastreamAndVersion(do *foxml.DigitalObject, name string, version int) (*foxml.Datastream, *foxml.DatastreamVersion) {
	var dsv *foxml.DatastreamVersion
	ds := getDatastream(do, name)
	if ds != nil {
		dsv = getDatastreamVersion(ds, version)
	}
	return ds, dsv
}

type Repository struct {
	objectStore, dsStore pool
}

// Return a new akubra repository with object info stored
// at objectPath and datastream contents stored at dsPath
func NewRepository(objectPath, dsPath string) Repository {
	obj := pool{Root: objectPath}
	ds := pool{Root: dsPath}
	obj.Format = guessFormat(objectPath)
	ds.Format = guessFormat(dsPath)
	return Repository{objectStore: obj, dsStore: ds}
}

func (r Repository) FindPid(pid string) (*FoxmlObject, error) {
	f, err := r.objectStore.getReader(pid)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	d, err := foxml.DecodeFoxml(f)
	if err != nil {
		return nil, err
	}
	return &FoxmlObject{
		repository: &r,
		obj:        &d,
	}, nil
}

type FoxmlObject struct {
	repository *Repository
	obj        *foxml.DigitalObject
}

func (fo FoxmlObject) Info() *fedoro.ObjectInfo {
	// TODO: change this to cache
	d := fo.obj
	return &fedoro.ObjectInfo{
		Pid:          d.Pid,
		Version:      d.Version,
		State:        d.State,
		Label:        d.Label,
		OwnerId:      d.OwnerId,
		CreatedDate:  d.CreatedDate,
		ModifiedDate: d.ModifiedDate,
	}
}

func (fo FoxmlObject) DsNames() []string {
	d := fo.obj
	result := make([]string, 0, 3)
	for _, ds := range d.Ds {
		result = append(result, ds.Id)
	}
	return result
}

func (fo FoxmlObject) DsInfo(dsid string, version int) *fedoro.DatastreamInfo {
	d := fo.obj
	ds := getDatastream(d, dsid)
	if ds == nil {
		return nil
	}
	dsv := getDatastreamVersion(ds, version)
	if dsv == nil {
		return nil
	}

	return &fedoro.DatastreamInfo{
		Name:         ds.Id,
		State:        ds.State,
		ControlGroup: ds.ControlGroup,
		Versionable:  ds.Versionable,
		NumVersions:  len(ds.Versions),
		Id:           dsv.Id,
		Label:        dsv.Label,
		Created:      dsv.Created,
		Mimetype:     dsv.Mimetype,
		Format_uri:   dsv.Format_uri,
		Size:         dsv.Size,
	}
}

func (fo FoxmlObject) DsContent(dsid string, version int) (io.ReadCloser, error) {
	d := fo.obj
	ds := getDatastream(d, dsid)
	if ds == nil {
		// TODO: return error!
		return nil, nil
	}
	dsv := getDatastreamVersion(ds, version)
	if dsv == nil {
		// TODO: return error!
		return nil, nil
	}

	if dsv.XmlContent.Contents != "" {
		return ioutil.NopCloser(strings.NewReader(dsv.XmlContent.Contents)), nil
	}

	switch ds.ControlGroup {
	case 'M':
		// Need to fetch contents from disk
		if dsv.ContentLocation.Kind == "INTERNAL" {
			return fo.repository.dsStore.getReader(dsv.ContentLocation.Ref)
		}
		// TODO: return error
		return nil, nil

	default:
		// TODO: this needs to return an error
		return nil, nil
	}
}
