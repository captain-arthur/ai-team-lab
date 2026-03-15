# Principles

Store **design principles, decision rules, and standards** that the AI team (and you) want to apply across projects.

---

## Domain principles: CAT / devcat (Cluster Acceptance Test)

**CAT·devcat 관련 작업을 할 때는 아래 문서를 먼저 읽고 제안·설계를 이 방향에 맞춘다.**

| 문서 | 내용 |
|------|------|
| **cat-vision.md** | CAT의 궁극 목표(실용적 CAT 시스템 완성), Conformance/기능성 + Performance/부하 관점, 수용/거부 판단 산출. |
| **cat-design-principles.md** | 도구 조합·의존 최소화·as-is 사용, 테스트=Job(scenario injector + SLI/SLO measurement + assertion), 결과 디렉터리 규약, 단순 오케스트레이션, ClusterLoader2 기반 출발점. |
| **sli-slo-philosophy.md** | Kubernetes scalability 철학, you promise we promise, 수용 테스트용 SLO 정당화, 소규모 vs 대규모 기준 차이. |
| **devcat-program-brief.md** | 실제 CAT 프로그램 방향·책임, devcat 저장소의 현재 상태(ClusterLoader2, perfdash, config/ol-test, results/), 열린 문제, 시각화·작업 모델(research→experiment→interpretation→devcat improvement). **CAT/devcat 작업의 단일 소스 오브 트루스.** |

- Manager brief, Researcher handoff, Architect 설계 시 CAT·devcat이 걸리면 위 원칙과 devcat-program-brief를 입력으로 사용한다. 제안은 devcat 현실을 존중하고 점진적 개선으로 연결한다.

---

## What to put here

- **Design principles** — e.g. "prefer simple over clever", "document decisions and rationale", "fail fast and visibly".
- **Decision rules** — e.g. "when choosing between build vs buy, prefer buy if an OSS option exists and is maintained".
- **Standards** — coding style, naming, or process rules that should be respected in future work.

---

## Format

- One file per theme or domain (e.g. `security.md`, `api-design.md`, `simplicity.md`).
- Use clear headings and short bullets so roles can scan quickly when starting a new task.
- Optionally add a **Source** line linking to the project that produced the principle.

---

## Example

```markdown
# API design principles

- Prefer small, focused endpoints over large multipurpose ones.
- Document errors and status codes in one place.
- Version in the URL path (e.g. /v1/...) from day one.
```
