---
area: security
phase: v0.1-blocker
---
# Make generated output safe for its HTML context

## Goal
Prevent values that are escaped as text from remaining dangerous in URL,
attribute, script, or style contexts.

## Acceptance criteria
- URL attributes reject unsafe schemes such as `javascript:` by default.
- Text, attribute, URL, script, and style contexts have explicit safe-value rules.
- Raw or trusted values require a visible, documented opt-in.
- Runtime HTTP tests cover malicious values in every supported context.
- Public security claims describe the actual guarantees and limitations.
