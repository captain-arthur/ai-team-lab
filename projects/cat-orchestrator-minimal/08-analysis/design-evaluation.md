# 설계 평가(장점/한계/확장 시 문제)

## 장점
- 과설계가 없다: 플러그인/분산/DB/Grafana에 의존하지 않고, 단일 CLI + 파일 기반으로 충분하다.
- CAT 최소 책임이 명확하다
  - Job Spec을 읽고
  - tool을 실행하고
  - adapter를 호출하고
  - 결과를 저장/확인한다
- PASS/FAIL 권위를 단일하게 고정한다
  - tool exit code = 최종 판정
  - CAT은 재판정하지 않는다

## 한계(솔직히)
- adapter의 “raw → cat-result” 변환 로직은 tool별로 필요하다.
- Ginkgo는 테스트 코드 내부에서 결과 파일을 생성할 수 있지만, k6/CL2는 raw 파싱/추출이 필요하다.
- CL2에 대해 어떤 SLI를 뽑을지(정의/필드 매핑)가 확정되어야 adapter가 안정적으로 동작한다.

## 확장 시 문제점(다음 단계에서 생길 수 있는 것)
- selected_sli 키 표준화가 흐려지면 adapter마다 ad-hoc 변환이 늘어난다.
- tool별 “exit code가 의미하는 바”가 정책과 다르면 전체 판정이 흔들릴 수 있다(그래서 권위=exit code를 고정).
- multi-job 전체 overall 판정(aggregator)이 필요해지면, warning/fail 정책을 별도로 정해야 한다.

## 왜 이 정도가 ‘최소 실무형’인가
- 우리가 이미 확보한 도구 모델(k6/CL2/Ginkgo)을 **같은 실행 계약(표준 Job + 표준 결과)**로 묶기 때문이다.
- 요구사항(Scenario/SLI/SLO/Persistence) 중 “파일 결과 통일”과 “PASS/FAIL 권위”를 가장 먼저 고정한다.

