package hueupnp

import "github.com/king-jam/gossdp"

/**
stolen from amazon-echo-ha-bridge

String discoveryTemplate = "HTTP/1.1 200 OK\r\n" +
			"CACHE-CONTROL: max-age=86400\r\n" +
			"EXT:\r\n" +
			"LOCATION: http://%s:%s/upnp/%s/setup.xml\r\n" +
			"OPT: \"http://schemas.upnp.org/upnp/1/0/\"; ns=01\r\n" +
			"01-NLS: %s\r\n" +
			"ST: urn:schemas-upnp-org:device:basic:1\r\n" +
			"USN: uuid:Socket-1_0-221438K0100073::urn:Belkin:device:**\r\n\r\n";
**/
/*
var upnpTemplate = template.Must(template.New("upnp").Parse(`HTTP/1.1 200 OK
CACHE-CONTROL: max-age=86400
EXT:
LOCATION: {{.location}}
OPT: "http://schemas.upnp.org/upnp/1/0/"; ns=01
ST: urn:schemas-upnp-org:device:basic:1
USN: uuid:{{.uuid}}::urn:Belkin:device:**

`))
*/
/*
HTTP/1.1 200 OK\r\n\
CACHE-CONTROL: max-age=100\r\n\
EXT:\r\n\
LOCATION: http://" + ip + ":" + port + "/description.xml\r\n\
SERVER: FreeRTOS/6.0.5, UPnP/1.0, IpBridge/0.1\r\n\
ST: upnp:rootdevice\r\n\
USN: uuid:2fa00080-d000-11e1-9b23-001f80007bbe::upnp:rootdevice\r\n
*/

// CreateUPNPResponder takes in the setupLocation http://[IP]:[POST]/upnp/setup.xml
func CreateUPNPResponder(setupLocation string, uuid string) *gossdp.Ssdp {
	s, err := gossdp.NewSsdp(nil)
	if err != nil {
		panic(err)
	}

	belkinDef := gossdp.AdvertisableServer{
		ServiceType: "urn:schemas-upnp-org:device:basic:1",
		DeviceUuid:  "uuid:" + uuid + "::urn:Belkin:device:**",
		Location:    setupLocation,
		MaxAge:      86400,
	}

	s.AdvertiseServer(belkinDef)

	hueDef := gossdp.AdvertisableServer{

		ServiceType: "upnp:rootdevice",
		DeviceUuid:  "uuid:" + uuid + "::upnp:rootdevice",
		Location:    setupLocation,
		MaxAge:      86400,
	}

	s.AdvertiseServer(hueDef)

	return s
}
