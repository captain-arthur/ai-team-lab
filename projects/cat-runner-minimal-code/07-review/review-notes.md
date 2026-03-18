# Review Notes (Ginkgo 범용성 + Evidence 준비)

## Ginkgo adapter를 실제로 붙여본 결과 구조가 유지되는지
유지된다.
- runner는 동일한 adapter contract로 tool을 선택하고 실행한다.
- Ginkgo는 stdout 파싱이 아니라 **테스트가 생성한 커스텀 raw JSON(`ginkgo-raw.json`)**를 CAT 입력으로 사용한다.
- adapter는 raw JSON에서 필요한 metric만 추출해 표준 `cat-result.json`으로 정규화한다.
- `final_pass_fail`은 exit code 기반으로 기록되어 CAT이 재판정하지 않는다.

따라서 “raw 포맷을 통일하지 않는다 → adapter가 흡수한다”라는 구조 원칙이 Ginkgo에서도 성립함을 실제 실행으로 확인했다.

## runner가 k6 전용이 아니었는지(코드 관점)
네.
- `main.go`의 tool switch가 `k6` 외에 `ginkgo` 케이스를 실제 adapter로 연결한다.
- Ginkgo adapter는 env 주입 + `go test` 실행 + raw JSON 파싱의 형태로 contract를 따른다.

즉 runner 공통 흐름이 유지되고, tool 특이 처리는 adapter 내부에 머문다.

## Evidence용 전처리 스키마가 과한지/적절한지
적절하다(최소).
- Evidence에 필요한 것은 “시각화/비교에 좋은 row table 형태”다.
- 제안한 schema는 metric별로 필요한 최소 컬럼만 flatten하며,
  CAT의 표준 출력 구조(`cat-result.json`)를 변형하거나 권위를 재판정하지 않는다.

## 아직 남은 모호한 점
- Evidence-ready 파생 파일의 파일 포맷(예: JSON rows vs CSV)은 Evidence ingest 방식에 따라 최종 결정이 필요하다.
- metric_name → metric_unit 매핑 규칙을 `cat-result.json` 또는 문서/코드 어느 쪽에서 단일화할지(이번 단계에선 문서 규칙 수준으로만 확정).
- CL2 XML parsing까지 포함한 “완전한 범용성 증명”은 다음 단계에서 실제 raw 샘플을 기반으로 보강할 필요가 있다.

## 다음 단계에서 무엇을 해야 하는지
1. (선택) `cat-result.json` → Evidence-ready 파일을 생성하는 별도 export 스크립트를 구현한다.
2. CL2에서 XML raw를 실제로 받았을 때 adapter parse/정규화 경로를 확장한다.
3. Evidence ingest 포맷(JSON/CSV) 확정 후, sample data로 end-to-end 시각화까지 확인한다.

