# Validation Plan — Central Kubernetes Operational Dashboard v1

**Project:** sre-monitoring-dashboard-design  
**Phase:** 05-implementation (validation)  
**Input:** grafana-dashboard-v1.json, promql-spec.md, panel-config.md, implementation-notes.md

이 문서는 **v1 Grafana 대시보드 구현 후보(grafana-dashboard-v1.json)** 를 **실제 Prometheus/Grafana 환경**에서 임포트·테스트하여 검증하는 **절차**를 정의한다. 검증 결과는 **validation-results-template.md** 에 기록한다.

---

## 1. 검증 대상

- **대상:** `05-implementation/grafana-dashboard-v1.json` (v1 implementation candidate).
- **목표:**  
  - 대시보드가 Grafana에 정상 임포트·표시되는지 확인.  
  - 각 A 우선순위 패널의 PromQL이 **실제 환경**에서 동작하는지, 필요 metric·라벨이 존재하는지, 출력이 **운영적으로 의미 있는지**, threshold가 적절한지 검증.  
- **산출물:** 검증 체크리스트 수행 후 **validation-results-template.md** 에 결과 기록. “Production usable” 여부 판정.

---

## 2. 임포트 절차 (Import steps)

### 2.1 사전 조건

1. **Grafana** 인스턴스에 접근 가능. (온프레미스 또는 호스팅)
2. **Prometheus** 가 Kubernetes 클러스터(또는 해당 환경) 메트릭을 스크래핑 중.
3. **kube-state-metrics**, **node_exporter** 가 Prometheus 타겟으로 등록되어 있음. (implementation-notes 참고)

### 2.2 임포트 단계

| 단계 | 작업 | 확인 사항 |
|------|------|-----------|
| 1 | Grafana 로그인 → **Dashboards** → **New** → **Import**. | — |
| 2 | **Upload JSON file** 선택 후 `grafana-dashboard-v1.json` 업로드. 또는 **Import via panel json** 에 JSON 내용 붙여넣기. | 파일이 유효한 JSON인지(이미 검증됨). |
| 3 | **Options** 에서 **UID** 충돌 시 기존 대시보드 덮어쓸지, 새 UID로 생성할지 선택. | `central-kubernetes-operational-v1` 이 이미 있으면 덮어쓰기 또는 이름 변경. |
| 4 | **Prometheus** datasource를 선택. (아래 3절 참고) | 변수 `datasource` 가 실제 Prometheus datasource를 가리키는지. |
| 5 | **Import** 실행. | 대시보드가 열리고 Row 3(Investigation)이 **접힌 상태**로 보이는지. |
| 6 | 상단 **Refresh** 를 **2m** (또는 1m)으로 설정. | 새로고침 시 패널이 갱신되는지. |

### 2.3 임포트 후 점검

- [ ] **Row 1 (Operational Confidence)** 에 패널 4개가 한 행에 보인다.  
- [ ] **Row 2 (Early Risk)** 에 패널 4개가 한 행에 보인다.  
- [ ] **Row 3 (Investigation)** 이 **collapsed** 로 되어 있고, 펼치면 패널 4개가 보인다.  
- [ ] 모든 패널에 **데이터소스**가 연결되어 있고, “No data” 또는 “Query error” 가 **datasource/메트릭 부재** 때문인지 구분 가능하다.

---

## 3. Datasource 설정 가정 (Datasource setup assumptions)

- **타입:** Prometheus.  
- **Grafana datasource UID:** 환경마다 다름. 대시보드 JSON의 `templating.list[].query` 가 `prometheus` 이므로, **이름이 "Prometheus"인 Prometheus datasource**가 있으면 변수 `datasource` 가 자동으로 해당 datasource를 선택한다.  
- **검증 시 확인:**  
  - [ ] Grafana에 Prometheus datasource가 등록되어 있고 **Working** 상태인지.  
  - [ ] 해당 datasource가 **kube-state-metrics**, **node_exporter** 를 스크래핑하는 Prometheus 서버를 가리키는지.  
  - [ ] 대시보드 상단 변수에서 **Datasource** 가 올바른 Prometheus로 선택되어 있는지.  
- **Multi-cluster:** 여러 클러스터를 하나의 Prometheus로 스크래핑하는 경우, 필요 시 `cluster` 라벨 등으로 필터하는 변수를 추가할 수 있으나, v1 검증 범위에서는 **단일 Prometheus** 가정.

---

## 4. 패널별 검증 체크리스트 (Panel-by-panel validation)

각 A 우선순위 패널에 대해 아래 **5가지**를 검증한다. 결과는 **validation-results-template.md** 에 기록.

1. **PromQL이 실제 환경에서 동작하는가?** — 쿼리 실행 시 에러 없이 결과(숫자/시리즈/테이블)가 반환되는지.  
2. **필요 metric이 존재하는가?** — Prometheus에 해당 metric이 수집되고 있는지. (Prometheus UI 또는 Grafana Explore에서 확인)  
3. **metric 이름·라벨 조정이 필요한가?** — kube-state-metrics/node_exporter 버전 차이로 metric 이름 또는 라벨이 다르면, 사용한 PromQL·필터를 기록.  
4. **패널 출력이 운영적으로 의미 있는가?** — 값이 “지금 안전한가?” / “조기 징후인가?” / “조사 대상” 판단에 쓸 수 있는 수준인지. (예: NotReady가 0이면 정상으로 해석 가능한지)  
5. **Threshold가 적절한가?** — 팀 정의 또는 권장(0=OK, 80%=경고 등)과 맞는지, 조정이 필요하면 기록.

### 4.1 Block 1 — Operational Confidence

| # | Panel name | 검증 포인트 (요약) |
|---|-------------|---------------------|
| 1 | **NotReady node count** | `kube_node_status_condition` 존재 여부. `condition="Ready",status="false"` 라벨. 결과 0 = 정상. threshold 0.01 (>&gt;0 빨강). |
| 2 | **Workload Pending pod count** | `kube_pod_status_phase{phase="Pending"}` 존재. 결과 0 = 정상. threshold 0.01. |
| 3 | **Excessive restarts** | `kube_pod_container_status_restarts_total` 존재. `increase(...[10m])` 가 retention 내에서 동작하는지. threshold N(예: 10) 팀 정의 반영 여부. |
| 4 | **Critical service endpoint empty** | **Endpoint metric 이름** 환경 확인. `kube_endpoint_address_available` vs `kube_endpoints_address_available` 등. available == 0 개수. critical namespace 필터 필요 시 적용. |

### 4.2 Block 2 — Early Risk

| # | Panel name | 검증 포인트 (요약) |
|---|-------------|---------------------|
| 5 | **Node CPU utilization** | `node_cpu_seconds_total{mode="idle"}`. `instance` 라벨. rate(5m). gauge 0–80–95–100 threshold. |
| 6 | **Node memory / OOM risk** | `node_memory_MemAvailable_bytes`, `node_memory_MemTotal_bytes`. 사용률 0–80–95–100. |
| 7 | **Node disk space** | `node_filesystem_avail_bytes`, `node_filesystem_size_bytes`. `mountpoint="/"`, `fstype` 필터(필요 시). 0–90–95–100. |
| 8 | **Pending pods trend** | `kube_pod_status_phase{phase="Pending"}`. time series로 추세가 보이는지, 또는 stat으로 현재값. |

### 4.3 Block 3 — Investigation

| # | Panel name | 검증 포인트 (요약) |
|---|-------------|---------------------|
| 9 | **CPU TOP10 nodes** | `topk(10, ...)` 결과가 테이블에 instance + value(CPU %)로 나오는지. 정렬 내림차순. |
| 10 | **Memory TOP10 nodes** | 동일. 메모리 % 컬럼. |
| 11 | **Restart TOP10 pods** | `increase(...[1h])`. retention ≥1h. namespace, pod 컬럼. |
| 12 | **Pending pods by workload** | `count by(namespace)(...)`. 테이블에 namespace + count. |

---

## 5. 자주 발생하는 실패와 대응 (Common failure cases)

| 현상 | 원인 | 대응 |
|------|------|------|
| **No data** (전체 또는 일부 패널) | Datasource 잘못 선택, Prometheus URL/네트워크 문제, 해당 metric 미수집 | Datasource 연결 상태 확인. Prometheus targets에서 kube-state-metrics, node_exporter up 여부 확인. |
| **Query error / invalid expression** | PromQL 문법 오류, **metric 이름 오타** 또는 해당 환경에 없는 metric | Grafana Explore에서 동일 쿼리 실행해 에러 메시지 확인. promql-spec의 metric 이름과 Prometheus에 실제 있는 metric 이름 비교. |
| **Critical service endpoint empty 쿼리 실패** | `kube_endpoint_address_available` 가 없고 `kube_endpoints_address_available` 등 다른 이름 사용 | implementation-notes §5. `curl .../metrics \| grep endpoint` 로 실제 이름 확인 후 쿼리·JSON 수정. |
| **Node disk space 0 또는 비정상** | `mountpoint="/"` 인 시리즈가 없음(다른 경로 사용), 또는 `fstype` 필터로 전부 제외됨 | `node_filesystem_*` 를 라벨 없이 조회해 mountpoint/fstype 확인. promql-spec의 필터를 환경에 맞게 수정. |
| **Restart TOP10 “No data” 또는 빈 테이블** | Prometheus retention &lt; 1h 이라 increase(1h) 구간이 없음 | retention 확인. 구간을 30m 등으로 줄이거나, increase 대신 `sum by(...)(kube_pod_container_status_restarts_total)` 등 현재값 기반 변형 검토. |
| **Excessive restarts 값이 항상 0** | increase 구간 내 데이터 없음, 또는 metric 이름/라벨 불일치 | 10m 구간이 retention 내인지 확인. metric이 정상 수집되는지 Explore에서 확인. |
| **Gauge/Stat 단위가 “percent” 인데 값이 0–1 범위로 나옴** | PromQL은 이미 *100 했는데 Grafana unit이 또 percent로 해석됨 | Panel 옵션에서 Unit을 **none** 또는 **short** 로 바꾸거나, PromQL에서 *100 제거 후 unit percent 유지. (일관되게 하나만 적용) |
| **Row 3 패널이 접혀 있어도 쿼리 실행됨** | Grafana 동작. collapsed row 내부 패널도 refresh 시 쿼리됨 | 구현 참고 사항. 부하가 크면 recording rule 또는 패널 비활성화 검토(implementation-notes). |

---

## 6. Metric / 라벨 불일치 점검 (Metric/label mismatch checks)

검증 시 아래를 **환경에서 실제로 확인**하고, spec과 다르면 **validation-results-template** 에 “실제 사용한 metric/라벨” 로 기록한다.

### 6.1 kube-state-metrics

| 검증 항목 | 확인 방법 | spec 기준 |
|-----------|-----------|-----------|
| 노드 조건 | `kube_node_status_condition` 존재, 라벨 `condition`, `status` | condition="Ready", status="false" |
| 파드 phase | `kube_pod_status_phase` 존재, 라벨 `phase` | phase="Pending" |
| 재시작 카운터 | `kube_pod_container_status_restarts_total` 존재 | — |
| Endpoint | endpoint 관련 metric 이름 | `kube_endpoint_address_available` 또는 환경별 이름 |

**확인 예:** Prometheus → **Status → Targets** 에서 kube-state-metrics up. **Explore** 에서 `{__name__=~"kube_.*"}` 또는 위 metric 이름으로 쿼리해 시리즈 존재 여부 확인.

### 6.2 node_exporter

| 검증 항목 | 확인 방법 | spec 기준 |
|-----------|-----------|-----------|
| CPU | `node_cpu_seconds_total` 존재, 라벨 `mode`, `instance` | mode="idle" |
| 메모리 | `node_memory_MemAvailable_bytes`, `node_memory_MemTotal_bytes` | instance |
| 디스크 | `node_filesystem_avail_bytes`, `node_filesystem_size_bytes` | mountpoint, fstype |

**확인 예:** Explore에서 `node_cpu_seconds_total`, `node_memory_MemAvailable_bytes` 등 조회. `instance` 라벨이 노드(또는 node_exporter)를 가리키는지 확인.

### 6.3 라벨 조정이 필요한 경우

- **instance vs node:** 일부 환경에서는 노드명이 `instance`가 아니라 `node` 라벨로 나올 수 있음. Table 패널에서 “노드 이름” 컬럼이 비어 있으면 라벨 매핑 확인.  
- **namespace:** Critical endpoint empty에서 critical 서비스만 보려면 `namespace=~"default|production"` 등으로 제한. 실제 critical 네임스페이스 목록을 팀에서 정한 뒤 반영.

---

## 7. Threshold 조정 노트 (Threshold tuning notes)

- **Block 1 (Stat):**  
  - NotReady, Workload Pending, Critical endpoint empty: **0 = OK, &gt;0 = Critical**. 변경 권장 없음.  
  - **Excessive restarts:** 설계상 “N 초과 시 비정상”. **N은 팀 정의**(예: 10, 20). 검증 시 현재 클러스터에서 정상 시 재시작 수를 잠깐 관찰한 뒤, N을 그에 맞게 설정.  
- **Block 2 (Gauge):**  
  - Node CPU, Node memory: **80% 경고, 95% 위험** 권장. 클러스터가 평소 70% 대면 80%는 적절. 평소 85% 대면 90/98 등으로 올릴 수 있음.  
  - Node disk: **90% 경고, 95% 위험** (사용률 기준). “여유 10% 미만”과 동일.  
  - Pending pods trend: Stat으로 쓸 때 **&gt;0 = 경고** 가능. Time series만 쓸 때는 참조선(예: 5, 10)을 두고 팀과 합의.  
- **Block 3 (Table):**  
  - CPU/Memory % 셀에 동일 80/95 threshold 적용 시 가독성 향상. Restart/Pending count는 팀에서 “이 정도면 조사” 기준을 정해 참조선으로 둘 수 있음.

검증 시 **실제 값 분포**를 보고 “너무 자주 경고” 또는 “거의 안 나옴”이면 threshold를 기록해 둔다.

---

## 8. Production 사용 가능 판정 기준 (Criteria for production usable)

다음이 모두 충족되면 **“Production usable”** 로 판정한다.

| # | 기준 | 확인 방법 |
|---|------|-----------|
| 1 | 대시보드가 Grafana에 **임포트되어** Row 1·2·3 구조대로 표시된다. | 시각 확인. |
| 2 | **모든 A 우선순위 패널(12개)** 에서 PromQL이 **에러 없이** 실행된다. | 각 패널 “Query” 또는 Explore에서 동일 쿼리 실행. |
| 3 | **필요 metric**이 Prometheus에 존재하며, metric/라벨 불일치가 있으면 **문서화·수정**되었다. | promql-spec / validation-results-template 에 실제 사용한 metric·라벨 기록. |
| 4 | **Block 1** 4개 패널이 “지금 안전한가?” 판단에 쓸 수 있는 **의미 있는 값**을 보여 준다. (0 = 정상, &gt;0 = 조사 필요 등) | 값 해석이 운영자에게 명확한지. |
| 5 | **Block 2** 4개 패널이 “조기 징후” 판단에 쓸 수 있는 **의미 있는 값**을 보여 준다. (%, 추세 등) | gauge/trend 해석 가능한지. |
| 6 | **Block 3** 4개 패널이 “어디를 조사할 것인가?”에 답할 수 있도록 **TOP10 또는 by-namespace** 결과를 보여 준다. | 테이블에 노드/파드/네임스페이스·값이 나오는지. |
| 7 | **Threshold**가 팀 정의 또는 권장과 맞거나, 조정 필요 사항이 **기록**되어 있다. | validation-results-template 에 threshold 메모. |
| 8 | **Refresh 1m 또는 2m** 시 Prometheus/Grafana 부하가 수용 가능하다. (선택: 지연·타임아웃 없음) | 부하 테스트 또는 일상 사용 관찰. |

**판정 결과:**  
- **Pass:** 위 기준 전부 충족 → 대시보드를 프로덕션에서 사용 가능하다고 판단.  
- **Pass with notes:** 일부 패널에서 metric 이름·라벨·threshold 조정을 했으나, 문서화되었고 운영적으로 사용 가능.  
- **Fail:** 하나 이상의 A 패널이 동작하지 않거나, 출력이 운영적으로 쓸 수 없음 → 원인 파악 후 수정·재검증.

---

## 9. 검증 수행 순서 요약

1. **준비:** Prometheus + kube-state-metrics + node_exporter 동작 확인. Grafana에 Prometheus datasource 등록.  
2. **임포트:** grafana-dashboard-v1.json 임포트. datasource 변수 연결.  
3. **패널별 검증:** validation-results-template.md 를 열고, 12개 A 패널 각각에 대해 “PromQL 동작 / metric 존재 / 조정 필요 / 의미 있음 / threshold” 5항목 기록.  
4. **실패 케이스 대응:** 5절·6절 참고해 metric 이름·라벨·필터 수정 후 재검증.  
5. **판정:** 8절 기준으로 Pass / Pass with notes / Fail 기록.  
6. **문서 반영:** 검증 결과를 validation-results-template.md 에 남기고, 필요 시 promql-spec.md, panel-config.md, grafana-dashboard-v1.json 에 환경별 수정 사항 반영.

---

*Validation plan v1. 검증 결과는 validation-results-template.md 에 기록한다.*
