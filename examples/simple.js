// All require paths are relative to the root dir (this will be fiexed)
var light = require('examples/util.js').getLight('Testing');

function getLights() {
  return {"1": light}
}

function getLight(id) {
  return light;
}

function setLightState(id, state) {
  var success = {};
  light.state.on = state.on;
  success["/lights/" + id + "/state/on"] = state.on;
  console.log("Light state is now " + state.on);
  return [{success: success}];
}

console.log('Script loaded');
