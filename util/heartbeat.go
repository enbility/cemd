package util

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

func IsHeartbeat(localEntity spineapi.EntityLocalInterface, payload spineapi.EventPayload) bool {
	//revive:disable-next-line
	switch payload.Data.(type) {
	case *model.DeviceDiagnosisHeartbeatDataType:
		return payload.Function == model.FunctionTypeDeviceDiagnosisHeartbeatData &&
			payload.EventType == spineapi.EventTypeDataChange &&
			payload.ChangeType == spineapi.ElementChangeUpdate &&
			payload.CmdClassifier != nil &&
			*payload.CmdClassifier == model.CmdClassifierTypeNotify
	default:
		return false
	}
}
