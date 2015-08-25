package externalqueue

import (
	"testing"
)

func TestQueueCreateExternalQueue_01(t *testing.T) {
	retObj := CreateExternalQueue("solidus")
	if retObj.Name() != "solidus" {
		t.Logf("TestQueueCreateExternalQueue_01 failed")
		t.Fail()
	}
}

func TestQueueCreateExternalQueue_02(t *testing.T) {
	retObj := CreateExternalQueue("clearinteract")
	if retObj.Name() != "clearinteract" {
		t.Logf("TestQueueCreateExternalQueue_02 failed")
		t.Fail()
	}
}

func TestQueueCreateExternalQueue_03(t *testing.T) {
	retObj := CreateExternalQueue("externalqueue")
	if retObj != nil {
		t.Logf("TestQueueCreateExternalQueue_03 failed")
		t.Fail()
	}
}
