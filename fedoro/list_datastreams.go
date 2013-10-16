
package fedoro

import (
    "encoding/xml"
    "os"
)

type dsType struct {
    Dsid string `xml:"dsid,attr"`
    Label string `xml:"label,attr"`
    MimeType string `xml:"mimeType,attr"`
}

type objectDatastreams struct {
    Pid string `xml:"pid,attr"`
    AsOfDateTime string `xml:"asOfDateTime,attr,omitempty"`
    BaseUrl string `xml:"baseURL,attr"`
    Datastream []dsType `xml:"datastream"`
}


func ListDatastreams(repo Repository, pid string) {
    var ods objectDatastreams

    ods.Pid = pid
    var s = [3]dsType {
        {Dsid: pid + ":DC", Label: "Dublin Core", MimeType: "text/xml"},
        {Dsid: pid + ":descMetadata", Label: "Descriptive Metadata", MimeType: "text/xml"},
        {Dsid: pid + ":content", Label: "Content", MimeType: "text/xml"},
    }
    ods.Datastream = s[:]

    e := xml.NewEncoder(os.Stdout)
    e.Encode(ods)
}
