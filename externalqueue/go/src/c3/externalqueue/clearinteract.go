package externalqueue

import (
	"c3/osm/workflow"
	"fmt"
	"time"
)

type ClearInteractQueue struct {
	clearInteractHost                            string
	includeCustomerDataInStampCentionItemStarted bool
	systemBaseURL                                string
}

func (o *ClearInteractQueue) Init() {
	wf.InitWorkflow(new(WorkflowClass))
	Feature.InitFeature(new(FeatureClass))
}

func (o *ClearInteractQueue) InitTest() {
	wf.InitWorkflow(new(WorkflowStub))
	Feature.InitFeature(new(FeatureStub))
}

func (o *ClearInteractQueue) invoke(url string, function string,
	paramMap map[string]string, retStruct interface{}) interface{} {
	message := "<?xml version='1.0' ?>" +
		"<soapenv:Envelope xmlns:xsi=\"http://www.w3.org/2001/" +
		"XMLSchema-instance\"" +
		" xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\"" +
		" xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\"" +
		" xmlns:ser=\"services.soap.uqf.clearit.se\">" +
		"<soapenv:Body>"
	message += fmt.Sprintf("<ser:%s soapenv:encodingStyle=\""+
		"http://schemas.xmlsoap.org/soap/encoding/\">", function)
	for key, value := range paramMap {
		message += toSOAP(key, value, "")
	}
	message += fmt.Sprintf("</ser:%s>", function)
	message += "</soapenv:Body>"
	message += "</soapenv:Envelope>"

	return sendRequest(url, " ", message, retStruct)
}

func (o *ClearInteractQueue) pullCentionItem(itemId int,
	retStruct interface{}) {
	invokeMap := map[string]string{
		"itemId": fmt.Sprintf("%d", itemId),
		"reason": "Pulled",
	}
	o.invoke("http://"+o.clearInteractHost+"/UqfCention",
		"pullCentionItem", invokeMap, retStruct)
}

func (o *ClearInteractQueue) queueCentionItem(itemId int,
	taskIdentifier string, url string, errandId string,
	area string, channel string, creationTime string,
	retStruct interface{}) int {
	if itemId <= 0 {
		itemId = -1
	}
	invokeMap := map[string]string{
		"itemId":         fmt.Sprintf("%d", itemId),
		"taskIdentifier": taskIdentifier,
		"url":            url,
		"errandId":       errandId,
		"area":           area,
		"channel":        channel,
		"creationTime":   creationTime,
	}
	o.invoke("http://"+o.clearInteractHost+"/UqfCention",
		"queuePersonalCentionItem", invokeMap, retStruct)
	return 0
}

func (o *ClearInteractQueue) queuePersonalCentionItem(itemId int,
	taskIdentifier string, agentIdentifier string, url string, errandId string,
	area string, channel string, creationTime string,
	retStruct interface{}) {
	invokeMap := map[string]string{
		"itemId":          fmt.Sprintf("%d", itemId),
		"taskIdentifier":  taskIdentifier,
		"agentIdentifier": agentIdentifier,
		"url":             url,
		"errandId":        errandId,
		"area":            area,
		"channel":         channel,
		"creationTime":    creationTime,
	}
	o.invoke("http://"+o.clearInteractHost+"/UqfCention",
		"queuePersonalCentionItem", invokeMap, retStruct)
}

func (o *ClearInteractQueue) endCentionItem(itemId int, reason string,
	retStruct interface{}) {
	invokeMap := map[string]string{
		"itemId": fmt.Sprintf("%d", itemId),
		"reason": reason,
	}
	o.invoke("http://"+o.clearInteractHost+"/UqfCention",
		"endCentionItem", invokeMap, retStruct)
	return

}

func (o *ClearInteractQueue) Name() string {
	return "clearinteract"
}

func (o *ClearInteractQueue) Configure(options map[string]interface{}) {
	if options["server-address"] != nil {
		o.clearInteractHost = fmt.Sprintf("%v", options["server-address"])
	}
	if options["include-customerdata-in-stampcentionitemstarted"] != nil {
		o.includeCustomerDataInStampCentionItemStarted =
			options["include-customerdata-in-stampcentionitemstarted"].(bool)
	}
	if options["open-errand-url"] != nil {
		o.systemBaseURL = fmt.Sprintf("%v", options["open-errand-url"])
	}
}

func (o *ClearInteractQueue) finaliseActions() {
}

func (o *ClearInteractQueue) AddErrand(errand *workflow.Errand,
	user *workflow.User, errandType int) bool {
	returnItemId := 0
	itemId := -1
	if errand.ExternalID > 0 {
		itemId = errand.ExternalID
	}
	taskIdentifier := ""
	area := ""
	if errand.TargetArea != nil {
		taskIdentifier = fmt.Sprintf("%d", errand.TargetArea.Id)
		area = fmt.Sprintf("%s/%s", errand.TargetArea.Name,
			errand.TargetArea.Id)
	}
	agentIdentifier := ""
	if user != nil {
		agentIdentifier = user.ExternalID
	}
	url := fmt.Sprintf("%s/clearinteract/login/-/errand/%d/item/"+
		"[at_external_id]/agent/[at_agent_id]/browser/[at_browser_type]",
		o.systemBaseURL, errand.Id)
	errandId := fmt.Sprintf("errand/%d", errand.Id)
	channel := errand.Service.Name
	creationTime := time.Unix(errand.TimestampArrive, 0).Format("06-01-02:15:04")
	retStruct := new(ClearRetEnv)
	action := wf.QueryClearInteractAction_byErrand(errand.Id)
	if action == nil {
		if user != nil {
			o.queuePersonalCentionItem(itemId, taskIdentifier,
				agentIdentifier, url, errandId, area,
				channel, creationTime, retStruct)
			if retStruct != nil {
				returnItemId = retStruct.Body.Return
			}
		} else {
			o.queueCentionItem(itemId, taskIdentifier, url,
				errandId, area, channel, creationTime, retStruct)
			if retStruct != nil {
				returnItemId = retStruct.Body.Return
			}
		}
	}

	if returnItemId > 0 {
		errand.SetQueuedInExternal(true)
		errand.SetExternalID(returnItemId)
		wf.SaveObject(errand)
		return true
	} else if returnItemId == -5 || returnItemId == -8 {
		action = new(workflow.ClearInteractAction)
	}

	if action != nil && errand.TargetArea != nil {
		action.SetSystemgroup(wf.QuerySystemGroup_minimalByAreaID(errand.TargetArea.Id).Id)
		action.SetType(workflow.ClearInteractAction_ACTION_QUEUE_ITEM)
		action.SetErrand(errand)
		action.SetItemId(itemId)
		action.SetTaskIdentifier(taskIdentifier)
		action.SetAgentIdentifier(agentIdentifier)
		action.SetUrl(url)
		action.SetErrandId(errandId)
		action.SetArea(area)
		action.SetChannel(channel)
		action.SetCreationTime(creationTime)
		wf.SaveObject(action)

		errand.SetQueuedInExternal(true)
		wf.SaveObject(errand)
		return true
	}
	return false
}

func (o *ClearInteractQueue) RemoveErrand(errand *workflow.Errand,
	actionType int, wasOpenedByExternalSystem bool) bool {
	if errand == nil {
		return false
	}

	action := wf.QueryClearInteractAction_byErrand(errand.Id)
	if action != nil {
		action.Delete()
		errand.SetQueuedInExternal(false)
		wf.SaveObject(errand)
		return true
	} else {
		reason := ""
		switch actionType {
		case workflow.Errand_ACTION_ANSWER:
			reason = "Answered"
		case workflow.Errand_ACTION_DELETE:
			reason = "Deleted"
		case workflow.Errand_ACTION_SAVE:
			reason = "Save"
		case workflow.Errand_ACTION_AGENT_FORWARD:
			reason = "Forwarded To Agent"
		case workflow.Errand_ACTION_AREA_FORWARD:
			reason = "Forwarded To Area"
		case workflow.Errand_ACTION_FOLDER_FORWARD:
			reason = "Forwarded To Folder"
		case workflow.Errand_ACTION_INBOX_RETURN:
			reason = "Return To Inbox"
		default:
			reason = "Pulled"
		}
		retStruct := new(ClearRetEnv)
		o.endCentionItem(errand.ExternalID, reason, retStruct)
		itemId := retStruct.Body.Return
		//log.Printf("itemId %v\n", itemId)
		if itemId == errand.ExternalID {
			errand.SetQueuedInExternal(false)
			wf.SaveObject(errand)
			return true
		}
	}
	return false
}

func (o *ClearInteractQueue) PullErrand(errand *workflow.Errand) bool {
	action := wf.QueryClearInteractAction_byErrand(errand.Id)
	if action != nil {
		action.Delete()
		errand.SetQueuedInExternal(false)
		wf.SaveObject(errand)
		return true
	} else {
		retStruct := new(ClearRetEnv)
		o.pullCentionItem(errand.ExternalID, retStruct)
		itemId := retStruct.Body.Return
		//-1 = errand not in CI queue
		if itemId == errand.ExternalID || itemId == -1 {
			errand.SetQueuedInExternal(false)
			wf.SaveObject(errand)
			return true
		}
	}
	return false
}
