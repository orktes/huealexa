[![Build Status](https://travis-ci.org/orktes/huessimo.svg?branch=master)](https://travis-ci.org/orktes/huessimo)

# huessimo

Expose any network connected device to Amazon Echo or Apple HomeKit.

```bash
go get -u github.com/orktes/huessimo
# OR Download binary from releases
huessimo -uuid="ac103f83-e6e9-41b8-6ae5-1ef6cbe0a021" -ip=10.0.1.4 -src examples/simple.js
# Lights provider by script/src are now available in Echo Smart Home
```

## License
Huessimo: See LICENSE file

hc (HomeKit go library) https://github.com/brutella/hc
