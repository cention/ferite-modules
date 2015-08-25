package externalqueue

import (
	"encoding/xml"
)

//Solidus related XMLs
type OpenMediaReqEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    OpenMediaReqSoapBody
}
type OpenMediaReqSoapBody struct {
	XMLName      xml.Name `xml:"Body"`
	MediaRequest Request  `xml:"GetOpenMediaRequestStatusByID>openMediaID"`
}
type Request struct {
	Name string `xml:",chardata"`
}

type CancelRequestEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    CancelRequestSoapBody
}
type CancelRequestSoapBody struct {
	XMLName       xml.Name     `xml:"Body"`
	CancelRequest CancelStruct `xml:"CancelRequest>request"`
}
type CancelStruct struct {
	CancelIfAllocated string `xml: "CancelIfAllocated"`
	DoNotReport       string `xml: "DoNotReport"`
	OpenMediaID       string `xml: "OpenMediaID"`
}

type ErrandRequestEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    ErrandRequestSoapBody
}
type ErrandRequestSoapBody struct {
	XMLName    xml.Name      `xml:"Body"`
	AddRequest ErrandRequest `xml:"AddRequest>request"`
}
type ErrandRequest struct {
	ForceToPreferredAgent string  `xml: "ForceToPreferredAgent"`
	PreferredAgentID      string  `xml: "PreferredAgentID"`
	PrivateData           string  `xml: "PrivateData"`
	ServiceGroupID        string  `xml: "ServiceGroupID"`
	ServiceGroupName      string  `xml: "ServiceGroupName"`
	TenantID              string  `xml: "TenantID"`
	TypeOfSession         string  `xml: "TypeOfSession"`
	IVRInfo               IvrInfo `xml: "IVRInfo"`
}
type IvrInfo struct {
	XMLName        xml.Name         `xml:"IVRInfo"`
	IVRInformation []IvrInformation `xml: "IVRInformation"`
}
type IvrInformation struct {
	Data  string `xml: "Data"`
	Label string `xml: "Label"`
}

type SetOptionsEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    SetOptionsSoapBody
}
type SetOptionsSoapBody struct {
	XMLName    xml.Name      `xml:"Body"`
	SetOptions OptionsStruct `xml:"SetOptions>request"`
}
type OptionsStruct struct {
	AgentActionOptions  string `xml: "AgentActionOptions"`
	AllowDifferentTypes string `xml: "AllowDifferentTypes"`
	CloseTabOptions     string `xml: "CloseTabOptions"`
	MaxNumberOfSessions string `xml: "MaxNumberOfSessions"`
	ResetAllOptions     string `xml: "ResetAllOptions"`
	TypeOfSession       string `xml: "TypeOfSession"`
}

type SolidusResponseEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    SolidusSoapResponseBody
}
type SolidusSoapResponseBody struct {
	XMLName       xml.Name    `xml:"Body"`
	RequestResult *ByStatusId `xml:"GetOpenMediaRequestStatusByIDResponse>GetOpenMediaRequestStatusByIDResult"`
}
type ByStatusId struct {
	Name              string        `xml:",chardata"`
	TimeStamp         StatusTime    `xml: "TimeStamp"`
	OpenMediaRequests MediaRequests `xml: "OpenMediaRequests"`
}
type StatusTime struct {
	Name string `xml:",chardata"`
}
type MediaRequests struct {
	AgentID   int    `xml:"AgentID"`
	TenantID  int    `xml:"TenantID"`
	Status    string `xml:"Status"`
	TimeStamp string `xml:"TimeStamp"`
	ID        int    `xml:"ID"`
	LogonID   int    `xml:"LogonID"`
}

type CancelRequestRespEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    CancelRequestRespBody
}
type CancelRequestRespBody struct {
	XMLName               xml.Name          `xml:"Body"`
	CancelRequestResponse CancelRequestResp `xml:"CancelRequestResponse>CancelRequestResult"`
}
type CancelRequestResp struct {
	CancelRequestResult bool
}

type SetOptionsRespEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    SetOptionsRespBody
}
type SetOptionsRespBody struct {
	XMLName            xml.Name    `xml:"Body"`
	SetOptionsResponse OptionsResp `xml:"SetOptionsResponse>SetOptionsResult"`
}
type OptionsResp struct {
	Value bool `xml:"boolean"`
}

type ErrandResponseEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    ErrandResponseSoapBody
}
type ErrandResponseSoapBody struct {
	XMLName            xml.Name       `xml:"Body"`
	AddRequestResponse ErrandResponse `xml:"AddRequestResponse>AddRequestResult"`
}
type ErrandResponse struct {
	RequestStatus string `xml:"RequestStatus"`
	EWT           string `xml:"EWT"`
	OpenMediaID   string `xml:"OpenMediaID"`
	QueuePosition string `xml:"QueuePosition"`
}

type SolidusCancelRequestAction struct {
	OpenMediaID       int
	CancelIfAllocated bool
}

//ClearInteract related XMLs
type ClearRetEnv struct {
	XMLName xml.Name                      `xml:"Envelope"`
	Body    ClearInteractPullItemResponse `xml:"Body"`
}

type ClearInteractPullItemResponse struct {
	XMLName xml.Name `xml:"Body"`
	Return  int      `xml:"Return"`
}

type CIPullEnvelope struct {
	XMLName xml.Name       `xml:"Envelope"`
	Body    CIPullSoapBody `xml:"Body"`
}

type CIPullSoapBody struct {
	Payload *ClearInteractPullItemRequest
}

type ClearInteractPullItemRequest struct {
	XMLName xml.Name `xml:"pullCentionItem"`
	ItemId  int      `xml:"itemId"`
	Reason  string   `xml:"reason"`
}

type CIEndEnvelope struct {
	XMLName xml.Name      `xml:"Envelope"`
	Body    CIEndSoapBody `xml:"Body"`
}

type CIEndSoapBody struct {
	Payload *CIEndItemRequest
}

type CIEndItemRequest struct {
	XMLName xml.Name `xml:"endCentionItem"`
	ItemId  int      `xml:"itemId"`
	Reason  string   `xml:"reason"`
}

type CIQPersonalEnvelope struct {
	XMLName xml.Name            `xml:"Envelope"`
	Body    CIQPersonalSoapBody `xml:"Body"`
}

type CIQPersonalSoapBody struct {
	Payload *CIQPersonalItemRequest
}

type CIQPersonalItemRequest struct {
	XMLName         xml.Name `xml:"queuePersonalCentionItem"`
	ItemId          int      `xml:"itemId"`
	TaskIdentifier  string   `xml:"taskIdentifier"`
	AgentIdentifier string   `xml:"agentIdentifier"`
	Url             string   `xml:"url"`
	ErandId         string   `xml:"erandId"`
	Area            string   `xml:"area"`
	Channel         string   `xml:"channel"`
	CreationTime    string   `xml:"creationTime"`
}
