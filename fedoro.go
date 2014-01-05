package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	_ "github.com/dbrower/fedoro/akubra"
	"github.com/dbrower/fedoro/fedoro"
)

const describeText = `<html><head><title>Describe Repository</title></head>
<body>
<h1>Fedoro</h1>
</body>
`

var dt = template.Must(template.New("describe").Parse(describeText))

func DescribeHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	dt.Execute(res, nil)
}

/*
 * Fedora API endpoints
 *
@describeRepositoryLite = [{label: "describeRepositoryLite",
                           command: "describe"}
]
@api_read_only_no_object = [
  {label: "findObjectsLite",
   command: "search?query=pid%7E*&pid=true"},
  {label: "findObjects",
   command: "objects?query=pid%7E*&pid=true"}
]
@api_read_only = [
  {label: "getObjectHistoryLite", command: "getObjectHistory/$PID"},
  {label: "getObjectProfileLite", command: "get/$PID"},
  {label: "listMethodsLite",      command: "listMethods/$PID?xml=true"},
  {label: "getDatastreamDisseminationLite", command: "get/$PID/DC"},
  {label: "listDatastreamsLite",  command: "listDatastreams/$PID?xml=true"},
  {label: "getObjectHistory",     command: "objects/$PID/versions?format=xml"},
  {label: "getObjectProfile",     command: "objects/$PID?format=xml"},
  {label: "listMethods",          command: "objects/$PID/methods?format=xml"},
  {label: "getDatastreamDissemination", command: "objects/$PID/datastreams/DC/content"},
  {label: "listDatastreams",      command: "objects/$PID/datastreams?format=xml"},
  {label: "export",               command: "objects/$PID/export"},
  {label: "getDatastream",        command: "objects/$PID/datastreams/DC?format=xml"},
  {label: "getDatastreamHistory", command: "objects/$PID/datastreams/DC/history?format=xml"},
  {label: "getDatastreams",       command: "objects/$PID/datastreams?profiles=true"},
  {label: "getObjectXML",         command: "objects/$PID/objectXML"},
  {label: "getRelationships",     command: "objects/$PID/relationships"},
  {label: "validate",             command: "objects/$PID/validate"}
]
@api_modify_no_object = [
  {label: "getNextPIDLite",       command: "management/getNextPID?xml=true"},
  {label: "getNextPID",           command: "objects/nextPID", method: :post}
]
@api_modify = [
  {label: "addDatastream",
   command: "objects/$PID/datastreams/test?controlGroup=M&dsLabel=test&checksumType=SHA-256&mimeType=text/plain",
   method: :post,
   post_data: "some-content"},
  {label: "addRelationship",
   command: "objects/$PID/relationships/new?predicate=http%3a%2f%2fwww.example.org%2frels%2fname&object=dublin%20core&isLiteral=true",
   method: :post},
  {label: "modifyDatastream",     command: "objects/$PID/datastreams/test?dsLabel=test-changed", method: :put},
  {label: "modifyObject",         command: "objects/$PID?label=test--new%20label", method: :put},
  {label: "setDatastreamVersionable",
   command: "objects/$PID/datastreams/test?versionable=true",
   method: :put},
  {label: "setDatastreamState",   command: "objects/$PID/datastreams/test?dsState=I", method: :put}
]
@api_oai = [{label: "oaiIdentify", command: "oai?verb=Identify"}]
@api_purge = [
  {label: "purgeDatastream",      command: "objects/$PID/datastreams/test", method: :delete},
  {label: "purgeObject",          command: "objects/$PID", method: :delete}
]
@api_softdelete = [
  # Fedora bug? cannot set the datastream state to D...think the xacml policy is confusing
  # the current ds state with the new ds state
  #should_only_work_admin "setDatastreamState D #{prefix}" "objects/#{prefix}:#{noid}/datastreams/test?dsState=D" put
  #should_not_work "get D datastream #{prefix}" "objects/#{prefix}:#{noid}/datastreams/test/content"
  # put the object in a D state and try to access
  {label: "set object to D state",command: "objects/$PID?state=D", method: :put}
]
@api_readsoftdelete = [
  {label: "get D object",         command: "objects/$PID?format=xml"},
  {label: "get ds from D object", command: "objects/$PID/datastreams/test"},
]

# unimplemented API tests:
#
# resumeFindObjectsLite # not implemented
# uploadFileLite # not implemented
# describeRepository # This entry has not been implemented by Fedora
# resumeFindObjects # not implemented
# getDissemination # not implemented
# compareDatastreamChecksum # not implemented
# ingest # not implemented
# purgeRelationship # not implemented
# upload # not implemented
*/

/* write APIs
  "management/getNextPID"
  POST "objects/nextPID"
  POST "objects/{pid}/datastreams/{dsid}
  POST objects/{pid}/relationships/new
  PUT  objects/{pid}/datastreams/{dsid}
  PUT  objects/{pid}
  PUT  objects/{pid}/datastreams/{dsid}

  GET  oai?verb=Identify

  DELETE objects/{pid}/datastreams/{dsid}
  DELETE objects/{pid}

# unimplemented API tests:
#
# resumeFindObjectsLite # not implemented
# uploadFileLite # not implemented
# describeRepository # This entry has not been implemented by Fedora
# resumeFindObjects # not implemented
# getDissemination # not implemented
# compareDatastreamChecksum # not implemented
# ingest # not implemented
# purgeRelationship # not implemented
# upload # not implemented

*/

type handlerEntry struct {
	path string
	f    http.HandlerFunc
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	fmt.Println("Starting Fedoro")

	repo := fedoro.NewMemRepo()
	do, _ := repo.NewObject(fedoro.ObjectInfo{Pid: "dummy:1234"})
	do.ReplaceContent("test", strings.NewReader("this is test content!!!"))
	fedoro.MainRepo = repo

	routeGet := []handlerEntry{
		{"/describe", fedoro.DescribeHandler},
		{"/objects/{pid}", fedoro.ObjectProfileHandler},
		{"/objects/{pid}/datastreams", fedoro.ListDatastreamsHandler},
		{"/objects/{pid}/datastreams/{dsid}", fedoro.GetDatastreamHandler},
		{"/objects/{pid}/datastreams/{dsid}/content", fedoro.DatastreamDisseminationHandler},
	}

	//r.HandleFunc("/objects/{pid}/datastreams/{dsid}/history")
	//r.HandleFunc("/objects/{pid}/export")
	//r.HandleFunc("/objects/{pid}/methods")
	//r.HandleFunc("/objects/{pid}/objectXML")
	//r.HandleFunc("/objects/{pid}/relationships")
	//r.HandleFunc("/objects/{pid}/validate")
	//r.HandleFunc("/objects/{pid}/versions")
	//r.HandleFunc("/search")
	//r.HandleFunc("/get/{pid}")
	//r.HandleFunc("/get/{pid}/{dsid}")
	//r.HandleFunc("/getObjectHistory/{pid}")
	//r.HandleFunc("/listDatastreams/{pid}")
	//r.HandleFunc("/listMethods/{pid}")

	r := mux.NewRouter()
	installHandlers(r, routeGet, "GET")
	installHandlers(r, routeGet, "HEAD")

	routePost := []handlerEntry{
		{"/objects/{pid}/datastreams/{dsid}", fedoro.AddDatastreamHandler},
		{"/objects/{pid}/relationships/new", notImplementedHandler},
		{"/objects/{pid}", fedoro.IngestHandler},
		{"/objects", fedoro.IngestHandler},
	}

	installHandlers(r, routePost, "POST")

	routePut := []handlerEntry{
		{"/objects/{pid}/datastreams/{dsid}", fedoro.ModifyDatastreamHandler},
		{"/objects/{pid}", notImplementedHandler},
		{"/objects/{pid}/datastreams/{dsid}", notImplementedHandler},
	}

	installHandlers(r, routePut, "PUT")

	r.NotFoundHandler = handlerWrapper(notFoundHandler)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func installHandlers(r *mux.Router, routes []handlerEntry, method string) {
	for _, entry := range routes {
		r.HandleFunc(entry.path, handlerWrapper(entry.f)).Methods(method)
	}
}

func handlerWrapper(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		f(w, r)
	}
}

func notImplementedHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/html")
	res.WriteHeader(http.StatusNotImplemented)
	res.Write([]byte("Not Implemented"))
}

func notFoundHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/html")
	res.WriteHeader(http.StatusNotFound)
	res.Write([]byte("Not Found"))
}
