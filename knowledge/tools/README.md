# Tools

Store **tool evaluations, recommendations, and usage notes** — libraries, CLIs, platforms, and services the team has researched or used.

---

## What to put here

- **Tool comparisons** — short summaries of options for a given job (e.g. "CLI task trackers: todo.txt vs Taskwarrior").
- **Recommendations** — when to use which tool, pros/cons, and caveats.
- **Usage notes** — gotchas, config tips, version constraints, or links to official docs.

---

## Format

- One file per tool, family, or comparison (e.g. `todo-txt-cli.md`, `kubernetes-testing-tools.md`).
- Include: name, purpose, when to use it, key limitations, and a link to official docs or repo.
- Optionally add **Source** (project that produced this insight).

---

## Example

```markdown
# todo.txt CLI

- **Purpose:** Manage tasks from the terminal with a single plain-text file.
- **When to use:** Local, single-user task lists; sync via Dropbox or similar.
- **Limitations:** Fewer features than Taskwarrior; pick a concrete implementation (e.g. topydo).
- **Docs:** https://github.com/todotxt/todo.txt
```
