# k6 옵션과 실무 연결(운영자 관점)

**Date:** 2026-03-18

## 목적
옵션을 “나열”이 아니라, **Kubernetes 운영/수용(CAT) 관점에서 어떤 선택을 의미하는지**로 연결한다.

## 1) 주요 옵션(핵심만)
- **vus / duration**: VU 기반 실행에서 동시성/실행 시간.
- **stages**: 램프/스파이크 같은 부하 패턴(시간에 따른 목표 변화).
- **executor**: 부하를 정의하는 축 선택(VU vs arrival-rate).
  - `constant-vus`, `ramping-vus`
  - `constant-arrival-rate`, `ramping-arrival-rate`
- **thresholds**: SLO 단언(PASS/FAIL) 정의. 실패 시 비0 종료 코드.
- **sleep / think time**: 사용자 여정의 대기. 처리량(throughput)과 iteration 시간에 직접 영향.
- **--summary-export**: raw 요약 결과 JSON 저장(후속 파싱/보관용).

## 2) 옵션의 “의미”와 실무 연결
### executor(가장 중요한 선택)
- **VU 기반이 필요한 상황**
  - “동시 접속자 수” 자체가 시나리오의 본질일 때(예: 동시 세션 수가 중요한 서비스)
  - 사용자 여정(여러 step + 대기)이 핵심이고, 요청률은 결과로 받아들일 때
- **arrival-rate 기반이 필요한 상황**
  - ingress/API 수용 테스트처럼 “RPS 목표”를 명시해야 할 때(트래픽 주입 기반)
  - CAT에서 “성능 목표를 RPS 기준으로 고정”하고 회귀를 감지하고 싶을 때
  - `dropped_iterations`를 통해 “목표를 못 맞춤(포화)”을 간접 신호로 읽고 싶을 때

### thresholds(판정 설계)
- 운영에서 필요한 건 “그래서 배포/승인할 수 있는가”이며, threshold는 이를 **기계적으로 단언**하게 만든다.
- 권장 최소 세트(HTTP ingress/API):
  - `http_req_duration: p(95)<X`
  - `http_req_failed: rate<Y`
  - (필요 시) 처리량 달성(under-drive 방지) 규칙 추가

### think time(sleep)
- **latency를 일부러 늘리는 옵션이 아니라**, 사용자 여정 속도를 늦춰 **throughput을 줄이는 모델링 도구**다.
- 따라서 “VU 기반 + think time” 조합에서 RPS가 떨어지는 것은 “성능 저하”가 아니라 “덜 보낸 것”일 수 있다.

### summary-export(결과 보관/융합의 관문)
- k6는 raw summary를 JSON으로 저장할 수 있다(`--summary-export`).
- CAT 관점에서는 raw를 그대로 쓰기보다:
  - raw(`k6-summary.json`)는 보관
  - 선택 SLI + 판정 + 메타데이터는 별도 표준 파일로 정규화(예: `cat-result.json`)

## 3) ingress/API 테스트에 연결(최소 가이드)
- **Ingress 성능 수용(CAT)**: arrival-rate 기반 + p95/실패율 threshold + (필요 시) dropped_iterations 관찰
- **API 회귀 감지**: unit 테스트(단일 엔드포인트) 형태로 짧게, threshold를 엄격히
- **트래픽 시뮬레이션(사용자 여정)**: VU 기반 + group/step + think time으로 현실성 부여(단, 처리량 해석 주의)
