import { execFile } from 'child_process';
import path from 'path';

export default function handler(req, res) {
  const cliPath = path.join(process.cwd(), '..', 'cmd', 'synnergy', 'main.go');
  const args = ['run', cliPath, 'integration', 'status', '--format', 'json'];
  execFile('go', args, { timeout: 15000 }, (err, stdout, stderr) => {
    if (err) {
      res.status(500).json({ error: stderr || err.message });
      return;
    }
    try {
      const payload = JSON.parse(stdout);
      res.status(200).json(payload);
    } catch (parseErr) {
      res
        .status(500)
        .json({ error: 'Failed to parse integration status', detail: parseErr.message, raw: stdout });
    }
  });
}
