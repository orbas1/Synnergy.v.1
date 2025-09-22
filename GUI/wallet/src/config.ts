export interface WalletGuiConfig {
  apiUrl: string;
  requestTimeoutMs: number;
  maxRetries: number;
  defaultCurrency: string;
}

const DEFAULT_CONFIG: WalletGuiConfig = {
  apiUrl: 'http://127.0.0.1:4000',
  requestTimeoutMs: 5000,
  maxRetries: 2,
  defaultCurrency: 'SYN'
};

function toPositiveInteger(value: string | undefined, fallback: number): number {
  if (!value) {
    return fallback;
  }

  const parsed = Number(value);
  if (Number.isFinite(parsed) && parsed > 0) {
    return Math.floor(parsed);
  }

  return fallback;
}

export function resolveConfig(overrides: Partial<WalletGuiConfig> = {}): WalletGuiConfig {
  const fromEnv: Partial<WalletGuiConfig> = {
    apiUrl: process.env.WALLET_API_URL || process.env.API_URL,
    requestTimeoutMs: toPositiveInteger(process.env.WALLET_REQUEST_TIMEOUT_MS, DEFAULT_CONFIG.requestTimeoutMs),
    maxRetries: toPositiveInteger(process.env.WALLET_REQUEST_RETRIES, DEFAULT_CONFIG.maxRetries),
    defaultCurrency: process.env.WALLET_CURRENCY
  };

  return {
    ...DEFAULT_CONFIG,
    ...fromEnv,
    ...overrides
  };
}

export const walletConfig = resolveConfig();
