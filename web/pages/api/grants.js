import { execFile } from 'child_process';
import path from 'path';

const CLI_MAIN = path.join(process.cwd(), '..', 'cmd', 'synnergy', 'main.go');

function runCLI(args) {
  return new Promise((resolve, reject) => {
    execFile('go', ['run', CLI_MAIN, ...args], { timeout: 15000 }, (err, stdout, stderr) => {
      if (err) {
        reject(new Error(stderr || err.message));
        return;
      }
      resolve(stdout);
    });
  });
}

function firstNonGasLine(output) {
  for (const line of output.split('\n')) {
    const trimmed = line.trim();
    if (!trimmed || trimmed.startsWith('gas cost:')) {
      continue;
    }
    return trimmed;
  }
  return output.trim();
}

function parseJSON(output) {
  const idxObj = output.indexOf('{');
  const idxArr = output.indexOf('[');
  let start = -1;
  if (idxObj !== -1 && idxArr !== -1) {
    start = Math.min(idxObj, idxArr);
  } else {
    start = Math.max(idxObj, idxArr);
  }
  if (start === -1) {
    return null;
  }
  const payload = output.slice(start).trim();
  try {
    return JSON.parse(payload);
  } catch (err) {
    return null;
  }
}

async function listGrants() {
  const [listRaw, statusRaw] = await Promise.all([
    runCLI(['syn3800', 'list']),
    runCLI(['syn3800', 'status']),
  ]);
  return {
    grants: parseJSON(listRaw) || [],
    status: parseJSON(statusRaw) || {},
  };
}

async function fetchGrant(id) {
  const [grantRaw, auditRaw] = await Promise.all([
    runCLI(['syn3800', 'get', String(id)]),
    runCLI(['syn3800', 'audit', String(id)]),
  ]);
  return {
    grant: parseJSON(grantRaw) || null,
    audit: parseJSON(auditRaw) || [],
  };
}

async function handleCreate(body) {
  const { beneficiary, name, amount, authorizers = [], wallet, password } = body || {};
  if (!beneficiary || !name || !amount) {
    throw new Error('beneficiary, name and amount are required');
  }
  if (!wallet || !password) {
    throw new Error('wallet and password are required');
  }
  const args = ['syn3800', 'create', beneficiary, name, String(amount), '--wallet', wallet, '--password', password];
  for (const auth of authorizers) {
    if (auth && auth.length) {
      args.push('--authorizer', auth);
    }
  }
  const output = await runCLI(args);
  return { message: firstNonGasLine(output) };
}

async function handleRelease(body) {
  const { id, amount, note = '', wallet, password } = body || {};
  if (!id || !amount || !wallet || !password) {
    throw new Error('id, amount, wallet and password are required');
  }
  const args = ['syn3800', 'release', String(id), String(amount)];
  if (note) {
    args.push(note);
  }
  args.push('--wallet', wallet, '--password', password);
  const output = await runCLI(args);
  return { message: firstNonGasLine(output) };
}

async function handleAuthorize(body) {
  const { id, wallet, password } = body || {};
  if (!id || !wallet || !password) {
    throw new Error('id, wallet and password are required');
  }
  const args = ['syn3800', 'authorize', String(id), '--wallet', wallet, '--password', password];
  const output = await runCLI(args);
  return { message: firstNonGasLine(output) };
}

export default async function handler(req, res) {
  try {
    if (req.method === 'GET') {
      const { id } = req.query || {};
      if (id) {
        const data = await fetchGrant(id);
        res.status(200).json(data);
        return;
      }
      const data = await listGrants();
      res.status(200).json(data);
      return;
    }
    if (req.method === 'POST') {
      const { action } = req.body || {};
      let result;
      switch (action) {
        case 'create':
          result = await handleCreate(req.body);
          break;
        case 'release':
          result = await handleRelease(req.body);
          break;
        case 'authorize':
          result = await handleAuthorize(req.body);
          break;
        default:
          res.status(400).json({ error: 'unknown action' });
          return;
      }
      const refreshed = await listGrants();
      res.status(200).json({ ...result, ...refreshed });
      return;
    }
    res.status(405).json({ error: 'Method not allowed' });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}
