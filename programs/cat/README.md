# Program: CAT (Cluster Acceptance Testing)

**목적:** CAT(Cluster Acceptance Testing)는 Kubernetes 클러스터의 수용 기준을 정의하고 구현하기 위한 프로그램이다. 이 프로그램에서 진행하는 작업은 SLI/SLO 정의, 측정 파이프라인, 그리고 클러스터가 정의된 수용 기준을 충족하는지 판단하는 evaluator를 산출한다.

---

## 해결하는 문제

- **클러스터 적합성 불명확:** 특정 사용 사례에 대해 “클러스터가 충분히 좋다”는 공통 정의가 없음.
- **임시 검증:** 공식 SLI/SLO나 반복 가능한 pass/fail 기준 없이 부하·스트레스 테스트만 수행.
- **회귀 위험:** 클러스터 설정이나 워크로드 변경이 안정성·성능을 저하시켜도 이를 막는 gate가 없음.
- **측정 공백:** 수용 판단에 필요한 metrics(pod startup, OOM, API latency 등)가 부재하거나 일관되지 않음.

CAT 작업은 SLI, SLO를 정의하고, 도구(예: devcat)가 산출할 수 있는 명확한 pass/fail(및 선택적 warn) 결과를 만들어 위 문제를 다룬다.

---

## 이 프로그램에 속하는 프로젝트 유형

CAT 아래에 포함되는 대표적인 프로젝트는 다음과 같다.

- **CL2 SLI/SLO 분석** — ClusterLoader2로 측정 가능한 SLI 후보를 식별하고, baseline을 수립하며, SLO threshold 제안(예: pod startup latency P99, OOM count, system pod restarts).
- **CAT evaluator 구현** — 실험 결과를 읽어 CAT specification에 따라 PASS / PASS_WITH_WARNINGS / FAIL을 출력하는 도구 구현 또는 확장.
- **Prometheus measurement 통합** — Prometheus 기반 metrics(예: API responsiveness, slow calls)를 CAT SLI 세트 및 evaluator에 통합.
- **Load test 검증** — 클러스터에 부하를 주는 시나리오를 설계·실행하고, 부하 하에서도 SLI/SLO가 만족되는지 검증.
- **CAT specification 정제** — CAT v1(또는 이후) spec 정의·갱신: stable vs provisional SLI, threshold, decision logic.

주요 산출물이 “더 나은 클러스터 수용 기준 또는 이를 평가하는 도구”인 프로젝트는 CAT에 속한다.

---

## 프로그램과 프로젝트의 관계

- **Program:** 장기 도메인(CAT). 관련 프로젝트를 묶고, 공유 컨텍스트·spec·향후 작업(예: chaos, capacity 등 다른 프로그램에 둘 수 있음)의 근거지를 제공한다.
- **Projects:** AI 팀 워크플로(manager → research → architecture → engineering → experiment → review → documentation)로 실행되는 구체적·유한한 작업. 각 프로젝트는 명확한 목표와 산출물을 가지며, spec, 코드, runbook, 분석 등을 만들어 프로그램에 피드백한다.

CAT 소속 프로젝트는 물리적으로 `programs/cat/` 아래에 둘 필요가 없다. `projects/<project-name>/`에 두고 README(또는 brief)에서 **CAT** 프로그램에 속함을 **선언**한다. 이 프로그램 README(본 문서)는 CAT가 무엇인지, 어떤 종류의 프로젝트를 포함하는지 설명한다.
