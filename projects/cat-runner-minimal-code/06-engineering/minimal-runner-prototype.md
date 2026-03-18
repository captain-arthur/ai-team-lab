# CAT Runner 최소 구현 프로토타입

이 문서는 “raw 포맷 차이를 adapter가 흡수하고, CAT 표준 출력은 `cat-result.json` 하나로 고정된다”를 코드 수준에서 보여주기 위한 최소 예시다.

## 1) 최소 Runner 코드(핵심 흐름)
Runner는 Job Spec을 읽고, tool을 선택한 뒤 adapter를 호출한다.

핵심 단계:
- `job.yaml` 로드
- `job.tool`로 adapter 선택
- `adapter.Run(job)`로 tool 실행(권위: tool exit code)
- `adapter.LocateRawResult(job)`로 raw 결과 위치 확보
- `adapter.ParseRawResult(job, raw)`로 표준 SLI 추출
- `adapter.BuildCatResult(job, parsed, raw, exitCode)`로 `cat-result.json` 생성/저장

실제 코드는 `07-engineering/code/main.go`에서 위 흐름을 그대로 수행한다.

## 2) k6 adapter (실제 동작 구현)
k6 adapter는 다음을 수행한다.

1. k6 실행
   - `k6 run --summary-export <rawPath> <entry>`
   - runner가 `TARGET_URL`, `MODE`, `TARGET_RPS`, `DURATION`, `SLO_P95_MS`, `SLO_FAIL_RATE`를 env로 주입
2. raw(JSON summary-export) 파싱
   - `metrics.http_req_duration["p(95)"]` -> `latency_p95_ms`
   - `metrics.http_req_failed["value"]` -> `error_rate`
   - `metrics.http_reqs["rate"]` -> `throughput_rps`
3. 표준 cat-result.json 생성
   - `final_pass_fail`은 **k6 프로세스 exit code**로만 결정한다.

실제 구현은 `07-engineering/code/k6_adapter.go`에 포함되어 있다.

## 3) Ginkgo / CL2 adapter (Ginkgo는 실제 구현, CL2는 stub)
이 최소 구현에서는 Ginkgo는 실제 raw JSON 경로를 통해 end-to-end 정규화까지 구현하고,
CL2는 구조적으로 수용 가능한 contract 수준(stub)만 둔다.

공통 인터페이스(함수 시그니처):
```go
type ToolAdapter interface {
  Run(job JobSpec) (exitCode int, err error)
  LocateRawResult(job JobSpec) (raw RawRef, err error)
  ParseRawResult(job JobSpec, raw RawRef) (parsed ParsedSLI, err error)
  BuildCatResult(job JobSpec, parsed ParsedSLI, raw RawRef, exitCode int) CatResult
}
```

Ginkgo의 raw 포맷 경로(구현 있음):
- raw format: 테스트가 생성하는 커스텀 JSON 파일(예: `ginkgo-raw.json`)
- locate/parse: raw 파일 우선(텍스트(stdout) 파싱은 기본 경로가 아니다)

CL2(clusterloader2) stub에서의 raw 포맷 가정:
  - raw format: JSON 또는 XML
  - 추후 구현 포인트: 측정되는 metric 이름과 XML/JSON parse 경로 분기

- `07-engineering/code/cl2_adapter_stub.go`

두 adapter는 동일한 인터페이스(`Run / LocateRawResult / ParseRawResult / BuildCatResult`)를 만족하지만, 현재는 `not implemented` 에러를 반환한다.

## 4) 실행 예시(k6 end-to-end)

### 준비물
- `k6` 바이너리가 로컬에 설치되어 있어야 한다.

### 입력(Job Spec 1개)
- `07-engineering/code/sample-job.yaml` 또는 프로젝트 루트의 `job.yaml`을 사용한다.

### 실행 명령어(예시)
아래 예시는 runner가 있는 디렉토리에서 실행한다.

```bash
cd projects/cat-runner-minimal-code/07-engineering/code
go run . ../../job.yaml
```

### 생성 결과
`job.output.dir` 아래에 `cat-result.json`이 생성된다. (예: `./results/http-basic-run/cat-result.json`)

### cat-result.json 예시
```json
{
  "test_name": "http-basic",
  "tool": "k6",
  "scenario_type": "http",
  "selected_sli": {
    "error_rate": 0,
    "latency_p95_ms": 222.52769999999998,
    "throughput_rps": 78.83383758268516
  },
  "slo_result": {
    "error_rate": {
      "measured": 0,
      "ok": true,
      "slo_max": 0.01
    },
    "latency_p95_ms": {
      "measured": 222.52769999999998,
      "ok": true,
      "slo_max": 450
    }
  },
  "final_pass_fail": "PASS",
  "exit_code": 0,
  "raw_result": {
    "format": "json",
    "path": "../../results/http-basic-run/k6-summary.json"
  },
  "timestamp": "2026-03-18T16:49:55.38079Z"
}
```

