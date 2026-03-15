# Project brief: Local CLI Task Tracker

**Role:** Manager  
**Project:** _sample

---

## Goal

Deliver a clear recommendation and usage documentation for a **local, CLI-only task tracker** that uses plain-text storage and requires minimal setup, so the user can adopt it without a server or GUI.

---

## Scope

| In scope | Out of scope |
|----------|--------------|
| Research and compare 2–3 CLI task tools | Building a custom app from scratch |
| Recommend one option with rationale | Web UI, mobile app, or server |
| Document how to install and use it | Cloud sync implementation (user will use Dropbox) |
| Optional: minimal wrapper script or config | Full automation or integrations |

---

## Constraints

- CLI only; no GUI or web UI.
- Plain-text or human-readable storage (e.g. markdown, JSON in a file).
- Minimal setup (ideally one binary or one script).
- Must run on the user's OS (assume macOS/Linux for the sample).

---

## Success criteria

- One recommended tool with a short comparison and rationale.
- A short "how to install and use" doc the user can follow.
- Final report summarizing the choice and next steps.

---

## Work breakdown

| # | Task | Owner | Depends on |
|---|------|--------|------------|
| 1 | Research 2–3 CLI task tools (plain-text, minimal) | Researcher | — |
| 2 | Recommend one and document rationale | Researcher | 1 |
| 3 | Describe "architecture" (single user, local, data format) | Architect | 2 |
| 4 | Write install/run steps; optional tiny script | Engineer | 3 |
| 5 | Review completeness and clarity | Reviewer | 4 |
| 6 | Final report and user README | Writer | 5 |

---

## Handoffs

### For Researcher

- **Questions to answer:** Which CLI task trackers support plain-text storage and minimal setup? Compare 2–3 (e.g. taskwarrior, todo.txt, or similar). What are pros/cons for a single-user, local-only workflow?
- **References:** Official docs or GitHub repos for each tool.

### For Architect

- **Boundaries:** Single user, local only; no server. Data format and file location matter for Dropbox sync.
- **Decisions:** Where to store data; whether to recommend a wrapper script or use the tool as-is.

### For Engineer

- **Deliverables:** Install steps, basic usage commands, optional wrapper script or config snippet. No full app.
- **Tech:** Whatever the recommended tool uses (e.g. shell, one binary).

### For Reviewer

- **Validate:** Does the recommendation match scope? Are install and usage steps clear and runnable? Any missing caveats?

### For Writer

- **Audience:** End user (technical, comfortable with terminal).
- **Format:** Final report (summary + recommendation) + short user README (how to install and use).
- **Must document:** Chosen tool, rationale, install, basic commands, and data location for sync.
