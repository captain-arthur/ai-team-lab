# Core Signal List — Cluster Health Monitoring Model v1

**Project:** sre-monitoring-cluster-health-model  
**Phase:** Engineering  
**Input:** Architecture (`03-architecture/architecture.md`)

이 문서는 Cluster Health Monitoring Model을 **실무에서 바로 쓸 수 있는 core signal 목록**으로 옮긴 것이다. Grafana dashboard 설계, alert 정책 논의, operational runbook 작성 시 직접 참고할 수 있도록 구성했다.

---

## 1. 레이어별 신호 수 (목표)

| Layer | 목표 개수 | 답하는 질문 |
|-------|-----------|-------------|
| **Cluster Health Summary** | **5~7개** | 지금 클러스터가 안전한가? |
| **Trend / Risk Indicators** | **5~6개** | 클러스터가 곧 불건강해질 징후는? |
| **Top Offenders / Drill-down** | **5~6개 뷰** (각 뷰당 TOP10) | 문제 시 어디를 먼저 조사할 것인가? |

Summary는 의도적으로 적게 두어, 한 화면에서 “전부 정상이면 안전”이라고 판단할 수 있게 한다.

---

## 2. Cluster Health Summary (5~7개)

운영자가 **“지금 클러스터가 건강한가?”**에 답하는 **최상위 신호만** 포함한다. 이 레이어의 신호가 모두 정상이면 “클러스터 안전”, 하나라도 비정상이면 조사 필요로 해석한다.

### 2.1 Summary 신호 목록

| # | Signal name | Category | Prometheus metric / PromQL example | Interpretation | 비고 |
|---|-------------|----------|------------------------------------|-----------------|------|
| 1 | NotReady node count | Node / Infrastructure | `count(kube_node_status_condition{condition="Ready",status="false"} == 1)` | threshold: 0 = 정상, >0 = 비정상 | |
| 2 | API server health (latency + errors) | Control Plane | P99: `histogram_quantile(0.99, sum(rate(apiserver_request_duration_seconds_bucket[5m])) by (le))`<br>5xx: `sum(rate(apiserver_request_total{code=~"5.."}[5m]))` | threshold: P99 <1s, 5xx 극소 = 정상 | |
| 3 | Scheduler pending pods | Control Plane | `scheduler_pending_pods` (또는 unschedulable 큐) | threshold: 0 = 정상, >0 = 비정상 | |
| 4 | Workload Pending pod count | Workload / Pod | `count(kube_pod_status_phase{phase="Pending"} == 1)` | threshold: 0(또는 허용 범위) = 정상 | 장시간 Pending만 쓰려면 라벨/시간 조건 추가 |
| 5 | Excessive restarts or Eviction | Workload / Pod | 재시작: `sum(increase(kube_pod_container_status_restarts_total[10m])) > N` (N=팀 정의)<br>Eviction: kubelet/eviction 관련 metric 있으면 사용 | threshold: 0 또는 임계치 이하 = 정상 | 일시적 1~2회는 제외하는 임계치 권장 |
| 6 | Critical service endpoint empty | Network / Ingress | `kube_endpoint_address_available{endpoint=~"<핵심서비스>"} == 0` 또는 `kube_endpoint_address_not_ready` 등 | threshold: 비어 있지 않음 = 정상 | “핵심 서비스”는 라벨/목록으로 정의 |

*Network에서 Ingress 5xx율을 Summary에 넣을 경우 7번째로 추가 가능. 환경에 따라 6번과 통합하거나 별도 행으로 둔다.*

### 2.2 Summary 신호 상세

각 Summary 신호에 대해: **why it matters**, **what operational question it answers**, **recommended interpretation**.

---

**S1. NotReady node count**

| 항목 | 내용 |
|------|------|
| **Signal name** | NotReady node count |
| **Category** | Node / Infrastructure |
| **Layer** | Cluster Health Summary |
| **Prometheus / PromQL** | `count(kube_node_status_condition{condition="Ready",status="false"} == 1)` |
| **Why it matters** | NotReady 노드는 스케줄 대상에서 제외된다. 노드 손실은 용량 감소·파드 재스케줄·서비스 영향으로 이어질 수 있다. |
| **Operational question** | 지금 스케줄 가능한 노드가 충분한가? 노드가 빠진 상태는 아닌가? |
| **Interpretation style** | **threshold** — 0 = 정상, >0 = 비정상(점검 필요). |

---

**S2. API server health (latency + errors)**

| 항목 | 내용 |
|------|------|
| **Signal name** | API server health (latency + errors) |
| **Category** | Control Plane |
| **Layer** | Cluster Health Summary |
| **Prometheus / PromQL** | P99: `histogram_quantile(0.99, sum(rate(apiserver_request_duration_seconds_bucket[5m])) by (le))`<br>5xx: `sum(rate(apiserver_request_total{code=~"5.."}[5m]))` |
| **Why it matters** | API server는 모든 제어면·워크로드 동작의 관문이다. 지연·5xx가 나면 kubectl, 스케줄러, controller 등 전반에 영향이 간다. |
| **Operational question** | 제어면이 지금 응답하는가? 에러나 심한 지연이 있는가? |
| **Interpretation style** | **threshold** — P99 <1s, 5xx 없음(또는 극소) = 정상. |

---

**S3. Scheduler pending pods**

| 항목 | 내용 |
|------|------|
| **Signal name** | Scheduler pending pods |
| **Category** | Control Plane |
| **Layer** | Cluster Health Summary |
| **Prometheus / PromQL** | `scheduler_pending_pods` (queue별로 있을 경우 unschedulable 등 확인) |
| **Why it matters** | 스케줄되지 못한 파드가 쌓이면 새 워크로드 기동·스케일이 불가능해진다. 0이어야 “클러스터가 워크로드를 받을 수 있는” 상태다. |
| **Operational question** | 지금 스케줄 대기 중인 파드가 있는가? |
| **Interpretation style** | **threshold** — 0 = 정상, >0 = 비정상. |

---

**S4. Workload Pending pod count**

| 항목 | 내용 |
|------|------|
| **Signal name** | Workload Pending pod count |
| **Category** | Workload / Pod |
| **Layer** | Cluster Health Summary |
| **Prometheus / PromQL** | `count(kube_pod_status_phase{phase="Pending"} == 1)` (필요 시 `kube_pod_status_reason` 등으로 장시간 Pending만 필터) |
| **Why it matters** | Pending이 많으면 스케줄 불가·리소스 부족·정책 문제 등이 있다. 사용자 관점에서 기동이 안 되는 워크로드가 있다는 직접 신호다. |
| **Operational question** | 스케줄 대기 중인 파드가 있는가? |
| **Interpretation style** | **threshold** — 0(또는 허용 범위) = 정상. |

---

**S5. Excessive restarts or Eviction**

| 항목 | 내용 |
|------|------|
| **Signal name** | Excessive restarts or Eviction |
| **Category** | Workload / Pod |
| **Layer** | Cluster Health Summary |
| **Prometheus / PromQL** | 재시작: `sum(increase(kube_pod_container_status_restarts_total[10m]))` — 팀에서 정한 N 초과 시 비정상.<br>Eviction: kubelet eviction metric(환경에 따라 이름 상이) 있으면 사용. |
| **Why it matters** | 10분 내 과도한 재시작·eviction은 워크로드 불안정·용량 부족을 의미한다. 일시적 1~2회(배포·드레인)는 임계치로 걸러 노이즈를 줄인다. |
| **Operational question** | 워크로드가 반복적으로 죽거나, 노드가 리소스 부족으로 파드를 쫓아내는가? |
| **Interpretation style** | **threshold** — 임계치 이하 = 정상. (예: 10분 내 전체 재시작 증가분 < N) |

---

**S6. Critical service endpoint empty**

| 항목 | 내용 |
|------|------|
| **Signal name** | Critical service endpoint empty |
| **Category** | Network / Ingress |
| **Layer** | Cluster Health Summary |
| **Prometheus / PromQL** | `kube_endpoint_address_available` 등으로 핵심 서비스만 필터 후, available == 0 인 서비스가 있는지. 또는 `kube_endpoint_address_not_ready` 등. (쿼리는 환경·라벨에 맞게 조정) |
| **Why it matters** | 엔드포인트가 비어 있으면 해당 서비스는 트래픽을 받지 못한다. “지금 트래픽이 들어갈 수 있는가?”에 대한 직접 답이다. |
| **Operational question** | 핵심 서비스 중 백엔드가 비어 있는 것이 있는가? |
| **Interpretation style** | **threshold** — 비어 있는 핵심 서비스 0개 = 정상. |

---

## 3. Trend / Risk Indicators (5~6개)

**“클러스터가 곧 불건강해질 수 있다”**는 징후를 보여 주는 신호. Summary가 정상이어도 Trend에서 경고가 나올 수 있다.

### 3.1 Trend/Risk 신호 목록

| # | Signal name | Category | Prometheus metric / PromQL example | Interpretation | 비고 |
|---|-------------|----------|------------------------------------|-----------------|------|
| T1 | Node CPU / memory utilization | Node, Capacity | node_exporter: CPU `100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m]))*100)`, 메모리 `(1 - node_memory_MemAvailable_bytes/node_memory_MemTotal_bytes)*100` | threshold: e.g. >80% = 리스크 | |
| T2 | API server inflight requests | Control Plane | `apiserver_current_inflight_requests` | trend: 상승 추세 = 리스크 | |
| T3 | Scheduler pending pods trend | Control Plane, Capacity | `scheduler_pending_pods` 시계열 또는 `increase(scheduler_pending_pods[15m])` | trend: 증가 = 리스크 | |
| T4 | Node disk space | Node / Infrastructure | `(node_filesystem_size_bytes - node_filesystem_avail_bytes) / node_filesystem_size_bytes * 100` (root 등), 또는 `node_filesystem_avail_bytes{mountpoint="/"}` | threshold: 여유 부족(예: <10%) = 리스크 | 선택 |
| T5 | Cluster resource headroom | Capacity | 노드 allocatable 대비 request/usage 비율, 또는 `scheduler_pending_pods`와 연계 | trend 또는 threshold: 여유 거의 없음 = 리스크 | Pending 추세와 중복 가능, 하나로 통합 가능 |

### 3.2 Trend/Risk 신호 상세 (요약)

| Signal | Why it matters | Operational question | Interpretation style |
|--------|----------------|---------------------|----------------------|
| **T1 Node CPU/memory** | 노드가 80% 근접하면 eviction·성능 저하 가능성 증가. | 노드가 포화 직전인가? | **threshold** (e.g. >80% 경고) |
| **T2 API server inflight** | inflight가 계속 높으면 API 지연·타임아웃으로 이어질 수 있음. | API가 포화되는 추세인가? | **trend** (상승 추세) |
| **T3 Pending trend** | Pending이 늘어나면 곧 스케줄 불가·용량 부족. | 대기 파드가 늘어나는가? | **trend** (증가 추세) |
| **T4 Node disk** | 디스크 풀 시 이미지 풀·로그 실패·노드 불안정. | 디스크 여유가 있는가? | **threshold** (여유 부족) |
| **T5 Resource headroom** | 여유가 거의 없으면 새 파드 스케줄 불가·eviction 증가. | 클러스터에 리소스 여유가 있는가? | **trend** 또는 **threshold** |

---

## 4. Top Offenders / Drill-down (5~6개 뷰)

Summary·Trend에서 “이상 있음”이 나왔을 때 **어떤 노드·워크로드가 원인 후보인지** TOP10으로 좁히는 뷰. 각 뷰는 **하나의 ranking 질문**에 대응한다.

### 4.1 Top Offenders 뷰 목록 및 근거

| # | View name | Category | Prometheus / PromQL (개념) | 근거 (왜 이 뷰인가) |
|---|-----------|----------|----------------------------|---------------------|
| O1 | **CPU TOP10 nodes** | Node, Capacity | node_exporter CPU 사용률을 노드별로 집계 후 `topk(10, ...)` | 노드 포화·eviction 원인 파악. Trend T1에서 “노드 사용률 높음”일 때 어느 노드가 가장 높은지 보여 줌. |
| O2 | **Memory TOP10 nodes** | Node, Capacity | node_exporter 메모리 사용률 노드별, `topk(10, ...)` | 메모리 압박·OOM·eviction 원인 파악. |
| O3 | **Restart TOP10 workloads** | Workload / Pod | `sum by (namespace, pod, deployment 등)(increase(kube_pod_container_status_restarts_total[1h]))` 후 `topk(10, ...)` | Summary S5 “과도한 재시작”일 때 어떤 워크로드가 가장 많은 재시작을 일으키는지. |
| O4 | **Pending pods by workload** | Workload / Pod | `count by (namespace, 등)(kube_pod_status_phase{phase="Pending"}==1)` 정렬 | Summary S4 “Pending 있음”일 때 어떤 워크로드가 Pending인지. TOP10 또는 “Pending 있는 워크로드 목록”. |
| O5 | **Error-heavy ingress / services** | Network / Ingress | Ingress controller 또는 서비스별 5xx rate / error rate, `topk(10, ...)` | 트래픽 경로에서 에러가 많은 ingress·서비스를 찾을 때. Summary S6에서 “엔드포인트 비어 있음” 외에 “에러율 높음”을 보려면 이 뷰로 드릴다운. |
| O6 | **Control plane component health** | Control Plane | scheduler, controller-manager 등 `/healthz` 또는 해당 metric별 상태 | Summary S2 “API server 비정상”일 때, 제어면 중 어떤 컴포넌트가 비정상인지. “뷰”는 TOP10이라기보다 “구성요소별 상태 표”에 가깝지만, 드릴다운 용도로 동일 레이어에 둠. |

*O1, O2는 “어느 노드가 부담을 주는가”에, O3, O4는 “어떤 워크로드가 문제인가”에, O5는 “어느 트래픽 경로가 문제인가”에 직접 답한다. O6은 제어면 이상 시 원인 좁히기용이다.*

### 4.2 Top Offenders 신호 상세 (공통 포맷)

| 항목 | 내용 |
|------|------|
| **Interpretation style** | **ranking** — 상위 N개(TOP10)만 표시. |
| **Usage** | Summary 또는 Trend에서 비정상/경고가 나온 뒤, 해당 카테고리와 연결된 Top Offenders 뷰를 열어 “어디를 먼저 조사할 것인가?”에 답한다. |

PromQL 예시(개념):

- **O1 CPU TOP10 nodes:** `topk(10, 100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m]))*100))`
- **O2 Memory TOP10 nodes:** `topk(10, (1 - node_memory_MemAvailable_bytes/node_memory_MemTotal_bytes)*100)` (또는 환경에 맞는 메모리 사용률)
- **O3 Restart TOP10:** `topk(10, sum by(namespace, pod)(increase(kube_pod_container_status_restarts_total[1h])))` (deployment 등 라벨 있으면 by에 포함)
- **O4 Pending by workload:** `count by(namespace)(kube_pod_status_phase{phase="Pending"}==1)` 등으로 정렬
- **O5 Error-heavy ingress:** 환경의 ingress controller metric에 따라 5xx rate by ingress/service, `topk(10, ...)`

---

## 5. 노이즈·보조 신호 (Summary에 넣지 말 것)

다음 신호들은 **Cluster Health Summary에는 포함하지 않는다**. 노이즈가 많거나, “지금 안전한가?”에 대한 직접 답이 아니기 때문이다. Drill-down·별도 패널·낮은 심각도 알림용으로만 사용할 것.

| Signal (개념) | 이유 | 권장 사용처 |
|---------------|------|-------------|
| **Readiness probe 실패** | 일시적 실패·의존성 지연으로 알림만 많을 수 있고, 실제 사용자 영향과 1:1이 아님. | Drill-down, 팀별 조사. Summary 제외. |
| **Calico/CNI 등 인프라 probe** | 플러그인 probe 실패가 알림을 만들지만 트래픽 장애와 무관할 수 있음. | 별도 패널·낮은 심각도. Summary 제외. |
| **일시적 파드 재시작 1~2회** | 배포·드레인으로 1~2회는 metric에 있지만 영향 없을 수 있음. | Summary에서는 “과도한 재시작”만(임계치로 걸러서) 사용. |
| **개별 Control plane healthz** | API server 하나로 Summary에서 제어면 건강을 대표하고, 개별 컴포넌트는 드릴다운(O6)에서만. | Top Offenders(O6) 또는 별도 상세 뷰. |

---

## 6. Dashboard 설계자용 요약표

향후 Grafana central dashboard를 만들 때, 아래 표만으로도 “어떤 패널을 어디에 둘지” 매핑할 수 있다.

| Layer | Signal / View | Category | Interpretation | 패널 배치 제안 |
|-------|----------------|----------|----------------|----------------|
| **Summary** | NotReady node count | Node | threshold | 상단 블록, 1행 |
| **Summary** | API server health | Control Plane | threshold | 상단 블록, 1행 |
| **Summary** | Scheduler pending pods | Control Plane | threshold | 상단 블록, 1행 |
| **Summary** | Workload Pending pod count | Workload | threshold | 상단 블록, 1행 |
| **Summary** | Excessive restarts or Eviction | Workload | threshold | 상단 블록, 1행 |
| **Summary** | Critical service endpoint empty | Network | threshold | 상단 블록, 1행 |
| **Trend** | Node CPU/memory utilization | Node, Capacity | threshold | 중간, Trend 섹션 |
| **Trend** | API server inflight | Control Plane | trend | 중간, Trend 섹션 |
| **Trend** | Pending trend | Control Plane, Capacity | trend | 중간, Trend 섹션 |
| **Trend** | Node disk (선택) | Node | threshold | 중간, Trend 섹션 |
| **Top Offenders** | CPU TOP10 nodes | Node, Capacity | ranking | 하단 또는 별도 탭 |
| **Top Offenders** | Memory TOP10 nodes | Node, Capacity | ranking | 하단 또는 별도 탭 |
| **Top Offenders** | Restart TOP10 workloads | Workload | ranking | 하단 또는 별도 탭 |
| **Top Offenders** | Pending by workload | Workload | ranking | 하단 또는 별도 탭 |
| **Top Offenders** | Error-heavy ingress/services | Network | ranking | 하단 또는 별도 탭 |
| **Top Offenders** | Control plane component health | Control Plane | 목록/상태 | 하단 또는 별도 탭 |

---

## 7. 정리

- **Cluster Health Summary:** 6개(선택 시 7개). threshold만으로 “지금 안전한가?” 판단. 신호 수는 의도적으로 5~7개로 제한.
- **Trend / Risk:** 4~5개. threshold 근접 + trend로 “곧 불건강해질 징후” 파악.
- **Top Offenders:** 6개 뷰. ranking(TOP10)으로 “어디를 먼저 볼 것인가?” 답함. CPU/메모리 노드, 재시작/Pending 워크로드, 에러 많은 ingress, 제어면 구성요소 상태.
- **노이즈·보조:** readiness probe 실패, 인프라 probe, 일시적 재시작 1~2회, 개별 control plane healthz는 Summary에 넣지 않고, drill-down·별도 패널·낮은 심각도로만 사용.

이 core signal list는 **sre-monitoring-dashboard-design**, **sre-monitoring-alert-policy**, **sre-monitoring-operational-runbooks** 프로젝트에서 직접 참조할 수 있다.

---

*Engineering phase output. 다음 단계: Experiment — Prometheus 조회 가능성 검증 및 제한 사항 문서화.*
