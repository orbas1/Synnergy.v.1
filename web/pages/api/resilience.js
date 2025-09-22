import { execFile } from 'child_process';
import path from 'path';

function extractJson(output) {
  const start = output.indexOf('{');
  if (start === -1) {
    return output;
  }
  return output.slice(start);
}

export default function handler(req, res) {
  const cliPath = path.join(process.cwd(), '..', 'cmd', 'synnergy', 'main.go');
  execFile(
    'go',
    ['run', cliPath, 'highavailability', 'report', '--json'],
    { timeout: 20000 },
    (err, stdout, stderr) => {
      if (err) {
        res.status(500).json({ error: stderr || err.message });
        return;
      }
      try {
        const parsed = JSON.parse(extractJson(stdout));
        res.status(200).json(parsed);
      } catch (parseErr) {
        res.status(200).json({ raw: stdout });
      }
    }
  );
}
