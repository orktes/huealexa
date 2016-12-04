package main

//go:generate go-bindata -prefix "vm" -pkg vm -ignore=\.go  -o vm/assets.go vm/...

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/orktes/huessimo/hueserver"
	"github.com/orktes/huessimo/hueupnp"
	"github.com/orktes/huessimo/vm"
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
	//log.SetOutput(ioutil.Discard)
	log.Printf("Available modules %s", strings.Join(vm.AssetNames(), ", "))

	genuuid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	scriptSrcPtr := flag.String("src", "", "Script source file location")
	uuidPtr := flag.String("uuid", "", "UUID for the HUE server (for example \""+genuuid.String()+"\")")
	portPtr := flag.String("port", "8989", "Port for the HUE server")
	upnpPortPtr := flag.String("upnp", "239.255.255.250:1900", "UPNP multicast addr for the HUE server")
	namePtr := flag.String("name", "fakeServer", "Name for the HUE server")
	ipPtr := flag.String("ip", getIPAddress(), "Interface to be used")

	flag.Parse()

	if *uuidPtr == "" {
		fmt.Printf("You must provide -uuid=\"%s\" (i just generated that for you) or something else\n", genuuid.String())
		return
	}

	if *ipPtr == "" {
		fmt.Printf(`You must provide -ip=\"\11.22.33.44\"\n`)
		return
	}

	runtime, err := vm.NewVM(*scriptSrcPtr)

	if err != nil {
		panic(err)
	}

	getLights := func() hueserver.LightList {
		callbackID, cbCh := vm.CreateAsyncNativeCallbackChannel()

		_, err = runtime.RunString(fmt.Sprintf(`
        require('registry')._getLights(
          require('async').createJSCallback(%d, true)
        );
    `, callbackID))
		if err != nil {
			panic(err)
		}

		str := (<-cbCh).Argument(0).String()

		list := &hueserver.LightList{}
		json.Unmarshal([]byte(str), list)
		return *list
	}

	getLight := func(id string) hueserver.Light {
		arg, merr := json.Marshal(id)
		if merr != nil {
			panic(merr)
		}
		callbackID, cbCh := vm.CreateAsyncNativeCallbackChannel()
		_, err = runtime.RunString(fmt.Sprintf(`
      require('registry')._getLight(
        %s,
        require('async').createJSCallback(%d, true)
      );
    `, string(arg), callbackID))
		if err != nil {
			panic(err)
		}

		str := (<-cbCh).Argument(0).String()

		light := &hueserver.Light{}
		json.Unmarshal([]byte(str), light)
		return *light
	}

	setLightState := func(id string, state hueserver.LightStateChange) hueserver.LightStateChangeResponse {
		arg, merr := json.Marshal(id)
		if merr != nil {
			panic(merr)
		}

		arg2, merr := json.Marshal(state)
		if merr != nil {
			panic(merr)
		}

		callbackID, cbCh := vm.CreateAsyncNativeCallbackChannel()
		_, err = runtime.RunString(fmt.Sprintf(`
      require('registry')._setLightState(
        %s,
        %s,
        require('async').createJSCallback(%d, true)
      );
    `,
			string(arg),
			string(arg2),
			callbackID,
		))
		if err != nil {
			panic(err)
		}

		str := (<-cbCh).Argument(0).String()

		resp := &hueserver.LightStateChangeResponse{}
		json.Unmarshal([]byte(str), resp)
		return *resp
	}

	go hueupnp.CreateUPNPResponder("http://"+*ipPtr+":"+*portPtr+"/upnp/setup.xml", *uuidPtr, *upnpPortPtr)

	srv := hueserver.NewServer(*uuidPtr, *ipPtr+":"+*portPtr, *namePtr, getLights, getLight, setLightState)
	err = srv.Start(":" + *portPtr)
	if err != nil {
		panic(err)
	}

}
