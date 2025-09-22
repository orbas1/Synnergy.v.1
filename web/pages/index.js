import { useEffect, useState } from "react";

const STORAGE_KEY = "syn-stage73-path";

export default function Home() {
  const [commands, setCommands] = useState([]);
  const [selected, setSelected] = useState("");
  const [flags, setFlags] = useState([]);
  const [inputs, setInputs] = useState({});
  const [extra, setExtra] = useState("");
  const [output, setOutput] = useState("");
  const [stage73Path, setStage73Path] = useState("");
  const [isRunning, setIsRunning] = useState(false);

  useEffect(() => {
    fetch("/api/commands")
      .then((res) => res.json())
      .then((data) => setCommands(data.commands || []))
      .catch((err) => setOutput(`Failed to load commands: ${err.message}`));
  }, []);

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }
    const stored = window.localStorage.getItem(STORAGE_KEY);
    if (stored) {
      setStage73Path(stored);
    }
  }, []);

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }
    if (stage73Path) {
      window.localStorage.setItem(STORAGE_KEY, stage73Path);
    } else {
      window.localStorage.removeItem(STORAGE_KEY);
    }
  }, [stage73Path]);

  useEffect(() => {
    if (selected) {
      fetch(`/api/help?cmd=${selected}`)
        .then((res) => res.json())
        .then((data) => {
          setFlags(data.flags || []);
          setInputs({});
        });
    } else {
      setFlags([]);
    }
  }, [selected]);

  const handleInput = (name, value) => {
    setInputs((prev) => ({ ...prev, [name]: value }));
  };

  const run = async () => {
    if (!selected) {
      setOutput("Select a command before running.");
      return;
    }
    const args = [];
    if (stage73Path.trim()) {
      args.push("--stage73-state", stage73Path.trim());
    }
    for (const f of flags) {
      const val = inputs[f.name];
      if (val && val.length) {
        args.push(`--${f.name}`);
        if (val !== "true" && val !== "false") {
          args.push(val);
        }
      }
    }
    if (extra.trim()) {
      args.push(...extra.trim().split(/\s+/));
    }
    setIsRunning(true);
    setOutput("Running...");
    try {
      const res = await fetch("/api/run", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ command: selected, args }),
      });
      const data = await res.json().catch(() => ({}));
      if (!res.ok) {
        setOutput(data.output || data.error || `Command failed (${res.status})`);
        return;
      }
      setOutput(data.output || data.error || "Command completed without output");
    } catch (err) {
      setOutput(`Failed to run command: ${err.message}`);
    } finally {
      setIsRunning(false);
    }
  };

  return (
    <main style={{ padding: 20 }}>
      <h1>Synnergy Control Panel</h1>
      <p>Select a command, fill in optional flags, and run.</p>
      <p>
        <a href="/regnode">Regulatory node console</a>
      </p>
      <div style={{ margin: "10px 0" }}>
        <label>
          Stage 73 state file
          <input
            style={{ marginLeft: 10, minWidth: 260 }}
            type="text"
            placeholder="/path/to/stage73_state.json"
            value={stage73Path}
            onChange={(e) => setStage73Path(e.target.value)}
          />
        </label>
        <p style={{ maxWidth: 600 }}>
          The CLI persists SYN3700, SYN3800, and SYN3900 data to this JSON file.
          Leave blank to use the default location on the server.
        </p>
      </div>
      <select value={selected} onChange={(e) => setSelected(e.target.value)}>
        <option value="">-- select command --</option>
        {commands.map((c) => (
          <option key={c.name} value={c.name}>
            {c.name} - {c.desc}
          </option>
        ))}
      </select>
      {flags.map((f) => (
        <div key={f.name} style={{ marginTop: 10 }}>
          <label>
            --{f.name}: {f.desc}
            <input
              style={{ marginLeft: 10 }}
              type="text"
              value={inputs[f.name] || ""}
              onChange={(e) => handleInput(f.name, e.target.value)}
            />
          </label>
        </div>
      ))}
      {selected && (
        <div style={{ marginTop: 10 }}>
          <label>
            Additional args
            <input
              style={{ marginLeft: 10 }}
              type="text"
              value={extra}
              onChange={(e) => setExtra(e.target.value)}
            />
          </label>
        </div>
      )}
      {selected && (
        <button style={{ marginTop: 20 }} onClick={run} disabled={isRunning}>
          {isRunning ? "Running" : "Run"}
        </button>
      )}
      <pre>{output}</pre>
    </main>
  );
}
