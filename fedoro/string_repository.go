
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
	return foxml.DigitalObject{}, error("Pid not found")
}

func (r StringRepository) Doit() int {
	return 42;
}

func ab(r Repository) {
	r.Doit()
}