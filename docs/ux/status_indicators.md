# Stage 135: System Status Indicators

Administrators need quick insight into whether the panel is connected and commands are executing correctly. Adding real‑time status indicators builds trust and simplifies troubleshooting.

## Key Indicators
- **Backend connectivity** – Display a green/red icon showing whether the server responds to health pings.
- **Command execution results** – Surface success, warning or failure banners after every run.
- **Node metrics** – Optionally show block height, peer count or latency for deeper diagnostics.

## Implementation Strategy
1. Poll a lightweight `/health` endpoint at regular intervals and update an icon in the header.
2. Capture exit codes from executed commands and map them to color‑coded notifications.
3. When metrics are available, stream them over WebSockets and present a compact dashboard or tooltip.

## Checklist
- [ ] Connectivity indicator updates within five seconds of a network change.
- [ ] Notifications include actionable messages and links to logs.
- [ ] Indicators respect high‑contrast and reduced‑motion preferences.

Consistent status feedback reassures operators that the system is functioning and highlights issues before they escalate.
