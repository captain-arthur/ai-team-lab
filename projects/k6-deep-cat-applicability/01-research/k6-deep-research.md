# k6 심층 리서치(왜 이런 도구인가)

**Date:** 2026-03-18

## 1) k6의 역사(누가/언제/왜)
- **누가**: Load Impact(스웨덴)에서 시작, 이후 Grafana Labs로 합류. ([`About us`](https://k6.io/about/))
- **언제**: 2016년 “새 오픈소스 도구를 바닥부터” 만들기 시작. ([`About us`](https://k6.io/about/))
- **왜**: “한 번 하는 부하 테스트”가 아니라, API 중심 환경에서 **지속적으로(자동화로) 성능 목표를 검증**하려는 수요가 커졌고, 이를 개발자 워크플로우/자동화 흐름에 넣기 위해. ([`About us`](https://k6.io/about/), [`Our beliefs`](https://k6.io/our-beliefs/))

## 2) k6의 설계 철학(왜 그렇게 설계했나)
핵심은 “부하 테스트를 **개발자 도구**로 만들겠다”이다. ([`Our beliefs`](https://k6.io/our-beliefs/))

- **기존 도구와의 차이(의도)**  
  - 테스트를 “툴 설정”이 아니라 **코드(Everything as code)**로 만들고 버전관리/리뷰/자동화에 올리기 위함.  
  - “목표 지향(goal-oriented)”으로 threshold를 통해 **실행이 곧 PASS/FAIL 신호**가 되도록 하기 위함. (CI가 비0 종료 코드로 실패를 감지)  

- **왜 JS 기반인가(의도)**  
  - 개발자 접근성을 최우선으로 두고, 테스트 스크립트를 “일반 코드”처럼 다루게 하기 위함(학습/리뷰/조합). (`Our beliefs`: Everything as code, Developer experience)

- **왜 VU 모델인가(의도)**  
  - “동시 사용자”라는 직관적 모델로 시나리오를 표현하고, 사용자 여정(여러 요청+대기)을 반복 실행하기 쉬움. (`Our beliefs`: scenario load test / unit load test 구분)

## 3) “k6”라는 이름의 의미
- 공식 문서(About/Beliefs)에서 **왜 ‘k6’인지의 어원/약어 의미는 명시적으로 확인되지 않았다**.  
- 확인 가능한 사실: k6는 Load Impact가 만든 오픈소스 도구 이름이며, 이후 회사 브랜딩과 결합되어 사용되었다. ([`About us`](https://k6.io/about/))

> 결론: “이름의 정확한 유래”는 본 POC의 판단(도구 적합성)에는 영향이 없으므로, 공식 출처가 확보되기 전까지는 **불명(Unknown)** 으로 취급한다.
