# Personal AI Team Workspace

A structured repository for running an **AI team workflow**: multiple AI roles collaborate on research, architecture, implementation planning, and documentation. The goal is not only software development but **structured problem solving** across research, design, analysis, and writing.

---

## Overview

This workspace supports:

- **Research** — Tool research, comparison, and evidence gathering  
- **Architecture design** — System and solution design  
- **Technical analysis** — Feasibility, trade-offs, risks  
- **Implementation planning** — Tasks, milestones, artifacts  
- **Documentation** — Reports, specs, and knowledge capture  

Different **AI roles** work in sequence (or in parallel where it makes sense), each with a clear prompt and output location.

---

## Repository Structure

```
ai-team-lab/
├── README.md                 # This file
├── WORKFLOW.md               # How tasks move through the AI team
├── .cursor/rules/            # Cursor operating rules — AI team behavior (see below)
├── tasks/                    # Task intake — define work before the workflow
│   ├── README.md
│   ├── intake-template.md    # Standard format for new tasks
│   └── example-kubernetes-cat.md
├── scripts/                  # Workflow runner and helpers
│   ├── README.md
│   └── run_workflow.py       # Initialize a project from a task file
├── prompts/                  # Role prompts for each AI role
│   ├── manager.md
│   ├── researcher.md
│   ├── architect.md
│   ├── engineer.md
│   ├── reviewer.md
│   └── writer.md
├── templates/                # Reusable document templates
│   ├── research.md
│   ├── architecture.md
│   └── final-report.md
├── knowledge/                # Knowledge memory — reusable learning from projects
│   ├── README.md
│   ├── principles/           # Design principles, decision rules
│   ├── tools/                # Tool evaluations and usage notes
│   ├── patterns/             # Reusable solution patterns
│   └── lessons/              # Lessons learned
├── projects/                 # One folder per project; outputs by phase
│   ├── README.md
│   └── _sample/              # Example project
│       ├── README.md
│       ├── 01-manager/
│       ├── 02-research/
│       ├── 03-architecture/
│       ├── 04-engineering/
│       ├── 05-review/
│       └── 06-documentation/
```

---

## Roles

| Role        | Focus                                      | Typical output                    |
|------------|---------------------------------------------|-----------------------------------|
| **Manager** | Task analysis, decomposition, prioritization | Brief, work breakdown, handoffs   |
| **Researcher** | Tool research, comparison, evidence       | Research notes, comparison tables |
| **Architect**  | System/solution design, constraints       | Architecture doc, diagrams        |
| **Engineer**   | Implementation plans, code, configs        | Specs, code, runbooks             |
| **Reviewer**   | Validation, critique, gaps                | Review notes, checklist           |
| **Writer**     | Documentation, reports, summaries         | Final report, docs, README        |

---

## Running a new AI team project

1. **Create a task** in `tasks/` using `tasks/intake-template.md` (e.g. `tasks/my-task.md`).
2. **Initialize the project** with the workflow runner. It will create the project folder, phase directories, and copy templates:

   ```bash
   python scripts/run_workflow.py tasks/my-task.md
   ```

   This creates `projects/my-task/` with `01-manager` through `06-documentation`, copies the task file as `task-intake.md`, and drops the research, architecture, and final-report templates into the right phases.
3. **Run each phase** using the prompts in `prompts/` and the templates in the phase folders. Start with the Manager (input: `task-intake.md`), then Researcher, Architect, Engineer, Reviewer, Writer, and finally **Knowledge Extraction**. See **WORKFLOW.md** for the full process.
4. **After a project is completed,** run **Knowledge Extraction**: fill in `templates/knowledge-extraction.md` (e.g. in `07-knowledge-extraction/`), then create or update files under `knowledge/principles/`, `knowledge/tools/`, `knowledge/patterns/`, and `knowledge/lessons/`. This turns the project into reusable learning so the AI team gets better over time. See `knowledge/README.md` for what goes in each directory.

The script is only a **workflow initializer** — it does not run AI or automate phases; it just prepares the workspace so you can execute the workflow step-by-step.

---

## Quick Start

1. **Create a task (intake):** Copy `tasks/intake-template.md` to a new file in `tasks/` (e.g. `tasks/my-task-name.md`). Fill in **Task Title**, **Problem Description**, **Goal**, **Scope**, **Expected Deliverables**, **Constraints**, **Priority**, and **Additional Context**. This is the standard way to define incoming work before it enters the workflow. See `tasks/README.md` and the example `tasks/example-kubernetes-cat.md`.
2. **Create a project**: Run `python scripts/run_workflow.py tasks/<task-name>.md` to scaffold the project and phase folders (see **Running a new AI team project** below). Or create `projects/<project-name>/` and phase folders manually.
3. **Run the workflow**: Follow `WORKFLOW.md` — give the Manager the task intake as input, then run Researcher, Architect, Engineer, Reviewer, and Writer, using the corresponding folders and templates.
4. **Extract knowledge after each project:** Use the Knowledge Extraction phase and `templates/knowledge-extraction.md`; store results in `knowledge/principles/`, `knowledge/tools/`, `knowledge/patterns/`, and `knowledge/lessons/`. See **WORKFLOW.md** and `knowledge/README.md`.
5. **Use templates**: Copy from `templates/` into the project phase folder (e.g. `02-research/`) when needed.

---

## Cursor and the AI team workflow

When you work in this repository with **Cursor**, the assistant is expected to follow the **operating rules** in `.cursor/rules/`. Those rules make Cursor behave as the default operator of this AI team workspace:

- Start from a task intake in `tasks/` when available
- Follow the workflow order in **WORKFLOW.md** and not skip phases unless you ask
- Keep each role within its responsibility (Manager, Researcher, Architect, etc.)
- Save outputs under `projects/<project-name>/0X-<phase>/` and use the templates
- After a project is completed, consider Knowledge Extraction and updates to `knowledge/`
- Prefer clarity, structure, and reuse over complexity

You can review or edit the rules in `.cursor/rules/` to adjust how Cursor behaves here.

---

## Design Principles

- **Simple** — Minimal structure; easy to navigate and extend.  
- **Readable** — Clear names, short docs, Markdown-only.  
- **Extensible** — Add roles, templates, or phases without breaking the workflow.

See **WORKFLOW.md** for the full task flow and **prompts/** for each role’s instructions.
