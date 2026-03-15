# Experiment notes: Core Signal Validation

**Project:** sre-monitoring-cluster-health-model  
**Phase:** Experiment  
**Input:** Architecture (`03-architecture/architecture.md`), Core signal list (`04-engineering/core-signal-list.md`)

이 실험의 목표는 **제안된 모니터링 신호가 현재 Prometheus/Grafana 환경에서 실제로 사용 가능한지** 확인하는 것이다. 최종 dashboard를 만드는 것이 아니라, **첫 번째 dashboard 반복을 현재 모니터링 환경으로 구축할 수 있는지** 판단하는 데 중점을 둔다.

---

## 1. 실험 범위

전체 신호를 모두 검증하지 않고, **첫 번째 대시보드에 쓸 수 있는 최소 집합** 위주로 검증 대상을 선정했다.

### 1.1 검증 대상 신호 (선택된 부분집합)

| Layer | 선택된 신호 | 비고 |
|-------|-------------|------|
| **Cluster Health Summary** | S1 NotReady node count, S2 API server health, S3 Scheduler pending pods, S4 Workload Pending pod count, S5 Excessive restarts, S6 Critical service endpoint empty | Summary 6개 전부 |
| **Trend / Risk** | T1 Node CPU/memory utilization, T2 API server inflight, T3 Scheduler pending trend, T4 Node disk | T5는 T3와 중복 가능하므로 제외 |
| **Top Offenders** | O1 CPU TOP10 nodes, O2 Memory TOP10 nodes, O3 Restart TOP10 workloads, O4 Pending pods by workload, O5 Error-heavy ingress, O6 Control plane component health | 6개 뷰 |

### 1.2 검증 방법 (실제 환경에서 수행 시)

- Prometheus UI 또는 Grafana Explore에서 해당 PromQL 실행.
- metric 이름 존재 여부: `{__name__=~"metric_name.*"}` 또는 직접 쿼리.
- 라벨 구조: `metric_name` 조회 후 라벨 목록 확인.
- 쿼리 비용: `[5m]`, `[1h]` 등 range selector 사용 시 응답 시간·리소스 사용 확인.
- 결과 해석: 반환값이 “지금 클러스터가 안전한가?” 등 운영 질문에 답하는지 판단.

*아래 검증표는 **일반적인 Kubernetes + Prometheus + kube-state-metrics + node_exporter** 환경을 가정한 예상 결과다. 실제 클러스터에서 실행한 뒤 “Exists?”, “Practical?”, “Known limitations”, “Recommendation”을 갱신하는 것을 권장한다.*

---

## 2. 검증 결과 표 (Validation table)

| Signal | Layer | Query / metric | Exists? | Practical to use? | Known limitations | Recommendation |
|--------|--------|-----------------|--------|-------------------|-------------------|----------------|
| **S1 NotReady node count** | Health Summary | `count(kube_node_status_condition{condition="Ready",status="false"} == 1)` | 예 (kube-state-metrics) | 예 | `status` 라벨 값이 `"true"`/`"false"` 또는 `"True"`/`"False"` 등 버전·exporter에 따라 다를 수 있음. | **즉시 사용 가능.** 환경에서 `kube_node_status_condition` 라벨 값을 확인한 뒤 `status="false"` 또는 해당 값으로 조정. |
| **S2 API server health (latency + 5xx)** | Health Summary | P99: `histogram_quantile(0.99, sum(rate(apiserver_request_duration_seconds_bucket[5m])) by (le))`<br>5xx: `sum(rate(apiserver_request_total{code=~"5.."}[5m]))` | 조건부 (API server scrape 필요) | 조건부 | Managed K8s(EKS, GKE, AKS 등)에서는 control plane을 사용자가 스크래핑하지 못할 수 있음. 그 경우 metric 자체가 없음. Self-managed 또는 API server `/metrics`를 스크래핑하는 환경에서만 존재. | **API server를 스크래핑 중이면 즉시 사용 가능.** 그렇지 않으면 **현재 사용 불가** — Summary에서 제외하거나, managed 서비스가 제공하는 control plane 지표로 대체 검토. |
| **S3 Scheduler pending pods** | Health Summary | `scheduler_pending_pods` (또는 queue별 metric) | 조건부 (scheduler scrape 필요) | 조건부 | kube-scheduler의 `/metrics`를 스크래핑해야 함. Managed K8s에서는 scheduler metric이 노출되지 않는 경우가 많음. | **Self-managed 등 scheduler를 스크래핑 중이면 즉시 사용 가능.** 그렇지 않으면 **현재 사용 불가** — S4 Pending pod count만으로 대체 가능. |
| **S4 Workload Pending pod count** | Health Summary | `count(kube_pod_status_phase{phase="Pending"} == 1)` | 예 (kube-state-metrics) | 예 | kube-state-metrics 표준 metric. `phase` 값이 소문자 `"Pending"` 등으로 일치하는지 확인. | **즉시 사용 가능.** |
| **S5 Excessive restarts** | Health Summary | `sum(increase(kube_pod_container_status_restarts_total[10m]))` | 예 (kube-state-metrics) | 예 | `increase`는 range vector 사용으로 10m 구간마다 평가. retention이 짧으면 오래된 구간은 빈 결과. Eviction은 별도 metric이 없을 수 있음. | **재시작 신호는 즉시 사용 가능.** Eviction은 **별도 metric/이벤트가 있으면 추가**, 없으면 v1에서는 제외. |
| **S5 Eviction** | Health Summary | kubelet eviction 관련 metric (이름 환경별 상이) | 흔히 없음 | 비실용적일 가능성 높음 | kubelet이 eviction 횟수를 Prometheus 형식으로 노출하지 않는 경우가 많음. 이벤트 기반으로만 파악 가능. | **현재 사용 불가**로 간주하고, v1 Summary에서는 재시작만 사용. Eviction은 후속에서 이벤트/다른 수단 검토. |
| **S6 Critical service endpoint empty** | Health Summary | `kube_endpoint_address_available` 또는 endpoint 주소 수 집계 | 조건부 | 수정 시 사용 가능 | kube-state-metrics의 endpoint metric 이름·라벨이 버전마다 다를 수 있음. “엔드포인트가 비어 있다”는 것을 “address 수 0” 또는 “not_ready 수만 있음” 등으로 유도해야 할 수 있음. | **metric 존재 시:** 라벨 구조 확인 후 “비어 있는 엔드포인트” 조건을 PromQL로 정의해 **수정 후 사용.** metric 없으면 **사용 불가.** |
| **T1 Node CPU / memory utilization** | Trend-Risk | CPU: `100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m]))*100)`<br>Memory: `(1 - node_memory_MemAvailable_bytes/node_memory_MemTotal_bytes)*100` | 예 (node_exporter) | 예 | `instance`가 노드명이 아니라 IP:port일 수 있어, Kubernetes 노드명과 매핑하려면 relabel 또는 `node_uname_info` 등과 조인 필요. | **즉시 사용 가능.** 노드명 표시가 필요하면 relabel 또는 조인 추가. |
| **T2 API server inflight** | Trend-Risk | `apiserver_current_inflight_requests` | 조건부 (API server scrape) | S2와 동일 | S2와 동일. Control plane 스크래핑이 되지 않으면 없음. | **S2 사용 가능 시 즉시 사용 가능.** 그렇지 않으면 **사용 불가.** |
| **T3 Scheduler pending trend** | Trend-Risk | `scheduler_pending_pods` 시계열, `increase(scheduler_pending_pods[15m])` 등 | 조건부 (scheduler scrape) | S3와 동일 | S3와 동일. | **S3 사용 가능 시 즉시 사용 가능.** 아니면 S4 Pending count 추세로 대체 가능. |
| **T4 Node disk space** | Trend-Risk | `(node_filesystem_size_bytes - node_filesystem_avail_bytes) / node_filesystem_size_bytes * 100` (mountpoint 필터) | 예 (node_exporter) | 예 | `mountpoint`, `fstype` 등으로 root/중요 볼륨만 필터 필요. 여러 파티션이면 집계 방식 정의 필요. | **즉시 사용 가능.** root 또는 중요 mountpoint만 쿼리하도록 제한 권장. |
| **O1 CPU TOP10 nodes** | Top Offenders | `topk(10, 100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m]))*100))` | 예 (node_exporter) | 예 | T1과 동일하게 instance–노드명 매핑 이슈 가능. | **즉시 사용 가능.** |
| **O2 Memory TOP10 nodes** | Top Offenders | `topk(10, (1 - node_memory_MemAvailable_bytes/node_memory_MemTotal_bytes)*100)` | 예 (node_exporter) | 예 | T1, O1과 동일. | **즉시 사용 가능.** |
| **O3 Restart TOP10 workloads** | Top Offenders | `topk(10, sum by(namespace, pod)(increase(kube_pod_container_status_restarts_total[1h])))` | 예 (kube-state-metrics) | 수정 시 사용 가능 | `pod`는 재시작 수가 많은 파드. deployment/워크로드 단위로 보려면 `kube_pod_owner` 등과 조인하거나 라벨이 있으면 `by (namespace, deployment)` 등으로 변경. | **즉시 사용 가능 (pod 단위).** deployment 등 상위 리소스 단위로 보려면 **수정 후 사용** (owner/라벨 조인). |
| **O4 Pending pods by workload** | Top Offenders | `count by(namespace)(kube_pod_status_phase{phase="Pending"}==1)` 또는 owner/라벨별 집계 | 예 (kube-state-metrics) | 예 | namespace별이면 그대로 사용 가능. deployment 등으로 나누려면 `kube_pod_owner` 등과 조인. | **즉시 사용 가능 (namespace 단위).** 더 세분화하면 수정 후 사용. |
| **O5 Error-heavy ingress / services** | Top Offenders | Ingress controller별 metric (예: nginx, traefik 등) 5xx rate by ingress/service | 조건부 | 조건부 | Ingress controller가 Prometheus metric을 노출해야 함. metric 이름·라벨이 제품마다 다름. | **Ingress controller metric이 있으면 수정 후 사용.** 없으면 **v1에서 사용 불가** 또는 “에러율 높은 서비스” 뷰는 제외. |
| **O6 Control plane component health** | Top Offenders | scheduler, controller-manager 등 `/healthz` 또는 해당 metric | 조건부 | 조건부 | Managed K8s에서는 제어면 컴포넌트를 직접 스크래핑할 수 없음. Self-managed에서만 의미 있음. | **Self-managed이고 스크래핑 중이면 즉시 사용 가능.** 그렇지 않으면 **현재 사용 불가.** |

---

## 3. 구분 요약: 즉시 사용 / 수정 후 사용 / 현재 불가

### 3.1 즉시 사용 가능 (현재 환경에 metric이 있다고 가정할 때)

- **S1** NotReady node count  
- **S4** Workload Pending pod count  
- **S5** Excessive restarts (재시작만; Eviction 제외)  
- **T1** Node CPU/memory utilization  
- **T4** Node disk space  
- **O1** CPU TOP10 nodes  
- **O2** Memory TOP10 nodes  
- **O3** Restart TOP10 (pod 단위)  
- **O4** Pending pods by workload (namespace 단위)

*조건: Prometheus가 kube-state-metrics, node_exporter를 스크래핑하고 있음.*

### 3.2 수정 후 사용 가능

- **S6** Critical service endpoint empty — endpoint metric 이름·라벨 구조에 맞게 “비어 있는 엔드포인트” PromQL 정의.
- **O3** Restart TOP10 — deployment/워크로드 단위로 보려면 owner 또는 라벨 조인 추가.
- **O4** Pending by workload — deployment 등 더 세분화하려면 owner/라벨 조인.
- **O5** Error-heavy ingress — Ingress controller metric이 있을 때, 해당 metric·라벨에 맞게 쿼리 작성.

### 3.3 현재 사용 불가 또는 조건부

- **S2, T2** API server health / inflight — **Managed K8s에서 control plane 미스크래핑 시** metric 없음. Self-managed 또는 API server scrape 환경에서만 사용 가능.
- **S3, T3** Scheduler pending (및 추세) — **scheduler를 스크래핑하지 않으면** metric 없음. S4 Pending count로 대체 가능.
- **S5 Eviction** — eviction 전용 Prometheus metric이 없는 경우가 많음. v1에서는 제외.
- **O5** Error-heavy ingress — Ingress controller metric 미노출 시 사용 불가.
- **O6** Control plane component health — Managed K8s 또는 제어면 미스크래핑 시 사용 불가.

---

## 4. 제한 사항 정리 (Limitations)

### 4.1 Metric 부재

- **Control plane (API server, scheduler):** Managed Kubernetes에서는 사용자가 control plane을 스크래핑하지 못해 `apiserver_*`, `scheduler_pending_pods` 등이 없을 수 있음.  
  → **대응:** Summary에서 S2, S3를 제외하고, S4 Pending count로 “스케줄 불가” 징후를 대표시키거나, managed 서비스가 제공하는 control plane 지표가 있으면 해당 지표로 대체 검토.
- **Eviction:** kubelet이 eviction 횟수를 Prometheus로 노출하지 않는 경우가 많음.  
  → **대응:** v1에서는 재시작만 사용하고, Eviction은 후속에서 이벤트·다른 수단 검토.
- **Ingress/서비스 에러율:** Ingress controller가 Prometheus metric을 노출하지 않으면 O5 사용 불가.  
  → **대응:** v1에서 O5를 제외하거나, 다른 에러 지표(예: 애플리케이션 metric)가 있으면 그에 맞게 정의.

### 4.2 라벨 구조 불일치

- **kube_node_status_condition:** `status` 값이 `"true"`/`"false"` 또는 `"True"`/`"False"` 등으로 다를 수 있음.  
  → **대응:** 실제 metric을 조회해 `status="false"` 등 올바른 값으로 쿼리 수정.
- **node_exporter `instance`:** IP:port 형태라 Kubernetes 노드명과 다를 수 있음.  
  → **대응:** relabel_configs로 노드명 매핑 또는 `node_uname_info` 등과 조인.
- **Endpoint metric:** kube-state-metrics 버전에 따라 endpoint 관련 metric 이름·라벨이 다름.  
  → **대응:** 문서·실제 metric 확인 후 “비어 있는 엔드포인트” 조건을 PromQL로 정의.

### 4.3 쿼리 비용·비실용성

- **histogram_quantile + rate(...[5m]):** API server P99 등은 range 쿼리로 부하가 있을 수 있음. 대시보드 새로고침 주기(예: 1분)를 너무 짧게 하지 않도록 권장.
- **increase(...[1h])** (O3 등): 1h range는 retention이 짧으면 구간 끝이 잘릴 수 있음. retention 내에서만 유의미.
- **topk(10, ...)** 자체는 상대적으로 가벼움. 다만 내부 식이 무거우면(많은 시리즈 스캔) 동일하게 부하 가능.  
  → **대응:** Prometheus retention·스크래핑 주기를 고려해 range 구간 선택; 대시보드 새로고침 주기 1–2분 이상 권장.

### 4.4 Retention 제한

- **increase(...[10m])**, **increase(...[1h])** 는 해당 구간이 retention 안에 있어야 함. retention이 1시간이면 1h range는 최근 1시간만 반영.
- **Trend** 신호(추세)는 일정 기간 시계열이 필요하므로, retention이 매우 짧으면 추세 판단이 어렵다.  
  → **대응:** retention을 확인하고, 필요 시 range 구간을 줄이거나 “현재 값만” threshold로 사용.

### 4.5 이론적으로 유용하나 현재 비실용적인 신호

- **S2/S3/T2/T3/O6 (control plane 전반):** Managed K8s에서 스크래핑 불가 시 모두 비실용.
- **S5 Eviction:** metric 부재로 v1에서는 비실용.
- **O5 (Error-heavy ingress):** Ingress controller metric 부재 시 비실용.

---

## 5. 첫 번째 대시보드 반복 가능 여부

### 5.1 결론

- **kube-state-metrics + node_exporter**만 있어도 **Summary 일부(S1, S4, S5 재시작), Trend(T1, T4), Top Offenders(O1, O2, O3, O4)** 는 구축 가능하다.
- **Control plane(API server, scheduler)** metric이 없으면 Summary에서 S2, S3를 빼고, **S4 Pending count**로 “스케줄·워크로드 관점 건강”을 대표시키는 구성이 현실적이다.
- **S6(endpoint empty)** 는 metric 존재·라벨 확인 후 수정하면 사용 가능하다.
- **Eviction, O5, O6** 은 환경에 따라 v1에서 제외하고, 후속 반복에서 metric 도입·대체 수단을 검토하면 된다.

따라서 **현재 모니터링 환경만으로도 첫 번째 Cluster Health 대시보드 반복은 가능**하다. 다만 control plane metric 유무에 따라 Summary 패널 구성(5~6개 중 4~6개 사용)을 조정하고, “즉시 사용 가능” vs “수정 후 사용” vs “현재 불가” 구분에 따라 우선순위를 두어 구현하는 것이 좋다.

### 5.2 권장 다음 단계

1. 실제 Prometheus에서 위 표의 쿼리를 실행해 “Exists?”, “Practical?”, “Known limitations”를 채우고 이 문서를 갱신.
2. Control plane(API server, scheduler) 스크래핑 가능 여부를 확인하고, 불가 시 Summary에서 S2, S3 제외 및 S4 강조 방안 확정.
3. Endpoint metric(S6)의 실제 이름·라벨을 확인한 뒤 “비어 있는 엔드포인트” PromQL을 작성해 core-signal-list에 반영.
4. **sre-monitoring-dashboard-design** 단계에서 “즉시 사용 가능” + “수정 후 사용” 신호만으로 첫 번째 central dashboard 레이아웃을 설계.

---

*Experiment phase output. 다음 단계: Review — “5–10분 안에 클러스터 건강 파악 가능한가?” 평가 및 모델·신호 목록 검토.*
