# 프로젝트 브리프: 범용 CAT Runner 검증 + Evidence 준비

## 이번 작업의 목적
이전 단계에서 확인한 “adapter 기반 raw 정규화 구조가 동작한다”를 **k6 전용이 아니라는 관점에서 범용성 있게 검증**한다.  
구체적으로 **Ginkgo adapter를 실제로 구현하고**, Ginkgo의 자연스러운 raw/result 경로를 사용해 동일한 `cat-result.json`으로 정규화되는지 end-to-end로 증명한다.  
동시에 CAT 결과를 Evidence 시각화로 연결하기 위해, **CAT 본질(표준 PASS/FAIL 권위)을 해치지 않는 최소 전처리 스키마**를 함께 확정한다.

## 왜 Ginkgo adapter 구현이 중요한가
Ginkgo는 k6처럼 “정해진 summary JSON을 내보내는” 형태가 아닐 수 있다.  
대신 테스트 출력(text)이나, 테스트가 생성한 커스텀 JSON 등 다양한 방식으로 raw을 남길 수 있다.  
따라서 Ginkgo를 실제로 붙여보는 것은 “raw 포맷이 달라도 adapter contract만 유지하면 CAT가 도구를 수용할 수 있는가?”를 검증하는 핵심 증거다.

## 왜 Evidence 시각화 스키마까지 같이 보는가
CAT이 `cat-result.json`을 내보내는 것만으로는 Evidence에서 바로 쓸 수 없을 수 있다.  
Evidence는 분석/시각화를 위한 **table(row/metric 중심) 형태의 파생 스키마**를 필요로 할 가능성이 높다.  
그러므로 CAT 원본을 뒤틀지 않는 선에서, “Evidence-ready 파생 결과”가 어떤 최소 구조를 가져야 하는지 함께 설계해, CAT 완성도를 높인다.

## 성공 기준
1. Ginkgo adapter가 실제로 runner에 붙고, `go test` 기반으로 Ginkgo custom CAT job을 실행한 뒤 `cat-result.json`을 생성한다.
2. runner 구조가 k6 전용이 아니라는 것을 코드와 실행 결과로 증명한다.
3. Evidence 시각화를 위해 필요한 최소 전처리 스키마(필드 요구사항 포함)가 정의된다.
4. `cat-result.json`과 Evidence-ready 스키마의 역할이 분리되며, CAT의 PASS/FAIL 권위가 손상되지 않는다.

