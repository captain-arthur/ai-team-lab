# core-documents-index.md (CAT + Monitoring 최종 선별)

이 문서는 “이름이 비슷한 문서가 많아지는” 문제를 막기 위한 **목록형 안내서**다. 각 섹션은 길어 보이지만, 실제로 필요한 건 소수의 문서뿐이다.

---

## 1) CAT 프로젝트 핵심 문서 (Top 5~10)
1. `projects/cat-orchestrator-minimal/02-architecture/cat-overall-architecture.md`  
   - 한 줄 설명: CAT 오케스트레이션의 최소 구성요소(Job Spec / Runner / Adapter / Result Store)와 고정 책임 경계를 정의한다.  
   - 왜 중요한가: “CAT이 무엇을 안 하는지(재판정 금지)”를 먼저 고정하지 않으면 이후 모든 문서가 흔들린다.  
   - 언제 봐야 하는가: CAT 시스템을 ‘계속 운영’할지 판단하기 전(시작/정리 단계).
2. `projects/cat-orchestrator-minimal/03-architecture/cat-job-spec.md`  
   - 한 줄 설명: tool-공통 입력(Job Spec) 스키마와 scenario 주입 규칙을 고정한다.  
   - 왜 중요한가: 어댑터 확장 시 “어느 필드까지가 계약인지”를 팀이 합의할 수 있게 한다.  
   - 언제 봐야 하는가: k6/Ginkgo/CL2를 추가하거나 교체할 때(변경 시점).
3. `projects/cat-orchestrator-minimal/05-architecture/adapter-design.md`  
   - 한 줄 설명: Adapter 책임 범위를 ‘변환/정규화만’으로 고정하고, PASS/FAIL 재판정 금지를 명문화한다.  
   - 왜 중요한가: adapter에 로직을 섞는 순간 CAT의 권위 모델(exit code)과 충돌한다.  
   - 언제 봐야 하는가: 누군가 “대시보드/SLI로 다시 판단하자”라고 말할 때.
4. `projects/cat-orchestrator-minimal/06-architecture/cat-runner-design.md`  
   - 한 줄 설명: 실행은 단순 CLI(실행→raw 위치 파악→adapter→`cat-result.json`)로 끝낸다는 흐름을 준수한다.  
   - 왜 중요한가: 운영 리스크(스케줄러/플러그인 과설계)를 줄인다.  
   - 언제 봐야 하는가: “러너가 점점 플랫폼이 되어가는지” 점검할 때.
5. `projects/cat-orchestrator-minimal/07-engineering/minimal-runner-prototype.md`  
   - 한 줄 설명: `cat run job.yaml`이 실제로 파일 규약으로 동작하는 최소 프로토타입을 보여준다.  
   - 왜 중요한가: 문서가 아니라 ‘실행 가능한 설계 예시’라 논쟁의 끝을 빠르게 낸다.  
   - 언제 봐야 하는가: 구현 착수 전/PR 리뷰 때.
6. `projects/cat-runner-minimal-code/04-architecture/cat-result-spec.md`  
   - 한 줄 설명: `cat-result.json` 최소 스키마와 PASS/FAIL 권위(출구코드 기반)를 정의한다.  
   - 왜 중요한가: CAT 산출물의 “표준 인터페이스”가 무엇인지 한 번에 확인할 수 있다.  
   - 언제 봐야 하는가: Evidence/정규화/대시보드 입력을 연결할 때.
7. `projects/cat-runner-minimal-code/05-architecture/adapter-contract.md`  
   - 한 줄 설명: Adapter가 raw 포맷을 흡수해 표준 결과로 번역하는 계약(실행/수집/파싱/생성)을 문서화한다.  
   - 왜 중요한가: 개발자가 “어디까지가 adapter의 책임인지”를 반복해서 확인하는 기준점이다.  
   - 언제 봐야 하는가: adapter 구현 변경 시(특히 Ginkgo/CL2).
8. `projects/cat-runner-minimal-code/05-architecture/evidence-ready-schema.md`  
   - 한 줄 설명: Evidence는 cat-result를 flatten한 파생 뷰를 사용해야 한다는 관계를 정의한다.  
   - 왜 중요한가: Evidence에 ‘권위 데이터’를 섞는 실수를 막는다.  
   - 언제 봐야 하는가: 시각화/누적 요구가 늘어날 때.
9. `projects/cat-result-normalization-strategy/normalization-strategy.md`  
   - 한 줄 설명: Raw→Adapter→Normalized(표준)로 나누고, selected SLI naming rule을 강제하는 가장 단순한 전략을 제시한다.  
   - 왜 중요한가: CAT 결과가 도구별로 달라지는 문제를 ‘정규화 레이어’로 고정한다.  
   - 언제 봐야 하는가: metric 드리프트/누적/비교 기능을 추가할 때.

---

## 2) Monitoring 프로젝트 핵심 문서 (Top 5~10)
1. `projects/sre-failure-driven-core-components-runbook-dashboard/01-architecture/dashboard-final-lock.md`  
   - 한 줄 설명: 3개 질문(Row)과 9개 패널 배치, 컴포넌트별 4개 골든 시그널을 최종 고정한다.  
   - 왜 중요한가: “이 대시보드가 무엇을 위해 존재하는지”를 최종 스코프에 묶는다.  
   - 언제 봐야 하는가: 대시보드가 커질 때(패널 추가/metric 추가 논쟁).
2. `projects/sre-failure-driven-core-components-runbook-dashboard/02-engineering/golden-signals-final.md`  
   - 한 줄 설명: 각 컴포넌트별 ‘정말 필요한 4개’ 신호가 무엇인지 운영 판단에 어떻게 쓰이는지 고정한다.  
   - 왜 중요한가: metric 나열로 빠지는 걸 막는다.  
   - 언제 봐야 하는가: 임계값 튜닝/신호 교체 시.
3. `projects/sre-failure-driven-core-components-runbook-dashboard/03-engineering/promql-spec-final.md`  
   - 한 줄 설명: 9개 패널 각각의 PromQL(0/1/2 상태 점수 반환)과 점수화 규칙을 고정한다.  
   - 왜 중요한가: “바로 구현 가능한 수준”의 쿼리 계약이다.  
   - 언제 봐야 하는가: Grafana JSON 수정 전/후에 변경되는 지점을 확인할 때.
4. `projects/sre-failure-driven-core-components-runbook-dashboard/04-engineering/grafana-dashboard.json`  
   - 한 줄 설명: import 가능한 Grafana 대시보드 JSON(패널 9개, Row 3개, datasource placeholder 포함)을 제공한다.  
   - 왜 중요한가: 실제 도입을 막는 문서 공백을 없앤다.  
   - 언제 봐야 하는가: Grafana에 바로 붙일 때(가장 자주).
5. `projects/sre-failure-driven-core-components-runbook-dashboard/05-engineering/grafana-build-notes.md`  
   - 한 줄 설명: import 후 무엇을 바꿔야 하는지(라벨/metric 차이)를 환경 관점에서 적는다.  
   - 왜 중요한가: “metric 이름/라벨 drift”가 발생했을 때 수정 지점을 빠르게 찾는다.  
   - 언제 봐야 하는가: 패널이 데이터 없음/에러일 때.
6. `projects/sre-failure-driven-core-components-runbook-dashboard/06-analysis/practical-check.md`  
   - 한 줄 설명: 최소 신호가 정말 최소인지/PromQL이 너무 복잡한지/운영 도입성을 점검한다.  
   - 왜 중요한가: ‘작동 여부’ 관점에서 최종 확인 목록이 된다.  
   - 언제 봐야 하는가: 운영 도입 직전(체크리스트).
7. `projects/sre-failure-driven-core-components-runbook-dashboard/07-engineering/runbook-starter-set.md`  
   - 한 줄 설명: Row 3 Action에 대응되는 runbook starter를 제공한다.  
   - 왜 중요한가: 대시보드가 “보기”에서 “행동”으로 연결되는 유일한 고리다.  
   - 언제 봐야 하는가: Runbook 연계 검증(실제 operator가 어떤 링크로 들어갈지).
8. `projects/sre-failure-driven-core-components-runbook-dashboard/07-documentation/final-report.md`  
   - 한 줄 설명: 최종 산출물이 무엇인지(골든 시그널/Row/PromQL/JSON 생성)와 결론을 고정한다.  
   - 왜 중요한가: “이 작업이 끝났는지”를 빠르게 판정한다.  
   - 언제 봐야 하는가: 최종 제출/회고 단계.

---

## 3) “이건 보지 마라” 문서 (중복/실효성 낮음)
CAT 측
- `projects/cat-runner-minimal-code/06-analysis/visualization-strategy-review.md`  
  - 왜: Evidence 관점 서술이 있지만, 실제 구현에는 위의 `evidence-ready-schema.md`가 더 직접적이다.
- `projects/cat-runner-minimal-code/08-review/review-notes.md` , `projects/cat-runner-minimal-code/07-review/review-notes.md`  
  - 왜: 리뷰/노트는 설계 계약서가 아니라 ‘잡음이 많은 중간 기록’에 가깝다.

Monitoring 측
- `projects/sre-failure-driven-core-components-runbook-dashboard/02-engineering/panel-minimal-spec.md`  
  - 왜: 최종은 `dashboard-final-lock.md` + `golden-signals-final.md` + `promql-spec-final.md`로 이미 고정됐다.
- `projects/sre-failure-driven-core-components-runbook-dashboard/05-engineering/dashboard-blueprint-final.md`  
  - 왜: 최종 Row/패널 규칙은 `dashboard-final-lock.md`로 더 단단하게 잠겼다.
- `projects/sre-failure-driven-core-components-runbook-dashboard/04-engineering/grafana-build-guide.md`  
  - 왜: “빌드 방식 설명” 문서지만, 실제로 구현은 JSON이 더 권위적이다.
- `projects/sre-failure-driven-core-components-runbook-dashboard/03-architecture/operational-methodology.md` , `projects/sre-failure-driven-core-components-runbook-dashboard/04-architecture/runbook-linkage-model.md`  
  - 왜: 철학/방법론이지만, 최소 도입에는 runbook starter와 final lock/panel 계약이 충분하다.
- `projects/sre-failure-driven-core-components-runbook-dashboard/01-architecture/dashboard-final-model.md`  
  - 왜: “모델 확정” 문서지만, 이미 final lock으로 최신화됐다.

---

## 4) 추천 학습 순서 (완전 처음 보는 사람 기준)
1. CAT: `projects/cat-orchestrator-minimal/02-architecture/cat-overall-architecture.md`  
2. CAT: `projects/cat-orchestrator-minimal/03-architecture/cat-job-spec.md`  
3. CAT: `projects/cat-orchestrator-minimal/05-architecture/adapter-design.md`  
4. CAT: `projects/cat-runner-minimal-code/04-architecture/cat-result-spec.md`  
5. CAT: `projects/cat-result-normalization-strategy/normalization-strategy.md`  
6. Monitoring: `projects/sre-failure-driven-core-components-runbook-dashboard/01-architecture/dashboard-final-lock.md`  
7. Monitoring: `projects/sre-failure-driven-core-components-runbook-dashboard/02-engineering/golden-signals-final.md`  
8. Monitoring: `projects/sre-failure-driven-core-components-runbook-dashboard/03-engineering/promql-spec-final.md`  
9. Monitoring: `projects/sre-failure-driven-core-components-runbook-dashboard/04-engineering/grafana-dashboard.json`  
10. Monitoring: `projects/sre-failure-driven-core-components-runbook-dashboard/05-engineering/grafana-build-notes.md`  

