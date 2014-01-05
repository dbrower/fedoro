package fedoro

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

func ModifyDatastreamHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	// TODO: sanitize values?
	pid := vars["pid"]
	dsid := vars["dsid"]

	/* This is complicated by a laxity in how Fedora processes these requests.
	 * It allows either (1) Multipart form data, taking the first file section
	 * as the datastream content; or (2) takes the complete request body as the
	 * datastream content, and only pulls the parameters from the url.
	 */
	// actually, fedora doesn't handle multipart PUTs at all

	// use the body as the datastream contents
	v, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		log.Printf("url.ParseQuery: %s", err)
	}
	do, err := MainRepo.FindPid(pid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s Not Found", pid)
		return
	}
	dsinfo := do.DsInfo(dsid, -1)

	// TODO: validate input

	if s := v.Get("controlGroup"); s != "" {
		dsinfo.ControlGroup = aToRune(s)
	}
	if s := v.Get("dsLabel"); s != "" {
		dsinfo.Label = s
	}
	if s := v.Get("versionable"); s != "" {
		dsinfo.Versionable = aToBool(s)
	}
	if s := v.Get("dsState"); s != "" {
		dsinfo.State = aToRune(s)
	}
	if s := v.Get("mimeType"); s != "" {
		dsinfo.Mimetype = s
	}

	err = do.UpdateDatastream(dsinfo)
	if err != nil {
		log.Println(err)
	}

	log.Printf("Modify Datastream, %v\n", do)

	err = do.ReplaceContent(dsid, req.Body)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(200)
}
