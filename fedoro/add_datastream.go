package fedoro

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type multiFileReader struct {
	form *multipart.Form
	f    *os.File
	pos  int
}

func (m *multiFileReader) Read(p []byte) {

}

func (m *multiFileReader) RemoveAll() {
	if m.f != nil {
		m.f.Close()
	}
	if m.form != nil {
		m.form.RemoveAll()
	}
}

func AddDatastreamHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	// TODO: sanitize values?
	pid := vars["pid"]
	dsid := vars["dsid"]

	/* This is complicated by a laxity in how Fedora processes these requests.
	 * It allows either (1) Multipart form data, taking the first file section
	 * as the datastream content; or (2) takes the complete request body as the
	 * datastream content, and only pulls the parameters from the url.
	 */
	var f io.Reader
	var v url.Values

	_, err := req.MultipartReader()
	switch {
	case err == nil:
		// a multipart form!
		// we will skip this for now

		// form, err := mp.ReadForm(2 << 20)
		// if err != nil {
		// }
		// defer form.RemoveAll()
		// /* fedora takes the first file part. We will take the first part
		//    in the map, but since maps are unordered it may not be the first
		//    in the stream. (But then, why pass more than one file to begin with? */
		// var fh []*multipart.FileHeader
		// for _, fh = range form.File {
		// 	break
		// }
		// if fh == nil {
		// 	// no files!
		// }
		// file, err := fh.Open()
		// defer file.Close()
		// f = file
		// v = url.Values(form.Value)

	case err == http.ErrNotMultipart:
		// use the body as the datastream contents
		f = req.Body
		v, err = url.ParseQuery(req.URL.RawQuery)
		if err != nil {
			log.Printf("url.ParseQuery: %s", err)
		}

	case err != nil:
		// something strange happened
		log.Println(err)
	}

	// Invariants:
	// v points to the correct hash for the parameters
	// f is the source of the data

	controlGroup := v.Get("controlGroup")
	dsLabel := v.Get("dsLabel")
	versionable := v.Get("versionable")
	dsState := v.Get("dsState")
	formatUri := v.Get("formatURI")
	mimetype := v.Get("mimeType")
	//dsLocation := v.Get("dsLocation")
	//altId := v.Get("dsLabel")
	//logMessage := v.Get("logMessage")
	//checksumType := v.Get("checksumType")
	//checksum := v.Get("checksum")

	do, err := MainRepo.FindPid(pid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s Not Found", pid)
		return
	}

	// TODO: validate input

	err = do.UpdateDatastream(&DatastreamInfo{
		Name:         dsid,
		State:        aToRune(dsState),
		ControlGroup: aToRune(controlGroup),
		Versionable:  aToBool(versionable),
		Label:        dsLabel,
		Created:      time.Now(),
		Mimetype:     mimetype,
		Format_uri:   formatUri,
	})
	if err != nil {
		log.Println(err)
	}

	log.Printf("Add Datastream, %v\n", do)

	err = do.ReplaceContent(dsid, f)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(201)
}

func aToRune(s string) rune {
	if len(s) > 0 {
		return rune(s[0])
	}
	return ' '
}

func aToBool(s string) bool {
	v, err := strconv.ParseBool(s)
	if err != nil {
		v = false
	}
	return v
}
