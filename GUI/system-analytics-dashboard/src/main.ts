import config from '../config/production';

export async function main(): Promise<string> {
  try {
    const res = await fetch(`${config.apiUrl}/metrics`);
    if (!res.ok) throw new Error('request failed');
    const data = await res.json();
    return `System analytics: ${data.status}`;
  } catch {
    return 'System analytics: ERROR';
  }
}

if (require.main === module) {
  main()
    .then((output) => console.log(output))
    .catch((err) => {
      console.error('Failed to start system-analytics-dashboard', err);
      process.exit(1);
    });
}
