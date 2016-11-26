var registry = require('registry');
var exec = require('process').exec;

registry.addDevice(registry.createDevice('Testing', registry.DimmableLamp, function (state) {
  console.log('State changed - ' + JSON.stringify(state));
}));

registry.addDevice(registry.createDevice('Second', registry.DimmableLamp, function (state) {
  console.log('State changed - ' + JSON.stringify(state));
}));

// Try out a bash command
console.log(exec("echo Script loaded"));
