package vm

import (
	"github.com/dop251/goja"
	uuid "github.com/nu7hatch/gouuid"
)

func (vm *VM) initUUID() {
	vm.Set("_native_uuid_v4", func(call goja.FunctionCall) goja.Value {
		genuuid, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}
		return vm.ToValue(genuuid.String())
	})
}
