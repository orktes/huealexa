var waitFor = require('http').waitFor;
var Deconz = require('devices/lighting/deconz');
var _ = require('lodash');


// Wait for server to be open
waitFor('10.0.1.3', 80);

var TIMEOUT = 1000 * 60 * 3;

// create bridge client
var deconz = new Deconz('http://10.0.1.3/api/3F430DA686');

var sensorLights = {
  "11": {
    lights: {
      "8": {
        presence: false,
        timeout: null,
        timeoutHandler: null,
        target: {
          bri: 80,
        }
      },
      "10": {
        presence: false,
        timeout: null,
        timeoutHandler: null,
        target: {
          bri: 80,
        }
      },
    },
  }
};

function handleLight(id, state) {
  if (state.presence) {
    clearTimeout(state.timeout);
    state.timeout = setTimeout(state.timeoutHandler, TIMEOUT);
    console.log('Was present');
    return;
  }

  state.presence = true;

  deconz.getLight(id, function (lightBeforeSet) {
    var stateBeforeSet = _.clone(lightBeforeSet.light.state);
    lightBeforeSet.setState(_.extend({}, stateBeforeSet, {
      on: true,
      bri: Math.max(stateBeforeSet.bri, state.target.bri)
    }), function () {});

    state.timeoutHandler = function () {
      state.presence = false;
      state.timeoutHandler = null;
      state.timeout = null;
      deconz.getLight(id, function (lightAfter) {
        if (lightAfter.light.state.bri !== lightBeforeSet.light.state.bri) {
          console.log('State has changed', lightAfter.light.state.bri);
          return;
        }

        lightAfter.setState(stateBeforeSet, function () {});
      });
    };
    state.timeout = setTimeout(state.timeoutHandler, TIMEOUT);
  });

}

function handlePresence(id) {
  _.each(sensorLights[id].lights, function (state, id) {
    handleLight(id, state);
  });
}

deconz.on('sensors_changed', function (id, state) {
  switch (id) {
    case "9":
      state.presence = state.buttonevent == 5002;
    case "11":
    case "8":
      if (state.presence) {
        handlePresence("11");
      }
    break;

  }
});

console.log('Script loaded');
