// This example exposes HUE lights trough home kit

var registry = require('registry');
var HomeKit = require('homekit').HomeKit;
var HueBridge = require('devices/lighting/hue').HueBridge;
var _ = require('lodash');


var homeKit = new HomeKit("12345679");
registry.addHandler(homeKit);

var hue = new HueBridge('http://10.0.1.3/api/3F430DA686');
// query lights from bridge
var lights = hue.getLights(function (data) {
  // Add lights to registry
  _.each(data, function (light) {
    registry.addDevice(light);
  });
});
