package fedoro

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func DatastreamDisseminationHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	// TODO: sanitize pid?
	pid := vars["pid"]
	dsid := vars["dsid"]
	// TODO: add asOfDateTime option

	log.Printf("DatastreamDisseminationHandler: pid = %v, dsid = %v\n", pid, dsid)

	object, err := MainRepo.FindPid(pid)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("Cannot find pid"))
		return
	}
	ds := object.DsInfo(dsid, -1)
	if ds == nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("Cannot find ds info for " + dsid))
		return
	}
	data, err := object.DsContent(dsid, -1)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("cannot find content"))
		return
	}
	mime := ds.Mimetype
	if mime == "" {
		mime = ""
	}
	res.Header().Add("Content-Type", mime)

	io.Copy(res, data)
	data.Close()
}
