# Task Intake: ClusterLoader2 부하 테스트와 SLI/SLO 분석

**Created:** 2025-03  
**Status:** Ready

---

## Task Title

ClusterLoader2 부하 테스트 기반 소규모 클러스터 수용 테스트를 위한 SLI/SLO 분석 및 정당화

---

## Problem Description

devcat은 ClusterLoader2 기반 부하 시나리오를 실행하고 결과를 results/에 저장하지만, **어떤 지표가 수용 테스트에 의미 있는지**, **그 지표를 어떻게 SLO로 정당화할지**가 아직 명확하지 않다. 도구가 주는 벤치마크 스타일 임계값을 그대로 쓰면 소규모 클러스터 수용 판단에 맞지 않을 수 있다. 이 프로젝트는 ClusterLoader2가 산출하는 지표를 분석하고, 소규모 Kubernetes 클러스터에서 실용적인 수용 테스트를 지원하기 위해 **어떤 메트릭이 중요하고, 왜 중요한지**, **SLI 후보와 SLO 값 정당화 방법**을 정리하는 것을 목표로 한다.

---

## Goal

ClusterLoader2 부하 테스트를 이해하고, **소규모 Kubernetes 클러스터에서의 실용적 수용 테스트**를 지원할 수 있도록, 의미 있는 메트릭·SLI 후보·SLO 정당화 방안을 도출한다. 결과는 devcat 시스템 진화와 research → experiment → interpretation → devcat improvement 작업 모델에 연결될 수 있어야 한다.

---

## Scope

- **In scope:**
  - ClusterLoader2가 산출하는 메트릭 종류·의미 파악.
  - 수용 테스트 관점에서 의미 있는 메트릭 식별 및 “왜 중요한가”(사용자·워크로드 관점) 정리.
  - SLI 후보 선정 및 Kubernetes scalability 철학(precise and well-defined, consistent, user-oriented, testable) 및 “you promise, we promise” 개념에 맞는 정당화 방안.
  - 소규모 클러스터용 SLO 값 정당화: user expectations, environment assumptions, repeatable measurement, practical acceptance meaning 활용.
  - 벤치마크 스타일 임계값이 소규모 수용 테스트에 부적합할 수 있음을 전제로 한 분석.
  - devcat 실험을 통해 SLO 값을 검증·정제하는 방안 고려.
  - (선택) 로컬 Docker/kind 클러스터를 이용한 검증 가능성 언급 또는 간단 실험.

- **Out of scope:**
  - devcat 저장소 코드 대규모 변경 또는 새 도구 전면 도입.
  - 대규모 벤치마크·성능 최적화 자체가 목표인 작업.
  - Conformance/기능 테스트 상세 설계(본 프로젝트는 부하·메트릭·SLO에 초점).

---

## Expected Deliverables

- ClusterLoader2 메트릭 중 수용 테스트에 의미 있는 것과 그 이유(사용자·워크로드 관점)를 정리한 문서.
- SLI 후보 목록 및 선정 근거.
- 소규모 클러스터용 SLO 값 정당화 방법( user expectations, environment assumptions, repeatable measurement, practical acceptance meaning 기준).
- Kubernetes scalability 철학 및 “you promise, we promise”를 반영한 정당화 원칙 요약.
- devcat 실험으로 SLO를 검증·정제하는 방안(실험 설계·해석 방법 제안).
- (선택) kind/Docker 로컬 검증 경로에 대한 짧은 안내 또는 제안.

---

## Constraints

- knowledge/principles에 정의된 CAT 원칙 및 devcat-program-brief를 따른다 (cat-vision, cat-design-principles, sli-slo-philosophy, devcat-program-brief).
- ClusterLoader2·devcat 현재 사용 방식(config.yaml, 클러스터별 오버라이드, results/)을 전제로 한다.
- 소규모 클러스터 수용 테스트 맥락을 유지한다(대규모 벤치마크 임계값 무비판 재사용 배제).

---

## Priority

High — devcat 진화와 실용 CAT 시스템 완성에 직접 연결되는 기초 작업이다.

---

## Additional Context

- **참조 원칙:** Kubernetes scalability philosophy (precise and well-defined, consistent, user-oriented, testable), “you promise, we promise”. knowledge/principles/sli-slo-philosophy.md 참고.
- **SLO 정당화:** 도구 기본값·벤치마크 스타일 임계값은 소규모 클러스터 수용 테스트에 그대로 쓰지 않는다. SLO 후보는 user expectations, environment assumptions, repeatable measurement, practical acceptance meaning으로 정당화한다.
- **작업 모델:** research → experiment → interpretation → devcat improvement. 본 태스크의 산출은 devcat 실험 설계·결과 해석·SLO 정제로 이어질 수 있어야 한다.
- **선택 검증:** 로컬 Docker/kind 클러스터로 메트릭·SLO 후보를 검증하는 경로가 실용적이면 포함한다.
- **산출 형식:** Markdown 문서; 필요 시 표·목록으로 정리. 플랫폼/SRE가 읽고 devcat 다음 단계에 활용할 수 있는 수준.
