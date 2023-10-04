function chooseBinary() {
        return `updateNetwork`
}
var  child_process   = require('node:child_process');
const binary = chooseBinary();
const mainScript = `${__dirname}/${binary}`;
var spawnSyncReturns = child_process.spawnSync(mainScript, { stdio: 'inherit' });
console.log('stdout:\n'+spawnSyncReturns.stdout)
