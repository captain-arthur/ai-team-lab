# Engineering: k6 통합 설계(실행 단위/판정/저장)

**Date:** 2026-03-18  
**목표:** CAT 최소요구 4요소를 k6로 “실제로 구현 가능”하게 만드는 최소 설계

---

## 1) 최소 실행 단위 정의(k6 기반 CAT test/job)
### 입력(필수)
- **test_name**: CAT에서의 테스트 이름(고유)
- **tool**: `k6` (고정)
- **script_path**: 실행할 k6 스크립트 경로
- **scenario_type**: 예) `constant-vus` / `constant-arrival-rate` / `ramping-arrival-rate`
- **target**: 예) `https://<ingress-host>/echo`
- **env(params)**: 스크립트가 소비하는 환경변수 맵(예: `TARGET_URL`, `MODE`, `THINK_TIME_S`, `SLO_P95_MS`)
- **selected_sli**: 저장할 SLI 키 목록(최소: p95 latency, error rate, (선택) throughput)
- **output_dir**: 결과를 저장할 디렉터리

### 실행(표준 커맨드)
- CAT Runner는 다음 형태로 k6를 실행한다.

```bash
mkdir -p "<output_dir>" && \
<ENV_KV_PAIRS> k6 run \
  --summary-export "<output_dir>/k6-summary.json" \
  "<script_path>" \
  | tee "<output_dir>/k6-output.txt"
```

### 출력(필수 파일)
- `k6-summary.json`: k6 raw summary (`--summary-export`)
- `k6-output.txt`: 콘솔 출력(사람/디버그용)
- `cat-result.json`: CAT 표준 결과(아래 스키마)

---

## 2) 결과 저장 설계(최소 스키마)
### CAT 표준 결과 파일: `cat-result.json` (최소)
반드시 포함 필드(요구사항 고정):
- `test_name`
- `tool`
- `scenario_type`
- `target`
- `selected_sli`
- `slo_result`
- `final_pass_fail`
- `timestamp`

제안 스키마(POC 수준, 도구 공통):
```json
{
  "test_name": "ingress-steady-a",
  "tool": "k6",
  "scenario_type": "constant-arrival-rate",
  "target": "https://<ingress-host>/echo",
  "timestamp": "2026-03-18T13:00:00Z",
  "selected_sli": {
    "http_req_duration_p95_ms": 243.9,
    "http_req_failed_rate": 0.0000,
    "http_reqs_rate_rps": 43.5
  },
  "slo_result": {
    "source": "k6_thresholds",
    "passed": true,
    "failed_thresholds": []
  },
  "final_pass_fail": "PASS",
  "artifacts": {
    "k6_summary_path": "k6-summary.json",
    "k6_output_path": "k6-output.txt"
  },
  "notes": {
    "prometheus": "N/A"
  }
}
```

원칙:
- `k6-summary.json`은 **raw**(정밀 분석/추후 파싱용)
- `cat-result.json`은 **판정용 최소 요약**(비교/누적/시각화용)

---

## 3) PASS/FAIL 설계(모호함 제거)
### 판정의 권위(우선순위)
1. **k6 종료 코드**(권위): 0=PASS, 비0=FAIL  
2. **cat-result.json**(기록): 어떤 SLO(=threshold)가 깨졌는지/측정값이 무엇인지

### 관계 정의
- CAT Runner는 k6 실행 종료 후:
  - `exit_code==0`이면 `final_pass_fail=PASS`
  - `exit_code!=0`이면 `final_pass_fail=FAIL`
- `slo_result.passed`는 `final_pass_fail`과 **항상 일치**해야 한다(불일치 금지).

### threshold 근거 기록
- 최소로는 “실패한 threshold 목록”을 기록한다.
- 실패한 threshold의 상세(예: `p(95)<300` vs 실제 값)는 k6 콘솔/summary에서 파싱해 `failed_thresholds`에 넣는다(구현 단계에서 래퍼가 수행).

---

## 4) Prometheus 활용 방식(보조 진단 vs 판정 승격)
### 기본 정책(권장)
- **Prometheus는 기본적으로 “진단”**이다. PASS/FAIL의 1차 권위는 k6에 둔다.
- 즉, `final_pass_fail`을 “Grafana에서 보고 바꾼다”는 흐름은 금지.

### 승격 조건(판정 입력으로 쓰는 경우)
다음 중 하나면 Prometheus 평가를 **2차 SLO 게이트**로 승격할 수 있다.
- 합격 기준이 “내부 SLI 포함”인 테스트(예: ingress가 느려지면 실패 + 동시에 노드 CPU 포화도 실패)
- 외부 SLI만으로는 “수용 불가” 판단이 과도하게 흔들리는 환경(정책으로 내부 조건을 함께 요구)

승격 시에도 원칙은 동일:
- **최종 결과는 파일(`cat-result.json`)로 저장**
- `slo_result.source`에 `k6_thresholds + prom_slo`처럼 합성 출처를 명시

---

## 5) 예시: 최소 CAT job 정의(ingress HTTP k6 테스트)
### 예시 Job(개념 YAML)
```yaml
test_name: ingress-steady-a
tool: k6
script_path: scripts/ingress-poc.js
scenario_type: constant-arrival-rate
target: https://<ingress-host>/echo
env:
  TARGET_URL: https://<ingress-host>/echo
  CASE: A
  TARGET_RPS: "200"
  SLO_P95_MS: "300"
  SLO_FAIL_RATE: "0.001"
selected_sli:
  - http_req_duration_p95_ms
  - http_req_failed_rate
  - http_reqs_rate_rps
output_dir: results/2026-03-18/ingress-steady-a
```

### 실행→결과 파일 생성 흐름(요약)
1. CAT Runner가 Job을 읽음(입력)
2. env를 주입해 `k6 run --summary-export ...` 실행(실행/측정)
3. 종료 코드로 PASS/FAIL 결정(판정)
4. `k6-summary.json` + `cat-result.json` 저장(저장)
5. FAIL이면 Prometheus는 “원인 분류”로만 추가 메모(기본 정책)
