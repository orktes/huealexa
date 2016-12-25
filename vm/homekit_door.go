package vm

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type HomeKitDoor struct {
	*accessory.Accessory
	Door *service.Door
}

// NewHomeKitDoor returns a door which implements service.Door.
func NewHomeKitDoor(info accessory.Info) *HomeKitDoor {
	acc := HomeKitDoor{}
	acc.Accessory = accessory.New(info, accessory.TypeDoor)
	acc.Door = service.NewDoor()

	acc.Door.CurrentPosition.SetValue(0)
	acc.Door.CurrentPosition.SetMinValue(0)
	acc.Door.CurrentPosition.SetMaxValue(4)
	acc.Door.CurrentPosition.SetStepValue(1)

	acc.AddService(acc.Door.Service)

	return &acc
}
