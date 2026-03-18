# Final Report: k6 심층 이해 및 실험 기반 CAT 적용 가능성 확장

**Date:** 2026-03-18  
**대상 버전:** k6 v1.6.1

## 1) k6에 대한 최종 이해 요약
- k6는 “개발자 워크플로우에 들어가는 목표지향(load testing should be goal-oriented) 부하 테스트”를 목표로 설계된 도구이며, threshold를 통해 **실행 자체가 PASS/FAIL 신호**가 되도록 만든다. ([`Our beliefs`](https://k6.io/our-beliefs/))
- 내부 동작은 “Scenario(executor) → VU/iteration 실행 → metric 샘플 생성 → 집계 → threshold 평가” 흐름으로 이해하는 게 정확하다.

## 2) k6의 핵심 강점(구조적)
- **시나리오 주입이 단순**: 코드 + env로 target/부하/SLO를 주입할 수 있어 자동화 친화적.
- **외부 SLI가 즉시 나온다**: p95/p99 지연, 실패율, 처리량을 요약으로 확보.
- **SLO 게이트 내장**: threshold 위반 시 FAIL과 종료 코드로 파이프라인 판정이 단단해진다.

## 3) k6의 구조적 한계
- 내부 SLI(클러스터 리소스/제어면) 기반의 “수용 판정 주체”가 아니다.
- 포화(saturation)는 외부 지표만으로는 원인 분류가 어렵고, dropped_iterations 같은 간접 신호로만 힌트를 준다.

## 4) k6 결과의 활용 방식(저장/융합)
- raw: `--summary-export` JSON(`k6-summary.json`)은 근거 데이터로 보관 가치가 크다.
- CAT/CL2 융합: raw를 동일하게 만들려 하지 말고, `selected_sli + 판정 + 메타데이터`를 **CAT 표준 결과 파일**로 정규화하면 된다. (`k6-vs-cl2-format.md`)

## 5) CAT 적용 시 장점/단점(운영 판단 중심)
- **장점**: ingress/API처럼 “클라이언트 관측 SLI”로 합격을 정할 수 있는 영역에서, k6 하나로 실행/측정/게이팅이 닫힌다.
- **단점**: 내부 안정성(OOM/재시작/자원 포화)까지 합격 조건에 넣으면 k6 단독으로는 불충분하며, 내부 신호(예: Prometheus/CL2 측정)를 합성해야 한다.

## Final Question
👉 **k6는 단순 트래픽 도구인가, 아니면 CAT의 핵심 구성요소가 될 수 있는가?**

**결론: k6는 “단순 트래픽 생성기”에 그치지 않고, 외부 관측 기반 CAT에서 ‘핵심 구성요소(Scenario+SLI+SLO 게이트)’가 될 수 있다.**  
단, 그 범위는 “외부 SLI로 수용을 단언할 수 있는 테스트”로 명확히 제한해야 하며, 내부 SLI 기반 수용은 CL2/Prometheus 등과 결합해 CAT에서 합성 판정을 해야 한다.
