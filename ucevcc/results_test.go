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

	errorMsg := spineapi.ResultMessage{
		DeviceRemote: s.remoteDevice,
		EntityRemote: s.mockRemoteEntity,
		FeatureLocal: localFeature,
		Result: &model.ResultDataType{
			ErrorNumber: eebusutil.Ptr(model.ErrorNumberTypeNoError),
		},
	}
	s.sut.HandleResult(errorMsg)

	errorMsg.EntityRemote = s.evEntity
	s.sut.HandleResult(errorMsg)

	errorMsg.Result = &model.ResultDataType{
		ErrorNumber: eebusutil.Ptr(model.ErrorNumberTypeGeneralError),
		Description: eebusutil.Ptr(model.DescriptionType("test error")),
	}
	errorMsg.MsgCounterReference = model.MsgCounterType(500)

	s.mockSender.
		EXPECT().
		DatagramForMsgCounter(errorMsg.MsgCounterReference).
		Return(model.DatagramType{}, errors.New("test")).Once()

	s.sut.HandleResult(errorMsg)

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

	s.sut.HandleResult(errorMsg)

}
