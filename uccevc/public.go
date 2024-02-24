package uccevc

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// helper

func (e *UCCEVC) isCompatibleEntity(entity spineapi.EntityRemoteInterface) bool {
	if entity == nil || entity.EntityType() != model.EntityTypeTypeEV {
		return false
	}

	return true
}
