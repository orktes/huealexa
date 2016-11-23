package hueupnp

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
	"text/template"
)

const (
	upnpMulticastAddress = "239.255.255.250:1900"
)

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

var upnpTemplate = template.Must(template.New("upnp").Parse(`HTTP/1.1 200 OK
CACHE-CONTROL: max-age=86400
EXT:
LOCATION: {{.location}}
OPT: "http://schemas.upnp.org/upnp/1/0/"; ns=01
ST: urn:schemas-upnp-org:device:basic:1
USN: uuid:{{.uuid}}::urn:Belkin:device:**

`))

// CreateUPNPResponder takes in the setupLocation http://[IP]:[POST]/upnp/setup.xml
func CreateUPNPResponder(setupLocation string, uuid string) {
	addr, err := net.ResolveUDPAddr("udp", upnpMulticastAddress)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("[UPNP] ListenMulticastUDP failed:", err)
	}

	l.SetReadBuffer(1024)

	for {
		b := make([]byte, 1024)
		n, src, err := l.ReadFromUDP(b)
		if err != nil {
			log.Fatal("[UPNP] ReadFromUDP failed:", err)
		}

		if strings.HasPrefix(string(b[:n]), "M-SEARCH * HTTP/1.1") && strings.Contains(string(b[:n]), "MAN: \"ssdp:discover\"") {
			c, err := net.DialUDP("udp", nil, src)
			if err != nil {
				log.Fatal("[UPNP] DialUDP failed:", err)
			}

			log.Println("[UPNP] discovery request from", src)

			b := &bytes.Buffer{}
			err = upnpTemplate.Execute(b, map[string]string{"location": setupLocation, "uuid": uuid})
			if err != nil {
				log.Fatal("[UPNP] execute template failed:", err)
			}
			fmt.Printf("[UPNP] Sending\n%s\nto %s\n", b.Bytes(), src)
			c.Write(b.Bytes())
		}
	}
}
