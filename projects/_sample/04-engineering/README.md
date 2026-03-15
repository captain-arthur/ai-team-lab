# Engineering: Local CLI task tracker (sample)

**Role:** Engineer  
**Project:** _sample

---

## Implementation plan

1. Pick one todo.txt CLI (e.g. **todo.txt-cli** or **topydo**) and document it.
2. Document install (e.g. `brew install todo-txt` or pip).
3. Document basic commands: add, list, complete, file location.
4. Optional: add a tiny wrapper script for "today" or "list by project."

---

## Chosen tool (example)

**todo.txt-cli** (or topydo) — use official install instructions for the user's OS.

---

## How to install (example)

```bash
# macOS with Homebrew
brew install todo-txt-cli
```

(Actual command may vary; Researcher/Engineer would verify.)

---

## Basic usage

```bash
# Add a task
t add "Review AI team workflow doc"

# List tasks
t ls

# Complete task #1
t do 1
```

Data file: `~/todo.txt` (or as configured). User can put this file in a Dropbox folder and point the tool at it via config if needed.

---

## Optional wrapper

A one-line alias or a small script could wrap `t ls` with a filter (e.g. "today"). Left as optional for the sample.

---

## Notes

- No code in this sample repo; real project would link or embed a minimal script.
- Engineer would verify install commands and paths for the user's OS.
