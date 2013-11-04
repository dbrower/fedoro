
package fedoro

import (
	"encoding/xml"
	"net/http"
)

type fedoraRepository struct {
	RepositoryName              string    `xml:"repositoryName"`
	RepositoryBaseUrl            string    `xml:"repositoryBaseURL"`
	RepositoryVersion          string    `xml:"repositoryVersion"`
	// boilerplate
	//Xmlns          string `xml:"xmlns,attr"`
	//Xsd            string `xml:"xmlns:xsd,attr"`
	//Xsi            string `xml:"xmlns:xsi,attr"`
	//SchemaLocation string `xml:"xsi:schemaLocation,attr"`
}


func DescribeHandler(res http.ResponseWriter, req *http.Request) {

    result := fedoraRepository{
        RepositoryName: "Fedora",
        RepositoryBaseUrl: "localhost:8080",
        RepositoryVersion: "3.7",
    }

	inXml := req.FormValue("xml")

	if aToBool(inXml) {
		res.Header().Add("Content-Type", "text/xml")
		e := xml.NewEncoder(res)
		e.Encode(result)
	} else {
		res.Write([]byte("add ?format=xml"))
	}
}
