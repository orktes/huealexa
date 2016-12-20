package vm

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

func (vm *VM) CreateAsyncJSFunction(orgCall func(goja.FunctionCall) goja.Value) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		go func(call goja.FunctionCall) {
			callbackID := call.Argument(len(call.Arguments) - 1).ToInteger()
			defer func() {
				if r := recover(); r != nil {
					_, err := vm.RunString(fmt.Sprintf(`require('async')._native_callback(%d, new Error('%s'), null, true);`, callbackID, r))
					if err != nil {
						panic(err)
					}
				}
			}()
			res := orgCall(goja.FunctionCall{Arguments: call.Arguments[:len(call.Arguments)-1]})
			data, merr := json.Marshal(res.Export())
			if merr != nil {
				panic(merr)
			}

			_, err := vm.RunString(fmt.Sprintf(`require('async')._native_callback(%d, null, %s, null, true);`, callbackID, string(data)))
			if err != nil {
				panic(err)
			}
		}(call)

		return goja.Null()
	}
}

func (vm *VM) SetAsyncFunction(name string, fn func(call goja.FunctionCall) goja.Value) {
	vm.Set(name+"_raw", vm.CreateAsyncJSFunction(fn))
	_, err := vm.RunString(fmt.Sprintf(`
    var %s = require('async').createAsyncFunction(%s);
  `, name, name+"_raw"))
	if err != nil {
		panic(err)
	}
}

func (vm *VM) initAsync() {
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
