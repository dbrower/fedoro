
package fedoro

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)


func AddDatastreamHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	// TODO: sanitize values?
	pid := vars["pid"]
	dsid := vars["dsid"]

	req.ParseMultipartForm(2000000000) // 2GB limit
	controlGroup := req.Form.Get("controlGroup")
	//dsLocation := req.Form.Get("dsLocation")
	//altId := req.Form.Get("dsLabel")
	dsLabel := req.Form.Get("dsLabel")
	versionable := req.Form.Get("versionable")
	dsState := req.Form.Get("dsState")
	formatUri := req.Form.Get("formatURI")
	//checksumType := req.Form.Get("checksumType")
	//checksum := req.Form.Get("checksum")
	mimetype := req.Form.Get("mimeType")
	//logMessage := req.Form.Get("logMessage")

	do, err := MainRepo.FindPid(pid)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("Not Found"))
		return
	}

	// TODO: validate input here


	do.UpdateDatastream(&DatastreamInfo{
		Name: dsid,
		State: rune(dsState[0]),
		ControlGroup: rune(controlGroup[0]),
		Versionable: aToBool(versionable),
		Label: dsLabel,
		Created:   time.Now(),
		Mimetype: mimetype,
		Format_uri: formatUri,
	})

	log.Println(req.MultipartForm.File)

	//do.ReplaceContent(dsid, r)

	res.WriteHeader(201)
}


func aToBool(s string) bool {
	switch s {
	case "true":
		return true
	}
	return false
}
