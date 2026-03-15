# Task Intake: Kubernetes Cluster Health Monitoring Model

**Created:** 2025-03-15  
**Status:** Ready

---

## Task Title

Design a practical **Kubernetes Cluster Health Monitoring Model** for daily operational awareness.

---

## Problem Description

운영자가 "지금 클러스터가 안전한가?", "곧 불안전해질 징후는 없는가?", "문제라면 어디를 먼저 봐야 하나?"를 빠르게 답할 수 있어야 한다. 현재도 Prometheus와 Grafana로 모니터링하고 수많은 dashboard가 있으나, 일상적인 운영 인지에는 항상 실용적이지 않다. 리더십 의도는 운영자가 매일 보고 **5–10분 안에** 클러스터 건강 상태를 파악할 수 있는 "중앙 dashboard"에 가까운 것을 만드는 것이다.

---

## Goal

**"운영 건강을 나타낼 수 있는 최소한의 신호 집합은 무엇인가?"**에 답하는, 실무에 바로 쓸 수 있는 **Cluster Health Monitoring Model**을 설계한다. 수십 개 패널을 훑을 필요 없이, 소수의 상위 신호로 클러스터 상태를 요약하고, 트렌드·초기 리스크·드릴다운(TOP10 등)을 구분할 수 있어야 한다.

---

## Scope

- **In scope:**  
  - Cluster Health Monitoring Model 설계(신호 분류, 상위 건강 신호, 초기 리스크 신호, drill-down/TOP10 뷰, Prometheus에서 실제 조회 가능한 신호).  
  - 핵심 신호 목록(core signal list) 및 각 신호의 의미·운영 질문·Prometheus metric/PromQL 예시.  
  - 현재 Prometheus 환경에서의 조회 가능성 실험 및 제한 사항 문서화.  
  - 최종 문서: **Kubernetes Cluster Health Monitoring Model** — 이후 dashboard 설계, alert 정책, runbook 프로젝트의 기반.
- **Out of scope:**  
  - 수많은 dashboard 제작 또는 대량 metric 나열.  
  - Prometheus/Grafana 자체 도입 또는 인프라 변경.  
  - Alert 정책·runbook 상세 설계(후속 프로젝트).

---

## Expected Deliverables

- **Cluster Health Monitoring Model** 문서: 신호 분류, 상위 건강 요약, 트렌드/리스크 지표, Top Offenders/Drill-down 뷰 정의.
- **Core signal list:** 각 신호별 signal name, Prometheus metric 또는 PromQL 예시, "왜 중요한가", "어떤 운영 질문에 답하는가", health summary / trend / top offender 구분.
- 실험 노트: 일부 신호의 Prometheus 조회 가능성, retention·query 복잡도·리소스 사용으로 인한 제한.
- 리뷰 요약: "5–10분 안에 클러스터 건강을 파악할 수 있는가?"에 대한 평가.

---

## Constraints

- **모든 것을 모니터링하지 않는다.** 최소 집합에 집중.
- 운영자가 수십 개 패널을 훑을 필요가 없어야 함.
- Prometheus 리소스 제약으로 장기·복잡 쿼리에 한계가 있음.
- 일부 신호(예: Calico readiness probe 실패)는 알림을 유발하지만 실제 서비스 영향과 직결되지 않을 수 있음 — 노이즈 고려.
- TOP10 형태의 고부하 워크로드 뷰 검토.
- 트렌드 및 초기 리스크 신호를 고려해야 함.

---

## Priority

High — 일상 운영 인지와 초기 리스크 감지를 위한 기반 모델.

---

## Additional Context

- **Program:** SRE Monitoring. 후속 프로젝트: `sre-monitoring-dashboard-design`, `sre-monitoring-alert-policy`, `sre-monitoring-operational-runbooks`.
- **Audience:** 클러스터를 운영하는 SRE·플랫폼 엔지니어.
- **Documentation:** 설명은 한국어, 기술 식별자(metric 이름, PromQL, 도구명 등)는 영어 유지.

---

*Template: `tasks/intake-template.md`*
