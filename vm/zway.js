var EventEmitter = require('eventemitter');
var _ = require('lodash');

var zWayId = 0;
var gateways = {};

var ZWayDevice = function (data) {
  _.extend(this, data);
};

ZWayDevice.prototype = _.extend({}, EventEmitter.prototype, {
  _setLevel: function (value) {
    this.metrics.level = value;
    this.emit('level_change', value);
  }
});

var ZWay = function (config) {
  this.id = zWayId++;
  this.devices = {};

  gateways[this.id] = this;

  _init_zway(this.id, JSON.stringify({
    Hostname: config.hostname,
    Port: config.port,
    PollTimeout: (config.poll_timeout || 2) * 1000000000, // Convert seconds to nanoseconds
    Username: config.username,
    Password: config.password
  }));
};

ZWay.prototype =  _.extend(EventEmitter.prototype, {
  _addDevice: function (device) {
    device = new ZWayDevice(device)
    this.devices[device.id] = device;
    this.emit('device_added', device);
  },
});

module.exports.ZWay = ZWay;
module.exports._device_added = function (id, device) {
  gateways[id]._addDevice(device);
};

module.exports._value_change = function (id, deviceId, value) {
  gateways[id].devices[deviceId]._setLevel(value);
};
