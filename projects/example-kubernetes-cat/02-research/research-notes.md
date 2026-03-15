# Research: Kubernetes Cluster Acceptance Test (CAT) 도구 및 접근법

**Project:** example-kubernetes-cat  
**Role:** Researcher  
**Input:** `01-manager/project-brief.md`  
**Date:** 2025-03

---

## Research questions (Manager handoff 기준)

1. 클러스터 acceptance/conformance 테스트를 지원하는 기존 도구·프로젝트는 무엇이 있는가? (sonobuoy, custom harness, OSS 등)
2. 사용 편의성, 테스트 정의 형식, 리포팅, “프로덕션 전 검증” 용도 적합성 측면에서 어떻게 비교되는가?
3. 기존 도구 사용 vs 얇은 커스텀 레이어(스크립트 + YAML) 간 트레이드오프는?

---

## Findings

### 1. Kubernetes 공식 Conformance 테스트

- **요지:** CNCF Certified Kubernetes Conformance Program에서 정의하는 표준은, Kubernetes e2e 테스트 스위트 중 `[Conformance]` 태그가 붙은 테스트들이다. GA 기능만 대상이며, 프로바이더 중립·비파괴적 실행 등 조건을 만족해야 한다.
- **출처:** [kubernetes/kubernetes - test/conformance](https://github.com/kubernetes/kubernetes/tree/master/test/conformance), [kubernetes/community - conformance-tests.md](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/conformance-tests.md), [cncf/k8s-conformance](https://github.com/cncf/k8s-conformance)
- **버전:** Kubernetes 1.24+ 기준 문서와 테스트 스위트 유지됨.

### 2. Sonobuoy

- **요지:** Kubernetes 클러스터에 표준화된 conformance/진단 테스트를 실행하는 도구. 테스트 어그리게이터 Pod를 배포해 노드 간 테스트를 조율하고, 공식 Kubernetes e2e 스위트를 실행한 뒤 결과를 tarball로 묶어 분석한다. 테스트 후 임시 네임스페이스·리소스는 정리되는 비파괴적 운영이 목표다.
- **테스트 모드:** Conformance(전체 CNCF conformance, 약 300개 테스트), Quick(핵심 테스트만, 수 분), Certified-conformance(인증 제출용). 기본은 `[Conformance]` 태그 테스트만 실행하고 disruptive 테스트는 제외.
- **사용 흐름:** `sonobuoy run --mode quick` → `sonobuoy status` / `sonobuoy wait` → `sonobuoy retrieve` → `sonobuoy results $results`. 전체 conformance는 보통 1–2시간 소요.
- **특징:** 플러그인 확장, 에어갭 클러스터·다중 Kubernetes 버전 지원. 결과는 JUnit XML 및 상세 로그 포함.
- **출처:** [sonobuoy.io](https://sonobuoy.io/), [Sonobuoy docs (e2eplugin, FAQ)](https://sonobuoy.io/docs/main/e2eplugin), [Understanding E2E Tests](https://sonobuoy.io/understanding-e2e-tests/)
- **버전:** v0.56.x (2024) 등 최신 릴리스 존재.

### 3. Hydrophone (kubernetes-sigs)

- **요지:** kubernetes-sigs의 **경량 conformance 테스트 러너**. Kubernetes 릴리스 팀이 배포하는 conformance 이미지를 사용한다. Sonobuoy보다 단순한 구조로, 공식 conformance 실행·제출 시 대안으로 언급됨.
- **출처:** [kubernetes-sigs/hydrophone](https://github.com/kubernetes-sigs/hydrophone), CNCF k8s-conformance instructions
- **용도:** “프로덕션 전 검증”보다는 공식 conformance 인증 제출에 가깝지만, 클러스터 검증 러너로 활용 가능.

### 4. KUTTL (Kubernetes Test TooL)

- **요지:** Kubernetes 리소스(선언적 YAML)로 테스트 케이스를 작성하는 **통합/이벤트 기반 테스트 하네스**. 오퍼레이터, Helm 차트, 앱 배포·동작 검증에 적합. mock control plane, kind, 실제 클러스터에서 실행 가능. 코드 없이 리소스만으로 e2e·통합 테스트 구성 가능. krew로 설치.
- **출처:** [kuttl.dev](https://kuttl.dev/docs/kuttl-test-harness.html), [kudobuilder/kuttl](https://github.com/kudobuilder/kuttl)
- **CAT와의 관계:** 공식 conformance 러너는 아님. “클러스터가 올바르게 설정되었는지”보다 “우리 앱/오퍼레이터가 이 클러스터에서 잘 동작하는지” 검증에 가깝다. 커스텀 acceptance 시나리오(예: DNS, StorageClass, 특정 API 동작)를 YAML로 정의할 때 보조 도구로 쓸 수 있음.

### 5. kubectl-validate (kubernetes-sigs)

- **요지:** SIG-CLI 산하 **로컬 리소스 검증 도구**. apiserver와 동일한 검증 로직을 사용해 매니페스트·CRD 검증. Kubernetes 1.23–1.27 네이티브 타입 및 클러스터/로컬 CRD 지원. kubectl 플러그인 또는 CLI로 실행.
- **출처:** [kubernetes-sigs/kubectl-validate](https://github.com/kubernetes-sigs/kubectl-validate)
- **CAT와의 관계:** 클러스터 “실행 시 검증”이 아니라 **선언적 리소스 문법·스키마 검증**에 특화. CAT 프레임워크에서 “배포 전 매니페스트 검증” 단계에 보조적으로 활용 가능.

### 6. 기존 도구 vs 얇은 커스텀 레이어

- **기존 도구(Sonobuoy/Hydrophone):** 공식 conformance·표준 스위트 재사용, CNCF 인증 경로와 일치, 유지보수 부담 적음. 반면 “우리만의” acceptance 항목(예: 내부 정책, 특정 StorageClass, DNS 네이밍)을 추가하려면 플러그인·래퍼가 필요.
- **얇은 커스텀 레이어(스크립트 + YAML):** kubectl, 간단한 스크립트, KUTTL 또는 단순 Job/테스트 Pod로 “API 가용성, 노드 준비, DNS, StorageClass 존재” 등만 검증. 구현·변경이 빠르고 팀 정의 체크리스트에 맞추기 쉽다. 대신 conformance 표준과 자동 동기화되지 않고, 테스트 품질·리포팅을 직접 설계해야 함.
- **하이브리드:** Sonobuoy(또는 Hydrophone)로 conformance를 돌리고, 별도 작은 스크립트/KUTTL로 팀별 acceptance 항목을 추가하는 방식이 현실적.

---

## Comparison

| 항목 | Sonobuoy | Hydrophone | KUTTL | 커스텀(스크립트+YAML) |
|------|----------|------------|-------|------------------------|
| **용도** | Conformance + 플러그인 확장 | 경량 conformance 러너 | 앱/오퍼레이터 통합 테스트 | 팀 정의 acceptance |
| **사용 난이도** | 중 (모드·플러그인 이해 필요) | 중하 (상대적 단순) | 중 (테스트 작성 패턴 학습) | 하 (kubectl 수준) |
| **테스트 정의** | 공식 e2e + 플러그인 | 공식 conformance 이미지 | Kubernetes 리소스(YAML) | 스크립트/Job/Helm 등 |
| **리포팅** | tarball, JUnit XML, 로그 | 제출용 출력 | 테스트별 성공/실패 | 직접 설계 |
| **프로덕션 전 검증 적합성** | 높음 (conformance + 선택 플러그인) | 높음 (conformance 중심) | 중 (시나리오별) | 중–높음 (요구사항에 따라) |
| **유지보수** | CNCF/커뮤니티 유지 | kubernetes-sigs | 커뮤니티 | 팀 전담 |
| **확장성** | 플러그인으로 커스텀 테스트 추가 가능 | conformance 위주 | 시나리오 무제한 추가 | 완전 자유 |

---

## Recommendation

- **1차 권장:** **Sonobuoy를 기반으로 하고, 팀별 acceptance 항목이 필요하면 소규모 커스텀 레이어(스크립트 또는 KUTTL 시나리오)를 추가하는 하이브리드**
  - 이유: 표준 conformance(API, 노드, 기본 동작)는 Sonobuoy로 커버하고, “우리 클러스터만의” 검증(특정 StorageClass, DNS, 정책)은 가벼운 스크립트나 KUTTL로 정의하면, 재사용·표준 준수와 유연성을 동시에 만족함.
- **대안:** 공식 conformance만 필요하고 Sonobuoy가 무겁다고 판단되면 **Hydrophone** 검토. “프로덕션 전 검증” 범위를 conformance + 소수 항목으로 한정할 때 적합.
- **커스텀만:** conformance 인증이 전혀 필요 없고, 검증 항목이 적고 단순하면 **스크립트 + YAML(또는 KUTTL)** 만으로도 가능. 장기적으로는 Sonobuoy 플러그인 또는 위 하이브리드로 통합하는 것을 권장.

---

## Summary for downstream roles (Architect / Engineer)

- **핵심 정리:** (1) 공식 CAT/conformance는 Kubernetes e2e의 `[Conformance]` 태그 + CNCF 프로그램으로 정의됨. (2) Sonobuoy가 표준 러너이며, Quick/Conformance/Certified-conformance 모드와 플러그인으로 확장 가능. (3) Hydrophone은 경량 대안. (4) KUTTL은 conformance 러너가 아니라 앱/오퍼레이터·커스텀 시나리오용. (5) “프로덕션 전 검증”에는 Sonobuoy(또는 Hydrophone) + 필요 시 소규모 커스텀 레이어가 적합.
- **오픈 포인트:** 팀에서 “반드시 통과해야 하는” 항목 목록(conformance만 vs 내부 정책·StorageClass·DNS 등)이 정리되어 있으면, Architect가 테스트 카테고리와 실행 모델을 더 구체화하기 쉬움. 에어갭·실행 시간 제약이 있으면 Sonobuoy Quick 모드 또는 Hydrophone 비중을 높이는 설계가 필요.
- **리스크/한계:** Sonobuoy 전체 conformance는 1–2시간 소요; CI/수동 “빠른 검증”에는 Quick 모드 또는 별도 경량 체크가 필요. Hydrophone은 문서·사례가 Sonobuoy보다 적을 수 있음.
