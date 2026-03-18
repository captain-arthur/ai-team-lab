# CAT Runner 아키텍처 (최소)

## 전체 흐름
Job Spec
→ Runner
→ Tool Adapter
→ raw result 수집(파일 위치 확보)
→ parser/adapter(필요 값 추출)
→ `cat-result.json` 생성(표준 출력)

## raw 포맷 다양성 명시
CAT은 raw 포맷을 통일하지 않는다. 도구별로 raw 결과 형태가 다를 수 있고, adapter가 그 차이를 흡수한다.

- k6: summary-export 기반 **JSON**
- Ginkgo: 테스트 출력 **텍스트** 또는 adapter가 선택하는 **커스텀 JSON 결과 파일**
- clusterloader2: 측정/리포트 **JSON 또는 XML**

## CAT의 책임(Runner/Adapter 경계)
CAT이 직접 하는 일은 아래로 고정한다.

1. tool 실행(또는 실행 스크립트 호출)
2. raw 결과 파일/산출물 위치 파악
3. adapter 호출(파싱/정규화 위임)
4. 표준 `cat-result.json` 생성 및 저장

## CAT가 하지 않는 것(명시적 제외)
아래는 CAT의 책임으로 두지 않는다.

- tool 내부/기본 결과 포맷을 바꾸도록 강제(예: k6 출력 포맷 변경)
- “재판정”(tool exit code 또는 test success 여부를 CAT이 뒤집지 않음)
- 시각화/DB 적재(범위 밖)

