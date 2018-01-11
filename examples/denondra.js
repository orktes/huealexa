var DenonDRA = require('devices/audio/denon_dra');


var dra = new DenonDRA("10.0.1.8:23");
dra.onconnect = function () {
    console.log("Set master volume");
    dra.setMasterVolume(70);
}

dra.onerror = function () {
    console.log("error");
};

dra.onclose = function () {
    console.log("closed");
}

dra.onupdate = console.log.bind(console);