package hueupnp

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"text/template"

	"golang.org/x/net/ipv4"
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

func createSocket() (*ipv4.PacketConn, net.PacketConn, error) {
	group := net.IPv4(239, 255, 255, 250)
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("net.Interfaces error: %s", err)
		return nil, nil, err
	}
	con, err := net.ListenPacket("udp4", "0.0.0.0:1900")
	if err != nil {
		log.Fatalf("net.ListenPacket error: %s", err)
		return nil, nil, err
	}
	p := ipv4.NewPacketConn(con)
	p.SetMulticastLoopback(true)
	didFindInterface := false
	for i, v := range interfaces {
		ef, err := v.Addrs()
		if err != nil {
			continue
		}
		hasRealAddress := false
		for k := range ef {
			asIp := net.ParseIP(ef[k].String())
			if asIp.IsUnspecified() {
				continue
			}
			hasRealAddress = true
			break
		}
		if !hasRealAddress {
			continue
		}
		err = p.JoinGroup(&v, &net.UDPAddr{IP: group})
		if err != nil {
			log.Printf("join group %d %s", i, err)
			continue
		}
		didFindInterface = true
	}
	if !didFindInterface {
		return nil, nil, errors.New("Unable to find a compatible network interface!")
	}

	return p, con, nil
}

// CreateUPNPResponder takes in the setupLocation http://[IP]:[POST]/upnp/setup.xml
func CreateUPNPResponder(setupLocation string, uuid string, upnpAddr string) {
	sock, rawCon, err := createSocket()
	if err != nil {
		panic(err)
	}

	defer sock.Close()
	defer rawCon.Close()

	for {
		b := make([]byte, 2048)
		n, src, err := rawCon.ReadFrom(b)
		if err != nil {
			log.Fatal("[UPNP] ReadFromUDP failed:", err)
		}

		if strings.HasPrefix(string(b[:n]), "M-SEARCH * HTTP/1.1") && strings.Contains(string(b[:n]), "MAN: \"ssdp:discover\"") {
			addr, err := net.ResolveUDPAddr("udp4", src.String())
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
			rawCon.WriteTo([]byte(b.String()), addr)
		}
	}

}
