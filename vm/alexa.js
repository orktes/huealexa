var path = require('path');
var fs = require('fs');
var http = require('http');
var _ = require('lodash');
var server = require('server');

var fullCredentialPath = path.join(env.data_dir, '.avs_credentials');

module.exports = {
  _renew: function () {
    var updated = new Date(this._credentials.updated);
    var timeLeft = updated.getTime() + (this._credentials.expires_in * 1000) - Date.now();
    if (timeLeft < 1000 * 60 * 5) {
      var response = http.request(this._credentials.auth_url, "POST", {
        refresh_token: this._credentials.refresh_token
      });

      try {
        if (response.access_token) {
          response.updated = new Date();
          this._setCredentials(_.extend({}, this._credentials, response));
        }
      } catch(e) {
        console.log('[ALEXA]: error: ' + e.message);
      }

    }
  },
  _persistCredentials: function () {
    fs.write(fullCredentialPath, JSON.stringify(this._credentials), 0644);
  },
  _setCredentials: function (data, dontSave) {
    this._credentials = data;
    if (!dontSave) {
      this._persistCredentials()
    }
    this._renew();
    console.log("[ALEXA]: Initialized");
  }
};

server.get('/alexa/auth', function (req, res) {
  var url = 'https://huealexaauth.herokuapp.com/?uuid=' + env.huealexa_uuid +  '&redirect=' + encodeURIComponent('http://' + req.base  + req.path);
  res.setHeader("Location", url);
  res.status = 302;
  res.end();
});

server.post('/alexa/auth', function (req, res) {
  var data = req.form();
  data.updated = new Date();
  res.end("All OK");
  module.exports._setCredentials(data);
});

try {
  var credentialStr = fs.read(fullCredentialPath);
  module.exports._setCredentials(JSON.parse(credentialStr), true);
} catch(e) {
  console.log(e);
  console.log("[ALEXA]: Not authenticated: open /avs/auth in your browser");
}
