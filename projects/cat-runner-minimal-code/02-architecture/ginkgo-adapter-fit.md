# Ginkgo Adapter Fit (범용성 검증용)

이 문서는 Ginkgo가 “k6처럼” 동작하는 도구가 아니라, Ginkgo에 자연스러운 raw/result 경로를 먼저 선택하고 그 선택이 CAT Runner 공통 계약과 어떻게 맞물리는지 정리한다.

## 1. Ginkgo raw/result 특성 정리

### Ginkgo는 raw 결과를 어떻게 남길 수 있는가
- stdout / go test 출력 텍스트
  - 예: 실패한 Expect 메시지, 테스트 실행 요약
  - 단점: CAT이 요구하는 “필요 값만 안정적으로 추출”하기 어렵고 파싱 규칙이 brittle해지기 쉽다.
- go test exit code
  - 실패하면 exit code가 1이 된다.
  - CAT PASS/FAIL 권위(재판정 금지) 구현에 적합하다.
- custom JSON(테스트 코드가 직접 파일로 생성)
  - Ginkgo 테스트 내부에서 `DeferCleanup` 등을 사용해 raw JSON 파일을 항상 생성할 수 있다.
  - CAT adapter가 “파일 우선”으로 안정적으로 파싱하기에 가장 단순하다.

### CAT 입력으로 삼는 것이 자연스러운 경로(선택)
가장 단순하고 안정적인 선택은 **커스텀 raw JSON 파일**이다.
- stdout 텍스트 파싱은 fallback으로만 고려하고,
- 기본 경로는 “테스트가 파일로 남긴 결과”를 사용한다.

### 어떤 방식이 가장 단순하고 안정적인가
- Ginkgo 테스트가 raw JSON을 생성한다.
- runner의 Ginkgo adapter는 그 raw JSON을 읽어 필요한 값만 뽑는다.
- pass/fail은 exit code로만 결정한다.

## 2. CAT Runner 공통 계약과의 적합성

CAT Runner 공통 흐름(계약)은 아래와 같다.
- `run(job)` → tool 실행 및 raw 산출물 생성
- `locate_raw_result(job)` → raw 결과 위치 반환
- `parse_raw_result(raw)` → 표준 SLI 추출
- `build_cat_result(parsed, exit_code)` → 표준 `cat-result.json` 생성

Ginkgo에 적용하면 아래 매핑이 자연스럽다.
- `run(job)`
  - runner는 `go test ./ginkgo_cat_job -run TestCATGinkgo` 형태로 Ginkgo suite를 실행한다.
  - SLO와 시나리오 주입 값은 env로 전달한다.
  - 테스트는 SLO 단언으로 exit code를 만들고, raw JSON 파일을 생성한다.
- `locate_raw_result(job)`
  - raw 파일을 `job.output.dir` 아래의 고정 경로(예: `ginkgo-raw.json`)로 찾는다.
  - Ginkgo가 stdout를 “CAT이 파싱해야 하는” 부담을 지지 않도록 파일을 표준 경로로 둔다.
- `parse_raw_result(raw)`
  - raw JSON에서 p95 latency, error rate, throughput 같은 “선택 SLI”만 읽는다.
  - CAT 표준 스키마로 정규화한다(표준 `selected_sli` + `slo_result`).
- `build_cat_result(...)`
  - `final_pass_fail`은 tool exit code를 그대로 `PASS/FAIL`로 기록한다.
  - CAT은 재판정하지 않는다.

## 3. Ginkgo adapter에서 책임/비책임 경계

### adapter가 책임지는 것
- Ginkgo suite를 실제로 실행한다(`go test` 호출).
- raw 결과 파일 위치를 찾는다.
- raw JSON을 읽고 표준 SLI/SLO evidence 형태로 정규화한다.
- exit code를 그대로 `cat-result.json.final_pass_fail`과 `exit_code`에 반영한다.

### adapter가 맡지 않는 것
- tool 내부 SLO 판정을 뒤집거나 재판정하지 않는다.
- stdout 텍스트를 의미적으로 파싱해 SLO를 대신 판단하지 않는다.
- Evidence 시각화를 위한 추가 변환(대신 파생 결과 export는 “선택적 단계”로 분리)한다.

