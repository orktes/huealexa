module.exports = function (mac, broadcastAddr, broadcastPort, bIface) {
    broadcastAddr = broadcastAddr || "255.255.255.255";
    broadcastPort = broadcastPort || 9;
    bIface = bIface || "";

    return _wol(mac, broadcastAddr, broadcastPort, bIface);
};