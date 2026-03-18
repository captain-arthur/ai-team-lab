# Final Report: k6를 CAT 시스템에 통합하는 설계

**Date:** 2026-03-18  
**Program:** CAT

## 1) 최종 설계 요약
- k6는 CAT에서 **외부 트래픽 시나리오 실행 + 외부 SLI 측정 + k6 threshold 기반 1차 SLO 게이트**를 담당한다.
- CAT는 k6의 raw 결과(`k6-summary.json`)를 보존하면서, 도구 공통 판정 파일(`cat-result.json`)로 **최종 PASS/FAIL을 파일로 고정 저장**한다.

## 2) k6가 CAT에서 맡는 역할
- **역할**: 외부 관측 기반 테스트의 실행 엔진(Scenario) + 외부 SLI 산출 + SLO(PASS/FAIL) 단언
- **경계**: 내부 SLI가 합격의 핵심인 테스트에서는 “부하 발생”까지만 담당(판정 주체 아님)

## 3) CAT 최소 요구 4요소 충족 방식(결론형)
- **시나리오 주입**: CAT Job 정의 → env/옵션으로 k6에 주입
- **SLI 측정**: k6 summary에서 selected_sli만 추출
- **SLO 평가**: k6 threshold + 종료 코드로 PASS/FAIL 단언(종료 코드가 권위)
- **결과 저장**: `k6-summary.json`(raw) + `cat-result.json`(표준 판정) 2파일 저장

## 4) k6 통합 시 기대 효과
- “외부 SLI 기반 수용” 테스트를 **도구 내부에서 바로 PASS/FAIL**로 게이팅 가능(자동화 친화).
- 결과가 파일(`cat-result.json`)로 남아 **비교/누적/시각화**가 가능해진다.

## 5) k6 통합의 한계
- 내부 SLI 기반 수용(제어면/노드/리소스)은 k6 단독으로 판정하기 어렵다.
- 외부 SLI 악화의 원인 분류는 k6가 아니라 내부 신호(예: Prometheus)가 필요해질 수 있다(단, 결과 저장의 대체재가 아님).

## 6) CL2와의 관계
- **k6는 CL2를 대체하지 않는다.** k6는 외부 SLI, CL2는 내부 SLI에 강하다.
- CAT는 둘의 결과를 동일한 `cat-result.json` 구조로 저장해 도구가 달라도 비교/누적한다.

## 7) 최종 판단: “k6는 CAT에 어떻게 통합되어야 하는가?”
- k6는 CAT에서 **외부 트래픽 기반 테스트 도구 슬롯**으로 통합한다.
- 실행 단위는 “k6 Job 1개 = k6 run 1회 + 결과 파일 2개 생성”으로 고정한다.
- PASS/FAIL은 **k6 종료 코드**로 단언하고, `cat-result.json`에 selected_sli와 실패 threshold 근거를 함께 저장한다.
