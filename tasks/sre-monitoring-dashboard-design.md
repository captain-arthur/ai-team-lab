# Task Intake: Central Kubernetes Operational Dashboard Design

**Created:** 2025-03-15  
**Status:** Ready

---

## Task Title

Central Kubernetes Operational Dashboard 설계 — 운영 확신(operational confidence)과 조기 리스크 감지를 위한 대시보드 디자인.

---

## Problem Description

Prometheus와 Grafana 기반 대시보드가 이미 많지만, 운영자는 여전히 **“이게 정상인가?”, “클러스터가 실제로 안전한가?”, “누가 확인해 줄 수 있나?”** 같은 질문을 자주 한다. 이는 **핵심적인 운영 문제**를 드러낸다. metric은 많지만 **명확한 운영 확신 모델(operational confidence model)** 이 부족하다. 목표는 “또 하나의 큰 대시보드”가 아니라, **운영자가 두 가지 질문에 확신 있게 답할 수 있는 Central Kubernetes Operational Dashboard**를 설계하는 것이다.

---

## Goal

**Central Kubernetes Operational Dashboard** 를 설계하여:

1. **지금 클러스터가 안전하다고 확신할 수 있는가?**
2. **클러스터가 곧 불안전해질 조기 징후가 보이는가?**

를 운영자가 답할 수 있게 하고, 이상이 보일 때 **어디를 조사할지** 빠르게 파악할 수 있게 한다. 대시보드는 **운영 확신(operational confidence)** 과 **조기 리스크 감지**에 초점을 두며, metric 수가 아니라 **판단 가능한 최소 신호**로 구성한다.

---

## Scope

- **In scope:**  
  - 기존 **Cluster Health Monitoring Model**(`projects/sre-monitoring-cluster-health-model`)을 기반으로 한 Central Dashboard **설계 문서** 작성.  
  - 대시보드 **레이아웃·신호 계층 구조** 정의(Summary / Trend-Risk / Drill-down).  
  - **조기 리스크 신호** 강조: CPU throttling, OOM risk, ingress controller stress, node resource pressure, pending pods 등.  
  - 5–10분 안에 클러스터 건강을 파악할 수 있는 구조, 대시보드 과부하 방지.  
- **Out of scope:**  
  - Grafana 패널 실제 구현·JSON export. (설계만, 구현은 별도 작업.)  
  - Alert 정책·runbook 상세(후속 프로젝트).

---

## Expected Deliverables

- **Central Kubernetes Operational Dashboard 설계 문서:**  
  - 운영 확신과 조기 리스크 감지를 위한 목표·원칙.  
  - 모니터링 모델과의 매핑(Summary / Trend / Top Offenders).  
  - 레이아웃(상단·중간·하단 또는 탭) 및 신호 계층.  
  - 조기 리스크 신호(CPU throttling, OOM risk, ingress stress, node pressure, pending pods) 배치·표현 방식.  
  - 5–10분 사용 흐름·조사 진입점.  
- (선택) 대시보드 스케치·와이어프레임 설명 또는 패널 목록.

---

## Constraints

- **운영 확신 중심:** metric 양이 아니라 “지금 안전한가?” / “곧 위험한가?”에 답할 수 있는지가 기준.
- **메인 뷰 신호 최소화:** 한 화면에 나열하는 신호 수를 줄여 “이상/정상”이 한눈에 보이게.
- **조기 리스크 강조:** CPU throttling, OOM risk, ingress controller stress, node resource pressure, pending pods 등을 눈에 띄게 배치.
- **모니터링 모델 준수:** `sre-monitoring-cluster-health-model` 의 Summary / Trend / Top Offenders 구조를 기반으로 하되, 대시보드 UX로 구체화.

---

## Priority

High — 운영 확신과 조기 리스크 감지는 일상 운영의 핵심.

---

## Additional Context

- **Program:** SRE Monitoring.
- **Foundation:** `projects/sre-monitoring-cluster-health-model` — Kubernetes Cluster Health Monitoring Model, core signal list, final-report, experiment 제한 사항.
- **후속:** sre-monitoring-alert-policy, sre-monitoring-operational-runbooks.
- **Documentation:** 설명은 한국어, 기술 식별자(metric, PromQL, Grafana 등)는 영어 유지.

---

*Template: `tasks/intake-template.md`*
