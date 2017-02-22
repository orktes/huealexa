package vm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/dop251/goja"
	"github.com/labstack/echo"
)

func (vm *VM) initServer() {
	callbackCounter := int64(0)
	bodyMap := map[int64]io.ReadCloser{}
	bodyMapMutex := sync.Mutex{}

	vm.Set("_add_server_handler", func(call goja.FunctionCall) goja.Value {
		id := atomic.AddInt64(&callbackCounter, 1)

		method := call.Argument(0).String()
		url := call.Argument(1).String()
		handler := func(c echo.Context) error {

			callbackID, cbCh := CreateAsyncNativeCallbackChannel()

			bodyMapMutex.Lock()
			bodyMap[callbackID] = c.Request().Body
			bodyMapMutex.Unlock()

			reqData, _ := json.Marshal(map[string]interface{}{
				"id":      callbackID,
				"headers": c.Request().Header,
				"path":    c.Request().URL.Path,
				"query":   c.Request().URL.Query(),
				"base":    vm.server.URLBase,
			})

			_, err := vm.RunString(fmt.Sprintf(`
        require('server')._request(
          %d,
          %s,
          require('async').createJSCallback(%d, true)
        );
      `,
				id,
				string(reqData),
				callbackID,
			))
			if err != nil {
				return err
			}

			str := (<-cbCh).Argument(0).String()

			var response struct {
				Body       string `json:"body"`
				StatusCode int    `json:"status_code"`
				Headers    []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"headers"`
			}

			err = json.Unmarshal([]byte(str), &response)
			if err != nil {
				return nil
			}

			for _, header := range response.Headers {
				c.Response().Header().Set(header.Key, header.Value)
			}

			c.Response().WriteHeader(response.StatusCode)
			io.Copy(c.Response(), strings.NewReader(response.Body))

			bodyMapMutex.Lock()
			bodyMap[callbackID].Close()
			delete(bodyMap, callbackID)
			bodyMapMutex.Unlock()

			return nil
		}

		switch method {
		case http.MethodGet:
			vm.server.GET(url, handler)
		case http.MethodPost:
			vm.server.POST(url, handler)
		case http.MethodPut:
			vm.server.PUT(url, handler)
		case http.MethodDelete:
			vm.server.DELETE(url, handler)
		case http.MethodHead:
		}

		return vm.ToValue(id)
	})

	vm.Set("_get_server_req_body", func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).ToInteger()
		bodyMapMutex.Lock()
		defer bodyMapMutex.Unlock()

		reader, ok := bodyMap[id]
		if !ok {
			panic(errors.New("Could not find body with given request ID " + strconv.Itoa(int(id))))
		}

		data, err := ioutil.ReadAll(reader)
		if err != nil {
			panic(err)
		}

		return vm.ToValue(string(data))
	})
}
