# Cursor rules for the AI team workspace

These rules guide Cursor to behave as the default operator of this repository: following the workflow, respecting role boundaries, and writing outputs in the right places.

| Rule file | Purpose |
|-----------|---------|
| `repository-behavior.mdc` | Start from task intake, prefer clarity and reuse, stay within repo structure. |
| `workflow-execution.mdc` | Run phases in order (Intake → Manager → … → Knowledge Extraction); don’t skip unless asked. |
| `role-boundaries.mdc` | Keep each role within its responsibility; use prompts in `prompts/`. |
| `output-locations.mdc` | Save outputs under `projects/<project>/0X-<phase>/`; use templates. |
| `devcat-execution.mdc` | For CAT/devcat work: allow running local commands to execute devcat experiments and inspect results. |
| `devcat-experiment-safety.mdc` | For devcat experiments: preflight checks, safe-default ClusterLoader2 options, timeout protection. |
| `knowledge-extraction-behavior.mdc` | After completion, run extraction and update `knowledge/`. |

All rules are set to `alwaysApply: true` so they apply in every conversation in this repo. Edit the `.mdc` files to change behavior.
