var registry = require('registry');
var waitFor = require('http').waitFor;
var HueBridge = require('devices/lighting/hue').HueBridge;
var _ = require('lodash');

// Wait for server to be open
waitFor('10.0.1.3', 80);

// create bridge client
var hue = new HueBridge('http://10.0.1.3/api/3F430DA686');
// query lights from bridge
var lights = hue.getLights(function (lights) {
  console.log('Lights returned');
  // Add lights to registry
  _.each(lights, function (light) {
    registry.addDevice(light);
  });
});

/*
hue.request('/groups', null, null, null, function (err, groups) {
  _.each(groups, function (group, id) {
    registry.addDevice(
      registry.createDevice(group.name, registry.DimmableLamp, function (state, cb) {
        hue.request('/groups/' + id + '/action', 'PUT', state, null, function () {
          cb();
        });
      })
    );
  });
});
*/
console.log('Script loaded');
