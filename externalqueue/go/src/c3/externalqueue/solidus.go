package externalqueue

import (
	"c3/osm/workflow"
	"fmt"
	"log"
	"strconv"
)

type Solidus struct {
	solidusHost               string
	systemBaseURL             string
	serverAddress             string
	soapActionURL             string
	QueuedRemoveErrandActions []*SolidusCancelRequestAction
}

func (o *Solidus) Init() {
	wf.InitWorkflow(new(WorkflowClass))
	Feature.InitFeature(new(FeatureClass))
}

func (o *Solidus) InitTest() {
	wf.InitWorkflow(new(WorkflowStub))
	Feature.InitFeature(new(FeatureStub))
}

func (o *Solidus) Name() string {
	return "solidus"
}

func (o *Solidus) SetDefaults(serverAddress string, soapActionURL string) {
	o.serverAddress = serverAddress
	o.soapActionURL = soapActionURL
}

func (o *Solidus) Configure(options map[string]interface{}) {
	if options["server-address"] != nil {
		o.solidusHost = fmt.Sprintf("%v", options["server-address"])
	}
	if options["open-errand-url"] != nil {
		o.systemBaseURL = fmt.Sprint("%v", options["open-errand-url"])
	}
}

func (o *Solidus) createRequest(url string, action string, function string,
	paramMap map[string]interface{}, result string, resultName string,
	retStruct interface{}) interface{} {
	message := "<?xml version='1.0' ?>"
	message += "<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\""
	message += " xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\""
	message += " xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\">"
	message += "<soap:Body>"
	message += fmt.Sprintf("<%s xmlns=\"http://tempuri.org/\">", function)

	if function != "GetOpenMediaRequestStatusByID" &&
		function != "GetServiceGroupStatus" &&
		function != "GetOMAgentsByGroup" {
		message += "<request>"
	}
	for key, value := range paramMap {
		nameSpace := ""
		if function != "GetOpenMediaRequestStatusByID" &&
			function != "GetServiceGroupStatus" &&
			function != "GetOMAgentsByGroup" {
			nameSpace = "http://schemas.datacontract.org/2004/07/Solidus.OpenMedia.Contracts.DataContracts"
		}
		message += toSOAP(key, fmt.Sprintf("%v", value), nameSpace)
	}

	if function != "GetOpenMediaRequestStatusByID" &&
		function != "GetServiceGroupStatus" &&
		function != "GetOMAgentsByGroup" {
		message += "</request>"
	}
	message += fmt.Sprintf("</%s>", function)
	message += "</soap:Body>"
	message += "</soap:Envelope>"
	sendRequest("http://"+url, action, message, retStruct)
	return retStruct
}

func (o *Solidus) cancelRequest(openMediaId int, cancelIfAllocated bool) bool {
	paramMap := make(map[string]interface{})
	paramMap["CancelIfAllocated"] = cancelIfAllocated
	paramMap["DoNotReport"] = false
	paramMap["OpenMediaID"] = openMediaId
	retStruct := new(CancelRequestRespEnvelope)
	o.createRequest(o.serverAddress, o.soapActionURL+"/CancelRequest",
		"CancelRequest",
		paramMap, "boolean", "CancelRequestResult", retStruct)

	return retStruct.Body.CancelRequestResponse.CancelRequestResult
}

func (o *Solidus) RemoveErrand(errand *workflow.Errand,
	wasOpenedByExternalSystem bool) bool {
	if wasOpenedByExternalSystem {
		action := new(SolidusCancelRequestAction)
		action.OpenMediaID = errand.SolidusOpenMediaID
		action.CancelIfAllocated = true
		o.QueuedRemoveErrandActions = append(o.QueuedRemoveErrandActions,
			action)
		errand.SetQueuedInExternal(false)
		errand.SetSolidusOpenMediaID(0)
		wf.SaveObject(errand)
		return true
	}
	return o.remove(errand, true)
}

func (o *Solidus) remove(errand *workflow.Errand, cancelIfAllocated bool) bool {
	log.Printf("Sending cancel request for Workflow.Errand(%d) (%d)",
		errand.Id, errand.SolidusOpenMediaID)
	if o.cancelRequest(errand.SolidusOpenMediaID, cancelIfAllocated) {
		log.Printf("Cancel request returned true for Workflow.Errand"+
			"(%d) (%d)", errand.Id, errand.SolidusOpenMediaID)
		errand.SetQueuedInExternal(false)
		errand.SetSolidusOpenMediaID(0)
		wf.SaveObject(errand)
		return true
	} else {
		log.Printf("Cancel request return false for Workflow.Errand"+
			"(%d) (%d)", errand.Id, errand.SolidusOpenMediaID)
		return false
	}
}

func (o *Solidus) getOpenMediaRequestStatus(mediaId int) *SolidusSoapResponseBody {
	paramMap := make(map[string]interface{})
	paramMap["openMediaID"] = mediaId
	retStruct := new(SolidusResponseEnvelope)
	o.createRequest(o.solidusHost,
		o.soapActionURL+"/GetOpenMediaRequestStatusByID",
		"GetOpenMediaRequestStatusByID",
		paramMap,
		"GetOpenMediaRequestStatusByIDResult",
		"GetOpenMediaRequestStatusByIDResult", retStruct)
	return &retStruct.Body
}

func (o *Solidus) PullErrand(errand *workflow.Errand) bool {
	failed := true
	status := o.getOpenMediaRequestStatus(errand.SolidusOpenMediaID)
	if status != nil {
		if status.RequestResult != nil {
			if status.RequestResult.OpenMediaRequests.Status == "Queued" {
				log.Printf("Status in Solidus for errand %d %d is 'Queued at "+
					" the ServiceGroup'", errand.Id, errand.SolidusOpenMediaID)
				if o.remove(errand, false) {
					log.Printf("Successfully cancelled errand %d in Solidus\n",
						errand.Id)
					failed = false
				} else {
					log.Printf("Unable to cancel errand %d %d in Solidus\n",
						errand.Id, errand.SolidusOpenMediaID)
				}
			} else if status.RequestResult.OpenMediaRequests.Status == "Complete" ||
				status.RequestResult.OpenMediaRequests.Status == "Cancelled" ||
				status.RequestResult.OpenMediaRequests.Status == "Failed" {
				log.Printf("Status in Solidus for errand %d %d is %s",
					errand.Id, errand.SolidusOpenMediaID,
					status.RequestResult.OpenMediaRequests.Status)
				failed = false
				errand.SetQueuedInExternal(false)
				errand.SetSolidusOpenMediaID(0)
				wf.SaveObject(errand)
			}
		} else {
			log.Printf("Got status from Solidus for errand %d %d but media " +
				"had no previous requests registered")
		}
	} else {
		log.Printf("Unable to get status from Solidus for errand %d %d",
			errand.Id, errand.SolidusOpenMediaID)
		failed = false
		errand.SetQueuedInExternal(false)
		errand.SetSolidusOpenMediaID(0)
		wf.SaveObject(errand)
	}
	return !failed
}

func (o *Solidus) setSettingOptions(setting *workflow.SolidusChannelSetting) bool {
	typeOfSession := setting.TypeOfSession
	agentActionOptions := setting.AgentActionOptions
	closeTabOptions := setting.CloseTabOptions
	maxNumberOfSessions := setting.MaxNumberOfSessions
	allowDifferentTypes := setting.AllowDifferentTypes
	resetAllOptions := false
	paramMap := make(map[string]interface{})
	paramMap["AgentActionOptions"] = agentActionOptions
	paramMap["AllowDifferentTypes"] = allowDifferentTypes
	paramMap["CloseTabOptions"] = closeTabOptions
	paramMap["MaxNumberOfSessions"] = maxNumberOfSessions
	paramMap["ResetAllOptions"] = resetAllOptions
	paramMap["TypeOfSession"] = typeOfSession
	retStruct := new(SetOptionsRespEnvelope)
	o.createRequest(o.serverAddress, o.soapActionURL+"/SetOptions",
		"SetOptions",
		paramMap, "boolean", "SetOptionsResult", retStruct)

	return retStruct.Body.SetOptionsResponse.Value
}

func (o *Solidus) addRequest(serviceGroupId int, preferredAgentId int,
	forceToPreferredAgent bool, privateData string,
	ivrInfo map[string]string, typeOfSession int) (*ErrandResponse, string) {
	//No idea how to implement trapping fault like in solidus.feh
	fault := ""
	paramMap := map[string]interface{}{
		"ForceToPreferredAgent": forceToPreferredAgent,
		"PreferredAgentID":      preferredAgentId,
		"PrivateData":           privateData,
		"ServiceGroupID":        serviceGroupId,
		"ServiceGroupName":      "",
		"TenantID":              -1,
	}
	if typeOfSession > 0 {
		paramMap["TypeOfSession"] = typeOfSession
	}
	ivrMsg := ""
	for key, value := range ivrInfo {
		data := fmt.Sprintf("%v", value)
		if len(data) > 20 {
			data = data[:16]
			data += "..."
		}
		ivrMsg += toSOAP(key, data, "")
	}
	if len(ivrMsg) > 0 {
		paramMap["IVRInfo"] = ivrMsg
	}
	retStruct := new(ErrandResponseEnvelope)
	o.createRequest(o.serverAddress, o.soapActionURL+"/AddRequest",
		"AddRequest", paramMap, "AddRequestResult", "AddRequestResult",
		retStruct)
	result := retStruct.Body.AddRequestResponse
	return &result, fault
}

//untested
func (o *Solidus) AddErrand(errand *workflow.Errand, user *workflow.User,
	numType int) bool {
	area := wf.LoadArea(errand.TargetArea.Id)
	systemGroup := wf.QuerySystemGroup_minimalByAreaID(errand.TargetArea.Id)
	setting := wf.QuerySolidusChannelSetting_byChannel(systemGroup.Id,
		errand.Service.Type)
	serviceGroupId := area.SolidusServiceGroupID
	privateData := fmt.Sprintf("%s/solidus/login/-/errand/%d", o.systemBaseURL,
		errand.Id)
	preferredAgentId := 0
	if user != nil {
		preferredAgentId = user.SolidusAgentID
	}
	forceToPreferredAgent := false
	from := "No Sender"
	if errand.Mail != nil && errand.Mail.From != nil {
		if len(errand.Mail.From.Name) > 0 {
			from = fmt.Sprintf("%s (%s)", errand.Mail.From.Name,
				errand.Mail.From.EmailAddress)
		} else {
			from = errand.Mail.From.EmailAddress
		}
	}
	subject := "No Subject"
	if errand.Message != nil {
		if len(errand.Message.Subject) > 0 {
			subject = errand.Message.Subject
		}
	}
	ivrInfo := map[string]string{
		"Errand":  fmt.Sprintf("%d", errand.Id),
		"From":    from,
		"Subject": subject,
	}
	typeOfSession := 0
	if setting != nil {
		typeOfSession = setting.TypeOfSession
	}

	sortReplies := false
	if user != nil {
		globalContext := fmt.Sprintf("system-group:%d", user.SystemGroup)
		Feature.SetGlobalContext(globalContext)
		sortReplies = Feature.Bool("sort-replies-first.default-yes")
	}
	if area.SolidusPrioritizedServiceGroupID > 0 &&
		(user != nil && (sortReplies && errand.Reply) ||
			numType == workflow.Errand_ACTION_EXTERNAL_EXPERT_ANSWER ||
			numType == workflow.Errand_ACTION_AGENT_FORWARD) {
		serviceGroupId = area.SolidusPrioritizedServiceGroupID
	}
	if setting != nil {
		if o.setSettingOptions(setting) == false {
			return false
		}
	}

	result, fault := o.addRequest(serviceGroupId, preferredAgentId,
		forceToPreferredAgent, privateData, ivrInfo, typeOfSession)

	// at the moment there's no code to return fault from addRequest()
	if fault == "InvalidDataException" || fault == "NoLicenseFault" ||
		fault == "RouterDisconnectedFault" || fault == "InvalidUserFault" {
		return true
	}

	if result != nil && len(result.OpenMediaID) > 0 {
		errand.SetQueuedInExternal(true)
		openMediaId, convErr := strconv.Atoi(result.OpenMediaID)
		if convErr == nil {
			errand.SetSolidusOpenMediaID(openMediaId)
		}
		wf.SaveObject(errand)
		return true
	}
	return false
}
