package lib

import (
	"fmt"
	"os/exec"

	"github.com/dop251/goja"
)

func initProcess(vm *goja.Runtime) {
	vm.Set("_native_exec", func(call goja.FunctionCall) goja.Value {
		cmd := call.Argument(0).String()

		fmt.Printf("[JS][SH]: %s\n", cmd)

		out, outErr := exec.Command("sh", "-c", cmd).Output()
		if outErr != nil {
			panic(outErr)
		}

		return vm.ToValue(string(out))
	})
}
