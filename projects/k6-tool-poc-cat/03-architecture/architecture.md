# Architecture: k6 동작 구조(도구 내부 모델)

**Date:** 2026-03-18

## k6 내부 동작 구조(요약)
- **입력**: JS 스크립트(시나리오/요청/체크/threshold)
- **실행 엔진**:
  - **Scenario executor**가 목표 부하(동시성 또는 도착률)를 만들고
  - **VU**들이 코드를 반복 실행하며 요청을 발생
- **관측/판정**:
  - 요청 결과에서 **metric 샘플**을 생성
  - 실행 종료 시점에 **집계(aggregations)** 후 threshold 평가로 **PASS/FAIL** 결정

## 요청 → 응답 → metric 생성 흐름
- Step 1: VU가 HTTP 요청 전송
- Step 2: 응답 수신(상태코드/타이밍)
- Step 3: k6가 기본 metric 샘플 기록
  - 예: `http_req_duration`(요청 지연), `http_req_failed`(실패 여부), `http_reqs`(요청 수)
- Step 4: 스크립트의 `check()`가 성공/실패를 metric(`checks`)로 기록

## metric 집계 방식(POC에서 사용하는 최소)
- **Latency**: `http_req_duration`의 `p(95)`(샘플 분포 기반)
- **Error rate**: `http_req_failed`의 `rate`(실패/전체 비율)
- **Throughput**: `http_reqs`의 `rate`(초당 요청 수, 관측된 달성치)

## threshold 동작 방식(PASS/FAIL)
- **정의**: metric 집계 값에 대한 조건(예: `p(95)<300`, `rate<0.001`)
- **평가**: 실행 종료 시(또는 설정 시 조기 중단) 조건을 평가
- **결과**:
  - 조건을 모두 만족하면 **PASS(종료 코드 0)**
  - 하나라도 위반하면 **FAIL(비영(非0) 종료 코드, 위반 metric을 출력)**
