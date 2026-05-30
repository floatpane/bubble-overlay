# Security Policy

## Supported Versions

Only the latest release of bubble-overlay is supported with security updates.

## Reporting a Vulnerability

If you discover a security vulnerability in bubble-overlay, please report it responsibly. **Do not open a public issue.**

Email us at [us@floatpane.com](mailto:us@floatpane.com) with:

- A description of the vulnerability
- Steps to reproduce the issue
- The potential impact
- Any suggested fixes (optional)

We will acknowledge your report within 48 hours and aim to provide a fix or mitigation plan within 7 days, depending on severity.

## Scope

This policy covers the bubble-overlay codebase and its official releases.

Of particular interest:

- Terminal escape-sequence injection through `Block` / `Line` output.
- Crashes or runaway memory on crafted inputs (very wide overlays, malformed ANSI in the base string, extreme row/col values).
- ANSI bleed-through where overlay styles affect surrounding cells.

Third-party dependencies are outside our direct control, but we will work to address reported issues in them as quickly as possible.

## Disclosure

We ask that you give us reasonable time to address the issue before disclosing it publicly. We are committed to crediting reporters in release notes (unless you prefer to remain anonymous).
