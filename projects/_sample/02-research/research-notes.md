# Research: CLI task trackers (sample)

**Role:** Researcher  
**Project:** _sample

---

## Research questions

1. Which CLI task trackers support plain-text or human-readable storage and minimal setup?
2. How do 2–3 options compare for a single-user, local-only workflow?
3. What are the trade-offs (complexity, format, sync-friendliness)?

---

## Findings

### Option 1: Taskwarrior

- **Finding:** Mature CLI task manager; data in plain-text under `~/.task`. Flexible, many commands.
- **Sources:** [taskwarrior.org](https://taskwarrior.org), GitHub taskwarrior/taskwarrior
- **Note:** Slightly more setup (optional config); very powerful. Data is text-based but in its own format.

### Option 2: todo.txt

- **Finding:** Format is a single `todo.txt` file; any editor or CLI (e.g. todo.txt-cli) can work with it. Minimal and sync-friendly.
- **Sources:** [todotxt.org](https://github.com/todotxt/todo.txt), various CLI implementations
- **Note:** Very simple; great for Dropbox (one file). Fewer features than Taskwarrior.

### Option 3: Simple custom script

- **Finding:** A small shell script + a markdown or text file is possible but reinvents the wheel; maintenance on the user.
- **Note:** Only recommend if no existing tool fits; for this request, prefer an existing tool.

---

## Comparison

| Criterion       | Taskwarrior   | todo.txt CLI   |
|----------------|---------------|----------------|
| Plain-text     | Yes (own dir) | Yes (one file) |
| Minimal setup  | Medium        | High           |
| Sync-friendly  | Yes (dir)     | Yes (one file) |
| Complexity     | Higher        | Lower          |
| Single binary  | Yes           | Depends on impl|

**Recommendation:** **todo.txt** (or a todo.txt-compatible CLI) for this user — emphasis on "minimal setup" and "plain-text storage" and Dropbox sync (one file is simpler).

**Rationale:** One file is easier to sync and back up; format is human-readable; many CLI implementations exist.

---

## Summary for downstream roles

- **Top findings:** todo.txt format is one file, human-readable, and very sync-friendly; Taskwarrior is more powerful but more setup and multi-file.
- **Open questions:** Which todo.txt CLI to recommend (e.g. topydo, todo.txt-cli) — Engineer can pick one and document.
- **Risks:** None critical; user can switch later since storage is plain text.
