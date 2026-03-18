# Final Report (범용 CAT Runner + Evidence-ready 준비)

## 1) 이번 작업이 증명한 것
- CAT Runner 구조가 k6 전용이 아니라, **Ginkgo에도 실제로 적용 가능**하다는 점을 end-to-end로 증명했다.
- 특히 Ginkgo는 stdout 텍스트 파싱이 아니라, 테스트가 생성한 **커스텀 raw JSON(`ginkgo-raw.json`)**을 안정적인 입력으로 삼았다.
- runner의 adapter는 raw JSON을 필요한 metric만 추출해 `cat-result.json`으로 정규화한다.

## 2) Ginkgo adapter가 실제로 CAT Runner 구조에 들어오는지
들어온다.
- runner는 `tool: ginkgo` 일 때 실제 `go test ./ginkgo_cat_job ...`를 실행하는 adapter를 호출한다.
- tool exit code로 PASS/FAIL 권위를 고정하고,
- raw JSON 파일을 parse해서 표준 `cat-result.json`을 생성한다.

즉 adapter contract(실행 → raw 위치 → parse → 표준 결과 생성)가 Ginkgo에서도 그대로 성립한다.

## 3) CAT 결과와 시각화 결과를 어떻게 분리할지
- CAT 원본: `cat-result.json`
  - 권위(최종 PASS/FAIL)는 여기에 고정된다.
- Evidence-ready(파생): `cat-result.json`을 metric row table 형태로 flatten한 파일
  - CAT의 권위/판정 로직을 재판정하지 않는다.

따라서 “표준 결과는 유지”하면서 “시각화 친화적 뷰만 파생”하는 분리를 유지한다.

## 4) Evidence 시각화에 가장 적절한 파생 스키마
최소로는 아래 row 기반 스키마가 적절하다.
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

`status`는 `cat-result.json.slo_result[metric_name].ok`에 기반해 `ok/fail/n/a`로 둔다.

## 5) 최종 결론(질문에 대한 답)

1) “CAT Runner 구조는 k6 전용이 아닌가?”
- 아니다. k6에 이어 Ginkgo도 같은 runner 흐름으로 실제 정규화까지 수행되어, adapter 기반 구조의 범용성이 확인됐다.

2) “Evidence 시각화를 위해 어떤 전처리 결과가 필요한가?”
- `cat-result.json`을 metric 중심 row table 형태로 flatten한 **Evidence-ready 파생 스키마**가 필요하다.

최종 답 한 줄:
**“CAT Runner는 tool별 raw를 adapter가 흡수해 표준 `cat-result.json`을 만들고, Evidence는 그 표준 결과를 metric row로 flatten한 파생 뷰를 사용해야 한다.”**

