import { useEffect, useState } from "react";

export default function Home() {
  const [commands, setCommands] = useState([]);
  const [selected, setSelected] = useState("");
  const [flags, setFlags] = useState([]);
  const [inputs, setInputs] = useState({});
  const [extra, setExtra] = useState("");
  const [output, setOutput] = useState("");
  const [integration, setIntegration] = useState(null);
  const [integrationError, setIntegrationError] = useState("");
  const [integrationLoading, setIntegrationLoading] = useState(true);
  const readinessKeys = [
    ["Security", "security"],
    ["Scalability", "scalability"],
    ["Privacy", "privacy"],
    ["Governance", "governance"],
    ["Interoperability", "interoperability"],
    ["Compliance", "compliance"],
  ];

  useEffect(() => {
    fetch("/api/commands")
      .then((res) => res.json())
      .then((data) => setCommands(data.commands || []));
    refreshIntegration();
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

  const refreshIntegration = () => {
    setIntegrationLoading(true);
    fetch("/api/integration")
      .then((res) => res.json())
      .then((data) => {
        if (data.error) {
          setIntegrationError(data.error);
          setIntegration(null);
        } else {
          setIntegration(data);
          setIntegrationError("");
        }
      })
      .catch((err) => {
        setIntegrationError(err.message);
        setIntegration(null);
      })
      .finally(() => setIntegrationLoading(false));
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

  return (
    <main style={{ padding: 20 }}>
      <h1>Synnergy Control Panel</h1>
      <p>Select a command, fill in optional flags, and run.</p>
      <p>
        <a href="/regnode">Regulatory node console</a>
      </p>
      <section style={{ marginBottom: 30 }}>
        <h2>Enterprise Integration Health</h2>
        {integrationLoading && <p>Loading integration diagnostics…</p>}
        {!integrationLoading && integrationError && (
          <p role="alert" style={{ color: "#b00020" }}>
            {integrationError}
          </p>
        )}
        {!integrationLoading && integration && (
          <div
            style={{
              border: "1px solid #ccc",
              borderRadius: 8,
              padding: 16,
              background: "#f7f9fc",
            }}
            aria-live="polite"
          >
            <p>
              <strong>Last updated:</strong> {integration.timestamp}
            </p>
            <ul>
              <li>
                <strong>VM:</strong> {integration.vm.mode} (running: {integration.vm.running ? "yes" : "no"},
                concurrency {integration.vm.concurrency})
              </li>
              <li>
                <strong>Node:</strong> {integration.node.id} — height {integration.node.block_height}, pending
                {" "}
                {integration.node.pending_transactions}
              </li>
              <li>
                <strong>Wallet:</strong> {integration.wallet.address}
              </li>
              <li>
                <strong>Consensus networks:</strong> {integration.consensus.networks} (relayers
                {" "}
                {integration.consensus.authorized_relays})
              </li>
              <li>
                <strong>Authority nodes:</strong> {integration.authority.registered}
              </li>
            </ul>
            <section style={{ marginTop: 12 }}>
              <h3>Enterprise readiness</h3>
              <ul>
                {readinessKeys.map(([label, key]) => {
                  const check = integration.enterprise?.[key];
                  if (!check) {
                    return (
                      <li key={key}>
                        <strong>{label}:</strong> not reported
                      </li>
                    );
                  }
                  return (
                    <li key={key}>
                      <strong>{label}:</strong> {check.detail} — healthy: {check.healthy ? "yes" : "no"}, latency {check.latency}
                    </li>
                  );
                })}
              </ul>
            </section>
            <details style={{ marginTop: 10 }}>
              <summary>Diagnostics log</summary>
              <ul>
                {Object.entries(integration.diagnostics || {}).map(([key, value]) => (
                  <li key={key}>
                    <strong>{key}:</strong> {value}
                  </li>
                ))}
              </ul>
            </details>
            {integration.issues && integration.issues.length > 0 && (
              <div role="alert" style={{ marginTop: 12, color: "#b00020" }}>
                <p>
                  <strong>Issues detected</strong>
                </p>
                <ul>
                  {integration.issues.map((issue, idx) => (
                    <li key={idx}>{issue}</li>
                  ))}
                </ul>
              </div>
            )}
            <button style={{ marginTop: 10 }} onClick={refreshIntegration}>
              Refresh diagnostics
            </button>
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
