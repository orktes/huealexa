package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/orktes/huessimo/hueserver"
	"github.com/orktes/huessimo/hueupnp"
)

func getIPAddress() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			panic(err)
		}

		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if !ip.IsLoopback() && ip.To4() != nil {
				return ip.String()
			}

		}
	}

	return ""
}

func main() {
	uuid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	scriptSrcPtr := flag.String("src", "", "Script source file location")
	uuidPtr := flag.String("uuid", "", "UUID for the HUE server (for example \""+uuid.String()+"\")")
	portPtr := flag.String("port", "8989", "Port for the HUE server")
	upnpPortPtr := flag.String("upnp", "239.255.255.250:1900", "UPNP multicast addr for the HUE server")
	namePtr := flag.String("name", "fakeServer", "Name for the HUE server")
	ipPtr := flag.String("ip", getIPAddress(), "Interface to be used")

	flag.Parse()

	if *uuidPtr == "" {
		fmt.Printf("You must provide -uuid=\"%s\" (i just generated that for you) or something else\n", uuid.String())
		return
	}

	registry := new(require.Registry)
	vm := goja.New()
	registry.Enable(vm)
	console.Enable(vm)

	vm.Set("exec", func(call goja.FunctionCall) goja.Value {
		cmd := call.Argument(0).String()

		fmt.Printf("[JS][SH]: %s\n", cmd)

		out, outErr := exec.Command("sh", "-c", cmd).Output()
		if outErr != nil {
			panic(err)
		}

		return vm.ToValue(string(out))
	})

	script, err := ioutil.ReadFile(*scriptSrcPtr)
	if err != nil {
		panic(err)
	}

	_, err = vm.RunString(string(script))
	if err != nil {
		panic(err)
	}

	getLights := func() hueserver.LightList {
		value, err := vm.RunString(`JSON.stringify(getLights());`)
		if err != nil {
			panic(err)
		}

		str := value.String()

		list := &hueserver.LightList{}
		json.Unmarshal([]byte(str), list)
		return *list
	}

	getLight := func(id string) hueserver.Light {
		arg, err := json.Marshal(id)
		if err != nil {
			panic(err)
		}
		value, err := vm.RunString(`JSON.stringify(getLight(` + string(arg) + `));`)
		if err != nil {
			panic(err)
		}

		str := value.String()

		light := &hueserver.Light{}
		json.Unmarshal([]byte(str), light)
		return *light
	}

	setLightState := func(id string, state hueserver.LightStateChange) hueserver.LightStateChangeResponse {
		arg, err := json.Marshal(id)
		if err != nil {
			panic(err)
		}

		arg2, err := json.Marshal(state)
		if err != nil {
			panic(err)
		}

		value, err := vm.RunString(`JSON.stringify(setLightState(` + string(arg) + `, ` + string(arg2) + `));`)
		if err != nil {
			panic(err)
		}

		str := value.String()

		resp := &hueserver.LightStateChangeResponse{}
		json.Unmarshal([]byte(str), resp)
		return *resp
	}

	go hueupnp.CreateUPNPResponder("http://"+*ipPtr+":"+*portPtr+"/upnp/setup.xml", *uuidPtr, *upnpPortPtr)

	srv := hueserver.NewServer(*uuidPtr, *ipPtr+":"+*portPtr, *namePtr, getLights, getLight, setLightState)
	srv.Start(":" + *portPtr)
}
