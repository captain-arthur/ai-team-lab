# CAT v1 명세 (Cluster Acceptance Test v1 Specification)

**Project:** sample-cl2-sli-slo-analysis  
**Input:** metric-baseline-analysis, sli-candidate-analysis, cat-sli-slo-hypothesis, slo-validation-analysis, latency-variance-analysis

**목적:** devcat에서 **즉시 사용 가능한** 첫 CAT 명세를 정의한다. Stable SLI(프로덕션 수용 기준)와 Provisional SLI(임시 임계값, 추후 정제)를 구분하고, 판정 로직과 산출 형식을 명시한다.

---

## 1. CAT v1 SLI 세트

### 1.1 구분: Stable vs Provisional

| 구분 | 의미 |
|------|------|
| **Stable SLI** | 검증 실험에서 **변동이 없고**, 수용 기준으로 신뢰할 수 있는 지표. **프로덕션 수용 기준(production-ready acceptance criteria)** 로 사용. |
| **Provisional SLI** | **임시 임계값**을 두고, 환경·실험이 쌓이면 **정제(refinement)** 할 지표. 현재는 로컬/kind 실험용 또는 “회귀 방지” 용도. |

### 1.2 Stable SLI (production-ready)

| SLI | metric source | measurement location | 단위 |
|-----|---------------|----------------------|------|
| **ClusterOOMsTracker (failures)** | ClusterOOMsTracker | `ClusterOOMsTracker_load_*.json` → `failures` 배열 길이 | 개수 |
| **SystemPodMetrics (시스템 파드 restartCount)** | SystemPodMetrics | `SystemPodMetrics_load_*.json` → 모든 `pods[].containers[].restartCount` 중 최댓값 | 개수 |

- 4회 검증 run에서 **항상 0**으로 관측되었고, “클러스터 안정성·워크로드 건강도”의 기준점으로 적합하다고 검증됨.

### 1.3 Provisional SLI (임시, 정제 대상)

| SLI | metric source | measurement location | 단위 |
|-----|---------------|----------------------|------|
| **PodStartupLatency (pod_startup, P99)** | CreatePhasePodStartupLatency | `PodStartupLatency_CreatePhasePodStartupLatency_load_*.json` → `dataItems[].labels.Metric == "pod_startup"` 인 항목의 `data.Perc99` | ms |
| **StatelessPodStartupLatency (pod_startup, P99)** | CreatePhasePodStartupLatency | `StatelessPodStartupLatency_CreatePhasePodStartupLatency_load_*.json` → 동일하게 `pod_startup`, `Perc99` | ms |

- 레이턴시 변동성이 크고, 원인 분리·환경 개선 후 **임계값을 조정**할 예정. 현재는 **실험 베이스라인 기반 임시 값**만 부여함.

### 1.4 v1에서 평가 제외 (N/A 처리)

| SLI | 비고 |
|-----|------|
| **APIResponsivenessPrometheus** | Prometheus 활성화 run에서만 수집. 수집되면 추후 명세에 포함 가능. v1에서는 **평가 제외** (N/A). |

---

## 2. SLO 임계값 정의

### 2.1 Stable SLI — 엄격한 수용 기준

| SLI | SLO 조건 | 판정 |
|-----|----------|------|
| **ClusterOOMsTracker failures** | **failures 개수 = 0** | 조건 미달 시 **FAIL** (수용 불가). |
| **SystemPodMetrics (max restartCount)** | **모든 시스템 파드의 restartCount 최댓값 = 0** | 조건 미달 시 **FAIL** (수용 불가). |

- **의미:** OOM 1건이라도 발생하거나, 시스템 파드가 1회라도 재시작하면 클러스터 수용 실패로 간주.

### 2.2 Provisional SLI — 실험 베이스라인 기반 임시 임계값

**⚠️ 아래 값은 모두 provisional(임시)이며, 추후 실험·환경에 따라 조정된다.**

| SLI | Provisional SLO (임시) | 근거 |
|-----|-------------------------|------|
| **PodStartupLatency (pod_startup, P99)** | **P99 ≤ 130 s** (130 000 ms) | 4 run 검증에서 worst observed 122.3 s. 약간의 여유를 두어 130 s. **Provisional.** |
| **StatelessPodStartupLatency (pod_startup, P99)** | **P99 ≤ 75 s** (75 000 ms) | 4 run 검증에서 worst observed 69.9 s. 여유를 두어 75 s. **Provisional.** |

- **용도:** 현재 kind·로컬 실험에서 “심한 회귀”를 걸러내는 수준. **프로덕션 SLO로 고정된 것이 아님.**

### 2.3 임계값 요약표

| SLI | 구분 | 임계값 | 비고 |
|-----|------|--------|------|
| ClusterOOMsTracker failures | Stable | = 0 | 엄격 수용 기준 |
| SystemPodMetrics max restartCount | Stable | = 0 | 엄격 수용 기준 |
| PodStartupLatency P99 | **Provisional** | ≤ 130 s | 임시, 정제 대상 |
| StatelessPodStartupLatency P99 | **Provisional** | ≤ 75 s | 임시, 정제 대상 |

---

## 3. CAT 판정 로직 (Decision Logic)

### 3.1 SLI별 판정 규칙

```
1. ClusterOOMsTracker (failures)
   - failures 개수 > 0  →  FAIL (stable)
   - failures 개수 = 0  →  PASS

2. SystemPodMetrics (시스템 파드 restartCount)
   - max(restartCount) > 0  →  FAIL (stable)
   - max(restartCount) = 0  →  PASS

3. PodStartupLatency (pod_startup, P99) [Provisional]
   - 측정값(ms) > 130_000  →  FAIL (provisional)
   - 측정값(ms) ≤ 130_000  →  PASS

4. StatelessPodStartupLatency (pod_startup, P99) [Provisional]
   - 측정값(ms) > 75_000   →  FAIL (provisional)
   - 측정값(ms) ≤ 75_000   →  PASS

5. 수집 불가 SLI (예: 해당 파일 없음)
   - 해당 SLI는 평가 제외 (N/A). Overall 판정 시 “평가한 SLI”만 사용.
```

### 3.2 Severity: Stable 실패 vs Provisional 실패

- **Stable SLI 1개라도 FAIL** → **최종 결과는 FAIL.** (클러스터 수용 불가.)
- **Stable SLI는 모두 PASS이고, Provisional SLI만 FAIL** → **최종 결과는 PASS_WITH_WARNINGS.** (수용은 하되, 레이턴시 개선·임계값 정제 권장.)

### 3.3 Overall 판정 규칙

| 조건 | Overall 결과 |
|------|----------------|
| 모든 평가 대상 SLI(Stable + Provisional) PASS | **PASS** |
| Stable SLI 중 하나라도 FAIL | **FAIL** |
| Stable SLI는 모두 PASS, Provisional SLI 중 하나라도 FAIL | **PASS_WITH_WARNINGS** |
| Stable SLI 중 수집 실패(N/A)가 있어 판정 불가 | **FAIL** 또는 **판정 보류** (정책에 따라 결정. 권장: 필수 Stable SLI 수집 실패 시 FAIL) |

---

## 4. devcat 수용 결과 산출 (Acceptance Result)

### 4.1 가능한 최종 결과 값

| 결과 | 의미 |
|------|------|
| **PASS** | 모든 평가 SLI가 임계값을 만족. 클러스터 수용. |
| **PASS_WITH_WARNINGS** | Stable SLI는 모두 만족하나, **Provisional SLI 중 하나 이상 불만족.** 수용은 하되, 레이턴시·provisional 임계값 정제 권장. |
| **FAIL** | **Stable SLI 중 하나 이상 불만족** (OOM 발생 또는 시스템 파드 재시작). 클러스터 수용 불가. |

### 4.2 산출 형식 (권장)

devcat run 후 **results/<run-id>/** 에 다음을 생성하는 것을 권장한다.

1. **SLI 측정값 요약**  
   - 파일 예: `clusterloader2/sli-measurements.json` 또는 `cat-result.json`  
   - 내용 예: `{ "PodStartupLatency_P99_ms": 103551, "ClusterOOMsTracker_failures": 0, "SystemPodMetrics_max_restartCount": 0, ... }`

2. **SLI별 판정**  
   - 파일 예: `clusterloader2/cat-slo-evaluation.json` 또는 동일 파일 내  
   - 내용 예: `{ "ClusterOOMsTracker": "PASS", "SystemPodMetrics": "PASS", "PodStartupLatency_P99": "PASS", "StatelessPodStartupLatency_P99": "FAIL" }`

3. **Overall 결과**  
   - 파일 예: `clusterloader2/cat-result.txt` 또는 `manifest.txt` 에 한 줄  
   - 값: `CAT_RESULT=PASS` | `CAT_RESULT=PASS_WITH_WARNINGS` | `CAT_RESULT=FAIL`

4. **판정 요약 (사람이 읽기 쉬운 형태)**  
   - 파일 예: `clusterloader2/cat-summary.md`  
   - 내용: Overall 결과, SLI별 측정값·임계값·PASS/FAIL, Stable vs Provisional 구분.

- 구현은 **runbook·스크립트**로 할 수 있으며, ClusterLoader2가 생성한 JSON만 읽어 위 규칙으로 비교·판정하면 된다.

### 4.3 판정 플로우 요약

```
Measure SLI (results/<run-id>/clusterloader2/ 에서 추출)
    ↓
Compare each SLI with threshold (Stable / Provisional)
    ↓
Per-SLI: PASS / FAIL / N/A
    ↓
If any Stable SLI = FAIL  →  CAT_RESULT = FAIL
Else if any Provisional SLI = FAIL  →  CAT_RESULT = PASS_WITH_WARNINGS
Else  →  CAT_RESULT = PASS
```

---

## 5. Provisional SLO 값의 향후 정제

Provisional SLI의 임계값은 **고정이 아니며**, 아래 방향으로 정제한다.

### 5.1 정제를 위한 입력

- **추가 실험:** 이미지 preload·warm-up run·시나리오 규모 축소 등으로 **레이턴시 변동 원인**을 줄인 뒤, 새 분포 수집.
- **환경 구분:** **로컬(kind)** vs **실제 클러스터(또는 스테이징)** 에 따라 **환경별 임계값**을 둘 수 있음. (latency-variance-analysis 권장.)
- **데이터 기반 조정:** N run의 P99 분포(중앙값·백분위)를 보고, “대부분의 정상 run이 통과”하도록 임계값을 올리거나 내림.

### 5.2 정제 절차 (권장)

1. **실험·환경 전제**를 문서에 명시 (노드 수, 이미지 preload 여부, 시나리오 이름 등).
2. **새 베이스라인**을 반복 run으로 수집하고, 변동성 분석을 갱신.
3. **Practical acceptance meaning**과 **user expectations**를 정당화 문서에 반영한 뒤, Provisional 임계값을 조정.
4. 조정된 값을 **이 명세의 §2.2·§2.3**에 반영하고, **Provisional** 라벨을 유지하다가, 충분히 안정화되면 **Stable** 로 승격하거나 **환경별 확정 SLO**로 문서화.

### 5.3 Stable SLI와의 관계

- **Stable SLI(OOM, 재시작)** 는 **기준점**으로 유지. 정제 대상이 아님.
- **Provisional SLI** 만 “임시 임계값 → 실험 → 정제” 사이클을 도는 것으로 정의.

---

## 6. 요약

- **CAT v1 SLI:** Stable 2개(ClusterOOMsTracker failures, SystemPodMetrics max restartCount), Provisional 2개(PodStartupLatency P99, StatelessPodStartupLatency P99).
- **Stable SLO:** failures = 0, max restartCount = 0. 미달 시 **FAIL**.
- **Provisional SLO:** P99 ≤ 130 s, P99 ≤ 75 s (실험 베이스라인 기반). **임시**이며 정제 대상.
- **판정:** PASS / PASS_WITH_WARNINGS / FAIL. Stable 실패 → FAIL; Stable 통과·Provisional 실패 → PASS_WITH_WARNINGS; 전부 통과 → PASS.
- **devcat 산출:** SLI 측정값·SLI별 판정·CAT_RESULT·요약 문서를 results/<run-id>/clusterloader2/ (또는 동일 run 디렉터리)에 생성.
- **Provisional 정제:** 추가 실험·환경 구분·정당화 문서에 따라 임계값을 조정하고, 본 명세를 갱신.

---

*CAT v1 명세. devcat에서 첫 사용 가능한 수용 테스트 명세.*
