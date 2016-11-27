package lib

import "github.com/dop251/goja"

func Register(vm *goja.Runtime) {
	initAsync(vm)
	initTimers(vm)
	initUUID(vm)
	initProcess(vm)
	initSSDP(vm)
}
