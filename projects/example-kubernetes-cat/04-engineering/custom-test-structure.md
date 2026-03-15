# 커스텀 테스트 구조 예시

Architecture의 **Custom test runner**를 스크립트 또는 KUTTL로 구현할 때의 디렉터리·규약 예시입니다. 팀 정의 acceptance(StorageClass, DNS, 노드 조건 등)를 이 구조에 맞춰 추가하면 Result collector·Reporter와 일관되게 연동할 수 있습니다.

---

## 디렉터리 구조 예시

```
custom-tests/
├── run-custom-tests.sh      # 진입점: 모든 커스텀 테스트 실행, 결과를 인자로 받은 디렉터리에 저장
├── 01-nodes-ready.sh        # 예: 노드 Ready 상태 체크
├── 02-dns.sh                # 예: CoreDNS 동작 체크 (nslookup 또는 Pod 기반)
├── 03-storageclass.sh       # 예: 필수 StorageClass 존재 여부
└── kuttl/                   # (선택) KUTTL 시나리오
    ├── tests/
    │   └── ...
    └── kuttl-test.yaml
```

- 실행 시 **출력 디렉터리**를 인자로 받아, 그 아래에 로그·요약 파일을 씁니다.

---

## 규약 (Result collector·Reporter와의 계약)

1. **출력 위치:** 호출자가 지정한 디렉터리(예: `results/cat/<run-id>/custom/`).
2. **파일:** 각 스크립트 또는 KUTTL run별로 로그 파일(예: `01-nodes-ready.log`, `02-dns.log`) 및 공통 `summary.txt`(exit code 요약) 등.
3. **Exit code:** 전체 커스텀 run의 성공 = 모든 하위 테스트 성공 시 0, 하나라도 실패 시 비 zero. 호출 스크립트가 이 exit code를 기록해 두면 Reporter가 “Custom: pass/fail” 판단에 사용.

---

## 스크립트 예시 (개념)

`run-custom-tests.sh` 의사 코드:

```bash
# Usage: ./run-custom-tests.sh <output_dir>
# output_dir 예: results/cat/<run-id>/custom
OUT_DIR="$1"
mkdir -p "$OUT_DIR"
FAIL=0

run_one() { local name=$1; shift; "$@" > "${OUT_DIR}/${name}.log" 2>&1 || FAIL=1; }
run_one "01-nodes-ready" ./01-nodes-ready.sh
run_one "02-dns"         ./02-dns.sh
run_one "03-storageclass" ./03-storageclass.sh

echo "custom_tests_exit_code=$FAIL" > "${OUT_DIR}/summary.txt"
exit $FAIL
```

- 개별 테스트 스크립트(예: `01-nodes-ready.sh`)는 `kubectl get nodes`, `kubectl get storageclass` 등으로 조건을 검사하고 성공 시 0, 실패 시 비 zero를 반환하도록 구현하면 됩니다.

---

## 개별 테스트 예시 (01-nodes-ready.sh 개념)

- **목적:** 모든 노드가 Ready인지 확인.
- **방식:** `kubectl get nodes`로 NotReady가 있으면 실패.

```bash
#!/usr/bin/env bash
# 01-nodes-ready.sh: 모든 노드가 Ready인지 확인
not_ready=$(kubectl get nodes --no-headers 2>/dev/null | grep -v " Ready " | wc -l)
[[ "$not_ready" -eq 0 ]] && exit 0 || exit 1
```

- 실제 환경에서는 JSON 출력 파싱, 타임아웃 등 보강 가능. 여기서는 구조만 예시.

---

## KUTTL 사용 시

- KUTTL 테스트는 `kubectl kuttl test --artifacts-dir <output_dir>` 로 실행하고, `output_dir`를 `results/cat/<run-id>/custom/kuttl` 등 규약 경로로 두면 됩니다.
- KUTTL의 JSON 리포트나 로그를 Reporter가 읽어 “Custom” 섹션에 반영하도록 팀에서 정하면 됩니다.

---

## 정리

- 커스텀 테스트는 **네임스페이스 격리**와 **실행 후 정리**(리소스 삭제)를 포함하는 것이 좋습니다 (Architecture 리스크 완화). Job/Pod를 쓰는 경우 테스트용 네임스페이스를 만들고, 테스트 종료 후 해당 네임스페이스 삭제하는 단계를 runbook에 넣을 수 있습니다.
