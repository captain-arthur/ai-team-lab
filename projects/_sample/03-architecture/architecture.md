# Architecture: Local CLI task tracker (sample)

**Role:** Architect  
**Project:** _sample

---

## Overview

- **Purpose:** Allow the user to manage personal tasks from the terminal with plain-text storage and minimal setup.
- **Scope (in):** One recommended tool, install/usage docs, optional wrapper.
- **Scope (out):** No custom app, no server, no GUI.

---

## Components

| Component        | Responsibility              | Boundaries        |
|-----------------|-----------------------------|--------------------|
| CLI tool        | Add/list/complete tasks     | Chosen by Research |
| Data file(s)    | Plain-text task storage     | User-owned, sync’able |
| Optional script| Wrapper or aliases          | Thin layer only   |

---

## Data / process flow

- User runs CLI commands → tool reads/writes local file(s).
- User syncs data directory (or single file) via Dropbox; no app-level sync logic.

---

## Design decisions

| ID | Decision           | Options considered | Chosen    | Rationale                    |
|----|--------------------|--------------------|-----------|-----------------------------|
| 1  | Storage format     | Taskwarrior vs todo.txt | todo.txt | One file, human-readable, sync-friendly |
| 2  | Wrapper script     | Yes vs no         | Optional  | Engineer can add minimal alias or script |
| 3  | Data location      | Default from tool | Use tool default | Simplicity; user can move if needed |

---

## Trade-offs and risks

- **Trade-offs:** Simplicity (todo.txt) vs power (Taskwarrior); we chose simplicity per request.
- **Risks:** None significant. If user outgrows todo.txt, migration is feasible (plain text).
