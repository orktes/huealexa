var cbs = {};
var ids = 0;

function createTimer(cb, timeout, c) {
  var id = ids++;
  cbs[id] = function () {
    cb();
    c--;
    if (c === 0) {
      delete cbs[id];
    } else {
      _native_set_timeout(id, timeout);
    }
  };

  _native_set_timeout(id, timeout);

  return id;
};

exports.setTimeout = function (cb, timeout) {
  return createTimer(cb, timeout, 1);
};

exports.setInterval = function (cb, interval) {
  return createTimer(cb, timeout, -1);
};

exports.clear = function (id) {
  delete cbs[id];
};

exports._native_callback = function (id) {
  var cb = cbs[id];
  cb && cb();
};
