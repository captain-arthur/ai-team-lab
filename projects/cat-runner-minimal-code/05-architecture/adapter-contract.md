# Adapter Contract (핵심)

각 tool adapter는 raw 포맷 차이를 CAT 표준 결과로 “번역”한다. CAT은 raw 포맷을 바꾸지 않고, adapter만 raw 포맷을 이해한다.

## 필수 인터페이스(문서 계약)
각 adapter는 아래 인터페이스를 따른다.

1. `run(job) -> raw 결과 생성`
2. `locate_raw_result(job) -> raw 결과 위치 반환`
3. `parse_raw_result(raw) -> 표준 SLI 추출`
4. `build_cat_result(parsed, exit_code) -> cat-result.json 생성`

## tool별 raw 포맷 차이 정리(흡수 책임)

### k6 adapter
- raw format: `json`
- input:
  - k6 `--summary-export` 결과 JSON 파일
- locate:
  - `output.dir/k6-summary.json` 같은 runner가 정한 경로를 기준으로 찾는다.
- parse:
  - `metrics`에서 `http_req_duration.p(95)`, `http_req_failed.value(=에러 비율)`, `http_reqs.rate` 같은 필요한 값만 추출한다.

### Ginkgo adapter
- raw format: `text` 또는 커스텀 `json`
- input:
  - Ginkgo 테스트 출력 텍스트 또는(가능하면) 결과 파일(예: `cat-result.json` 또는 별도 custom json)
- parse:
  - **우선순위는 결과 파일 우선**
  - 출력 텍스트 파싱은 fallback에 둔다(텍스트 포맷은 더 불안정하므로 지양).

### clusterloader2(CL2) adapter
- raw format: `json` 또는 `xml`
- input:
  - measurement/report 파일(JSON 또는 XML)
- parse:
  - 필요한 metric만 발췌해서 SLI로 정규화한다.
  - **구조 수준에서는 JSON/XML 둘 다 커버** 가능하도록 adapter 내부 parse 경로를 분기한다.

## 중요한 제한(과책임 방지)
- adapter는 “raw -> 표준 SLI/SLO 결과”까지만 책임진다.
- CAT은 시각화/DB 적재/추가적인 재판정을 하지 않는다.

