# Review Notes (최소 구현 검토)

## 이 최소 구현의 장점
- “raw 포맷 통일”을 목표로 하지 않고, adapter를 통해 **필요 값만 추출**하는 구조를 실제 실행으로 검증했다. (k6 end-to-end)
- `cat-result.json` 스키마가 고정되어 있어, tool이 달라도 CAT의 다음 단계(정규화/시각화/축적)에 동일한 표준 입력을 줄 수 있다.
- PASS/FAIL 권위는 tool exit code로 고정하여, CAT이 “재판정”하지 않는다는 원칙이 코드에서 유지된다.

## 아직 부족한 점
- Ginkgo/CL2 adapter는 현재 stub이라서 “실제 raw 파싱 결과가 표준 출력으로 수렴되는지”에 대한 런타임 증명은 없다.
- k6 adapter는 현재 `http_*` 메트릭에 강하게 결합되어 있다(다른 시나리오 유형 확장 시 매핑 규칙이 필요).
- raw 결과 파일 경로 정책은 runner가 정하는 경로(`k6-summary.json`)에 의존한다. tool이 항상 그 위치/이름을 제공하도록 adapter가 더 엄격히 관리해야 한다.

## raw 포맷 차이를 정말 흡수할 수 있는가?
이번 구현에서 “흡수”의 핵심은 다음처럼 나뉜다.
- CAT: raw 위치만 알고 adapter 호출
- Adapter: raw(JSON/XML/text)을 읽고, 필요한 SLI/SLO 정보만 골라 표준 결과로 변환

k6는 실제로 이 흐름이 동작한다. 즉 “구조적 가능성”은 확인됐다.  
Ginkgo/CL2는 stub이지만, adapter contract가 동일하므로 (구현만 하면) raw 차이를 흡수할 수 있는 형태로 남아있다.

## CL2 XML까지 고려했을 때 구조가 유지되는가?
네. adapter contract 수준에서는 CL2가 `json` 또는 `xml` raw를 제공할 수 있음을 전제한다.
- CAT 표준은 raw 내용을 알 필요가 없다.
- CL2 adapter만 XML/JSON parse 전략을 분기하면 되므로 구조가 유지된다.
다만 현재는 XML parse가 실제로 구현되진 않았다.

## 다음 단계에서 구현해야 할 것
- Ginkgo adapter:
  - raw 텍스트 출력 또는 커스텀 json 중 “어떤 우선순위/어떤 파일 포맷”을 표준적으로 선택할지 정하고 parse 구현
- CL2 adapter:
  - XML raw parsing 경로와 JSON parsing 경로를 둘 다 구현
- 스키마 확장:
  - SLI key set과 메트릭 매핑을 scenario_type 별로 명확히 정의(과설계 없이 최소 규칙만)

