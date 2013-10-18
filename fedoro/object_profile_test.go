package fedoro

import (
	"fmt"
	"strings"
	"testing"

	"bitbucket.org/ww/goraptor"
)

func TestObjectModels(t *testing.T) {
	p := goraptor.NewParser("rdfxml")
	defer p.Free()

	source := strings.NewReader(`<rdf:RDF xmlns:ns0="info:fedora/fedora-system:def/model#" xmlns:ns1="info:fedora/fedora-system:def/relations-external#" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
<rdf:Description rdf:about="info:fedora/vecnet:3b5918592">
<ns0:hasModel rdf:resource="info:fedora/afmodel:GenericFile"/>
<ns1:isPartOf rdf:resource="info:fedora/vecnet:3b591858s"/>
</rdf:Description>
</rdf:RDF>`)

	ch := p.Parse(source, "http://localhost")
	for {
		b, ok := <-ch
		if !ok {
			break
		}
		fmt.Println(b)
		if b.Predicate.String() == "info:fedora/fedora-system:def/model#hasModel" {
			fmt.Println("@@@ Is a Model")
		}
	}
}
