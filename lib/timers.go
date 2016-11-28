package lib

import (
	"fmt"
	"time"

	"github.com/dop251/goja"
)

func initTimers(vm *VM) {
	vm.Set("_native_set_timeout", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).ToInteger()
		timeout := call.Argument(1).ToInteger()

		go func(id, timeout int64) {
			time.Sleep(time.Duration(timeout) * time.Millisecond)
			vm.RunString(fmt.Sprintf("require('timers')._native_callback(%d)", id))
		}(id, timeout)

		return vm.ToValue(nil)
	})

	vm.RunString(`
    var setTimeout = require('timers').setTimeout;
    var setInterval = require('timers').setInterval;
    var clearTimeout = require('timers').clear;
    var clearInterval = require('timers').clear;
  `)
}
