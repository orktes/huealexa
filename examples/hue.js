var registry = require('registry');
var waitFor = require('http').waitFor;
var HueBridge = require('devices/lighting/hue').HueBridge;
var _ = require('lodash');

// Wait for server to be open
waitFor('10.0.1.3', 80);

// create bridge client
var hue = new HueBridge('http://10.0.1.3/api/3F430DA686');
// query lights from bridge
var lights = hue.getLights();

// Add lights to registry
_.each(lights, function (light) {
  registry.addDevice(light);
});

// Also add Viera TV
require('devices/tv/viera').discoverDevices(function (devices) {
  if (devices.length > 0) {
    registry.addDevice(devices[0].createDevice('TV'));
    registry.addDevice(devices[0].createDevice('Apple TV', 'NRC_CHG_INPUT-ONOFF', 'NRC_TV-ONOFF'));
  }
});
