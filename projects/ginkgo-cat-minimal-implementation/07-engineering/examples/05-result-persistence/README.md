# Example 05: 결과 파일 저장(cat-result.json)

## 목적
- Expect 실패여도 증거를 남기도록 `WriteCatResult`를 “먼저” 쓰는 패턴을 익힌다.
- custom CAT 스타일의 최소 흐름을 그대로 따른다:
  - Scenario 실행
  - SLI 계산
  - SLO 게이트(PASS/FAIL)
  - cat-result.json 저장

## 실행
```bash
cd projects/ginkgo-cat-minimal-implementation/07-engineering/examples/05-result-persistence
go test -v ./...
```

## 시나리오 주입(선택)
- `DELAY_MS` 기본 20
- `FAIL_EVERY` 기본 999999(실패 없음)
- `REQUESTS` 기본 10
- `SLO_LATENCY_MAX_AVG_MS` 기본 300
- `SLO_ERROR_RATE_MAX` 기본 0.0

## 출력
- `results/result-persistence/cat-result.json`

## CAT 관점 연결
- CAT 최소 조건 중 “Result Persistence”를 예제로 완결한다.
