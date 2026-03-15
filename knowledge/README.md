# Knowledge Memory System

The **knowledge** directory is the shared memory of the AI team. Reusable insights from completed projects are extracted and stored here so future tasks can benefit from what the team has learned.

---

## Directory structure

| Directory | Purpose |
|-----------|---------|
| **principles/** | Design principles, decision rules, and standards that should inform future work. |
| **tools/** | Tool evaluations, recommendations, and usage notes (libraries, CLIs, platforms). |
| **patterns/** | Reusable solution patterns — how we solve recurring problems (e.g. auth, config, testing). |
| **lessons/** | Lessons learned: what worked, what didn’t, and what to do differently next time. |

Each subdirectory has its own README describing what to store and how to name files.

---

## When to use

- **After a project:** Run the **Knowledge Extraction** stage (see **WORKFLOW.md**). Use `templates/knowledge-extraction.md` to capture key findings, patterns, tool insights, and lessons. Then create or update files in `knowledge/principles`, `knowledge/tools`, `knowledge/patterns`, and `knowledge/lessons` from that extraction.
- **Before a project:** When starting a new task, check these folders for relevant principles, tools, patterns, and lessons. Feed them into the Manager brief or Researcher handoff so the team reuses past learning.
- **Domain-specific work:** For **CAT / devcat** related tasks, read `knowledge/principles/cat-vision.md`, `cat-design-principles.md`, `sli-slo-philosophy.md`, and **`devcat-program-brief.md`** first and align proposals with them. The program brief describes current direction, devcat reality, open problems, and the research→experiment→devcat improvement model. See `knowledge/principles/README.md` for the list.

---

## File naming

- Use short, descriptive names in kebab-case: e.g. `api-design.md`, `kubernetes-testing.md`.
- One topic or theme per file. Link between files in `knowledge/` when relevant.
