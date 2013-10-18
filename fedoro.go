
package main

import (
    "fmt"
    "net/http"
    "html/template"
    "log"

    "github.com/gorilla/mux"

	"github.com/dbrower/fedoro/fedoro"
	"github.com/dbrower/fedoro/akubra"
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

/* Read Only APIs

    /describe
    /get/:pid
    /get/:pid/:ds
    /getObjectHistory/:pid
    /listDatastreams/:pid
    /listMethods/:pid
    /objects
    /objects/:pid
    /objects/:pid/datastreams
    /objects/:pid/datastreams/:ds
    /objects/:pid/datastreams/:ds/content
    /objects/:pid/datastreams/:ds/history
    /objects/:pid/export
    /objects/:pid/methods
    /objects/:pid/objectXML
    /objects/:pid/relationships
    /objects/:pid/validate
    /objects/:pid/versions
    /search
*/



func main() {
    fmt.Println("Starting Fedoro")

	fedoro.MainRepo = akubra.NewRepository("fedoro/test-repo", "fedoro/test-repo")


    r := mux.NewRouter()
    r.HandleFunc("/describe", DescribeHandler).Methods("GET", "HEAD")
    r.HandleFunc("/objects/{pid}/datastreams", fedoro.ListDatastreamsHandler).Methods("GET", "HEAD")
	err := http.ListenAndServe(":8080", r)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
