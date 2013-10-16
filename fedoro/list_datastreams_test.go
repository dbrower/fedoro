
package fedoro

import (
    "testing"
)

func TestListDatastreams(t *testing.T) {
    var r Repository
    
    ListDatastreams(r, "12345")
}
