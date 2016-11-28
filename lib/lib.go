package lib

import (
	"sync"

	"github.com/dop251/goja"
)

type VM struct {
	*goja.Runtime
	sync.Mutex
}

func (vm *VM) Lock() {
	vm.Mutex.Lock()
}

func (vm *VM) Unlock() {
	vm.Mutex.Unlock()
}

func (vm *VM) RunString(str string) (goja.Value, error) {
	vm.Lock()
	defer vm.Unlock()

	return vm.Runtime.RunString(str)
}

func (vm *VM) RunScript(name, value string) (goja.Value, error) {
	vm.Lock()
	defer vm.Unlock()

	return vm.Runtime.RunScript(name, value)
}

func Register(vm *VM) {
	initAsync(vm)
	initTimers(vm)
	initUUID(vm)
	initProcess(vm)
	initSSDP(vm)
}
