# Adapter 설계(책임 범위 고정)

## 공통 책임(모든 tool adapter)
- 입력: `job spec` + `output.dir`
- 실행: tool 호출(또는 호출되는 명령을 명확히 정의)
- 수집: tool raw 결과 경로를 찾는다
- 변환: raw → `cat-result.json` 생성(표준 스키마)
- 로직 원칙: **추가 평가(재판정) 금지**  
  - PASS/FAIL은 tool exit code로 결정하며, CAT 재판정은 하지 않는다.

## [k6 adapter]
- 실행
  - `k6 run --summary-export <output.dir>/k6-summary.json <script_path>`
- 입력(환경)
  - job spec의 `scenario`를 adapter가 `__ENV`/CLI env로 매핑
- 출력(raw)
  - `k6-summary.json`, `k6-output.txt`
- 처리(변환)
  - `k6-summary.json`에서 `selected_sli`에 필요한 metric만 추출
  - `cat-result.json` 생성

## [Ginkgo adapter]
- 실행
  - `go test <package-path> -v` (또는 `ginkgo run`)
- 입력
  - job spec env를 `SCENARIO_*`, `SLO_*` 형태로 테스트에 주입(환경변수)
- 출력(raw)
  - 테스트 코드가 `cat-result.json` 또는 artifacts를 생성(이번 설계에서는 “cat-result를 adapter가 먼저 만들도록”도 가능)
- 처리(변환)
  - `cat-result.json`이 있으면 수집/경로 기록
  - 없으면 adapter가 최소 변환(예: 결과 파일 위치/exit code 기록)만 수행

## [CL2 adapter]
- 실행
  - `clusterloader2 ...` (config.yaml + overrides)
- 입력
  - job spec의 cluster scenario config를 config 파일/오버라이드로 매핑
- 출력(raw)
  - 측정 JSON들(측정기별 결과)
- 처리(변환)
  - CAT `selected_sli`에 필요한 값만 추출해 `cat-result.json` 생성

## adapter가 “하지 않는 것”(중요)
- SLO를 다시 계산해 PASS/FAIL을 바꾸지 않는다.
- Grafana/Prometheus에서 “보고 나서” 결정을 바꾸지 않는다.
- 추가 데이터베이스 적재/시각화는 CAT 범위 밖.

