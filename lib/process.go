package lib

import (
	"log"
	"os/exec"

	"github.com/dop251/goja"
)

func initProcess(vm *goja.Runtime) {
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
	VMSetAsyncFunction(vm, "_native_exec_async", fn)
}
