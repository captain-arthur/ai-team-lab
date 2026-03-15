# Patterns

Store **reusable solution patterns** — how the team approaches recurring problems so future projects can copy or adapt them.

---

## What to put here

- **Architecture patterns** — e.g. "single-file config for sync-friendly apps", "layer boundaries for our services".
- **Process patterns** — e.g. "how we do acceptance tests for new clusters", "how we document API changes".
- **Implementation patterns** — recurring code or config structures that work well.

---

## Format

- One file per pattern or theme (e.g. `cluster-acceptance-testing.md`, `plain-text-storage.md`).
- Describe: problem, approach, steps or structure, and when to use (or not use) it.
- Optionally add **Source** (project that produced this pattern).

---

## Example

```markdown
# Cluster acceptance testing

- **Problem:** Verify a new or updated cluster before production.
- **Approach:** Define a small set of checks (API, nodes, DNS, storage), run them from a script or harness, report pass/fail.
- **When to use:** After provisioning or major upgrades. Prefer existing tools (e.g. sonobuoy) when they fit.
```
