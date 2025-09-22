export interface HttpClientOptions {
  baseUrl: string;
  timeoutMs: number;
  maxRetries: number;
  headers?: Record<string, string>;
}

export class HttpError extends Error {
  constructor(public readonly status: number, public readonly body: string) {
    super(`HTTP ${status}: ${body || 'Request failed'}`);
    this.name = 'HttpError';
  }
}

const DEFAULT_HEADERS = {
  Accept: 'application/json',
  'Content-Type': 'application/json'
} as const;

export class HttpClient {
  private readonly headers: Record<string, string>;

  constructor(private readonly options: HttpClientOptions) {
    this.headers = {
      ...DEFAULT_HEADERS,
      ...(options.headers ?? {})
    };
  }

  async get<T>(path: string): Promise<T> {
    return this.request<T>('GET', path);
  }

  async post<T>(path: string, body?: unknown): Promise<T> {
    return this.request<T>('POST', path, body);
  }

  private buildUrl(path: string): string {
    return new URL(path, this.options.baseUrl).toString();
  }

  private async request<T>(method: string, path: string, body?: unknown, attempt = 0): Promise<T> {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), this.options.timeoutMs);

    try {
      const response = await fetch(this.buildUrl(path), {
        method,
        body: body === undefined ? undefined : JSON.stringify(body),
        headers: this.headers,
        signal: controller.signal
      });

      if (!response.ok) {
        const errorBody = await response.text();
        throw new HttpError(response.status, errorBody || response.statusText);
      }

      if (response.status === 204) {
        return undefined as T;
      }

      const text = await response.text();
      if (!text) {
        return undefined as T;
      }

      return JSON.parse(text) as T;
    } catch (error) {
      if (error instanceof HttpError) {
        throw error;
      }

      if (attempt < this.options.maxRetries) {
        const backoffMs = Math.min(this.options.timeoutMs, 1000) * (attempt + 1);
        await new Promise((resolve) => setTimeout(resolve, backoffMs));
        return this.request<T>(method, path, body, attempt + 1);
      }

      if (error instanceof Error) {
        throw new Error(`HTTP request failed: ${error.message}`);
      }

      throw new Error('HTTP request failed due to an unknown error');
    } finally {
      clearTimeout(timeout);
    }
  }
}
