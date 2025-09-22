import { promises as fs } from "fs";
import crypto from "crypto";

import {
  ensureStage73Storage,
  withStage73Lock,
} from "../../lib/stage73.js";

export default async function handler(req, res) {
  if (req.method !== "GET") {
    res.status(405).json({ error: "Method not allowed" });
    return;
  }
  try {
    const statePath = await ensureStage73Storage();
    let status = 200;
    let payload = {
      snapshot: null,
      digest: null,
      updatedAt: null,
      size: 0,
    };
    await withStage73Lock(async () => {
      let raw;
      try {
        raw = await fs.readFile(statePath, "utf8");
      } catch (err) {
        if (err.code === "ENOENT") {
          status = 200;
          payload = { ...payload, message: "No Stage 73 snapshot present" };
          return;
        }
        throw err;
      }
      if (!raw || raw.trim().length === 0) {
        payload = { ...payload, message: "Snapshot empty" };
        return;
      }
      let parsed;
      try {
        parsed = JSON.parse(raw);
      } catch (err) {
        status = 500;
        payload = { error: "Snapshot is corrupt", detail: err.message };
        return;
      }
      const canonical = JSON.stringify(parsed);
      const digest = crypto
        .createHash("sha256")
        .update(canonical)
        .digest("hex");
      const stat = await fs.stat(statePath);
      payload = {
        snapshot: parsed,
        digest,
        updatedAt: stat.mtime.toISOString(),
        size: canonical.length,
      };
    });
    res.status(status).json(payload);
  } catch (err) {
    res.status(500).json({ error: err.message || "failed to read snapshot" });
  }
}
