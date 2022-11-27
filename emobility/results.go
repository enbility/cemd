package emobility

import (
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

func (e *EMobilityImpl) HandleResult(errorMsg spine.ResultMessage) {
	if errorMsg.EntityRemote == e.evseEntity {
		// handle errors coming from the remote EVSE entity
		switch errorMsg.FeatureLocal.Type() {
		case model.FeatureTypeTypeDeviceDiagnosis:
			e.handleResultDeviceDiagnosis(errorMsg)
		}

	} else if e.evEntity != nil && errorMsg.EntityRemote == e.evEntity {
		// handle errors coming from a remote EV entity
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

	if resultMsg.DeviceRemote.IsHeartbeatMsgCounter(resultMsg.MsgCounterReference) {
		resultMsg.DeviceRemote.Stopheartbeat()

		// something is horribly wrong, disconnect and hope a new connection will fix it
		e.service.DisconnectSKI(resultMsg.DeviceRemote.Ski(), string(*resultMsg.Result.Description))
	}
}
