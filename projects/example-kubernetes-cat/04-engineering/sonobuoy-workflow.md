# Sonobuoy 실행 워크플로 예시

Architecture의 **Conformance runner**를 Sonobuoy로 수행할 때의 명령 흐름과 옵션입니다.

---

## 워크플로 개요

1. `sonobuoy run` — 클러스터에 테스트 Pod 스케줄, 실행 시작.
2. `sonobuoy status` / `sonobuoy wait` — 진행 확인 및 완료 대기.
3. `sonobuoy retrieve` — 결과 tarball을 로컬로 가져오기.
4. `sonobuoy results <tarball>` — 요약 확인 (선택).
5. tarball을 결과 수집 규약 경로에 저장 (예: `results/cat/<run-id>/sonobuoy/`).

---

## 모드 선택

| 모드 | 명령 | 용도 | 소요 시간 대략 |
|------|------|------|------------------|
| Quick | `sonobuoy run --mode quick` | 빠른 게이트, CI | 수 분 |
| Conformance | `sonobuoy run --mode conformance` | 전체 conformance | 1–2시간 |
| Certified-conformance | `sonobuoy run --mode certified-conformance` | CNCF 인증 제출용 | 1–2시간 |

- 프로덕션 전 “빠른 검증”에는 `quick`; 주기적 전체 검증에는 `conformance` 권장.

---

## 명령 예시 (복사용)

```bash
# 1) 실행 시작 (Quick)
sonobuoy run --mode quick

# 2) 완료 대기 (필수)
sonobuoy wait

# 3) 결과 가져오기 (현재 디렉터리에 .tar.gz 생성)
sonobuoy retrieve

# 4) 생성된 tarball을 규약 경로로 이동 (RUN_ID, RESULTS_ROOT는 runbook에서 정의)
# mv sonobuoy_*.tar.gz "${RESULTS_ROOT}/sonobuoy/"

# 5) 요약 확인 (선택)
sonobuoy results sonobuoy_*.tar.gz
```

---

## kubeconfig 및 네임스페이스

- Sonobuoy는 현재 `KUBECONFIG`(또는 `--kubeconfig`)가 가리키는 클러스터에 테스트를 실행합니다.
- 테스트는 임시 네임스페이스 등을 사용하며, 완료 후 정리됩니다. 실행 전 `kubectl get nodes` 등으로 클러스터 접근을 확인하세요.

---

## 버전 호환

- 사용할 Sonobuoy 버전과 클러스터 Kubernetes 버전을 runbook 또는 팀 문서에 명시하는 것을 권장합니다 (Architecture 리스크 완화). 예: Sonobuoy v0.56.x, Kubernetes 1.24+.
