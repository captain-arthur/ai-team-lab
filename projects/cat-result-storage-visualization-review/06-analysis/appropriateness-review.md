# InfluxDB 중심 구조의 적절성 분석( CAT 관점 )

**Date:** 2026-03-19

## 1) InfluxDB 중심 구조가 정말 적절한가?
- 결론: **“필수”는 아니다.**
- 이유
  - CAT의 1차 소비는 결국 `cat-result.json` 기반의 PASS/FAIL과 selected_sli 비교/누적이다.
  - 시계열은 원인 분류/추세 관측에 유용하지만, 모든 run에서 필수로 넣을 필요는 없다.

## 2) k6 친화적으로만 흘러가버릴 위험
- 옵션 B(InfluxDB 중심)로 가면 위험이 커진다.
  - k6는 Influx output이 비교적 자연스럽지만
  - CL2/Ginkgo는 “어떤 시간축 metric을 어떤 태그/필드로 넣을지”가 강제로 커진다.
- 따라서 권장 전략은
  - **혼합(옵션 C의 방향)**: `cat-result.json`은 고정, 시계열은 필요할 때만.

## 3) CL2 / Ginkgo를 억지로 맞추게 되는가?
- “InfluxDB에 동일한 시계열 스키마를 강제”하면 억지 정합이 된다.
- 이를 피하는 방법
  - CL2/Ginkgo는 파일 기반 selected_sli로 충분히 CAT 본체를 충족
  - Influx는 k6에서만 먼저(또는 제한된 metric만) 적용

## 4) CAT의 본질(단순 PASS/FAIL + 결과 저장)을 해치는가?
- 옵션 B는 해칠 가능성이 있다.
  - PASS/FAIL을 관측 저장소에서 재판정하려는 유혹이 생기고,
  - “CAT의 파일 기반 단언”이 흐려질 수 있다.
- 방어 규칙
  - PASS/FAIL은 항상 tool exit code(또는 테스트 성공/실패)로 고정
  - InfluxDB는 관측/디버깅 보조에 머문다.

## 5) 지금 단계에서의 추천(최종)
- 추천: **혼합 구조(옵션 C 방향)**
  - `cat-result.json`(파일)은 고정/권위 유지
  - 시계열은 “k6에서 먼저, 필요할 때만” 추가
  - Evidence는 기본 리포트는 `cat-result` 기반, 원인/추세는 필요 시 시계열 소스를 추가

