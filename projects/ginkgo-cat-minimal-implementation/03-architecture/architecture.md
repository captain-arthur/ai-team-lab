# Ginkgo custom CAT Job 흐름(최소)

## 1) 실행 흐름(요구 흐름 고정)
Scenario 실행
→ SLI 수집
→ SLO 단언
→ PASS/FAIL 결정
→ cat-result.json 저장

## 2) PASS/FAIL 권위
- Ginkgo suite의 테스트 성공/실패가 go test의 종료 코드로 반영된다.
- CAT는 재판정하지 않고, **테스트 실패 = FAIL**로만 기록한다.
- `cat-result.json`은 기록용(근거 데이터 포함)으로 남긴다.

## 3) SLI 수집 위치(명확화)
- SLI는 코드 내부에서 측정한다.
  - `httptest` 서버에 HTTP 요청을 N회 보내고
  - 각 요청의 latency를 직접 time measurement으로 수집한다.
  - 상태코드/에러로 error 여부를 계산해 error rate를 만든다.
- 외부 조회(Prometheus/k6 등)는 이번 최소 구현 범위에서 사용하지 않는다.

