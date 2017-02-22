var registry = require('registry');
var exec = require('process').exec;
var server = require('server');

registry.addDevice(registry.createDevice('Testing', registry.DimmableLamp, function (state, cb) {
  console.log('State changed - ' + JSON.stringify(state));
  return cb();
}));

registry.addDevice(registry.createDevice('Second', registry.DimmableLamp, function (state, cb) {
  console.log('State changed - ' + JSON.stringify(state));
  return cb();
}));


server.get('/foobar', function (req, res) {
  res.headers["content-type"] = "text/plain";
  res.end("You can create custom http routes also");
});
