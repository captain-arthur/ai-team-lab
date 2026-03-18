# Evidence-ready 스키마 제안(파생 결과)

## 1. `cat-result.json`과 Evidence 시각화용 전처리 결과의 관계

- `cat-result.json`
  - CAT 표준 출력(권위): tool exit code를 기준으로 `final_pass_fail`/`exit_code`를 기록한다.
  - adapter가 raw 포맷 차이를 흡수해서 `selected_sli`와 `slo_result` evidence를 담는다.
  - CAT은 여기서 멈춘다(시각화/DB 적재 범위 밖).
- Evidence 시각화용 전처리 파일(파생 결과)
  - Evidence가 다루기 쉬운 “row 기반 metric 테이블” 형태로 `cat-result.json`을 flatten한다.
  - CAT 원본 구조를 뒤집거나 재판정하지 않고, “읽기 쉬운 형태”로만 바꾼다.

정리하면:
`raw(tool)` → adapter → `cat-result.json` → (선택적 파생 단계) → `evidence-ready(테이블/CSV/rows)`

## 2. Evidence에 필요한 최소 전처리 스키마(제안)

### 형태: row 기반 metric 테이블(가장 단순)
각 `cat-result.json.selected_sli`의 (metric_name, metric_value) 쌍을 한 row로 flatten한다.

예: `evidence-metrics.json`(또는 `evidence-metrics.csv`) 형태
```json
[
  {
    "test_name": "",
    "tool": "",
    "scenario_type": "",
    "metric_name": "",
    "metric_value": 0,
    "metric_unit": "",
    "status": "", 
    "timestamp": "",
    "raw_result_path": "",
    "tags": {}
  }
]
```

## 3. 필드 필요/불필요 판단(요구 필드 포함)

### 반드시 필요한 필드(요구사항 그대로 채택)
- `test_name`
- `tool`
- `scenario_type`
- `metric_name`
- `metric_value`
- `metric_unit`
- `status`
- `timestamp`
- `raw_result_path`
- `tags(선택)`

### 무엇을 `status`로 둘지(정의)
- Evidence 관점에서는 “해당 metric이 SLO 조건을 만족했는가”가 가장 유용하다.
- 따라서 `status`는 다음 규칙을 추천한다.
  - `cat-result.json.slo_result[metric_name].ok == true`  → `status: "ok"`
  - `cat-result.json.slo_result[metric_name].ok == false` → `status: "fail"`
  - `slo_result`에 해당 metric 정보가 없으면 → `status: "n/a"`

### 불필요한 필드(최소화)
- `cat-result.json.final_pass_fail`과 `slo_result` 전체 구조는 Evidence에 “필수”가 아니다.
  - 왜냐하면 metric별 `status`와 `metric_value`가 이미 차트를 그릴 수 있기 때문이다.
- 다만 “테스트 단위 최종 PASS/FAIL”을 Evidence에서 바로 필터링하고 싶다면, 파생 스키마 row에 추가하는 것은 선택이다(이번 최소 제안에서는 고정값으로 두지 않는다).

## 4. 왜 이 스키마가 Evidence에 적합한가

Evidence는 보통 “필터/그룹/비교”에 강점이 있고, row 기반 metric 테이블은 다음 장점을 갖는다.
- `metric_name`을 컬럼으로 둔 뒤 여러 지표를 같은 형식으로 누적 가능
- `metric_value`, `status`로 PASS/FAIL 혹은 SLO 미달을 쉽게 시각화 가능
- `tool`, `scenario_type`, `test_name`을 기준으로 비교/추세 분석 가능
- `raw_result_path`는 evidence 추적성(원본 재확인)을 제공

## 5. 왜 CAT 본질(명확한 PASS/FAIL)을 해치지 않는가

- CAT의 `final_pass_fail`은 tool exit code에 의해 결정되고, 그 값은 `cat-result.json`에 이미 고정된다.
- Evidence-ready 파생 결과는 `cat-result.json`을 읽어 flatten하는 과정만 수행한다.
- 즉 Evidence-ready 스키마는 “시각화용 파생 뷰”일 뿐, CAT의 권위를 재판정하지 않는다.

