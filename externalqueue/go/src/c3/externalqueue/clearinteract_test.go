package externalqueue

import (
	"c3/osm/workflow"
	"testing"
)

func TestCIPullErrand_01(t *testing.T) {
	//test is no graceful way to shutdown http server as of go 1.4
	//so we start the server here and no need to do this again
	// in other test functions
	go mainServer(12345)
	ci := new(ClearInteractQueue)
	ci.InitTest()
	configureMap := map[string]interface{}{
		"server-address":                                  "localhost:12345",
		"include-customerdata-in-stampcentionitemstarted": false,
		"open-errand-url":                                 "http://localhost",
	}
	ci.Configure(configureMap)
	errand := workflow.NewErrand()
	errand.Id = 1
	TestServerReturnId = 5
	errand.ExternalID = 5
	retVal := ci.PullErrand(errand)
	if retVal == false {
		t.Logf("TestPullErrand_01 failed")
		t.Fail()
	}
}

func TestCIPullErrand_02(t *testing.T) {
	ci := new(ClearInteractQueue)
	ci.InitTest()
	configureMap := map[string]interface{}{
		"server-address":                                  "localhost:12345",
		"include-customerdata-in-stampcentionitemstarted": false,
		"open-errand-url":                                 "http://localhost",
	}
	ci.Configure(configureMap)
	errand := workflow.NewErrand()
	errand.Id = 1
	TestServerReturnId = 2
	errand.ExternalID = 5
	retVal := ci.PullErrand(errand)
	if retVal == true {
		t.Logf("TestPullErrand_02 failed")
		t.Fail()
	}
}

func TestCIRemoveErrand_01(t *testing.T) {
	ci := new(ClearInteractQueue)
	ci.InitTest()
	configureMap := map[string]interface{}{
		"server-address":                                  "localhost:12345",
		"include-customerdata-in-stampcentionitemstarted": false,
		"open-errand-url":                                 "http://localhost",
	}
	ci.Configure(configureMap)
	errand := workflow.NewErrand()
	errand.Id = 1
	TestServerReturnId = 2
	errand.ExternalID = 2
	retVal := ci.RemoveErrand(errand, workflow.Errand_ACTION_ANSWER, false)
	if retVal == false {
		t.Logf("TestRemoveErrand_01 failed")
		t.Fail()
	}
}

func TestCIRemoveErrand_02(t *testing.T) {
	ci := new(ClearInteractQueue)
	ci.InitTest()
	configureMap := map[string]interface{}{
		"server-address":                                  "localhost:12345",
		"include-customerdata-in-stampcentionitemstarted": false,
		"open-errand-url":                                 "http://localhost",
	}
	ci.Configure(configureMap)
	errand := workflow.NewErrand()
	errand.Id = 1
	TestServerReturnId = 3
	errand.ExternalID = 2
	retVal := ci.RemoveErrand(errand, workflow.Errand_ACTION_ANSWER, false)
	if retVal == true {
		t.Logf("TestRemoveErrand_02 failed")
		t.Fail()
	}
}

func TestCIAddErrand_01(t *testing.T) {
	ci := new(ClearInteractQueue)
	ci.InitTest()
	configureMap := map[string]interface{}{
		"server-address":                                  "localhost:12345",
		"include-customerdata-in-stampcentionitemstarted": false,
		"open-errand-url":                                 "http://localhost",
	}
	//set up fake objsrv return
	TestServerReturnId = 3
	StubClearInteractAction = workflow.NewClearInteractAction()
	StubSystemGroup = workflow.NewSystemGroup()
	StubSystemGroup.Id = 1

	ci.Configure(configureMap)
	errand := workflow.NewErrand()
	errand.TargetArea = workflow.NewArea()
	errand.TargetArea.Id = 50
	errand.TargetArea.Name = "Disneyland"
	errand.Service = workflow.NewService()
	errand.Service.Name = "TestCIAddErrand_01"
	errand.Id = 3
	errand.ExternalID = 2
	user := workflow.NewUser()
	user.Id = 1
	user.ExternalID = "100"
	retVal := ci.AddErrand(errand, user, 0)
	if retVal == false {
		t.Logf("TestCIAddErrand_01 failed")
		t.Fail()
	}
}

func TestCIAddErrand_02(t *testing.T) {
	ci := new(ClearInteractQueue)
	ci.InitTest()
	configureMap := map[string]interface{}{
		"server-address":                                  "localhost:12345",
		"include-customerdata-in-stampcentionitemstarted": false,
		"open-errand-url":                                 "http://localhost",
	}
	//set up fake objsrv return
	TestServerReturnId = 3
	StubClearInteractAction = workflow.NewClearInteractAction()
	StubSystemGroup = workflow.NewSystemGroup()
	StubSystemGroup.Id = 1

	ci.Configure(configureMap)
	errand := workflow.NewErrand()
	errand.TargetArea = nil
	errand.Service = workflow.NewService()
	errand.Service.Name = "TestCIAddErrand_01"
	errand.Id = 3
	errand.ExternalID = 2
	user := workflow.NewUser()
	user.Id = 1
	user.ExternalID = "100"
	retVal := ci.AddErrand(errand, user, 0)
	if retVal == true {
		t.Logf("TestCIAddErrand_02 failed")
		t.Fail()
	}
}

func TestCIAddErrand_03(t *testing.T) {
	ci := new(ClearInteractQueue)
	ci.InitTest()
	configureMap := map[string]interface{}{
		"server-address":                                  "localhost:12345",
		"include-customerdata-in-stampcentionitemstarted": false,
		"open-errand-url":                                 "http://localhost",
	}
	//set up fake objsrv return
	TestServerReturnId = 3
	StubClearInteractAction = nil
	StubSystemGroup = workflow.NewSystemGroup()
	StubSystemGroup.Id = 1

	ci.Configure(configureMap)
	errand := workflow.NewErrand()
	errand.TargetArea = workflow.NewArea()
	errand.TargetArea.Id = 50
	errand.TargetArea.Name = "Disneyland"
	errand.Service = workflow.NewService()
	errand.Service.Name = "TestCIAddErrand_03"
	errand.Id = 3
	errand.ExternalID = 2
	user := workflow.NewUser()
	user.Id = 1
	user.ExternalID = "100"
	retVal := ci.AddErrand(errand, user, 0)
	if retVal == false {
		t.Logf("TestCIAddErrand_03 failed")
		t.Fail()
	}
}
func TestCIAddErrand_04(t *testing.T) {
	ci := new(ClearInteractQueue)
	ci.InitTest()
	configureMap := map[string]interface{}{
		"server-address":                                  "localhost:12345",
		"include-customerdata-in-stampcentionitemstarted": false,
		"open-errand-url":                                 "http://localhost",
	}
	//set up fake objsrv return
	TestServerReturnId = 4
	StubClearInteractAction = nil
	StubSystemGroup = workflow.NewSystemGroup()
	StubSystemGroup.Id = 1

	ci.Configure(configureMap)
	errand := workflow.NewErrand()
	errand.TargetArea = workflow.NewArea()
	errand.TargetArea.Id = 50
	errand.TargetArea.Name = "Disneyland"
	errand.Service = workflow.NewService()
	errand.Service.Name = "TestCIAddErrand_04"
	errand.Id = 3
	errand.ExternalID = 2
	retVal := ci.AddErrand(errand, nil, 0)
	if retVal == false {
		t.Logf("TestCIAddErrand_04 failed")
		t.Fail()
	}
}
