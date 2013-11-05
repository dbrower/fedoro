package fedoro

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Handles object ingest
func IngestHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	// TODO: sanitize values?
	pid := vars["pid"]

	label := req.FormValue("label")
	namespace := req.FormValue("namespace")
	ownerId := req.FormValue("ownerId")
	//format := req.FormValue("format")
	//encoding := req.FormValue("encoding")
	//logMessage := req.FormValue("logMessage")
	//ignoreMime := req.FormValue("ignoreMime")

	// TODO: validate input here

	obj := ObjectInfo{}
	obj.Pid = pid
	if len(pid) == 0 || pid == "new" {
		if len(namespace) == 0 {
			namespace = "dummy"
		}
		obj.Pid = namespace + ":" + "1"
	}
	if len(label) > 0 {
		obj.Label = label
	}
	if len(ownerId) > 0 {
		obj.OwnerId = ownerId
	}
	obj.State = "A"

	_, err := MainRepo.NewObject(obj)
	if err != nil {
		log.Println(err)
		res.WriteHeader(403)
		res.Write([]byte(err.Error()))
		return
	}

	res.Header().Add("Location", "http://localhost:8080/objects/"+obj.Pid)

	res.WriteHeader(201)
	res.Write([]byte(obj.Pid))
}
