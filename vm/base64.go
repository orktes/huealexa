package vm

import (
	"encoding/base64"

	"github.com/dop251/goja"
)

func (vm *VM) initBase64() {
	btoa := func(call goja.FunctionCall) goja.Value {
		input := call.Argument(0).String()
		output := base64.StdEncoding.EncodeToString([]byte(input))
		return vm.ToValue(output)
	}

	atob := func(call goja.FunctionCall) goja.Value {
		input := call.Argument(0).String()
		if output, err := base64.StdEncoding.DecodeString(input); err == nil {
			return vm.ToValue(string(output))
		}

		return vm.ToValue("")
	}

	vm.Set("btoa", btoa)
	vm.Set("atob", atob)
}
