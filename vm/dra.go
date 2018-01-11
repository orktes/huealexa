package vm

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dop251/goja"
	"github.com/orktes/go-dra"
)

var draMap = map[string]*dra.DRA{}

func (vm *VM) initDRA() {

	for _, d := range draMap {
		d.Close()
	}

	sendError := func(id string, err error) {
		errorJSON, err := json.Marshal(err.Error())
		if err != nil {
			log.Println(err)
			return
		}

		_, err = vm.RunString(fmt.Sprintf(`require('devices/audio/denon_dra')._error("%s", %s);`, id, errorJSON))
		if err != nil {
			log.Println(err)
		}
	}

	sendUpdate := func(id string, message string) {
		messageJSON, err := json.Marshal(message)
		if err != nil {
			log.Println(err)
			return
		}

		_, err = vm.RunString(fmt.Sprintf(`require('devices/audio/denon_dra')._update("%s", %s);`, id, messageJSON))
		if err != nil {
			log.Println(err)
		}
	}

	sendClose := func(id string) {
		_, err := vm.RunString(fmt.Sprintf(`require('devices/audio/denon_dra')._close("%s");`, id))
		if err != nil {
			log.Println(err)
		}
	}

	sendConnect := func(id string) {
		_, err := vm.RunString(fmt.Sprintf(`require('devices/audio/denon_dra')._connect("%s");`, id))
		if err != nil {
			log.Println(err)
		}
	}

	vm.Set("_init_dra", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).String()
		addr := call.Argument(1).String()
		go func() {
			d, err := dra.NewFromAddr(addr)
			d.OnUpdate = make(chan string, 10)
			if err != nil {
				sendError(id, err)
				return
			}

			draMap[id] = d

			sendConnect(id)
			defer d.Close()
			defer sendClose(id)

			for update := range d.OnUpdate {
				sendUpdate(id, update)
			}
		}()

		return goja.Null()
	})

	vm.Set("_dra_master_volume", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).String()
		volume := call.Argument(1).ToInteger()

		if d, ok := draMap[id]; ok {
			d.SetMasterVolume(int(volume))
		}

		return goja.Null()
	})

	vm.Set("_dra_power", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).String()
		power := call.Argument(1).ToBoolean()

		if d, ok := draMap[id]; ok {
			d.SetPower(power)
		}

		return goja.Null()
	})

	vm.Set("_dra_mute", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).String()
		val := call.Argument(1).ToBoolean()

		if d, ok := draMap[id]; ok {
			d.SetMute(val)
		}

		return goja.Null()
	})

	vm.Set("_close_dra", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).String()

		if ws, ok := draMap[id]; ok {
			err := ws.Close()
			if err != nil {
				panic(err)
			}

			sendClose(id)
		}

		return goja.Null()
	})

}
