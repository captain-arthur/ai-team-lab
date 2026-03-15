# Engineering: Kubernetes CAT 프레임워크 구현 가이드

**Project:** example-kubernetes-cat  
**Role:** Engineer  
**Input:** 01-manager, 02-research, 03-architecture

이 디렉터리는 **구현 가이드 및 예제**만 제공합니다. 프로덕션 수준의 완전한 구현이 아니라, Architecture에 맞춰 CAT를 도입할 때 참고할 runbook·디렉터리 규약·스크립트 예시입니다.

---

## 구현 계획 (Implementation plan)

| 순서 | 단계 | 설명 | 의존성 |
|------|------|------|--------|
| 1 | 사전 요건 정리 | Sonobuoy 설치, kubeconfig, 결과 디렉터리 규약 확인 | — |
| 2 | Sonobuoy 실행 워크플로 적용 | run → wait → retrieve → 결과 저장 경로 맞추기 | 1 |
| 3 | 커스텀 테스트 구조 정의 | 스크립트 또는 KUTTL 디렉터리 구조, exit code·로그 규약 | 1 |
| 4 | 한 번의 CAT run 오케스트레이션 | Sonobuoy → 커스텀 → 수집 → 리포팅 순서를 runbook/스크립트로 고정 | 2, 3 |
| 5 | 결과 수집·리포팅 규약 | `results/cat/<run-id>/` 아래 sonobuoy/, custom/, report.md 구조 | 4 |

---

## 산출물 목록

| 파일 | 용도 |
|------|------|
| `runbook.md` | 예제 CAT 실행 runbook (단계별 절차) |
| `sonobuoy-workflow.md` | Sonobuoy 실행 워크플로 및 명령 예시 |
| `custom-test-structure.md` | 커스텀 테스트 디렉터리·스크립트 구조 예시 |
| `result-structure.md` | 결과 디렉터리 규약 및 예시 트리 |
| `run-cat-example.sh` | 예제 CAT 실행 스크립트 (Sonobuoy + 커스텀 + 수집·리포팅 흐름) |

---

## 실행 방법 요약

- **Runbook:** `runbook.md`를 따라 수동으로 단계 실행.
- **스크립트:** `run-cat-example.sh`는 예제 흐름만 보여 주며, 실제 환경에 맞게 `KUBECONFIG`, `RESULTS_DIR`, Sonobuoy 모드 등을 수정한 뒤 사용.
- **사전 요건:** Sonobuoy CLI 설치, 대상 클러스터에 대한 kubeconfig, 쓸 수 있는 로컬 디렉터리(결과 저장용).

---

## 구현 노트 (Implementation notes)

- **가정:** Kubernetes 1.24+, Sonobuoy v0.56.x 호환 가정. 클러스터에 Sonobuoy가 네임스페이스·리소스를 생성할 수 있는 권한 필요.
- **Architecture 준수:** Conformance runner(Sonobuoy) + Custom runner(스크립트 예시) + Result collector/Reporter 규약을 03-architecture와 동일하게 맞춤.
- **TODO/한계:** 실제 팀 환경에서는 커스텀 테스트 목록 확장, CI 연동, 보관 정책(로그 로테이션)을 추가로 정의해야 함. Reporter의 자동 요약 생성(예: Sonobuoy JUnit 파싱 → report.md)은 별도 스크립트/도구로 구현 가능하며, 본 예제에서는 수동 작성 예시만 포함.
- **이탈 사항:** 없음. Architecture의 결과 경로·실행 순서를 그대로 반영함.
