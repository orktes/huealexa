package lib

import "github.com/dop251/goja"

func Register(vm *goja.Runtime) {
	initUUID(vm)
	initProcess(vm)
	initSSDP(vm)
}
