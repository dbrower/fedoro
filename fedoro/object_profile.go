package fedoro

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/ww/goraptor"
	"github.com/gorilla/mux"

	"github.com/dbrower/fedoro/foxml"
)

type objectProfile struct {
	Pid              string    `xml:"pid,attr"`
	Label            string    `xml:"objLabel"`
	OwnerId          string    `xml:"objOwnerId"`
	Models           []string  `xml:"objModels>models"`
	CreateDate       time.Time `xml:"objCreateDate"`
	LastModDate      time.Time `xml:"objLastModDate"`
	DissIndexViewURL string    `xml:",omitempty"`
	DissItemViewURL  string    `xml:",omitempty"`
	State            string    `xml:"objState"`
	// boilerplate
	Xmlns          string `xml:"xmlns,attr"`
	Xsd            string `xml:"xmlns:xsd,attr"`
	Xsi            string `xml:"xmlns:xsi,attr"`
	SchemaLocation string `xml:"xsi:schemaLocation,attr"`
}

func ObjectModels(do foxml.DigitalObject) []string {
	rdf := goraptor.NewParser("guess")
	defer rdf.Free()

	result := make([]string, 0, 3)

	ds := foxml.GetDatastream(do, "RELS-EXT")
	content := r.GetMostRecentContent(ds)
	// assume most recent version is last
	dsv := ds.Versions[len(ds.Versions)-1]

	ch := rdf.Parse(do.GetContent("RELS-EXT", -1), "http://localhost")
	for statement := range ch {
		m := statement.Predicate.String()
		if m == "info:fedora/fedora-system:def/model#hasModel" {
			result = append(result, statement.Object.String())
		}
	}
	return result
}

func ObjectProfile(r Repository, pid string) (*objectProfile, error) {
	object, err := r.FindPid(pid)
	if err != nil {
		// TODO
		fmt.Println(err)
		return nil, err
	}

	result := &objectProfile{
		Pid:         pid,
		Label:       object.Label,
		State:       object.State,
		OwnerId:     object.OwnerId,
		CreateDate:  object.CreatedDate,
		LastModDate: object.ModifiedDate,
	}

	result.Xmlns = "http://www.fedora.info/definitions/1/0/access/"
	result.Xsd = "http://www.w3.org/2001/XMLSchema"
	result.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	result.SchemaLocation = "http://www.fedora.info/definitions/1/0/access/ http://www.fedora-commons.org/definitions/1/0/listDatastreams.xsd"

	return result, nil
}

func ObjectProfileHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	// TODO: sanitize pid?
	pid := vars["pid"]

	result, err := ObjectProfile(MainRepo, pid)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	req.ParseForm()
	format := req.Form.Get("format")

	if format == "xml" {
		res.Header().Add("Content-Type", "text/xml")
		e := xml.NewEncoder(res)
		e.Encode(result)
	}
}
