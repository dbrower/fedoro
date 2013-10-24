package fedoro

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type dsType struct {
	Dsid     string `xml:"dsid,attr"`
	Label    string `xml:"label,attr"`
	MimeType string `xml:"mimeType,attr"`
}

type objectDatastreams struct {
	Pid          string   `xml:"pid,attr"`
	AsOfDateTime string   `xml:"asOfDateTime,attr,omitempty"`
	BaseUrl      string   `xml:"baseURL,attr"`
	Datastream   []dsType `xml:"datastream"`
	// boilerplate
	Xmlns          string `xml:"xmlns,attr"`
	Xsd            string `xml:"xmlns:xsd,attr"`
	Xsi            string `xml:"xmlns:xsi,attr"`
	SchemaLocation string `xml:"xsi:schemaLocation,attr"`
}

func ListDatastreams(r Repository, pid string) (*objectDatastreams, error) {
	object, err := r.FindPid(pid)
	if err != nil {
		// TODO
		fmt.Println(err)
		return nil, err
	}

	result := &objectDatastreams{Pid: pid}
	dsNames := object.DsNames()
	for _, name := range dsNames {
		if name == "AUDIT" {
			continue
		}
		info := object.DsInfo(name, -1)
		a := dsType{Dsid: info.Name, Label: info.Label, MimeType: info.Mimetype}
		result.Datastream = append(result.Datastream, a)
	}

	result.Xmlns = "http://www.fedora.info/definitions/1/0/access/"
	result.Xsd = "http://www.w3.org/2001/XMLSchema"
	result.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	result.SchemaLocation = "http://www.fedora.info/definitions/1/0/access/ http://www.fedora-commons.org/definitions/1/0/listDatastreams.xsd"

	return result, nil
}

func ListDatastreamsHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	// TODO: sanitize pid?
	pid := vars["pid"]

	result, err := ListDatastreams(MainRepo, pid)
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
