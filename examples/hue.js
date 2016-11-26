var exec = require('process').exec;
var registry = require('registry');

registry._getLights = function () {
  // Call real hue bridge
  return JSON.parse(exec("curl 'http://10.0.1.3/api/3F430DA686/lights'"));
}

registry._getLight = function (id) {
  // Call real hue bridge
  return JSON.parse(exec("curl 'http://10.0.1.3/api/3F430DA686/lights/" + id + "'"));
}

registry._setLightState = function (id, state) {
  // Call real hue bridge
  return JSON.parse(exec("curl 'http://10.0.1.3/api/3F430DA686/lights/" + id + "/state' -X PUT --data-binary '" + JSON.stringify(state) + "'"));
}


console.log('Script loaded');
