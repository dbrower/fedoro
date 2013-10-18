
package fedoro

import (
	"github.com/dbrower/fedoro/foxml"
)

type StringRepository struct {
	items []foxml.DigitalObject
}

func (r StringRepository) FindPid(pid string) (foxml.DigitalObject, error) {
	for _, d := range r.items {
		if d.Pid == pid {
			return d, nil
		}
	}
	// should be an error instead of nil
	return foxml.DigitalObject{}, nil
}

