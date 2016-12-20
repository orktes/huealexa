
var _ = require('lodash');
var devices = {};
var deviceId = 1;
var handlers = [];

function setLightState(id, state, callback) {
  devices[id].setState(state, function (response) {
    if (response) {
      return callback(response);
    }
    var success = {};
    for (var key in state) {
      success["/lights/" + id + "/state/" + key] = state[key];
    }
    return callback([{success: success}]);
  });
}

exports._getLights = function getLights(callback) {
  return callback(devices);
}

exports._getLight = function getLight(id, callback) {
  return callback(devices[id]);
}

exports._setLightState = function (id, state, callback) {
  setLightState(id, state, function (response) {
    callback(response);
    _.each(handlers, function (handler) {
      handler.setLightState(id, response, state);
    });
  });
}

exports.createDevice = function (name, type, setStateCallback) {
  var device = {
    hascolor:true,
    name: name,
    pointsymbol:{},
    state: {
      on: true,
      bri: 0,
      reachable: true
     },
    type: "Color temperature light",
    swversion:"V1.03.07"
   };
  return {
    toJSON: function () {
      return device;
    },
    setState: function (state, callback) {
      setStateCallback(state, function (response) {
        for (var key in state) {
          device.state[key] = state[key];
        }
        callback(response);
      });

    }
  };
};

exports.addDevice = function (device) {
  console.log("[REGISTRY]: Adding device " + device.toJSON().name);
  var id = deviceId++;
  devices[id] = device;
  _.each(handlers, function (handler) {
    handler.addDevice(id, device);
  });
  return id;
};

exports.removeDevice = function (device) {
  console.log("[REGISTRY]: Removing device " + device.toJSON().name);
  for (var id in devices) {
    if (device[id] === device) {
      delete device[id];
      return;
    }
  }
};

exports.addHandler = function (handler) {
  console.log("[REGISTRY]: New handler added " + handler.toString());
  handler.addLightStateListener(function (id, state, callback) {
    setLightState(id, state, callback);
  });
  handlers.push(handler);
  _.each(devices, function (device, id) {
    handler.addDevice(id, device);
  });
};
