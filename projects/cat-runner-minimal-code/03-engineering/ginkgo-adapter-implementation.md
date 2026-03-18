# Ginkgo Adapter 구현 (실행 증명 포함)

이 문서는 “k6 전용이 아닌가?”를 증명하기 위해, runner 흐름 안에서 Ginkgo custom CAT job을 실제로 실행하고 raw 결과 → `cat-result.json` 정규화까지 end-to-end로 확인한 과정을 정리한다.

## 1. 구현 위치(요약)
- runner의 Ginkgo adapter: `07-engineering/code/ginkgo_adapter_stub.go` (stub에서 실제 구현으로 전환됨)
- Ginkgo custom CAT job(단 1개):
  - `07-engineering/code/ginkgo_cat_job/cat_ginkgo_job_test.go`
- sample job spec:
  - `07-engineering/code/sample-ginkgo-job.yaml`

## 2. runner에서 호출 가능한 형태(어댑터 contract 준수)
runner는 공통 인터페이스로 adapter를 호출한다.
- `run(job)`
  - adapter는 `go test ./ginkgo_cat_job -run TestCATGinkgo -count=1`을 실행한다.
  - scenario injection / SLO 값은 env로 전달한다.
  - tool exit code가 PASS/FAIL 권위가 된다.
- `locate_raw_result(job)`
  - adapter는 `job.output.dir` 아래 `ginkgo-raw.json`을 raw로 간주한다.
  - 주의: `go test`는 패키지 디렉토리를 working dir로 쓰기 때문에, adapter는 output dir을 절대 경로로 env 주입한다.
- `parse_raw_result(raw)`
  - raw JSON에서 `latency_p95_ms`, `error_rate`, `throughput_rps`를 읽어 표준 `selected_sli`로 만든다.
  - 측정값과 SLO 임계값의 `ok` 비교는 evidence로만 기록한다(`slo_result`).
- `build_cat_result(...)`
  - `final_pass_fail`은 tool exit code를 기준으로 기록한다.

## 3. 실제 end-to-end 실행 방법

### job.yaml 예시(샘플)
`07-engineering/code/sample-ginkgo-job.yaml`이 사용된다.
핵심은 `tool: ginkgo`이고, SLO 및 시나리오 주입 값은 `scenario.config`로 전달한다.

### 실행 명령어
```bash
cd /Users/hooni/Documents/github/ai-team-lab/projects/cat-runner-minimal-code/07-engineering/code
go run . ./sample-ginkgo-job.yaml
```

### 결과 파일 위치
- `job.output.dir` 아래 `cat-result.json` 생성
- 예:
  - `07-engineering/code/sample-run/results/ginkgo-custom-http/cat-result.json`
- 또한 raw:
  - `07-engineering/code/sample-run/results/ginkgo-custom-http/ginkgo-raw.json`

## 4. 실제 생성 결과 예시

### raw result 예시(`ginkgo-raw.json`)
```json
{
  "delay_ms": 10,
  "error_rate": 0,
  "fail_every": 999999,
  "latency_p95_ms": 12.199,
  "requests": 25,
  "throughput_rps": 85.09360083357691
}
```

### 표준 결과(`cat-result.json`) 예시
```json
{
  "test_name": "ginkgo-custom-http",
  "tool": "ginkgo",
  "scenario_type": "http",
  "selected_sli": {
    "error_rate": 0,
    "latency_p95_ms": 12.199,
    "throughput_rps": 85.09360083357691
  },
  "slo_result": {
    "error_rate": {
      "measured": 0,
      "ok": true,
      "slo_max": 0.01
    },
    "latency_p95_ms": {
      "measured": 12.199,
      "ok": true,
      "slo_max": 50
    }
  },
  "final_pass_fail": "PASS",
  "exit_code": 0,
  "raw_result": {
    "format": "json",
    "path": "/Users/hooni/Documents/github/ai-team-lab/projects/cat-runner-minimal-code/07-engineering/code/sample-run/results/ginkgo-custom-http/ginkgo-raw.json"
  },
  "timestamp": "2026-03-18T17:01:23.445996Z"
}
```

## 5. 관찰(범용성 관점)
k6에 이어 Ginkgo도 같은 runner 흐름으로 들어온다.
- CAT은 tool별 raw를 알 필요가 없다.
- adapter가 raw 포맷 차이를 흡수한다.
- tool exit code만이 `final_pass_fail` 권위를 가진다.

