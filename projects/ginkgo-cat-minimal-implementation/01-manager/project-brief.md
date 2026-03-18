# Ginkgo 실전 활용 완전 정복 가이드

**Date:** 2026-03-18  
**Program:** CAT

## 목표(이번 결과물의 끝 지점)
- 이 문서/코드 묶음만 보고 Ginkgo/Gomega로 **CAT 스타일 테스트(Job)**를 스스로 만들 수 있는 수준까지 도달한다.
- 특히
  - Ginkgo 구조/철학을 읽는 법
  - lifecycle를 안전하게 쓰는 법
  - Gomega의 Eventually/Consistently를 상태 검증에 쓰는 법
  - CAT 최소 블록(Scenario/SLI/SLO/Persistence)을 코드로 조립하는 법
  - custom CAT 시나리오 2개를 실제로 구현하는 법
  를 “패턴+재사용 가능한 코드”로 제공한다.

## 이번 작업이 증명/해결하려는 것
- Ginkgo가 custom CAT Job의 실행 엔진이 될 수 있는지(이번 구현 범위를 넘어 실무 적용 가능성까지).
- 결과 파일 저장(`cat-result.json`)과 PASS/FAIL 단언이 테스트 코드 내부에서 완결되는지.

## 성공 기준(다음 조건이면 완료)
- 문서만 보고 따라가면 예제 5개가 실행 가능한 수준이다.
- custom CAT 시나리오 2개가 Go 코드로 조립되어 있고, 각각 JSON 결과가 생성된다.
- “안티패턴”을 알고, lifecycle/비동기 검증/전역 상태 오염을 피할 수 있다.
