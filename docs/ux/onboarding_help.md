# Stage 131: Onboarding and Contextual Guidance

New contributors often land on the web control panel with little understanding of the CLI syntax that powers it. A purposeful onboarding flow shortens the learning curve and prevents misconfiguration.

## Goals
- Introduce the command model and required flags.
- Demonstrate a successful round‑trip command before the user explores advanced options.
- Reduce support requests by answering common questions in context.

## Recommended Patterns
- **Interactive tour** – A step‑by‑step overlay that highlights form fields and explains their purpose.
- **Inline help** – Tooltips and "more info" links beside every flag or argument.
- **Starter templates** – Pre‑filled examples for popular tasks so the first run is copy‑free.

## Implementation Tips
1. Trigger the tour automatically on first visit and allow users to relaunch it from the help menu.
2. Persist completion state in local storage so returning users are not forced through the tutorial.
3. Provide a prominent link to full documentation and community channels for deeper learning.

## Checklist
- [ ] All mandatory fields include helper text and validation examples.
- [ ] The tour can be navigated with keyboard controls for accessibility.
- [ ] Users can skip or exit onboarding at any time without losing progress.

Effective onboarding turns the CLI panel from a wall of options into an inviting toolkit.
