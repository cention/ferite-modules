package externalqueue

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"strings"
)

var listen net.Listener
var TestServerReturnId int
var OpenMediaStatus string
var CancelRequestResult bool
var EWT string
var OpenMediaId int
var QueuePosition string
var Status string
var SetOptionBool bool

type myHandler struct {
}

func (o *myHandler) clearInteractResponseMessage(retCode int) string {
	message := "<?xml version='1.0' ?>" +
		"<soapenv:Envelope xmlns:xsi=\"http://www.w3.org/2001/" +
		"XMLSchema-instance\"" +
		" xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\"" +
		" xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\"" +
		" >" +
		"<soapenv:Body>"
	message += toSOAP("Return", fmt.Sprintf("%v", retCode), "")
	message += "</soapenv:Body>"
	message += "</soapenv:Envelope>"
	return message
}

func (o *myHandler) respGetOpenMediaReqStatus(status string) string {
	retEnvelope := new(SolidusResponseEnvelope)
	retEnvelope.Body.RequestResult = new(ByStatusId)
	retEnvelope.Body.RequestResult.OpenMediaRequests.Status = status
	output, err := xml.MarshalIndent(retEnvelope, " ", "    ")
	if err == nil {
		return string(output)
	}
	return ""
}

func (o *myHandler) respOpenMediaCancelReq(result bool) string {
	retEnvelope := new(CancelRequestRespEnvelope)
	retEnvelope.Body.CancelRequestResponse.CancelRequestResult = result
	output, err := xml.MarshalIndent(retEnvelope, " ", "    ")
	if err == nil {
		return string(output)
	} else {
		fmt.Println(err)
	}
	return ""
}

func (o *myHandler) respOpenMediaAddReq(ewt string, openMediaId int,
	pos string, status string) string {
	retEnvelope := new(ErrandResponseEnvelope)
	retEnvelope.Body.AddRequestResponse.EWT = ewt
	retEnvelope.Body.AddRequestResponse.OpenMediaID =
		fmt.Sprintf("%d", openMediaId)
	retEnvelope.Body.AddRequestResponse.QueuePosition = pos
	retEnvelope.Body.AddRequestResponse.RequestStatus = status
	output, err := xml.MarshalIndent(retEnvelope, " ", "    ")
	if err == nil {
		return string(output)
	} else {
		fmt.Println(err)
	}
	return ""
}

func (o *myHandler) respSetOptions(respBool bool) string {
	retEnvelope := new(SetOptionsRespEnvelope)
	retEnvelope.Body.SetOptionsResponse.Value = respBool
	output, err := xml.MarshalIndent(retEnvelope, " ", "    ")
	if err == nil {
		return string(output)
	} else {
		fmt.Println(err)
	}
	return ""
}

func (o *myHandler) respAddRequest(ewt string, openMediaId int,
	pos string, status string) string {
	retEnvelope := new(ErrandResponseEnvelope)
	retEnvelope.Body.AddRequestResponse.EWT = ewt
	retEnvelope.Body.AddRequestResponse.OpenMediaID =
		fmt.Sprintf("%d", openMediaId)
	retEnvelope.Body.AddRequestResponse.QueuePosition = pos
	retEnvelope.Body.AddRequestResponse.RequestStatus = status
	output, err := xml.MarshalIndent(retEnvelope, " ", "    ")
	if err == nil {
		return string(output)
	} else {
		fmt.Println(err)
	}
	return ""
}

func (o *myHandler) processRequest(s string) string {

	envelope := new(CIPullEnvelope)
	clearBody := new(ClearInteractPullItemRequest)
	envelope.Body.Payload = clearBody
	retMsg := ""

	_, err := doUnmarshalString(s, envelope)
	if err != nil {
		retMsg = o.clearInteractResponseMessage(0)
	} else {
		//fmt.Printf("xmlRequest ==>%v\n", envelope.Body.Payload)
		retMsg = o.clearInteractResponseMessage(TestServerReturnId)
	}
	return retMsg
}

func (o *myHandler) endCentionItem(s string) string {

	envelope := new(CIEndEnvelope)
	endBody := new(CIEndItemRequest)
	envelope.Body.Payload = endBody
	retMsg := ""

	//fmt.Printf("s ==> %s\n", s)
	_, err := doUnmarshalString(s, envelope)
	if err != nil {
		retMsg = o.clearInteractResponseMessage(0)
	} else {
		//fmt.Printf("xmlRequest ==>%v\n", envelope.Body.Payload.Reason)
		retMsg = o.clearInteractResponseMessage(TestServerReturnId)
	}
	return retMsg
}

func (o *myHandler) queuePersonalCentionItem(s string) string {

	envelope := new(CIQPersonalEnvelope)
	retMsg := ""

	//fmt.Printf("s ==> %s\n", s)
	_, err := doUnmarshalString(s, envelope)
	if err != nil {
		retMsg = o.clearInteractResponseMessage(0)
	} else {
		//fmt.Printf("xmlRequest ==>%v\n", envelope.Body.Payload.CreationTime)
		retMsg = o.clearInteractResponseMessage(TestServerReturnId)
	}
	return retMsg
}

func (o *myHandler) processGetOpenMediaReqStatus(r *http.Request) string {
	envelope := new(OpenMediaReqEnvelope)
	retMsg := ""
	_, err := doUnmarshal(r.Body, envelope)
	if err == nil {
		retMsg = o.respGetOpenMediaReqStatus(OpenMediaStatus)
	} else {
		fmt.Println(err)
	}
	//fmt.Printf("RetMsg ==>%s\n", retMsg)
	return retMsg
}

func (o *myHandler) processOpenMediaCancelReq(r *http.Request) string {

	envelope := new(CancelRequestEnvelope)
	retMsg := ""
	_, err := doUnmarshal(r.Body, envelope)
	if err == nil {
		retMsg = o.respOpenMediaCancelReq(CancelRequestResult)
	} else {
		fmt.Println(err)
	}
	//fmt.Printf("RetMsg ==>%s\n", retMsg)
	return retMsg
}

func (o *myHandler) processOpenMediaAddReq(r *http.Request) string {

	envelope := new(ErrandRequestEnvelope)
	retMsg := ""
	_, err := doUnmarshal(r.Body, envelope)
	if err == nil {
		retMsg = o.respOpenMediaAddReq(EWT, OpenMediaId, QueuePosition, Status)
	} else {
		fmt.Println(err)
	}
	//fmt.Printf("RetMsg ==>%s\n", retMsg)
	return retMsg
}

func (o *myHandler) processCIRequest(r *http.Request) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	s := buf.String()
	retMsg := ""
	if strings.Contains(s, //r.Header.Get("SOAPAction"),
		"pullCentionItem") {
		retMsg = o.processRequest(s)
	} else if strings.Contains(s, //r.Header.Get("SOAPAction"),
		"endCentionItem") {
		retMsg = o.endCentionItem(s)
	} else if strings.Contains(s, //r.Header.Get("SOAPAction"),
		"queuePersonalCentionItem") {
		retMsg = o.queuePersonalCentionItem(s)
	}
	return retMsg
}

func (o *myHandler) processSetOptions(r *http.Request) string {

	envelope := new(SetOptionsEnvelope)
	retMsg := ""
	_, err := doUnmarshal(r.Body, envelope)
	if err == nil {
		retMsg = o.respSetOptions(SetOptionBool)
	} else {
		fmt.Println(err)
	}
	//fmt.Printf("RetMsg ==>%s\n", retMsg)
	return retMsg
}

func (o *myHandler) processAddRequest(r *http.Request) string {

	envelope := new(ErrandRequestEnvelope)
	retMsg := ""
	_, err := doUnmarshal(r.Body, envelope)
	if err == nil {
		retMsg = o.respAddRequest(EWT, OpenMediaId, QueuePosition, Status)
	} else {
		fmt.Println(err)
	}
	fmt.Printf("RetMsg ==>%s\n", retMsg)
	return retMsg
}

func (o *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	retMsg := ""
	if strings.Contains(r.URL.Path, "quittest") {
		fmt.Println("quitting")
		listen.Close()
		return
	} else if strings.Contains(r.URL.Path, "UqfCention") {
		//clear it stuffs has UqfCention in it's path
		retMsg = o.processCIRequest(r)
	} else if strings.Contains(r.Header.Get("SOAPAction"),
		"SolidusGoTest/GetOpenMediaRequestStatusByID") {
		retMsg = o.processGetOpenMediaReqStatus(r)
	} else if strings.Contains(r.Header.Get("SOAPAction"),
		"CancelRequest") {
		retMsg = o.processOpenMediaCancelReq(r)
	} else if strings.Contains(r.Header.Get("SOAPAction"),
		"AddRequest") {
		retMsg = o.processOpenMediaAddReq(r)
	} else if strings.Contains(r.Header.Get("SOAPAction"),
		"SetOptions") {
		retMsg = o.processSetOptions(r)
	}
	//fmt.Printf("retMsg --> %s\n", retMsg)
	if len(retMsg) > 0 {
		fmt.Fprintf(w, "%s", retMsg)
	} else {
		fmt.Fprintf(w, "hello world")
	}
}

func mainServer(portNo int) {
	//http.HandleFunc("/", handler)
	//http.ListenAndServe(":12345", nil)
	port := ":12345"
	if portNo > 0 {
		port = fmt.Sprintf(":%d", portNo)
	}
	listen, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	s := &http.Server{
		Addr:    port,
		Handler: &myHandler{},
	}
	s.Serve(listen)
}
