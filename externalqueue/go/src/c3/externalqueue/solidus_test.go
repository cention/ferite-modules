package externalqueue

import (
	"c3/osm/workflow"
	"testing"
)

func TestSolidusPullErrand_01(t *testing.T) {
	go mainServer(12346)
	sol := new(Solidus)
	sol.InitTest()
	configureMap := map[string]interface{}{
		"server-address":  "localhost:12346",
		"open-errand-url": "http://localhost",
	}
	sol.Configure(configureMap)
	sol.SetDefaults("localhost:12346", "SolidusGoTest")
	errand := workflow.NewErrand()
	errand.Id = 1
	OpenMediaStatus = "Queued"
	CancelRequestResult = true
	errand.ExternalID = 5
	errand.SolidusOpenMediaID = 6
	retVal := sol.PullErrand(errand)
	if retVal == false {
		t.Logf("TestSolidusPullErrand_01 failed")
		t.Fail()
	}
}

func TestSolidusPullErrand_02(t *testing.T) {
	go mainServer(12346)
	sol := new(Solidus)
	sol.InitTest()
	configureMap := map[string]interface{}{
		"server-address":  "localhost:12346",
		"open-errand-url": "http://localhost",
	}
	sol.Configure(configureMap)
	sol.SetDefaults("localhost:12346", "SolidusGoTest")
	errand := workflow.NewErrand()
	errand.Id = 1
	OpenMediaStatus = "Queued"
	CancelRequestResult = false
	errand.ExternalID = 5
	errand.SolidusOpenMediaID = 6
	retVal := sol.PullErrand(errand)
	if retVal == true {
		t.Logf("TestSolidusPullErrand_02 failed")
		t.Fail()
	}
}

func TestSolidusPullErrand_03(t *testing.T) {
	sol := new(Solidus)
	sol.InitTest()
	configureMap := map[string]interface{}{
		"server-address":  "localhost:12346",
		"open-errand-url": "http://localhost",
	}
	sol.Configure(configureMap)
	sol.SetDefaults("localhost:12346", "SolidusGoTest")
	errand := workflow.NewErrand()
	errand.Id = 1
	OpenMediaStatus = "Cancelled"
	CancelRequestResult = true
	errand.ExternalID = 5
	errand.SolidusOpenMediaID = 6
	retVal := sol.PullErrand(errand)
	if retVal == false {
		t.Logf("TestSolidusPullErrand_03 failed")
		t.Fail()
	}
}

func TestSolidusRemErrand_01(t *testing.T) {
	sol := new(Solidus)
	sol.InitTest()
	configureMap := map[string]interface{}{
		"server-address":  "localhost:12346",
		"open-errand-url": "http://localhost",
	}
	sol.Configure(configureMap)
	sol.SetDefaults("localhost:12346", "SolidusGoTest")
	errand := workflow.NewErrand()
	errand.Id = 2
	OpenMediaStatus = "Queued"
	CancelRequestResult = true
	errand.ExternalID = 5
	errand.SolidusOpenMediaID = 6
	retVal := sol.RemoveErrand(errand, true)
	if retVal == false || len(sol.QueuedRemoveErrandActions) != 1 {
		t.Logf("TestSolidusRemErrand_01 failed")
		t.Fail()
	}
}

func TestSolidusRemErrand_02(t *testing.T) {
	sol := new(Solidus)
	sol.InitTest()
	configureMap := map[string]interface{}{
		"server-address":  "localhost:12346",
		"open-errand-url": "http://localhost",
	}
	sol.Configure(configureMap)
	sol.SetDefaults("localhost:12346", "SolidusGoTest")
	errand := workflow.NewErrand()
	errand.Id = 2
	OpenMediaStatus = "Queued"
	CancelRequestResult = false
	errand.ExternalID = 5
	errand.SolidusOpenMediaID = 7
	retVal := sol.RemoveErrand(errand, false)
	if retVal == true || len(sol.QueuedRemoveErrandActions) != 0 {
		t.Logf("TestSolidusRemErrand_02 failed")
		t.Fail()
	}
}

func TestSolidusAddErrand_01(t *testing.T) {
	sol := new(Solidus)
	sol.InitTest()
	configureMap := map[string]interface{}{
		"server-address":  "localhost:12346",
		"open-errand-url": "http://localhost",
	}
	SetOptionBool = false
	StubChannelSetting = workflow.NewSolidusChannelSetting()
	StubChannelSetting.Id = 4
	StubChannelSetting.TypeOfSession = 2
	StubChannelSetting.AgentActionOptions = 3
	StubChannelSetting.CloseTabOptions = 4
	StubChannelSetting.MaxNumberOfSessions = 1
	StubChannelSetting.AllowDifferentTypes = 0
	StubSystemGroup = workflow.NewSystemGroup()
	StubSystemGroup.Id = 7
	sol.Configure(configureMap)
	sol.SetDefaults("localhost:12346", "SolidusGoTest")
	errand := workflow.NewErrand()
	errand.Id = 2
	errand.ExternalID = 5
	errand.SolidusOpenMediaID = 7
	errand.TargetArea = workflow.NewArea()
	errand.TargetArea.Id = 2
	StubArea = errand.TargetArea
	errand.Service = workflow.NewService()
	errand.Service.Type = workflow.Errand_SERVICE_EMAIL
	user := workflow.NewUser()
	retVal := sol.AddErrand(errand, user, 0)
	user.SolidusAgentID = 8
	if retVal == true {
		t.Logf("TestSolidusAddErrand_01 failed")
		t.Fail()
	}
}

func TestSolidusAddErrand_02(t *testing.T) {
	sol := new(Solidus)
	sol.InitTest()
	configureMap := map[string]interface{}{
		"server-address":  "localhost:12346",
		"open-errand-url": "http://localhost",
	}
	SetOptionBool = true
	OpenMediaId = 90
	Status = "ok"
	StubChannelSetting = workflow.NewSolidusChannelSetting()
	StubChannelSetting.Id = 4
	StubChannelSetting.TypeOfSession = 2
	StubChannelSetting.AgentActionOptions = 3
	StubChannelSetting.CloseTabOptions = 4
	StubChannelSetting.MaxNumberOfSessions = 1
	StubChannelSetting.AllowDifferentTypes = 0
	StubSystemGroup = workflow.NewSystemGroup()
	StubSystemGroup.Id = 7
	sol.Configure(configureMap)
	sol.SetDefaults("localhost:12346", "SolidusGoTest")
	errand := workflow.NewErrand()
	errand.Id = 2
	errand.ExternalID = 5
	errand.SolidusOpenMediaID = 7
	errand.TargetArea = workflow.NewArea()
	errand.TargetArea.Id = 2
	StubArea = errand.TargetArea
	StubArea.SolidusServiceGroupID = 700
	errand.Service = workflow.NewService()
	errand.Service.Type = workflow.Errand_SERVICE_EMAIL
	errand.Mail = workflow.NewMail()
	errand.Mail.From = workflow.NewMailOrigin()
	errand.Mail.From.Name = "jake"
	errand.Mail.Subject = "TestSolidusAddErrand_02"
	user := workflow.NewUser()
	user.SolidusAgentID = 1007
	retVal := sol.AddErrand(errand, user, 0)
	user.SolidusAgentID = 8
	if retVal == false {
		t.Logf("TestSolidusAddErrand_02 failed")
		t.Fail()
	}
}
