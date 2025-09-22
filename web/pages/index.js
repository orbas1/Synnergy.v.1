import { useEffect, useState } from "react";

export default function Home() {
  const [commands, setCommands] = useState([]);
  const [selected, setSelected] = useState("");
  const [flags, setFlags] = useState([]);
  const [inputs, setInputs] = useState({});
  const [extra, setExtra] = useState("");
  const [output, setOutput] = useState("");
  const [orchestrator, setOrchestrator] = useState(null);
  const [orchError, setOrchError] = useState("");
  const [orchLoading, setOrchLoading] = useState(true);

  useEffect(() => {
    fetch("/api/commands")
      .then((res) => res.json())
      .then((data) => setCommands(data.commands || []));
  }, []);

  useEffect(() => {
    refreshOrchestrator();
  }, []);

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
    const args = [];
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
    const res = await fetch("/api/run", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ command: selected, args }),
    });
    const data = await res.json();
    setOutput(data.output || data.error);
  };

  const refreshOrchestrator = () => {
    setOrchLoading(true);
    fetch("/api/orchestrator")
      .then((res) => {
        if (!res.ok) {
          throw new Error("Failed to load orchestrator status");
        }
        return res.json();
      })
      .then((data) => {
        setOrchestrator(data);
        setOrchError("");
      })
      .catch((err) => {
        setOrchestrator(null);
        setOrchError(err.message);
      })
      .finally(() => setOrchLoading(false));
  };

  return (
    <main style={{ padding: 20 }}>
      <h1>Synnergy Control Panel</h1>
      <p>Select a command, fill in optional flags, and run.</p>
      <p>
        <a href="/regnode">Regulatory node console</a>
        {" | "}
        <a href="/warfare">Warfare command center</a>
      </p>
      <section style={{ marginTop: 20, marginBottom: 20 }}>
        <h2>Enterprise Orchestrator Status</h2>
        <button onClick={refreshOrchestrator} disabled={orchLoading}>
          {orchLoading ? "Refreshing..." : "Refresh"}
        </button>
        {orchError && <p style={{ color: "red" }}>{orchError}</p>}
        {orchestrator && (
          <div style={{ marginTop: 10 }}>
            <p>
              <strong>Timestamp:</strong> {orchestrator.timestamp}
            </p>
            <p>
              <strong>VM Mode:</strong> {orchestrator.vmMode} (running {""}
              {orchestrator.vmRunning ? "yes" : "no"}, concurrency {""}
              {orchestrator.vmConcurrency})
            </p>
            <p>
              <strong>Consensus networks:</strong> {orchestrator.consensusNetworks}
            </p>
            <p>
              <strong>Authority nodes:</strong> {orchestrator.authorityNodes}
            </p>
            <p>
              <strong>Wallet:</strong> {orchestrator.walletAddress}
            </p>
            <p>
              <strong>Ledger height:</strong> {orchestrator.ledgerHeight}
            </p>
            {Array.isArray(orchestrator.missingOpcodes) && orchestrator.missingOpcodes.length > 0 ? (
              <p style={{ color: "orange" }}>
                Missing opcode documentation: {orchestrator.missingOpcodes.join(", ")}
              </p>
            ) : (
              <p>Opcode documentation is complete.</p>
            )}
          </div>
        )}
      </section>
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
        <button style={{ marginTop: 20 }} onClick={run}>
          Run
        </button>
      )}
      <pre>{output}</pre>
    </main>
  );
}
