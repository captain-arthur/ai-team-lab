# Review Notes: k6 CAT 통합 설계 검토

**Date:** 2026-03-18

## 강점
- CAT 최소요구 4요소(주입/측정/판정/저장)를 **k6 실행 단위**로 직접 매핑했고, 특히 **결과 파일(`cat-result.json`)을 필수**로 고정했다.
- PASS/FAIL의 권위를 “종료 코드”로 고정해, 사람이 대시보드 보고 바꾸는 흐름을 차단했다.
- CL2와의 관계를 “대체 NO / 보완 YES”로 단정하고, 경계를 외부/내부 SLI로 분리했다.

## 아직 남은 모호한 부분(구현 전에 정해야 함)
- `k6-summary.json`에서 “failed thresholds 상세”를 어떤 방식으로 파싱해 `cat-result.json`에 넣을지(최소 구현 룰 필요).
- `http_reqs_rate_rps` 같은 “달성 처리량”을 SLO에 포함할지 여부(도입 시 under-drive 판정 정책 필요).
- Prometheus를 2차 게이트로 승격할 때 overall 판정(PASS_WITH_WARNINGS 등)을 도입할지 정책 결정.

## 바로 구현 가능한 수준인가?
- **YES(최소 구현 기준)**: 래퍼/어댑터가 (1) k6 실행, (2) 종료 코드 캡처, (3) summary에서 selected_sli 추출, (4) `cat-result.json` 생성까지만 하면 된다.

## 구현 전에 추가 확인이 필요한 것
- k6 버전별 `--summary-export` JSON 구조가 실무에서 파싱 안정적인지(필드 경로 고정성).
- CAT가 표준 결과 스키마를 도구 공통으로 확장할 계획이 있는지(현재는 최소 필드만 정의).
