import { promises as fs } from "fs";
import path from "path";

const STATE_KEY = Symbol.for("synnergy.stage73.statePath");
const LOCK_KEY = Symbol.for("synnergy.stage73.lock");

function resolveStatePath() {
  const envPath = process.env.STAGE73_STATE_PATH;
  if (envPath && envPath.trim().length > 0) {
    return path.isAbsolute(envPath)
      ? envPath
      : path.join(process.cwd(), envPath);
  }
  return path.join(process.cwd(), "..", ".synnergy", "stage73_state.json");
}

export function getStage73StatePath() {
  if (!globalThis[STATE_KEY]) {
    globalThis[STATE_KEY] = resolveStatePath();
  }
  return globalThis[STATE_KEY];
}

export async function ensureStage73Storage() {
  const statePath = getStage73StatePath();
  await fs.mkdir(path.dirname(statePath), { recursive: true });
  return statePath;
}

export function withStage73Lock(fn) {
  if (!globalThis[LOCK_KEY]) {
    globalThis[LOCK_KEY] = Promise.resolve();
  }
  const pending = globalThis[LOCK_KEY];
  const next = pending.then(() => fn());
  globalThis[LOCK_KEY] = next.catch(() => {});
  return next;
}
