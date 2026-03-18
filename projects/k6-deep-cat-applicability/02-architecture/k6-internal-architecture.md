# k6 내부 아키텍처(어떻게 동작하는가)

**Date:** 2026-03-18  
**대상 버전:** k6 v1.6.1 (로컬 실행 기준)

## 1) 구성요소(핵심 개념)
- **VU (Virtual User)**: 스크립트를 실행하는 실행자(동시성 단위).
- **Iteration**: VU가 “한 번” 실행한 사용자 여정(스크립트 함수 1회 수행).
- **Scenario**: “어떤 executor로 어떤 부하 형태를 언제까지 만들지”를 정의한 실행 단위.
- **Executor**: scenario를 실제 부하로 바꾸는 방식(예: `constant-vus`, `constant-arrival-rate`, `ramping-*`).

## 2) 실행 매커니즘(스크립트→실행→metric→집계→threshold)
1. **스크립트 로딩(init 단계)**  
   - 옵션(options: scenarios/thresholds 등) 확정
   - 환경변수(`__ENV`)로 target/부하/임계값 주입
2. **Scenario 실행(executor 작동)**  
   - VU 기반: “VU 수”를 목표로 올리고, 각 VU가 iteration을 반복(요청률은 결과로 결정)
   - arrival-rate 기반: “목표 iteration rate”를 맞추려고 VU를 늘리며 추종(달성 실패 시 dropped_iterations 발생)
3. **요청/응답 처리**  
   - 요청 단위로 timing/상태가 수집되며 metric 샘플이 생성됨
4. **metric 집계(요약 값 생성)**  
   - Trend 계열(예: `http_req_duration`)은 p95/p99 등 백분위를 포함해 집계
   - Rate/Counter(예: `http_req_failed`, `http_reqs`)는 rate/count로 요약
5. **threshold 평가(SLO 단언)**  
   - 집계 값이 조건을 만족하면 PASS, 위반하면 FAIL(+비0 종료 코드)

## 3) 핵심 원리(운영자가 오해하기 쉬운 지점)
### concurrency 모델
- VU는 “동시 실행자”이며, **VU 기반 시나리오에서는 처리량이 목표가 아니라 결과**다.
- think time(`sleep`)이 있으면 VU당 iteration이 느려져 처리량이 내려간다.

### arrival-rate vs VU 모델 차이(의미)
- **VU 모델**: “동시 사용자 수”를 고정/변화시키며, RPS는 시스템+스크립트에 의해 결정됨.
- **arrival-rate 모델**: “목표 RPS/iteration rate”를 주입하고, 달성하려고 VU를 자동 확장함.  
  - 달성 실패는 `dropped_iterations` 같은 신호로 드러난다(포화의 간접 신호).

### metric aggregation 방식(요약 JSON 기준)
- `http_req_duration`: p(90)/p(95)/p(99) 등 백분위로 지연 분포를 읽는다.
- `http_req_failed`: rate(실패율)로 읽는다(0이면 실패 샘플이 없음).
- `http_reqs`: rate(달성 처리량)로 읽는다.
- `dropped_iterations`: arrival-rate가 목표를 못 맞춘 정도(“요청을 못 보낸 것”)를 읽는다.
