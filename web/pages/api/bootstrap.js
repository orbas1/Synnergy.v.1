import { execFile } from 'child_process';
import path from 'path';

export default function handler(req, res) {
  if (req.method !== 'POST') {
    res.status(405).json({ error: 'Method not allowed' });
    return;
  }
  const {
    nodeId = 'web-node',
    address = '',
    consensus = 'Synnergy-PBFT',
    governance = 'SYN-Gov',
    replicate = true,
    regulator = true,
    authorities = [],
  } = req.body || {};

  const cliPath = path.join(process.cwd(), '..', 'cmd', 'synnergy', 'main.go');
  const args = ['run', cliPath, '--json', 'orchestrator', 'bootstrap'];
  if (nodeId) {
    args.push('--node-id', String(nodeId));
  }
  if (address) {
    args.push('--address', String(address));
  }
  if (consensus) {
    args.push('--consensus', String(consensus));
  }
  if (governance) {
    args.push('--governance', String(governance));
  }
  if (!replicate) {
    args.push('--replicate=false');
  }
  if (!regulator) {
    args.push('--regulator=false');
  }
  (Array.isArray(authorities) ? authorities : [authorities]).forEach((entry) => {
    if (entry && entry.length) {
      args.push('--authority', String(entry));
    }
  });

  execFile('go', args, { timeout: 20000 }, (err, stdout, stderr) => {
    if (err) {
      res.status(500).json({ error: stderr || err.message });
      return;
    }
    try {
      const parsed = JSON.parse(stdout);
      res.status(200).json(parsed);
    } catch (parseErr) {
      res.status(200).json({ raw: stdout });
    }
  });
}
