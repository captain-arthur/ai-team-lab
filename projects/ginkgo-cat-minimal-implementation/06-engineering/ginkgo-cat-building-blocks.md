# CAT 관점 Ginkgo 빌딩 블록(최소)

목표는 “Scenario/SLI/SLO/Result persistence”를 Ginkgo 코드 구조로 재사용 가능하게 만드는 것이다.

## 1) Scenario Injection을 코드화하는 법
### 권장 구조
- `scenarioParamsFromEnv()` 같은 함수로 입력을 먼저 만든다.
- scenario는 Go 코드에서 자유롭게 조립(HTTP handler, 타이머, 상태 머신).

### 스켈레톤
```go
type ScenarioParams struct {
  DelayMS   int
  FailEvery int
  Requests  int
}

func scenarioParamsFromEnv() ScenarioParams {
  return ScenarioParams{
    DelayMS:   envInt("SCENARIO_DELAY_MS", 50),
    FailEvery: envInt("SCENARIO_FAIL_EVERY", 999999),
    Requests:  envInt("SCENARIO_REQUESTS", 20),
  }
}
```

## 2) SLI Measurement를 모듈화하는 법
### 권장 구조
- `measureSLI(params) -> (selectedSLI, errors)`처럼 출력만 반환한다.
- measurement 함수는 assertion을 몰라야 한다(순수 계산).

### 예제(평균 latency + error rate)
```go
type Sample struct {
  latency time.Duration
  ok bool
}

func measureSLI(url string, params ScenarioParams) (avgMs float64, errorRate float64) {
  samples := make([]Sample, 0, params.Requests)
  // 반복 요청 + samples 수집
  // avg / errorRate 계산 후 반환
  return avgMs, errorRate
}
```

## 3) SLO Assertion 구조화
### 핵심 규칙
- SLO는 “boolean 게이트”로 계산하고
- 그 boolean이 `final_pass_fail`과 1:1로 연결되게 한다.

### 예제
```go
latOk := avgMs <= sloLatencyMaxMs
errOk := errorRate <= sloErrorRateMax

finalPassFail := "FAIL"
exitCode := 1
if latOk && errOk {
  finalPassFail = "PASS"
  exitCode = 0
}

Expect(finalPassFail).To(Equal("PASS"), "latOk=%v errOk=%v", latOk, errOk)
```

## 4) Result Persistence 공통화
### 권장 방식(증거 보존)
- `WriteCatResult()`는 Expect 전에 “한 번” 수행한다.
- 마지막에 `defer`로 재시도/안전장치를 둔다.

### 예제
```go
resultPath := filepath.Join("results", "job-x", "cat-result.json")
_ = catutil.WriteCatResult(resultPath, cat)
defer func() { _ = catutil.WriteCatResult(resultPath, cat) }()
```

## CAT용 결론(실전 요약)
- `Injection(입력 조립)` → `Measurement(순수 계산)` → `Evaluation(boolean gate)` → `Persistence(JSON 저장)`  
이 순서가 코드 구조로 고정되면, custom CAT 시나리오가 “프레임워크 제약”이 아니라 “Go 코드 조립”으로 확장된다.

