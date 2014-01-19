package fedoro

import (
	"bytes"
	"fmt"
	"io"

	//"bitbucket.org/ww/goraptor"
	"github.com/deltamobile/goraptor"
)

type relationship struct {
	predicate, object string
}

type RelationshipIO interface {
	read(r io.Reader) ([]relationship, error)
	write(w io.Writer, rels []relationship) error
}

func ObjectModels(do DigitalObject) []string {
	result := make([]string, 0, 3)

	content, err := do.DsContent("RELS-EXT", -1)
	if err != nil {
		// TODO: handle error. propagate it or empty array?
		return result
	}
	defer content.Close()

	z := simpleRdf{}

	rels, err := z.read(content)
	for _, r := range rels {
		if r.predicate == "info:fedora/fedora-system:def/model#hasModel" {
			result = append(result, r.object)
		}
	}
	return result
}

func addObjectModel(do DigitalObject, model string) error {
	content, err := do.DsContent("RELS-EXT", -1)
	if err != nil {
		return err
	}
	z := simpleRdf{}
	rels, err := z.read(content)
	content.Close()
	for _, r := range rels {
		if r.predicate == "info:fedora/fedora-system:def/model#hasModel" &&
			r.object == model {
			// already has this model
			return nil
		}
	}
	rels = append(rels, relationship{predicate: "info:fedora/fedora-system:def/model#hasModel",
		object: model,
	})
	var b bytes.Buffer
	z.write(&b, rels)
	err = do.ReplaceContent("RELS-EXT", &b)

	return err
}

type simpleRdf struct {
}

func (s simpleRdf) read(r io.Reader) ([]relationship, error) {
	var err error
	result := make([]relationship, 0, 3)
	for {
		var predicate, object string
		n, err := fmt.Fscanf(r, "%q %q\n", &predicate, &object)
		if n != 2 || err != nil {
			break
		}
		result = append(result, relationship{
			predicate: predicate,
			object:    object,
		})
	}
	return result, err
}

func (s simpleRdf) write(w io.Writer, rels []relationship) error {
	var err error
	for _, rel := range rels {
		_, err := fmt.Fprintf(w, "%q %q\n", rel.predicate, rel.object)
		if err != nil {
			break
		}
	}
	return err
}

type RaptorXML struct {
}

func (rxml *RaptorXML) read(r io.Reader) ([]relationship, error) {
	rdf := goraptor.NewParser("guess")
	defer rdf.Free()

	result := make([]relationship, 0, 3)

	ch := rdf.Parse(r, "http://localhost")
	for statement := range ch {
		result = append(result, relationship{
			predicate: statement.Predicate.String(),
			object:    statement.Object.String(),
		})
	}
	return result, nil
}

func (rxml *RaptorXML) write(w io.Writer, rels []relationship) error {
	/* This is ugly. it seems the goraptor package does not support
	 * the creation of new RDF statements...at least not in a way that
	 * I can figure out. So, temporarily we will store the relationships
	 * in a pseudo N3 format: each line is a relationship and is written with
	 * the syntax <predicate> <space> <object> <newline>
	 * very simple and probably broken.
	 */
	// make statement asserting additional object model
	//u := goraptor.Uri("info:fedora/fedora-system:def/model#hasModel")
	//newStatement := goraptor.Statement{
	//	Subject:   &goraptor.Literal{Value: "do.Pid"},
	//	Predicate: &u,
	//	Object:    &goraptor.Literal{Value: model},
	//	Graph:     new(goraptor.Blank),
	//}
	//serializer := goraptor.NewSerializer("xml")
	//if serializer == nil {
	//	log.Println("serializer nil")
	//}
	//defer serializer.Free()

	//serializer.Add(&newStatement)

	//rdf := goraptor.NewParser("guess")
	//defer rdf.Free()

	//ch := rdf.Parse(content, "http://localhost")
	//serializer.AddN(ch)

	// save to datastream
	//str, err := serializer.Serialize(nil, "")
	//if err != nil {
	//return err
	//}

	//do.ReplaceContent("RELS-EXT", strings.NewReader(str))

	return nil
}
