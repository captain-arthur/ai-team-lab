# k6 실험: 골든 시그널 기반으로 읽기(실행 근거 포함)

**Date:** 2026-03-18  
**k6 버전:** v1.6.1  
**실험 위치:** `04-engineering/`  

## 1) 실험 목표
- k6 결과만으로 **골든 시그널**(Latency / Error rate / Throughput / (간접) Saturation)을 읽을 수 있는지 검증한다.
- SLI→SLO(threshold)→PASS/FAIL이 “운영자가 의사결정”에 바로 연결되는지 확인한다.

## 2) 실험 시나리오(직접 실행)
- **대상**: 공개 HTTP 엔드포인트 `https://test.k6.io/`
- **부하 모델**: `constant-arrival-rate` (목표 처리량 주입)
- **목표**: 80 iters/s, 30s
- **스크립트**: `scripts/golden-signals-http.js`

## 3) 측정 대상(골든 시그널 매핑)
- **Latency**: `http_req_duration` p(95), p(99)
- **Error rate**: `http_req_failed` rate + `checks` rate
- **Throughput**: `http_reqs` rate (달성치)
- **Saturation(간접)**:
  - `dropped_iterations` > 0 이면 “목표 도착률을 못 맞춘 구간이 존재”(주입 부하 대비 포화/제한 가능성)
  - `vus`가 예상보다 크게 튀면 “목표 도착률을 맞추기 위해 동시 실행자 확장”(자원/대기 증가 가능성)

## 4) 실행(명령어/환경)
- **환경**: 로컬(macOS), 외부 인터넷 접근

```bash
cd projects/k6-deep-cat-applicability/04-engineering

TARGET_URL=https://test.k6.io/ \
MODE=arrival TARGET_RPS=80 DURATION=30s \
SLO_P95_MS=400 SLO_FAIL_RATE=0.01 \
k6 run --summary-export results/golden-arrival-summary.json scripts/golden-signals-http.js \
  | tee results/golden-arrival-output.txt
```

## 5) 결과 파일(필수)
- `results/golden-arrival-summary.json` (k6 raw summary)
- `results/golden-arrival-output.txt` (콘솔 로그)

## 6) 결과 해석(수치→SLO 비교→판정)
### 핵심 관찰(로그 발췌)
- threshold: `p(95)<400` **PASS**, `rate<0.01` **PASS**
- `http_req_duration p(95)=217.95ms`, `p(99)=261.11ms`
- `http_req_failed rate=0.00%`
- `http_reqs rate=157.99/s` *(k6는 요청 1회보다 더 많은 하위 요청이 발생할 수 있어, iters/s와 다를 수 있음)*
- `dropped_iterations=16` (30초 동안 일부 iteration 목표 미달)

### 운영자가 이 결과로 무엇을 판단할 수 있는가
- **수용(PASS/FAIL)**: threshold가 만족되므로 “외부 SLI 기준”으로는 **PASS**라고 단언 가능(종료 코드 0).
- **성능 상태**: p95 218ms 수준으로 “지연 분포”가 목표(400ms) 이하에 들어옴.
- **포화 징후**: dropped_iterations가 0이 아니므로, “목표 도착률을 항상 만족한 것은 아님” →  
  동일 환경에서 RPS 상향 시 **포화/제한이 더 명확히 드러날 가능성**.  
  (실제 클러스터 적용 시에는 이 시점에 내부 신호(Prometheus)로 원인을 좁히는 게 합리적)

## 7) PASS/FAIL 판단 과정(명문화)
1. 선택 SLI를 확인한다(p95 latency, 실패율, 처리량, dropped_iterations).
2. threshold(SLO)와 비교한다.
3. threshold 위반이면 FAIL, 모두 만족이면 PASS.
4. 결과는 **로그 + summary JSON**으로 보존한다(후속 CAT 표준 포맷 변환 가능).
