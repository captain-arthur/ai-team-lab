# Research Notes: k6가 CAT 최소요구 4요소를 어디까지 충족하는가

**Date:** 2026-03-18

## 결론 요약
- k6는 CAT 최소요구 4요소를 **기본적으로 모두 충족**할 수 있다.  
  단, “결과 저장”은 k6의 `--summary-export`만으로도 가능하지만, CAT가 원하는 최소 스키마(도구 공통 결과)는 **어댑터(래퍼)**로 표준화하는 게 안전하다.

## 1) 시나리오 주입(Scenario Injection)
- **충족(YES)**: 스크립트 + 환경변수(`__ENV.*`) + 실행 옵션으로 타깃/부하/모드 주입 가능.
- **CAT에서 필요한 보완**
  - “무엇이 시나리오 입력인지”를 CAT 레벨에서 고정(예: `target_url`, `scenario_type`, `load_model`, `duration`, `env`).

## 2) SLI 측정(SLI Measurement)
- **충족(YES)**: `http_req_duration(p95)`, `http_req_failed(rate)`, `http_reqs(rate)` 등 외부 SLI를 k6가 직접 산출.
- **CAT에서 필요한 보완**
  - CAT가 “선택된 SLI”만 뽑아 저장하도록 규칙화(너무 많은 metric 저장 방지).

## 3) SLO 평가(SLO Evaluation)
- **충족(YES)**: threshold 위반 시 FAIL + 비0 종료 코드로 “단언” 가능.
- **CAT에서 필요한 보완**
  - CAT의 최종 판정은 **종료 코드(권위) + 결과 파일(근거)**로 고정한다.
  - 내부 SLO(예: 자원 포화)까지 판정에 넣고 싶다면, k6가 아니라 **별도 평가 단계**(Prometheus 등)로 합성해야 한다.

## 4) 최종 결과 저장(Result Persistence)
- **부분 충족(기능은 YES, 표준화는 필요)**:
  - k6는 `--summary-export`로 JSON 저장 가능(“raw summary”).
  - 하지만 CAT는 도구 공통 필드(`test_name`, `final_pass_fail` 등)가 필요하므로 **CAT 결과 파일**을 별도로 생성하는 게 맞다.

## 보완 원칙(최소)
- **k6는 k6가 잘하는 것만**: 외부 트래픽 실행 + 외부 SLI + k6 threshold 기반 1차 PASS/FAIL
- **CAT 공통 결과 저장은 어댑터가 담당**: k6 raw summary + CAT 표준 결과(`cat-result.json`) 동시 저장
