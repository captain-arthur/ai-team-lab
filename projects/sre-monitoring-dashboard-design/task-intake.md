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

를 운영자가 답할 수 있게 하고, 이상이 보일 때 **어디를 조사할지** 빠르게 파악할 수 있게 한다. 대시보드는 **운영 확신**과 **조기 리스크 감지**에 초점을 두며, metric 수가 아니라 **판단 가능한 최소 신호**로 구성한다.

---

## Scope

- **In scope:** 기존 Cluster Health Monitoring Model 기반 Central Dashboard 설계 문서, 레이아웃·신호 계층 정의, 조기 리스크 신호 강조(CPU throttling, OOM risk, ingress stress, node pressure, pending pods), 5–10분 파악 구조, 과부하 방지.
- **Out of scope:** Grafana 패널 실제 구현·JSON. Alert 정책·runbook 상세(후속 프로젝트).

---

## Expected Deliverables

- Central Kubernetes Operational Dashboard 설계 문서: 목표·원칙, 모니터링 모델 매핑, 레이아웃·신호 계층, 조기 리스크 배치·표현, 5–10분 사용 흐름·조사 진입점. (선택: 스케치·패널 목록.)

---

## Constraints

- 운영 확신 중심. 메인 뷰 신호 최소화. 조기 리스크 강조. 모니터링 모델(sre-monitoring-cluster-health-model) 준수.

---

## Priority

High.

---

## Additional Context

- **Program:** SRE Monitoring. **Foundation:** projects/sre-monitoring-cluster-health-model. **Documentation:** 한국어 설명, 기술 식별자 영어.

---

*Template: `tasks/intake-template.md`*
