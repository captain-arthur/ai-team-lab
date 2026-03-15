# Implementation-Ready Panel Specification — Central Kubernetes Operational Dashboard v1

**Project:** sre-monitoring-dashboard-design  
**Phase:** Engineering (implementation-ready)  
**Input:** Architecture, Panel design, Final design document

이 문서는 **Grafana + Prometheus** 로 v1 대시보드를 구현할 때 필요한 **패널 단위 명세**를 정리한 것이다. 구현자가 추가 설계 없이 패널을 만들 수 있도록, 패널 이름·목적·레이어·PromQL 예시·패널 타입·쿼리 비용·환경 의존성을 명시한다.  
**설명은 한국어**, **기술 식별자(metric, PromQL, Grafana 패널 타입)는 영어**로 통일한다.

---

## 1. 목표와 제약

### 1.1 운영 목표

대시보드는 다음 두 가지에 **확신 있게** 답할 수 있어야 한다.

1. **지금 클러스터가 안전한가?** (Operational Confidence)  
2. **클러스터가 곧 불안전해질 조기 징후가 보이는가?** (Early Risk)

조기 리스크로 특히 다루는 신호: **ingress/nginx-controller stress**, **CPU throttling**, **OOM risk**, **node resource pressure**, **pending pods**.

### 1.2 v1 쿼리 제약

- Prometheus 리소스는 무제한이 아니므로, **비용이 크거나 긴 구간 쿼리는 v1에서 피한다.**  
- 각 패널에 **expected query cost (low / medium / high)** 를 부여하여, high 비용 패널은 선택(Optional) 또는 제한된 새로고침으로 구현한다.

---

## 2. 패널 우선순위 분류 (Prioritized Panel List)

모든 후보 패널을 다음 세 등급으로 분류한다.

| 등급 | 의미 | v1 메인 뷰 포함 |
|------|------|------------------|
| **A. Required for v1** | 대부분의 클러스터에서 사용 가능하고, “지금 안전한가?” / “조기 징후가 있는가?”에 필수. | **반드시 포함** (환경에서 metric 있으면). |
| **B. Optional for v1** | control plane / ingress / throttling 등 특정 metric에 의존. 있으면 포함, 없으면 생략. | **환경 가능 시 포함**. |
| **C. Deferred / future** | v1 이후(alert·runbook·추가 신호)에서 다룰 수 있는 항목. | v1에서 구현하지 않거나, 드릴다운만 최소 구현. |

### 2.1 Block 1: Operational Confidence

| Panel name | 우선순위 | 비고 |
|------------|----------|------|
| NotReady node count | **A** | 대부분 클러스터에서 kube-state-metrics만으로 가능. |
| Workload Pending pod count | **A** | 동일. |
| Excessive restarts | **A** | 동일. |
| Critical service endpoint empty | **A** | endpoint metric 이름만 환경별 조정. |
| API server health | **B** | control plane metric 필요. Managed K8s에서는 보통 미제공. |
| Scheduler pending pods | **B** | control plane metric 필요. |

### 2.2 Block 2: Early Risk

| Panel name | 우선순위 | 비고 |
|------------|----------|------|
| Node CPU utilization | **A** | node_exporter 표준. |
| Node memory / OOM risk | **A** | 동일. |
| Node disk space | **A** | 동일. |
| Pending pods trend | **A** | Pending count 또는 짧은 구간 추세. |
| CPU throttling risk | **B** | 컨테이너 CPU throttling metric 필요(제품·버전 의존). |
| Ingress stress | **B** | Ingress controller / ingressgateway / nginx metric 필요. |

### 2.3 Block 3: Investigation / Top Offenders

| Panel name | 우선순위 | 비고 |
|------------|----------|------|
| CPU TOP10 nodes | **A** | Block 2 Node CPU 이상 시 드릴다운. |
| Memory TOP10 nodes | **A** | Block 2 OOM risk 이상 시 드릴다운. |
| Restart TOP10 pods | **A** | Block 1 Excessive restarts 시 드릴다운. |
| Pending pods by workload | **A** | Block 1 Pending 시 드릴다운. |
| Error-heavy ingress / services | **B** | Ingress metric 필요. |
| Control plane component health | **B** | Self-managed 등 control plane metric 있을 때만. |

---

## 3. 최소 패널 집합 (Minimal Set for Operational Confidence)

운영자가 **“지금 클러스터가 안전하다”** / **“조기 경고 징후가 보인다”** 를 확신 있게 말하려면, 아래 **최소 집합**이 메인 뷰에 반드시 있어야 한다.

### 3.1 “지금 안전한가?” (Current safety confidence)

**Block 1에서 다음 4개는 필수(A).**  
control plane metric이 있는 환경에서는 API server health, Scheduler pending 2개를 추가하면 6개.

| # | Panel name | 우선순위 |
|---|------------|----------|
| 1 | NotReady node count | A |
| 2 | Workload Pending pod count | A |
| 3 | Excessive restarts | A |
| 4 | Critical service endpoint empty | A |
| 5 | API server health | B |
| 6 | Scheduler pending pods | B |

**판단 규칙:** 위 패널이 **전부 정상(threshold 이내)** 이면 “클러스터 안전”. **하나라도 비정상**이면 “조사 필요”.

### 3.2 “조기 경고 징후가 있는가?” (Early warning / future risk confidence)

**Block 2에서 다음 4개는 필수(A).**  
CPU throttling, Ingress stress는 metric 있을 때 Optional(B)로 추가.

| # | Panel name | 우선순위 |
|---|------------|----------|
| 1 | Node CPU utilization | A |
| 2 | Node memory / OOM risk | A |
| 3 | Node disk space | A |
| 4 | Pending pods trend | A |
| 5 | CPU throttling risk | B |
| 6 | Ingress stress | B |

**판단 규칙:** Block 2에서 **경고/위험(예: 80% 초과, 여유 &lt;10%)** 이 하나라도 있으면 “곧 불안정해질 수 있음” → 용량·정리·확장 검토 또는 Block 3 드릴다운.

### 3.3 요약

- **최소 메인 뷰:** Block 1 필수 4개 + Block 2 필수 4개 = **8개 패널**로 “안전한가?” / “조기 징후가 있는가?”에 답 가능.  
- **권장 메인 뷰:** Block 1에 B 2개(API server, Scheduler pending) 포함 시 6개, Block 2에 B 2개(CPU throttling, Ingress stress) 포함 시 6개 → **최대 12개**.  
- Block 3(Investigation)은 **드릴다운 전용**이며, “최소 집합”에는 메인 뷰 패널 수에 포함하지 않는다.

---

## 4. 패널별 상세 명세 (Panel Specification)

아래는 **각 패널마다** 구현에 필요한 항목을 정리한 것이다.  
**표기:**  
- **Signal type:** `current safety confidence` / `future risk confidence` / `investigation`  
- **Query cost:** `low` / `medium` / `high`  
- **Environment dependency:** `works in most clusters` / `requires control plane metrics` / `requires ingress metrics` / `requires workload-specific labels`

---

### 4.1 Block 1: Operational Confidence

#### P1 — NotReady node count

| 항목 | 내용 |
|------|------|
| **Panel name** | NotReady node count |
| **Purpose** | Ready가 아닌 노드가 0인지 확인. 0이면 스케줄 용량 유지. |
| **Layer** | Operational Confidence |
| **Prometheus metric / PromQL example** | `count(kube_node_status_condition{condition="Ready",status="false"} == 1)` |
| **Panel type** | Stat |
| **Why it matters** | NotReady 노드가 있으면 용량 감소·스케줄 실패 가능. 한 눈에 “정상/비정상” 판단. |
| **Signal type** | current safety confidence |
| **Expected query cost** | low |
| **Environment dependency** | works in most clusters (kube-state-metrics) |

**Grafana:** Value mappings — 0 = OK (green), &gt;0 = Critical (red). Threshold: 0.

---

#### P2 — API server health

| 항목 | 내용 |
|------|------|
| **Panel name** | API server health |
| **Purpose** | API 서버 응답·에러 여부. P99 &lt;1s, 5xx 극소 = 정상. |
| **Layer** | Operational Confidence |
| **Prometheus metric / PromQL example** | P99: `histogram_quantile(0.99, sum(rate(apiserver_request_duration_seconds_bucket[5m])) by (le))`<br>5xx: `sum(rate(apiserver_request_total{code=~"5.."}[5m]))` |
| **Panel type** | Stat (또는 P99 + 5xx 두 개 Stat) |
| **Why it matters** | 제어면 비정상 시 클러스터 전체에 영향. |
| **Signal type** | current safety confidence |
| **Expected query cost** | medium (histogram_quantile + rate 5m) |
| **Environment dependency** | requires control plane metrics |

**Grafana:** Threshold 예: P99 &lt;1s 정상. 5xx &lt; 0.1/s 정상. Refresh 1–2분 권장.

---

#### P3 — Scheduler pending pods

| 항목 | 내용 |
|------|------|
| **Panel name** | Scheduler pending pods |
| **Purpose** | 스케줄러 대기 파드 수. 0이어야 워크로드를 받을 수 있음. |
| **Layer** | Operational Confidence |
| **Prometheus metric / PromQL example** | `scheduler_pending_pods` (또는 환경에 맞는 metric 이름) |
| **Panel type** | Stat |
| **Why it matters** | 스케줄러가 처리하지 못하는 파드가 있으면 용량·리소스 이슈. |
| **Signal type** | current safety confidence |
| **Expected query cost** | low |
| **Environment dependency** | requires control plane metrics |

**Grafana:** 0 = OK, &gt;0 = Critical.

---

#### P4 — Workload Pending pod count

| 항목 | 내용 |
|------|------|
| **Panel name** | Workload Pending pod count |
| **Purpose** | 사용자 관점 “기동 안 되는 파드” 수. 0 = 정상. |
| **Layer** | Operational Confidence |
| **Prometheus metric / PromQL example** | `count(kube_pod_status_phase{phase="Pending"} == 1)` |
| **Panel type** | Stat |
| **Why it matters** | Pending이 있으면 스케줄 실패·리소스 부족·이미지 풀 등 이슈 가능. |
| **Signal type** | current safety confidence |
| **Expected query cost** | low |
| **Environment dependency** | works in most clusters |

**Grafana:** 0 = OK, &gt;0 = Critical.

---

#### P5 — Excessive restarts

| 항목 | 내용 |
|------|------|
| **Panel name** | Excessive restarts |
| **Purpose** | 10분 내 과도한 재시작 여부. 임계치 이하 = 정상. |
| **Layer** | Operational Confidence |
| **Prometheus metric / PromQL example** | `sum(increase(kube_pod_container_status_restarts_total[10m]))` (threshold N은 팀 정의) |
| **Panel type** | Stat |
| **Why it matters** | 재시작 폭증은 워크로드 불안정·OOM·이미지 문제 등. |
| **Signal type** | current safety confidence |
| **Expected query cost** | low (10m increase, 단 retention 고려) |
| **Environment dependency** | works in most clusters |

**Grafana:** Threshold N 미만 = OK, N 이상 = Critical. N은 팀에서 정의(예: 10).

---

#### P6 — Critical service endpoint empty

| 항목 | 내용 |
|------|------|
| **Panel name** | Critical service endpoint empty |
| **Purpose** | 핵심 서비스 백엔드가 비어 있지 않은지. 비어 있지 않음 = 정상. |
| **Layer** | Operational Confidence |
| **Prometheus metric / PromQL example** | `kube_endpoint_address_available` 등으로 핵심 서비스 필터 후, 비어 있는 엔드포인트 수 또는 0/1. 예: `count(kube_endpoints_address_not_ready{namespace="default"} > 0)` 등 환경에 맞게 조정. |
| **Panel type** | Stat |
| **Why it matters** | 엔드포인트가 비면 트래픽 불가. |
| **Signal type** | current safety confidence |
| **Expected query cost** | low |
| **Environment dependency** | works in most clusters (endpoint metric 이름·라벨은 환경별 확인) |

**Grafana:** 0 = OK, &gt;0 = Critical. 대상 서비스는 변수 또는 하드코드.

---

### 4.2 Block 2: Early Risk

#### T1 — Node CPU utilization

| 항목 | 내용 |
|------|------|
| **Panel name** | Node CPU utilization |
| **Purpose** | 노드 CPU 포화 직전(예: 80% 이상) 여부. 조기 리스크. |
| **Layer** | Early Risk |
| **Prometheus metric / PromQL example** | 클러스터 평균: `100 - (avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)`<br>80% 초과 노드 수: `count(100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m]))*100) > 80)` |
| **Panel type** | Gauge 또는 Stat |
| **Why it matters** | CPU 포화 시 eviction·성능 저하 리스크. |
| **Signal type** | future risk confidence |
| **Expected query cost** | low |
| **Environment dependency** | works in most clusters (node_exporter) |

**Grafana:** Thresholds: 0–80 green, 80–95 yellow, 95–100 red.

---

#### T2 — Node memory / OOM risk

| 항목 | 내용 |
|------|------|
| **Panel name** | Node memory / OOM risk |
| **Purpose** | 메모리 압박 = OOM·eviction 위험. “OOM risk” 강조. |
| **Layer** | Early Risk |
| **Prometheus metric / PromQL example** | 사용률: `(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100`<br>80% 초과 노드 수: `count((1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100 > 80)` |
| **Panel type** | Gauge 또는 Stat |
| **Why it matters** | 메모리 부족 시 OOM kill·eviction. |
| **Signal type** | future risk confidence |
| **Expected query cost** | low |
| **Environment dependency** | works in most clusters (node_exporter) |

**Grafana:** Thresholds: 0–80 green, 80–95 yellow, 95–100 red.

---

#### T3 — Node disk space

| 항목 | 내용 |
|------|------|
| **Panel name** | Node disk space |
| **Purpose** | 디스크 사용률. 여유 &lt;10% 등이면 경고. |
| **Layer** | Early Risk |
| **Prometheus metric / PromQL example** | 사용률: `(1 - (node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"})) * 100`<br>또는 `(node_filesystem_size_bytes - node_filesystem_avail_bytes) / node_filesystem_size_bytes * 100` (mountpoint="/" 등) |
| **Panel type** | Gauge 또는 Stat |
| **Why it matters** | 디스크 부족 시 이미지 풀·로그 실패·노드 불안정. |
| **Signal type** | future risk confidence |
| **Expected query cost** | low |
| **Environment dependency** | works in most clusters (node_exporter) |

**Grafana:** Thresholds: 0–90 green, 90–95 yellow, 95–100 red (또는 “여유 10%” 기준으로 반전).

---

#### T4 — Pending pods trend

| 항목 | 내용 |
|------|------|
| **Panel name** | Pending pods trend |
| **Purpose** | Pending 수 또는 짧은 구간 추세. 증가 시 스케줄·용량 리스크. |
| **Layer** | Early Risk |
| **Prometheus metric / PromQL example** | 현재값: `count(kube_pod_status_phase{phase="Pending"} == 1)`<br>추세(15분 증가): `increase(count(kube_pod_status_phase{phase="Pending"}==1)[15m])` (또는 time-series로 같은 쿼리) |
| **Panel type** | Time series 또는 Stat |
| **Why it matters** | Pending이 늘어나면 곧 용량·스케줄 문제로 이어질 수 있음. |
| **Signal type** | future risk confidence |
| **Expected query cost** | low (Stat) / low–medium (time-series 짧은 구간) |
| **Environment dependency** | works in most clusters |

**Grafana:** Time series면 “증가 추세” 가시화. Stat이면 “현재 Pending 수” + threshold(예: &gt;0 경고).

---

#### T5 — CPU throttling risk (Optional)

| 항목 | 내용 |
|------|------|
| **Panel name** | CPU throttling risk |
| **Purpose** | 컨테이너가 CPU limit에 막혀 throttling 되는지. 있으면 성능 저하·조기 리스크. |
| **Layer** | Early Risk |
| **Prometheus metric / PromQL example** | 환경 의존. cAdvisor: `container_cpu_cfs_throttled_seconds_total`, `container_cpu_cfs_periods_total` 등으로 throttling 비율 계산. 예: `sum(rate(container_cpu_cfs_throttled_seconds_total[5m])) / sum(rate(container_cpu_cfs_periods_total[5m])) * 100` (해당 metric 있을 때). |
| **Panel type** | Stat 또는 Gauge |
| **Why it matters** | CPU throttling이 높으면 워크로드 지연·불만. |
| **Signal type** | future risk confidence |
| **Expected query cost** | low–medium |
| **Environment dependency** | requires workload-specific labels / cAdvisor 또는 해당 metric 노출 환경 |

**Grafana:** metric 없으면 패널 비표시 또는 Node CPU로 대체 안내.

---

#### T6 — Ingress stress (Optional)

| 항목 | 내용 |
|------|------|
| **Panel name** | Ingress stress |
| **Purpose** | Ingress / ingressgateway / nginx-controller 부하·지연·에러율. 높으면 트래픽 경로 리스크. |
| **Layer** | Early Risk |
| **Prometheus metric / PromQL example** | 제품별 상이. Istio: `istio_request_duration_milliseconds`, `istio_requests_total`. Nginx: `nginx_ingress_controller_requests`, `nginx_ingress_controller_request_duration_seconds`. 예: P99 지연 `histogram_quantile(0.99, sum(rate(istio_request_duration_milliseconds_bucket[5m])) by (le))` 또는 5xx rate. |
| **Panel type** | Stat 또는 Gauge |
| **Why it matters** | Ingress 부하·에러가 높으면 사용자 영향. |
| **Signal type** | future risk confidence |
| **Expected query cost** | medium (histogram_quantile 사용 시) |
| **Environment dependency** | requires ingress metrics |

**Grafana:** 환경에 맞는 metric 이름으로 교체. 없으면 패널 비표시.

---

### 4.3 Block 3: Investigation / Top Offenders

#### O1 — CPU TOP10 nodes

| 항목 | 내용 |
|------|------|
| **Panel name** | CPU TOP10 nodes |
| **Purpose** | Block 2 “Node CPU 높음”일 때 어느 노드가 1위인지 드릴다운. |
| **Layer** | Investigation |
| **Prometheus metric / PromQL example** | `topk(10, 100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100))` |
| **Panel type** | Table |
| **Why it matters** | 원인 노드 빠르게 좁히기. |
| **Signal type** | investigation |
| **Expected query cost** | low |
| **Environment dependency** | works in most clusters |

**Grafana:** 컬럼: instance (또는 node), value (%). 정렬: value 내림차순.

---

#### O2 — Memory TOP10 nodes

| 항목 | 내용 |
|------|------|
| **Panel name** | Memory TOP10 nodes |
| **Purpose** | Block 2 “OOM risk”일 때 어느 노드가 메모리 압박 1위인지. |
| **Layer** | Investigation |
| **Prometheus metric / PromQL example** | `topk(10, (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100)` |
| **Panel type** | Table |
| **Why it matters** | OOM·eviction 원인 노드 파악. |
| **Signal type** | investigation |
| **Expected query cost** | low |
| **Environment dependency** | works in most clusters |

**Grafana:** 컬럼: instance, value (%). 정렬: value 내림차순.

---

#### O3 — Restart TOP10 pods

| 항목 | 내용 |
|------|------|
| **Panel name** | Restart TOP10 pods |
| **Purpose** | Block 1 “Excessive restarts”일 때 어떤 워크로드가 재시작 1위인지. |
| **Layer** | Investigation |
| **Prometheus metric / PromQL example** | `topk(10, sum by(namespace, pod)(increase(kube_pod_container_status_restarts_total[1h])))` (deployment 단위로 하려면 `kube_pod_owner` 등과 조인) |
| **Panel type** | Table |
| **Why it matters** | 재시작 원인 워크로드 좁히기. |
| **Signal type** | investigation |
| **Expected query cost** | medium (1h increase, retention 고려) |
| **Environment dependency** | works in most clusters; deployment 단위는 requires workload-specific labels |

**Grafana:** 컬럼: namespace, pod (또는 deployment), restarts. 정렬: restarts 내림차순.

---

#### O4 — Pending pods by workload

| 항목 | 내용 |
|------|------|
| **Panel name** | Pending pods by workload |
| **Purpose** | Block 1 “Pending 있음”일 때 어떤 네임스페이스/워크로드가 Pending인지. |
| **Layer** | Investigation |
| **Prometheus metric / PromQL example** | `count by(namespace)(kube_pod_status_phase{phase="Pending"} == 1)` (정렬) 또는 owner/라벨별로 by 추가. |
| **Panel type** | Table |
| **Why it matters** | Pending 원인 워크로드·네임스페이스 파악. |
| **Signal type** | investigation |
| **Expected query cost** | low |
| **Environment dependency** | works in most clusters; workload 단위는 requires workload-specific labels |

**Grafana:** 컬럼: namespace (및 workload), pending_count. 정렬: pending_count 내림차순.

---

#### O5 — Error-heavy ingress / services (Optional)

| 항목 | 내용 |
|------|------|
| **Panel name** | Error-heavy ingress / services |
| **Purpose** | 에러율이 높은 ingress·서비스 TOP10. |
| **Layer** | Investigation |
| **Prometheus metric / PromQL example** | 환경별. 예: `topk(10, sum by(destination_service_name)(rate(istio_requests_total{response_code=~"5.."}[5m])) / sum by(destination_service_name)(rate(istio_requests_total[5m])) * 100)` |
| **Panel type** | Table |
| **Why it matters** | Ingress/서비스별 에러 원인 좁히기. |
| **Signal type** | investigation |
| **Expected query cost** | medium |
| **Environment dependency** | requires ingress metrics |

**Grafana:** metric 있을 때만 표시.

---

#### O6 — Control plane component health (Optional)

| 항목 | 내용 |
|------|------|
| **Panel name** | Control plane component health |
| **Purpose** | API server 비정상일 때 제어면 중 어떤 컴포넌트가 비정상인지. |
| **Layer** | Investigation |
| **Prometheus metric / PromQL example** | scheduler, controller-manager 등 `/healthz` 또는 해당 metric별 상태. 환경별. |
| **Panel type** | Table |
| **Why it matters** | Self-managed에서 제어면 장애 원인 좁히기. |
| **Signal type** | investigation |
| **Expected query cost** | low–medium |
| **Environment dependency** | requires control plane metrics |

**Grafana:** Self-managed 등에서만 사용. Managed K8s에서는 비표시.

---

## 5. 환경별 적용 요약

| 환경 | Block 1 | Block 2 | Block 3 |
|------|---------|---------|---------|
| **대부분 클러스터 (kube-state-metrics + node_exporter)** | P1, P4, P5, P6 (4개) | T1, T2, T3, T4 (4개) | O1, O2, O3, O4 (4개) |
| **Control plane metric 있음** | + P2, P3 → 6개 | 동일 | + O6 → 5개 |
| **Ingress metric 있음** | 동일 | + T6 → 5~6개 | + O5 |
| **CPU throttling metric 있음** | 동일 | + T5 | 동일 |

---

## 6. 구현 시 체크리스트 (Grafana Implementer)

- [ ] **메인 뷰:** Block 1 Row + Block 2 Row만 기본 표시. 패널 수 8~12개(필수 8 + 선택 4).  
- [ ] **Block 3:** 별도 Row이며 **기본 접기(collapsed)** 또는 별도 탭. 메인 뷰 패널 수에 포함하지 않음.  
- [ ] **Refresh:** 1–2분 권장. `histogram_quantile` 사용 패널(P2, T6 등)은 부하 고려.  
- [ ] **Threshold:** Block 1은 0 = OK, &gt;0 = Critical 등 팀 정의. Block 2는 80% 경고, 95% 위험 등 통일.  
- [ ] **Value mappings:** Stat 패널에 0 = OK (green), &gt;0 = Critical (red) 등 적용.  
- [ ] **환경 변수:** 클러스터·네임스페이스 선택이 필요하면 Grafana variable 추가. 최소화 권장.  
- [ ] **PromQL:** 환경에서 실제 metric 이름·라벨 확인 후 endpoint, ingress, scheduler 등 조정.  
- [ ] **고비용 쿼리:** P99·histogram_quantile 사용 패널은 새로고침 주기 늘리거나 Optional로 두어 v1 부하 제한.

---

*Implementation-ready panel specification v1. Grafana 구현 시 이 문서만으로 패널 구성·PromQL·타입·비용·환경 의존성을 적용할 수 있다.*
