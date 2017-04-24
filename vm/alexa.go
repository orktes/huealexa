package vm

import (
	"encoding/json"
	"fmt"

	"github.com/dop251/goja"
	avs "github.com/rogertalk/go-avs"
)

var alexaDownChannelOpen = false

func (vm *VM) initAlexa() {
	// TODO figure out how to close the Downchannel if JS context reloads
	vm.Set("_alexa_downchannel", func(call goja.FunctionCall) goja.Value {
		if alexaDownChannelOpen {
			return goja.Undefined()
		}

		accessToken := call.Argument(0).String()
		go func() {
			alexaDownChannelOpen = true
			defer func() {
				vm.RunString(`require('alexa')._downchannel_closed();`)
				alexaDownChannelOpen = false
			}()

			println(accessToken)
			directives, err := avs.CreateDownchannel(accessToken)
			if err != nil {
				panic(err)
			}

			println("Listening to down channel")
			for directive := range directives {
				println("Got directive")
				b, _ := json.Marshal(directive)
				vm.RunString(fmt.Sprintf(`require('alexa')._downchannel(%s);`, string(b)))
			}
		}()

		return goja.Undefined()
	})
}
