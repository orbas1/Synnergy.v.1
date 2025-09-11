# Error Catalogue

Common error codes and messages returned by Synnergy components. Each error implements Go's `error` interface and may wrap lower‑level causes.

| Code | Message | Description | Retry |
|------|---------|-------------|------|
| `ERR_INVALID_TX` | invalid transaction | Transaction failed basic validation (nonce, signature, funds). | no |
| `ERR_UNAUTHORIZED` | unauthorized | Caller lacks required permissions or role. | no |
| `ERR_NOT_FOUND` | not found | Requested resource does not exist in the ledger or registry. | no |
| `ERR_DUPLICATE` | duplicate entry | An item already exists with the same identifier. | yes |
| `ERR_NETWORK` | network error | Temporary network failure; peers may be unreachable. | yes |
| `ERR_TIMEOUT` | operation timed out | Operation exceeded its allotted time. | yes |
| `ERR_INTERNAL` | internal error | Unexpected failure in a module; check logs for details. | no |

## Usage

Errors follow the pattern `errors.Wrap(code, fmt.Errorf("detail: %w", err))` so callers can match on the code while retaining context.

```go
if errors.Code(err) == errors.ERR_INVALID_TX {
    // handle accordingly
}
```

See `pkg/errors` for helper utilities and additional codes.

