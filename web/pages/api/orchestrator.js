import { execFile } from 'child_process';
import path from 'path';

export default function handler(req, res) {
  const cliPath = path.join(process.cwd(), '..', 'cmd', 'synnergy', 'main.go');
  execFile(
    'go',
    ['run', cliPath, 'orchestrator', 'status', '--json'],
    { timeout: 15000 },
    (err, stdout, stderr) => {
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
    }
  );
}
