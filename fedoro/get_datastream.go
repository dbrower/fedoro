package fedoro

import (
	"encoding/xml"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type DatastreamProfile struct {
	Pid          string    `xml:"pid,attr"`
	Dsid         string    `xml:"dsID,attr"`
	Label        string    `xml:"dsLabel"`
	Version      string    `xml:"dsVersionID"`
	CreateDate   time.Time `xml:"dsCreateDate"`
	State        string    `xml:"dsState"`
	Mimetype     string    `xml:"dsMIME"`
	FormatUri    string    `xml:"dsFormatURI"`
	ControlGroup string    `xml:"dsControlGroup"`
	Size         int       `xml:"dsSize"`
	Versionable  bool      `xml:"dsVersionable"`
	InfoType     string    `xml:"dsInfoType"`
	Location     string    `xml:"dsLocation"`
	LocationType string    `xml:"dsLocationType"`
	ChecksumType string    `xml:"dsChecksumType"`
	Checksum     string    `xml:"dsChecksum"`

	// boilerplate
	Xmlns          string `xml:"xmlns,attr"`
	Xsd            string `xml:"xmlns:xsd,attr"`
	Xsi            string `xml:"xmlns:xsi,attr"`
	SchemaLocation string `xml:"xsi:schemaLocation,attr"`
}

func filloutDatastreamProfile(r Repository, pid string, dsid string) (*DatastreamProfile, error) {
	object, err := r.FindPid(pid)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	dsinfo := object.DsInfo(dsid, -1)

	result := &DatastreamProfile{
		Pid:          pid,
		Dsid:         dsid,
		Label:        dsinfo.Label,
		Version:      dsid + "." + strconv.Itoa(dsinfo.NumVersions-1),
		CreateDate:   dsinfo.Created,
		State:        string(dsinfo.State),
		Mimetype:     dsinfo.Mimetype,
		FormatUri:    dsinfo.Format_uri,
		ControlGroup: string(dsinfo.ControlGroup),
		Size:         dsinfo.Size,
		Versionable:  dsinfo.Versionable,
		InfoType:     "",
		Location:     "",
		LocationType: "",
		ChecksumType: "DISABLED",
		Checksum:     "none",
	}

	result.Xmlns = "http://www.fedora.info/definitions/1/0/management/"
	result.Xsd = "http://www.w3.org/2001/XMLSchema"
	result.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	result.SchemaLocation = "http://www.fedora.info/definitions/1/0/management/ http://www.fedora-commons.org/definitions/1/0/datastreamProfile.xsd"

	return result, nil
}

func GetDatastreamHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	// TODO: sanitize pid?
	pid := vars["pid"]
	dsid := vars["dsid"]
	// TODO: add asOfDateTime option
	// TODO: add validateChecksum

	log.Printf("getDatastreamHandler: pid = %v, dsid = %v\n", pid, dsid)

	result, err := filloutDatastreamProfile(MainRepo, pid, dsid)
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
