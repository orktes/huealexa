var devices = {};
var deviceId = 1;

exports._getLights = function getLights() {
  return devices;
}

exports._getLight = function getLight(id) {
  return devices[id];
}

exports._setLightState = function setLightState(id, state) {
  var response = devices[id].setState(state);
  if (response) {
    return response;
  }
  var success = {};
  for (var key in state) {
    success["/lights/" + id + "/state/" + key] = state[key];
  }
  return [{success: success}];
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
    setState: function (state) {
      setStateCallback(state, device);
      for (var key in state) {
        device.state[key] = state[key];
      }
    }
  };
};

exports.addDevice = function (device) {
  var id = deviceId++;
  devices[id] = device;
  return id;
};

exports.removeDevice = function (device) {
  for (var id in devices) {
    if (device[id] === device) {
      delete device[id];
      return;
    }
  }
};
