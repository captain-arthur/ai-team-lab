# Lessons

Store **lessons learned** — what worked, what didn’t, and what to do differently next time. This is the team’s retrospective memory.

---

## What to put here

- **What worked** — practices or choices that paid off and should be repeated.
- **What didn’t** — pitfalls, dead ends, or trade-offs that didn’t work out.
- **What to do differently** — concrete changes for future projects (process, scope, or tools).

---

## Format

- One file per project or theme (e.g. `project-xyz-lessons.md`, `research-phase-lessons.md`).
- Keep entries short and actionable. Avoid blame; focus on behavior and decisions.
- Optionally add **Source** (project or phase) and **Date** so context is clear later.

---

## Example

```markdown
# Lessons: Kubernetes CAT project

- **Worked:** Starting with a small set of checks and expanding; reusing sonobuoy where possible.
- **Didn’t:** Assuming one tool would cover all environments without checking versions first.
- **Next time:** Lock tool versions in the design doc; add a "supported versions" section early.
```
