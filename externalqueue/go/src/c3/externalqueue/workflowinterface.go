package externalqueue

import (
	"c3/osm/workflow"
)

type WorkflowInterface interface {
	SaveObject(theObject WorkflowSaveInterface) error
	QueryClearInteractAction_byErrand(id int) *workflow.ClearInteractAction
	QuerySystemGroup_minimalByAreaID(id int) *workflow.SystemGroup
	LoadArea(areaId int) *workflow.Area
	QuerySolidusChannelSetting_byChannel(sGroup int, sType int) *workflow.SolidusChannelSetting
}

type WorkflowSaveInterface interface {
	Save() error
}

//WorkflowObject uses the WorkflowInterface interface
type WorkflowObject struct {
	wf WorkflowInterface
}

func (o *WorkflowObject) InitWorkflow(theObject WorkflowInterface) {
	o.wf = theObject
}

func (o *WorkflowObject) SaveObject(theObject WorkflowSaveInterface) error {
	return o.wf.SaveObject(theObject)
}

func (o *WorkflowObject) QueryClearInteractAction_byErrand(id int) *workflow.ClearInteractAction {
	return o.wf.QueryClearInteractAction_byErrand(id)
}

func (o *WorkflowObject) QuerySystemGroup_minimalByAreaID(id int) *workflow.SystemGroup {
	return o.wf.QuerySystemGroup_minimalByAreaID(id)
}

func (o *WorkflowObject) LoadArea(id int) *workflow.Area {
	return o.wf.LoadArea(id)
}

func (o *WorkflowObject) QuerySolidusChannelSetting_byChannel(sGroup int,
	sType int) *workflow.SolidusChannelSetting {
	return o.wf.QuerySolidusChannelSetting_byChannel(sGroup, sType)
}

//WorkflowClass implements the WorkflowInterface
type WorkflowClass struct {
}

func (o *WorkflowClass) SaveObject(theObject WorkflowSaveInterface) error {
	return theObject.Save()
}

func (o *WorkflowClass) QueryClearInteractAction_byErrand(id int) *workflow.ClearInteractAction {
	action, err := workflow.QueryClearInteractAction_byErrand(id)
	if err != nil {
		return nil
	}
	return action
}

func (o *WorkflowClass) QuerySystemGroup_minimalByAreaID(id int) *workflow.SystemGroup {
	sysGroup, err := workflow.QuerySystemGroup_minimalByAreaID(id)
	if err != nil {
		return nil
	}
	return sysGroup
}

func (o *WorkflowClass) LoadArea(id int) *workflow.Area {
	area, err := workflow.LoadArea(id)
	if err != nil {
		return nil
	}
	return area
}

func (o *WorkflowClass) QuerySolidusChannelSetting_byChannel(sGroup int,
	sType int) *workflow.SolidusChannelSetting {
	setting, err := workflow.QuerySolidusChannelSetting_byChannel(sGroup, sType)
	if err != nil {
		return nil
	}
	return setting
}
