var ZWay = require('zway').ZWay;
var HomeKit = require('homekit').HomeKit;

var zway = new ZWay({
  hostname: '10.0.1.3',
  port: '8083',
  username: 'admin',
  password: ''
});

var homeKit = new HomeKit("12345679");

var types = {
  temperature: 'temperature_sensor',
  'door-window': 'door',
  luminosity: 'light_sensor',
};

var converters =  {
  temperature_sensor: function (val) {
    return {temp: val};
  },
  door: function (val) {
    return {on: val === 'on'};
  },
  light_sensor: function (val) {
    return {lux: val}; // TODO this is just wrong
  }
}

zway.on('device_added', function (device) {
  var type = types[device.probeType];
  if (!type) {
    return;
  }

  var converter = converters[type];

  homeKit.addDevice(device.id, {
    homekit_type: type,
    name: device.metrics.title,
    state: converter(device.metrics.level)
  });

  device.on('level_change', function (value) {
    console.log('Level change', value);
    homeKit.setDeviceState(device.id, converter(value));
  });
});
