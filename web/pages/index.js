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
  const [bootstrapForm, setBootstrapForm] = useState({
    nodeId: "web-node",
    address: "",
    consensus: "Synnergy-PBFT",
    governance: "SYN-Gov",
    replicate: true,
    regulator: true,
    authorities: "",
  });
  const [bootstrapResult, setBootstrapResult] = useState(null);
  const [bootstrapError, setBootstrapError] = useState("");
  const [bootstrapLoading, setBootstrapLoading] = useState(false);

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

  const updateBootstrap = (name, value) => {
    setBootstrapForm((prev) => ({ ...prev, [name]: value }));
  };

  const toggleBootstrap = (name) => {
    setBootstrapForm((prev) => ({ ...prev, [name]: !prev[name] }));
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

  const runBootstrap = () => {
    setBootstrapLoading(true);
    setBootstrapError("");
    setBootstrapResult(null);
    const authorities = bootstrapForm.authorities
      ? bootstrapForm.authorities
          .split(",")
          .map((v) => v.trim())
          .filter((v) => v.length > 0)
      : [];
    fetch("/api/bootstrap", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        nodeId: bootstrapForm.nodeId,
        address: bootstrapForm.address,
        consensus: bootstrapForm.consensus,
        governance: bootstrapForm.governance,
        replicate: bootstrapForm.replicate,
        regulator: bootstrapForm.regulator,
        authorities,
      }),
    })
      .then(async (res) => {
        const text = await res.text();
        if (!res.ok) {
          try {
            const parsed = JSON.parse(text);
            throw new Error(parsed.error || "Bootstrap failed");
          } catch (err) {
            throw new Error(text || "Bootstrap failed");
          }
        }
        try {
          return JSON.parse(text);
        } catch (err) {
          return { raw: text };
        }
      })
      .then((data) => {
        setBootstrapResult(data);
      })
      .catch((err) => {
        setBootstrapError(err.message);
      })
      .finally(() => setBootstrapLoading(false));
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
              <strong>Bootstrap nodes:</strong> {orchestrator.bootstrapNodes}
            </p>
            <p>
              <strong>Replication:</strong> {orchestrator.replicationActive ? "active" : "paused"}
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
      <section style={{ marginTop: 20, marginBottom: 20 }}>
        <h2>Stage 79 Enterprise Bootstrap</h2>
        <p>
          Provision a ledger-backed node with authority, consensus, replication and privacy controls
          without leaving the control panel.
        </p>
        <div style={{ display: "grid", gap: 8, maxWidth: 600 }}>
          <label>
            Node ID
            <input
              type="text"
              value={bootstrapForm.nodeId}
              onChange={(e) => updateBootstrap("nodeId", e.target.value)}
              style={{ marginLeft: 10 }}
            />
          </label>
          <label>
            Address
            <input
              type="text"
              value={bootstrapForm.address}
              onChange={(e) => updateBootstrap("address", e.target.value)}
              style={{ marginLeft: 10 }}
              placeholder="host:port"
            />
          </label>
          <label>
            Consensus profile
            <input
              type="text"
              value={bootstrapForm.consensus}
              onChange={(e) => updateBootstrap("consensus", e.target.value)}
              style={{ marginLeft: 10 }}
            />
          </label>
          <label>
            Governance profile
            <input
              type="text"
              value={bootstrapForm.governance}
              onChange={(e) => updateBootstrap("governance", e.target.value)}
              style={{ marginLeft: 10 }}
            />
          </label>
          <label>
            Authorities (comma separated address=role)
            <input
              type="text"
              value={bootstrapForm.authorities}
              onChange={(e) => updateBootstrap("authorities", e.target.value)}
              style={{ marginLeft: 10 }}
              placeholder="authority1=ops,authority2=audit"
            />
          </label>
          <label>
            <input
              type="checkbox"
              checked={bootstrapForm.replicate}
              onChange={() => toggleBootstrap("replicate")}
            />
            Enable ledger replication
          </label>
          <label>
            <input
              type="checkbox"
              checked={bootstrapForm.regulator}
              onChange={() => toggleBootstrap("regulator")}
            />
            Attach regulatory validation
          </label>
        </div>
        <button style={{ marginTop: 10 }} onClick={runBootstrap} disabled={bootstrapLoading}>
          {bootstrapLoading ? "Bootstrapping..." : "Run bootstrap"}
        </button>
        {bootstrapError && <p style={{ color: "red" }}>{bootstrapError}</p>}
        {bootstrapResult && (
          <pre style={{ marginTop: 10 }}>
            {JSON.stringify(bootstrapResult, null, 2)}
          </pre>
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
