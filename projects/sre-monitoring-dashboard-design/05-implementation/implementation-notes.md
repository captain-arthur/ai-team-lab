# Implementation Notes — Central Kubernetes Operational Dashboard v1

**Project:** sre-monitoring-dashboard-design  
**Phase:** 05-implementation  
**Input:** operational-confidence-theory, panel-design, implementation-ready-panel-spec, central-kubernetes-operational-dashboard-design

이 문서는 **Grafana + Prometheus** 로 v1 대시보드를 실제 배포할 때 고려해야 할 **실무적 사항**을 정리한다. 쿼리 비용, 새로고침 주기, control plane/ingress metric이 없는 환경, 그리고 가능한 최적화를 중심으로 서술한다.

---

## 1. Prometheus 쿼리 비용 (Query cost)

### 1.1 비용이 낮은 쿼리 (v1 전부 해당)

- **Instant query + count/sum:** NotReady node count, Workload Pending, Critical endpoint empty, Node memory, Node disk, Pending by workload — 모두 **instant** 이며 시리즈 수가 수백~수천 수준이면 부하는 작다.  
- **rate(...[5m]) + avg/topk(10):** Node CPU, CPU TOP10, Memory TOP10 — **5분 range** 한 번이므로 상대적으로 가볍다.  
- **increase(...[10m]) 한 번:** Excessive restarts — 10분 구간이면 retention 범위 내에서 무리 없음.  

**권장:** v1에서 **high 비용** 쿼리(긴 range, histogram_quantile 다수, 대량 시리즈 스캔)는 **도입하지 않았다.** Optional(B) 패널인 API server P99, Ingress histogram 등은 별도 도입 시 비용을 고려한다.

### 1.2 중간 비용 — Restart TOP10 (1h increase)

- **Restart TOP10 pods** 는 `increase(...[1h])` 를 사용한다. 파드·컨테이너 수가 많으면 시리즈 수가 커질 수 있어 **medium** 으로 분류했다.  
- **대응:** Prometheus retention이 1h 미만이면 구간을 30m 등으로 줄이거나, **재시작 수** 대신 **현재 재시작 카운터 값**으로 대체하는 변형을 검토한다.  
- 대시보드 새로고침을 **2m** 로 두면 1h increase 쿼리는 2분마다 한 번만 실행되므로, 일반적인 클러스터 규모에서는 수용 가능하다.

### 1.3 쿼리 수와 동시 실행

- 메인 뷰(Block 1 + Block 2) **8개 패널** + Block 3 **4개 패널** = **12개 쿼리**가 새로고침 시마다 실행된다.  
- Block 3이 **접혀 있을 때**에도 Grafana는 해당 row 내 패널 쿼리를 실행할 수 있으므로, 필요하면 **변수로 “Block 3 쿼리 실행” 토글**을 두거나, Block 3 row를 완전히 비활성화하는 방식은 Grafana 기본 기능만으로는 제한적이다. 실무에서는 **2m refresh** 로 12개 쿼리 동시 실행을 가정하고, Prometheus 스크래핑 부하와 함께 관찰한다.

---

## 2. 대시보드 새로고침 주기 (Refresh interval)

- **권장:** **1m** 또는 **2m**.  
- **이유:**  
  - Block 1·2는 “지금 안전한가?” / “조기 징후인가?”에 답하는 데 **1~2분 지연**이면 충분하다.  
  - 더 짧은 주기(30s 등)는 Prometheus에 불필요한 부하를 주고, **rate(5m)** 구간과도 맞지 않는다.  
- **사용자 선택:** 대시보드 상단에서 refresh를 “Off” / “1m” / “2m” / “5m” 중 선택 가능하게 두면, 부하가 큰 환경에서는 5m으로 늘릴 수 있다.

---

## 3. Control plane metric이 없는 환경 (Managed K8s 등)

- **EKS, GKE, AKS** 등에서는 사용자가 **API server, scheduler** metric을 스크래핑하지 못하는 경우가 많다.  
- **대응:**  
  - **Block 1** 에는 A 우선순위 **4개만** 사용한다: NotReady node count, Workload Pending pod count, Excessive restarts, Critical service endpoint empty.  
  - API server health, Scheduler pending pods 패널은 **이번 v1 JSON에는 포함하지 않았으며**, B 우선순위이므로 **self-managed 등 control plane metric이 있는 환경**에서만 별도 패널로 추가한다.  
  - “제어면이 정상인가?”는 **managed 서비스 콘솔·알림·지원 채널**로 확인하도록 runbook·대시보드 설명에 명시한다.  
- **grafana-dashboard-v1.json** 은 A 패널만 포함하므로, Managed K8s에서도 **그대로** 임포트해 사용할 수 있다.

---

## 4. Ingress metric이 없는 환경

- **Ingress stress**, **Error-heavy ingress** 는 B 우선순위이며, **nginx-ingress, Istio** 등 Ingress controller별 metric에 의존한다.  
- **v1 JSON에는 A 패널만** 넣었으므로 Ingress 패널은 **포함되지 않았다.**  
- Ingress metric을 나중에 도입할 때는:  
  - 해당 controller의 **request duration, 5xx rate** 등 metric 이름을 환경에서 확인한 뒤, promql-spec과 panel-config에 **환경별 절**을 추가하고,  
  - **histogram_quantile** 사용 시 쿼리 비용이 medium이므로 **refresh 2m** 유지 또는 해당 패널만 더 긴 주기로 두는 것을 권장한다.

---

## 5. Endpoint metric 이름 차이 (Critical service endpoint empty) — 검증 반영

- **kube-state-metrics** 버전에 따라 **endpoint** 관련 metric 이름이 다르다.  
  - **단수:** `kube_endpoint_address_available`  
  - **복수 (첫 검증에서 동작):** `kube_endpoints_address_available`  
- **grafana-dashboard-v1.json** 기본값은 **복수형** `kube_endpoints_address_available` 로 되어 있다(validation-results-first-pass 반영). No data 또는 Query error 시 **단수형**으로 패널 쿼리를 바꾼다.  
- **대응:** `curl http://kube-state-metrics:8080/metrics | grep -E "endpoint.*available"` 로 실제 metric 이름 확인 후 교체. critical 서비스만 보려면 `namespace=~"default|production"` 등 필터 추가.

---

## 6. 가능한 최적화 (Optimizations)

### 6.1 Recording rules (선택)

- **Node CPU / memory / disk** 의 클러스터 평균이나 **80% 초과 노드 수** 등을 **주기적으로 미리 계산**해 두면, 대시보드 쿼리가 **instant** 한 번으로 끝나도록 할 수 있다.  
- 예: `avg(100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m]))*100))` 를 1m 주기 recording rule로 두고, 대시보드는 해당 rule metric만 조회.  
- v1에서는 **필수는 아니며**, Prometheus 부하가 문제가 될 때 검토한다.

### 6.2 Block 3 쿼리 지연 로딩

- Block 3(Investigation)은 **이상 시**에만 보는 드릴다운이다.  
- Grafana 기본만으로는 “row가 펼쳐질 때만 쿼리 실행”을 구현하기 어렵지만, **변수**로 “Investigation 보기” on/off를 두고, off일 때 Block 3 패널의 쿼리에서 **항상 0을 반환하는 조건**을 넣어 “실제 쿼리 스킵”을 흉내 내는 방법은 있다.  
- 실무에서는 **2m refresh + 12개 쿼리**로 먼저 운영해 보며, 부하가 크면 recording rule 또는 패널 수 축소를 검토한다.

### 6.3 Excessive restarts threshold (N)

- **promql-spec** 과 **panel-config** 에서 “N”은 **팀 정의**로 두었다.  
- 초기값 예: **10** (10분 내 재시작 10회 초과 시 비정상).  
- 클러스터 규모·워크로드 특성에 따라 **5~20** 등으로 조정하고, 알림 정책과 맞춘다.

### 6.4 디스크 metric — mountpoint / fstype (검증 반영)

- **node_exporter** 에서 `mountpoint="/"` 만 쓰면 여러 fstype(tmpfs, overlay 등)이 섞일 수 있다. **promql-spec** 에서는 `fstype!~"tmpfs|overlay"` 로 필터.  
- **환경별 조정:** 실제 노드의 **root 디스크** mountpoint가 `/` 가 아니면(예: `/var` 등) **Node disk space** 패널 쿼리에서 `mountpoint="/"` 를 `mountpoint="/var"` 등으로 변경. 첫 검증에서 “조건부” 로 기록됨(validation-results-first-pass).

---

## 7. 검증 결과 반영 (Validation findings — first pass)

첫 검증 패스(validation-results-first-pass.md) 결과를 반영한 **환경별 소규모 수정** 사항이다.

### 7.1 P4 — Endpoint metric 이름 차이

- **현상:** Critical service endpoint empty 패널에서 `kube_endpoint_address_available` (단수) 미존재 시 Query error / No data.
- **적용:** grafana-dashboard-v1.json 및 promql-spec에서 **복수형** `kube_endpoints_address_available` 을 기본(또는 대안)으로 문서화·반영. 환경에 따라 단수형으로 되돌릴 수 있음.
- **문서:** promql-spec §2.4, 위 §5.

### 7.2 T3 — Node disk space mountpoint

- **현상:** root 디스크가 `mountpoint="/"` 가 아닌 환경(예: `/var`)에서는 0 또는 비정상 값.
- **적용:** promql-spec §3.3에 **mountpoint 환경별 조정** 명시. 해당 환경에서는 패널 쿼리에서 `mountpoint="/var"` 등으로 변경.
- **문서:** implementation-notes §6.4.

### 7.3 Excessive restarts — Threshold N 팀 정의

- **현상:** JSON 기본 threshold 10이 모든 클러스터에 맞지 않을 수 있음.
- **적용:** **Threshold N은 반드시 팀에서 정의** 후 Grafana에 반영(예: N=10~20). promql-spec §2.3 Aggregation strategy에 명시.
- **문서:** panel-config, implementation-notes §6.3.

---

## 8. 배포 체크리스트 (요약)

- [ ] Prometheus에 **kube-state-metrics**, **node_exporter** 타겟이 등록되어 있는지 확인.  
- [ ] **Endpoint** metric 이름을 환경에서 확인 후 Critical service endpoint empty 쿼리 반영.  
- [ ] **Excessive restarts** threshold N 팀 정의 후 Grafana threshold 설정.  
- [ ] 대시보드 **refresh** 1m 또는 2m 설정.  
- [ ] **datasource** 변수가 실제 Prometheus datasource UID와 연결되는지 확인.  
- [ ] Managed K8s인 경우 Block 1이 4개 패널만 있으면 됨(현재 JSON과 동일). 제어면 건강은 별도 채널 안내.  
- [ ] Block 3 row가 **collapsed by default** 로 열리면 “Investigation”이 필요할 때만 펼쳐서 사용.  
- [ ] **첫 검증 반영:** P4 endpoint metric 이름(복수형 기본), T3 mountpoint(환경별), Excessive restarts N(팀 정의) 확인.

---

*Implementation notes v1. 검증: validation-results-first-pass.md.* 운영 시 발생하는 이슈는 이 문서를 기반으로 runbook·alert 정책과 연계해 보완한다.*
