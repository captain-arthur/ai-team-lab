# Project brief: Kubernetes Cluster Health Monitoring Model

**Project:** sre-monitoring-cluster-health-model  
**Source:** `task-intake.md`  
**Role:** Manager  
**Program:** SRE Monitoring

---

## 1. Project brief

### Goal

운영자가 **5–10분 안에** “클러스터가 지금 안전한가?”, “곧 불안전해질 징후는?”, “문제 시 어디를 먼저 볼 것인가?”를 답할 수 있도록, **실무에 쓸 수 있는 Kubernetes Cluster Health Monitoring Model**을 설계한다. 핵심 질문: **“운영 건강을 나타낼 수 있는 최소한의 신호 집합은 무엇인가?”**

### Scope

| In scope | Out of scope |
|----------|--------------|
| Cluster Health Monitoring Model 설계(신호 분류, 상위 건강 신호, 트렌드/리스크, drill-down/TOP10) | 수많은 dashboard 제작 또는 대량 metric 나열 |
| 모델의 세 개념 그룹: Cluster Health Summary, Trend/Risk Indicators, Top Offenders/Drill-down | Prometheus/Grafana 도입 또는 인프라 변경 |
| 카테고리 정의: Node/Infrastructure, Control Plane, Workload/Pod, Capacity/Resource pressure, Network/ingress | Alert 정책·runbook 상세 설계(후속 프로젝트) |
| Core signal list: signal name, Prometheus metric/PromQL 예시, “왜 중요한가”, “어떤 운영 질문에 답하는가”, summary/trend/top-offender 구분 | |
| Prometheus에서 실제 조회 가능한 신호만 포함·검증 | |
| 실험: 일부 신호의 조회 가능성, retention·query 복잡도·리소스 제한 문서화 | |
| 최종 문서: **Kubernetes Cluster Health Monitoring Model** (이후 dashboard/alert/runbook 프로젝트 기반) | |

### Constraints

- **최소 집합:** 모든 것을 모니터링하지 않는다. signal-to-noise가 높은 소수 신호에 집중.
- **운영자 UX:** 수십 개 패널을 훑을 필요가 없어야 함.
- **Prometheus 제약:** 장기·복잡 쿼리 제한; retention·리소스 사용 고려.
- **노이즈:** 일부 신호(예: Calico readiness probe 실패)는 알림을 만들 수 있으나 실제 서비스 영향과 직결되지 않을 수 있음 — 모델에서 구분·설명.
- **TOP10 뷰:** 고부하 워크로드 등 drill-down은 TOP10 스타일 검토.
- **트렌드·초기 리스크:** “곧 불안전해질 수 있다”는 징후를 나타내는 신호를 포함.

### Success criteria

- **Cluster Health Monitoring Model** 문서가 존재하며, 세 그룹(Health Summary, Trend/Risk, Top Offenders)과 카테고리(Node, Control Plane, Workload, Capacity, Network)가 정의되어 있다.
- **Core signal list**가 있고, 각 신호에 대해 signal name, Prometheus metric 또는 PromQL 예시, “왜 중요한가”, “어떤 운영 질문에 답하는가”, health summary / trend / top offender 구분이 명시되어 있다.
- 일부 신호에 대해 현재 Prometheus 환경에서의 조회 가능성이 실험되고, 제한 사항(retention, query 복잡도, 리소스)이 문서화되어 있다.
- 리뷰에서 “운영자가 이 모델(또는 이를 반영한 dashboard)을 보고 5–10분 안에 클러스터 건강을 파악할 수 있는가?”에 대한 평가가 문서화되어 있다.

---

## 2. Work breakdown

| # | Sub-task | Owner | Depends on |
|---|----------|--------|------------|
| 1 | Research: 프로덕션 Kubernetes 모니터링, Prometheus+Grafana 운영 dashboard, SRE 모니터링 관행, 클러스터 건강 지표; 신호 **카테고리** 중심 | Researcher | — |
| 2 | Architecture: Cluster Health Monitoring Model 설계 — 세 그룹(Health Summary, Trend/Risk, Top Offenders) 및 카테고리(Node, Control Plane, Workload, Capacity, Network) | Architect | 1 |
| 3 | Engineering: Core signal list 작성 — signal name, Prometheus metric/PromQL, 의미, 운영 질문, summary/trend/top-offender 구분; 소수·고신호 집합 유지 | Engineer | 2 |
| 4 | Experiment: 일부 신호의 현재 Prometheus 조회 가능성 검증; retention·query 복잡도·리소스 제한 문서화 | Experimenter | 3 |
| 5 | Review: “5–10분 안에 클러스터 건강 파악 가능한가?” 평가; 모델·신호 목록 완전성·실용성 검토 | Reviewer | 4 |
| 6 | Documentation: 최종 문서 **Kubernetes Cluster Health Monitoring Model** 작성 — 이후 dashboard/alert/runbook 프로젝트 기반 | Writer | 5 |
| 7 | Knowledge Extraction: 모델·패턴·제한 사항을 `knowledge/`에 추출 | Knowledge Extraction | 6 |

---

## 3. Handoff notes

### For Researcher

- **답할 질문:**
  - 모니터링 신호는 어떻게 **카테고리화**하는가? (node, control plane, workload, capacity, network 등)
  - “클러스터가 건강하다”를 나타내는 신호는 무엇인가?
  - “클러스터가 곧 불건강해질 수 있다”를 나타내는 신호는 무엇인가?
  - 노이즈가 많고 실행 가능하지 않은(noisy, not actionable) 신호는 어떤 것인가?
- **조사할 대상:** Kubernetes 모니터링 문서, Prometheus + Grafana 운영용 dashboard, SRE 모니터링 관행, 프로덕션 환경에서 쓰는 cluster health indicator. **대량 metric 나열보다 신호 카테고리**에 초점.
- **Output:** `02-research/`에 research notes. 신호 카테고리, 건강/리스크/노이즈 구분, 권장 사항 요약.

### For Architect

- **설계 범위:** **Cluster Health Monitoring Model** 한 편. 세 가지 개념 그룹을 반드시 구분한다.
  - **Cluster Health Summary:** 지금 클러스터가 건강한지를 직접 나타내는 신호.
  - **Trend / Risk Indicators:** 클러스터가 곧 불건강해질 수 있음을 나타낼 수 있는 신호.
  - **Top Offenders / Drill-down:** 어떤 워크로드·노드가 부담을 주는지 파악하는 신호(TOP10 스타일).
- **카테고리 정의:** Node/Infrastructure, Control Plane, Workload/Pod, Capacity/Resource pressure, Network/ingress 등.
- **제약 반영:** “모든 것을 모니터링하지 않음”, Prometheus 조회 가능성, 노이즈 신호 구분, TOP10 뷰·트렌드 고려.
- **Output:** `03-architecture/`에 Cluster Health Monitoring Model 설계 문서.

### For Engineer

- **산출물:** **Core signal list** — 첫 번째 실용 dashboard 구축에 쓸 수 있는 수준. 각 신호마다:
  - signal name
  - Prometheus metric 또는 PromQL 예시
  - why it matters(왜 중요한가)
  - what operational question it answers(어떤 운영 질문에 답하는가)
  - 소속: health summary / trend signal / top offender view
- **원칙:** 신호 수를 과도하게 늘리지 않는다. **작고, signal-to-noise가 높은 모니터링 모델**을 목표로 한다.
- **Output:** `04-engineering/`에 core signal list 및 필요 시 보조 노트.

### For Experimenter

- **목표:** 소수 신호가 **현재 Prometheus 설정**에서 실제로 조회 가능한지 확인한다.
- **문서화할 제한:** Prometheus retention, query 복잡도, 리소스 사용으로 인한 제약. 각 제한이 어떤 신호·쿼리에 영향을 주는지 명시.
- **Output:** `05-experiment/`에 실험 노트 및 제한 사항 문서.

### For Reviewer

- **핵심 평가 질문:** “운영자가 이 dashboard(모델)를 보고 **5–10분 안에** 클러스터 건강을 이해할 수 있는가?”
- **검토 항목:** 모델이 리더십 질문(안전한가 / 곧 불안전해질 징후는 / 어디를 먼저 볼 것인가)에 답하는가, core signal list의 완전성·실용성, 실험에서 드러난 제한이 문서에 반영되었는가. 누락·리스크·개선 제안 기록.
- **Output:** `06-review/`에 review summary 및 체크리스트.

### For Writer

- **최종 문서 제목:** **Kubernetes Cluster Health Monitoring Model**
- **내용:** 모델 개요, 세 그룹(Health Summary, Trend/Risk, Top Offenders), 카테고리, core signal list 요약(또는 본문/부록), Prometheus 제한 사항, 사용 가이드(5–10분 파악 방법). 이 문서는 이후 프로젝트 **sre-monitoring-dashboard-design**, **sre-monitoring-alert-policy**, **sre-monitoring-operational-runbooks**의 기반이 된다.
- **대상 독자:** 클러스터를 운영하는 SRE·플랫폼 엔지니어.
- **문서화 언어:** 설명은 한국어, 기술 식별자(metric, PromQL, 도구명 등)는 영어 유지.
- **Output:** `07-documentation/`에 최종 문서 및 필요 시 요약·부록.
