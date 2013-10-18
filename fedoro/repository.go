package fedoro

import (
	//	"time"
	//    "github.com/dbrower/fedoro/akubra"
	"github.com/dbrower/fedoro/foxml"
)

var (
	MainRepo Repository
)

type Repository interface {
	FindPid(string) (foxml.DigitalObject, error)
}
