var _ = require('lodash');

var webSocketId = 0;
var webSocketMap = {};

var WebSocket = module.exports = function (url) {
  var id = webSocketId++;
  webSocketMap[id] = this;
  _init_websocket(id+"", url);
};

function addNativeApis(Class, apis) {
  _.each(apis, function (fn, key) {
    Class.prototype[key] = fn;
    Class[key] = function (id) {
      var args = Array.prototype.slice.call(arguments);
      var ws = webSocketMap[id];
      args.shift();
      ws[key].apply(ws, args);
    };
  });
}

addNativeApis(WebSocket, {
  _message: function (msg) {
    if (this.onmessage) {
      this.onmessage({data: msg});
    }
  },

  _error: function (err) {
    if (this.onerror) {
      this.onerror(new Error(err));
    }
  },

  _connect: function () {
    if (this.onconnect) {
      this.onconnect();
    }
  },

  _close: function () {
    if (this.onclose) {
      this.onclose();
    }
  }
});
