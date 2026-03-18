# k6 Adapter Spec (CAT → k6 실행/수집/저장)

**Date:** 2026-03-18  
**목표:** CAT가 k6를 “표준 인터페이스”로 실행하고, 결과를 표준 파일로 수집

---

## 1) 입력 스펙(CAT → k6)
### CAT Job 입력(예시)
```yaml
test_name: http-basic
tool: k6
scenario:
  type: http
  target: https://example.com
  load_model: constant-arrival-rate
  rps: 100
  duration: 60s
slo:
  latency_p95_ms: 300
  error_rate: 0.01
selected_sli:
  - latency_p95_ms
  - error_rate
  - throughput_rps
artifacts_dir: results/2026-03-18/http-basic
script:
  kind: builtin
  name: http-basic
```

### 입력이 k6에 매핑되는 방식(규칙)
- **Scenario → k6 script**
  - `scenario.type=http` 이면 HTTP용 스크립트 사용(내장 script 또는 지정된 script_path)
  - `load_model`에 따라 executor 선택
- **Target → env**
  - `scenario.target` → `TARGET_URL`
- **Load → env**
  - `rps` → `TARGET_RPS`
  - `duration` → `DURATION`
- **SLO → thresholds(env)**
  - `latency_p95_ms` → `SLO_P95_MS`
  - `error_rate` → `SLO_FAIL_RATE`

---

## 2) 실행 방식(표준)
### 실행 명령(표준형)
```bash
mkdir -p "<artifacts_dir>" && \
TARGET_URL="<scenario.target>" \
MODE="arrival" \
TARGET_RPS="<scenario.rps>" \
DURATION="<scenario.duration>" \
SLO_P95_MS="<slo.latency_p95_ms>" \
SLO_FAIL_RATE="<slo.error_rate>" \
k6 run \
  --summary-export "<artifacts_dir>/k6-summary.json" \
  "<script_path>" \
  | tee "<artifacts_dir>/k6-output.txt"
```

### script 선택 방식(최소)
- CAT는 job에 `script.kind/name` 또는 `script_path`를 준다.
- Adapter는 `script.name=http-basic`이면 저장소 내 `03-engineering/scripts/http-basic.js`를 사용한다.

---

## 3) 결과 수집 구조(필수 2종)
### (1) k6 raw 결과(증거)
- `k6-summary.json` (`--summary-export`)
- `k6-output.txt` (콘솔 로그, 디버그/설명용)

### (2) CAT 표준 결과(소비)
- `cat-result.json` (표준 결과)

표준 결과 예시(요구 필드 포함):
```json
{
  "test_name": "http-basic",
  "tool": "k6",
  "scenario_type": "http",
  "target": "https://example.com",
  "sli": {
    "latency_p95_ms": 210,
    "error_rate": 0.001,
    "throughput_rps": 98.7
  },
  "slo": {
    "latency_p95_ms": 300,
    "error_rate": 0.01
  },
  "final_pass_fail": "PASS",
  "exit_code": 0,
  "timestamp": "2026-03-18T13:00:00Z"
}
```

---

## 4) PASS/FAIL 정의(중요, 단언 규칙)
- **기준**: k6 **exit code**
  - `0` → `final_pass_fail=PASS`
  - `!=0` → `final_pass_fail=FAIL`
- **threshold 위반 → FAIL** (k6가 이미 단언)
- **CAT는 재판정하지 않음**
  - `cat-result.json`은 “판정 결과”를 기록만 한다.

---

## 5) 변환 로직(Adapter 핵심)
### 입력
- `k6-summary.json`
- CAT Job 메타데이터(`test_name`, `scenario_type`, `target`, `slo`, timestamp)
- k6 exit code

### 출력
- `cat-result.json`

### metric 매핑(최소)
- `latency_p95_ms` ← `metrics.http_req_duration.p(95)` *(단위 ms)*
- `error_rate` ← `metrics.http_req_failed.value` *(0~1)*
- `throughput_rps` ← `metrics.http_reqs.rate`

### 저장 규칙
- raw는 그대로 보존(`k6-summary.json`)
- 표준 결과는 **선택된 SLI만** 포함(필드 폭발 방지)
