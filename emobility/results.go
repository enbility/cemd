package emobility

import (
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

func (e *EMobilityImpl) HandleResult(errorMsg spine.ResultMessage) {
	isEvse := errorMsg.EntityRemote == e.evseEntity
	isEv := e.evEntity != nil && errorMsg.EntityRemote == e.evEntity

	if isEvse || isEv {
		// handle errors coming from the remote EVSE entity
		switch errorMsg.FeatureLocal.Type() {
		case model.FeatureTypeTypeDeviceDiagnosis:
			e.handleResultDeviceDiagnosis(errorMsg)
		}

	}
}

// Handle DeviceDiagnosis Results
func (e *EMobilityImpl) handleResultDeviceDiagnosis(resultMsg spine.ResultMessage) {
	// is this an error for a heartbeat message?
	if *resultMsg.Result.ErrorNumber == model.ErrorNumberTypeNoError {
		return
	}

	// check if this is for a cached notify message
	datagram, err := resultMsg.DeviceRemote.Sender().DatagramForMsgCounter(resultMsg.MsgCounterReference)
	if err != nil {
		return
	}

	if len(datagram.Payload.Cmd) > 0 &&
		datagram.Payload.Cmd[0].DeviceDiagnosisHeartbeatData != nil {
		// something is horribly wrong, disconnect and hope a new connection will fix it
		e.service.DisconnectSKI(resultMsg.DeviceRemote.Ski(), string(*resultMsg.Result.Description))
	}
}
