import { HttpClient, HttpError } from './httpClient';

describe('HttpClient', () => {
  const fetchMock = jest.fn();
  let originalFetch: typeof fetch | undefined;

  beforeEach(() => {
    originalFetch = global.fetch;
    (global as unknown as { fetch: typeof fetch }).fetch = fetchMock as unknown as typeof fetch;
    fetchMock.mockReset();
  });

  afterEach(() => {
    if (originalFetch) {
      (global as unknown as { fetch: typeof fetch }).fetch = originalFetch;
    } else {
      // @ts-expect-error cleanup mock fetch when not present
      delete global.fetch;
    }
  });

  function createResponse(body: unknown, init?: Partial<Response>): Response {
    return {
      ok: init?.ok ?? true,
      status: init?.status ?? 200,
      statusText: init?.statusText ?? 'OK',
      text: async () => (body === undefined ? '' : JSON.stringify(body))
    } as unknown as Response;
  }

  it('performs GET requests and parses JSON', async () => {
    fetchMock.mockResolvedValueOnce(createResponse({ value: 1 }));

    const client = new HttpClient({ baseUrl: 'http://example.com', timeoutMs: 100, maxRetries: 0 });
    const result = await client.get<{ value: number }>('/wallet');

    expect(fetchMock).toHaveBeenCalledWith('http://example.com/wallet', expect.objectContaining({ method: 'GET' }));
    expect(result.value).toBe(1);
  });

  it('retries failed requests before succeeding', async () => {
    fetchMock
      .mockRejectedValueOnce(new Error('network'))
      .mockResolvedValueOnce(createResponse({ ok: true }));

    const client = new HttpClient({ baseUrl: 'http://example.com', timeoutMs: 10, maxRetries: 1 });
    const result = await client.get<{ ok: boolean }>('/retry');

    expect(fetchMock).toHaveBeenCalledTimes(2);
    expect(result.ok).toBe(true);
  });

  it('throws HttpError for non-OK responses', async () => {
    fetchMock.mockResolvedValueOnce(createResponse('boom', { ok: false, status: 500, statusText: 'Server Error' }));

    const client = new HttpClient({ baseUrl: 'http://example.com', timeoutMs: 10, maxRetries: 0 });

    await expect(client.get('/fail')).rejects.toBeInstanceOf(HttpError);
  });
});
