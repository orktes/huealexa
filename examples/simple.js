var registry = require('registry');
var exec = require('process').exec;

registry.addDevice(registry.createDevice('Testing', registry.DimmableLamp, function (state, cb) {
  console.log('State changed - ' + JSON.stringify(state));
  return cb();
}));

registry.addDevice(registry.createDevice('Second', registry.DimmableLamp, function (state, cb) {
  console.log('State changed - ' + JSON.stringify(state));
  return cb();
}));
