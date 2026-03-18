# Final Report (최소 CAT Runner: adapter 기반 raw 정규화)

## 1) 이번 구현이 증명한 것
- 서로 다른 raw 포맷(JSON / XML / text) 자체를 CAT이 “통일”하려 하지 않아도,
  **adapter가 필요한 값만 추출**하면 `cat-result.json`으로 정규화하는 CAT Runner 구조가 실제로 동작한다.
- 특히 k6는 end-to-end로 실행되어,
  `k6-summary.json(표준 raw JSON)` → `cat-result.json(표준 CAT 스키마)` 변환이 완료된다.

## 2) 아직 stub인 부분
- Ginkgo adapter: raw 텍스트 출력/커스텀 json 중 어떤 경로를 기본으로 삼을지에 대한 parse 구현이 아직 없다.
- clusterloader2(CL2) adapter: json/xml raw 둘 다 고려한 parse 구현이 아직 없다.

즉, 현재는 “같은 adapter contract를 통해 구현 가능”한 구조적 준비를 완료했고, 런타임 파싱 증명은 k6만 포함한다.

## 3) raw 포맷 다양성을 adapter가 흡수하는 구조의 타당성
이번 minimal implementation에서 확인한 타당성은 다음 이유로 합리적이다.
- CAT 표준 출력이 `cat-result.json` 1개로 고정되면,
  raw 포맷 차이는 adapter가 감당하는 것이 가장 단순한 책임 분리다.
- PASS/FAIL 권위는 tool exit code로 고정하면, CAT이 재판정 로직을 가지지 않아도 된다.
- CL2 XML/JSON처럼 raw 형태가 확실히 다른 경우에도 adapter 내부에서만 분기하면 되므로, 구조가 깨지지 않는다.

## 4) CAT Runner 최소 구조에 대한 판단
결론적으로, “CAT Runner 최소 구조는 실제로 가능하다”는 판단이다.
이번 프로젝트는 특히 아래 2가지를 동시에 만족한다.
- tool 실행과 raw 위치 파악은 CAT(Runner)이 담당
- raw → SLI/SLO 추출 → 표준 결과 생성은 adapter가 담당

## 5) 최종 결론: CAT Runner는 raw 결과를 어떻게 다뤄야 하는가?
서로 다른 raw 결과 포맷(JSON/XML/text)을 가진 tool들을,
CAT Runner는 **raw를 통일하려 하지 말고 adapter가 필요한 값만 추출해 `cat-result.json`으로 정규화**해야 한다.

즉 최종 답은 다음 한 문장으로 정리된다.

**“CAT Runner는 raw 포맷을 adapter가 흡수하고, 표준 출력은 `cat-result.json`로만 고정해야 한다.”**

