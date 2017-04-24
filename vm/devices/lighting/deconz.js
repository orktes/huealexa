var _ = require('lodash');
var HueBridge = require('devices/lighting/hue').HueBridge;
var WebSocket = require('websocket');
var EventEmitter = require('eventemitter');

var Deconz = module.exports = function (api) {
  HueBridge.call(this, api);
  this._connectWebSocket();
};

_.extend(Deconz.prototype, EventEmitter.prototype);
_.extend(Deconz.prototype, HueBridge.prototype);
_.extend(Deconz.prototype, {
  _connectWebSocket: function () {
    var self = this;
    this.request('/config', null, null, null, function (err, config) {
      if (err != null) {
        return;
      }

      var port = config.websocketport;
      var host = config.ipaddress;
      var ws = new WebSocket('ws://' + host + ':' + port);
      ws.onmessage = function(msg) {
        var event = JSON.parse(msg.data);
        self.emit(event.r + '_' + event.e, event.id, event.state);
      };
      ws.onconnect = function () {
        console.log('[DECONZ] WS Connection created');
      };
    });
  }
});
