import { useCallback, useEffect, useMemo, useState } from "react";

const POLL_INTERVAL_MS = 15000;

function formatDate(value) {
  if (!value) {
    return "--";
  }
  try {
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return value;
    }
    return date.toLocaleString();
  } catch (err) {
    return value;
  }
}

function Section({ title, children }) {
  return (
    <section style={{ marginBottom: "2rem" }}>
      <h2>{title}</h2>
      {children}
    </section>
  );
}

function Status({ result }) {
  if (!result) {
    return null;
  }
  const color = result.success ? "#0a7c4a" : "#b00020";
  return (
    <div role="status" style={{ marginTop: "0.5rem", color }}>
      {result.message && result.message.length > 240 ? (
        <details>
          <summary>{result.success ? "Command executed" : "Command failed"}</summary>
          <pre style={{ whiteSpace: "pre-wrap" }}>{result.message}</pre>
        </details>
      ) : (
        <pre style={{ whiteSpace: "pre-wrap", margin: 0 }}>{result.message}</pre>
      )}
    </div>
  );
}

function CommandForm({ title, description, command, buildArgs, fields, onSuccess }) {
  const initialValues = useMemo(() => {
    const values = {};
    for (const field of fields) {
      values[field.name] = field.defaultValue ?? "";
    }
    return values;
  }, [fields]);
  const [values, setValues] = useState(initialValues);
  const [submitting, setSubmitting] = useState(false);
  const [result, setResult] = useState(null);

  const handleChange = (name, value) => {
    setValues((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setSubmitting(true);
    setResult(null);
    try {
      const args = buildArgs(values);
      if (!Array.isArray(args)) {
        throw new Error("Invalid command arguments");
      }
      const response = await fetch("/api/run", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ command, args }),
      });
      const data = await response.json();
      if (!response.ok) {
        throw new Error(data.output || data.error || "Command failed");
      }
      setResult({ success: true, message: data.output || "Command executed" });
      if (onSuccess) {
        await onSuccess();
      }
    } catch (err) {
      setResult({ success: false, message: err.message });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div style={{ border: "1px solid #d0d7de", borderRadius: 8, padding: "1rem", marginBottom: "1.5rem" }}>
      <h3 style={{ marginTop: 0 }}>{title}</h3>
      {description && <p style={{ marginTop: 0 }}>{description}</p>}
      <form onSubmit={handleSubmit}>
        {fields.map((field) => (
          <div key={field.name} style={{ display: "flex", flexDirection: "column", marginBottom: "0.75rem" }}>
            <label htmlFor={`${command}-${field.name}`} style={{ fontWeight: 600 }}>
              {field.label}
            </label>
            {field.type === "textarea" ? (
              <textarea
                id={`${command}-${field.name}`}
                name={field.name}
                aria-required={field.required}
                required={field.required}
                placeholder={field.placeholder}
                value={values[field.name]}
                onChange={(event) => handleChange(field.name, event.target.value)}
                style={{ padding: "0.5rem", borderRadius: 4, border: "1px solid #d0d7de", minHeight: "4rem" }}
              />
            ) : (
              <input
                id={`${command}-${field.name}`}
                name={field.name}
                type={field.type || "text"}
                aria-required={field.required}
                required={field.required}
                placeholder={field.placeholder}
                value={values[field.name]}
                onChange={(event) => handleChange(field.name, event.target.value)}
                style={{ padding: "0.5rem", borderRadius: 4, border: "1px solid #d0d7de" }}
              />
            )}
            {field.help && (
              <small style={{ color: "#57606a" }}>{field.help}</small>
            )}
          </div>
        ))}
        <button
          type="submit"
          disabled={submitting}
          style={{
            padding: "0.6rem 1.2rem",
            borderRadius: 6,
            border: "none",
            backgroundColor: "#0969da",
            color: "#fff",
            cursor: submitting ? "wait" : "pointer",
          }}
        >
          {submitting ? "Running…" : `Run ${command}`}
        </button>
      </form>
      <Status result={result} />
    </div>
  );
}

function SnapshotView({ snapshot }) {
  if (!snapshot) {
    return <p>No Stage 73 snapshot has been persisted yet.</p>;
  }
  return (
    <div style={{ display: "grid", gap: "1.5rem" }}>
      {snapshot.index && (
        <div>
          <h3>SYN3700 Index</h3>
          <p>
            {snapshot.index.name} ({snapshot.index.symbol}) — controllers: {snapshot.index.controllers.length}
          </p>
          <table style={{ width: "100%", borderCollapse: "collapse" }}>
            <thead>
              <tr>
                <th style={{ textAlign: "left", borderBottom: "1px solid #d0d7de", paddingBottom: 4 }}>Token</th>
                <th style={{ textAlign: "right", borderBottom: "1px solid #d0d7de", paddingBottom: 4 }}>Weight</th>
                <th style={{ textAlign: "right", borderBottom: "1px solid #d0d7de", paddingBottom: 4 }}>Drift</th>
              </tr>
            </thead>
            <tbody>
              {snapshot.index.components.map((component) => (
                <tr key={component.token}>
                  <td style={{ padding: "0.25rem 0" }}>{component.token}</td>
                  <td style={{ textAlign: "right" }}>{component.weight}</td>
                  <td style={{ textAlign: "right" }}>{component.drift}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      <div>
        <h3>SYN3800 Grants</h3>
        <p>
          Total grants: {snapshot.grants.summary.total}, completed: {snapshot.grants.summary.completed}
        </p>
        {snapshot.grants.records.length === 0 ? (
          <p>No grants recorded.</p>
        ) : (
          <table style={{ width: "100%", borderCollapse: "collapse" }}>
            <thead>
              <tr>
                <th style={{ textAlign: "left", borderBottom: "1px solid #d0d7de" }}>ID</th>
                <th style={{ textAlign: "left", borderBottom: "1px solid #d0d7de" }}>Beneficiary</th>
                <th style={{ textAlign: "right", borderBottom: "1px solid #d0d7de" }}>Amount</th>
                <th style={{ textAlign: "right", borderBottom: "1px solid #d0d7de" }}>Released</th>
                <th style={{ textAlign: "left", borderBottom: "1px solid #d0d7de" }}>Status</th>
              </tr>
            </thead>
            <tbody>
              {snapshot.grants.records.map((grant) => (
                <tr key={grant.id}>
                  <td>{grant.id}</td>
                  <td>{grant.beneficiary}</td>
                  <td style={{ textAlign: "right" }}>{grant.amount}</td>
                  <td style={{ textAlign: "right" }}>{grant.released}</td>
                  <td>{grant.status}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
      <div>
        <h3>SYN3900 Benefits</h3>
        <p>
          Total benefits: {snapshot.benefits.summary.total}, claimed: {snapshot.benefits.summary.claimed}
        </p>
        {snapshot.benefits.records.length === 0 ? (
          <p>No benefits registered.</p>
        ) : (
          <table style={{ width: "100%", borderCollapse: "collapse" }}>
            <thead>
              <tr>
                <th style={{ textAlign: "left", borderBottom: "1px solid #d0d7de" }}>ID</th>
                <th style={{ textAlign: "left", borderBottom: "1px solid #d0d7de" }}>Recipient</th>
                <th style={{ textAlign: "right", borderBottom: "1px solid #d0d7de" }}>Amount</th>
                <th style={{ textAlign: "left", borderBottom: "1px solid #d0d7de" }}>Status</th>
              </tr>
            </thead>
            <tbody>
              {snapshot.benefits.records.map((benefit) => (
                <tr key={benefit.id}>
                  <td>{benefit.id}</td>
                  <td>{benefit.recipient}</td>
                  <td style={{ textAlign: "right" }}>{benefit.amount}</td>
                  <td>{benefit.status}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
      <div>
        <h3>SYN4200 Charity Campaigns</h3>
        {(!snapshot.charity || snapshot.charity.length === 0) && <p>No campaigns recorded.</p>}
        {snapshot.charity && snapshot.charity.length > 0 && (
          <ul>
            {snapshot.charity.map((campaign) => (
              <li key={campaign.symbol}>
                {campaign.symbol}: raised {campaign.raised} / goal {campaign.goal} — {campaign.purpose}
              </li>
            ))}
          </ul>
        )}
      </div>
      <div>
        <h3>SYN4700 Legal Tokens</h3>
        {(!snapshot.legal || snapshot.legal.length === 0) && <p>No legal instruments registered.</p>}
        {snapshot.legal && snapshot.legal.length > 0 && (
          <table style={{ width: "100%", borderCollapse: "collapse" }}>
            <thead>
              <tr>
                <th style={{ textAlign: "left", borderBottom: "1px solid #d0d7de" }}>ID</th>
                <th style={{ textAlign: "left", borderBottom: "1px solid #d0d7de" }}>Owner</th>
                <th style={{ textAlign: "left", borderBottom: "1px solid #d0d7de" }}>Status</th>
                <th style={{ textAlign: "left", borderBottom: "1px solid #d0d7de" }}>Expires</th>
              </tr>
            </thead>
            <tbody>
              {snapshot.legal.map((token) => (
                <tr key={token.id}>
                  <td>{token.id}</td>
                  <td>{token.owner}</td>
                  <td>{token.status}</td>
                  <td>{formatDate(token.expiry)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
      {snapshot.utility && (
        <div>
          <h3>SYN500 Utility Token</h3>
          <p>
            {snapshot.utility.name} ({snapshot.utility.symbol}) owned by {snapshot.utility.owner}
          </p>
          <p>
            Total grants: {snapshot.utility.telemetry.grants},
            active: {snapshot.utility.telemetry.active},
            total usage: {snapshot.utility.telemetry.usage}
          </p>
        </div>
      )}
    </div>
  );
}

export default function Stage73Console() {
  const [snapshot, setSnapshot] = useState(null);
  const [digest, setDigest] = useState("--");
  const [updatedAt, setUpdatedAt] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);

  const loadSnapshot = useCallback(async () => {
    setLoading(true);
    try {
      const response = await fetch("/api/stage73");
      const data = await response.json();
      if (!response.ok) {
        throw new Error(data.error || data.detail || "Failed to load snapshot");
      }
      setSnapshot(data.snapshot);
      setDigest(data.digest || "--");
      setUpdatedAt(data.updatedAt || null);
      setError(null);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadSnapshot();
    const interval = setInterval(loadSnapshot, POLL_INTERVAL_MS);
    return () => clearInterval(interval);
  }, [loadSnapshot]);

  return (
    <main style={{ padding: "1.5rem", maxWidth: 1200, margin: "0 auto" }}>
      <h1>Stage 73 Enterprise Console</h1>
      <p>
        Manage the SYN3700 index, SYN3800 grant program, SYN3900 benefits, SYN4200 charity campaigns,
        SYN4700 legal instruments and SYN500 utility token from a single, fault-tolerant workspace.
      </p>
      <p>
        <a href="/">Return to the command catalog</a>
      </p>
      <Section title="Snapshot overview">
        <p>
          Current digest: <strong>{digest}</strong> (updated {formatDate(updatedAt)}).
          {" "}
          {loading && <span>Refreshing…</span>}
        </p>
        {error ? <p style={{ color: "#b00020" }}>{error}</p> : <SnapshotView snapshot={snapshot} />}
      </Section>

      <Section title="SYN3700 Index controls">
        <CommandForm
          title="Initialise index token"
          description="Create the enterprise index and register the primary controller wallet."
          command="syn3700"
          fields={[
            { name: "name", label: "Token name", required: true, placeholder: "Institutional Index" },
            { name: "symbol", label: "Token symbol", required: true, placeholder: "IDX" },
            {
              name: "controllerPath",
              label: "Controller wallet path",
              required: true,
              placeholder: "/var/synnergy/wallets/controller.wallet",
            },
            {
              name: "controllerPassword",
              label: "Controller wallet password",
              required: true,
              placeholder: "********",
              type: "password",
            },
          ]}
          buildArgs={(values) => [
            "init",
            "--name",
            values.name.trim(),
            "--symbol",
            values.symbol.trim(),
            "--controller",
            `${values.controllerPath.trim()}:${values.controllerPassword}`,
          ]}
          onSuccess={loadSnapshot}
        />
        <CommandForm
          title="Add index component"
          description="Register a component token with weighting and drift thresholds."
          command="syn3700"
          fields={[
            { name: "token", label: "Component token", required: true, placeholder: "AAA" },
            { name: "weight", label: "Weight", required: true, placeholder: "0.5" },
            { name: "drift", label: "Allowed drift", required: true, placeholder: "0.1" },
            {
              name: "walletPath",
              label: "Controller wallet path",
              required: true,
              placeholder: "/var/synnergy/wallets/controller.wallet",
            },
            {
              name: "walletPassword",
              label: "Controller wallet password",
              required: true,
              placeholder: "********",
              type: "password",
            },
          ]}
          buildArgs={(values) => [
            "add",
            values.token.trim(),
            "--weight",
            values.weight.trim(),
            "--drift",
            values.drift.trim(),
            "--wallet",
            values.walletPath.trim(),
            "--password",
            values.walletPassword,
          ]}
          onSuccess={loadSnapshot}
        />
      </Section>

      <Section title="SYN3800 Grant management">
        <CommandForm
          title="Create grant"
          description="Open a grant with optional authorizer wallet registration."
          command="syn3800"
          fields={[
            { name: "beneficiary", label: "Beneficiary address", required: true, placeholder: "wallet-address" },
            { name: "name", label: "Program name", required: true, placeholder: "STEM scholarship" },
            { name: "amount", label: "Total amount", required: true, placeholder: "1000" },
            {
              name: "authorizerPath",
              label: "Authorizer wallet path",
              required: false,
              placeholder: "/var/synnergy/wallets/controller.wallet",
              help: "Optional path:password pair to auto-register the first authorizer.",
            },
            {
              name: "authorizerPassword",
              label: "Authorizer wallet password",
              required: false,
              placeholder: "********",
              type: "password",
            },
          ]}
          buildArgs={(values) => {
            const args = [
              "create",
              values.beneficiary.trim(),
              values.name.trim(),
              values.amount.trim(),
            ];
            if (values.authorizerPath.trim() && values.authorizerPassword) {
              args.push("--authorizer", `${values.authorizerPath.trim()}:${values.authorizerPassword}`);
            }
            return args;
          }}
          onSuccess={loadSnapshot}
        />
        <CommandForm
          title="Release grant funds"
          description="Execute an authorised release with audit note support."
          command="syn3800"
          fields={[
            { name: "id", label: "Grant ID", required: true, placeholder: "1" },
            { name: "amount", label: "Amount to release", required: true, placeholder: "100" },
            { name: "note", label: "Audit note", required: false, placeholder: "Phase one milestone" },
            { name: "walletPath", label: "Authorizer wallet path", required: true, placeholder: "/path/to/wallet" },
            {
              name: "walletPassword",
              label: "Authorizer wallet password",
              required: true,
              placeholder: "********",
              type: "password",
            },
          ]}
          buildArgs={(values) => {
            const args = [
              "release",
              values.id.trim(),
              values.amount.trim(),
            ];
            if (values.note.trim()) {
              args.push(values.note.trim());
            }
            args.push("--wallet", values.walletPath.trim(), "--password", values.walletPassword);
            return args;
          }}
          onSuccess={loadSnapshot}
        />
        <CommandForm
          title="Authorize grant wallet"
          description="Grant release permissions to an additional controller wallet."
          command="syn3800"
          fields={[
            { name: "id", label: "Grant ID", required: true, placeholder: "1" },
            { name: "walletPath", label: "Authorizer wallet path", required: true, placeholder: "/path/to/wallet" },
            {
              name: "walletPassword",
              label: "Authorizer wallet password",
              required: true,
              placeholder: "********",
              type: "password",
            },
          ]}
          buildArgs={(values) => [
            "authorize",
            values.id.trim(),
            "--wallet",
            values.walletPath.trim(),
            "--password",
            values.walletPassword,
          ]}
          onSuccess={loadSnapshot}
        />
      </Section>

      <Section title="SYN3900 Benefit programs">
        <CommandForm
          title="Register benefit"
          description="Create a benefit allocation and optionally enroll an approver wallet."
          command="syn3900"
          fields={[
            { name: "recipient", label: "Recipient wallet", required: true, placeholder: "wallet-address" },
            { name: "program", label: "Program name", required: true, placeholder: "Housing assistance" },
            { name: "amount", label: "Benefit amount", required: true, placeholder: "500" },
            {
              name: "approverPath",
              label: "Approver wallet path",
              required: false,
              placeholder: "/var/synnergy/wallets/approver.wallet",
              help: "Optional path:password pair for the initial approver registration.",
            },
            {
              name: "approverPassword",
              label: "Approver wallet password",
              required: false,
              placeholder: "********",
              type: "password",
            },
          ]}
          buildArgs={(values) => {
            const args = [
              "register",
              values.recipient.trim(),
              values.program.trim(),
              values.amount.trim(),
            ];
            if (values.approverPath.trim() && values.approverPassword) {
              args.push("--approver", `${values.approverPath.trim()}:${values.approverPassword}`);
            }
            return args;
          }}
          onSuccess={loadSnapshot}
        />
        <CommandForm
          title="Claim benefit"
          description="Allow a recipient to claim their allocation with wallet attestation."
          command="syn3900"
          fields={[
            { name: "id", label: "Benefit ID", required: true, placeholder: "1" },
            { name: "walletPath", label: "Recipient wallet path", required: true, placeholder: "/path/to/wallet" },
            {
              name: "walletPassword",
              label: "Recipient wallet password",
              required: true,
              placeholder: "********",
              type: "password",
            },
          ]}
          buildArgs={(values) => [
            "claim",
            values.id.trim(),
            "--wallet",
            values.walletPath.trim(),
            "--password",
            values.walletPassword,
          ]}
          onSuccess={loadSnapshot}
        />
        <CommandForm
          title="Approve benefit"
          description="Approve a claimed benefit once governance checks complete."
          command="syn3900"
          fields={[
            { name: "id", label: "Benefit ID", required: true, placeholder: "1" },
            { name: "walletPath", label: "Approver wallet path", required: true, placeholder: "/path/to/wallet" },
            {
              name: "walletPassword",
              label: "Approver wallet password",
              required: true,
              placeholder: "********",
              type: "password",
            },
          ]}
          buildArgs={(values) => [
            "approve",
            values.id.trim(),
            "--wallet",
            values.walletPath.trim(),
            "--password",
            values.walletPassword,
          ]}
          onSuccess={loadSnapshot}
        />
      </Section>

      <Section title="SYN4200 Charity contributions">
        <CommandForm
          title="Record donation"
          description="Capture a donation to a registered campaign."
          command="syn4200_token"
          fields={[
            { name: "symbol", label: "Campaign symbol", required: true, placeholder: "HELP" },
            { name: "from", label: "Donor address", required: true, placeholder: "wallet-address" },
            { name: "amount", label: "Amount", required: true, placeholder: "25" },
            { name: "purpose", label: "Purpose / Note", required: false, placeholder: "Emergency relief" },
          ]}
          buildArgs={(values) => [
            "donate",
            values.symbol.trim(),
            "--from",
            values.from.trim(),
            "--amt",
            values.amount.trim(),
            "--purpose",
            values.purpose.trim(),
          ]}
          onSuccess={loadSnapshot}
        />
      </Section>

      <Section title="SYN4700 Legal instruments">
        <CommandForm
          title="Create legal token"
          description="Register a legal agreement with signatories and supply controls."
          command="syn4700"
          fields={[
            { name: "id", label: "Document ID", required: true, placeholder: "AGR-001" },
            { name: "name", label: "Document name", required: true, placeholder: "Service Agreement" },
            { name: "symbol", label: "Token symbol", required: true, placeholder: "AGR" },
            { name: "doctype", label: "Document type", required: true, placeholder: "contract" },
            { name: "hash", label: "Document hash", required: true, placeholder: "sha256:..." },
            { name: "owner", label: "Owning wallet", required: true, placeholder: "wallet-address" },
            {
              name: "expiry",
              label: "Expiry (ISO date)",
              required: true,
              placeholder: "2025-01-01T00:00:00Z",
              help: "ISO-8601 timestamp converted to seconds automatically.",
            },
            { name: "supply", label: "Token supply", required: true, placeholder: "100" },
            {
              name: "parties",
              label: "Parties (comma separated)",
              required: true,
              placeholder: "alice,bob",
            },
          ]}
          buildArgs={(values) => {
            const expiry = Date.parse(values.expiry);
            if (Number.isNaN(expiry)) {
              throw new Error("Invalid expiry timestamp");
            }
            const parties = values.parties
              .split(",")
              .map((entry) => entry.trim())
              .filter((entry) => entry.length > 0);
            if (parties.length === 0) {
              throw new Error("At least one party required");
            }
            const args = [
              "create",
              "--id",
              values.id.trim(),
              "--name",
              values.name.trim(),
              "--symbol",
              values.symbol.trim(),
              "--doctype",
              values.doctype.trim(),
              "--hash",
              values.hash.trim(),
              "--owner",
              values.owner.trim(),
              "--expiry",
              Math.floor(expiry / 1000).toString(),
              "--supply",
              values.supply.trim(),
            ];
            for (const party of parties) {
              args.push("--party", party);
            }
            return args;
          }}
          onSuccess={loadSnapshot}
        />
        <CommandForm
          title="Sign legal token"
          description="Attach a digital signature for a participating party."
          command="syn4700"
          fields={[
            { name: "id", label: "Document ID", required: true, placeholder: "AGR-001" },
            { name: "party", label: "Party", required: true, placeholder: "alice" },
            { name: "signature", label: "Signature", required: true, placeholder: "0xabc123" },
          ]}
          buildArgs={(values) => [
            "sign",
            values.id.trim(),
            values.party.trim(),
            values.signature.trim(),
          ]}
          onSuccess={loadSnapshot}
        />
        <CommandForm
          title="Log dispute"
          description="Record a dispute and optional resolution notes."
          command="syn4700"
          fields={[
            { name: "id", label: "Document ID", required: true, placeholder: "AGR-001" },
            { name: "reason", label: "Dispute reason", required: true, placeholder: "Breach of clause" },
            { name: "resolution", label: "Resolution status", required: true, placeholder: "pending" },
          ]}
          buildArgs={(values) => [
            "dispute",
            values.id.trim(),
            values.reason.trim(),
            values.resolution.trim(),
          ]}
          onSuccess={loadSnapshot}
        />
      </Section>

      <Section title="SYN500 Utility services">
        <CommandForm
          title="Create utility token"
          description="Launch the SYN500 token used to meter enterprise services."
          command="syn500"
          fields={[
            { name: "name", label: "Token name", required: true, placeholder: "Service Credit" },
            { name: "symbol", label: "Token symbol", required: true, placeholder: "UTL" },
            { name: "owner", label: "Owner wallet", required: true, placeholder: "wallet-address" },
            { name: "decimals", label: "Decimals", required: true, placeholder: "4" },
            { name: "supply", label: "Supply", required: true, placeholder: "100000" },
          ]}
          buildArgs={(values) => [
            "create",
            "--name",
            values.name.trim(),
            "--symbol",
            values.symbol.trim(),
            "--owner",
            values.owner.trim(),
            "--dec",
            values.decimals.trim(),
            "--supply",
            values.supply.trim(),
          ]}
          onSuccess={loadSnapshot}
        />
        <CommandForm
          title="Grant usage tier"
          description="Assign a service tier with usage limits and renewal window."
          command="syn500"
          fields={[
            { name: "address", label: "Wallet address", required: true, placeholder: "wallet-address" },
            { name: "tier", label: "Tier", required: true, placeholder: "1" },
            { name: "max", label: "Max usage", required: true, placeholder: "10" },
            { name: "window", label: "Window", required: false, placeholder: "1h" },
          ]}
          buildArgs={(values) => [
            "grant",
            values.address.trim(),
            "--tier",
            values.tier.trim(),
            "--max",
            values.max.trim(),
            "--window",
            values.window.trim() || "1h",
          ]}
          onSuccess={loadSnapshot}
        />
        <CommandForm
          title="Record usage"
          description="Log a unit of usage for a granted wallet and update telemetry."
          command="syn500"
          fields={[
            { name: "address", label: "Wallet address", required: true, placeholder: "wallet-address" },
          ]}
          buildArgs={(values) => ["use", values.address.trim()]}
          onSuccess={loadSnapshot}
        />
      </Section>
    </main>
  );
}
