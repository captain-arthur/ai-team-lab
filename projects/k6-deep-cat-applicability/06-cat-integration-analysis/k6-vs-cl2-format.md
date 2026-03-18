# k6 vs clusterloader2: 결과 구조 비교와 CAT 표준 포맷 융합 판단

**Date:** 2026-03-18  
**근거(k6)**: `04-engineering/results/golden-arrival-summary.json`  
**근거(CL2/CAT)**: `projects/sample-cl2-sli-slo-analysis/08-cat-v1-spec/cat-v1-specification.md`

## 1) k6 vs CL2 결과 구조(요약 비교)
| 항목 | k6 | clusterloader2(CL2) + CAT v1 권장 |
|---|---|---|
| 실행 단위 | `k6 run` 1회 | CL2 run 1회 |
| raw 산출 | `--summary-export` JSON(지표 집계) + 콘솔 로그 | 측정 JSON 다수(측정기별 파일) + 로그 |
| SLI | 외부 관측 SLI(요청 지연/실패율/처리량 등) | 내부 관측 SLI(OOM, 시스템 파드 재시작, pod startup latency 등) |
| SLO 평가 | k6 threshold(실행 중/종료 시) + 종료 코드 | CAT 로직이 측정 JSON을 읽어 SLI별 판정 + overall 산출 |
| 최종 결과 저장(권장) | 별도 표준 파일 필요 | CAT v1에서 `cat-result`/요약 파일 생성을 권장 |

## 2) 공통 요소(융합 포인트)
둘 다 CAT 표준 포맷으로 수렴시킬 때 필요한 공통 요소는 동일하다.
- **메타데이터**: test_name, tool, scenario_type, target(or cluster), timestamp
- **SLI 측정값**: selected_sli(키-값)
- **SLO 평가 결과**: SLI별 PASS/FAIL(또는 최소 overall PASS/FAIL 근거)
- **최종 판정**: PASS / FAIL (+ 정책에 따라 PASS_WITH_WARNINGS)
- **증거(artifacts)**: raw 결과 파일 경로

## 3) 차이점(표준화 시 주의)
- k6는 “외부 SLI”가 raw summary에 한 파일로 뭉쳐 나오지만, CL2는 측정기별 JSON이 다수로 흩어진다.
- k6의 PASS/FAIL은 threshold+종료 코드로 즉시 나오고, CL2는 “측정→비교→판정”이 별도 단계다.
- 따라서 CAT 표준화는:
  - k6: raw summary + (래퍼가) 표준 결과 생성
  - CL2: 측정 JSON들 + (평가기) 표준 결과 생성
  로 접근이 다르다.

## 4) CAT 표준 포맷으로 융화 가능한가?
👉 **YES.**
- 이유: CAT v1 명세 자체가 “SLI 측정값 요약 + SLI별 판정 + overall 결과”를 파일로 남기는 형태를 권장하며, k6도 동일하게 (1) raw summary 보관, (2) 표준 결과 파일 생성으로 맞출 수 있다.
- 핵심은 도구별 raw를 동일하게 만들려는 게 아니라, **도구 공통 결과 파일(`cat-result.json`)**로 수렴시키는 것이다.

## 5) evidance 같은 도구로 시각화 가능한가?
👉 **가능(조건부)**.
- 조건: `cat-result.json`이 테스트 단위로 누적 저장되고, selected_sli와 final_pass_fail이 일관된 키로 정규화되어 있어야 한다.
- k6 raw summary는 시각화 도구에 직접 넣기엔 필드가 많고 변동 여지가 있어, **표준 결과 파일을 1차 시각화 소스로** 두는 게 안정적이다.
