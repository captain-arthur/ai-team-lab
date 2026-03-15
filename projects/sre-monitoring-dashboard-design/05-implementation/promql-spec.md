# PromQL Specification — Central Kubernetes Operational Dashboard v1

**Project:** sre-monitoring-dashboard-design  
**Phase:** 05-implementation  
**Input:** implementation-ready-panel-spec, panel-design

이 문서는 **A 우선순위 패널만** 대상으로, **배포 가능한 수준의 PromQL**과 필요 metric·라벨·집계·구간·쿼리 비용을 정의한다.  
일반적인 Kubernetes 모니터링 스택: **Prometheus**, **kube-state-metrics**, **node_exporter**. Ingress controller metric은 선택.

---

## 1. 전제 조건

- **데이터소스:** Prometheus (datasource 이름은 환경에 맞게 `Prometheus` 또는 변수로 지정).
- **메트릭 소스:** kube-state-metrics, node_exporter가 노드·파드에 배포되어 있다고 가정.
- **비용:** 극단적으로 비싼 쿼리(긴 range, 대량 시리즈 스캔)는 사용하지 않는다.

---

## 2. Block 1 — Operational Confidence (A 우선순위 4개)

### 2.1 NotReady node count

| 항목 | 내용 |
|------|------|
| **Panel name** | NotReady node count |
| **Production-ready PromQL** | `count(kube_node_status_condition{condition="Ready",status="false"} == 1)` |
| **Required metrics** | `kube_node_status_condition` (kube-state-metrics) |
| **Required labels** | `condition`, `status` |
| **Aggregation strategy** | 클러스터 전체 count. 추가 필터(예: node role)는 환경별로 적용. |
| **Time range** | instant query (기본). range 불필요. |
| **Expected query cost** | **low** |

---

### 2.2 Workload Pending pod count

| 항목 | 내용 |
|------|------|
| **Panel name** | Workload Pending pod count |
| **Production-ready PromQL** | `count(kube_pod_status_phase{phase="Pending"} == 1)` |
| **Required metrics** | `kube_pod_status_phase` (kube-state-metrics) |
| **Required labels** | `phase` |
| **Aggregation strategy** | 클러스터 전체 Pending 파드 수. namespace 필터는 변수로 선택 시 적용. |
| **Time range** | instant query. |
| **Expected query cost** | **low** |

---

### 2.3 Excessive restarts

| 항목 | 내용 |
|------|------|
| **Panel name** | Excessive restarts |
| **Production-ready PromQL** | `sum(increase(kube_pod_container_status_restarts_total[10m]))` |
| **Required metrics** | `kube_pod_container_status_restarts_total` (kube-state-metrics) |
| **Required labels** | (없어도 동작. namespace/pod 등은 선택) |
| **Aggregation strategy** | 클러스터 전체 10분 내 재시작 횟수 합계. **Threshold N은 반드시 팀에서 정의** 후 Grafana threshold에 반영(예: N=10~20). |
| **Time range** | `[10m]` (increase). Prometheus retention 내에서만 유의미. |
| **Expected query cost** | **low** (10m range, 시리즈 수에 비례) |

---

### 2.4 Critical service endpoint empty

| 항목 | 내용 |
|------|------|
| **Panel name** | Critical service endpoint empty |
| **Production-ready PromQL** | **환경별 선택.** 단수형: `count(kube_endpoint_address_available == 0)`. 복수형(첫 검증에서 동작): `count(kube_endpoints_address_available == 0)`. 핵심 네임스페이스만: `count(kube_endpoints_address_available{namespace=~"default\|production"} == 0)` |
| **Required metrics** | `kube_endpoint_address_available` 또는 `kube_endpoints_address_available` (kube-state-metrics 버전에 따라 상이). |
| **Required labels** | `namespace`, `endpoint` (metric별로 상이. 환경에서 확인) |
| **Aggregation strategy** | “available 주소가 0인 endpoint” 개수. critical만 보려면 `namespace` 필터. |
| **Time range** | instant query. |
| **Expected query cost** | **low** |

**환경별 metric 이름 (첫 검증 반영):** kube-state-metrics 버전에 따라 **단수** `kube_endpoint_address_available` 또는 **복수** `kube_endpoints_address_available` 사용. No data/Query error 시 `curl <kube-state-metrics>/metrics | grep -E "endpoint.*available"` 로 실제 이름 확인 후 교체. grafana-dashboard-v1.json 기본값은 복수형(검증 통과 환경).

---

## 3. Block 2 — Early Risk (A 우선순위 4개)

### 3.1 Node CPU utilization

| 항목 | 내용 |
|------|------|
| **Panel name** | Node CPU utilization |
| **Production-ready PromQL** | 클러스터 평균: `avg(100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100))`  
  또는 80% 초과 노드 수: `count(100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80)` |
| **Required metrics** | `node_cpu_seconds_total` (node_exporter) |
| **Required labels** | `mode` (idle 필수), `instance` |
| **Aggregation strategy** | `avg by(instance)(rate(...[5m]))` 로 노드별 CPU 사용률 후, 클러스터 평균 또는 count(>80). |
| **Time range** | `[5m]` (rate). |
| **Expected query cost** | **low** |

---

### 3.2 Node memory / OOM risk

| 항목 | 내용 |
|------|------|
| **Panel name** | Node memory / OOM risk |
| **Production-ready PromQL** | 클러스터 평균 사용률: `avg((1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100)`  
  또는 80% 초과 노드 수: `count((1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100 > 80)` |
| **Required metrics** | `node_memory_MemAvailable_bytes`, `node_memory_MemTotal_bytes` (node_exporter) |
| **Required labels** | `instance` |
| **Aggregation strategy** | 노드별 (1 - MemAvailable/MemTotal)*100 후 avg 또는 count(>80). |
| **Time range** | instant query (스냅샷). |
| **Expected query cost** | **low** |

---

### 3.3 Node disk space

| 항목 | 내용 |
|------|------|
| **Panel name** | Node disk space |
| **Production-ready PromQL** | 클러스터 평균 사용률 (root): `avg((1 - (node_filesystem_avail_bytes{mountpoint="/",fstype!~"tmpfs|overlay"} / node_filesystem_size_bytes{mountpoint="/",fstype!~"tmpfs|overlay"})) * 100)`. **root가 다른 경로인 환경:** `mountpoint="/var"` 등으로 변경. |
| **Required metrics** | `node_filesystem_avail_bytes`, `node_filesystem_size_bytes` (node_exporter) |
| **Required labels** | `mountpoint`, `instance`. `fstype` 필터로 tmpfs/overlay 제외 권장. |
| **Aggregation strategy** | **mountpoint는 환경별 조정.** 기본 "/". 일부 노드/OS는 root가 `/var` 등이면 해당 값으로 교체. |
| **Time range** | instant query. |
| **Expected query cost** | **low** |

---

### 3.4 Pending pods trend

| 항목 | 내용 |
|------|------|
| **Panel name** | Pending pods trend |
| **Production-ready PromQL** | 현재값 (Stat용): `count(kube_pod_status_phase{phase="Pending"} == 1)`  
  시계열용 (Time series): 동일 쿼리를 시계열로 표시.  
  또는 15분 증가량: `increase(count(kube_pod_status_phase{phase="Pending"}==1)[15m])` — 주의: count는 instant이므로 increase는 시계열 컨텍스트에서만 유의미. 실무에서는 “현재 Pending 수”를 시계열로 그리는 것이 일반적. |
| **Required metrics** | `kube_pod_status_phase` (kube-state-metrics) |
| **Required labels** | `phase` |
| **Aggregation strategy** | 클러스터 전체 Pending 수. Time series 패널이면 시간에 따른 동일 instant 쿼리 반복. |
| **Time range** | instant (Stat) 또는 대시보드 time range (Time series). |
| **Expected query cost** | **low** (Stat), **low–medium** (Time series, 구간에 따라) |

**권장:** v1에서는 **현재 Pending 수**를 **Time series**로 표시해 “추세”를 보이거나, **Stat**으로 현재값만 표시 후 threshold(예: >0 경고) 적용.

---

## 4. Block 3 — Investigation / Top Offenders (A 우선순위 4개)

### 4.1 CPU TOP10 nodes

| 항목 | 내용 |
|------|------|
| **Panel name** | CPU TOP10 nodes |
| **Production-ready PromQL** | `topk(10, 100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100))` |
| **Required metrics** | `node_cpu_seconds_total` (node_exporter) |
| **Required labels** | `mode`, `instance` |
| **Aggregation strategy** | 노드별 CPU 사용률 계산 후 topk(10). Table에서는 instance + value 컬럼으로 표시. |
| **Time range** | `[5m]` (rate). |
| **Expected query cost** | **low** |

---

### 4.2 Memory TOP10 nodes

| 항목 | 내용 |
|------|------|
| **Panel name** | Memory TOP10 nodes |
| **Production-ready PromQL** | `topk(10, (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100)` |
| **Required metrics** | `node_memory_MemAvailable_bytes`, `node_memory_MemTotal_bytes` (node_exporter) |
| **Required labels** | `instance` |
| **Aggregation strategy** | 노드별 메모리 사용률 후 topk(10). |
| **Time range** | instant query. |
| **Expected query cost** | **low** |

---

### 4.3 Restart TOP10 pods

| 항목 | 내용 |
|------|------|
| **Panel name** | Restart TOP10 pods |
| **Production-ready PromQL** | `topk(10, sum by(namespace, pod)(increase(kube_pod_container_status_restarts_total[1h])))` |
| **Required metrics** | `kube_pod_container_status_restarts_total` (kube-state-metrics) |
| **Required labels** | `namespace`, `pod` (container 단위이면 sum by로 pod로 묶음) |
| **Aggregation strategy** | namespace, pod별 1h 재시작 증가량 합계 후 topk(10). |
| **Time range** | `[1h]`. Prometheus retention이 1h 이상이어야 유의미. |
| **Expected query cost** | **medium** (1h range, 시리즈 수 많을 수 있음) |

---

### 4.4 Pending pods by workload

| 항목 | 내용 |
|------|------|
| **Panel name** | Pending pods by workload |
| **Production-ready PromQL** | `count by(namespace)(kube_pod_status_phase{phase="Pending"} == 1)`  
  (정렬은 Grafana Table 옵션에서 value 내림차순) |
| **Required metrics** | `kube_pod_status_phase` (kube-state-metrics) |
| **Required labels** | `namespace`, `phase` |
| **Aggregation strategy** | namespace별 Pending 파드 수. deployment 등으로 세분화하려면 `kube_pod_owner` 등과 조인(환경별). |
| **Time range** | instant query. |
| **Expected query cost** | **low** |

---

## 5. 요약 표

| Block | Panel name | Primary metric(s) | Time range | Query cost |
|-------|------------|-------------------|------------|------------|
| 1 | NotReady node count | kube_node_status_condition | instant | low |
| 1 | Workload Pending pod count | kube_pod_status_phase | instant | low |
| 1 | Excessive restarts | kube_pod_container_status_restarts_total | 10m | low |
| 1 | Critical service endpoint empty | kube_endpoint_address_available | instant | low |
| 2 | Node CPU utilization | node_cpu_seconds_total | 5m | low |
| 2 | Node memory / OOM risk | node_memory_* | instant | low |
| 2 | Node disk space | node_filesystem_* | instant | low |
| 2 | Pending pods trend | kube_pod_status_phase | instant / dashboard | low |
| 3 | CPU TOP10 nodes | node_cpu_seconds_total | 5m | low |
| 3 | Memory TOP10 nodes | node_memory_* | instant | low |
| 3 | Restart TOP10 pods | kube_pod_container_status_restarts_total | 1h | medium |
| 3 | Pending pods by workload | kube_pod_status_phase | instant | low |

---

*PromQL spec v1. A-priority panels only. B/C 패널 및 환경별 변형(control plane, ingress metric)은 implementation-notes.md 참고.*
