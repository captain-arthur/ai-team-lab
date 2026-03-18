# Custom CAT Scenario 01: 안정 구간 유지(stability window)

## 무엇을 검증하는가
- 어떤 상태가 “나중에 ready가 되고” 이후 안정적으로 유지되는지 검증한다.

## Go 코드로 어떻게 custom 되는가
- 서버의 ready 여부를 `atomic.Bool` 상태로 직접 제어한다.
- test는
  - 최초 ready 응답까지의 latency(=SLI 1)
  - ready 관측 구간의 error rate(=SLI 2)
  를 계산해 SLO를 게이트한다.

## 실행
```bash
cd projects/ginkgo-cat-minimal-implementation/08-engineering/custom-cat-scenarios/01-stability-window
go test -v ./...
```

## 시나리오 입력(환경변수, 선택)
- `WARMUP_MS` (기본 120)
- `STABILITY_MS` (기본 400)
- `REQUESTS_IN_WINDOW` (기본 10)
- `SLO_READY_MAX_MS` (기본 250)
- `SLO_ERROR_RATE_MAX` (기본 0.0)

## 결과 파일
- `results/stability-window/cat-result.json`

## CAT 최소 조건 매핑
- Scenario Injection: WARMUP/requests 입력
- SLI Measurement: ready latency + stability 구간 error rate
- SLO Assertion: ready_ok && error_ok
- Result Persistence: cat-result.json 저장
