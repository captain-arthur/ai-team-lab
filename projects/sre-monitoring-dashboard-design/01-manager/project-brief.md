# Project brief: Central Kubernetes Operational Dashboard Design

**Project:** sre-monitoring-dashboard-design  
**Source:** `task-intake.md`  
**Role:** Manager  
**Program:** SRE Monitoring

---

## 1. Project brief

### Goal

**Central Kubernetes Operational Dashboard** 를 설계하여, 운영자가 **운영 확신(operational confidence)** 을 갖고 다음 두 질문에 답할 수 있게 한다.

1. **지금 클러스터가 안전하다고 확신할 수 있는가?**
2. **클러스터가 곧 불안전해질 조기 징후가 보이는가?**

또한 이상이 보일 때 **어디를 조사할지** 빠르게 찾을 수 있도록 한다. 목표는 “또 하나의 큰 대시보드”가 아니라, **metric 양이 아닌 판단 가능한 최소 신호**로 구성된 **Central Operational Dashboard** 의 설계 문서를 만드는 것이다.

### Scope

| In scope | Out of scope |
|----------|--------------|
| Cluster Health Monitoring Model 기반 Central Dashboard **설계 문서** | Grafana 패널 실제 구현·JSON export |
| 레이아웃·신호 계층(Summary / Trend-Risk / Drill-down) 정의 | Alert 정책·runbook 상세(후속) |
| 조기 리스크 신호 강조: CPU throttling, OOM risk, ingress controller stress, node resource pressure, pending pods | |
| 5–10분 내 클러스터 건강 파악 구조, 대시보드 과부하 방지 | |
| (선택) 대시보드 스케치·와이어프레임 또는 패널 목록 | |

### Constraints

- **운영 확신 중심:** “지금 안전한가?” / “곧 위험한가?”에 답할 수 있는지가 기준. metric 수가 아님.
- **메인 뷰 신호 최소화:** 한 화면에 나열하는 신호를 줄여 “이상/정상”이 한눈에 보이게.
- **조기 리스크 강조:** CPU throttling, OOM risk, ingress stress, node pressure, pending pods 등을 눈에 띄게 배치.
- **모니터링 모델 준수:** `projects/sre-monitoring-cluster-health-model` 의 Kubernetes Cluster Health Monitoring Model(Summary / Trend / Top Offenders)을 기반으로 설계. core-signal-list, final-report, experiment 제한 사항을 반영.

### Success criteria

- **Central Kubernetes Operational Dashboard 설계 문서**가 존재하며, 다음을 포함한다:  
  - 운영 확신·조기 리스크 감지 목표 및 설계 원칙.  
  - 모니터링 모델과의 매핑(어떤 신호를 어느 영역에 배치하는지).  
  - 레이아웃(상단 Summary, 중간 Trend-Risk, 하단/탭 Drill-down) 및 신호 계층.  
  - 조기 리스크 신호(CPU throttling, OOM risk, ingress stress, node pressure, pending pods)의 배치·표현 방식.  
  - 5–10분 사용 흐름 및 이상 시 조사 진입점.  
- 설계만으로도 이후 Grafana 구현 단계에서 “어디에 어떤 패널을 둘지” 결정할 수 있는 수준.

---

## 2. Work breakdown

| # | Sub-task | Owner | Depends on |
|---|----------|--------|------------|
| 1 | 입력 정리: Cluster Health Monitoring Model(final-report, core-signal-list, experiment 제한) 요약 및 “설계에 반영할 신호·제약” 목록 작성 | Researcher 또는 Architect | — |
| 2 | 설계 원칙·목표 정리: 운영 확신, 조기 리스크, 신호 최소화, 5–10분 파악 | Architect | 1 |
| 3 | 레이아웃·신호 계층 설계: Summary 블록 / Trend-Risk 블록 / Drill-down 블록(또는 탭), 각 블록별 패널 역할 | Architect | 2 |
| 4 | 조기 리스크 신호 배치: CPU throttling, OOM risk, ingress stress, node pressure, pending pods를 Trend-Risk(및 필요 시 Summary)에 어떻게 넣을지 | Engineer 또는 Architect | 3 |
| 5 | 운영 흐름·조사 진입점: 5–10분 사용 시나리오, “이상 시 어디로 들어가는지” 설명 | Writer 또는 Architect | 3, 4 |
| 6 | 설계 문서 통합: Central Kubernetes Operational Dashboard 설계 문서 최종화 | Writer | 5 |
| 7 | (선택) Knowledge Extraction: 설계 패턴·원칙을 `knowledge/`에 추출 | Knowledge Extraction | 6 |

*이 프로젝트는 모니터링 모델이 이미 있으므로 Research 단계를 짧게 하거나, Architect가 모델 요약과 설계를 한꺼번에 할 수 있다. Experiment 단계는 생략 가능(모델 검증은 sre-monitoring-cluster-health-model에서 완료).*

### Handoff: 설계에 사용할 입력

- **`projects/sre-monitoring-cluster-health-model/06-documentation/final-report.md`** — Kubernetes Cluster Health Monitoring Model 요약, 세 가지 질문, Summary/Trend/Top Offenders, core signals, 제약, v1 한계.
- **`projects/sre-monitoring-cluster-health-model/04-engineering/core-signal-list.md`** — Summary 6개, Trend 4~5개, Top Offenders 6개 뷰, Dashboard 설계자용 요약표(패널 배치 제안).
- **`projects/sre-monitoring-cluster-health-model/05-experiment/experiment-notes.md`** — 즉시 사용 가능/수정 후/불가 신호, Prometheus 제한, Managed K8s 시 축소 구성.

---

## 3. Handoff notes

### For Architect (또는 설계 주도 역할)

- **설계 범위:** Central Dashboard **한 편**의 설계. Grafana 구현은 하지 않음.
- **반드시 반영:**  
  - 모니터링 모델의 **세 레이어**(Cluster Health Summary / Trend-Risk / Top Offenders)를 대시보드 **레이아웃**으로 옮기기.  
  - **조기 리스크** 강조: CPU throttling, OOM risk, ingress controller stress, node resource pressure, pending pods. 이들은 모델의 Trend-Risk(및 일부 Summary)에 대응하므로, “Early risk” 또는 “Trend / Risk” 블록에서 **눈에 띄게** 배치할 것.  
  - **메인 뷰 신호 최소화:** Summary는 5~7개(또는 환경별 4~5개)만 상단에. 한눈에 “전부 정상 = 안전”이 보이게.
- **환경 차이:** Experiment에서 control plane metric이 없는 환경을 고려했으므로, 설계 문서에 **“Managed K8s 등에서 Summary 패널 구성 조정”** 옵션을 한 줄이라도 넣어 두면 좋음.
- **Output:** 설계 문서 초안(목표·원칙·레이아웃·신호 계층·조기 리스크·운영 흐름).

### For Writer

- **최종 문서 제목:** Central Kubernetes Operational Dashboard 설계 문서 (또는 동일한 제목의 단일 설계 문서).
- **포함 내용:** 목표·원칙, 모니터링 모델 매핑, 레이아웃·신호 계층, 조기 리스크 배치·표현, 5–10분 사용 흐름·조사 진입점, (선택) 스케치·패널 목록.  
- **대상 독자:** 대시보드를 구현할 플랫폼/SRE 엔지니어, 운영자.
- **문서화 언어:** 설명은 한국어, 기술 식별자(metric, PromQL, Grafana 등)는 영어.
- **Output:** `07-documentation/` 또는 프로젝트 내 설계 문서가 들어갈 폴더에 최종 설계 문서.

---

*Project brief. Foundation: sre-monitoring-cluster-health-model.*
