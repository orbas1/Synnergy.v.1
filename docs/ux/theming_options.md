# Stage 130: Theming and Customization

A single light theme can cause eye strain and fails to meet user preferences. Offering configurable themes improves usability and brand alignment.

## Theme Variants
- **Light and dark modes** as the primary options.
- **High‑contrast theme** for accessibility.
- Allow custom color palettes through CSS variables for organizations that wish to rebrand the panel.

## Implementation Details
1. Define a base set of CSS variables for colors, spacing and typography.
2. Provide alternate theme files that override the variables and toggle them by adding a class on the `<body>` element.
3. Store the selected theme in local storage and respect the user's system preference via `prefers-color-scheme`.

## User Controls
- Include a theme switcher in the header or settings panel.
- Offer live preview when customizing colors and fonts.
- Expose an import/export option so teams can share theme presets.

## Checklist
- [ ] Theme choices persist across sessions.
- [ ] Custom themes never reduce text contrast below WCAG guidelines.
- [ ] Documentation explains how to add new theme variants.

Flexible theming makes the CLI panel comfortable in any environment and adaptable to diverse branding requirements.
