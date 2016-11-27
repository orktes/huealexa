package lib

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/dop251/goja"
)

var nativeAsyncResponseCounter = int64(0)
var nativeAsyncResponseCallbacks = map[int64]func(call goja.FunctionCall){}

func CreateAsyncNativeCallback(callback func(call goja.FunctionCall)) int64 {
	// TODO make thread sage
	id := atomic.AddInt64(&nativeAsyncResponseCounter, 1)
	nativeAsyncResponseCallbacks[id] = callback
	return id
}

func CreateAsyncNativeCallbackChannel() (int64, chan goja.FunctionCall) {
	ch := make(chan goja.FunctionCall, 1)
	id := CreateAsyncNativeCallback(func(call goja.FunctionCall) {
		ch <- call
	})
	return id, ch
}

func CreateAsyncJSFunction(orgCall func(goja.FunctionCall) goja.Value, vm *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		go func(call goja.FunctionCall) {
			callbackID := call.Argument(0).ToInteger()
			defer func() {
				if r := recover(); r != nil {
					vm.RunString(fmt.Sprintf(`require('async')._native_callback(%d, new Error('%s'));`, callbackID, r))
				}
			}()
			res := orgCall(goja.FunctionCall{Arguments: call.Arguments[:len(call.Arguments)-1]})
			data, merr := json.Marshal(res.Export())
			if merr != nil {
				panic(merr)
			}

			vm.RunString(fmt.Sprintf(`require('async')._native_callback(%d, null, %s);`, callbackID, string(data)))
		}(call)

		return goja.Null()
	}
}

func VMSetAsyncFunction(vm *goja.Runtime, name string, fn func(call goja.FunctionCall) goja.Value) {
	vm.Set(name+"_raw", CreateAsyncJSFunction(fn, vm))
	vm.RunString(fmt.Sprintf(`
    var %s = require('async').createAsyncFunction(%s);
  `, name, name+"_raw"))
}

func initAsync(vm *goja.Runtime) {
	vm.Set("_native_async_response", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).ToInteger()

		if cb, ok := nativeAsyncResponseCallbacks[id]; ok {
			cb(goja.FunctionCall{Arguments: call.Arguments[1:]})
		}
		return vm.ToValue(nil)
	})
}
