package vm

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type HomeKitHumiditySensor struct {
	*accessory.Accessory

	HumiditySensor *service.HumiditySensor
}

// NewHomeKitHumiditySensor returns a HomeKitHumiditySensor which implements service.HumiditySensor.
func NewHomeKitHumiditySensor(info accessory.Info, current, min, max, step float64) *HomeKitHumiditySensor {
	acc := HomeKitHumiditySensor{}
	acc.Accessory = accessory.New(info, accessory.TypeUnknown)
	acc.HumiditySensor = service.NewHumiditySensor()

	acc.HumiditySensor.CurrentRelativeHumidity.SetMinValue(min)
	acc.HumiditySensor.CurrentRelativeHumidity.SetMaxValue(max)
	acc.HumiditySensor.CurrentRelativeHumidity.SetStepValue(step)
	acc.HumiditySensor.CurrentRelativeHumidity.SetValue(current)

	acc.AddService(acc.HumiditySensor.Service)

	return &acc
}
