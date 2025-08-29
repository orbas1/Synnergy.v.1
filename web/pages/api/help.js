import { exec } from 'child_process';
import path from 'path';

export default function handler(req, res) {
  const { cmd } = req.query || {};
  if (!cmd) {
    res.status(400).json({ error: 'Missing cmd' });
    return;
  }
  const cliPath = path.join(process.cwd(), '..', 'cmd', 'synnergy', 'main.go');
  exec(`go run ${cliPath} ${cmd} --help`, { timeout: 5000 }, (err, stdout, stderr) => {
    if (err) {
      res.status(500).json({ error: stderr || err.message });
      return;
    }
    const lines = stdout.split('\n');
    const flags = [];
    let inFlags = false;
    for (const line of lines) {
      if (line.startsWith('Flags:')) {
        inFlags = true;
        continue;
      }
      if (inFlags) {
        if (!line.trim()) break;
        const m = line.trim().match(/^(?:-\S+,\s*)?(--\S+)(?:\s+\S+)?\s+(.*)$/);
        if (m) {
          const name = m[1].replace(/^--/, '');
          const desc = m[2] || '';
          flags.push({ name, desc });
        }
      }
    }
    res.status(200).json({ flags });
  });
}
