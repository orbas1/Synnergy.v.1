import { execFile } from 'child_process';
import path from 'path';

export default function handler(req, res) {
  if (req.method !== 'POST') {
    return res.status(405).json({ output: 'Method not allowed' });
  }
  const { command, args = [] } = req.body || {};
  if (!command) {
    return res.status(400).json({ output: 'Missing command' });
  }
  const cliPath = path.join(process.cwd(), '..', 'cmd', 'synnergy', 'main.go');
  const finalArgs = ['run', cliPath, command, ...args];
  execFile('go', finalArgs, { timeout: 10000 }, (err, stdout, stderr) => {
    if (err) {
      res.status(500).json({ output: stderr || err.message });
      return;
    }
    res.status(200).json({ output: stdout });
  });
}
