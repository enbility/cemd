package cem

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// Handle events from eebus-go library
func (h *CemImpl) HandleEvent(payload spine.EventPayload) {
	switch payload.EventType {
	case spine.EventTypeSubscriptionChange:
		switch payload.Data.(type) {
		case model.SubscriptionManagementRequestCallType:
			h.subscriptionRequestHandling(payload)
		}
	}
}

// Handle subscription requests
func (h *CemImpl) subscriptionRequestHandling(payload spine.EventPayload) {
	data := payload.Data.(model.SubscriptionManagementRequestCallType)

	// Heartbeat subscription requests
	if *data.ServerFeatureType == model.FeatureTypeTypeDeviceDiagnosis {
		remoteDevice := h.myService.RemoteDeviceForSki(payload.Ski)
		if remoteDevice == nil {
			fmt.Println("No remote device found for SKI:", payload.Ski)
			return
		}

		senderAddr := h.myService.LocalDevice().FeatureByTypeAndRole(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer).Address()
		destinationAddr := payload.Feature.Address()
		if senderAddr == nil || destinationAddr == nil {
			fmt.Println("No sender or destination address found for SKI:", payload.Ski)
			return
		}

		switch payload.ChangeType {
		case spine.ElementChangeAdd:
			// start sending heartbeats
			remoteDevice.StartHeartbeatSend(senderAddr, destinationAddr)
		case spine.ElementChangeRemove:
			// stop sending heartbeats
			remoteDevice.Stopheartbeat()
		}
	}
}
