# sample-cl2-sli-slo-analysis

ClusterLoader2 기반 SLI/SLO 분석 및 CAT v1 specification 작업이다. CL2 결과에서 측정 가능한 SLI 후보를 식별하고, baseline을 수립하며, SLO 가설을 검증한 뒤, devcat에서 사용 가능한 첫 CAT specification을 산출한다.

**Program:** 이 프로젝트는 **CAT** (Cluster Acceptance Testing) 프로그램에 속한다. 프로그램의 목적과 대표 프로젝트 유형은 `programs/cat/README.md`를 참고한다.

---

## 컨텍스트

- **입력:** Task intake(또는 이에 상응하는 문서), Manager brief, Research, Architecture, Engineering runbook.
- **단계:** Manager → Research → Architecture → Engineering → Experiment → Review → (validation) → CAT v1 spec.
- **주요 산출물:** SLI 후보 분석, metric baseline, SLO 가설, validation run, latency variance 분석, `08-cat-v1-spec/cat-v1-specification.md`.

상세 산출물은 단계별 폴더(`01-manager/` ~ `08-cat-v1-spec/`)를 참고한다.
