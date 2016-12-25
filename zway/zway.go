package zway

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const dataAPI = "/ZWaveAPI/Data/"
const autoAPI = "/ZAutomation/api/v1/"
const runAPI = "/ZWaveAPI/Run/"
const cookieName = "ZWAYSession"

var ErrorCookieNotFound = errors.New("Cookie not found")

type DevicesResponse struct {
	Data struct {
		Devices          []*Device `json:"devices"`
		StructureChanged bool      `json:"structureChanged"`
		UpdateTime       float64   `json:"updateTime"`
	} `json:"data"`
}

type ZWay struct {
	sync.RWMutex
	Config      Config
	Cookie      *http.Cookie
	Devices     Devices
	updateTime  int
	stopChannel chan bool
}

func (zWay *ZWay) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, fmt.Sprintf("http://%s:%s%s", zWay.Config.Hostname, zWay.Config.Port, path), body)
}

func (zWay *ZWay) getDevicesUpdates() (err error) {
	url := fmt.Sprintf("%sdevices", autoAPI)
	if zWay.updateTime > 0 {
		url = fmt.Sprintf("%s?since=%d", url, zWay.updateTime)
	}
	req, err := zWay.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.AddCookie(zWay.Cookie)
	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return
	}
	defer rsp.Body.Close()
	decoder := json.NewDecoder(rsp.Body)
	resp := &DevicesResponse{}
	err = decoder.Decode(&resp)
	if err != nil {
		return
	}

	if resp.Data.UpdateTime > 0 {
		zWay.Lock()
		zWay.updateTime = int(resp.Data.UpdateTime)
		zWay.Unlock()
	}

	if len(resp.Data.Devices) > 0 {
		for _, device := range resp.Data.Devices {
			existingDevice, ok := zWay.Devices.GetDevice(device.ID)
			if ok {
				existingDevice.Update(device.Metrics.Level, device.UpdateTime)
			} else {
				zWay.Devices.AddDevice(device)
			}
		}
	}

	return
}

func (zWay *ZWay) Poll(interval time.Duration) {
	for {
		err := zWay.getDevicesUpdates()
		if err != nil {
			log.Printf("[ZWAY]: %s\n", err.Error())
		}
		select {
		case <-time.After(interval):
			continue
		case <-zWay.stopChannel:
			return
		}
	}
}

func (zWay *ZWay) Auth() (err error) {
	url := fmt.Sprintf("%slogin", autoAPI)
	login := fmt.Sprintf("{\"login\": \"%s\", \"password\": \"%s\"}",
		zWay.Config.Username, zWay.Config.Password)
	req, err := zWay.NewRequest("POST", url, strings.NewReader(login))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return
	}

	cookies := rsp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == cookieName && cookie.Path == "/" {
			zWay.Cookie = cookie
			return
		}
	}

	return ErrorCookieNotFound
}

func (zWay *ZWay) Stop() {
	zWay.stopChannel <- true
}

func New(config Config) (zWay *ZWay, err error) {
	zWay = &ZWay{Config: config, stopChannel: make(chan bool)}
	err = zWay.Auth()
	if err != nil {
		return
	}
	go zWay.Poll(config.PollTimeout)
	return
}
