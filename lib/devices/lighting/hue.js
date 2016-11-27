var _ = require('lodash');
var request = require('http').request;

function HueLight(light, bridge) {
  this.light = light;
  this.bridge = bridge;
}

HueLight.prototype = {
  toJSON: function () {
    return this.light;
  },
  setState: function (state) {
    var self = this;
    var result = this.bridge.request(
      '/lights/' + this.light.id + '/state',
      'PUT',
      state
    );
    _.each(state, function (value, key) {
      self.light.state[key] = value;
    });
    return result;
  }
};

function HueBridge(api) {
  this.api = api;
}

HueBridge.prototype = {
  request: function (path, method, data, headers) {
    return request(this.api + path, method, data, headers);
  },
  getLights: function () {
    var self = this;
    var response = this.request('/lights');
    return _.map(response, function (light, id) {
      light.id = id;
      return new HueLight(light, self);
    });
  }
}

exports.HueBridge = HueBridge;
