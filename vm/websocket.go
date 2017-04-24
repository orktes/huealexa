package vm

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dop251/goja"
	"github.com/gorilla/websocket"
)

var webSocketMap = map[string]*websocket.Conn{}

func (vm *VM) initWebSocket() {

	for _, webSocket := range webSocketMap {
		webSocket.Close()
	}

	sendError := func(id string, err error) {
		errorJSON, err := json.Marshal(err.Error())
		if err != nil {
			log.Println(err)
			return
		}

		_, err = vm.RunString(fmt.Sprintf(`require('websocket')._error("%s", %s);`, id, errorJSON))
		if err != nil {
			log.Println(err)
		}
	}

	sendMessage := func(id string, message string) {
		messageJSON, err := json.Marshal(message)
		if err != nil {
			log.Println(err)
			return
		}

		_, err = vm.RunString(fmt.Sprintf(`require('websocket')._message("%s", %s);`, id, messageJSON))
		if err != nil {
			log.Println(err)
		}
	}

	sendClose := func(id string) {
		_, err := vm.RunString(fmt.Sprintf(`require('websocket')._close("%s");`, id))
		if err != nil {
			log.Println(err)
		}
	}

	sendConnect := func(id string) {
		_, err := vm.RunString(fmt.Sprintf(`require('websocket')._connect("%s");`, id))
		if err != nil {
			log.Println(err)
		}
	}

	vm.Set("_init_websocket", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).String()
		url := call.Argument(1).String()
		go func() {
			c, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				sendError(id, err)
				return
			}

			webSocketMap[id] = c

			sendConnect(id)
			defer c.Close()
			defer sendClose(id)

			for {
				messageType, message, err := c.ReadMessage()
				if err != nil {
					sendError(id, err)
					return
				}
				if messageType == websocket.CloseMessage {
					return
				}

				sendMessage(id, string(message))
			}
		}()

		return goja.Null()
	})

	vm.Set("_close_websocket", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).String()

		if ws, ok := webSocketMap[id]; ok {
			err := ws.Close()
			if err != nil {
				panic(err)
			}

			sendClose(id)
		}

		return goja.Null()
	})

}
