package ucevcc

import (
	"errors"

	eebusutil "github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

func (s *UCEVCCSuite) Test_Results() {
	localDevice := s.service.LocalDevice()
	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)
	localFeature := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)

	errorMsg := spineapi.ResponseMessage{
		DeviceRemote: s.remoteDevice,
		EntityRemote: s.evEntity,
		FeatureLocal: localFeature,
		Data:         eebusutil.Ptr(model.MsgCounterType(0)),
	}
	s.sut.HandleResponse(errorMsg)

	errorMsg = spineapi.ResponseMessage{
		EntityRemote: s.evEntity,
		FeatureLocal: localFeature,
		Data:         eebusutil.Ptr(model.MsgCounterType(0)),
	}
	s.sut.HandleResponse(errorMsg)

	errorMsg = spineapi.ResponseMessage{
		DeviceRemote: s.remoteDevice,
		EntityRemote: s.mockRemoteEntity,
		FeatureLocal: localFeature,
		Data: &model.ResultDataType{
			ErrorNumber: eebusutil.Ptr(model.ErrorNumberTypeNoError),
		},
	}
	s.sut.HandleResponse(errorMsg)

	errorMsg.EntityRemote = s.evEntity
	s.sut.HandleResponse(errorMsg)

	errorMsg.Data = &model.ResultDataType{
		ErrorNumber: eebusutil.Ptr(model.ErrorNumberTypeGeneralError),
		Description: eebusutil.Ptr(model.DescriptionType("test error")),
	}
	errorMsg.MsgCounterReference = model.MsgCounterType(500)

	s.mockSender.
		EXPECT().
		DatagramForMsgCounter(errorMsg.MsgCounterReference).
		Return(model.DatagramType{}, errors.New("test")).Once()

	s.sut.HandleResponse(errorMsg)

	datagram := model.DatagramType{
		Payload: model.PayloadType{
			Cmd: []model.CmdType{
				{
					DeviceDiagnosisHeartbeatData: &model.DeviceDiagnosisHeartbeatDataType{},
				},
			},
		},
	}
	s.mockSender.
		EXPECT().
		DatagramForMsgCounter(errorMsg.MsgCounterReference).
		Return(datagram, nil).Once()

	s.sut.HandleResponse(errorMsg)
}
