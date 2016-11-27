var _ = require('lodash');

var cbs = {};
var ids = 0;

exports._native_callback = function (id, err, value, remove) {
  var cb = cbs[id];
  cb(err, value);
  if (remove) {
    delete cbs[id];
  }
};

exports.createJSCallback = function (id, stringify) {
  return function () {
    var args = Array.prototype.slice.call(arguments);
    if (stringify) {
      args = _.map(args, function (arg) {
        return JSON.stringify(arg);
      });
    }
    _native_async_response.apply(null, [id].concat(args));
  };
};

exports.clear = function (id) {
  delete cbs[id];
};

exports.createNativeCallback = function (cb) {
  var id = ids++;
  cbs[id] = cb;
  return id;
};

exports.createAsyncFunction = function (fn) {
  return function () {
    var args = Array.prototype.slice.call(arguments);
    var cb = args[args.length - 1];
    if (typeof cb !== 'function') {
      throw new Error('Last argument should be a function')
    }
    args[args.length - 1] = exports.createNativeCallback(cb);
    fn.apply(null, args);
  };
};
