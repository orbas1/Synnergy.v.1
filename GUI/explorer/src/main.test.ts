import { parseHead } from './main';
import assert from 'assert';

const info = parseHead('42 0xabc');
assert.strictEqual(info.height, 42);
assert.strictEqual(info.hash, '0xabc');
console.log('parseHead test passed');
