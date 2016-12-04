package vm

import (
	"github.com/dop251/goja"
	"github.com/huin/goupnp"
)

func initSSDP(vm *VM) {
	VMSetAsyncFunction(vm, "_native_ssdp_discover_devices", func(call goja.FunctionCall) goja.Value {
		search := call.Argument(0).String()
		responses, err := goupnp.DiscoverDevices(search)
		if err != nil {
			panic(err)
		}
		return vm.ToValue(responses)
	})
}
