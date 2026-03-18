# CAT 결과 정규화 전략(가장 단순/확장 가능한 방식)

## 1) 목표
- 도구(k6/CL2/Ginkgo)의 raw 결과가 달라도,
  CAT가 비교/누적/시각화에 쓰는 “정규화 결과”는 항상 동일한 구조로 남긴다.
- 전처리/ETL을 만들지 않고, adapter가 얇게 유지되도록 한다.

## 2) Raw vs Normalized 분리
- Raw(Result of tool)
  - k6: `k6-summary.json`(또는 콘솔/artefact)
  - CL2: measurement 요약 파일/로그/JUnit 등
  - Ginkgo: 테스트 코드가 이미 만드는 artefact
- Normalized(Result of adapter)
  - CAT 표준 파일: `cat-result.json`
  - CAT는 이 파일만 Evidence/비교/누적에 사용한다.

## 3) 레이어(고정)
Tool Result → Adapter → Normalized Result

## 4) 표준 결과 구조(단순)
### cat-result.json (권장 최소)
```json
{
  "test_name": "ingress-basic",
  "tool": "k6|cl2|ginkgo",
  "scenario_type": "http|custom|cluster",
  "metrics": [
    { "name": "latency_p95_ms", "value": 210.0, "unit": "ms", "tags": { "case": "A" } },
    { "name": "error_rate", "value": 0.001, "unit": "ratio" }
  ],
  "status": "PASS|FAIL",
  "timestamp": "2026-03-19T00:00:00Z",
  "raw_result_path": "./results/run-001/raw/..."
}
```

### 왜 이렇게 단순한가
- `metrics[]`로 확장 가능(필드 고정이 아니라 항목 고정)
- Evidence/비교/누적이 “run 단위 + metric name/value”로 가능
- InfluxDB/시계열 모델로 강제되지 않는다.

## 5) SLI naming rule(필수)
CAT selected SLI 키는 아래 3개만 “최소 공통”으로 강제한다(이름/형식):
- `latency_p95_ms`
- `error_rate`  (0~1 ratio)
- `throughput_rps` (요청/초 or 작업/초)

추가 SLI가 필요하면 “이 규칙을 따르는 새 key”를 추가한다.

## 6) PASS/FAIL 권위(중복 재판정 금지)
- PASS/FAIL은 tool/테스트 runner의 exit code 또는 테스트 성공/실패를 권위로 한다.
- adapter는 `status`를 기록만 한다(재판정 로직 금지).

## 7) Evidence 연결(단순 전략)
- Evidence는 기본적으로 `cat-result.json`에서 selected SLI만 뽑아 시각화한다.
- Evidence가 JSON 파싱이 어렵다면, runner가 동일 내용을 `cat-result-table.csv`로 단순 파생 생성한다.
  - 이 csv는 “정규화 결과의 복사본”이며, ETL로 간주하지 않는다(스키마 고정/변환 없음 수준).

