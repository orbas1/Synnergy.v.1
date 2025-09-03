import { parseDeploy } from './main';
import assert from 'assert';

const addr = parseDeploy('contract: 0xabc');
assert.strictEqual(addr, '0xabc');
console.log('parseDeploy test passed');
