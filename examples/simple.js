var registry = require('registry');
var exec = require('process').exec;

registry.addDevice(registry.createDevice('Testing', registry.DimmableLamp, function (state) {
  console.log('State changed - ' + JSON.stringify(state));
}));

registry.addDevice(registry.createDevice('Second', registry.DimmableLamp, function (state) {
  console.log('State changed - ' + JSON.stringify(state));
}));

// Try out a bash command
console.log(exec("echo Script loaded"));

var registry = require('registry');
require('devices/tv/viera').discoverDevices(function (devices) {
  if (devices.length > 0) {
    registry.addDevice(devices[0].createDevice('TV'));
    registry.addDevice(devices[0].createDevice('Apple TV', 'NRC_CHG_INPUT-ONOFF', 'NRC_TV-ONOFF'));
  }
});
