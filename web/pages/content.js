import { useEffect, useState } from 'react';

export default function Content() {
  const [items, setItems] = useState([]);
  useEffect(() => {
    fetch('/api/run', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ command: 'content_node', args: ['list', '--json'] })
    })
      .then(res => res.json())
      .then(data => {
        try {
          const parsed = JSON.parse(data.output || '[]');
          setItems(parsed);
        } catch {
          setItems([]);
        }
      });
  }, []);
  return (
    <main style={{ padding: 20 }}>
      <h1>Content Nodes</h1>
      <ul>
        {items.map(m => (
          <li key={m.id}>{m.id} - {m.name}</li>
        ))}
      </ul>
    </main>
  );
}
