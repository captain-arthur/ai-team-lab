# 메트릭 베이스라인 분석

**Project:** sample-cl2-sli-slo-analysis  
**Input:** 실험 결과(experiment-notes.md), ClusterLoader2 measurement 산출물, SIG-Scalability SLO 문서 맥락, sli-candidate-analysis.md

**목적:** ClusterLoader2 메트릭의 **의미**를 이해하고, 현재 실험에서 관측된 값을 **베이스라인**으로 정리한다.  
**범위:** 수용 임계값(SLO)을 정하는 단계가 아님. 메트릭 이해와 베이스라인 확립만 수행.

---

## 1. 분석 대상 메트릭 개요

실험에서 **실제 산출물이 있는** 메트릭과, **Prometheus 활성화 시 수집 가능**한 메트릭을 대상으로 한다.

| 메트릭 | 실험 산출 여부 | 비고 |
|--------|----------------|------|
| PodStartupLatency / StatelessPodStartupLatency / StatefulPodStartupLatency | ✅ 있음 | CreatePhasePodStartupLatency, Run `20260315-181750` |
| ClusterOOMsTracker | ✅ 있음 | 동일 run |
| SystemPodMetrics | ✅ 있음 | 동일 run |
| APIResponsivenessPrometheus | ❌ 해당 run에서 미수집 | Prometheus 비활성 run; Precheck 통과 후 full run에서 수집 가능 |
| ResourceUsageSummary | ✅ 있음 | Phase·리소스 보조 해석용 |

---

## 2. 메트릭별 상세 분석

### 2.1 PodStartupLatency (CreatePhasePodStartupLatency)

**무엇을 측정하는가**  
파드 생성 요청부터 파드가 Ready 조건을 만족할 때까지 걸리는 시간. 하위 구간으로 `pod_startup`(전체), `run_to_watch`(클라이언트가 watch 시작하기까지), `schedule_to_watch`(스케줄링부터 watch까지) 등이 있다.

**ClusterLoader2에서 생성 위치**  
- Measurement: **CreatePhasePodStartupLatency** (reconcile-objects 단계에서 수집).  
- 모듈: `measurements.yaml`의 PodStartupLatency, StatelessPodStartupLatency, StatefulPodStartupLatency.  
- 산출 파일: `PodStartupLatency_CreatePhasePodStartupLatency_load_*.json`, `StatelessPodStartupLatency_*_*.json`, `StatefulPodStartupLatency_*_*.json`.  
- 데이터 경로: `dataItems[].labels.Metric`, `dataItems[].data.Perc50|Perc90|Perc99` (단위: ms).

**SIG-Scalability에서 중요한 이유**  
- [Kubernetes SIG-Scalability SLO](https://github.com/kubernetes/community/blob/master/sig-scalability/slos/pod_startup_latency.md): **P99 pod startup latency < 5s**를 대규모 클러스터(100노드, 3000파드, prepulled image) 기준으로 제시.  
- “사용자가 체감하는 배포·스케일링 품질”과 직결되며, scalable한 클러스터는 노드·파드 수에 대해 선형 이하로 증가해야 한다는 기대가 있음.

**일반적인 벤치마크 환경에서의 전형적 범위**  
- 대규모: 100노드·3000파드 수준에서 P99 < 5s 목표.  
- 소규모: 노드·파드 수가 적을 때는 수 초 이내가 기대되나, 이미지 풀·리소스 경합에 따라 수십 초까지 나올 수 있음.  
- 공식 ClusterLoader2 CreatePhasePodStartupLatency threshold는 **1h**(대규모 saturation “전부 기동될 때까지” 대기용)로, 수용 판단용으로는 부적합하다는 것이 Research·Architecture에서 정리됨.

**현재 실험에서 나온 값 (Run `20260315-181750`, kind 3노드)**  
- **PodStartupLatency (전체), metric `pod_startup`:**  
  - Perc50: 8991.7 ms (~9.0 s), Perc90: 44905.0 ms (~44.9 s), Perc99: 64120.3 ms (~64.1 s).  
- **StatelessPodStartupLatency, metric `pod_startup`:**  
  - Perc50: 5319.4 ms (~5.3 s), Perc90: 10612.1 ms (~10.6 s), Perc99: 10727.9 ms (~10.7 s).  
- **StatefulPodStartupLatency:** 동일 run에서 파일 존재; 상세 값은 전체 PodStartupLatency와 유사한 구간으로 해석 가능.

**해석**  
- P99가 SIG-Scalability 목표(5s)보다 현저히 큼(전체 ~64s, stateless ~10.7s). kind·소규모 환경, 리소스·이미지 풀 등으로 인한 지연으로 보임.  
- Stateless만 보면 P50은 ~5.3s로 5s 근처이나, P99는 10.7s로 수용 목표로 쓰려면 “현재 환경의 베이스라인”으로 기록하고, 임계값은 별도 정당화 후 정할 필요가 있음.

---

### 2.2 StatelessPodStartupLatency / StatefulPodStartupLatency

**무엇을 측정하는가**  
Create phase 동안 생성된 파드 중 **stateless** / **stateful** 로 구분한 pod startup latency. 동일한 “파드 생성 → Ready” 구간을 workload 유형별로 집계한 것이다.

**ClusterLoader2에서 생성 위치**  
- CreatePhasePodStartupLatency의 labelSelector=load 하에, stateless/stateful 분류에 따라 별도 JSON으로 집계.  
- 파일: `StatelessPodStartupLatency_CreatePhasePodStartupLatency_load_*.json`, `StatefulPodStartupLatency_CreatePhasePodStartupLatency_load_*.json`.

**SIG-Scalability·수용 테스트에서 중요한 이유**  
- Pod startup latency와 동일한 “배포·스케일 체감” 의미.  
- Stateful 워크로드는 볼륨·초기화 등으로 stateless보다 지연이 클 수 있어, 유형별로 베이스라인을 보는 것이 해석에 도움이 됨.

**현재 실험에서 나온 값**  
- Stateless `pod_startup`: Perc50 5.3 s, Perc90 10.6 s, Perc99 10.7 s (위와 동일).  
- Stateful: 동일 run에서 파일 존재; 정량은 전체 PodStartupLatency와 같은 JSON 구조로 확인 가능.

**해석**  
- Stateless 구간이 전체보다 낮은 편. 현재 환경의 “워크로드 유형별 베이스라인”으로 활용 가능.

---

### 2.3 ClusterOOMsTracker

**무엇을 측정하는가**  
테스트 구간 동안 클러스터에서 발생한 **OOM(Out of Memory) kill** 이벤트를 추적한다. `failures`(이번 run에서 발생), `ignored`(무시 정책에 따른 것), `past`(과거 이벤트)로 구분된다.

**ClusterLoader2에서 생성 위치**  
- Measurement: **ClusterOOMsTracker** (TestMetrics 계열).  
- 산출 파일: `ClusterOOMsTracker_load_*.json`.  
- 구조: `{ "failures": [], "ignored": [], "past": [] }` 형태의 목록.

**SIG-Scalability·수용 테스트에서 중요한 이유**  
- OOM은 워크로드 안정성과 직결. “이 클러스터에서 워크로드가 메모리 부족 없이 동작하는가”를 보는 결함 지표.  
- 수용 기준으로는 보통 “failures가 비어 있음(OOM 0건)”을 기대.

**일반적인 벤치마크 환경에서의 전형적 범위**  
- 정상: `failures`가 비어 있음.  
- 리소스 부족·과도한 부하 시: failures에 OOM 이벤트가 기록됨.

**현재 실험에서 나온 값 (Run `20260315-181750`)**  
- `failures: []`, `ignored: []`, `past: []` — OOM 발생 없음.

**해석**  
- 현재 run·환경에서는 OOM이 없던 베이스라인. 수용 기준 후보로 “OOM 0건”을 두기 좋은 메트릭.

---

### 2.4 SystemPodMetrics

**무엇을 측정하는가**  
시스템 네임스페이스(kube-system 등)의 **시스템 파드**(coredns, etcd, kube-apiserver, kube-proxy, kube-scheduler 등)에 대한 메트릭. 파드·컨테이너별 `restartCount`, `lastRestartReason` 등을 수집한다.

**ClusterLoader2에서 생성 위치**  
- Measurement: **TestMetrics** 계열의 시스템 파드 메트릭 수집.  
- 산출 파일: `SystemPodMetrics_load_*.json`.  
- 구조: `pods[].containers[].restartCount`, `lastRestartReason` 등.

**SIG-Scalability·수용 테스트에서 중요한 이유**  
- 컨트롤 플레인·시스템 컴포넌트 안정성 지표. 재시작이 많으면 플랫폼 불안정 또는 리소스/설정 문제 신호.  
- “클러스터 건강도”와 “워크로드가 기대하는 플랫폼 상태”를 반영.

**일반적인 벤치마크 환경에서의 전형적 범위**  
- 정상: restartCount 0 또는 허용된 수준.  
- 비정상: 특정 컴포넌트의 반복 재시작은 조사 대상.

**현재 실험에서 나온 값 (Run `20260315-181750`)**  
- 시스템 파드(coredns, etcd, kube-apiserver, kube-proxy, kube-scheduler 등)의 `restartCount`가 모두 **0**, `lastRestartReason`은 `""`.

**해석**  
- 현재 run에서 시스템 파드 재시작이 없던 베이스라인. “시스템 파드 재시작 0회”는 CAT 수용 기준 후보로 적합.

---

### 2.5 APIResponsivenessPrometheus

**무엇을 측정하는가**  
Prometheus에 수집된 API 서버 메트릭을 이용해 **API 호출별 레이턴시**(verb/resource별)와 **슬로우 콜** 개수를 산출한다. APIResponsivenessPrometheusSimple은 단순화된 레이턴시 쿼리를 사용한다.

**ClusterLoader2에서 생성 위치**  
- Measurement: **APIResponsivenessPrometheus**, **APIResponsivenessPrometheusSimple**.  
- Prometheus scrape가 필요하며, apiserver 메트릭(예: `apiserver_request_duration_seconds`)을 쿼리.  
- 산출 파일: Prometheus 활성화·gather 완료 시 `APIResponsivenessPrometheus_*_load_*.json`, `APIResponsivenessPrometheus_simple_*_*.json` 등.

**SIG-Scalability에서 중요한 이유**  
- API 반응성은 kubectl·컨트롤러·오퍼레이터 동작에 직결. 슬로우 콜이 많으면 사용자·워크로드가 “멈춤”처럼 체감.  
- 대규모 클러스터에서 API 레이턴시·슬로우 콜 수에 대한 목표가 있음(문서·perf 테스트에서 verb/resource별 임계값 참조).

**일반적인 벤치마크 환경에서의 전형적 범위**  
- verb/resource별 P99 레이턴시 목표(수백 ms ~ 수 초 수준, 리소스에 따라 상이).  
- 슬로우 콜: 0 또는 허용 개수 이하를 목표로 하는 설정이 많음.

**현재 실험에서 나온 값**  
- Run `20260315-181750`은 **Prometheus 비활성**이라 해당 measurement 미수집.  
- Prometheus measurement readiness precheck(§10–§11) 통과 후 full run(6443 scrape port)이 완료되면 동일 시나리오에서 산출 가능.  
- **베이스라인:** 해당 run 완료 전까지는 “수집 가능한 메트릭”으로만 분류하고, 수치 베이스라인은 추후 run 결과로 보완.

**해석**  
- API 반응성은 클러스터 건강도·사용자 체감과 직결되어 CAT 수용 기준 후보에 적합.  
- 현재는 “Prometheus 활성화 run에서 측정 가능” 단계로 두고, 실제 값이 쌓이면 본 문서의 베이스라인 표에 추가할 것.

---

### 2.6 ResourceUsageSummary

**무엇을 측정하는가**  
테스트 단계별 **리소스 사용 요약**(CPU/메모리 등)과 phase 구간 정보. 단계 타이머가 불완전한 환경에서도 phase 구간 해석에 쓸 수 있다.

**ClusterLoader2에서 생성 위치**  
- TestMetrics·ResourceUsage 수집의 요약.  
- 산출 파일: `ResourceUsageSummary_load_*.json`.

**수용 테스트에서의 의미**  
- 주 SLI보다는 “왜 레이턴시가 이렇게 나왔는가”를 해석하는 **보조** 용도.  
- Phase duration SLI가 불완전할 때 구간 길이 참고용.

**현재 실험에서 나온 값**  
- Run `20260315-181750`에서 파일 존재. 상세 숫자는 파일 열어 확인 가능.  
- 본 단계에서는 “phase·리소스 보조 해석용”으로만 기록하고, CAT 수용 기준 후보에서는 제외해도 됨.

---

## 3. 메트릭 베이스라인 테이블

| Metric | Source | Meaning | Experiment value (Run `20260315-181750`) | Typical Kubernetes benchmark range | Interpretation |
|--------|--------|---------|----------------------------------------|------------------------------------|----------------|
| **PodStartupLatency (pod_startup, 전체)** | CreatePhasePodStartupLatency | 파드 생성 → Ready 소요 시간 | P50 9.0 s, P90 44.9 s, P99 64.1 s | 대규모: P99 < 5s 목표; 소규모: 수 초~수십 초 가능 | 현재 환경 P99가 목표보다 큼; kind·리소스 등 환경 요인 반영. 베이스라인으로 확립 후 임계값 별도 정당화 필요. |
| **StatelessPodStartupLatency (pod_startup)** | CreatePhasePodStartupLatency | Stateless 파드 생성 → Ready | P50 5.3 s, P90 10.6 s, P99 10.7 s | 위와 유사, stateless가 보통 더 짧음 | Stateless만 보면 P50은 5s 근처; P99는 10.7s로 수용 목표로 쓰려면 환경별 정당화 필요. |
| **StatefulPodStartupLatency (pod_startup)** | CreatePhasePodStartupLatency | Stateful 파드 생성 → Ready | (동일 run에서 파일 존재; 전체와 유사 구간으로 해석) | Stateless보다 클 수 있음 | 유형별 베이스라인으로 활용. |
| **ClusterOOMsTracker (failures)** | ClusterOOMsTracker | 테스트 구간 OOM 발생 건수 | failures: [] (0건) | 정상: 0 | OOM 없음. CAT 수용 기준 후보로 적합. |
| **SystemPodMetrics (restartCount)** | TestMetrics / SystemPodMetrics | 시스템 파드 재시작 횟수 | 전체 0 | 정상: 0 또는 허용 수준 | 재시작 없음. CAT 수용 기준 후보로 적합. |
| **APIResponsivenessPrometheus** | APIResponsivenessPrometheus (Prometheus) | API 호출 레이턴시·슬로우 콜 | 미수집 (Prometheus 비활성 run) | verb/resource별 P99·슬로우 콜 목표 존재 | Prometheus 활성 full run 완료 후 베이스라인 추가. CAT 후보에 적합. |
| **ResourceUsageSummary** | TestMetrics / ResourceUsage | 단계별 리소스·구간 요약 | 파일 존재 (상세 값 생략) | — | 보조 해석용; CAT 주 기준으로는 미사용. |

---

## 4. CAT 수용 기준 후보로 적합한 메트릭

다음 네 가지에 초점을 두고, “수용 기준으로 쓰기 적합한가”를 판단했다.

- **클러스터 건강도 반영:** 컨트롤 플레인·시스템·리소스 안정성.  
- **사용자·워크로드 체감 반영:** 배포·API·네트워크 등 체감 품질.  
- **run 간 안정성:** 동일 시나리오에서 재현 가능하게 측정 가능.  
- **devcat에서 측정 가능:** ClusterLoader2·Prometheus 산출로 실제 수집 가능.

### 4.1 적합한 메트릭 (권장 CAT 후보)

| 메트릭 | 클러스터 건강도 | 워크로드 체감 | 안정성 | devcat 측정 | 비고 |
|--------|-----------------|---------------|--------|-------------|------|
| **PodStartupLatency** (pod_startup) | △ (스케줄·리소스 반영) | ✅ 직접 | ✅ JSON으로 재현 가능 | ✅ 확인됨 | 임계값은 환경·시나리오 정당화 후 정의. |
| **ClusterOOMsTracker** (failures) | ✅ | ✅ | ✅ | ✅ 확인됨 | OOM 0건을 수용 기준으로 두기 좋음. |
| **SystemPodMetrics** (restartCount) | ✅ | △ (간접) | ✅ | ✅ 확인됨 | 시스템 파드 재시작 0회 등. |
| **APIResponsivenessPrometheus** | ✅ | ✅ | ✅ (Prometheus 수집 시) | ✅ Precheck 통과, full run 시 수집 가능 | 수치 베이스라인은 Prometheus run 결과로 보완. |

### 4.2 보조·선택 메트릭

| 메트릭 | 비고 |
|--------|------|
| **StatelessPodStartupLatency / StatefulPodStartupLatency** | PodStartupLatency의 유형별 분해; 필요 시 수용 기준을 유형별로 세분화할 때 사용. |
| **ResourceUsageSummary** | Phase·리소스 해석용; 수용 기준보다는 “이유 분석”용. |
| **Phase duration** | kind에서 SchedulingMetrics 실패로 불완전; 측정 안정화 후 CAT 후보 검토 가능. |

### 4.3 아직 베이스라인 미확립

| 메트릭 | 비고 |
|--------|------|
| **InClusterNetworkLatency** | Prometheus 의존; full run 후 값 확보 시 베이스라인·CAT 적합성 재검토. |
| **NetworkProgrammingLatency** | 동일. |
| **KubeProxyIptablesRestoreFailures** | 동일; 0 기대 지표로 CAT 후보 가능. |
| **APIAvailability** | 옵션·산출 경로 확인 후 베이스라인 확립. |

---

## 5. 요약

- **베이스라인 확립됨:** PodStartupLatency(전체·Stateless·Stateful), ClusterOOMsTracker, SystemPodMetrics.  
  - 현재 실험 값은 위 베이스라인 표에 반영됨.  
  - **수용 임계값(SLO)은 이 단계에서 정의하지 않음.** 이해와 베이스라인만 정리.
- **Prometheus 의존 메트릭:** APIResponsivenessPrometheus 등은 Precheck 통과 후 full run 완료 시 베이스라인 표에 추가할 것.
- **CAT 수용 기준 후보:** Pod startup latency, OOM 0건, 시스템 파드 재시작 0회, API 반응성(수집 가능 시)이 “클러스터 건강도·워크로드 체감·안정성·측정 가능성” 관점에서 적합하다고 판단됨.  
  최종 SLO 임계값은 별도 단계에서 환경·시나리오·user expectations 등을 정당화한 뒤 정하는 것이 좋다.

---

*메트릭 베이스라인 분석. SLO 임계값 정의는 하지 않음. Review 단계는 보류.*
