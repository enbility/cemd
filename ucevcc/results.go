package ucevcc

import (
	"fmt"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

func (e *UCEVCC) HandleResult(errorMsg api.ResultMessage) {
	// before SPINE 1.3 the heartbeats are on the EVSE entity
	if errorMsg.EntityRemote == nil ||
		(errorMsg.EntityRemote.EntityType() != model.EntityTypeTypeEV &&
			errorMsg.EntityRemote.EntityType() != model.EntityTypeTypeEVSE) {
		return
	}

	// handle errors coming from the remote EVSE entity
	switch errorMsg.FeatureLocal.Type() {
	case model.FeatureTypeTypeDeviceDiagnosis:
		e.handleResultDeviceDiagnosis(errorMsg)
	}
}

// Handle DeviceDiagnosis Results
func (e *UCEVCC) handleResultDeviceDiagnosis(resultMsg api.ResultMessage) {
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
		errorText := fmt.Sprintf("Error Code: %d", resultMsg.Result.ErrorNumber)
		if resultMsg.Result.Description != nil {
			errorText = fmt.Sprintf("%s - %s", errorText, string(*resultMsg.Result.Description))
		}
		e.service.DisconnectSKI(resultMsg.DeviceRemote.Ski(), errorText)
	}
}
