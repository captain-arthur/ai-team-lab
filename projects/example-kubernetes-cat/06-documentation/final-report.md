# Final Report: Kubernetes Cluster Acceptance Test (CAT) 프레임워크

**Project:** example-kubernetes-cat  
**Role:** Writer  
**Input:** 01-manager ~ 05-review  
**Date:** 2025-03

---

## Executive summary

본 프로젝트는 **Kubernetes Cluster Acceptance Test(CAT) 프레임워크**의 설계·권장안·실행 가이드를 산출했다. 목표는 “프로덕션 승격 전 클러스터 검증”을 표준 conformance와 팀 정의 항목으로 반복 가능하게 하는 것이다. Research에서 **Sonobuoy + 소규모 커스텀 레이어(스크립트/KUTTL)** 하이브리드를 권장했고, Architecture에서 테스트 카테고리·실행 모델·결과 수집·리포팅을 정의했으며, Engineering에서 runbook·예제 스크립트·디렉터리 규약을 제공했다. Review 결과 **Pass with comments**이며, 실사용 전 버전 명시·Reporter 자동화·커스텀 예제 보강 등 소규모 보완을 권장한다.

---

## 1. CAT 프레임워크 개요

CAT 프레임워크는 **한 번의 run**으로 다음을 수행한다.

- **Conformance:** Kubernetes 공식 conformance 테스트(Sonobuoy) 실행.
- **Custom:** 팀 정의 검증(노드 Ready, DNS, StorageClass, 정책 등)을 스크립트 또는 KUTTL로 실행.
- **수집:** Sonobuoy tarball과 커스텀 로그를 한 디렉터리 트리(`results/cat/<run-id>/`)에 모음.
- **리포팅:** 요약 리포트(`report.md`)로 Conformance·Custom 결과와 **Overall pass/fail**을 제공.

**Overall passed**는 Conformance 성공 **및** Custom 성공일 때만 인정한다. 플랫폼/SRE는 이 결과를 보고 프로덕션 승격 여부를 판단한다.

---

## 2. 목적 및 범위

### 목적

- 프로덕션으로 올리기 전에 클러스터가 **기대대로 구성·동작하는지** 검증할 수 있게 한다.
- 테스트 정의 방식, 실행 순서, 결과 수집·리포팅 방식을 **일관된 규약**으로 둔다.
- 팀이 **기존 OSS 도구(Sonobuoy) + 필요한 만큼의 커스텀**으로 도입할 수 있게 한다.

### 범위 (in scope)

- 프레임워크 설계: 테스트 카테고리, 실행 모델, 결과 수집, 리포팅.
- 기존 도구 비교 및 권장안(Sonobuoy + 커스텀 레이어).
- 구현 가이드: runbook, Sonobuoy 워크플로, 커스텀 테스트 구조·결과 디렉터리 규약, 예제 스크립트.

### 범위 (out of scope)

- 전체 프로덕션 수준 구현 코드, 클러스터 프로비저닝/배포, CI/CD 상세 설계. 표준 Kubernetes 1.24+ 가정; proprietary 확장 미포함.

### 성공 기준 (Manager 기준)

- 설계 문서로 테스트 유형·흐름·도구 옵션이 정의되어 있음. ✅  
- “기존 도구 vs 소규모 커스텀”에 대한 명확한 권장과 근거가 있음. ✅  
- 플랫폼/SRE가 참고해 도입할 수 있는 짧은 보고서·runbook·예제가 있음. ✅  

---

## 3. 아키텍처 요약

### 테스트 카테고리

| 카테고리 | 검증 목적 | 수행 방식 |
|----------|-----------|-----------|
| Conformance | Kubernetes API·핵심 동작 표준 준수 | Sonobuoy (quick / conformance / certified-conformance) |
| API & Control plane | API server 가용성, 기본 CRUD | Sonobuoy e2e + 필요 시 kubectl 스크립트 |
| Nodes & Scheduling | 노드 Ready, 스케줄링 | Sonobuoy e2e + 선택적 노드 체크 스크립트 |
| DNS & Network | CoreDNS, 서비스 디스커버리 | Sonobuoy 또는 커스텀 Job/Pod |
| Storage | PVC/PV, StorageClass | Sonobuoy e2e + 팀 정의 스크립트/KUTTL |
| Custom / Team | 팀 정책, CRD, 내부 요구사항 | 스크립트, KUTTL, 또는 Sonobuoy 플러그인 |

### 구성 요소

- **Conformance runner:** Sonobuoy CLI. run → wait → retrieve → tarball을 규약 경로에 저장.
- **Custom test runner:** 스크립트 또는 KUTTL. exit code·로그를 규약 경로에 저장.
- **Result collector:** `results/cat/<run-id>/` 아래 sonobuoy/, custom/로 한 run당 한 트리 유지.
- **Reporter:** 수집 결과를 바탕으로 `report.md` 생성(요약, Conformance·Custom pass/fail, Overall).

### 실행·데이터 흐름

1. Run ID·결과 디렉터리 준비 → 2. Sonobuoy 실행(모드 선택) → 3. 커스텀 테스트 실행 → 4. 수집(이미 동일 트리) → 5. report.md 작성 → 6. 결과 검토·판단.

### 주요 설계 결정

- **Conformance 도구:** Sonobuoy 기반(Research 권장). Hydrophone은 경량 대안으로 유지.
- **팀 정의 테스트:** Sonobuoy와 별도로 스크립트/KUTTL 사용(유지보수·확장 용이).
- **결과 위치:** 로컬/CI 단일 디렉터리 규약(`results/cat/<run-id>/`).
- **전체 pass 기준:** Conformance 성공 **및** Custom 성공.
- **리포팅:** 요약 리포트(Markdown) + Sonobuoy JUnit 유지.

---

## 4. 프레임워크 실행 방법

### 사전 요건

- 대상 클러스터에 대한 **KUBECONFIG** 설정(`kubectl get nodes` 등으로 확인).
- **Sonobuoy CLI** 설치(예: [sonobuoy.io](https://sonobuoy.io), `brew install sonobuoy`).
- 결과를 저장할 로컬 디렉터리 생성 권한(예: `./results`).
- (선택) 커스텀 테스트용 스크립트 또는 KUTTL 준비.

### 권장 환경 (Review 반영)

- **Kubernetes:** 1.24+ (표준 API 기준).
- **Sonobuoy:** v0.56.x 또는 호환 버전. 사용 전 `sonobuoy version` 확인 권장.

### 실행 옵션

- **Runbook으로 수동 실행:** `04-engineering/runbook.md`의 Step 1~6을 순서대로 수행. Run ID 생성, Sonobuoy run/wait/retrieve, 커스텀 테스트, 수집, report.md 작성.
- **예제 스크립트:** `04-engineering/run-cat-example.sh` 사용. 인자로 `quick` 또는 `conformance` 전달.  
  - 사용 전 `KUBECONFIG`, `RESULTS_BASE` 등 환경에 맞게 수정.  
  - 커스텀 테스트가 없으면 해당 단계는 skip되고 Custom exit = 0으로 처리됨.

```bash
# 예: Quick 모드로 한 번 실행
./run-cat-example.sh quick
```

- **Sonobuoy만 단독 실행:** `04-engineering/sonobuoy-workflow.md` 참고. run → wait → retrieve → tarball을 규약 경로에 저장.

---

## 5. 예시 워크플로

1. **준비:** `export RUN_ID=$(date +%Y%m%d-%H%M%S)` 및 `results/cat/${RUN_ID}/sonobuoy`, `custom` 디렉터리 생성.
2. **Sonobuoy:** `sonobuoy run --mode quick` → `sonobuoy wait` → `sonobuoy retrieve` → tarball을 `results/cat/${RUN_ID}/sonobuoy/`로 이동.
3. **커스텀:** `./custom-tests/run-custom-tests.sh results/cat/${RUN_ID}/custom` 실행(해당 스크립트가 있을 경우). 로그·exit code를 custom/ 아래에 저장.
4. **리포팅:** 수집된 Sonobuoy 결과·커스텀 로그를 바탕으로 `results/cat/${RUN_ID}/report.md` 작성. Conformance·Custom·Overall pass/fail 기록.
5. **판단:** report.md와 로그를 보고 프로덕션 승격 여부 또는 재검증 대상 결정.

커스텀 테스트가 없으면 3단계 생략, Custom = pass로 간주하면 된다.

---

## 6. 기대 산출물

- **디렉터리:** `results/cat/<run-id>/`
  - `sonobuoy/sonobuoy_<run-id>.tar.gz` (및 선택적 압축 해제 내용)
  - `custom/*.log`, `custom/summary.txt` (및 선택적 KUTTL 아티팩트)
  - `report.md`
- **report.md 내용:** Run ID, Conformance 결과 요약, Custom exit code·로그 경로, **Overall: pass / fail**.
- **Overall pass:** Conformance 성공 **및** Custom 성공일 때만 pass.

---

## 7. 알려진 한계 및 제약

- **Reporter 자동 판단:** Conformance 성공 여부는 현재 Sonobuoy 출력을 사람이 해석하는 전제다. Overall pass를 자동으로 만들려면 `sonobuoy results` 또는 JUnit 파싱 스크립트가 추가로 필요하다.
- **Sonobuoy 실패 시:** `run-cat-example.sh`는 `set -e`로 Sonobuoy 단계 실패 시 즉시 종료된다. Custom 단계·report 생성은 수행되지 않는다. 실패 시에는 runbook에 따라 수동으로 partial 결과·report에 fail을 기록하는 절차를 따르는 것이 좋다.
- **커스텀 테스트 예제:** 구조·규약만 문서화되어 있고, 실행 가능한 `custom-tests/run-custom-tests.sh` 등 실제 파일은 04-engineering에 포함되어 있지 않다. 도입 시 팀이 custom-test-structure.md를 참고해 최소 한두 개 스크립트를 추가하는 것을 권장한다.
- **버전 호환:** Sonobuoy와 Kubernetes 버전 불일치 시 테스트 실패 가능. 권장 버전(Kubernetes 1.24+, Sonobuoy v0.56.x)을 팀 문서에 명시하고, 필요 시 runbook에 버전 확인 단계를 넣는 것이 좋다.
- **결과 보관:** 보관 기간·오래된 run 정리 정책은 팀 운영 정책으로 별도 정의한다. result-structure.md에 “예: 30일 후 삭제 또는 아카이브” 수준의 가이드만 있다.
- **커스텀 테스트 리소스 정리:** 테스트용 네임스페이스·리소스는 테스트 정의 시 생성·삭제 단계를 포함할 것을 권장한다. 구체 예시는 custom-test-structure.md 또는 runbook에 추가하면 Architecture에서 제기한 리스크를 완화할 수 있다.

---

## 8. Reviewer 의견 및 사용자 권장 사항

Reviewer는 전반적으로 **Pass with comments**로 평가했으며, 실사용 전 아래 사항을 반영할 것을 권장한다.

1. **버전 명시:** README 또는 runbook 상단에 “권장: Kubernetes 1.24+, Sonobuoy v0.56.x”를 적고, 필요 시 `sonobuoy version` / `kubectl version` 확인 단계를 runbook에 추가.
2. **Reporter 자동화:** Conformance pass/fail을 자동 판단해 Overall을 쓰는 작은 스크립트를 추가하거나, runbook에 “자동화 시 JUnit/sonobuoy results 파싱 필요”를 명시.
3. **커스텀 테스트 최소 예제:** `run-custom-tests.sh`, `01-nodes-ready.sh` 수준의 실제 파일을 두어 한 번에 돌려 볼 수 있게 함.
4. **Runbook 보강:** (1) 커스텀 테스트가 없을 때 Step 3 생략·Custom = pass 명시. (2) Sonobuoy 실패 시 partial 결과 저장·수동 report 작성 절차 한 줄 추가.
5. **커스텀 테스트 정리:** 테스트용 네임스페이스 생성·삭제 예시를 runbook 또는 custom-test-structure에 추가.
6. **보관 정책:** 보관 기간·오래된 `results/cat/*` 정리 방법을 result-structure 또는 운영 가이드에 한 줄이라도 추가.

위 항목은 모두 **Minor**이며, 당장 차단되는 사항은 없다. Writer 단계에서 “한계 및 follow-up”으로 문서화했고, 필요 시 Engineer 보완으로 이어가면 된다.

---

## 9. 산출물 및 참조

### 프로젝트 내 산출물

| 단계 | 경로 | 주요 내용 |
|------|------|-----------|
| Manager | `01-manager/project-brief.md` | 목표, 범위, 워크 브레이크다운, 역할별 handoff |
| Research | `02-research/research-notes.md` | 도구 비교, Sonobuoy·Hydrophone·KUTTL, 권장안 |
| Architecture | `03-architecture/architecture.md` | 테스트 카테고리, 구성 요소, 실행·수집·리포팅, 설계 결정 |
| Engineering | `04-engineering/` | README, runbook, sonobuoy-workflow, custom-test-structure, result-structure, run-cat-example.sh |
| Review | `05-review/review-notes.md` | 검토 요약, 체크리스트, 이슈·권장 사항 |

### 사용자 문서로 참고할 파일

- **실행 절차:** `04-engineering/runbook.md`
- **Sonobuoy만 실행:** `04-engineering/sonobuoy-workflow.md`
- **커스텀 테스트 구조:** `04-engineering/custom-test-structure.md`
- **결과 디렉터리 규약:** `04-engineering/result-structure.md`

---

## 10. Follow-up 및 권장 다음 단계

- **팀 도입 시:** 위 “알려진 한계”와 “Reviewer 의견”을 반영해 버전 명시·runbook 보강·커스텀 예제 1~2개·(선택) Reporter 파싱 스크립트를 추가.
- **CI 연동:** 본 설계는 CI/CD 상세를 scope 밖으로 두었으나, runbook·스크립트·결과 규약을 그대로 활용해 CI job에서 CAT run을 호출하고, `report.md` 또는 exit code로 게이트를 거는 방식으로 확장 가능.
- **Knowledge Extraction:** 프로젝트 완료 후 `07-knowledge-extraction` 단계에서 CAT 설계 패턴·Sonobuoy/커스텀 하이브리드·클러스터 검증 교훈을 `knowledge/`에 정리할 수 있다.

---

*Writer 역할 산출물. Knowledge Extraction 단계는 별도로 진행.*
