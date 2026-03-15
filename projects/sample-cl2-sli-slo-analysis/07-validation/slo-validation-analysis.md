# SLO 검증 실험 분석

**Project:** sample-cl2-sli-slo-analysis  
**Input:** cat-sli-slo-hypothesis.md, 동일 시나리오·오버라이드로 수행한 반복 실험, SLI 추출 결과

**목적:** CAT v1 SLI 세트에 대한 **초기 SLO 가설**을 반복 실험으로 검증하고, 측정값 분포·변동성·가설 적절성(too strict / reasonable / too loose)을 정리한다.

---

## 1. 검증 실험 설정

- **시나리오:** `scenarios/load/config.yaml`
- **오버라이드:** `overrides/ol-test.yaml` (단일 오버라이드, Prometheus 비활성)
- **환경:** kind 3노드, 동일 kubeconfig
- **실행 횟수:** 4 runs  
  - 기존 베이스라인 run: `20260315-181750`  
  - SLO 검증용 3회: `slo-validation-1`, `slo-validation-2`, `slo-validation-3`
- **ClusterLoader2:** devcat `bin/cl2`, PROVIDER=local, report-dir=`results/<run-id>/clusterloader2`
- **APIResponsivenessPrometheus:** 본 실험은 Prometheus 비활성으로 실행하여 **미수집(N/A)**. 검증 대상에서 제외.

---

## 2. Run별 SLI 측정값 (분포 테이블)

### 2.1 PodStartupLatency (pod_startup, P99)

단위: ms (아래 요약 표에서는 초(s)로 환산해 병기).

| Run | PodStartupLatency P99 (ms) | PodStartupLatency P99 (s) |
|-----|----------------------------|---------------------------|
| 20260315-181750 | 64120.3 | 64.1 |
| slo-validation-1 | 74468.6 | 74.5 |
| slo-validation-2 | 122281.9 | 122.3 |
| slo-validation-3 | 103551.0 | 103.6 |

### 2.2 StatelessPodStartupLatency (pod_startup, P99)

| Run | StatelessPodStartupLatency P99 (ms) | (s) |
|-----|-------------------------------------|-----|
| 20260315-181750 | 10727.9 | 10.7 |
| slo-validation-1 | 24966.0 | 25.0 |
| slo-validation-2 | 69906.2 | 69.9 |
| slo-validation-3 | 51722.1 | 51.7 |

### 2.3 ClusterOOMsTracker (failures 개수)

| Run | OOM failures |
|-----|---------------|
| 20260315-181750 | 0 |
| slo-validation-1 | 0 |
| slo-validation-2 | 0 |
| slo-validation-3 | 0 |

### 2.4 SystemPodMetrics (시스템 파드 최대 restartCount)

| Run | max restartCount |
|-----|------------------|
| 20260315-181750 | 0 |
| slo-validation-1 | 0 |
| slo-validation-2 | 0 |
| slo-validation-3 | 0 |

### 2.5 APIResponsivenessPrometheus

| Run | 슬로우 콜 / P99 |
|-----|-----------------|
| (전 run) | N/A (Prometheus 비활성) |

### 2.6 통합 분포 표 (요약)

| Run | PodStartupLatency P99 (s) | Stateless P99 (s) | OOM failures | System restartCount |
|-----|---------------------------|-------------------|--------------|---------------------|
| 20260315-181750 | 64.1 | 10.7 | 0 | 0 |
| slo-validation-1 | 74.5 | 25.0 | 0 | 0 |
| slo-validation-2 | 122.3 | 69.9 | 0 | 0 |
| slo-validation-3 | 103.6 | 51.7 | 0 | 0 |

---

## 3. 변동성 분석

### 3.1 PodStartupLatency (pod_startup, P99)

| 항목 | 값 |
|------|-----|
| **typical value** | 중앙값 약 **89 s** (74.5와 103.6 사이), 평균 약 **91.1 s**. 4 run 모두 60 s 이상. |
| **worst observed** | **122.3 s** (slo-validation-2). |
| **best observed** | **64.1 s** (20260315-181750). |
| **stability across runs** | **변동 큼.** 64 s ~ 122 s 구간에서 분포. 이미지 풀·노드 부하·스케줄링 타이밍 등에 따라 run 간 차이가 큰 것으로 해석. |

### 3.2 StatelessPodStartupLatency (pod_startup, P99)

| 항목 | 값 |
|------|-----|
| **typical value** | 4 run 값: 10.7, 25.0, 69.9, 51.7 s. 중앙값 약 **38 s**, 평균 약 **39.3 s**. |
| **worst observed** | **69.9 s** (slo-validation-2). |
| **best observed** | **10.7 s** (20260315-181750). |
| **stability across runs** | **변동 매우 큼.** 10.7 s ~ 69.9 s. Stateless만으로도 run 간 편차가 커서, 단일 run으로 SLO 판단하기엔 불안정. |

### 3.3 ClusterOOMsTracker (failures)

| 항목 | 값 |
|------|-----|
| **typical value** | **0**. |
| **worst observed** | **0**. |
| **stability** | 4 run 모두 동일. **안정.** |

### 3.4 SystemPodMetrics (restartCount)

| 항목 | 값 |
|------|-----|
| **typical value** | **0**. |
| **worst observed** | **0**. |
| **stability** | 4 run 모두 동일. **안정.** |

---

## 4. 초기 SLO 가설에 대한 평가

cat-sli-slo-hypothesis.md에 정의된 **initial CAT SLO hypothesis**를 위 실험 결과와 비교한 판정이다.

### 4.1 PodStartupLatency (pod_startup, P99)

- **초기 가설:** P99 ≤ 30 s  
- **실험 결과:** 4 run 모두 **30 s 초과** (64.1, 74.5, 122.3, 103.6 s).  
- **판정:** **Too strict (너무 엄격)**  
  - 현재 kind·3노드·동일 시나리오에서 30 s 이하 P99를 만족한 run이 없음.  
  - 수용 기준으로 그대로 두면 “항상 REJECT”에 가깝게 됨.  
- **권장:** 이 환경·시나리오에 맞춰 임계값을 **완화** 검토. 예: P99 ≤ 90 s 또는 ≤ 120 s로 두고, 추가 run으로 통과 비율을 본 뒤 60 s·75 s 등으로 조정.

### 4.2 StatelessPodStartupLatency (pod_startup, P99)

- **초기 가설:** P99 ≤ 15 s  
- **실험 결과:** 4 run 중 **1 run만 15 s 이하** (10.7 s). 나머지 25.0, 69.9, 51.7 s.  
- **판정:** **Too strict (너무 엄격)**  
  - 대부분의 run이 15 s를 초과.  
- **권장:** 환경·시나리오를 전제로 **완화** 검토. 예: P99 ≤ 60 s 또는 ≤ 70 s. 또는 Stateless를 v1 필수 SLI에서 제외하고 PodStartupLatency(전체)만 두고 조정할 수 있음.

### 4.3 ClusterOOMsTracker (failures)

- **초기 가설:** failures = 0  
- **실험 결과:** 4 run 모두 **0**.  
- **판정:** **Reasonable (적절)**  
  - 현재 관측으로는 가설이 달성 가능하고, “OOM 0건”은 수용 기준으로 유지하는 것이 타당함.

### 4.4 SystemPodMetrics (시스템 파드 restartCount)

- **초기 가설:** 전체 0  
- **실험 결과:** 4 run 모두 **0**.  
- **판정:** **Reasonable (적절)**  
  - “시스템 파드 재시작 0회”는 유지해도 됨.

### 4.5 APIResponsivenessPrometheus

- **초기 가설:** 슬로우 콜 = 0 (및 추후 P99)  
- **실험 결과:** **미수집 (N/A)**.  
- **판정:** **평가 보류.**  
  - Prometheus 활성화 run으로 데이터가 쌓이면 동일 방식으로 분포·가설 적절성 평가 필요.

### 4.6 SLO 가설 평가 요약표

| SLI | initial SLO hypothesis | 실험 결과 (요약) | 판정 |
|-----|-------------------------|------------------|------|
| PodStartupLatency P99 | ≤ 30 s | 64.1, 74.5, 122.3, 103.6 s (전부 초과) | **Too strict** |
| StatelessPodStartupLatency P99 | ≤ 15 s | 10.7, 25.0, 69.9, 51.7 s (1/4 만 충족) | **Too strict** |
| ClusterOOMsTracker failures | = 0 | 0 (전 run) | **Reasonable** |
| SystemPodMetrics restartCount | = 0 | 0 (전 run) | **Reasonable** |
| APIResponsivenessPrometheus | 슬로우 콜 = 0 | N/A | **평가 보류** |

---

## 5. 결론 및 다음 단계

- **레이턴시 SLO:**  
  - PodStartupLatency P99 ≤ 30 s, Stateless P99 ≤ 15 s는 **현재 환경에서 모두 너무 엄격**하다.  
  - 반복 실험에서 P99가 64 s ~ 122 s, Stateless 10.7 s ~ 69.9 s로 넓게 퍼져 있으므로, **임계값 완화** 또는 “이 환경에서는 레이턴시 SLO를 별도 정당화” 후 수치 조정이 필요하다.
- **결함/안정성 SLO:**  
  - OOM 0건, 시스템 파드 재시작 0회는 **그대로 유지해도 무방**하다. (Reasonable.)
- **API 반응성:**  
  - Prometheus 활성화 run으로 슬로우 콜·P99를 수집한 뒤, 동일한 검증 절차(분포 테이블·변동성·가설 평가)를 적용할 것.
- **추가 실험 권장:**  
  - 동일 시나리오로 run을 더 쌓아 분포를 보강하고,  
  - “P99 ≤ 90 s” 등 완화된 가설에 대해 통과 비율을 확인한 뒤,  
  - practical acceptance meaning과 환경 가정을 문서에 반영해 최종 SLO 값을 확정하는 단계를 권장한다.

---

*SLO 검증 실험 분석. 초기 가설은 검증을 통해 조정이 필요함.*
