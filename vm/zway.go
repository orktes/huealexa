package vm

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dop251/goja"
	"github.com/orktes/huessimo/zway"
)

var zWayMap = map[string]*zway.ZWay{}

func (vm *VM) initZWay() {

	for _, zWay := range zWayMap {
		zWay.Stop()
	}

	vm.Set("_init_zway", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).String()
		zwayConfigStr := call.Argument(1).String()

		config := &zway.Config{}
		json.Unmarshal([]byte(zwayConfigStr), config)

		zWay, err := zway.New(*config)
		if err != nil {
			panic(err)
		}

		zWayMap[id] = zWay

		zWay.Devices.AddNewDeviceListener(func(device *zway.Device) {
			deviceStr, err := json.Marshal(device)
			if err != nil {
				log.Println(err)
				return
			}

			_, err = vm.RunString(fmt.Sprintf(`require('zway')._device_added("%s", %s);`, id, deviceStr))
			if err != nil {
				log.Println(err)
			}

			device.Metrics.Level.AddValueChangeListener(func(level *zway.Level) {
				valueStr, err := json.Marshal(level)
				if err != nil {
					log.Println(err)
					return
				}
				_, err = vm.RunString(fmt.Sprintf(`require('zway')._value_change("%s", "%s", %s);`, id, device.ID, valueStr))
				if err != nil {
					log.Println(err)
				}
			})
		})

		return goja.Null()
	})

}
