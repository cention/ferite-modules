package externalqueue

import (
	"bytes"
	"c3/osm/workflow"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var wf WorkflowObject
var Feature FeatureObject

type ExternalQueueInterface interface {
	Init()
	Name() string
	Configure(options map[string]interface{})
	AddErrand(errand *workflow.Errand,
		user *workflow.User, errandType int) bool
	PullErrand(errand *workflow.Errand) bool
}

type Envelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    SoapBody `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
}

type SoapBody struct {
	Payload interface{}
}

func CreateExternalQueue(queueType string) ExternalQueueInterface {
	var retObj ExternalQueueInterface
	if queueType == "solidus" {
		retObj = new(Solidus)
	} else if queueType == "clearinteract" {
		retObj = new(ClearInteractQueue)
	}
	if retObj != nil {
		retObj.Init()
	}
	return retObj
}

func toSOAP(name string, value string, ns string) string {
	if len(ns) > 0 {
		ns = fmt.Sprintf(" xmlns=\"%s\"", ns)
	}
	if len(value) > 0 {
		return fmt.Sprintf("<%s%s>%s</%s>", name, ns, value, name)
	}
	return fmt.Sprintf("<%s%s />", name, ns)
}

func sendRequest(url string, action string, content string,
	retStruct interface{}) interface{} {
	req, err := http.NewRequest("POST", url, strings.NewReader(content))
	if err != nil {
		log.Println(err)
		return nil
	}
	req.Method = "POST"
	req.Header.Add("ContentType", "text/xml; charset=utf-8")
	req.Header.Add("SOAPAction", action)
	client := http.Client{}

	//log.Printf("sending content: %s", content)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("http Do error: %v", err)
		return nil
	}

	doUnmarshal(resp.Body, retStruct)
	resp.Body.Close()
	return retStruct
}

func doUnmarshal(r io.ReadCloser, v interface{}) (interface{}, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	s := buf.String()

	return doUnmarshalString(s, v)
}

func doUnmarshalString(s string, v interface{}) (interface{}, error) {

	//log.Printf("s ==>%s\n", s)
	err := xml.Unmarshal([]byte(s), v)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	return v, err
}

//Unable to get this to work as there is no graceful way to shutdown
//http.Server and no way of killing of a goroutine manually
func sendQuitTest(url string, content string) {
	req, err := http.NewRequest("POST", url, strings.NewReader(content))
	if err != nil {
		log.Println(err)
		return
	}
	req.Method = "POST"
	req.Header.Add("ContentType", "text/xml; charset=utf-8")
	req.Header.Add("SOAPAction", " ")
	client := http.Client{}

	client.Do(req)
}
