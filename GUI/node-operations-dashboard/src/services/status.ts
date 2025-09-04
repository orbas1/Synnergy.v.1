export interface NodeStatusResponse {
  status: string;
}

/**
 * Fetches node status from the API defined by STATUS_URL.
 * Falls back to a local endpoint and returns `ERROR` on network issues.
 */
export async function fetchNodeStatus(
  url: string = process.env.STATUS_URL || 'http://localhost:8080/status'
): Promise<string> {
  try {
    const res = await fetch(url);
    if (!res.ok) {
      throw new Error(`Request failed with status ${res.status}`);
    }
    const data = (await res.json()) as Partial<NodeStatusResponse>;
    return data.status ?? 'UNKNOWN';
  } catch (err) {
    // Log for observability but do not leak internals to caller
    console.error('fetchNodeStatus error', err);
    return 'ERROR';
  }
}
