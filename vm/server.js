var _ = require('lodash');

var handlers = {};

function Response(callback) {
  this.headers = [];
  this.body = "";
  this.status = 200;
  this.callback = callback;
}

Response.prototype = {
  write: function (data) {
    this.body += (data || "");
  },
  end: function (data) {
    data && this.write(data);
    this.callback({
      body: this.body,
      headers: this.headers,
      status_code: this.status
    });
  },
  setHeader: function (key, value) {
    this.headers.push({key: key, value: value});
  }
};

function Request(method, data) {
  this.headers = _.map(data.headers, function (value, key) {
    return {
      value: _.last(value),
      key: key
    };
  });

  this.base = data.base;
  this.id = data.id;
  this.method = method;
  this.path = data.path;
  this.query = _.map(data.query, function (value, key) {
    return {
      value: _.last(value),
      key: key
    };
  });
}

Request.prototype = {
  body: function () {
    return _get_server_req_body(this.id);
  },
  json: function () {
    return JSON.parse(this.body());
  },
  form: function () {
    var data = this.body();
    var result = {};
    var vars = data.split('&');
    for (var i = 0; i < vars.length; i++) {
      var pair = vars[i].split('=');
      result[decodeURIComponent(pair[0])] = decodeURIComponent(pair[1]);
    }
    return result;
  }
};

function createHandler(method) {
  return function (url, callback) {
    var id = _add_server_handler(method, url);
    handlers[id] = function (req, resCallback) {
      callback(new Request(method, req), new Response(resCallback));
    };
  }
}

module.exports = {
  _request: function (id, res, callback) {
    handlers[id](res, callback);
  },
  get: createHandler("GET"),
  post: createHandler("POST"),
  delete: createHandler("DELETE"),
  put: createHandler("PUT")
};
