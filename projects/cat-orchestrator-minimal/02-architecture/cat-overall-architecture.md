# CAT Orchestrator 최소 전체 구조

**Date:** 2026-03-19

## 1) CAT 구성요소
- **Job Spec**
  - 테스트 입력(시나리오/target/load/selected_sli/SLO/아티팩트 경로)을 “도구 공통 계약”으로 정의
- **Runner**
  - Job Spec을 읽고, `tool`을 선택해 adapter를 실행하며 raw 결과를 모은다
- **Adapter**
  - 도구별 raw 결과를 표준 `cat-result.json`으로 변환한다(로직 추가 금지, 변환만)
- **Result Store**
  - raw 결과(`k6-summary.json`, `ginkgo` artifacts, `CL2` measurement JSON 등)와 `cat-result.json`을 파일로 저장
- **(선택) Aggregator**
  - 여러 Job의 overall PASS/FAIL을 합성(이번 설계에서는 “가능”만 언급, 구현 범위 밖)

## 2) 전체 흐름(고정)
```
Job 정의(job.yaml)
  ↓
CAT Runner 실행
  ↓ (도구 선택)
Tool Adapter 실행
  ↓
Tool raw 결과 생성(요약/리포트/measurement)
  ↓
Adapter가 k6/Ginkgo/CL2 raw → cat-result.json 변환
  ↓
Result Store에 저장
  ↓
(옵션) 전체 PASS/FAIL 판단(Aggregator)
```

## 3) CAT 책임 vs 외부 책임
### CAT가 하는 것
- Job Spec을 읽고
- 도구를 실행하고
- adapter를 호출하고
- 결과를 지정된 디렉터리에 저장하도록 파일 규약을 지킨다

### CAT가 하지 않는 것
- 도구 내부 로직(예: k6 metric 집계, CL2 측정기/SLI 생성, Ginkgo SLI 계산)을 CAT가 알지 않는다
- CAT가 PASS/FAIL을 재판정하지 않는다(표준으로 “tool exit code”를 권위로 고정)

