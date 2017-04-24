var _ = require('lodash');
var requestAsync = require('http').requestAsync;

function HueLight(light, bridge) {
  this.light = light;
  this.bridge = bridge;
}

HueLight.prototype = {
  toJSON: function () {
    return this.light;
  },
  setState: function (state, cb) {
    var self = this;
    this.bridge.request(
      '/lights/' + this.light.id + '/state',
      'PUT',
      state,
      null,
      function (err, result) {
        _.each(state, function (value, key) {
          self.light.state[key] = value;
        });
        return cb(result);
      }
    );

  }
};

function HueBridge(api) {
  this.api = api;
}

HueBridge.prototype = {
  request: function (path, method, data, headers, callback) {
    return requestAsync(this.api + path, method, data, headers, callback);
  },
  getLight: function (id, callback) {
    var self = this;
    this.request('/lights/' + id, null, null, null, function (err, light) {
      light.id = id;
      callback(new HueLight(light, self));
    });
  },
  getLights: function (callback) {
    var self = this;
    this.request('/lights', null, null, null, function (err, response) {
      callback(_.map(response, function (light, id) {
        light.id = id;
        return new HueLight(light, self);
      }));
    });
  }
}

exports.HueBridge = HueBridge;
