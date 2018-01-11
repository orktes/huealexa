package dra

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type handler struct {
	prefix  string
	handler func(data string)
}

// DRA struct representing a single Denon DRA amplifier
type DRA struct {
	OnUpdate chan string

	conn     net.Conn
	handlers []handler

	// state
	system struct {
		power        bool
		masterVolume int
		mute         bool
		input        string
		sleepOn      bool
		sleepValue   int
	}
}

func NewFromConn(conn net.Conn) *DRA {
	dra := &DRA{conn: conn}

	dra.initHandlers()
	dra.queryInitialValues()

	go dra.listen()

	return dra
}

func NewFromAddr(addr string) (*DRA, error) {
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	return NewFromConn(conn), nil
}

func (dra *DRA) Close() error {
	return dra.conn.Close()
}

func (dra *DRA) Send(cmd ...string) error {
	w := bufio.NewWriter(dra.conn)
	for _, p := range cmd {
		w.Write([]byte(p))
	}

	w.Write([]byte{'\r'})

	return w.Flush()
}

// SetMasterVolume sets DRA master volume value should be between 0 - 91
func (dra *DRA) SetMasterVolume(value int) error {
	return dra.Send("MV", fmt.Sprintf("%d", value))
}

// GetMasterVolume returns master volume
func (dra *DRA) GetMasterVolume() int {
	return dra.system.masterVolume
}

// SetMute sets DRA mute
func (dra *DRA) SetMute(value bool) error {
	if value {
		return dra.Send("MU", "ON")
	}
	return dra.Send("MU", "OFF")
}

// GetMute return mute
func (dra *DRA) GetMute() bool {
	return dra.system.mute
}

// SetPower sets power mode of dra (true = on, false = standby)
func (dra *DRA) SetPower(value bool) error {
	if value {
		return dra.Send("PW", "ON")
	}
	return dra.Send("PW", "STANDBY")
}

// GetPower returns power status
func (dra *DRA) GetPower() bool {
	return dra.system.power
}

func (dra *DRA) initHandlers() {
	dra.registerBoolPointer("PW", &dra.system.power, "ON")
	dra.registerIntPointer("MV", &dra.system.masterVolume)
	dra.registerBoolPointer("MU", &dra.system.mute, "ON")
	dra.registerStrPointer("SI", &dra.system.input)

	dra.registerHandler("SLP", func(value string) {
		dra.system.sleepOn = value != "OFF"
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			dra.system.sleepValue = int(intVal)
		}
	})
}

func (dra *DRA) queryInitialValues() {
	dra.Send("PW?")
	dra.Send("MV?")
	dra.Send("MU?")
	dra.Send("SI?")
	dra.Send("SLP?")
}

func (dra *DRA) listen() {
	reader := bufio.NewReader(dra.conn)
	for {
		cmd, err := reader.ReadString('\r')
		if err != nil {
			break
		}

		for _, handler := range dra.handlers {
			if strings.HasPrefix(cmd, handler.prefix) {
				handler.handler(cmd[len(handler.prefix) : len(cmd)-1])
				if dra.OnUpdate != nil {
					dra.OnUpdate <- handler.prefix
				}
				break
			}
		}
	}

	if dra.OnUpdate != nil {
		close(dra.OnUpdate)
	}
}

func (dra *DRA) registerHandler(prefix string, handlerFunc func(data string)) {
	dra.handlers = append(dra.handlers, handler{prefix, handlerFunc})
}

func (dra *DRA) registerIntPointer(prefix string, ptr *int) {
	dra.registerHandler(prefix, func(value string) {
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			*ptr = int(intVal)
		}
	})
}

func (dra *DRA) registerStrPointer(prefix string, ptr *string) {
	dra.registerHandler(prefix, func(value string) {
		*ptr = value
	})
}

func (dra *DRA) registerBoolPointer(prefix string, ptr *bool, trueValue string) {
	dra.registerHandler(prefix, func(value string) {
		*ptr = value == trueValue
	})
}
