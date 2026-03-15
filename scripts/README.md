# Scripts

Lightweight automation for the AI team workspace.

---

## run_workflow.py

**Workflow initializer.** Given a task file in `tasks/`, it:

1. Creates a new project directory under `projects/` (name derived from the task filename).
2. Creates the six phase directories: `01-manager`, `02-research`, `03-architecture`, `04-engineering`, `05-review`, `06-documentation`.
3. Copies the relevant templates into the phase folders (`research.md`, `architecture.md`, `final-report.md`).
4. Copies the task file into the project as `task-intake.md`.
5. Writes a minimal `README.md` in the project.

It does **not** run AI or execute phases; it only scaffolds the workspace so you can run each phase (e.g. with an AI assistant) using the prompts and templates.

**Usage:**

```bash
python scripts/run_workflow.py tasks/my-task.md
```

Or from the repo root:

```bash
python scripts/run_workflow.py tasks/example-kubernetes-cat.md
```

Requires Python 3.6+ (standard library only).
