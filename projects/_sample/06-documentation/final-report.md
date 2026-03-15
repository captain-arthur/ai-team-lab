# Final report: Local CLI task tracker (sample)

**Role:** Writer  
**Project:** _sample  
**Note:** This is the **sample project** for the AI team workspace; content is illustrative.

---

## Executive summary

We researched CLI task trackers with plain-text storage and minimal setup, recommended the **todo.txt** approach with a compatible CLI (e.g. todo.txt-cli or topydo), and documented install and basic usage. The outcome is a clear path for the user to adopt a single-file, sync-friendly task list from the terminal.

---

## Goal and scope

- **Goal:** Recommend and document a local, CLI-only task tracker with plain-text storage and minimal setup.
- **In scope:** Research, comparison, recommendation, install/usage docs.
- **Out of scope:** Custom app, server, GUI, or cloud sync implementation.
- **Success criteria:** One recommended tool with rationale and a short how-to.

---

## Approach

- **Research:** Compared Taskwarrior and todo.txt-style tools; chose todo.txt for one-file simplicity and sync-friendliness.
- **Design:** Single-user, local-only; data in one file; optional wrapper script.
- **Implementation:** Chose a todo.txt CLI, documented install and basic commands, noted data file location.

---

## Outcomes

- **Deliverables:** This report, research notes, architecture summary, engineering README, and review notes (all in `projects/_sample/`).
- **Status:** Sample complete; real project would verify exact install and paths per OS.
- **Known limitations:** Install command and data path are examples; user should confirm for their OS and chosen CLI.

---

## Key decisions and rationale

- **todo.txt over Taskwarrior:** Better fit for "minimal setup" and "one file" for Dropbox; human-readable format.
- **Use existing CLI:** No custom app; lower maintenance and faster adoption.

---

## Follow-up and recommendations

- **Next steps:** User picks a todo.txt CLI (e.g. todo.txt-cli, topydo), runs install, creates `~/todo.txt` or configures path, and optionally moves file to Dropbox.
- **Open issues:** Verify brew/pip (or other) install and default file path for chosen CLI and OS.
- **Knowledge base:** Optional — add "CLI task tracking: todo.txt vs Taskwarrior" to `knowledge/references.md` for future projects.

---

## References

- Project folder: `projects/_sample/`
- Architecture: `03-architecture/architecture.md`
- Implementation: `04-engineering/README.md`
- Review: `05-review/review-notes.md`
