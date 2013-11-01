package fedoro

import (
	"bitbucket.org/ww/goraptor"
)

func ObjectModels(do DigitalObject) []string {
	rdf := goraptor.NewParser("guess")
	defer rdf.Free()

	result := make([]string, 0, 3)

	content, err := do.DsContent("RELS-EXT", -1)
	if err != nil {
		// TODO: handle error. proprogate it or empty array?
		return nil
	}
	defer content.Close()

	ch := rdf.Parse(content, "http://localhost")
	for statement := range ch {
		m := statement.Predicate.String()
		if m == "info:fedora/fedora-system:def/model#hasModel" {
			result = append(result, statement.Object.String())
		}
	}
	return result
}
