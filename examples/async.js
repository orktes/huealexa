var process = require('process');

require('http').requestAsync('http://google.fi', null, null, null, function (err, response) {
  console.log(err, response);
});
process.execAsync('echo goo', function (error, response) {
  console.log('Here2');
  console.log('returned', response);
});
console.log('Here1');
process.exec('sleep 1s');


console.log('Before http response')
