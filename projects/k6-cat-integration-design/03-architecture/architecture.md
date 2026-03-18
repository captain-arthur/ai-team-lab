# Architecture: k6 CAT 통합 설계(최소 요구 충족)

**Date:** 2026-03-18  
**Program:** CAT

## 1) k6의 CAT 내 역할(무엇을 맡기고, 무엇을 맡기지 않는가)
- **k6를 쓰는 테스트 유형(권장)**: “클러스터 외부 관측 기반” 테스트
  - 예: ingress/HTTP, API 게이트웨이, 서비스 엔드포인트의 p95/실패율 기반 수용
- **k6를 쓰지 않는 테스트 유형**: “내부 SLI가 합격의 핵심”인 테스트
  - 예: 제어면/노드/리소스/재시작 등 내부 신호가 1차 판정 기준인 경우(이 경우 k6는 부하 발생 도구로만 보조)

## 2) CAT 최소 4요소 매핑
| CAT 요구 | k6에서의 대응 | CAT 통합에서의 고정 규칙 |
|---|---|---|
| **Scenario Injection** | 스크립트 + `__ENV` + 실행 옵션 | CAT Job 정의에서 `target`, `scenario_type`, `params(env)`를 명시하고, 실행 시 환경변수로 주입 |
| **SLI Measurement** | k6 metric 집계(예: p95/failed/rate) | CAT가 저장할 “선택된 SLI” 리스트를 고정(필드명 표준화) |
| **SLO Evaluation** | k6 threshold 평가 + 종료 코드 | **종료 코드=권위**, 결과 파일에는 “어떤 threshold가 깨졌는지” 근거 포함 |
| **Result Persistence** | `--summary-export`(raw) | raw와 별개로 CAT 표준 결과 파일(`cat-result.json`)을 반드시 생성 |

## 3) CL2(clusterloader2)와의 관계(대체/보완/경계)
- **대체?** NO  
- **보완?** YES  
- **경계**
  - k6: 외부 트래픽 + 외부 SLI + 외부 SLO 게이트(HTTP 관점 수용)
  - CL2: 내부 부하 + 내부 SLI + 내부 SLO 게이트(클러스터 자체 수용)
  - CAT는 둘을 “동일한 결과 스키마”로 수렴시켜 비교/누적한다.

## 4) 최소 통합 흐름(입력→실행→측정→판정→저장)
### 입력
- CAT Job(도구= k6) 정의:
  - `test_name`, `script_path`, `target`, `scenario_type`, `params(env)`, `selected_sli`, `slo_policy`

### 실행
- CAT Runner가 k6를 실행:
  - `k6 run ... --summary-export <raw_summary.json> <script>`

### 측정
- k6 raw summary에서 “selected_sli”만 추출

### 판정
- **1차 판정(필수)**: k6 종료 코드 기반 PASS/FAIL
- **2차 판정(선택)**: (정책상 필요할 때만) Prometheus 등 내부 SLO를 추가 평가해 overall을 합성

### 저장(필수)
- 최소 2개 파일을 남긴다:
  1. `k6-summary.json` (raw: `--summary-export`)
  2. `cat-result.json` (표준: 도구 공통 결과 스키마)
