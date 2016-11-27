var registry = require('registry');
require('devices/tv/viera').discoverDevices(function (devices) {
  if (devices.length > 0) {
    console.log('Found a television')
    registry.addDevice(devices[0].createDevice('TV'));
    registry.addDevice(devices[0].createDevice('Apple TV', 'NRC_CHG_INPUT-ONOFF', 'NRC_TV-ONOFF'));
  }
});
