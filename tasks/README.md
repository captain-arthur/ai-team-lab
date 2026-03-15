# Task Intake

Tasks enter the AI team workflow from this directory. Each task is a **structured intake document** that the Manager uses as the primary input.

---

## How to create a new task

1. **Copy the template:** Copy `intake-template.md` into a new file in `tasks/`.
2. **Name the file:** Use a short, descriptive name in kebab-case, e.g. `api-migration-plan.md`, `auth-design-review.md`.
3. **Fill in all sections:** Task Title, Problem Description, Goal, Scope, Expected Deliverables, Constraints, Priority, Additional Context.
4. **Start the workflow:** Create a matching project under `projects/<project-name>/` and run the Manager phase with this intake as input. See **WORKFLOW.md** and **README.md** for the full flow.

---

## File naming

- Use kebab-case: `my-task-name.md`.
- Keep names short and recognizable so they map clearly to `projects/<project-name>/`.

---

## Example

See **example-kubernetes-cat.md** for a filled-in example task.
