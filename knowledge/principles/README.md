# Principles

Store **design principles, decision rules, and standards** that the AI team (and you) want to apply across projects.

---

## What to put here

- **Design principles** — e.g. "prefer simple over clever", "document decisions and rationale", "fail fast and visibly".
- **Decision rules** — e.g. "when choosing between build vs buy, prefer buy if an OSS option exists and is maintained".
- **Standards** — coding style, naming, or process rules that should be respected in future work.

---

## Format

- One file per theme or domain (e.g. `security.md`, `api-design.md`, `simplicity.md`).
- Use clear headings and short bullets so roles can scan quickly when starting a new task.
- Optionally add a **Source** line linking to the project that produced the principle.

---

## Example

```markdown
# API design principles

- Prefer small, focused endpoints over large multipurpose ones.
- Document errors and status codes in one place.
- Version in the URL path (e.g. /v1/...) from day one.
```
