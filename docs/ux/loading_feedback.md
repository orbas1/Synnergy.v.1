# Stage 129: Loading Indicators and Progress Feedback

Commands may take several seconds to execute and currently leave the screen frozen. Clear feedback reassures users that work is in progress and helps them distinguish between slow operations and failures.

## Techniques
- **Spinners or skeleton screens** for short operations such as fetching help text.
- **Progress bars** that advance as the backend reports completion percentages.
- **Disable repeated actions** by greying out buttons until the command completes.

## Implementation Tips
1. Emit progress events from the server at logical checkpoints and consume them through WebSockets or Server‑Sent Events.
2. Provide textual updates alongside visual indicators, e.g., "Executing…", "Broadcasting transaction…", "Waiting for block confirmation…".
3. Ensure indicators are accessible by announcing state changes to screen readers via `aria-live` regions.

## Checklist
- [ ] Every network request shows a spinner or progress bar within 500 ms.
- [ ] The interface prevents duplicate submissions while a command is running.
- [ ] Users receive a clear success or failure message when execution finishes.

By setting expectations and preventing accidental re‑runs, loading feedback keeps the user engaged and protects backend resources.
