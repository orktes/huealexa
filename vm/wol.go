package vm

import (
	"github.com/dop251/goja"
	wol "github.com/sabhiram/go-wol"
)

func (vm *VM) initWOL() {
	vm.Set("_wol", func(call goja.FunctionCall) goja.Value {
		macAddr := call.Argument(0).String()
		broadcastIP := call.Argument(1).String()
		broadcastPort := call.Argument(2).String()
		broadcastInterface := call.Argument(3).String()

		err := wol.SendMagicPacket(macAddr, broadcastIP+":"+broadcastPort, broadcastInterface)
		if err != nil {
			panic(err)
		}

		return vm.ToValue(nil)
	})
}
