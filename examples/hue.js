function getLights() {
  // Call real hue bridge
  return JSON.parse(exec("curl 'http://10.0.1.3/api/3F430DA686/lights'"));
}

function getLight(id) {
  // Call real hue bridge
  return JSON.parse(exec("curl 'http://10.0.1.3/api/3F430DA686/lights/" + id + "'"));
}

function setLightState(id, state) {
  // Call real hue bridge
  return JSON.parse(exec("curl 'http://10.0.1.3/api/3F430DA686/lights/" + id + "/state' -X PUT --data-binary '" + JSON.stringify(state) + "'"));
}


console.log('Script loaded');
