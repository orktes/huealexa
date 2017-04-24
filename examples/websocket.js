var WebSocket = require('websocket');

var host = '10.0.1.3';
var port = 8080;

var ws = new WebSocket('ws://' + host + ':' + port);
ws.onmessage = function(msg) {
    console.log(msg.data);
}
