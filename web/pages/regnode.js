import { useState } from "react";

export default function RegNode() {
  const [from, setFrom] = useState("");
  const [amount, setAmount] = useState("");
  const [flagAddr, setFlagAddr] = useState("");
  const [flagReason, setFlagReason] = useState("");
  const [logAddr, setLogAddr] = useState("");
  const [output, setOutput] = useState("");

  const run = async (args) => {
    const res = await fetch("/api/run", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ command: "regnode", args }),
    });
    const data = await res.json();
    setOutput(data.output || data.error);
  };

  return (
    <main style={{ padding: 20 }}>
      <h1>Regulatory Node Console</h1>
      <section style={{ marginBottom: 20 }}>
        <h2>Approve Transaction</h2>
        <input
          placeholder="from"
          value={from}
          onChange={(e) => setFrom(e.target.value)}
        />
        <input
          placeholder="amount"
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
          style={{ marginLeft: 10 }}
        />
        <button onClick={() => run(["approve", from, amount])}>Approve</button>
      </section>
      <section style={{ marginBottom: 20 }}>
        <h2>Flag Address</h2>
        <input
          placeholder="address"
          value={flagAddr}
          onChange={(e) => setFlagAddr(e.target.value)}
        />
        <input
          placeholder="reason"
          value={flagReason}
          onChange={(e) => setFlagReason(e.target.value)}
          style={{ marginLeft: 10 }}
        />
        <button onClick={() => run(["flag", flagAddr, flagReason])}>
          Flag
        </button>
      </section>
      <section style={{ marginBottom: 20 }}>
        <h2>View Logs</h2>
        <input
          placeholder="address"
          value={logAddr}
          onChange={(e) => setLogAddr(e.target.value)}
        />
        <button
          onClick={() => run(["logs", logAddr])}
          style={{ marginLeft: 10 }}
        >
          Fetch Logs
        </button>
      </section>
      <pre>{output}</pre>
    </main>
  );
}
