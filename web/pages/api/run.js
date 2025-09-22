import { execFile } from "child_process";
import { promisify } from "util";
import path from "path";

import {
  ensureStage73Storage,
  getStage73StatePath,
  withStage73Lock,
} from "../../lib/stage73.js";

const execFileAsync = promisify(execFile);

function sanitizeArgs(args = []) {
  if (!Array.isArray(args)) {
    return [];
  }
  const result = [];
  for (let i = 0; i < args.length; i += 1) {
    const value = args[i];
    if (typeof value !== "string") {
      continue;
    }
    if (value === "--stage73-state") {
      i += 1; // skip attempted override value
      continue;
    }
    result.push(value);
  }
  return result;
}

export default async function handler(req, res) {
  if (req.method !== "POST") {
    res.status(405).json({ output: "Method not allowed" });
    return;
  }
  const { command, args = [] } = req.body || {};
  if (!command || typeof command !== "string") {
    res.status(400).json({ output: "Missing command" });
    return;
  }
  const cliPath = path.join(process.cwd(), "..", "cmd", "synnergy", "main.go");
  try {
    const statePath = await ensureStage73Storage();
    const runResult = await withStage73Lock(async () => {
      const finalArgs = [
        "run",
        cliPath,
        "--stage73-state",
        statePath,
        command,
        ...sanitizeArgs(args),
      ];
      const { stdout, stderr } = await execFileAsync("go", finalArgs, {
        timeout: 15000,
        cwd: path.join(process.cwd(), ".."),
        env: {
          ...process.env,
          STAGE73_STATE_PATH: getStage73StatePath(),
        },
      });
      return { stdout, stderr };
    });
    const output = runResult.stdout?.trim();
    const error = runResult.stderr?.trim();
    res.status(200).json({ output: output || error || "" });
  } catch (err) {
    const message = err.stderr?.toString().trim() || err.message || "execution failed";
    res.status(500).json({ output: message });
  }
}
