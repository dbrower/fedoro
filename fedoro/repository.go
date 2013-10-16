
package fedoro

import (
    "github.com/dbrower/fedoro/akubra"
    "github.com/dbrower/fedoro/foxml"
)

type Repository struct {
    objectStore, dsStore akubra.Pool
}

func NewRepository(object, ds akubra.Pool) Repository {
    return Repository{objectStore: object, dsStore: ds}
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

