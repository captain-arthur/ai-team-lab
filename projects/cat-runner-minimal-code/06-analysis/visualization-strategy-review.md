# Visualization Strategy Review (CAT vs Evidence)

## 1. CAT 최종 결과와 시각화 데이터는 같은가, 다른가
- 같은 점: 둘 다 `cat-result.json`에서 출발한다(어댑터가 만든 표준 결과).
- 다른 점: Evidence 시각화용 데이터는 “flatten된 파생 결과”다.
  - `cat-result.json`은 한 테스트 단위의 구조적 evidence(`selected_sli`, `slo_result`)를 가진다.
  - Evidence-ready는 metric 중심 row table로 변환된다.

즉 Evidence를 위한 “보기/분석 친화적 형태”로 바뀌며, 권위 데이터(최종 PASS/FAIL)는 건드리지 않는다.

## 2. `cat-result.json`만으로 충분한가?
작은 규모에서는 가능할 수 있다.
하지만 Evidence 관점에서,
- `selected_sli`가 map이므로 시각화/필터링 SQL/쿼리 구성에 번거로움이 생기고
- metric별로 status(ok/fail)를 group-by 하려면 파싱/flatten 로직이 필요하다.

따라서 최소한의 파생(전처리) 단계가 Evidence 친화성을 크게 높인다.

## 3. 별도 “Evidence-ready” 파일이 필요한가?
필요하다고 판단한다.
- CAT 원본을 보존하면서 Evidence에서 바로 쓰는 형태를 제공할 수 있기 때문이다.
- “adapter 단계에서 Evidence까지 책임지는” 과설계를 피할 수 있다.

## 4. 필요하다면 어떤 시점에 생성하는 것이 가장 자연스러운가?
추천 우선순위:
1. runner 이후 별도 export 단계(권장)
   - `cat-result.json`이 생성된 뒤, 별도의 `export-evidence-ready`가 읽고 변환한다.
   - CAT의 책임(정규화/표준 출력)을 깨지 않는다.
2. runner 내 optional step
   - 아주 작은 프로젝트에서는 runner 옵션으로 파생 파일을 같이 만들 수도 있다.
3. adapter 단계
   - adapter는 raw→표준 결과에 집중해야 하므로, 시각화 요구를 끌어오면 책임이 커진다.

이번 최소 구현 범위에서는 1번(별도 export 단계)이 가장 자연스럽다.

## 5. 최종 추천: 3종 역할 관계를 명확히
- raw 결과: tool이 남기는 원본(JSON/XML/text)
- `cat-result.json`: CAT 표준 출력(권위: tool exit code 기반 PASS/FAIL)
- Evidence-ready 결과: `cat-result.json`을 flatten한 파생 뷰(시각화용)

정리:
raw은 보관 증거, `cat-result.json`은 표준 인터페이스, Evidence-ready는 분석 편의 파생 결과다.

