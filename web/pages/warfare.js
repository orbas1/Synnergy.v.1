import { useEffect, useState } from "react";

export default function WarfareConsole() {
  const [command, setCommand] = useState("");
  const [commander, setCommander] = useState("");
  const [privateKey, setPrivateKey] = useState("");
  const [asset, setAsset] = useState("");
  const [location, setLocation] = useState("");
  const [status, setStatus] = useState("");
  const [tactical, setTactical] = useState("");
  const [events, setEvents] = useState([]);
  const [since, setSince] = useState(0);
  const [output, setOutput] = useState("");

  const run = async (subcommand, args) => {
    const res = await fetch("/api/run", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ command: "warfare", args: [subcommand, ...args] }),
    });
    const data = await res.json();
    setOutput(data.output || data.error || "");
    return data;
  };

  const executeCommand = async () => {
    const args = ["command", command];
    if (commander) {
      args.push("--commander", commander);
    }
    if (privateKey) {
      args.push("--private", privateKey);
    }
    const data = await run(args[0], args.slice(1));
    if (data.output) {
      setOutput(data.output);
    }
  };

  const trackLogistics = async () => {
    const args = ["track", asset, location, status];
    const data = await run(args[0], args.slice(1));
    if (data.output) {
      setOutput(data.output);
    }
  };

  const shareTactical = async () => {
    const data = await run("share", [tactical]);
    if (data.output) {
      setOutput(data.output);
    }
  };

  const fetchEvents = async () => {
    const args = ["events"];
    if (since > 0) {
      args.push("--since", String(since));
    }
    const result = await run(args[0], args.slice(1));
    if (result.output) {
      try {
        const parsed = JSON.parse(result.output);
        setEvents(parsed);
      } catch (err) {
        setEvents([{ error: result.output }]);
      }
    }
  };

  useEffect(() => {
    fetchEvents();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <main style={{ padding: 20 }}>
      <h1>Warfare Node Console</h1>
      <section style={{ marginBottom: 20 }}>
        <h2>Execute Command</h2>
        <input
          placeholder="command"
          value={command}
          onChange={(e) => setCommand(e.target.value)}
        />
        <input
          placeholder="commander (optional)"
          value={commander}
          onChange={(e) => setCommander(e.target.value)}
          style={{ marginLeft: 10 }}
        />
        <input
          placeholder="private key (hex)"
          value={privateKey}
          onChange={(e) => setPrivateKey(e.target.value)}
          style={{ marginLeft: 10 }}
        />
        <button onClick={executeCommand} style={{ marginLeft: 10 }}>
          Execute
        </button>
      </section>

      <section style={{ marginBottom: 20 }}>
        <h2>Track Logistics</h2>
        <input
          placeholder="asset"
          value={asset}
          onChange={(e) => setAsset(e.target.value)}
        />
        <input
          placeholder="location"
          value={location}
          onChange={(e) => setLocation(e.target.value)}
          style={{ marginLeft: 10 }}
        />
        <input
          placeholder="status"
          value={status}
          onChange={(e) => setStatus(e.target.value)}
          style={{ marginLeft: 10 }}
        />
        <button onClick={trackLogistics} style={{ marginLeft: 10 }}>
          Record
        </button>
      </section>

      <section style={{ marginBottom: 20 }}>
        <h2>Broadcast Tactical Update</h2>
        <input
          placeholder="message"
          value={tactical}
          onChange={(e) => setTactical(e.target.value)}
        />
        <button onClick={shareTactical} style={{ marginLeft: 10 }}>
          Share
        </button>
      </section>

      <section style={{ marginBottom: 20 }}>
        <h2>Events</h2>
        <div>
          <label>
            Since sequence
            <input
              type="number"
              value={since}
              onChange={(e) => setSince(Number(e.target.value))}
              style={{ marginLeft: 10 }}
            />
          </label>
          <button onClick={fetchEvents} style={{ marginLeft: 10 }}>
            Refresh
          </button>
        </div>
        <pre style={{ background: "#111", color: "#0f0", padding: 10 }}>
          {JSON.stringify(events, null, 2)}
        </pre>
      </section>

      <section>
        <h2>Raw Output</h2>
        <pre>{output}</pre>
      </section>
    </main>
  );
}
