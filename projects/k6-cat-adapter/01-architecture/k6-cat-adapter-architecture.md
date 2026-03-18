# k6 CAT Adapter 아키텍처(표준 실행 구조)

**Date:** 2026-03-18  
**목표:** k6를 CAT에 “즉시 연결” 가능한 플러그인 형태로 표준화

## 1) CAT 구조 내 k6 위치
- **관계**: CAT(오케스트레이터/표준화) → k6(외부 실행 도구)
- **k6가 담당하는 테스트 유형(권장)**: HTTP/HTTPS 기반 외부 관측 테스트(예: ingress/API)
- **k6가 담당하지 않는 것**: 내부 SLI 중심 합격 판정(그건 CL2/Prometheus 등 다른 도구 영역)

## 2) 실행 흐름(표준)
```
CAT Job(입력)
  ↓ (시나리오/목표/SLO를 주입)
k6 실행(k6 run)
  ↓
k6 raw 결과 생성(--summary-export: k6-summary.json + 콘솔 로그)
  ↓
CAT Adapter 변환(k6-summary.json → cat-result.json)
  ↓
최종 PASS/FAIL 파일 저장(cat-result.json)
```

## 3) CAT 최소 조건 매핑(4요소)
- **Scenario Injection**
  - CAT Job의 `scenario`/`slo`/`target` → k6 스크립트 입력(환경변수 `__ENV` + 옵션)
- **SLI Measurement**
  - k6 `metrics`(예: `http_req_duration.p(95)`, `http_req_failed.value`, `http_reqs.rate`) → CAT 표준 SLI 필드로 정규화
- **SLO Evaluation**
  - k6 `thresholds` + **k6 exit code** → PASS/FAIL 단언
  - CAT는 **재판정하지 않는다**(기록만 한다)
- **Result Persistence**
  - raw: `k6-summary.json`
  - 표준: `cat-result.json` (CAT 공통 스키마)

## 핵심 설계 원칙(모호함 제거)
1. **판정 권위는 k6 exit code**다.
2. `cat-result.json`은 “판정 결과”를 **기록**하고, 후속 비교/누적/시각화를 위해 필드를 정규화한다.
3. k6 raw summary는 보관(증거), 표준 결과는 소비(판정/대시보드/리그레션).
