package lib

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/dop251/goja"
)

var nativeAsyncResponseCounter = int64(0)
var nativeAsyncResponseCallbacks = map[int64]func(call goja.FunctionCall){}
var nativeAsyncResponseCallbacksLock = &sync.Mutex{}

func CreateAsyncNativeCallback(callback func(call goja.FunctionCall)) int64 {
	// TODO make thread sage
	id := atomic.AddInt64(&nativeAsyncResponseCounter, 1)
	nativeAsyncResponseCallbacksLock.Lock()
	defer nativeAsyncResponseCallbacksLock.Unlock()
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

func CreateAsyncJSFunction(orgCall func(goja.FunctionCall) goja.Value, vm *VM) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		go func(call goja.FunctionCall) {
			callbackID := call.Argument(len(call.Arguments) - 1).ToInteger()
			defer func() {
				if r := recover(); r != nil {
					vm.RunString(fmt.Sprintf(`require('async')._native_callback(%d, new Error('%s'), null, true);`, callbackID, r))
				}
			}()
			res := orgCall(goja.FunctionCall{Arguments: call.Arguments[:len(call.Arguments)-1]})
			data, merr := json.Marshal(res.Export())
			if merr != nil {
				panic(merr)
			}

			vm.RunString(fmt.Sprintf(`require('async')._native_callback(%d, null, %s, null, true);`, callbackID, string(data)))
		}(call)

		return goja.Null()
	}
}

func VMSetAsyncFunction(vm *VM, name string, fn func(call goja.FunctionCall) goja.Value) {
	vm.Set(name+"_raw", CreateAsyncJSFunction(fn, vm))
	vm.RunString(fmt.Sprintf(`
    var %s = require('async').createAsyncFunction(%s);
  `, name, name+"_raw"))
}

func initAsync(vm *VM) {
	vm.Set("_native_async_response", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).ToInteger()
		nativeAsyncResponseCallbacksLock.Lock()
		defer nativeAsyncResponseCallbacksLock.Unlock()
		if cb, ok := nativeAsyncResponseCallbacks[id]; ok {
			cb(goja.FunctionCall{Arguments: call.Arguments[1:]})
		}
		return vm.ToValue(nil)
	})
}
