# Reversing and Cancelling Transactions

## Introduction
Neto Solaris designs the Synnergy Network to provide enterprises with programmable, controllable digital value transfers.  Beyond simple fund movement, organisations require mechanisms to halt or undo operations when errors, disputes, or regulatory interventions arise.  The network therefore supports two complementary capabilities: **cancelling scheduled transactions** and **reversing confirmed transactions**.  Both functions preserve ledger integrity while allowing the ecosystem to recover from mistakes or fraudulent activity.

## Scheduled Transactions and Cancellation
The platform allows a payment to be staged for execution at a future time.  A
`ScheduledTransaction` wraps a normal transaction with an `ExecuteAt` timestamp
and a mutable `Canceled` flag.

### Data Model and Behaviour
`ScheduledTransaction` consists of:

| Field | Description |
|-------|-------------|
| `Tx` | Pointer to the underlying transaction payload. |
| `ExecuteAt` | Absolute UNIX time at which the payment becomes valid. |
| `Canceled` | Boolean indicating whether the schedule has been revoked prior to execution. |

`ScheduleTransaction` merely creates this wrapper; orchestration services are
responsible for persisting and dispatching the transaction when `ExecuteAt`
elapses.  The cancellation logic is intentionally minimal: `CancelTransaction`
marks the record as cancelled when called before the deadline and returns a
boolean so callers can confirm success.  After the timestamp passes, attempts to
cancel gracefully fail.

### Workflow
1. **Scheduling** – `ScheduleTransaction` wraps a transaction with the target execution time.
2. **Cancellation window** – Prior to the scheduled time, participants may invoke `CancelTransaction`.  The routine returns `true` when the cancellation succeeds and marks the schedule as `Canceled`.
3. **Execution** – Uncancelled schedules execute at or after the specified timestamp.

### Interfaces
- **CLI** – The `tx control` command exposes `schedule` and `cancel` actions for
  operators who manage timed payments.
- **Virtual Machine** – SNVM opcodes `core_transaction_control_ScheduleTransaction`
  and `core_transaction_control_CancelTransaction` allow smart contracts to
  orchestrate timed transfers programmatically.

## Immediate Transaction Reversal
For transactions already committed to the ledger, a direct reversal is available
when the recipient retains the transferred amount.  `ReverseTransaction` works
as follows:

1. **Ledger lock** – the ledger mutex is acquired to isolate concurrent
   modifications.
2. **Balance check** – the recipient must hold at least the original amount; the
   fee is refunded from block rewards so only the amount is required on the
   recipient side.
3. **State adjustment** – the recipient’s balance is debited and the sender’s
   balance is credited with the amount plus fee.
4. **Unlock** – the mutex is released and the operation completes.

If the recipient lacks sufficient funds the routine aborts with an error,
preventing negative balances or monetary creation.

## Authority‑Mediated Reversals
In many jurisdictions a reversal requires oversight.  Synnergy introduces an authority‑driven workflow enabling trusted nodes to arbitrate disputes.

### Requesting a Reversal
`RequestReversal` freezes the recipient’s balance for the transaction amount plus an additional fee supplied by the requester.  This guarantees the funds remain available during the review process and records a `ReversalRequest` containing the transaction, fee and request time.

#### ReversalRequest Structure

| Field | Description |
|-------|-------------|
| `Tx` | The transaction under dispute. |
| `RequestedAt` | Timestamp when the reversal was initiated. |
| `Fee` | Return-gas amount paid by the requester and reimbursed if the reversal succeeds. |
| `votes` | Map of authority identifiers to boolean approval decisions. |

### Voting and Finalisation
Authority nodes cast decisions using `Vote`.  When `FinalizeReversal` observes
the required approval threshold, it executes a compensating transfer from the
frozen funds back to the sender and credits the fee to the recipient to cover
gas costs.  The operation must complete within a 30‑day `reversalWindow`; if the
deadline elapses `RejectReversal` releases the frozen balance and records that
the request expired.

#### Processing Sequence
1. **Freeze** – `RequestReversal` deducts the amount+fee from the recipient and
   places it in the ledger’s `frozen` map.
2. **Vote** – authorities register approvals or rejections via `Vote`.
3. **Finalize** – once approvals meet the required quorum, `FinalizeReversal`
   re-credits the recipient’s frozen balance, constructs a compensating
   transaction and applies it to reimburse the sender.
4. **Reject** – if quorum is not reached within 30 days, `RejectReversal`
   restores the frozen funds to the recipient and marks the request failed.

#### Example Authority Workflow
```go
req, err := core.RequestReversal(ledger, tx, fee)
if err != nil { /* handle */ }
req.Vote("authority-1", true)
req.Vote("authority-2", true)
if err := core.FinalizeReversal(ledger, req, 2); err != nil {
    // reversal failed or quorum not met
}
```

## Ledger State and Fund Freezing
The ledger maintains both spendable balances and a `frozen` map that holds amounts reserved for pending reversals.  All operations lock the ledger’s state to maintain consistency under concurrency.  By segregating frozen funds the system prevents double‑spending while a dispute is active.

## Security and Governance Considerations
- **Multi‑party approval** – Reversal finalisation requires an explicit quorum of
  authority votes, providing checks and balances against unilateral rollbacks.
- **Time‑boxed review** – The 30‑day window motivates timely adjudication and
  avoids indefinite freezes.
- **Auditability** – Every action—scheduling, cancelling, requesting reversals,
  voting, and finalisation—produces deterministic state transitions suitable for
  auditing.
- **Error handling** – Both scheduling and reversal APIs return explicit error
  values enabling callers to react programmatically to insufficient balances or
  expired requests.

## Edge Cases and Failure Handling
- **Expired cancellation** – Calling `CancelTransaction` after `ExecuteAt`
  results in a graceful failure, leaving the schedule intact for execution.
- **Insufficient funds** – `ReverseTransaction` and `RequestReversal` both abort
  when the recipient lacks the required balance, ensuring the ledger never
  creates value.
- **Unmet quorum** – If the approval threshold is not reached within the
  `reversalWindow`, `RejectReversal` automatically releases frozen funds and the
  request is marked failed.

## Usage Examples
```bash
# Schedule a future transfer
synnergy tx control schedule alice bob 100 1 0 $(date -d "+1 hour" +%s)

# Cancel the scheduled transfer before execution
synnergy tx control cancel alice bob 100 1 0 $(date -d "+1 hour" +%s)

# Apply and immediately reverse a transaction on a fresh ledger
synnergy tx control reverse alice bob 100 1 0
```
Smart contracts and off‑chain services may call the corresponding SNVM opcodes
to achieve the same outcomes within automated workflows.

## Related Transaction Controls
- **Private transactions** – `ConvertToPrivate` encrypts a transaction using
  AES‑GCM so only holders of the key can inspect details.  `PrivateTransaction`
  objects can be decrypted and applied after inspection by the intended
  counterparties.
- **Receipts** – `GenerateReceipt` and `ReceiptStore` provide verifiable records
  of transaction outcomes, allowing auditors to trace reversals or cancellations
  long after execution.

## Best Practices
- Confirm recipient details before scheduling transfers to minimise reversal requests.
- Use cancellation for routine corrections; reserve authority‑mediated reversals for exceptional cases such as fraud or regulatory orders.
- Monitor outstanding reversal requests to ensure they conclude within the allowed window.

## Conclusion
Through these mechanisms Neto Solaris equips the Synnergy Network with comprehensive transaction control.  Organisations can schedule payments confidently, revoke them when circumstances change, and rely on a governed reversal process when settlement errors must be rectified.  The combination of programmability and oversight delivers a resilient financial substrate for modern digital commerce.
