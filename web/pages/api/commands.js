import { exec } from 'child_process';
import path from 'path';

export default function handler(req, res) {
  const cliPath = path.join(process.cwd(), '..', 'cmd', 'synnergy', 'main.go');
  exec(`go run ${cliPath} --help`, { timeout: 5000 }, (err, stdout, stderr) => {
    if (err) {
      res.status(500).json({ error: stderr || err.message });
      return;
    }
    const lines = stdout.split('\n');
    const commands = [];
    let collecting = false;
    for (const line of lines) {
      if (line.startsWith('Available Commands:')) {
        collecting = true;
        continue;
      }
      if (collecting) {
        const trimmed = line.trim();
        if (!trimmed) break;
        const match = trimmed.match(/^(\S+)\s+(.+)$/);
        if (match) {
          commands.push({ name: match[1], desc: match[2] });
        }
      }
    }
    res.status(200).json({ commands });
  });
}
