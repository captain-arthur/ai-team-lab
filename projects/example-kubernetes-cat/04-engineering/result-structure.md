# 결과 디렉터리 구조 (Result collection 규약)

Architecture의 **Result collector** 및 **Reporter**와 맞추기 위한 결과 디렉터리 규약 예시입니다.

---

## 한 Run당 단일 트리

한 번의 CAT run은 하나의 `run-id`로 구분하고, 해당 run의 모든 결과를 한 디렉터리 트리 아래에 둡니다.

```
results/
└── cat/
    └── <run-id>/                    # 예: 20250315-143022 또는 CI build-id
        ├── sonobuoy/                # Conformance runner 산출물
        │   ├── sonobuoy_<run-id>.tar.gz
        │   └── (선택) 압축 해제한 내용: plugins/ e2e.log 등
        ├── custom/                  # Custom test runner 산출물
        │   ├── 01-nodes-ready.log
        │   ├── 02-dns.log
        │   ├── 03-storageclass.log
        │   ├── summary.txt          # exit code 요약 등
        │   └── (선택) kuttl/        # KUTTL artifacts
        └── report.md                # Reporter 산출: 요약 리포트
```

---

## 디렉터리별 용도

| 경로 | 담당 | 내용 |
|------|------|------|
| `results/cat/<run-id>/sonobuoy/` | Conformance runner | Sonobuoy `retrieve` tarball; 필요 시 압축 해제한 JUnit, e2e.log |
| `results/cat/<run-id>/custom/` | Custom test runner | 팀 정의 테스트 로그, summary, (선택) KUTTL 아티팩트 |
| `results/cat/<run-id>/report.md` | Reporter | 카테고리별 pass/fail, 전체 CAT 결과, 로그 경로 요약 |

---

## run-id 규약 예시

- **수동 실행:** `date +%Y%m%d-%H%M%S` (예: `20250315-143022`).
- **CI:** 빌드 ID 또는 job run 번호.
- **규칙:** 동일 run-id 아래에 해당 run의 Sonobuoy·Custom·report만 두면 됨.

---

## 보관 정책

- Architecture에 따라 보관 기간·아카이브는 운영 정책으로 별도 정의합니다. 예: 30일 후 삭제, 또는 오래된 `results/cat/<run-id>/` 디렉터리를 주기적으로 압축·이동하는 스크립트를 두는 방식.
