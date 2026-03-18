# Final Report: k6 CAT 표준 실행 구조(Adapter 설계)

**Date:** 2026-03-18

## 1) k6의 CAT 내 역할
- k6는 CAT 밖에서 실행되는 **외부 실행 도구**이며, HTTP류 “외부 관측” 테스트의
  - 시나리오 실행(부하 발생)
  - SLI 산출(요약 metric)
  - SLO 단언(threshold→exit code)
  를 담당한다.

## 2) 실행 → 결과 수집 방식(요약)
- CAT Job 입력(표준 필드)
  - `test_name`, `scenario(target/rps/duration)`, `slo(p95/error_rate)`, `selected_sli`, `artifacts_dir`
- CAT Runner 실행
  - `k6 run --summary-export <dir>/k6-summary.json <script> | tee <dir>/k6-output.txt`
  - 종료 코드 수집(권위)
- Adapter 변환
  - `k6-summary.json`에서 selected SLI를 뽑아 `cat-result.json` 생성
  - `final_pass_fail`은 종료 코드와 동일

## 3) 확장성(다른 도구에도 적용 가능한가?)
- **가능(YES)**: CAT 표준은 “도구 실행/원시 결과/정규화 결과” 3단만 고정하면 된다.
  - fio 등도 동일하게 `fio raw` → `adapter` → `cat-result.json`로 편입 가능.

## 4) 최종 결론
👉 **“k6는 CAT에서 어떤 형태(역할/인터페이스)로 존재해야 하는가?”**

- k6는 CAT의 “내부 모듈”이 아니라, **외부 도구 플러그인**으로 존재해야 한다.
- CAT는 k6를 **표준 Job 스펙으로 실행**하고,
  - raw(`k6-summary.json`)
  - 표준(`cat-result.json`)
  을 파일로 저장해 누적/비교/시각화 가능하게 해야 한다.
- PASS/FAIL은 k6 threshold 결과(=exit code)로 **단언**하고, CAT는 **재판정하지 않는다**.
