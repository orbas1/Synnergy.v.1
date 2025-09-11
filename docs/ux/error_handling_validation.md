# Stage 128: Error Handling and Input Validation

Users currently receive little feedback when a command fails or when required flags are missing. Robust validation prevents mistakes and makes recovery intuitive.

## Client‑Side Validation
- Mark required fields and enforce format checks before submission.
- Provide inline messages near the offending field instead of generic alerts.
- Use pattern matching or dropdowns to restrict inputs to valid values.

## Server‑Side Safeguards
- Re‑validate all parameters on the server and return structured error objects.
- Include remediation hints such as "check available balance" or "use hex‑encoded address".
- Log validation failures for auditing and to improve future UI hints.

## Error Display Patterns
- Show a summary banner for fatal errors and highlight specific fields below.
- Offer retry actions when possible and preserve user input to avoid re‑typing.
- For long‑running commands, stream partial output so users can see where it failed.

## Checklist
- [ ] Every form field has validation rules with descriptive error messages.
- [ ] Backend errors are surfaced to the UI with human‑readable explanations.
- [ ] Users can recover from errors without refreshing the page.

Clear validation and error handling turn confusing failures into learning opportunities and smoother workflows.
