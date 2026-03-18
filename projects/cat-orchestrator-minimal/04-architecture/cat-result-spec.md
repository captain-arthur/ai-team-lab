# CAT 표준 결과 스키마(cat-result.json) (고정)

```json
{
  "test_name": "",
  "tool": "",
  "scenario_type": "",
  "selected_sli": {},
  "slo_result": {},
  "final_pass_fail": "",
  "exit_code": 0,
  "raw_result_path": "",
  "timestamp": ""
}
```

## 1) SLI 표준화 전략(최소)
- `selected_sli`는 “CAT이 알아야 하는 SLI 키만” 들어간다.
- 각 tool adapter는 자기 raw에서 아래 형태로 정규화한다.
  - latency: `avg_latency_ms` 또는 `latency_p95_ms`
  - error rate: `error_rate`
  - throughput: `throughput_rps`
- metric 이름(`slo[].metric`)은 `selected_sli` 키와 1:1로 매핑된다.

## 2) PASS/FAIL 권위(중요, 단정)
- **권위 = tool exit code**
  - tool이 종료 코드 0이면 `final_pass_fail="PASS"`
  - 0이 아니면 `final_pass_fail="FAIL"`
- CAT는 재판정하지 않는다.
- `slo_result`는 기록용(왜 PASS/FAIL인지 근거)을 제공한다.

## 3) k6/Ginkgo/CL2 매핑 방식(요약)
- **k6**
  - raw: `k6-summary.json`
  - adapter가 p95/failed rate를 뽑아 `selected_sli` 채움
  - final은 k6 exit code로 결정(중복 재판정 금지)
- **Ginkgo**
  - raw: 테스트 코드가 이미 cat-result.json을 기록(또는 adapter가 파일을 생성)
  - final은 `go test` exit code로 결정
- **CL2**
  - raw: 측정 JSON들(ClusterLoader2 결과 디렉터리)
  - adapter가 필요한 SLI만 추출해 `selected_sli` 채움
  - final은 CL2 실행 exit code로 결정(또는 assertion이 들어있는 runner exit code)

