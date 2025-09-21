import { useEffect, useState } from "react";

export default function Home() {
  const [commands, setCommands] = useState([]);
  const [selected, setSelected] = useState("");
  const [flags, setFlags] = useState([]);
  const [inputs, setInputs] = useState({});
  const [extra, setExtra] = useState("");
  const [output, setOutput] = useState("");

  useEffect(() => {
    fetch("/api/commands")
      .then((res) => res.json())
      .then((data) => setCommands(data.commands || []));
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

  return (
    <main style={{ padding: 20 }}>
      <h1>Synnergy Control Panel</h1>
      <section style={{ marginBottom: 24 }}>
        <h2>Opportunity Orb Overview</h2>
        <p>
          The orb showcases Synnergy&apos;s emerging professional ecosystem—a global map of
          potential collaborations we are curating before formally onboarding clients. Each
          node represents a prospective partnership, specialist exchange, or market pathway
          that demonstrates how talent and expertise can intersect within the Synnergy
          network.
        </p>
        <ul>
          <li>
            <strong>Worldwide perspective.</strong> The constantly shifting layout mirrors the
            breadth of industries, time zones, and disciplines we are preparing to connect
            across the globe.
          </li>
          <li>
            <strong>Dynamic connection cues.</strong> Wobbling nodes and luminous trails
            illustrate opportunities gaining momentum, signaling where meaningful
            professional interactions can take shape next.
          </li>
          <li>
            <strong>Responsive clarity.</strong> The lightweight canvas adapts instantly to
            desktops, tablets, and phones so stakeholders can explore the opportunity space
            from any device.
          </li>
        </ul>
      </section>
      <p>Select a command, fill in optional flags, and run.</p>
      <p>
        <a href="/regnode">Regulatory node console</a>
      </p>
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
