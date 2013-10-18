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
	"net/url"
	"os"
	"strings"

	"github.com/dbrower/fedoro/foxml"
)

var _ = fmt.Println

// The basic type which defines an akubra storage pool.
type Pool struct {
	// The absolute file path to the root directory of this pool
	Root string
	// The format of the sub-directories. Consists of a sequence of hash-marks
	// with optional embedded forward slashes. E.g. "##" or "#/##".
	// Defaults to ""
	Format string
}

type Repository struct {
	objectStore, dsStore Pool
}

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
func GuessFormat(root string) string {
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
	subformat := GuessFormat(root + "/" + fiList[0].Name())
	if subformat != "" {
		return segement + "/" + subformat
	}
	return segement
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
func (p Pool) resolveName(id string) string {
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

func (p Pool) GetReader(id string) (io.ReadCloser, error) {
	path := p.Root + "/" + p.resolveName(id)
	return os.Open(path)
}

// Return a new akubra repository with object info stored
// at objectPath and datastream contents stored at dsPath
func NewRepository(objectPath, dsPath string) Repository {
	obj := Pool{Root: objectPath}
	ds := Pool{Root: dsPath}
	obj.Format = GuessFormat(objectPath)
	ds.Format = GuessFormat(dsPath)
	return Repository{objectStore: obj, dsStore: ds}
}

func (r Repository) FindPid(pid string) (foxml.DigitalObject, error) {
	var d foxml.DigitalObject
	f, err := r.objectStore.GetReader(pid)
	if err != nil {
		return d, err
	}
	defer f.Close()
	d, err = foxml.DecodeFoxml(f)
	return d, err
}
