import { exec } from "child_process";
import path from "path";

export default function handler(req, res) {
  const cliPath = path.join(process.cwd(), "..", "cmd", "synnergy", "main.go");
  exec(`go run ${cliPath} modules list --json`, { timeout: 10000 }, (err, stdout, stderr) => {
    if (err) {
      res.status(500).json({ error: stderr || err.message });
      return;
    }
    const trimmed = stdout.trim();
    const start = trimmed.indexOf("[");
    const end = trimmed.lastIndexOf("]");
    if (start === -1 || end === -1) {
      res.status(500).json({ error: "unexpected modules output" });
      return;
    }
    let modules;
    try {
      modules = JSON.parse(trimmed.slice(start, end + 1));
    } catch (parseErr) {
      res.status(500).json({ error: `failed to parse modules: ${parseErr.message}` });
      return;
    }
    res.status(200).json({ modules });
  });
}
