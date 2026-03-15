# 예제 CAT 실행 Runbook

프로덕션 승격 전 또는 정기 검증 시 한 번의 Cluster Acceptance Test(CAT) run을 수행하기 위한 단계별 절차입니다. Architecture의 실행 순서(Conformance → Custom → 수집 → 리포팅)를 따릅니다.

---

## 사전 요건

- 대상 클러스터에 대한 `KUBECONFIG` 설정 완료 (`kubectl get nodes` 등으로 확인).
- Sonobuoy CLI 설치 ([sonobuoy.io](https://sonobuoy.io) 또는 `brew install sonobuoy` 등).
- 결과를 저장할 로컬 디렉터리 생성 권한 (예: `./results`).
- (선택) 커스텀 테스트용 스크립트 또는 KUTTL 테스트가 준비되어 있음.

---

## Step 1: Run ID 및 결과 디렉터리 준비

- 이번 run을 구분할 ID 생성. 예: `$(date +%Y%m%d-%H%M%S)` 또는 CI 빌드 ID.
- 결과 루트 디렉터리 생성. Architecture 규약 예: `./results/cat/<run-id>/`.

```bash
export RUN_ID=$(date +%Y%m%d-%H%M%S)
export RESULTS_ROOT="./results/cat/${RUN_ID}"
mkdir -p "${RESULTS_ROOT}"/{sonobuoy,custom}
```

---

## Step 2: Sonobuoy 실행 (Conformance)

- **Quick 모드**(수 분): 빠른 게이트용.
- **Conformance 모드**(1–2시간): 전체 conformance 스위트.

**2.1** Sonobuoy run 시작 (Quick 예시)

```bash
sonobuoy run --mode quick
```

**2.2** 진행 상황 확인 (선택)

```bash
sonobuoy status
```

**2.3** 완료 대기

```bash
sonobuoy wait
```

**2.4** 결과 tarball 가져오기

```bash
out=$(sonobuoy retrieve)
mv "$out" "${RESULTS_ROOT}/sonobuoy/sonobuoy_${RUN_ID}.tar.gz"
```

**2.5** (선택) Sonobuoy 결과 요약 확인

```bash
sonobuoy results "${RESULTS_ROOT}/sonobuoy/sonobuoy_${RUN_ID}.tar.gz"
```

---

## Step 3: 커스텀 테스트 실행

- 팀 정의 테스트(StorageClass, DNS, 노드 조건 등)를 스크립트 또는 KUTTL로 실행.
- 출력·로그를 `RESULTS_ROOT/custom/` 아래에 규약대로 저장. exit code는 스크립트/runbook에서 기록.

**3.1** 예: 스크립트 기반 커스텀 테스트

```bash
./custom-tests/run-custom-tests.sh "${RESULTS_ROOT}/custom" || true
```

- 실패 시에도 수집·리포팅을 위해 `|| true`로 exit code만 기록하고 계속 진행할 수 있음. 최종 pass/fail 판단은 Reporter 단계에서 Conformance + Custom 둘 다 반영.

**3.2** 예: KUTTL 실행 시

```bash
kubectl kuttl test --report json --artifacts-dir "${RESULTS_ROOT}/custom/kuttl" ./custom-tests/kuttl/ || true
```

- `--artifacts-dir`로 결과를 규약 경로에 저장.

---

## Step 4: 결과 수집 (Collector)

- Step 2에서 이미 `RESULTS_ROOT/sonobuoy/`에 tarball 저장됨.
- Step 3에서 `RESULTS_ROOT/custom/`에 로그·아티팩트 저장됨.
- 추가로 tarball 압축 해제하여 JUnit 등 중간 산출물을 같은 run 아래 두고 싶다면:

```bash
tar -xzf "${RESULTS_ROOT}/sonobuoy/sonobuoy_${RUN_ID}.tar.gz" -C "${RESULTS_ROOT}/sonobuoy/"
```

- 한 run당 하나의 `RESULTS_ROOT` 트리로 정리되는지 확인.

---

## Step 5: 리포팅 (Reporter)

- 수집된 결과를 바탕으로 요약 리포트 작성. 위치: `RESULTS_ROOT/report.md`.
- Conformance: Sonobuoy 결과 요약(또는 JUnit 파싱). Custom: 커스텀 실행 로그·exit code 요약.
- **전체 CAT passed** 조건: Conformance 성공 **및** Custom 성공.

**5.1** 예: 수동 report.md 초안

- `report.md`에 다음 형식으로 작성 예시:

```markdown
# CAT Run Report: <run-id>
- Conformance: pass/fail (Sonobuoy: <path>)
- Custom: pass/fail (logs: custom/)
- **Overall: pass / fail**
```

- 자동화하려면 Sonobuoy `results` 출력 파싱 + 커스텀 로그/exit code를 읽는 작은 스크립트를 별도로 두면 됨 (본 runbook에서는 수동 예시만 제시).

---

## Step 6: 정리 및 판단

- `RESULTS_ROOT` 경로와 `report.md`를 팀 규약(아티팩트 저장소 등)에 맞게 보관.
- report.md의 **Overall** 결과와 실패 항목을 보고 프로덕션 승격 여부 또는 재검증 대상 결정.

---

## 롤백/재실행

- Sonobuoy는 비파괴적 실행이며, 실행 후 리소스 정리됨. 커스텀 테스트는 네임스페이스 격리·삭제 단계를 테스트 정의에 포함할 것.
- 재실행 시 새 `RUN_ID`로 위 단계를 다시 수행하면 됨.
