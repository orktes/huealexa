exports.discoverDevices = function (search, callback) {
  callback(_native_ssdp_discover_devices(search))
};
