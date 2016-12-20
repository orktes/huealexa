package vm

import (
	"log"
	"os/exec"

	"github.com/dop251/goja"
)

func (vm *VM) initProcess() {
	fn := func(call goja.FunctionCall) goja.Value {
		cmd := call.Argument(0).String()

		log.Printf("[JS][SH]: %s\n", cmd)

		out, outErr := exec.Command("sh", "-c", cmd).Output()
		if outErr != nil {
			panic(vm.ToValue(outErr.Error()))
		}

		return vm.ToValue(string(out))
	}
	vm.Set("_native_exec", fn)
	vm.SetAsyncFunction("_native_exec_async", fn)
}
