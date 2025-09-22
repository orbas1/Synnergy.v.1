import { useEffect, useState } from "react";

function TextInput({ label, value, onChange, type = "text", placeholder = "" }) {
  return (
    <label style={{ display: "block", marginBottom: 8 }}>
      <span style={{ display: "block", fontWeight: 600 }}>{label}</span>
      <input
        type={type}
        value={value}
        placeholder={placeholder}
        onChange={(e) => onChange(e.target.value)}
        style={{ width: "100%", padding: 8, marginTop: 4 }}
      />
    </label>
  );
}

export default function GrantsConsole() {
  const [grants, setGrants] = useState([]);
  const [status, setStatus] = useState({});
  const [audit, setAudit] = useState([]);
  const [selectedId, setSelectedId] = useState("");
  const [message, setMessage] = useState("");
  const [busy, setBusy] = useState(false);
  const [createForm, setCreateForm] = useState({ beneficiary: "", name: "", amount: "", authorizer: "", wallet: "", password: "" });
  const [releaseForm, setReleaseForm] = useState({ id: "", amount: "", note: "", wallet: "", password: "" });
  const [authorizeForm, setAuthorizeForm] = useState({ id: "", wallet: "", password: "" });

  useEffect(() => {
    refresh();
  }, []);

  const refresh = async () => {
    setBusy(true);
    setMessage("");
    try {
      const res = await fetch("/api/grants");
      const data = await res.json();
      setGrants(Array.isArray(data.grants) ? data.grants : []);
      setStatus(data.status || {});
    } catch (err) {
      setMessage(err.message);
    } finally {
      setBusy(false);
    }
  };

  const runAction = async (payload) => {
    setBusy(true);
    setMessage("");
    try {
      const res = await fetch("/api/grants", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      const data = await res.json();
      if (!res.ok) {
        throw new Error(data.error || "Action failed");
      }
      setGrants(Array.isArray(data.grants) ? data.grants : []);
      setStatus(data.status || {});
      setMessage(data.message || "Action completed");
    } catch (err) {
      setMessage(err.message);
    } finally {
      setBusy(false);
    }
  };

  const loadAudit = async (id) => {
    if (!id) return;
    setBusy(true);
    setMessage("");
    try {
      const res = await fetch(`/api/grants?id=${id}`);
      const data = await res.json();
      if (!res.ok) {
        throw new Error(data.error || "Unable to load audit");
      }
      setAudit(Array.isArray(data.audit) ? data.audit : []);
      setSelectedId(id);
    } catch (err) {
      setMessage(err.message);
    } finally {
      setBusy(false);
    }
  };

  const handleCreate = async (e) => {
    e.preventDefault();
    if (!createForm.wallet.trim() || !createForm.password) {
      setMessage("Wallet and password are required for grant creation");
      return;
    }
    const authorizers = createForm.authorizer
      .split(",")
      .map((item) => item.trim())
      .filter(Boolean);
    await runAction({
      action: "create",
      beneficiary: createForm.beneficiary.trim(),
      name: createForm.name.trim(),
      amount: createForm.amount.trim(),
      authorizers,
      wallet: createForm.wallet.trim(),
      password: createForm.password,
    });
    setCreateForm({
      beneficiary: "",
      name: "",
      amount: "",
      authorizer: "",
      wallet: createForm.wallet.trim(),
      password: "",
    });
  };

  const handleRelease = async (e) => {
    e.preventDefault();
    await runAction({
      action: "release",
      id: releaseForm.id.trim(),
      amount: releaseForm.amount.trim(),
      note: releaseForm.note.trim(),
      wallet: releaseForm.wallet.trim(),
      password: releaseForm.password,
    });
    setReleaseForm({ id: "", amount: "", note: "", wallet: "", password: "" });
  };

  const handleAuthorize = async (e) => {
    e.preventDefault();
    await runAction({
      action: "authorize",
      id: authorizeForm.id.trim(),
      wallet: authorizeForm.wallet.trim(),
      password: authorizeForm.password,
    });
    setAuthorizeForm({ id: "", wallet: "", password: "" });
  };

  return (
    <main style={{ padding: 24, maxWidth: 960, margin: "0 auto" }}>
      <h1>Enterprise Grant Lifecycle Console</h1>
      <p>
        Manage SYN3800 grant programs with auditable authorisations, encrypted wallet controls,
        and instant telemetry. This console fronts the hardened CLI, consensus-aware virtual machine
        and regulatory workflows added in Stage 85.
      </p>
      <p>
        <a href="/">Back to command launcher</a>
      </p>
      {message && (
        <div
          role="status"
          style={{
            marginBottom: 16,
            padding: 12,
            border: "1px solid #555",
            background: "#111",
            color: "#0f0",
          }}
        >
          {message}
        </div>
      )}
      {busy && <p>Processing request…</p>}

      <section style={{ marginTop: 24 }}>
        <h2>Create grant</h2>
        <form onSubmit={handleCreate}>
          <TextInput
            label="Beneficiary"
            value={createForm.beneficiary}
            onChange={(v) => setCreateForm((prev) => ({ ...prev, beneficiary: v }))}
          />
          <TextInput
            label="Program name"
            value={createForm.name}
            onChange={(v) => setCreateForm((prev) => ({ ...prev, name: v }))}
          />
          <TextInput
            label="Amount"
            value={createForm.amount}
            onChange={(v) => setCreateForm((prev) => ({ ...prev, amount: v }))}
            type="number"
          />
          <TextInput
            label="Authorizer wallets (path:password, comma separated)"
            value={createForm.authorizer}
            placeholder="memory-wallet-0:pass"
            onChange={(v) => setCreateForm((prev) => ({ ...prev, authorizer: v }))}
          />
          <TextInput
            label="Creator wallet path"
            value={createForm.wallet}
            onChange={(v) => setCreateForm((prev) => ({ ...prev, wallet: v }))}
            placeholder="memory-wallet-creator"
          />
          <TextInput
            label="Creator wallet password"
            value={createForm.password}
            onChange={(v) => setCreateForm((prev) => ({ ...prev, password: v }))}
            type="password"
          />
          <button type="submit" disabled={busy} style={{ marginTop: 12 }}>
            Create grant
          </button>
        </form>
      </section>

      <section style={{ marginTop: 24 }}>
        <h2>Release funds</h2>
        <form onSubmit={handleRelease}>
          <TextInput
            label="Grant ID"
            value={releaseForm.id}
            onChange={(v) => setReleaseForm((prev) => ({ ...prev, id: v }))}
          />
          <TextInput
            label="Amount"
            value={releaseForm.amount}
            onChange={(v) => setReleaseForm((prev) => ({ ...prev, amount: v }))}
            type="number"
          />
          <TextInput
            label="Release note"
            value={releaseForm.note}
            onChange={(v) => setReleaseForm((prev) => ({ ...prev, note: v }))}
            placeholder="phase one milestone"
          />
          <TextInput
            label="Wallet path"
            value={releaseForm.wallet}
            onChange={(v) => setReleaseForm((prev) => ({ ...prev, wallet: v }))}
            placeholder="memory-wallet-0"
          />
          <TextInput
            label="Wallet password"
            value={releaseForm.password}
            onChange={(v) => setReleaseForm((prev) => ({ ...prev, password: v }))}
            type="password"
          />
          <button type="submit" disabled={busy} style={{ marginTop: 12 }}>
            Release funds
          </button>
        </form>
      </section>

      <section style={{ marginTop: 24 }}>
        <h2>Authorize wallet</h2>
        <form onSubmit={handleAuthorize}>
          <TextInput
            label="Grant ID"
            value={authorizeForm.id}
            onChange={(v) => setAuthorizeForm((prev) => ({ ...prev, id: v }))}
          />
          <TextInput
            label="Wallet path"
            value={authorizeForm.wallet}
            onChange={(v) => setAuthorizeForm((prev) => ({ ...prev, wallet: v }))}
          />
          <TextInput
            label="Wallet password"
            value={authorizeForm.password}
            onChange={(v) => setAuthorizeForm((prev) => ({ ...prev, password: v }))}
            type="password"
          />
          <button type="submit" disabled={busy} style={{ marginTop: 12 }}>
            Authorize wallet
          </button>
        </form>
      </section>

      <section style={{ marginTop: 32 }}>
        <h2>Active grants</h2>
        {grants.length === 0 ? (
          <p>No grants recorded yet.</p>
        ) : (
          <table style={{ width: "100%", borderCollapse: "collapse" }}>
            <thead>
              <tr>
                <th>ID</th>
                <th>Name</th>
                <th>Beneficiary</th>
                <th>Status</th>
                <th>Released / Amount</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {grants.map((grant) => (
                <tr key={grant.id} style={{ borderTop: "1px solid #333" }}>
                  <td>{grant.id}</td>
                  <td>{grant.name}</td>
                  <td>{grant.beneficiary}</td>
                  <td>{grant.status}</td>
                  <td>
                    {grant.released} / {grant.amount}
                  </td>
                  <td>
                    <button type="button" onClick={() => loadAudit(grant.id)} disabled={busy}>
                      View audit
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        <div style={{ marginTop: 16 }}>
          <strong>Totals:</strong> {status.total || 0} grants, {status.pending || 0} pending, {status.active || 0} active, {status.completed || 0} completed.
        </div>
      </section>

      <section style={{ marginTop: 32 }}>
        <h2>Audit trail</h2>
        {selectedId ? <p>Timeline for grant {selectedId}</p> : <p>Select a grant to inspect its events.</p>}
        {audit.length === 0 && selectedId ? <p>No events recorded.</p> : null}
        <ul>
          {audit.map((evt, index) => (
            <li key={`${evt.timestamp}-${index}`}>
              <strong>{evt.type}</strong> — {evt.timestamp} — {evt.signer || "system"}
              {evt.amount ? ` — ${evt.amount}` : ""}
              {evt.note ? ` — ${evt.note}` : ""}
            </li>
          ))}
        </ul>
      </section>
    </main>
  );
}
