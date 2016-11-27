var ssdp = require('ssdp');
var _ = require('lodash');
var registry = require('registry');
var http = require('http');


function VieraTV(data) {
  this.data = data;
  this.host = data.Root.URLBase.Host;
}

VieraTV.prototype = {
  sendAction: function (action) {
    return this.sendRequest('command', 'X_SendKey', '<X_KeyEvent>'+action+'</X_KeyEvent>');
  },
  sendRequest: function(type, action, command, options) {
    var url, urn;
    if (typeof this.host === 'undefined') return;
    if (type === 'command') {
      url = '/nrc/control_0';
      urn = 'panasonic-com:service:p00NetworkControl:1';
    } else if (type === 'render') {
      url = '/dmr/control_0';
      urn = 'schemas-upnp-org:service:RenderingControl:1';
    }

    var body = "<?xml version='1.0' encoding='utf-8'?> \
    <s:Envelope xmlns:s='http://schemas.xmlsoap.org/soap/envelope/' s:encodingStyle='http://schemas.xmlsoap.org/soap/encoding/'> \
    <s:Body> \
    <u:" + action + " xmlns:u='urn:" + urn + "'> \
    " + command + " \
    </u:" + action + "> \
    </s:Body> \
    </s:Envelope>";

    return http.request('http://' + this.host + url, 'POST', body, {
      'SOAPACTION': '"urn:'+urn+'#'+action+'"'
    });
  },
  createDevice: function (name, onAction, offAction) {
    var vieraTv = this;
    return registry.createDevice(name, registry.DimmableLamp, function (state, cb) {
      if ('on' in state) {
        if (state.on) {
          vieraTv.sendAction(onAction || 'NRC_POWER-ONOFF');
        } else {
          vieraTv.sendAction(offAction || 'NRC_POWER-ONOFF');
        }
      }
      cb();
    });
  }
};

module.exports.discoverDevices = function (callback) {
  ssdp.discoverDevices('urn:panasonic-com:service:p00NetworkControl:1', function (devices) {
    callback(_.map(devices, function (device) {
      return new VieraTV(device);
    }));
  });
}
