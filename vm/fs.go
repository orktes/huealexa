package vm

import (
	"io/ioutil"
	"os"

	"github.com/dop251/goja"
)

func (vm *VM) initFS() {
	write := func(call goja.FunctionCall) goja.Value {
		file := call.Argument(0).String()
		data := call.Argument(1).String()
		flag := call.Argument(2).ToInteger()

		err := ioutil.WriteFile(file, []byte(data), os.FileMode(flag))
		if err != nil {
			panic(vm.ToValue(err.Error()))
		}

		return goja.Undefined()
	}

	read := func(call goja.FunctionCall) goja.Value {
		file := call.Argument(0).String()

		data, err := ioutil.ReadFile(file)
		if err != nil {
			panic(vm.ToValue(err.Error()))
		}

		return vm.ToValue(string(data))
	}

	vm.Set("_fs_read", read)
	vm.SetAsyncFunction("_fs_read_async", read)

	vm.Set("_fs_write", write)
	vm.SetAsyncFunction("_fs_write_async", write)
}
