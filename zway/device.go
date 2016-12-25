package zway

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

type Level struct {
	sync.RWMutex
	Value interface{}
	listener
}

func (l *Level) SetValue(val interface{}) {
	l.Lock()
	defer l.Unlock()

	l.Value = val
	l.emit(l)
}

func (l *Level) AddValueChangeListener(fn func(*Level)) int64 {
	return l.addListener(func(val interface{}) {
		fn(val.(*Level))
	})
}

func (l *Level) RemoveValueChangeListener(id int64) {
	l.removeListener(id)
}

func (l *Level) ToFloat() (float64, error) {
	l.Lock()
	defer l.Unlock()

	if val, ok := l.Value.(float64); ok {
		return val, nil
	}

	return 0, errors.New("Can't convert to float64")
}

func (l *Level) String() string {
	return fmt.Sprintf("%+v", l.Value)
}

func (l *Level) ToBoolean() (bool, error) {
	l.Lock()
	defer l.Unlock()

	if val, ok := l.Value.(bool); ok {
		return val, nil
	}

	switch l.String() {
	case "ok":
		return true, nil
	case "true":
		return true, nil
	case "off":
		return false, nil
	case "false":
		return false, nil
	}

	return false, errors.New("Can't convert to bool")
}

func (l *Level) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &l.Value)
}

func (l *Level) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Value)
}

type DeviceType struct {
	Name string
}

func (dt *DeviceType) UnmarshalJSON(data []byte) error {
	dt.Name = string(data[1 : len(data)-1])
	return nil
}

func (dt *DeviceType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", dt.Name)), nil
}

type ProbeType struct {
	Name string
}

func (bt *ProbeType) UnmarshalJSON(data []byte) error {
	bt.Name = string(data[1 : len(data)-1])
	return nil
}

func (bt *ProbeType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", bt.Name)), nil
}

type Metrics struct {
	ProbeTitle string `json:"probeTitle"`
	ScaleTitle string `json:"scaleTitle"`
	Title      string `json:"title"`
	Level      Level  `json:"level"`
}

type Device struct {
	sync.RWMutex
	DeviceType        DeviceType `json:"deviceType"`
	Tags              []string   `json:"tags"`
	CreatorID         int        `json:"creatorId"`
	HasHistory        bool       `json:"hasHistory"`
	ID                string     `json:"id"`
	Location          int        `json:"location"`
	Metrics           Metrics    `json:"metrics"`
	PermanentlyHidden bool       `json:"permanently_hidden"`
	ProbeType         ProbeType  `json:"probeType"`
	Visibility        bool       `json:"visibility"`
	UpdateTime        float64    `json:"updateTime"`
}

func (d *Device) Update(level Level, updateTime float64) {
	d.Lock()
	defer d.Unlock()

	if updateTime > d.UpdateTime {
		d.Metrics.Level.SetValue(level.Value)
	}
}

type Devices struct {
	sync.RWMutex
	listener
	devices map[string]*Device
}

func (devices *Devices) AddDevice(device *Device) {
	devices.Lock()
	defer devices.Unlock()

	if devices.devices == nil {
		devices.devices = map[string]*Device{}
	}

	devices.devices[device.ID] = device
	devices.emit(device)
}

func (devices *Devices) GetDevice(id string) (device *Device, found bool) {
	devices.Lock()
	defer devices.Unlock()

	if devices.devices == nil {
		return nil, false
	}

	device, found = devices.devices[id]
	return
}

func (devices *Devices) AddNewDeviceListener(fn func(*Device)) int64 {
	return devices.addListener(func(val interface{}) {
		fn(val.(*Device))
	})
}

func (devices *Devices) RemoveNewDeviceListener(id int64) {
	devices.removeListener(id)
}
