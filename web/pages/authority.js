import { useEffect, useState } from 'react';

export default function Authority() {
  const [nodes, setNodes] = useState([]);
  useEffect(() => {
    fetch('/api/run', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ command: 'authority', args: ['list', '--json'] })
    })
      .then(res => res.json())
      .then(data => {
        try {
          const parsed = JSON.parse(data.output || '[]');
          setNodes(parsed);
        } catch {
          setNodes([]);
        }
      });
  }, []);
  return (
    <main style={{ padding: 20 }}>
      <h1>Authority Nodes</h1>
      <ul>
        {nodes.map(n => (
          <li key={n.address}>{n.address} - {n.role} ({n.votes} votes)</li>
        ))}
      </ul>
    </main>
  );
}
