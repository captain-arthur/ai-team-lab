# k6 결과 포맷 분석(summary-export JSON)

**Date:** 2026-03-18  
**샘플 근거:** `04-engineering/results/golden-arrival-summary.json`

## 1) 최상위 구조
`--summary-export` JSON은 크게 2개 블록으로 읽으면 된다.
- **`metrics`**: 모든 요약 지표(집계 결과). CAT가 SLI를 뽑는 핵심 영역.
- **`root_group`**: check/group 구조 요약(어떤 check가 몇 번 pass/fail 했는지).

## 2) `metrics` 구조(유형별 패턴)
### Trend(분포형) 예: `http_req_duration`
- 필드 예: `avg`, `min`, `med`, `max`, `p(90)`, `p(95)`, `p(99)`
- **threshold 평가 결과**가 포함될 수 있음:
  - `thresholds: { "p(95)<400": false }`
  - 주의: 여기의 boolean은 “임계치 위반 여부”를 직접 의미하지 않는 것처럼 보일 수 있어(표현이 직관적이지 않음), CAT 판정 권위는 **종료 코드/콘솔 THRESHOLDS 섹션**을 기준으로 두는 편이 안전하다.

### Rate/Counter 예: `http_reqs`, `iterations`, `dropped_iterations`
- 필드 예: `count`, `rate`
- CAT에서 throughput/달성 여부를 계산할 때 사용.

### Gauge 예: `vus`, `vus_max`
- 필드 예: `value`, `min`, `max`
- arrival-rate 모델에서 VU 확장 정도(간접 포화/비용 신호)를 읽는 용도.

### Boolean-like 예: `http_req_failed`, `checks`
- 필드 예:
  - `http_req_failed`: `value`(0~1), `passes`, `fails`
  - `checks`: `value`(성공률), `passes`, `fails`

## 3) CAT에서 “selected SLI”로 뽑기 좋은 최소 필드
- `metrics.http_req_duration.p(95)` → `http_req_duration_p95_ms`
- `metrics.http_req_failed.value` → `http_req_failed_rate`
- `metrics.http_reqs.rate` → `http_reqs_rate_rps`
- (포화 간접) `metrics.dropped_iterations.count` 또는 `rate`

## 4) CAT 저장 관점 결론
- k6 raw summary는 **근거 데이터**로 보관 가치가 높다.
- 다만 CAT 표준 결과 파일은 별도로 만들고(`cat-result.json`), raw에서 selected SLI만 정규화해 저장하는 구성이 가장 안전하다.
