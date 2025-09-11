import { useEffect, useState } from 'react';

export default function DAO() {
  const [daos, setDaos] = useState([]);
  const [form, setForm] = useState({ dao: '', admin: '', addr: '', role: '' });
  const [status, setStatus] = useState('');
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
      <h2>Update Member Role</h2>
      <form onSubmit={e => {
        e.preventDefault();
        fetch('/api/run', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ command: 'dao-members', args: ['update', form.dao, form.admin, form.addr, form.role, '--json'] })
        })
          .then(res => res.json())
          .then(data => setStatus(data.output || ''));
      }}>
        <input placeholder="DAO ID" value={form.dao} onChange={e => setForm({ ...form, dao: e.target.value })} />
        <input placeholder="Admin" value={form.admin} onChange={e => setForm({ ...form, admin: e.target.value })} />
        <input placeholder="Member" value={form.addr} onChange={e => setForm({ ...form, addr: e.target.value })} />
        <input placeholder="Role" value={form.role} onChange={e => setForm({ ...form, role: e.target.value })} />
        <button type="submit">Update</button>
      </form>
      {status && <pre>{status}</pre>}
    </main>
  );
}

