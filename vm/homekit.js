var EventEmitter = require('eventemitter');
var _ = require('lodash');

var eventEmitter = new EventEmitter();

var HomeKit = function (pin) {
  this.pin = pin;
  this.states = {};
};

HomeKit.prototype = {
  addLightStateListener: function (callback) {
    var self = this;
    function transform(key, mutator) {
      return function (id, value) {
        if (!self.states[id]) {
          self.states[id] = {
            trigger: _.debounce(function () {
                var state = _.reduce(self.states[id].states, function (fullState, partial) {
                  _.each(partial, function (value, key) {
                    fullState[key] = value;
                  });
                  return fullState;
                }, {});
                delete self.states[id];
                callback(id, state, function () {
                   // NO OP For now
                 });
              }, 300),
            states: []
          }
        }

        var state = {};
        state[key] = mutator ? mutator(value) : value;
        self.states[id].states.push(state);
        self.states[id].trigger();
      };
    }

    eventEmitter.on('light_on_change', transform('on'));
    eventEmitter.on('light_bri_change', transform('bri', function (value) {
      return Math.round(255 * (value / 100));
    }));
    eventEmitter.on('light_sat_change', transform('sat', function (value) {
      return Math.round(255 * (value / 100));
    }));
    eventEmitter.on('light_hue_change', transform('hue', function (value) {
      return Math.floor(value * 182.0444);
    }));
  },
  setLightState: function (id, response, state) {
    if ('on' in state) {
      _set_homekit_device_on(id, state.on);
    }

    if ('bri' in state) {
      _set_homekit_device_bri(id, Math.round((state.bri / 255) * 100));
    }

    if ('sat' in state) {
      _set_homekit_device_sat(id, Math.round((state.sat / 255) * 100));
    }

    if ('hue' in state) {
      _set_homekit_device_hue(id, Math.round(state.hue / 182.0444));
    }
  },
  addDevice: function (id, device) {
    // Only support lights for now
    var deviceData = device.toJSON();
    var info = {
      Name: deviceData.name,
      SerialNumber: deviceData.uniqueid,
      Manufacturer: deviceData.manufacturername,
      Model: deviceData.modelid,
    };

    _add_homekit_device(id, "lightbulb", this.pin, JSON.stringify(info));
    this.setLightState(id, {}, deviceData.state);
  },
  toString: function () {
    return "HomeKit";
  }
};


module.exports = {
  _remote_on_change: eventEmitter.emit.bind(eventEmitter, 'light_on_change'),
  _remote_bri_change: eventEmitter.emit.bind(eventEmitter, 'light_bri_change'),
  _remote_sat_change: eventEmitter.emit.bind(eventEmitter, 'light_sat_change'),
  _remote_hue_change: eventEmitter.emit.bind(eventEmitter, 'light_hue_change'),
  HomeKit: HomeKit
}
