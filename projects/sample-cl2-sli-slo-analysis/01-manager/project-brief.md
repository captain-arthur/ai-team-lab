# Project brief: ClusterLoader2 부하 테스트 기반 SLI/SLO 분석

**Project:** sample-cl2-sli-slo-analysis  
**Source:** `tasks/sample-cl2-sli-slo-analysis.md`  
**Role:** Manager  
**Reference:** knowledge/principles (cat-vision, cat-design-principles, sli-slo-philosophy, devcat-program-brief)

---

## 1. Project brief

### Goal

ClusterLoader2 부하 테스트가 산출하는 메트릭을 분석하고, **소규모 Kubernetes 클러스터 수용 테스트**에 의미 있는 메트릭·SLI 후보·SLO 정당화 방안을 도출하여, devcat 실험 설계·SLO 검증·정제로 이어지는 구체적 기초를 마련한다.

### Scope

| In scope | Out of scope |
|----------|--------------|
| ClusterLoader2 메트릭 종류·의미 파악 및 수용 테스트 관점에서 의미 있는 것 식별 | devcat 저장소 대규모 코드 변경·새 도구 전면 도입 |
| “왜 중요한가”(사용자·워크로드 관점) 정리 | 대규모 벤치마크·성능 최적화 자체가 목표인 작업 |
| SLI 후보 선정, Kubernetes scalability 철학·“you promise, we promise” 반영 | Conformance/기능 테스트 상세 설계 |
| 소규모 클러스터용 SLO 정당화(user expectations, environment assumptions, repeatable measurement, practical acceptance meaning) | |
| 벤치마크 스타일 임계값의 소규모 수용 부적합성 전제 | |
| devcat 실험을 통한 SLO 검증·정제 방안(실험 설계·해석 제안) | |
| (선택) kind/Docker 로컬 검증 경로 안내 | |

### Constraints

- knowledge/principles의 CAT 원칙·devcat-program-brief 준수.
- ClusterLoader2·devcat 현재 사용 방식(config.yaml, 클러스터별 오버라이드, results/) 전제.
- 소규모 클러스터 수용 테스트 맥락 유지; 대규모 벤치마크 임계값 무비판 재사용 배제.

### Success criteria

- ClusterLoader2 메트릭 중 수용 테스트에 의미 있는 것과 그 이유가 문서로 정리되어 있다.
- SLI 후보 목록과 선정 근거가 있다.
- 소규모 클러스터용 SLO 값 정당화 방법이 네 가지 기준(user expectations, environment assumptions, repeatable measurement, practical acceptance meaning)으로 정리되어 있다.
- Kubernetes scalability 철학 및 “you promise, we promise”를 반영한 정당화 원칙이 요약되어 있다.
- devcat 실험으로 SLO를 검증·정제하는 방안(실험 설계·해석)이 제안되어 있다.
- 산출물이 research → experiment → interpretation → devcat improvement 모델과 연결 가능한 형태로 되어 있다.

---

## 2. Work breakdown

| # | Sub-task | Owner | Depends on |
|---|----------|--------|------------|
| 1 | ClusterLoader2 메트릭 조사: 산출 지표 종류·의미·문서/소스 정리 | Researcher | — |
| 2 | 수용 테스트 관점 의미 메트릭·SLI 후보·정당화 기준 정리; 벤치마크 임계값 한계 반영 | Researcher | 1 |
| 3 | SLI/SLO 정당화 체계와 devcat 실험 연동 구조 설계(실험 설계·결과 해석·SLO 정제 흐름) | Architect | 2 |
| 4 | 분석 요약·SLI 후보표·SLO 정당화 템플릿·(선택) 실험 runbook/스크립트 초안 | Engineer | 3 |
| 5 | CAT 원칙·intake 질문·devcat 현실 반영 여부 검증 | Reviewer | 4 |
| 6 | 최종 보고서·사용자 문서(메트릭·SLI·SLO 정당화·실험 방안 요약) | Writer | 5 |
| 7 | knowledge/에 SLI/SLO·ClusterLoader2·devcat 실험 관련 인사이트 추출 | Knowledge Extraction | 6 |

---

## 3. Role handoffs

### 핵심 질문 (모든 단계가 이 질문에 답하도록 유도)

- **Q1.** ClusterLoader2가 산출하는 메트릭 중 **수용 테스트에 의미 있는 것은 무엇인가?**
- **Q2.** 그 메트릭들이 **사용자·워크로드 관점에서 왜 중요한가?**
- **Q3.** 어떤 메트릭을 **SLI 후보**로 둘 것인가?
- **Q4.** **소규모 클러스터**에서 SLO 값을 어떻게 **정당화**할 것인가?
- **Q5.** **devcat 실험**으로 그 SLO를 어떻게 **검증·정제**할 수 있는가?

---

### For Researcher

- **답할 질문:** Q1, Q2, Q3. ClusterLoader2 메트릭 종류·출처(문서, 결과 파일 형식), 수용 테스트 관점에서 “의미 있다”의 기준, 사용자·워크로드 관점 정당화, SLI 후보 후보군과 선정 근거.
- **참고:** Kubernetes scalability philosophy (precise and well-defined, consistent, user-oriented, testable), “you promise, we promise”. sli-slo-philosophy.md의 user expectations, environment assumptions, repeatable measurement, practical acceptance meaning.
- **명시:** 벤치마크 스타일 임계값이 소규모 수용 테스트에 부적합할 수 있음을 전제로 분석.
- **출력:** `02-research/`에 메트릭 분석·의미 있는 메트릭·SLI 후보·정당화 기준 초안.

### For Architect

- **설계 범위:** SLI/SLO 정당화 “체계”(어떤 기준으로 SLO를 선정·기록할지), devcat 실험과의 연동 구조(실험 설계 → ClusterLoader2 실행 → 결과 해석 → SLO 검증·정제 흐름). 구현 코드나 devcat 리포 구조 변경 상세는 Engineer로.
- **결정 사항:** SLI 후보를 SLO로 넘길 때의 정당화 템플릿 형태, 실험 한 번의 입력/출력/판단 기준 개념 구조.
- **비기능:** devcat-program-brief의 현재 현실(ClusterLoader2, config.yaml, results/)을 전제로, 점진적 개선 가능한 구조.
- **출력:** `03-architecture/`에 정당화 체계·실험 연동 구조 문서.

### For Engineer

- **산출물:** 메트릭·SLI 후보 요약표, SLO 정당화용 템플릿(또는 체크리스트), devcat 실험 설계·해석 방법을 담은 runbook 또는 단계 설명. (선택) kind/Docker 로컬 검증 절차 요약. 코드는 “예제·스크립트 초안” 수준, devcat 대규모 변경 제외.
- **기술 힌트:** ClusterLoader2 결과 형식, devcat results/ 구조, 기존 config.yaml·오버라이드 패턴 유지.
- **출력:** `04-engineering/`에 문서·템플릿·(선택) 스크립트/runbook.

### For Reviewer

- **검증:** intake의 다섯 가지 핵심 질문(Q1~Q5)에 대한 답이 산출물에 반영되어 있는지, CAT 원칙·devcat-program-brief·sli-slo-philosophy와 일치하는지, 소규모 클러스터·벤치마크 임계값 무비판 재사용 배제가 지켜졌는지, devcat 실험 연동이 구체적으로 제안되었는지.
- **출력:** `05-review/`에 검토 요약·체크리스트·이슈·제안.

### For Writer

- **독자:** 플랫폼/SRE, devcat 다음 단계(실험·SLO 정제)를 수행할 사람.
- **형식:** 최종 보고서(Markdown), 요약·메트릭·SLI·SLO 정당화·실험 방안·한계·follow-up. 필요 시 사용자용 요약 또는 “다음에 할 일” 안내.
- **필수 포함:** Q1~Q5에 대한 답 요약, Kubernetes scalability·“you promise, we promise” 반영 여부, devcat 실험 검증·정제 방안 요약, Reviewer 의견 반영(한계·권장 사항).
- **출력:** `06-documentation/`에 final report 및 사용자 문서.

### For Knowledge Extraction

- **추출 대상:** SLI/SLO 정당화 원칙, ClusterLoader2 메트릭·수용 테스트 관점 요약, 소규모 클러스터 SLO 정당화 패턴, devcat 실험 연동 시 유의점. knowledge/principles, knowledge/patterns, knowledge/lessons 중 적절한 위치에 반영.
- **출력:** `07-knowledge-extraction/` 및 knowledge/ 하위 업데이트.

---

## 4. devcat 연동

- **후속 단계(Research, Architecture, Engineering)** 에서는, **devcat 저장소를 사용한 실험**을 설계할 수 있다.
- 예: devcat에서 ClusterLoader2 테스트를 실행(config.yaml, 클러스터별 오버라이드)하고, results/에 쌓인 결과를 분석해 메트릭·SLI 후보·SLO 후보 값을 검증하거나 정제하는 실험 절차·해석 방법을 제안할 수 있다.
- 이때 설계는 devcat-program-brief의 **현재 현실**(ClusterLoader2, perfdash, config.yaml, ol-test.yaml, results/)을 전제로 하며, **점진적 개선**으로 이어지도록 한다. 연구만으로 끝나지 않고 **research → experiment → interpretation → devcat improvement** 흐름에 맞춘다.
