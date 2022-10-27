package usecases

type EVCommunicationStandardType string

const (
	EVCommunicationStandardTypeUnknown      EVCommunicationStandardType = "unknown"
	EVCommunicationStandardTypeISO151182ED1 EVCommunicationStandardType = "iso15118-2ed1"
	EVCommunicationStandardTypeISO151182ED2 EVCommunicationStandardType = "iso15118-2ed2"
	EVCommunicationStandardTypeIEC61851     EVCommunicationStandardType = "iec61851"
)

// Interface for the evCC use case for CEM device
type EMobilityI interface {
	// handle device state updates from the remote EV entity
	HandleEVEntityState(ski string, failure bool)
}
