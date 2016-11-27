var process = require('process');

console.log('Here');
process.execAsync('sleep 5s; echo goo', function (error, response) {
  console.log('Here2');
  console.log('returned', response);
});
console.log('Here1');
process.exec('sleep 1s');
