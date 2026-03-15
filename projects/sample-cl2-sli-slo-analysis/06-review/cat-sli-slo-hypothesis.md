# CAT 초기 SLI·SLO 가설

**Project:** sample-cl2-sli-slo-analysis  
**Input:** metric-baseline-analysis.md, sli-candidate-analysis.md, 실험 결과, Kubernetes SIG-Scalability SLO 철학

**목적:** 메트릭 베이스라인을 바탕으로 **CAT(Cluster Acceptance Test)용 초기 SLO 가설**을 제안하고, **최소 CAT v1 SLI 세트**와 **수용 판정 모델**을 정리한다.

**중요:** 아래 SLO 값은 **초기 가설(initial hypotheses)** 이며, **추가 실험을 통한 검증이 필요**하다. 환경·시나리오·실측 분포가 쌓이면 정당화 문서와 함께 조정할 것.

---

## 1. CAT 적합 메트릭별 초기 SLO 가설

베이스라인 분석에서 CAT 수용 기준 후보로 적합하다고 판단한 메트릭마다, **실험 베이스라인 값**, **Kubernetes 벤치마크 참고**, **초기 CAT SLO 가설**, **근거**를 표로 정리한다.

| SLI | baseline value (실험) | Kubernetes benchmark reference | initial CAT SLO hypothesis | reasoning |
|-----|------------------------|---------------------------------|----------------------------|-----------|
| **PodStartupLatency** (pod_startup, P99) | 전체 P99 64.1 s, Stateless P99 10.7 s (Run `20260315-181750`, kind 3노드) | SIG-Scalability: P99 < 5s (100노드·3000파드, prepulled image). 소규모는 수 초~수십 초 가능. | **P99 ≤ 30 s** (Create phase, pod_startup metric) | 현재 베이스라인이 5s 목표를 크게 상회함. 소규모·kind·이미지 풀 등 환경을 반영해 “실용적 수용 의미”를 갖도록 **30s**를 초기 가설로 둠. 5s는 장기 목표로 두고, 추가 실험으로 분포가 안정화되면 15s·20s 등으로 조정 검토. |
| **StatelessPodStartupLatency** (pod_startup, P99) | P99 10.7 s | 위와 동일 | **P99 ≤ 15 s** (Stateless만) | Stateless가 전체보다 짧음. 전체보다 엄격한 15s를 초기 가설로 두어, “stateless 워크로드는 상대적으로 빠르게 기동”을 반영. |
| **ClusterOOMsTracker** (failures) | failures: [] (0건) | 정상: OOM 0건 기대 | **failures 개수 = 0** | 워크로드 안정성·클러스터 건강도의 결함 지표. 0이 아니면 수용 실패로 두는 것이 일반적. |
| **SystemPodMetrics** (시스템 파드 restartCount) | 전체 0 | 정상: 0 또는 허용 수준 | **모든 시스템 파드 restartCount = 0** (또는 시나리오 내 허용 상한 0) | 컨트롤 플레인·시스템 컴포넌트 안정성. v1에서는 “재시작 없음”을 수용 기준으로 둠. |
| **APIResponsivenessPrometheus** (슬로우 콜·레이턴시) | 미수집 (Prometheus run 대기) | verb/resource별 P99·슬로우 콜 수 목표; 슬로우 콜 0 또는 허용 개수 이하 | **슬로우 콜 수 = 0** (또는 허용 개수 이하). P99 레이턴시는 Prometheus run 베이스라인 확보 후 가설 추가 | API 반응성은 사용자·워크로드 체감과 직결. SIG-Scalability에서도 슬로우 콜·레이턴시 목표를 둠. 수치 베이스라인이 없어 “슬로우 콜 0”을 우선 가설로 두고, P99는 추후 보완. |

---

## 2. 최소 CAT v1 SLI 세트

**클러스터 건강도**와 **워크로드 준비도**를 가장 잘 대표하는 **4~5개 SLI**로 최소 CAT v1 세트를 정의한다.

### 2.1 선정 결과: CAT v1 SLI 세트 (5개)

| # | SLI | metric source | 선정 근거 |
|---|-----|---------------|-----------|
| 1 | **PodStartupLatency** (pod_startup, P99) | CreatePhasePodStartupLatency | **워크로드 체감 직접 반영.** 배포·스케일 시 “파드가 언제 쓸 수 있게 되는가”를 측정. SIG-Scalability 핵심 SLO 중 하나. devcat에서 이미 산출 확인됨. |
| 2 | **ClusterOOMsTracker** (failures) | ClusterOOMsTracker | **클러스터 건강도·워크로드 안정성.** OOM 발생 시 수용 불가로 두기 쉬우며, 측정·판정이 명확함(0 vs 비어 있지 않음). devcat에서 산출 확인됨. |
| 3 | **SystemPodMetrics** (시스템 파드 restartCount) | SystemPodMetrics | **클러스터 건강도.** 컨트롤 플레인·시스템 컴포넌트가 테스트 구간 동안 재시작 없이 동작하는지를 반영. devcat에서 산출 확인됨. |
| 4 | **APIResponsivenessPrometheus** (슬로우 콜, 필요 시 P99) | APIResponsivenessPrometheus | **사용자·워크로드 체감.** API 반응성은 kubectl·컨트롤러 동작에 직결. Prometheus precheck 통과로 수집 가능; v1에서는 “슬로우 콜 0” 등으로 수용 판단. |
| 5 | **StatelessPodStartupLatency** (pod_startup, P99) — 선택 | CreatePhasePodStartupLatency | **워크로드 유형별 세분화.** PodStartupLatency와 동일 소스이지만 stateless만 보면 “일반적인 무상태 워크로드” 수용을 더 엄격하게 볼 수 있음. v1에서 4개로 줄이면 이 항목은 제외 가능. |

- **4개로 최소화할 경우:** 위 1~4만 두고, 5(StatelessPodStartupLatency)는 “PodStartupLatency에 포함된 정보의 부분집합”으로 보아 제외해도 됨.  
- **5개로 둘 경우:** Stateless 전용 P99를 두어 “stateless 워크로드 기동 품질”을 명시적으로 수용 기준에 넣음.

### 2.2 v1에서 제외·후순위로 둔 메트릭

| 메트릭 | 이유 |
|--------|------|
| StatefulPodStartupLatency | PodStartupLatency와 동일 계열; v1에서는 전체 pod_startup 하나로 대표. 필요 시 v2에서 유형별로 확장. |
| ResourceUsageSummary | 보조 해석용; 수용 통과/실패의 주 기준으로 쓰지 않음. |
| Phase duration | kind에서 SchedulingMetrics 불안정; 측정 확립 후 v2에서 검토. |
| InClusterNetworkLatency, NetworkProgrammingLatency, KubeProxyIptablesRestoreFailures | Prometheus 의존·베이스라인 미확립; 수집 안정화 후 v2 후보. |
| APIAvailability | 옵션·산출 경로 확인 필요; v2에서 검토. |

---

## 3. CAT 평가 모델 (Cluster Acceptance Decision)

devcat이 **클러스터 수용 여부**를 결정하는 방식은 아래와 같이 단순 비교 모델로 둔다.

### 3.1 단계별 흐름

```
1. 측정 (Measure SLI)
   → ClusterLoader2 실행 후 results/<run-id>/clusterloader2/ 에서
     각 SLI에 해당하는 산출 파일·키를 읽어 SLI 측정값을 추출한다.

2. SLO와 비교 (Compare with SLO)
   → CAT v1 SLI 세트에 정의된 “initial CAT SLO hypothesis”와
     측정값을 비교한다.
   → 비교 규칙:
     - “낮을수록 좋은” 지표 (레이턴시): measured ≤ threshold 이면 해당 SLI PASS.
     - “0이어야 함” 지표 (OOM 건수, 슬로우 콜 수, 재시작 수): measured == 0 (또는 ≤ 허용 상한) 이면 PASS.
     - “높을수록 좋은” 지표 (가용성 % 등): measured ≥ threshold 이면 PASS.

3. 메트릭별 PASS/FAIL
   → 각 SLI마다 PASS 또는 FAIL을 부여한다.
   → FAIL인 경우 (선택) 측정값 vs 임계값을 기록해 해석·개선에 활용한다.

4. 전체 클러스터 수용 판정 (Overall cluster acceptance)
   → CAT v1 SLI 세트의 **모든** SLI가 PASS이면 → **Cluster ACCEPT**.
   → **하나라도** FAIL이면 → **Cluster REJECT** (또는 “조건부 수용·재검증 권장”).
```

### 3.2 수집 불가 SLI 처리

- **해당 run에서 산출이 없는 SLI** (예: Prometheus 비활성 run에서 APIResponsivenessPrometheus):  
  - **N/A**로 두고, “평가 대상에서 제외”하거나,  
  - “수집 불가 시 FAIL” 정책을 두어 수용 전제 조건(측정 파이프라인 준비)을 강제할 수 있다.  
- v1 권장: **수집 가능한 SLI만** 평가 대상으로 하고, N/A는 Overall 판정에서 “해당 SLI 제외”로 처리.  
  - 단, “CAT v1 SLI 세트 중 수집 가능한 것이 4개 미만이면 전체 판정을 보류”하는 식의 최소 조건을 둘 수 있음.

### 3.3 출력 예시

- **SLI별:** `PodStartupLatency P99: 10.7s (threshold 30s) → PASS`  
  `ClusterOOMsTracker failures: 0 → PASS`  
  `SystemPodMetrics restarts: 0 → PASS`  
  `APIResponsivenessPrometheus slow calls: N/A (미수집) → N/A`  
- **Overall:** `Cluster ACCEPT` (평가한 SLI 전부 PASS) 또는 `Cluster REJECT` (하나라도 FAIL), 또는 `판정 보류` (N/A 과다 등).

---

## 4. 초기 SLO 가설 요약표 (CAT v1)

| SLI | measurement location | initial SLO hypothesis | 비고 |
|-----|----------------------|------------------------|------|
| PodStartupLatency (pod_startup, P99) | `PodStartupLatency_*_load_*.json` → dataItems, Metric=pod_startup, Perc99 | **P99 ≤ 30 s** | 소규모·kind 환경 반영. 추가 실험으로 조정. |
| ClusterOOMsTracker (failures) | `ClusterOOMsTracker_load_*.json` → failures | **failures = 0** | |
| SystemPodMetrics (시스템 파드 restartCount) | `SystemPodMetrics_load_*.json` → pods[].containers[].restartCount | **전체 0** | |
| APIResponsivenessPrometheus (슬로우 콜) | `APIResponsivenessPrometheus*_*_load_*.json` (Prometheus run 시) | **슬로우 콜 = 0** (또는 허용 개수 이하) | 베이스라인 확보 후 P99 가설 추가 가능. |
| StatelessPodStartupLatency (pod_startup, P99) — 선택 | `StatelessPodStartupLatency_*_*.json` → pod_startup, Perc99 | **P99 ≤ 15 s** | v1에서 4개만 쓰면 제외 가능. |

---

## 5. 검증 필요 사항 (Validation)

아래 SLO 값은 **초기 가설**이며, 다음을 통해 검증·조정해야 한다.

- **추가 실험:** 동일 시나리오·환경에서 반복 실행해 측정값 분포(P50/P90/P99)를 확인하고, “대부분의 run이 통과” 또는 “합리적 비율이 통과”하도록 임계값 조정.
- **환경 가정 명시:** “kind 3노드”, “prepulled image 여부”, “create phase 파드 수” 등을 전제로 한 정당화 문서 보완.
- **User expectations·Practical acceptance meaning:** “이 SLO를 만족하면 이 클러스터를 수용한다”가 팀·운영에서 말이 되도록, 필요 시 임계값 완화 또는 강화.
- **APIResponsivenessPrometheus:** Prometheus 활성화 full run 완료 후 실제 슬로우 콜·P99 값을 반영해 가설 보완.

---

*CAT 초기 SLI·SLO 가설. SLO 값은 초기 가설이며, 추가 실험을 통한 검증이 필요함.*
