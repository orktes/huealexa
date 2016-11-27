var exec = require('process').exec;
var _ = require('lodash');

exports.waitFor = function (host, port, wait) {
  wait = wait || 1;
  while (true) {
    try {
      exec("nc -z " + host + " " + port);
      return;
    } catch(e) {
      exec("sleep " + wait + "s");
    }
  }
};

// TODO do actual http
exports.request = function(url, method, body, headers) {
  method = method || 'GET';

  var cmd = "curl '" + url + "' -X " + method;
  if (method !== 'GET' && body) {
      cmd += " -d "+JSON.stringify(typeof body === 'object' ? JSON.stringify(body) : body)
  }

  if (headers) {
    _.each(headers, function (value, key) {
      cmd += ' -H \'' + key + ': ' + value +  '\''
    });
  }

  var result = exec(cmd);
  try {
    result = JSON.parse(result);
  } catch (e) {}

  return result;
};
