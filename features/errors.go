package features

import "errors"

// ErrMetadataNotAvailable indicates that the meta data information is not available
// e.g. decsriptions, constraints, ...
var ErrMetadataNotAvailable = errors.New("meta data not available")

// ErrDataNotAvailable indicates that no data set is yet available
var ErrDataNotAvailable = errors.New("data not available")

// ErrDataForMetadataKeyNotFound indicates that no data item is found for the given key
var ErrDataForMetadataKeyNotFound = errors.New("data for key not found")
