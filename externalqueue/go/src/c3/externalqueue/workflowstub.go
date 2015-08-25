package externalqueue

import (
	"c3/osm/workflow"
)

//WorkflowStub implements the WorkflowInterface
//should only be used for testing
type WorkflowStub struct {
}

var StubClearInteractAction *workflow.ClearInteractAction
var StubSystemGroup *workflow.SystemGroup
var StubArea *workflow.Area
var StubChannelSetting *workflow.SolidusChannelSetting

func InitTest() {
	StubClearInteractAction = nil
	StubSystemGroup = nil
}

func (o *WorkflowStub) SaveObject(theObject WorkflowSaveInterface) error {
	return nil
}

func (o *WorkflowStub) QueryClearInteractAction_byErrand(id int) *workflow.ClearInteractAction {
	return StubClearInteractAction
}

func (o *WorkflowStub) QuerySystemGroup_minimalByAreaID(id int) *workflow.SystemGroup {
	return StubSystemGroup
}

func (o *WorkflowStub) LoadArea(id int) *workflow.Area {
	return StubArea
}

func (o *WorkflowStub) QuerySolidusChannelSetting_byChannel(sGroup int,
	sType int) *workflow.SolidusChannelSetting {
	return StubChannelSetting
}
