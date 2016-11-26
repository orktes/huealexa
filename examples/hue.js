var request = require('http').request;
var registry = require('registry');

registry._getLights = function () {
  // Call real hue bridge
  return request("http://10.0.1.3/api/3F430DA686/lights");
}

registry._getLight = function (id) {
  // Call real hue bridge
  return request("http://10.0.1.3/api/3F430DA686/lights/" + id);
}

registry._setLightState = function (id, state) {
  // Call real hue bridge
  return request("http://10.0.1.3/api/3F430DA686/lights/" + id + "/state", 'PUT', state);
}


console.log('Script loaded');
