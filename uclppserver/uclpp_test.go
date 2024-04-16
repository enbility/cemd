package uclppserver

import (
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UCLPPServerSuite) Test_UpdateUseCaseAvailability() {
	s.sut.UpdateUseCaseAvailability(true)
}

func (s *UCLPPServerSuite) Test_IsUseCaseSupported() {
	data, err := s.sut.IsUseCaseSupported(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	data, err = s.sut.IsUseCaseSupported(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), false, data)

	ucData := &model.NodeManagementUseCaseDataType{
		UseCaseInformation: []model.UseCaseInformationDataType{
			{
				Actor: eebusutil.Ptr(model.UseCaseActorTypeEnergyGuard),
				UseCaseSupport: []model.UseCaseSupportType{
					{
						UseCaseName:      eebusutil.Ptr(model.UseCaseNameTypeLimitationOfPowerProduction),
						UseCaseAvailable: eebusutil.Ptr(true),
						ScenarioSupport:  []model.UseCaseScenarioSupportType{1, 2, 3, 4},
					},
				},
			},
		},
	}

	nodemgmtEntity := s.remoteDevice.Entity([]model.AddressEntityType{0})
	nodeFeature := s.remoteDevice.FeatureByEntityTypeAndRole(nodemgmtEntity, model.FeatureTypeTypeNodeManagement, model.RoleTypeSpecial)
	fErr := nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, ucData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.IsUseCaseSupported(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), true, data)
}
