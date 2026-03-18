# Engineering(계획): 최소 custom CAT 시나리오 실험 스펙 후보(구현 없음)

**Date:** 2026-03-18

본 문서는 “실제로 다음 단계(구현/실험)로 넘어가기 위한 최소 스펙”만 적는다.

## 공통 실험 스펙(모든 후보에 공통)
- **CAT 최소 조건 매핑**
  - Scenario Injection: 외부 입력(config/env/flags)으로 scenario 조건을 바꾼다
  - SLI Measurement: 테스트가 수집/산출하는 값(키-값)으로 SLI를 만든다
  - SLO Evaluation: assertion(단언)이 PASS/FAIL로 결정되게 만든다
  - Result Persistence: `results/<run-id>/cat-result.json` 같은 파일로 남긴다
- **PASS/FAIL 권위**
  - 최종 결과는 “테스트 실행 종료 코드 + cat-result.json”을 함께 쓴다
  - CAT는 재판정하지 않는다(기록 중심)

## 후보 A: 내부 측정 + 단언 + JSON 저장(완전 커스텀)
- 시나리오
  - 테스트 입력으로 `target_mode`(예: “fast/slow”), `sample_size`를 받는다
- SLI
  - 예: “요청 지연(시뮬레이션 또는 Go 함수 실행 시간)”의 p95/p99 또는 성공률
- SLO
  - 예: p95 < X ms, error_rate < Y
- 왜 이게 중요한가
  - Ginkgo가 “프레임워크 제약 없이” SLI 산출/단언/파일 저장을 묶어낼 수 있는지 직접 본다
- k6/CL2 대비 차이
  - k6/CL2가 측정 모델을 제공한다면, 여기서는 “측정/판정/저장”을 테스트 코드가 전부 통제한다(비용/제약 확인)

## 후보 B: Kubernetes 상태 기반 scenario + SLO 단언 + 결과 파일
- 시나리오
  - 외부 입력으로 “어떤 오브젝트를 만들고/기다릴지”와 “성공 조건(예: ready time 상한)”을 받는다
- SLI
  - 예: ready latency, ready 도달 여부, 실패 이벤트 수(측정 가능한 상태 기반)
- SLO
  - 예: ready time p95 <= X, fail_events == 0 등
- 왜 이게 중요한가
  - Ginkgo가 cluster 관측 기반 scenario를 테스트로 통합할 수 있는지 본다
- k6/CL2 대비 차이
  - k6는 외부 관측(HTTP)을 잘하고, CL2는 내부 부하를 잘함. Ginkgo는 “완전 커스텀 관측 로직” 축을 확인한다

## 후보 C: k6/CL2 결과 파일을 입력으로 합성 PASS/FAIL
- 시나리오
  - 입력: `k6-summary.json` 또는 CL2 산출 JSON
  - 작업: SLI 추출 → SLO 비교 → overall PASS/FAIL 단언
- SLI
  - k6/CL2가 산출한 기존 SLI를 그대로 매핑(정규화 단계 포함)
- SLO
  - SLO는 후보 입력에 포함(또는 공통 config에서 읽음)
- 왜 이게 중요한가
  - “도구 조합”을 프레임워크 레이어에서 할 수 있는지 확인(플랫폼 기반성 판단 핵심)

