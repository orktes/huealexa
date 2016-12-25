package vm

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type HomeKitLightSensor struct {
	*accessory.Accessory
	LightSensor *service.LightSensor
}

// NewHomeKitLightSensor returns a light sensor which implements service.LightSensor.
func NewHomeKitLightSensor(info accessory.Info) *HomeKitLightSensor {
	acc := HomeKitLightSensor{}
	acc.Accessory = accessory.New(info, accessory.TypeUnknown)
	acc.LightSensor = service.NewLightSensor()

	acc.LightSensor.CurrentAmbientLightLevel.SetMaxValue(100000)
	acc.LightSensor.CurrentAmbientLightLevel.SetMinValue(0)
	acc.LightSensor.CurrentAmbientLightLevel.SetStepValue(0.0001)

	acc.AddService(acc.LightSensor.Service)

	return &acc
}
