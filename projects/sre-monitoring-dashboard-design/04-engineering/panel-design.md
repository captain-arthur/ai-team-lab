# Central Kubernetes Operational Dashboard — Panel Design

**Project:** sre-monitoring-dashboard-design  
**Phase:** Engineering  
**Input:** Architecture (`03-architecture/architecture.md`), Kubernetes Cluster Health Monitoring Model (final-report, core-signal-list)

이 문서는 대시보드 아키텍처를 **Grafana 패널 단위 설계**로 옮긴 것이다. 각 블록(Block 1 Operational Confidence, Block 2 Early Risk, Block 3 Investigation)별 **패널 목록**, **패널 타입**, **metric/PromQL**, **존재 이유**를 정의하여 구현 단계에서 그대로 사용할 수 있게 한다. **메인 뷰**는 Block 1 + Block 2만 포함하며, **총 10~14개 이하**로 유지한다. Block 3(Investigation)은 **드릴다운 전용**으로, 메인 뷰를 과부하시키지 않는다.

---

## 1. 메인 뷰와 드릴다운 구분

| 구분 | 포함 블록 | 패널 수 | 용도 |
|------|-----------|---------|------|
| **메인 뷰 (Primary view)** | Block 1 + Block 2 | **10~14개 이하** | 5–10분 일상 점검. “지금 안전한가?” / “곧 위험한가?” 판단. |
| **드릴다운 (Drill-down)** | Block 3 | 4~6개 뷰 | 이상 시 “어디를 조사할 것인가?” — TOP10 등. 탭·접기·별도 섹션. |

---

## 2. Block 1: Operational Confidence (운영 확신)

**목적:** “지금 클러스터가 안전하다고 확신할 수 있는가?”에 답한다. 전부 정상이면 “안전”, 하나라도 비정상이면 “조사 필요”.

**배치:** 메인 뷰 최상단. 한 행(또는 두 행)에 배치해 스크롤 없이 볼 수 있게 한다.

**패널 수:** 4~6개(환경별). Managed K8s 등에서 control plane metric이 없으면 P2, P3 제외 → 4개.

### 2.1 Block 1 패널 목록

| Panel name | Category | Layer | Metric / PromQL example | Panel type | Why the panel exists |
|------------|----------|--------|--------------------------|------------|----------------------|
| **NotReady node count** | Node / Infrastructure | Confidence | `count(kube_node_status_condition{condition="Ready",status="false"} == 1)` | **stat** | NotReady 노드가 0이어야 스케줄 용량이 유지됨. 0 = 정상, >0 = 비정상. 한 눈에 상태 판단. |
| **API server health** | Control Plane | Confidence | P99: `histogram_quantile(0.99, sum(rate(apiserver_request_duration_seconds_bucket[5m])) by (le))`<br>5xx: `sum(rate(apiserver_request_total{code=~"5.."}[5m]))` | **stat** (또는 2개 stat: P99 + 5xx) | 제어면 응답·에러 여부. P99 <1s, 5xx 극소 = 정상. **Control plane 미사용 환경에서는 패널 제외.** |
| **Scheduler pending pods** | Control Plane | Confidence | `scheduler_pending_pods` | **stat** | 스케줄 대기 파드 0이어야 클러스터가 워크로드를 받을 수 있음. **Control plane 미사용 환경에서는 패널 제외.** |
| **Workload Pending pod count** | Workload / Pod | Confidence | `count(kube_pod_status_phase{phase="Pending"} == 1)` | **stat** | 사용자 관점 “기동 안 되는 파드” 수. 0 = 정상. |
| **Excessive restarts** | Workload / Pod | Confidence | `sum(increase(kube_pod_container_status_restarts_total[10m]))` (임계치 N과 비교해 threshold 설정) | **stat** | 10분 내 과도한 재시작은 워크로드 불안정·용량 이슈. 임계치 이하 = 정상. |
| **Critical service endpoint empty** | Network / Ingress | Confidence | `kube_endpoint_address_available` 등으로 핵심 서비스 필터 후, 비어 있는 엔드포인트 수 또는 0/1 | **stat** | 핵심 서비스 백엔드가 비어 있으면 트래픽 불가. 비어 있지 않음 = 정상. |

**Grafana 표시 권장:**  
- 각 패널에 **threshold** 설정: 정상 구간(녹색), 경고/비정상(빨강 또는 노랑).  
- **Value mappings** 로 “0 = OK”, “>0 = Critical” 등 한눈에 상태가 보이게 한다.  
- **Panel type:** 현재 값과 상태만 보면 되므로 **stat**이 적합. 숫자보다 “정상/비정상”이 먼저 보이도록 옵션(색, 아이콘) 사용.

---

## 3. Block 2: Early Risk (조기 리스크)

**목적:** “클러스터가 곧 불안전해질 조기 징후가 보이는가?”에 답한다. 노드 압박, OOM 위험, 디스크, pending 추세, (가능 시) CPU throttling, ingress stress를 강조.

**배치:** 메인 뷰에서 Block 1 직하단. 스크롤 한 번 이내.

**패널 수:** 4~6개. 메인 뷰 합계(Block 1 + Block 2)가 10~14개 이하가 되도록 Block 2는 4~6개로 제한.

### 3.1 Block 2 패널 목록

| Panel name | Category | Layer | Metric / PromQL example | Panel type | Why the panel exists |
|------------|----------|--------|--------------------------|------------|----------------------|
| **Node CPU utilization** | Node, Capacity | Risk | `100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m]))*100)` → 클러스터 전체 평균 또는 “80% 초과 노드 수” | **gauge** (또는 **stat**) | 노드 CPU 포화 직전(예: 80% 이상)이면 eviction·성능 저하 리스크. 조기 리스크 가시화. |
| **Node memory / OOM risk** | Node, Capacity | Risk | `(1 - node_memory_MemAvailable_bytes/node_memory_MemTotal_bytes)*100` → 평균 또는 “80% 초과 노드 수” | **gauge** (또는 **stat**) | 메모리 압박 = OOM·eviction 위험. “OOM risk” 또는 “메모리 압박” 라벨로 강조. |
| **Node disk space** | Node / Infrastructure | Risk | `(node_filesystem_size_bytes - node_filesystem_avail_bytes) / node_filesystem_size_bytes * 100` (mountpoint="/" 등) | **gauge** 또는 **stat** | 디스크 여유 부족 시 이미지 풀·로그 실패·노드 불안정. 여유 <10% 등 threshold로 경고. |
| **Pending pods trend** | Workload, Capacity | Risk | `count(kube_pod_status_phase{phase="Pending"} == 1)` 시계열 또는 `increase(...[15m])` | **time-series** (또는 **stat** “현재 Pending 수”) | Pending이 늘어나면 스케줄·용량 부족 징후. 추세 또는 현재값으로 “조기 리스크” 표시. |
| **CPU throttling risk** (선택) | Workload | Risk | 컨테이너 CPU throttling metric(환경에 따라 상이) | **stat** 또는 **gauge** | 파드가 CPU limit에 막혀 있으면 성능 저하·조기 리스크. **metric 없으면 패널 제외 또는 Node CPU로 대체.** |
| **Ingress stress** (선택) | Network / Ingress | Risk | Ingress controller 지연·에러율·부하(제품별 metric) | **stat** 또는 **gauge** | Ingress 부하·에러가 높으면 트래픽 경로 리스크. **metric 있을 때만 추가.** |

**Grafana 표시 권장:**  
- **gauge** 사용 시 **thresholds:** 예: 0–80 녹색, 80–95 노랑, 95–100 빨강.  
- **Block 2 제목:** “Early Risk” 또는 “조기 리스크”로 블록 라벨을 두어, Summary가 정상이어도 이 블록에서 “곧 나빠질 수 있음”을 인지하게 함.  
- **메인 뷰 패널 수:** Block 1이 4~6개, Block 2가 4~6개이면 **합계 8~12개.** 14개 초과하지 않도록 Block 2에서 선택 패널(CPU throttling, Ingress stress)은 환경에 따라 생략 가능.

---

## 4. Block 3: Investigation / Top Offenders (조사 / 드릴다운)

**목적:** “문제가 있다면 어디를 먼저 조사할 것인가?”에 답한다. TOP10 스타일로 노드·워크로드·ingress 등을 나열.

**배치:** **메인 뷰에 포함하지 않음.** 탭·접기·아래쪽 별도 섹션으로 “이상 시” 또는 “더 보기”로 진입.  
**메인 뷰 패널 수에 포함하지 않음.**

### 4.1 Block 3 패널 목록

| Panel name | Category | Layer | Metric / PromQL example | Panel type | Why the panel exists |
|------------|----------|--------|--------------------------|------------|----------------------|
| **CPU TOP10 nodes** | Node, Capacity | Investigation | `topk(10, 100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m]))*100))` | **table** | 노드 포화·eviction 원인 파악. Block 2 “Node CPU 높음”일 때 어느 노드가 1위인지 드릴다운. |
| **Memory TOP10 nodes** | Node, Capacity | Investigation | `topk(10, (1 - node_memory_MemAvailable_bytes/node_memory_MemTotal_bytes)*100)` | **table** | 메모리 압박·OOM 원인 노드. Block 2 “OOM risk”일 때 드릴다운. |
| **Restart TOP10 pods** | Workload / Pod | Investigation | `topk(10, sum by(namespace, pod)(increase(kube_pod_container_status_restarts_total[1h])))` (또는 deployment 등 by 추가) | **table** | Block 1 “Excessive restarts”일 때 어떤 워크로드가 재시작 1위인지. |
| **Pending pods by workload** | Workload / Pod | Investigation | `count by(namespace)(kube_pod_status_phase{phase="Pending"}==1)` 정렬 (또는 owner/라벨별) | **table** | Block 1 “Pending 있음”일 때 어떤 네임스페이스/워크로드가 Pending인지. |
| **Error-heavy ingress / services** (선택) | Network / Ingress | Investigation | Ingress controller 5xx rate by ingress/service, `topk(10, ...)` | **table** | 에러율이 높은 ingress·서비스 TOP10. **metric 있을 때만.** |
| **Control plane component health** (선택) | Control Plane | Investigation | scheduler, controller-manager 등 `/healthz` 또는 해당 metric별 상태 | **table** | API server 비정상일 때 제어면 중 어떤 컴포넌트가 비정상인지. **Self-managed 등에서만.** |

**Grafana 표시 권장:**  
- **table** 로 컬럼: 노드명(또는 instance), 값(CPU% / 메모리% / 재시작 수 / Pending 수 등).  
- **정렬:** 값 내림차순. TOP10만 표시.  
- **드릴다운 전용:** 이 블록 전체를 **Grafana 탭** 또는 **접을 수 있는 row**로 두어, 기본 뷰에는 Block 1·2만 보이게 한다.  
- **메인 뷰 과부하 방지:** Block 3 패널은 **메인 뷰 패널 수(10~14)에 포함하지 않는다.**

---

## 5. 패널 수 요약

| Block | 패널 수 | 메인 뷰 포함 | 비고 |
|-------|---------|--------------|------|
| **Block 1 (Operational Confidence)** | 4~6개 | **예** | S2, S3는 control plane 미사용 시 제외 → 최소 4개 |
| **Block 2 (Early Risk)** | 4~6개 | **예** | CPU throttling, Ingress stress는 선택. 합계 10~14 이하 유지 |
| **Block 3 (Investigation)** | 4~6개 | **아니오** | 드릴다운 전용. TOP10 table. |

**메인 뷰 총 패널 수:** Block 1 + Block 2 = **8~12개**(표준), 최대 **14개 이하** 권장.  
**Block 3:** 4~6개 패널이지만 **메인 뷰에 항상 노출하지 않음.**

---

## 6. Grafana 레이아웃 권장

### 6.1 행 구성

```
Row 1: [제목] Operational Confidence
  → Panel 1 (NotReady) | Panel 2 (API health)* | Panel 3 (Scheduler pending)* | Panel 4 (Pending pods) | Panel 5 (Restarts) | Panel 6 (Endpoint)
  * 환경에 따라 숨김

Row 2: [제목] Early Risk
  → Panel 7 (Node CPU) | Panel 8 (Node memory/OOM) | Panel 9 (Node disk) | Panel 10 (Pending trend) | [Panel 11 CPU throttling]* | [Panel 12 Ingress stress]*
  * 선택

Row 3 (또는 별도 탭): [제목] Investigation / Top Offenders
  → Panel 13 (CPU TOP10) | Panel 14 (Memory TOP10) | Panel 15 (Restart TOP10) | Panel 16 (Pending by workload) | [Panel 17 Error-heavy ingress]* | [Panel 18 Control plane]*
  * 선택. 이 Row는 접기 또는 탭으로 기본 숨김
```

### 6.2 구현 시 참고

- **Grafana variables:** 클러스터·네임스페이스 선택이 필요하면 변수 추가. 메인 뷰 단순화를 위해 변수 수는 최소화.
- **Refresh:** 메인 뷰(Block 1·2)는 1–2분 간격 새로고침 권장. Prometheus 부하 고려.
- **Thresholds:** Block 1은 “0 = 정상, >0 = 비정상” 등 팀 정의 threshold. Block 2는 “80% 초과 = 경고” 등 조기 리스크 기준을 동일하게 적용.

---

## 7. 환경별 변형

| 환경 | Block 1 | Block 2 | Block 3 |
|------|----------|----------|---------|
| **Full (control plane 스크래핑 가능)** | 6개 (P1–P6) | 4~6개 (Node CPU/메모리/디스크, Pending trend, 선택 2) | 4~6개 |
| **Managed K8s (control plane 미사용)** | 4개 (P1, P4, P5, P6. P2·P3 제외) | 4개 (T1 CPU/메모리, T4 디스크, Pending trend. T2·T3 대체 또는 제외) | 4개 (O6 제외) |

---

## 8. 요약

- **Block 1 (Operational Confidence):** 4~6개 **stat** 패널. NotReady, API health, Scheduler pending, Workload Pending, Excessive restarts, Critical endpoint. 메인 뷰 최상단.
- **Block 2 (Early Risk):** 4~6개 **gauge/stat/time-series** 패널. Node CPU, Node memory/OOM risk, Node disk, Pending trend, (선택) CPU throttling, Ingress stress. 메인 뷰 Block 1 직하.
- **Block 3 (Investigation):** 4~6개 **table** 패널. CPU TOP10 nodes, Memory TOP10 nodes, Restart TOP10 pods, Pending by workload, (선택) Error-heavy ingress, Control plane component. **드릴다운 전용**, 메인 뷰에 포함하지 않음.
- **메인 뷰:** Block 1 + Block 2만. **총 10~14개 이하** 패널. 5–10분 일상 점검·운영 확신·조기 리스크 가시성·과부하 방지를 만족하도록 유지.

이 문서를 기준으로 Grafana에서 패널을 생성하고, PromQL·threshold·패널 타입을 연결하면 된다. Review 단계에서 “5–10분 파악”, “운영 확신”, “조기 리스크”, “과부하 방지” 관점으로 검증할 수 있다.

---

*Engineering phase output. 다음 단계: Review — 설계 검증.*
