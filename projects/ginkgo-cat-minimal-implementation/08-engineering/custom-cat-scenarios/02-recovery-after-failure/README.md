# Custom CAT Scenario 02: 실패 후 회복(recovery after failure)

## 무엇을 검증하는가
- 일정 시간 동안 실패 상태가 지속된 뒤, 시스템이 회복하는지를
  - 전체 error rate
  - 회복 시간(연속 성공 N회 도달 시점)
  으로 평가한다.

## Go 코드로 어떻게 custom 되는가
- HTTP handler가 `failUntil` 시점을 기준으로 500/200을 전환한다.
- recovery_time은 “첫 실패 이후 연속 성공 N회”라는 정의로 Go 코드에서 직접 계산한다.

## 실행
```bash
cd projects/ginkgo-cat-minimal-implementation/08-engineering/custom-cat-scenarios/02-recovery-after-failure
go test -v ./...
```

## 시나리오 입력(환경변수, 선택)
- `FAIL_FOR_MS` (기본 300)
- `TOTAL_MS` (기본 1200)
- `INTERVAL_MS` (기본 50)
- `SUCCESS_STREAK` (기본 3)
- `SLO_MAX_ERROR_RATE` (기본 0.35)
- `SLO_MAX_RECOVERY_MS` (기본 800)

## 결과 파일
- `results/recovery-after-failure/cat-result.json`

## CAT 최소 조건 매핑
- Scenario Injection: failUntil/요청 주기 입력
- SLI Measurement: error rate + recovery time
- SLO Assertion: error_ok && recovery_ok
- Result Persistence: cat-result.json 저장
