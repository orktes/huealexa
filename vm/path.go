package vm

import (
	"path"

	"github.com/dop251/goja"
)

func (vm *VM) initPath() {
	join := func(call goja.FunctionCall) goja.Value {
		elements := make([]string, len(call.Arguments))

		for _, arg := range call.Arguments {
			elements = append(elements, arg.String())
		}

		return vm.ToValue(path.Join(elements...))
	}

	vm.Set("_path_join", join)
}
