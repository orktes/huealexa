var light = {
  "etag":"03e7d1d9643d16c750427155032ce2d5",
  "hascolor":true,
  "manufacturer":"OSRAM",
  "modelid":"Classic A60 TW",
  "name":"Testing",
  "pointsymbol":{},
  "state":{
    "on":true,
    "reachable":true
   },
   "swversion":"V1.03.07",
   "type":"Color temperature light",
   "uniqueid":"84:18:26:00:00:CA:56:E3-03"
 };

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
