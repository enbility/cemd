package uclpcserver

import (
	eebusutil "github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

func (s *UCLPCServerSuite) Test_Events() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.HandleEvent(payload)

	payload.Entity = s.monitoredEntity
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeEntityChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.CmdClassifier = eebusutil.Ptr(model.CmdClassifierTypeWrite)
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.Function = model.FunctionTypeLoadControlLimitListData
	payload.Data = eebusutil.Ptr(model.LoadControlLimitListDataType{})
	s.sut.HandleEvent(payload)

	payload.LocalFeature = s.loadControlFeature
	s.sut.HandleEvent(payload)

	payload.Function = model.FunctionTypeDeviceConfigurationKeyValueListData
	payload.Data = eebusutil.Ptr(model.DeviceConfigurationKeyValueListDataType{})
	s.sut.HandleEvent(payload)

	payload.LocalFeature = s.deviceConfigurationFeature
	s.sut.HandleEvent(payload)
}
