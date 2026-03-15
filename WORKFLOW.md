# AI Team Workflow

This document describes how a task moves through the AI team: from intake to final documentation.

---

## Flow Overview

```
Task Intake → Manager → Research → Architecture → Engineering → Experiment → Review → Documentation → Knowledge Extraction
```

Each phase has a **prompt** (in `prompts/`), optional **templates** (in `templates/`), and a **project folder** (e.g. `projects/<name>/01-manager/`). Tasks enter the workflow as structured intake documents in `tasks/`. For **CAT/devcat** work, the **Experiment** phase runs devcat experiments and produces interpretation notes before Review. The final stage, **Knowledge Extraction**, writes reusable insights into `knowledge/` so the team learns from completed work.

---

## Stage 0: Task Intake

**Goal:** Define the task in a standard format before the Manager phase.

- **Input:** Raw idea, request, or problem statement from the user.
- **Template:** `tasks/intake-template.md`
- **Output:** A task file in `tasks/<task-name>.md` (e.g. `tasks/example-kubernetes-cat.md`).

**Sections to complete:** Task Title, Problem Description, Goal, Scope, Expected Deliverables, Constraints, Priority, Additional Context. See **README.md** for how to create new tasks.

**Handoff:** The intake document is the primary input for the Manager. Create a project folder under `projects/<project-name>/` and start the Manager phase with this intake.

---

## Project–Program relationship

Every project should **declare which program it belongs to** (e.g. CAT, SRE Monitoring). This links the project to a long-term work stream and makes it easier to find related work.

- **How to declare:** In the project’s `README.md` (or in the Manager’s project brief), state the program name, e.g. “This project belongs to the **CAT** program” or “Program: SRE Monitoring.”
- **Where programs live:** Program descriptions and scope live under `programs/<program-name>/README.md`. Projects themselves stay under `projects/<project-name>/`; they are not moved into `programs/`.
- **Why it matters:** Programs support multiple independent work streams (CAT, SRE Monitoring, and future ones such as chaos engineering or capacity planning). Declaring the program keeps the repo organized and clarifies which domain a project contributes to.

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

**Handoff:** Implementation artifacts and runbooks for Experimenter to execute (for CAT/devcat) or for Reviewer to validate.

---

## Phase 5: Experiment (CAT/devcat)

**Goal:** Execute devcat experiments, collect real ClusterLoader2 results, analyze metrics, evaluate SLO candidates, and produce interpretation notes.

- **Input:** Manager brief + Research + Architecture + Engineering (runbook, SLI/SLO definitions).
- **Prompt:** `prompts/experimenter.md`
- **Safety:** Experiment phase follows **Experiment Safety Rules** (`.cursor/rules/devcat-experiment-safety.mdc`): **preflight environment checks** (kubeconfig, API server, StorageClass), **safe experiment execution** (safe-default ClusterLoader2 options, timeout protection), and **result inspection**.
- **Output:** `projects/<project>/05-experiment/`
  - Experiment notes (what was run, results directory layout, metrics found)
  - SLI measurement extraction (actual metrics produced by ClusterLoader2)
  - Mapping of results to SLI candidates from architecture
  - SLO evaluation (if applicable) and interpretation notes

**Fast feedback (3-minute start report and progress updates):**  
Devcat experiments often run 10–30 minutes. To avoid long silent periods, the Experimenter **must**:

1. **Start report within 3 minutes** of experiment start. Include: experiment started confirmation, **run id**, command being executed, scenario used, override used, **expected duration**, **current phase**.
2. **Progress updates every ~5 minutes** during long runs. Each update: **current phase** (create / load / gather), whether Prometheus targets are up (if applicable), whether obvious failures appeared in logs.
3. **Full experiment notes** (results layout, metrics, SLI mapping, interpretation) are written **only after** the experiment completes.

This rule applies especially to devcat experiments; the goal is fast feedback instead of waiting silently for the full run to finish.

**When to run:** Only for projects that involve CAT or devcat (ClusterLoader2-based acceptance testing). For other projects, skip this phase and proceed from Engineering to Review.

**Handoff:** Experiment findings inform Reviewer and Writer; interpretation notes feed SLO refinement and devcat improvement.

---

## Phase 6: Reviewer

**Goal:** Validate correctness, completeness, and consistency; identify gaps.

- **Input:** All prior outputs (Manager through Engineer, and Experiment when applicable).
- **Prompt:** `prompts/reviewer.md`
- **Output:** `projects/<project>/06-review/`
  - Review summary (pass/fail, confidence)
  - Checklist results, issues, suggestions
  - List of fixes or follow-up tasks

**Handoff:** Review shapes final documentation and any rework.

---

## Phase 7: Writer

**Goal:** Produce final documentation and reports.

- **Input:** All prior outputs including Review (and Experiment when applicable).
- **Prompt:** `prompts/writer.md`
- **Template:** `templates/final-report.md`
- **Output:** `projects/<project>/07-documentation/`
  - Final report or summary
  - User-facing docs, README, runbooks
  - References to knowledge base

**Handoff:** Deliverable is ready for the user. Proceed to Knowledge Extraction to capture reusable learning.

---

## Phase 8: Knowledge Extraction

**Goal:** Extract reusable knowledge from the completed project and store it in the knowledge memory system so future projects can benefit.

- **Input:** All prior outputs (Manager through Writer), especially final report and documentation.
- **Template:** `templates/knowledge-extraction.md`
- **Output:** 
  - A filled extraction document in `projects/<project>/08-knowledge-extraction/` (e.g. `knowledge-extraction.md`).
  - New or updated files in `knowledge/principles/`, `knowledge/tools/`, `knowledge/patterns/`, and `knowledge/lessons/` based on that extraction.

**What to extract:** Key findings, reusable patterns, tool insights, lessons learned, and references. Use the template sections; then create or append to files in the appropriate `knowledge/` subdirectory. See **README.md** and `knowledge/README.md` for what belongs in each.

**Handoff:** The knowledge base is updated. Future tasks can consult `knowledge/` when starting (Manager, Researcher) so the team reuses past learning.

---

## How to Run the Workflow

1. **Create a task (intake):** Copy `tasks/intake-template.md` to `tasks/<task-name>.md`, fill in all sections (Title, Problem, Goal, Scope, Deliverables, Constraints, Priority, Context). See **README.md** and `tasks/README.md`.
2. **Create project:** `projects/<project-name>/` with a `README.md` (goal, context); the project name can match or derive from the task name.
3. **Run phases in order:** Start with Manager using the task intake as input. Then Researcher, Architect, Engineer, and (for CAT/devcat work) Experiment. Then Reviewer, Writer, and finally Knowledge Extraction. For each phase, use the corresponding prompt and write outputs into the phase folder.
4. **Use templates:** Copy templates from `templates/` into the phase folder when needed.
5. **Iterate if needed:** Reviewer may trigger rework in Engineer or Architect; re-run from that phase and then Reviewer/Writer again.
6. **Knowledge Extraction:** After Documentation, fill in `templates/knowledge-extraction.md` (e.g. in `08-knowledge-extraction/`), then create or update files in `knowledge/principles/`, `knowledge/tools/`, `knowledge/patterns/`, and `knowledge/lessons/` so the team can reuse what was learned.

---

## Optional Shortcuts

- **Small tasks:** Intake → Manager → Engineer → Writer (skip Research/Architecture/Review if scope is tiny). Still run Knowledge Extraction if something reusable was learned.
- **Research-only:** Intake → Manager → Researcher → Writer (no implementation).
- **Design-only:** Intake → Manager → Researcher → Architect → Writer (no implementation or review).

Adapt the pipeline to the task; the folder structure and prompts stay the same. You can still run Manager with a raw user request instead of a formal intake file when doing ad-hoc work. Whenever a project completes, run Knowledge Extraction and update `knowledge/` so the team keeps learning.
