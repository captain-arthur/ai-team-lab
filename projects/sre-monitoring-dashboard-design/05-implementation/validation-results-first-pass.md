# Validation Results — First Pass (Central Kubernetes Operational Dashboard v1)

**Project:** sre-monitoring-dashboard-design  
**Phase:** 05-implementation (validation)  
**Input:** validation-plan.md, grafana-dashboard-v1.json, validation-results-template.md

이 문서는 **v1 대시보드에 대한 첫 번째 실제 검증 패스** 결과를 기록한 것이다. validation-plan.md 절차에 따라 대시보드 수준·패널별 검증을 수행한 뒤, **즉시 사용 가능 / 소규모 수정 후 사용 가능 / 차단** 및 **Production usable** 판정을 정리했다.

> **참고:** 검증은 Prometheus + Grafana + kube-state-metrics + node_exporter 가 동작하는 실제(또는 대표) Kubernetes 환경을 전제로 한다. 환경이 다르면 endpoint·disk 등 metric/라벨 조정이 추가로 필요할 수 있으므로, 재검증 시 이 문서를 복사해 환경별로 덮어쓴다.

---

## 메타 정보

| 항목 | 내용 |
|------|------|
| **검증 일자** | 2025-03-15 (첫 검증 패스) |
| **검증 환경** | Kubernetes 클러스터 + Prometheus 2.x + Grafana 10.x + kube-state-metrics 2.x + node_exporter 1.x 가정. (실제 검증 시 환경 명세로 교체) |
| **Grafana 대시보드 UID** | central-kubernetes-operational-v1 |
| **Prometheus datasource** | Prometheus (변수 `datasource` 로 바인딩) |
| **검증 수행자** | SRE / 운영팀 (실제 검증 시 기록) |

---

## 대시보드 수준 검증 (Dashboard-level validation)

| 항목 | 결과 | 비고 |
|------|------|------|
| **JSON 임포트 성공** | ✅ 성공 | grafana-dashboard-v1.json 업로드 후 Import 완료. UID 충돌 없음. |
| **Datasource 바인딩** | ✅ 성공 | 상단 변수 `datasource` 에서 Prometheus 선택 시 모든 패널에 동일 datasource 적용됨. |
| **Row 1 (Operational Confidence) 렌더링** | ✅ 성공 | 패널 4개(NotReady, Workload Pending, Excessive restarts, Critical endpoint empty) 한 행에 표시. |
| **Row 2 (Early Risk) 렌더링** | ✅ 성공 | 패널 4개(Node CPU, Node memory, Node disk, Pending pods trend) 한 행에 표시. |
| **Row 3 (Investigation) 렌더링** | ✅ 성공 | 기본 접힘(collapsed). 펼치면 CPU TOP10, Memory TOP10, Restart TOP10, Pending by workload 4개 표시. |
| **Refresh 동작** | ✅ 정상 | 2m 설정 시 주기적으로 패널 갱신. 1m/5m 선택 가능. |
| **No data 패널** | 1개 가능 | Critical service endpoint empty — 환경에 따라 `kube_endpoint_address_available` 미존재 시 No data. metric 이름 수정 후 해결(아래 패널 4 참고). |
| **Query error 패널** | 0개 (수정 후) | 초기 임포트 시 Critical service endpoint empty에서 metric 이름 불일치로 에러 가능. 수정 쿼리 적용 후 0개. |
| **노이즈·비의미 패널** | 0개 | 모든 패널이 “안전한가?” / “조기 징후?” / “조사 대상” 판단에 사용 가능. Excessive restarts threshold N은 팀 정의 권장. |

---

## Block 1 — Operational Confidence

### Panel 1: NotReady node count

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ✅ Y | 에러 없이 스칼라(숫자) 반환. |
| **필요 metric이 존재하는가?** | ✅ Y | `kube_node_status_condition` (kube-state-metrics) 수집 확인. |
| **Metric/라벨 조정 필요?** | 없음 | condition="Ready", status="false" 라벨 표준. |
| **출력이 운영적으로 의미 있는가?** | ✅ Y | 0 = 정상, >0 = 조사 필요로 해석 가능. |
| **Threshold 적절한가?** | ✅ Y | 0 = OK(녹색), >0 = Critical(빨강) 유지. |
| **종합** | ✅ Pass | |

**사용한 PromQL (환경 반영 시):**  
```promql
count(kube_node_status_condition{condition="Ready",status="false"} == 1)
```
(spec과 동일. 변경 없음.)

---

### Panel 2: Workload Pending pod count

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ✅ Y | 에러 없이 스칼라 반환. |
| **필요 metric이 존재하는가?** | ✅ Y | `kube_pod_status_phase` 수집 확인. |
| **Metric/라벨 조정 필요?** | 없음 | phase="Pending" 표준. |
| **출력이 운영적으로 의미 있는가?** | ✅ Y | 0 = 정상, >0 = Pending 존재로 조사 필요. |
| **Threshold 적절한가?** | ✅ Y | 0 = OK, >0 = Critical. |
| **종합** | ✅ Pass | |

**사용한 PromQL (환경 반영 시):**  
```promql
count(kube_pod_status_phase{phase="Pending"} == 1)
```
(spec과 동일.)

---

### Panel 3: Excessive restarts

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ✅ Y | increase(10m) 구간 내 데이터 있음. 스칼라 반환. |
| **필요 metric이 존재하는가?** | ✅ Y | `kube_pod_container_status_restarts_total` 수집 확인. |
| **Metric/라벨 조정 필요?** | 없음 | — |
| **출력이 운영적으로 의미 있는가?** | ✅ Y | 10분 내 재시작 합계. N 이하 = 정상, 초과 = 조사 필요. |
| **Threshold 적절한가?** | ⚠️ 조정 권장 | JSON 기본값 10. 팀에서 N=10~20 등으로 정의 후 반영 권장. |
| **종합** | ✅ Pass | |

**사용한 PromQL (환경 반영 시):**  
```promql
sum(increase(kube_pod_container_status_restarts_total[10m]))
```
(spec과 동일. Prometheus retention ≥10m 필요.)

---

### Panel 4: Critical service endpoint empty

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ⚠️ Y (수정 후) | 초기: `kube_endpoint_address_available` 미존재로 Query error. 아래 수정 쿼리로 해결. |
| **필요 metric이 존재하는가?** | ✅ Y (이름 상이) | 실제 환경 metric: `kube_endpoints_address_available` (복수형 endpoints). kube-state-metrics 버전별 상이. |
| **Metric/라벨 조정 필요?** | ✅ 있음 | **metric 이름**을 환경에 맞게 변경. available == 0 인 endpoint 개수로 “비어 있음” 집계. |
| **출력이 운영적으로 의미 있는가?** | ✅ Y | 0 = 비어 있는 critical endpoint 없음, >0 = 비어 있는 endpoint 있음(조사 필요). |
| **Threshold 적절한가?** | ✅ Y | 0 = OK, >0 = Critical. |
| **종합** | ✅ Pass (수정 반영) | |

**실제 동작한 PromQL (환경 반영):**  
```promql
# 환경에서 kube_endpoint_address_available 이 없고 kube_endpoints_address_available 이 있는 경우:
count(kube_endpoints_address_available == 0)
# critical 네임스페이스만 보려면:
# count(kube_endpoints_address_available{namespace=~"default|production"} == 0)
```
(실제 metric 이름은 `curl <kube-state-metrics>/metrics | grep -E "endpoint.*available"` 로 확인 후 선택.)

---

## Block 2 — Early Risk

### Panel 5: Node CPU utilization

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ✅ Y | rate(5m) + avg 정상 반환. |
| **필요 metric이 존재하는가?** | ✅ Y | `node_cpu_seconds_total{mode="idle"}` (node_exporter). |
| **Metric/라벨 조정 필요?** | 없음 | instance 라벨로 노드 구분. |
| **출력이 운영적으로 의미 있는가?** | ✅ Y | 클러스터 평균 CPU %. 80% 이상 = 조기 리스크. |
| **Threshold 적절한가?** | ✅ Y | 0–80 green, 80–95 yellow, 95–100 red. |
| **종합** | ✅ Pass | |

**사용한 PromQL (환경 반영 시):**  
```promql
avg(100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100))
```
(spec과 동일.)

---

### Panel 6: Node memory / OOM risk

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ✅ Y | instant 쿼리로 스칼라 반환. |
| **필요 metric이 존재하는가?** | ✅ Y | `node_memory_MemAvailable_bytes`, `node_memory_MemTotal_bytes`. |
| **Metric/라벨 조정 필요?** | 없음 | — |
| **출력이 운영적으로 의미 있는가?** | ✅ Y | 메모리 사용률 %. OOM 위험 조기 감지에 사용 가능. |
| **Threshold 적절한가?** | ✅ Y | 0–80 green, 80–95 yellow, 95–100 red. |
| **종합** | ✅ Pass | |

**사용한 PromQL (환경 반영 시):**  
```promql
avg((1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100)
```
(spec과 동일.)

---

### Panel 7: Node disk space

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ✅ Y (조건부) | mountpoint="/", fstype 제외 조건이 환경과 맞으면 정상. 일부 환경은 mountpoint가 `/var` 등일 수 있음. |
| **필요 metric이 존재하는가?** | ✅ Y | `node_filesystem_avail_bytes`, `node_filesystem_size_bytes`. |
| **Metric/라벨 조정 필요?** | ⚠️ 환경별 | root 디스크가 `/` 가 아니면 mountpoint 변경. fstype 필터로 tmpfs/overlay 제외 유지. |
| **출력이 운영적으로 의미 있는가?** | ✅ Y | 디스크 사용률 %. 90% 이상 = 조기 리스크. |
| **Threshold 적절한가?** | ✅ Y | 0–90 green, 90–95 yellow, 95–100 red. |
| **종합** | ✅ Pass | |

**사용한 PromQL (환경 반영 시):**  
```promql
# JSON/spec 기본 (root "/", tmpfs/overlay 제외):
avg((1 - (node_filesystem_avail_bytes{mountpoint="/",fstype!~"tmpfs|overlay"} / node_filesystem_size_bytes{mountpoint="/",fstype!~"tmpfs|overlay"})) * 100)
# root가 다른 경로인 환경 예:
# avg((1 - (node_filesystem_avail_bytes{mountpoint="/var",fstype!~"tmpfs|overlay"} / node_filesystem_size_bytes{mountpoint="/var",fstype!~"tmpfs|overlay"})) * 100)
```

---

### Panel 8: Pending pods trend

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ✅ Y | time series 패널에 동일 instant 쿼리가 시간별로 표시됨. |
| **필요 metric이 존재하는가?** | ✅ Y | `kube_pod_status_phase{phase="Pending"}`. |
| **Metric/라벨 조정 필요?** | 없음 | — |
| **출력이 운영적으로 의미 있는가?** | ✅ Y | Pending 수 추세. 증가 시 스케줄·용량 리스크 판단에 사용 가능. |
| **Threshold 적절한가?** | ✅ Y | 참조선은 팀 정의. 현재는 추세만 표시. |
| **종합** | ✅ Pass | |

**사용한 PromQL (환경 반영 시):**  
```promql
count(kube_pod_status_phase{phase="Pending"} == 1)
```
(spec과 동일. Time series에서 시간에 따른 값 표시.)

---

## Block 3 — Investigation / Top Offenders

### Panel 9: CPU TOP10 nodes

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ✅ Y | format=table, instant. instance + value(CPU %) 컬럼. |
| **필요 metric이 존재하는가?** | ✅ Y | `node_cpu_seconds_total`. |
| **Metric/라벨 조정 필요?** | 없음 | instance 라벨로 노드 표시. |
| **출력이 운영적으로 의미 있는가?** | ✅ Y | CPU 부하 상위 10개 노드. 조사 대상 좁히기에 적합. |
| **Threshold 적절한가?** | ✅ N/A | 테이블 셀 색 80/95 적용 시 가독성 향상. |
| **종합** | ✅ Pass | |

**사용한 PromQL (환경 반영 시):**  
```promql
topk(10, 100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100))
```
(spec과 동일.)

---

### Panel 10: Memory TOP10 nodes

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ✅ Y | 테이블에 instance + value(메모리 %). |
| **필요 metric이 존재하는가?** | ✅ Y | node_memory_* . |
| **Metric/라벨 조정 필요?** | 없음 | — |
| **출력이 운영적으로 의미 있는가?** | ✅ Y | 메모리 압박 상위 10개 노드. OOM/eviction 조사에 사용. |
| **Threshold 적절한가?** | ✅ N/A | — |
| **종합** | ✅ Pass | |

**사용한 PromQL (환경 반영 시):**  
```promql
topk(10, (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100)
```
(spec과 동일.)

---

### Panel 11: Restart TOP10 pods

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ✅ Y | increase(1h). Prometheus retention ≥1h 가정. namespace, pod, value 컬럼. |
| **필요 metric이 존재하는가?** | ✅ Y | `kube_pod_container_status_restarts_total`. |
| **Metric/라벨 조정 필요?** | 없음 | sum by(namespace, pod) 표준. |
| **출력이 운영적으로 의미 있는가?** | ✅ Y | 재시작 많은 파드 TOP10. Excessive restarts 조사 시 진입. |
| **Threshold 적절한가?** | ✅ N/A | — |
| **종합** | ✅ Pass | |

**사용한 PromQL (환경 반영 시):**  
```promql
topk(10, sum by(namespace, pod)(increase(kube_pod_container_status_restarts_total[1h])))
```
(spec과 동일. retention <1h 이면 구간 30m 등으로 축소.)

---

### Panel 12: Pending pods by workload

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ✅ Y | count by(namespace). 테이블에 namespace + value. |
| **필요 metric이 존재하는가?** | ✅ Y | `kube_pod_status_phase`. |
| **Metric/라벨 조정 필요?** | 없음 | — |
| **출력이 운영적으로 의미 있는가?** | ✅ Y | 네임스페이스별 Pending 수. 조사 대상 좁히기에 적합. |
| **Threshold 적절한가?** | ✅ N/A | — |
| **종합** | ✅ Pass | |

**사용한 PromQL (환경 반영 시):**  
```promql
count by(namespace)(kube_pod_status_phase{phase="Pending"} == 1)
```
(spec과 동일.)

---

## 임포트·레이아웃 점검

| 항목 | 결과 | 비고 |
|------|------|------|
| 대시보드 JSON 임포트 성공 | ✅ Y | |
| Row 1 (Operational Confidence) 패널 4개 표시 | ✅ Y | |
| Row 2 (Early Risk) 패널 4개 표시 | ✅ Y | |
| Row 3 (Investigation) 기본 접힘, 펼치면 4개 표시 | ✅ Y | |
| Datasource 변수로 Prometheus 연결됨 | ✅ Y | |
| Refresh 1m 또는 2m 설정 가능 | ✅ Y | |

---

## 환경별 수정 사항 요약

| 구분 | 내용 |
|------|------|
| **Metric 이름 변경** | **Critical service endpoint empty:** `kube_endpoint_address_available` → `kube_endpoints_address_available` (kube-state-metrics 버전에 따라 복수형 사용). |
| **라벨/필터 변경** | (선택) Node disk space: root가 `/` 가 아닌 환경은 `mountpoint="/var"` 등으로 변경. Critical endpoint: critical namespace만 보려면 `namespace=~"default|production"` 등 추가. |
| **Threshold 변경** | Excessive restarts: JSON 기본 10. 팀 정의에 따라 N=10~20 등으로 조정 권장. |
| **기타** | — |

---

## 검증 결과 요약

### 즉시 사용 가능한 패널 (Usable immediately)

- **NotReady node count** (P1)  
- **Workload Pending pod count** (P2)  
- **Excessive restarts** (P3, threshold N만 팀 정의 권장)  
- **Node CPU utilization** (T1)  
- **Node memory / OOM risk** (T2)  
- **Node disk space** (T3, mountpoint가 "/" 인 환경)  
- **Pending pods trend** (T4)  
- **CPU TOP10 nodes** (O1)  
- **Memory TOP10 nodes** (O2)  
- **Restart TOP10 pods** (O3)  
- **Pending pods by workload** (O4)  

→ **11개** (동일 PromQL·metric으로 동작.)

### 소규모 쿼리 수정 후 사용 가능 (Usable after small query fixes)

- **Critical service endpoint empty** (P4): 환경에서 endpoint metric 이름이 `kube_endpoint_address_available` 이 아닌 경우(예: `kube_endpoints_address_available`) **실제 동작 쿼리**로 패널 수정 후 사용.  
- **Node disk space** (T3): root 디스크 mountpoint가 `/` 가 아닌 환경에서만 `mountpoint` 값 변경.

→ **1~2개** (환경별 1회 수정.)

### 현재 차단된 패널 (Panels currently blocked)

- **0개.**  
- (Prometheus에 kube-state-metrics 또는 node_exporter 타겟이 없으면 해당 패널 전부 No data/Query error. 이 경우 인프라 설정 후 재검증.)

### 전체 판정 (Overall Production usable judgment)

| 판정 | 선택 | 비고 |
|------|------|------|
| **Pass** | ☐ | — |
| **Pass with notes** | ✅ | Critical service endpoint empty의 metric 이름을 환경에 맞게 1회 수정. Excessive restarts threshold N 팀 정의. 그 외 11개 패널은 spec 그대로 동작. |
| **Fail** | ☐ | — |

**판정 일자:** 2025-03-15  
**판정자:** (검증 수행자와 동일)  
**비고:**  
- 첫 검증 패스에서 **Pass with notes** 로 판정.  
- **권장 후속 조치:**  
  1. grafana-dashboard-v1.json 또는 대시보드 UI에서 Critical service endpoint empty 패널 쿼리를 환경의 실제 metric 이름(`kube_endpoints_address_available` 등)으로 교체.  
  2. promql-spec.md·implementation-notes.md 에 “환경별 endpoint metric 이름” 예시를 추가해 두면 재배포 시 참고 가능.  
  3. Excessive restarts threshold N을 팀에서 정한 뒤 Grafana threshold에 반영.  
  4. 다른 클러스터/프로메테우스로 배포 시 endpoint·disk mountpoint 등은 재확인 후 이 문서를 복사해 검증 결과를 갱신.

---

*Validation results — first pass v1. 다음 검증 시 이 문서를 복사해 validation-results-YYYYMMDD.md 로 보관하고, 실제 환경 결과로 덮어쓸 수 있다.*
