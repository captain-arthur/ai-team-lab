# devcat Program Brief (devcat 프로그램 브리프)

이 문서는 **실제 CAT 프로그램 방향**과 **devcat 저장소의 현재 상태**를 AI 팀 워크스페이스 안에서의 단일 소스 오브 트루스로 둔다. CAT 관련 작업·제안은 이 브리프와 cat-vision, cat-design-principles, sli-slo-philosophy에 맞춘다.

---

## 1. 내가 책임지는 것

- **실용적인 CAT 시스템**을 목표로 하고 있다. 논문·이론이 아니라 “실행하고, 결과 보고, 수용/거부 판단할 수 있는” 시스템이다.
- CAT 시스템은 최종적으로 **Conformance/기능성**과 **Performance/부하** 두 관점을 모두 지원해야 한다.
- **현재는** 연구와 실용적 방향 정립, 특히 **Performance/부하 수용 테스트** 쪽에 초점을 두고 있다. Conformance·기능성은 이후 단계에서 더 다룬다.

---

## 2. 현재 CAT 방향

- **완벽한 단일 CAT 도구는 없다.** 여러 도구를 실용적으로 조합한다.
- **테스트 구조:** 하나의 테스트는 **Job**으로 본다. Job 구성:
  - **Scenario injector** — 부하·시나리오 주입
  - **SLI/SLO measurement** — 지표 측정·SLO 준수 여부
  - **Assertion** — pass/fail 단언
- **결과**는 **정의된 디렉터리**에 저장하고, 가능하면 **자동 시각화**가 되게 한다.
- **오케스트레이션**은 단순하게. Prow처럼 복잡한 것은 정말 필요할 때만.
- **기능 테스트**는 나중에 Ginkgo 같은 도구를 쓸 수 있으나, 지금 당장의 초점은 아니다.

---

## 3. 현재 devcat 현실

- **devcat**은 실 구현이 이루어지는 **실제 구현 저장소**다 (이 AI 팀 워크스페이스와는 별도).
- **현재 사용 방식:**
  - 로컬 바이너리: **ClusterLoader2**, **perfdash** 등.
  - ClusterLoader2 기반 **부하 시나리오**를 실행한다.
  - 부하 시나리오를 복사해 두고, **config.yaml**을 실행하며 클러스터별 오버라이드 파일(예: **ol-test.yaml**)을 쓴다.
  - 결과는 **results/** 아래에 저장된다.
- **시각화:** **perfdash**를 쓰고 있으나, **assertion이나 PASS/FAILED 상태**를 명확히 보여 주지 않는다. 숫자·차트 위주다.

---

## 4. 내가 실제로 원하는 것

- **연구·문서화에서 멈추지 않는다.** 실제로 테스트를 **실행**하고, 결과를 **검사**하고, “어떤 지표가 왜 중요한지” 설명할 수 있을 만큼 **숙련**된 뒤, devcat을 점진적으로 **실용 CAT 시스템**으로 키우고 싶다.
- **AI 팀**은 “모호한 방향”을 **구체적인 구현·검증**으로 옮기는 데 기여해야 한다. 제안이 devcat의 현재 상태와 연결되고, 다음 실험·개선으로 이어져야 한다.
- 가능하면 **로컬 Docker/kind 기반 테스트**도 실용적인 검증 경로로 고려한다.

---

## 5. 현재 열린 문제

- **ClusterLoader2 결과 중 수용 테스트에 정말 중요한 지표는 무엇인가?**
- **그 지표들이 왜 중요한가?** (사용자 관점, 운영 관점에서 정당화)
- **소규모 클러스터용 SLI/SLO를 어떻게 선정·정당화할 것인가?**
- **넓은 벤치마크 임계값을 무비판적으로 재사용하지 않으려면?**
- **Assertion을 어떤 형태로 표현할 것인가?** (예: JUnit XML, JSON, 스키마)
- **perfdash만으로는 부족한 부분을 어떻게 시각화할 것인가?**

이 문제들은 “연구만으로 끝나지 않고, 실험·해석·devcat 개선”으로 이어져야 한다.

---

## 6. 향후 시각화 방향

- **Evidence**를 탐색할 가치가 있다. 코드 기반 대시보드(code-based dashboards)를 제공할 가능성이 있기 때문이다.
- **향후 대시보드**는 perfdash처럼 **숫자 결과**를 보여 주되, **XML/테스트 결과 기반의 PASS/FAILED 표현**도 포함하는 형태를 염두에 둘 수 있다.
- Evidence는 **가능한 방향 중 하나**일 뿐이며, **최종 확정 선택이 아니다.** 다른 옵션과 비교·실험 후 결정한다.

---

## 7. 향후 작업의 작업 모델

- **CAT 관련 태스크는 연구에서 멈추지 않는다.** 아래 흐름으로 연결한다:

  **research → experiment → interpretation → devcat improvement**

- **제안**은 devcat의 **현재 현실**(ClusterLoader2, config.yaml, ol-test.yaml, results/, perfdash)을 존중하고, **점진적으로 개선**하는 방향이어야 한다. 한 번에 완전히 새로 짜는 것보다, “다음에 할 수 있는 한 단계”가 명확한 것이 좋다.
- AI 팀은 Manager/Researcher/Architect/Engineer 단계에서 이 브리프와 cat-vision, cat-design-principles, sli-slo-philosophy를 참고해, **실제 devcat 진화**에 기여하는 산출물을 낸다.

---

## 요약

| 항목 | 내용 |
|------|------|
| **책임** | 실용 CAT 시스템, Conformance + Performance; 현재는 연구·방향·Performance/부하 초점. |
| **방향** | 다중 도구 조합, Job(scenario + SLI/SLO + assertion), 결과 디렉터리·시각화, 단순 오케스트레이션. |
| **devcat 현실** | ClusterLoader2, perfdash, config.yaml + ol-test.yaml, results/, assertion/PASS·FAIL 미명확. |
| **원하는 것** | 실행·검사·숙련·devcat 진화, AI 팀은 모호→구체, kind/Docker 검증 경로 고려. |
| **열린 문제** | 중요 지표·정당화, 소규모 SLO, assertion 표현, 시각화 개선. |
| **시각화** | Evidence 등 코드 기반 대시보드·PASS/FAIL 표현 탐색; 최종 선택 아님. |
| **작업 모델** | research → experiment → interpretation → devcat improvement, 점진적 개선. |
