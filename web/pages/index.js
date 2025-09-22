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
  const [treasury, setTreasury] = useState(null);
  const [treasuryError, setTreasuryError] = useState("");
  const [treasuryLoading, setTreasuryLoading] = useState(true);
  const [treasuryForm, setTreasuryForm] = useState({
    operator: "",
    issueAddress: "",
    issueAmount: "",
    burnAddress: "",
    burnAmount: "",
    transferAddress: "",
    transferAmount: "",
    authorizeAddress: "",
    revokeAddress: "",
  });
  const [treasuryActionStatus, setTreasuryActionStatus] = useState("");
  const [treasuryActionError, setTreasuryActionError] = useState("");
  const [treasuryActionLoading, setTreasuryActionLoading] = useState(false);

  useEffect(() => {
    fetch("/api/commands")
      .then((res) => res.json())
      .then((data) => setCommands(data.commands || []));
  }, []);

  useEffect(() => {
    refreshOrchestrator();
    refreshTreasury();
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

  const refreshTreasury = () => {
    setTreasuryLoading(true);
    fetch("/api/treasury")
      .then((res) => {
        if (!res.ok) {
          throw new Error("Failed to load treasury telemetry");
        }
        return res.json();
      })
      .then((data) => {
        setTreasury(data);
        setTreasuryError("");
      })
      .catch((err) => {
        setTreasury(null);
        setTreasuryError(err.message);
      })
      .finally(() => setTreasuryLoading(false));
  };

  const updateTreasuryForm = (field, value) => {
    setTreasuryForm((prev) => ({ ...prev, [field]: value }));
  };

  const performTreasuryAction = async (action) => {
    setTreasuryActionStatus("");
    setTreasuryActionError("");
    setTreasuryActionLoading(true);
    const args = ["telemetry", "--json"];
    const form = treasuryForm;
    if (form.operator.trim()) {
      args.push("--operator", form.operator.trim());
    }
    try {
      switch (action) {
        case "issue": {
          if (!form.issueAddress.trim() || !form.issueAmount.trim()) {
            throw new Error("Issue requires address and amount");
          }
          args.push("--issue", `${form.issueAddress.trim()}:${form.issueAmount.trim()}`);
          break;
        }
        case "burn": {
          if (!form.burnAddress.trim() || !form.burnAmount.trim()) {
            throw new Error("Burn requires address and amount");
          }
          args.push("--burn", `${form.burnAddress.trim()}:${form.burnAmount.trim()}`);
          break;
        }
        case "transfer": {
          if (!form.transferAddress.trim() || !form.transferAmount.trim()) {
            throw new Error("Transfer requires destination and amount");
          }
          args.push("--transfer", `${form.transferAddress.trim()}:${form.transferAmount.trim()}`);
          break;
        }
        case "authorize": {
          if (!form.authorizeAddress.trim()) {
            throw new Error("Provide an operator to authorize");
          }
          args.push("--authorize-operator", form.authorizeAddress.trim());
          break;
        }
        case "revoke": {
          if (!form.revokeAddress.trim()) {
            throw new Error("Provide an operator to revoke");
          }
          args.push("--revoke-operator", form.revokeAddress.trim());
          break;
        }
        default:
          throw new Error("Unsupported action");
      }

      const res = await fetch("/api/run", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ command: "coin", args }),
      });
      const data = await res.json();
      if (!res.ok) {
        throw new Error(data.output || "Treasury action failed");
      }
      if (data.output) {
        try {
          JSON.parse(data.output);
        } catch (err) {
          // Ignore parse errors and keep raw output for troubleshooting
        }
      }
      setTreasuryActionStatus(`Action ${action} completed`);
      refreshTreasury();
    } catch (err) {
      setTreasuryActionError(err.message);
    } finally {
      setTreasuryActionLoading(false);
    }
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
      <section style={{ marginTop: 20, marginBottom: 20 }}>
        <h2>Stage 80 Treasury Telemetry</h2>
        <button onClick={refreshTreasury} disabled={treasuryLoading}>
          {treasuryLoading ? "Refreshing..." : "Refresh"}
        </button>
        {treasuryError && <p style={{ color: "red" }}>{treasuryError}</p>}
        {treasury && (
          <div style={{ marginTop: 10 }}>
            <p>
              <strong>Wallet:</strong> {treasury.wallet}
            </p>
            <p>
              <strong>Supply:</strong> minted {treasury.minted}, burned {treasury.burned},
              circulating {treasury.circulating}
            </p>
            <p>
              <strong>Ledger height:</strong> {treasury.ledgerHeight}
            </p>
            {treasury.health && (
              <p>
                <strong>Health:</strong>{" "}
                {Object.entries(treasury.health)
                  .map(([key, value]) => `${key}: ${value}`)
                  .join(", ")}
              </p>
            )}
            <p>
              <strong>Consensus bridges:</strong> {treasury.consensusNetworks}
            </p>
            <p>
              <strong>Authority nodes:</strong> {treasury.authorityNodes}
            </p>
            {Array.isArray(treasury.missingOpcodes) && treasury.missingOpcodes.length > 0 ? (
              <p style={{ color: "orange" }}>
                Missing opcodes: {treasury.missingOpcodes.join(", ")}
              </p>
            ) : (
              <p>Opcode documentation is complete.</p>
            )}
            <p>
              <strong>Operators:</strong>{" "}
              {Array.isArray(treasury.operators) && treasury.operators.length > 0
                ? treasury.operators.join(", ")
                : "Treasury guardian only"}
            </p>
            {treasury.gasCoverage && (
              <div>
                <strong>Gas coverage:</strong>
                <table style={{ marginTop: 8, borderCollapse: "collapse" }}>
                  <thead>
                    <tr>
                      <th style={{ border: "1px solid #ccc", padding: "4px 8px" }}>Opcode</th>
                      <th style={{ border: "1px solid #ccc", padding: "4px 8px" }}>Gas</th>
                    </tr>
                  </thead>
                  <tbody>
                    {Object.entries(treasury.gasCoverage)
                      .sort(([a], [b]) => a.localeCompare(b))
                      .map(([name, gas]) => (
                        <tr key={name}>
                          <td style={{ border: "1px solid #ccc", padding: "4px 8px" }}>{name}</td>
                          <td style={{ border: "1px solid #ccc", padding: "4px 8px" }}>{gas}</td>
                        </tr>
                      ))}
                  </tbody>
                </table>
              </div>
            )}
            {Array.isArray(treasury.auditTrail) && treasury.auditTrail.length > 0 && (
              <div style={{ marginTop: 12 }}>
                <strong>Recent audit trail:</strong>
                <table style={{ marginTop: 8, borderCollapse: "collapse" }}>
                  <thead>
                    <tr>
                      <th style={{ border: "1px solid #ccc", padding: "4px 8px" }}>Time</th>
                      <th style={{ border: "1px solid #ccc", padding: "4px 8px" }}>Type</th>
                      <th style={{ border: "1px solid #ccc", padding: "4px 8px" }}>Amount</th>
                      <th style={{ border: "1px solid #ccc", padding: "4px 8px" }}>Digest</th>
                    </tr>
                  </thead>
                  <tbody>
                    {treasury.auditTrail
                      .slice(-10)
                      .reverse()
                      .map((evt, idx) => (
                        <tr key={`${evt.digest}-${idx}`}>
                          <td style={{ border: "1px solid #ccc", padding: "4px 8px" }}>{evt.timestamp}</td>
                          <td style={{ border: "1px solid #ccc", padding: "4px 8px" }}>{evt.type}</td>
                          <td style={{ border: "1px solid #ccc", padding: "4px 8px" }}>{evt.amount}</td>
                          <td style={{ border: "1px solid #ccc", padding: "4px 8px", maxWidth: 160, wordBreak: "break-all" }}>
                            {evt.digest}
                          </td>
                        </tr>
                      ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}
        <div style={{ marginTop: 15 }}>
          <h3>Execute Treasury Actions</h3>
          <div style={{ display: "flex", flexDirection: "column", gap: 8, maxWidth: 420 }}>
            <label>
              Operator context
              <input
                type="text"
                value={treasuryForm.operator}
                onChange={(e) => updateTreasuryForm("operator", e.target.value)}
                placeholder="optional operator address"
              />
            </label>
            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 8 }}>
              <label>
                Issue address
                <input
                  type="text"
                  value={treasuryForm.issueAddress}
                  onChange={(e) => updateTreasuryForm("issueAddress", e.target.value)}
                />
              </label>
              <label>
                Issue amount
                <input
                  type="text"
                  value={treasuryForm.issueAmount}
                  onChange={(e) => updateTreasuryForm("issueAmount", e.target.value)}
                />
              </label>
            </div>
            <button onClick={() => performTreasuryAction("issue")} disabled={treasuryActionLoading}>
              {treasuryActionLoading ? "Processing..." : "Issue"}
            </button>
            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 8 }}>
              <label>
                Burn address
                <input
                  type="text"
                  value={treasuryForm.burnAddress}
                  onChange={(e) => updateTreasuryForm("burnAddress", e.target.value)}
                />
              </label>
              <label>
                Burn amount
                <input
                  type="text"
                  value={treasuryForm.burnAmount}
                  onChange={(e) => updateTreasuryForm("burnAmount", e.target.value)}
                />
              </label>
            </div>
            <button onClick={() => performTreasuryAction("burn")} disabled={treasuryActionLoading}>
              {treasuryActionLoading ? "Processing..." : "Burn"}
            </button>
            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 8 }}>
              <label>
                Transfer destination
                <input
                  type="text"
                  value={treasuryForm.transferAddress}
                  onChange={(e) => updateTreasuryForm("transferAddress", e.target.value)}
                />
              </label>
              <label>
                Transfer amount
                <input
                  type="text"
                  value={treasuryForm.transferAmount}
                  onChange={(e) => updateTreasuryForm("transferAmount", e.target.value)}
                />
              </label>
            </div>
            <button onClick={() => performTreasuryAction("transfer")} disabled={treasuryActionLoading}>
              {treasuryActionLoading ? "Processing..." : "Transfer"}
            </button>
            <label>
              Authorize operator
              <input
                type="text"
                value={treasuryForm.authorizeAddress}
                onChange={(e) => updateTreasuryForm("authorizeAddress", e.target.value)}
              />
            </label>
            <button onClick={() => performTreasuryAction("authorize")} disabled={treasuryActionLoading}>
              {treasuryActionLoading ? "Processing..." : "Authorize"}
            </button>
            <label>
              Revoke operator
              <input
                type="text"
                value={treasuryForm.revokeAddress}
                onChange={(e) => updateTreasuryForm("revokeAddress", e.target.value)}
              />
            </label>
            <button onClick={() => performTreasuryAction("revoke")} disabled={treasuryActionLoading}>
              {treasuryActionLoading ? "Processing..." : "Revoke"}
            </button>
          </div>
          {treasuryActionStatus && (
            <p style={{ color: "green" }}>{treasuryActionStatus}</p>
          )}
          {treasuryActionError && (
            <p style={{ color: "red" }}>{treasuryActionError}</p>
          )}
        </div>
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
