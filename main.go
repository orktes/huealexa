package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/orktes/huessimo/hueserver"
	"github.com/orktes/huessimo/hueupnp"
	"github.com/robertkrimen/otto"
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

	vm := otto.New()

	vm.Set("exec", func(call otto.FunctionCall) otto.Value {
		cmd := call.Argument(0).String()

		fmt.Printf("[JS][SH]: %s\n", cmd)

		out, err := exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			panic(err)
		}

		val, _ := otto.ToValue(string(out))
		return val
	})

	vm.Set("print", func(call otto.FunctionCall) otto.Value {
		fmt.Printf("[JS] %s.\n", call.Argument(0).String())
		return otto.Value{}
	})

	script, err := ioutil.ReadFile(*scriptSrcPtr)
	if err != nil {
		panic(err)
	}

	_, err = vm.Run(string(script))
	if err != nil {
		panic(err)
	}

	getLights := func() hueserver.LightList {
		value, err := vm.Run(`JSON.stringify(getLights());`)
		if err != nil {
			panic(err)
		}

		str, err := value.ToString()
		if err != nil {
			panic(err)
		}

		list := &hueserver.LightList{}
		json.Unmarshal([]byte(str), list)
		return *list
	}

	getLight := func(id string) hueserver.Light {
		arg, err := json.Marshal(id)
		if err != nil {
			panic(err)
		}
		value, err := vm.Run(`JSON.stringify(getLight(` + string(arg) + `));`)
		if err != nil {
			panic(err)
		}

		str, err := value.ToString()
		if err != nil {
			panic(err)
		}

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

		value, err := vm.Run(`JSON.stringify(setLightState(` + string(arg) + `, ` + string(arg2) + `));`)
		if err != nil {
			panic(err)
		}

		str, err := value.ToString()
		if err != nil {
			panic(err)
		}

		resp := &hueserver.LightStateChangeResponse{}
		json.Unmarshal([]byte(str), resp)
		return *resp
	}

	go hueupnp.CreateUPNPResponder("http://"+*ipPtr+":"+*portPtr+"/upnp/setup.xml", *uuidPtr, *upnpPortPtr)

	srv := hueserver.NewServer(*uuidPtr, *ipPtr+":"+*portPtr, *namePtr, getLights, getLight, setLightState)
	srv.Start(":" + *portPtr)
}
