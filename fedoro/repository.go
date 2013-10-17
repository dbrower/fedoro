
package fedoro

import (
//	"time"
//    "github.com/dbrower/fedoro/akubra"
    "github.com/dbrower/fedoro/foxml"
)

type Repository interface {
    FindPid(string) (foxml.DigitalObject, error)
}

