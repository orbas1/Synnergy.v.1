import { useEffect, useState } from 'react';

export default function DAO() {
  const [daos, setDaos] = useState([]);
  useEffect(() => {
    fetch('/api/run', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ command: 'dao', args: ['list', '--json'] })
    })
      .then(res => res.json())
      .then(data => {
        try {
          const parsed = JSON.parse(data.output || '[]');
          setDaos(parsed);
        } catch {
          setDaos([]);
        }
      });
  }, []);
  return (
    <main style={{ padding: 20 }}>
      <h1>Decentralised Autonomous Organisations</h1>
      <ul>
        {daos.map(d => (
          <li key={d.ID}>{d.ID} - {d.Name}</li>
        ))}
      </ul>
    </main>
  );
}

