# Project brief: Kubernetes Cluster Acceptance Test (CAT) framework

**Project:** example-kubernetes-cat  
**Source:** `task-intake.md`  
**Role:** Manager

---

## 1. Project brief

### Goal

Deliver a clear, implementable design for a Kubernetes Cluster Acceptance Test framework—how tests are defined, run, what they validate, and how results are reported—so a team can adopt it (via existing tools or a small custom layer).

### Scope

| In scope | Out of scope |
|----------|--------------|
| Framework design (concepts, test categories, execution model, reporting) | Full implementation in code |
| Comparison of existing tools (e.g. sonobuoy, custom scripts, test harnesses) | Managing or deploying clusters |
| Documentation and a concrete recommendation | CI/CD integration details (note as follow-up only) |
| Optional: minimal example test definition or runbook | Proprietary or non-standard Kubernetes extensions |

### Constraints

- **Technology:** Standard Kubernetes (1.24+); no proprietary extensions.
- **Approach:** Prefer reusing existing OSS tools over building from scratch.
- **Audience:** Design and docs must be understandable by platform/SRE engineers.

### Success criteria

- A design document exists that defines test types, flow, and tooling options.
- A clear recommendation: use an existing tool vs. a small custom layer, with rationale.
- A short report (and optionally a minimal example or runbook) that a platform/SRE engineer can use to implement or adopt the framework.

---

## 2. Work breakdown

| # | Sub-task | Owner | Depends on |
|---|----------|--------|------------|
| 1 | Research existing CAT/conformance tools (sonobuoy, others) and compare fit for our use case | Researcher | — |
| 2 | Recommend approach: existing tool vs small custom layer; document pros/cons | Researcher | 1 |
| 3 | Design CAT framework: test categories, execution model, reporting, data flow | Architect | 2 |
| 4 | Produce implementation plan, example test definition or runbook, and any minimal config/spec | Engineer | 3 |
| 5 | Review design and docs for correctness, completeness, and usability | Reviewer | 4 |
| 6 | Final report and user-facing documentation | Writer | 5 |
| 7 | Extract reusable knowledge (patterns, tools, lessons) into `knowledge/` | Knowledge Extraction | 6 |

---

## 3. Handoff notes

### For Researcher

- **Questions to answer:**
  - What existing tools or projects support cluster acceptance/conformance testing (e.g. sonobuoy, custom harnesses, OSS projects)?
  - How do they compare for: ease of use, test definition format, reporting, and fit for “validate before production” (not only conformance)?
  - What are the trade-offs between using an existing tool vs. a thin custom layer (e.g. script + YAML)?
- **References to check:** Official Kubernetes testing/conformance docs, sonobuoy docs/repo, any well-known OSS CAT or conformance tools.
- **Output:** Research notes and comparison in `02-research/`; clear recommendation with rationale.

### For Architect

- **Design boundaries:** Framework design only (concepts, test categories, execution, reporting). No implementation code; no cluster lifecycle or CI/CD design. Assume standard Kubernetes 1.24+.
- **Key decisions to make:** Test categories (e.g. API, nodes, DNS, storage, optional security); execution model (who runs what, where); reporting format and where results live.
- **Non-functional:** Design must be implementable by a small team; prefer alignment with Researcher’s recommended tool(s).
- **Output:** Architecture/design document in `03-architecture/`.

### For Engineer

- **Deliverables:** Implementation plan (how to implement or adopt the framework), optional minimal example test definition (e.g. YAML or script), and a short runbook. No full codebase—only what’s needed to make the design actionable.
- **Tech stack hints:** Use standard Kubernetes tooling (kubectl, maybe sonobuoy or similar if recommended). Plain config/YAML and scripts; language agnostic unless one is clearly better.
- **Output:** Artifacts and notes in `04-engineering/`.

### For Reviewer

- **Validate:** Does the design meet the goal and scope? Are research recommendation and architecture aligned? Are deliverables (design doc, recommendation, example/runbook) complete and usable by platform/SRE engineers? Any gaps or risks?
- **Output:** Review summary and checklist in `05-review/`.

### For Writer

- **Audience:** Platform and SRE engineers who provision or validate clusters.
- **Format:** Short final report (summary, recommendation, next steps) and user-facing documentation (how to adopt/run the framework). Markdown; optional Mermaid or ASCII for flow.
- **Must document:** Final design summary, tool/recommendation, how to run or adopt the framework, and any limitations or follow-up (e.g. CI/CD).
- **Output:** Final report and docs in `06-documentation/`.
