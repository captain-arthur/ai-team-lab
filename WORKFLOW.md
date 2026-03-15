# AI Team Workflow

This document describes how a task moves through the AI team: from intake to final documentation.

---

## Flow Overview

```
Task Intake → Manager → Research → Architecture → Engineering → Review → Documentation → Knowledge Extraction
```

Each phase has a **prompt** (in `prompts/`), optional **templates** (in `templates/`), and a **project folder** (e.g. `projects/<name>/01-manager/`). Tasks enter the workflow as structured intake documents in `tasks/`. The final stage, **Knowledge Extraction**, writes reusable insights into `knowledge/` so the team learns from completed work.

---

## Stage 0: Task Intake

**Goal:** Define the task in a standard format before the Manager phase.

- **Input:** Raw idea, request, or problem statement from the user.
- **Template:** `tasks/intake-template.md`
- **Output:** A task file in `tasks/<task-name>.md` (e.g. `tasks/example-kubernetes-cat.md`).

**Sections to complete:** Task Title, Problem Description, Goal, Scope, Expected Deliverables, Constraints, Priority, Additional Context. See **README.md** for how to create new tasks.

**Handoff:** The intake document is the primary input for the Manager. Create a project folder under `projects/<project-name>/` and start the Manager phase with this intake.

---

## Phase 1: Manager

**Goal:** Understand the request, decompose it, and define handoffs.

- **Input:** Task intake document from `tasks/` (or, for ad-hoc work, a user request or problem statement).
- **Prompt:** `prompts/manager.md`
- **Output:** Stored in `projects/<project>/01-manager/`
  - Project brief (goal, scope, constraints)
  - Work breakdown (sub-tasks, dependencies)
  - Handoff notes for Researcher, Architect, etc.

**Handoff:** Manager output is the single source of truth for what the project is and what each role should focus on.

---

## Phase 2: Researcher

**Goal:** Gather evidence, compare options, and summarize findings.

- **Input:** Manager brief + research questions from Manager.
- **Prompt:** `prompts/researcher.md`
- **Template:** `templates/research.md`
- **Output:** `projects/<project>/02-research/`
  - Research notes, sources, comparison tables
  - Recommendations and open questions

**Handoff:** Researcher output informs Architect (options, constraints) and Engineer (tools, APIs, limits).

---

## Phase 3: Architect

**Goal:** Propose system/solution design and key decisions.

- **Input:** Manager brief + Research findings.
- **Prompt:** `prompts/architect.md`
- **Template:** `templates/architecture.md`
- **Output:** `projects/<project>/03-architecture/`
  - Architecture document (components, boundaries, data flow)
  - Decisions, trade-offs, risks
  - Optional diagrams (described or linked)

**Handoff:** Architecture is the reference for implementation and review.

---

## Phase 4: Engineer

**Goal:** Produce implementation artifacts (plans, code, configs).

- **Input:** Manager brief + Research + Architecture.
- **Prompt:** `prompts/engineer.md`
- **Output:** `projects/<project>/04-engineering/`
  - Implementation plan or task list
  - Code, configs, scripts, runbooks
  - Notes on gaps or follow-ups

**Handoff:** Concrete artifacts for Reviewer to validate and Writer to document.

---

## Phase 5: Reviewer

**Goal:** Validate correctness, completeness, and consistency; identify gaps.

- **Input:** All prior outputs (Manager through Engineer).
- **Prompt:** `prompts/reviewer.md`
- **Output:** `projects/<project>/05-review/`
  - Review summary (pass/fail, confidence)
  - Checklist results, issues, suggestions
  - List of fixes or follow-up tasks

**Handoff:** Review shapes final documentation and any rework.

---

## Phase 6: Writer

**Goal:** Produce final documentation and reports.

- **Input:** All prior outputs including Review.
- **Prompt:** `prompts/writer.md`
- **Template:** `templates/final-report.md`
- **Output:** `projects/<project>/06-documentation/`
  - Final report or summary
  - User-facing docs, README, runbooks
  - References to knowledge base

**Handoff:** Deliverable is ready for the user. Proceed to Knowledge Extraction to capture reusable learning.

---

## Phase 7: Knowledge Extraction

**Goal:** Extract reusable knowledge from the completed project and store it in the knowledge memory system so future projects can benefit.

- **Input:** All prior outputs (Manager through Writer), especially final report and documentation.
- **Template:** `templates/knowledge-extraction.md`
- **Output:** 
  - A filled extraction document in `projects/<project>/07-knowledge-extraction/` (e.g. `knowledge-extraction.md`).
  - New or updated files in `knowledge/principles/`, `knowledge/tools/`, `knowledge/patterns/`, and `knowledge/lessons/` based on that extraction.

**What to extract:** Key findings, reusable patterns, tool insights, lessons learned, and references. Use the template sections; then create or append to files in the appropriate `knowledge/` subdirectory. See **README.md** and `knowledge/README.md` for what belongs in each.

**Handoff:** The knowledge base is updated. Future tasks can consult `knowledge/` when starting (Manager, Researcher) so the team reuses past learning.

---

## How to Run the Workflow

1. **Create a task (intake):** Copy `tasks/intake-template.md` to `tasks/<task-name>.md`, fill in all sections (Title, Problem, Goal, Scope, Deliverables, Constraints, Priority, Context). See **README.md** and `tasks/README.md`.
2. **Create project:** `projects/<project-name>/` with a `README.md` (goal, context); the project name can match or derive from the task name.
3. **Run phases in order:** Start with Manager using the task intake as input. Then Researcher, Architect, Engineer, Reviewer, Writer, and finally Knowledge Extraction. For each phase, use the corresponding prompt and write outputs into the phase folder.
4. **Use templates:** Copy templates from `templates/` into the phase folder when needed.
5. **Iterate if needed:** Reviewer may trigger rework in Engineer or Architect; re-run from that phase and then Reviewer/Writer again.
6. **Knowledge Extraction:** After Documentation, fill in `templates/knowledge-extraction.md` (e.g. in `07-knowledge-extraction/`), then create or update files in `knowledge/principles/`, `knowledge/tools/`, `knowledge/patterns/`, and `knowledge/lessons/` so the team can reuse what was learned.

---

## Optional Shortcuts

- **Small tasks:** Intake → Manager → Engineer → Writer (skip Research/Architecture/Review if scope is tiny). Still run Knowledge Extraction if something reusable was learned.
- **Research-only:** Intake → Manager → Researcher → Writer (no implementation).
- **Design-only:** Intake → Manager → Researcher → Architect → Writer (no implementation or review).

Adapt the pipeline to the task; the folder structure and prompts stay the same. You can still run Manager with a raw user request instead of a formal intake file when doing ad-hoc work. Whenever a project completes, run Knowledge Extraction and update `knowledge/` so the team keeps learning.
