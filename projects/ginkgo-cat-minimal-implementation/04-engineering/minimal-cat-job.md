# 04-Engineering: Ginkgo 기반 minimal custom CAT Job(1개)

## 1) 구현 파일
- 테스트 코드: `ginkgo_cat_test.go`
- 결과 스키마: `cat_result.go`
- 실행 결과 파일: `results/cat-result.json`

## 2) custom 시나리오(Go 코드 내부)
- `httptest` 서버를 띄운 뒤 `/`에 HTTP 요청을 `SCENARIO_REQUESTS` 회 보냄
- 서버는
  - 응답 지연: `SCENARIO_DELAY_MS` ms
  - 오류 주입: 요청 카운터가 `SCENARIO_FAIL_EVERY`의 배수이면 500 반환

### 입력(환경변수)
- `SCENARIO_DELAY_MS` (기본 50)
- `SCENARIO_FAIL_EVERY` (기본 999999: 실패 없음)
- `SCENARIO_REQUESTS` (기본 20)
- `SLO_LATENCY_MAX_MS` (기본 300)
- `SLO_ERROR_RATE_MAX` (기본 0.0)

## 3) SLI 측정(코드 내부 계산, 2개)
- `avg_latency_ms`: N회 요청 latency 평균(밀리초)
- `error_rate`: (성공 조건 미충족 횟수 / N)
  - 성공 조건: `err==nil` 이고 `HTTP 200`인 경우만 OK

## 4) SLO 단언(PASS/FAIL 명확화)
- `latency_ok`: `avg_latency_ms <= SLO_LATENCY_MAX_MS`
- `error_ok`: `error_rate <= SLO_ERROR_RATE_MAX`
- `final_pass_fail`:
  - 둘 다 true면 `PASS`, 하나라도 false면 `FAIL`
- 단언: `Expect(cat.FinalPassFail).To(Equal("PASS"))`
  - 이 assertion이 깨지면 테스트 FAIL로 종료됨(=exit code는 비0)

## 5) 실행 방법(실제 동작 확인)
```bash
cd projects/ginkgo-cat-minimal-implementation/04-engineering
go test -v ./...
```

## 6) 실행 로그 일부(실제 출력 발췌)
```text
=== RUN   TestCATGinkgo
Running Suite: CAT Ginkgo Minimal Suite - /Users/hooni/Documents/github/ai-team-lab/projects/ginkgo-cat-minimal-implementation/04-engineering
...
SUCCESS!  --  1 Passed | 0 Failed
--- PASS: TestCATGinkgo (1.04s)
PASS
ok  	ginkgo-cat-minimal	1.40s
```

## 7) 생성된 결과 파일(cat-result.json 예시, 실제 생성)
```json
{
  "test_name": "ginkgo-basic-test",
  "tool": "ginkgo",
  "scenario_type": "custom-http",
  "selected_sli": {
    "avg_latency_ms": 51.25,
    "error_rate": 0
  },
  "slo_result": {
    "latency_ok": true,
    "error_ok": true,
    "latency_slo_ms": 300,
    "error_slo_rate": 0
  },
  "final_pass_fail": "PASS",
  "exit_code": 0,
  "timestamp": "2026-03-18T15:25:13.540854Z",
  "scenario_params": {
    "delay_ms": 50,
    "fail_every": 999999,
    "requests": 20
  }
}
```

