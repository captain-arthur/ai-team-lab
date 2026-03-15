# Projects

Each project is a folder under `projects/` with a consistent structure. Outputs from each AI team phase are stored in numbered subfolders.

---

## Naming

- Use a short, snake_case or kebab-case name: e.g. `api-migration`, `auth-design`, `q1-report`.
- Reserve `_sample` as the example project; do not use it for real work.

---

## Standard structure

```
projects/<project-name>/
├── README.md           # Project goal, context, and (optionally) original request
├── 01-manager/         # Brief, work breakdown, handoffs
├── 02-research/        # Research notes, comparisons
├── 03-architecture/    # Architecture doc, decisions
├── 04-engineering/     # Implementation plan, code, configs
├── 05-review/          # Review summary, checklist, issues
└── 06-documentation/   # Final report, user docs
```

You can skip a phase (e.g. no formal research or no implementation); leave that folder empty or add a one-line note like "Skipped — not needed for this task."

---

## Starting a new project

1. Create `projects/<project-name>/`.
2. Add `projects/<project-name>/README.md` with the goal and any context.
3. Run the workflow from **Manager** (`prompts/manager.md`), then fill each phase folder as you go.
4. Copy templates from `templates/` into the relevant phase folder when needed.

See **WORKFLOW.md** at the repo root for the full flow. See **\_sample** for an example.
