package fedoro

import (
	"encoding/xml"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

func ObjectProfile(r Repository, pid string) (*objectProfile, error) {
	object, err := r.FindPid(pid)
	if err != nil {
		// TODO
		log.Println(err)
		return nil, err
	}

	info := object.Info()

	result := &objectProfile{
		Pid:         pid,
		Label:       info.Label,
		State:       info.State,
		OwnerId:     info.OwnerId,
		CreateDate:  info.CreatedDate,
		LastModDate: info.ModifiedDate,
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
	// TODO: add asOfDateTime

	log.Printf("ObjectProfileHandler: pid = %v", pid)

	result, err := ObjectProfile(MainRepo, pid)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("Not Found"))
		return
	}

	req.ParseForm()
	format := req.Form.Get("format")

	if format == "xml" {
		res.Header().Add("Content-Type", "text/xml")
		e := xml.NewEncoder(res)
		e.Encode(result)
	} else {
		res.Write([]byte("add ?format=xml"))
	}
}
