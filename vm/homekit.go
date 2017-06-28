package vm

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/dop251/goja"
)

var homeKitTransports = []hc.Transport{}

func (vm *VM) initHomeKit() {

	for _, transport := range homeKitTransports {
		transport.Stop()
	}
	homeKitTransports = []hc.Transport{}

	devices := map[string]interface{}{}

	vm.Set("_add_homekit_device", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).String()
		deviceType := call.Argument(1).String()
		pin := call.Argument(2).String()
		deviceInfoStr := call.Argument(3).String()

		info := &accessory.Info{}
		json.Unmarshal([]byte(deviceInfoStr), info)

		switch deviceType {
		case "lightbulb":
			acc := accessory.NewLightbulb(*info)
			config := hc.Config{Pin: pin, StoragePath: vm.dataDir + "/homekit/" + info.Name}
			t, err := hc.NewIPTransport(config, acc.Accessory)
			if err != nil {
				log.Panic(err)
			}

			acc.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
				var state = "true"
				if !on {
					state = "false"
				}
				_, err := vm.RunString(fmt.Sprintf(`require('homekit')._remote_on_change("%s", %s);`, id, state))
				if err != nil {
					log.Println(err)
				}
			})

			acc.Lightbulb.Brightness.OnValueRemoteUpdate(func(bri int) {
				_, err := vm.RunString(fmt.Sprintf(`require('homekit')._remote_bri_change("%s", %d);`, id, bri))
				if err != nil {
					log.Println(err)
				}
			})

			acc.Lightbulb.Hue.OnValueRemoteUpdate(func(hue float64) {
				_, err := vm.RunString(fmt.Sprintf(`require('homekit')._remote_hue_change("%s", %f);`, id, hue))
				if err != nil {
					log.Println(err)
				}
			})

			acc.Lightbulb.Saturation.OnValueRemoteUpdate(func(sat float64) {
				_, err := vm.RunString(fmt.Sprintf(`require('homekit')._remote_sat_change("%s", %f);`, id, sat))
				if err != nil {
					log.Println(err)
				}
			})

			hc.OnTermination(func() {
				t.Stop()
			})

			homeKitTransports = append(homeKitTransports, t)
			devices[id] = acc

			go t.Start()
		case "temperature_sensor":
			// Reasonable values for example for DHT22 sensors
			acc := accessory.NewTemperatureSensor(*info, 0, -40, 80, 0.1)
			config := hc.Config{Pin: pin, StoragePath: vm.dataDir + "/homekit/" + info.Name}
			t, err := hc.NewIPTransport(config, acc.Accessory)
			if err != nil {
				log.Panic(err)
			}

			hc.OnTermination(func() {
				t.Stop()
			})

			homeKitTransports = append(homeKitTransports, t)
			devices[id] = acc

			go t.Start()
		case "humidity_sensor":
			acc := NewHomeKitHumiditySensor(*info, 0, 0, 100, 0.1)
			config := hc.Config{Pin: pin, StoragePath: vm.dataDir + "/homekit/" + info.Name}
			t, err := hc.NewIPTransport(config, acc.Accessory)
			if err != nil {
				log.Panic(err)
			}

			hc.OnTermination(func() {
				t.Stop()
			})

			homeKitTransports = append(homeKitTransports, t)
			devices[id] = acc

			go t.Start()
		case "door":
			acc := NewHomeKitDoor(*info)
			config := hc.Config{Pin: pin, StoragePath: vm.dataDir + "/homekit/" + info.Name}
			t, err := hc.NewIPTransport(config, acc.Accessory)
			if err != nil {
				log.Panic(err)
			}

			hc.OnTermination(func() {
				t.Stop()
			})

			homeKitTransports = append(homeKitTransports, t)
			devices[id] = acc

			go t.Start()

		case "light_sensor":
			acc := NewHomeKitLightSensor(*info)
			config := hc.Config{Pin: pin, StoragePath: vm.dataDir + "/homekit/" + info.Name}
			t, err := hc.NewIPTransport(config, acc.Accessory)
			if err != nil {
				log.Panic(err)
			}

			hc.OnTermination(func() {
				t.Stop()
			})

			homeKitTransports = append(homeKitTransports, t)
			devices[id] = acc

			go t.Start()
		}

		return goja.Null()
	})

	vm.Set("_set_homekit_device_on", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).String()
		value := call.Argument(1).ToBoolean()
		device := devices[id]

		switch device.(type) {
		case *accessory.Lightbulb:
			device.(*accessory.Lightbulb).Lightbulb.On.SetValue(value)
		case *HomeKitDoor:
			intValue := 0
			if value {
				intValue = 1
			}

			device.(*HomeKitDoor).Door.CurrentPosition.SetValue(intValue)
		}
		return goja.Null()
	})

	vm.Set("_set_homekit_device_bri", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).String()
		value := call.Argument(1).ToInteger()
		device := devices[id]

		switch device.(type) {
		case *accessory.Lightbulb:
			device.(*accessory.Lightbulb).Lightbulb.Brightness.SetValue(int(value))
		}
		return goja.Null()
	})

	vm.Set("_set_homekit_device_hue", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).String()
		value := call.Argument(1).ToFloat()
		device := devices[id]

		switch device.(type) {
		case *accessory.Lightbulb:
			device.(*accessory.Lightbulb).Lightbulb.Hue.SetValue(value)
		}
		return goja.Null()
	})

	vm.Set("_set_homekit_device_temperature", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).String()
		value := call.Argument(1).ToFloat()
		device := devices[id]

		switch device.(type) {
		case *accessory.Thermometer:
			device.(*accessory.Thermometer).TempSensor.CurrentTemperature.SetValue(value)
		}
		return goja.Null()
	})

	vm.Set("_set_homekit_device_humidity", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).String()
		value := call.Argument(1).ToFloat()
		device := devices[id]

		switch device.(type) {
		case *HomeKitHumiditySensor:
			device.(*HomeKitHumiditySensor).HumiditySensor.CurrentRelativeHumidity.SetValue(value)
		}
		return goja.Null()
	})

	vm.Set("_set_homekit_device_lux", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).String()
		value := call.Argument(1).ToFloat()
		device := devices[id]

		switch device.(type) {
		case *HomeKitLightSensor:
			device.(*HomeKitLightSensor).LightSensor.CurrentAmbientLightLevel.SetValue(value)
		}
		return goja.Null()
	})

	vm.Set("_set_homekit_device_sat", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).String()
		value := call.Argument(1).ToFloat()
		device := devices[id]

		switch device.(type) {
		case *accessory.Lightbulb:
			device.(*accessory.Lightbulb).Lightbulb.Saturation.SetValue(value)
		}
		return goja.Null()
	})

}
