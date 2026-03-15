# Validation Results — Central Kubernetes Operational Dashboard v1

**Project:** sre-monitoring-dashboard-design  
**Phase:** 05-implementation (validation)  
**Input:** validation-plan.md, grafana-dashboard-v1.json

이 문서는 **실제 Prometheus/Grafana 환경**에서 v1 대시보드를 검증한 **결과를 기록**하는 템플릿이다. 검증 시 각 패널별로 아래 표를 채우고, 마지막에 **Production usable** 판정을 기록한다.

---

## 메타 정보 (검증 시 채우기)

| 항목 | 내용 |
|------|------|
| **검증 일자** | (예: 2025-03-15) |
| **검증 환경** | (예: EKS 클러스터, Prometheus 2.x, Grafana 10.x, kube-state-metrics 2.x, node_exporter 1.x) |
| **Grafana 대시보드 UID** | (임포트 후 UID 또는 동일: central-kubernetes-operational-v1) |
| **Prometheus datasource** | (사용한 datasource 이름/UID) |
| **검증 수행자** | (이름 또는 팀) |

---

## Block 1 — Operational Confidence

### Panel 1: NotReady node count

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ☐ Y / ☐ N / ☐ N/A | (에러 메시지, 수정한 쿼리 등) |
| **필요 metric이 존재하는가?** | ☐ Y / ☐ N / ☐ N/A | metric: `kube_node_status_condition` |
| **Metric/라벨 조정 필요?** | ☐ 없음 / ☐ 있음 | (실제 사용한 metric·라벨) |
| **출력이 운영적으로 의미 있는가?** | ☐ Y / ☐ N / ☐ N/A | (0=정상 해석 가능 여부) |
| **Threshold 적절한가?** | ☐ Y / ☐ N / ☐ 조정함 | (0=OK, >0=Critical. 조정 시 값 기록) |
| **종합** | ☐ Pass / ☐ Fail | |

**사용한 PromQL (환경 반영 시):**  
```
(기록)
```

---

### Panel 2: Workload Pending pod count

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ☐ Y / ☐ N / ☐ N/A | |
| **필요 metric이 존재하는가?** | ☐ Y / ☐ N / ☐ N/A | metric: `kube_pod_status_phase` |
| **Metric/라벨 조정 필요?** | ☐ 없음 / ☐ 있음 | |
| **출력이 운영적으로 의미 있는가?** | ☐ Y / ☐ N / ☐ N/A | |
| **Threshold 적절한가?** | ☐ Y / ☐ N / ☐ 조정함 | |
| **종합** | ☐ Pass / ☐ Fail | |

**사용한 PromQL (환경 반영 시):**  
```
(기록)
```

---

### Panel 3: Excessive restarts

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ☐ Y / ☐ N / ☐ N/A | (10m increase, retention 확인) |
| **필요 metric이 존재하는가?** | ☐ Y / ☐ N / ☐ N/A | metric: `kube_pod_container_status_restarts_total` |
| **Metric/라벨 조정 필요?** | ☐ 없음 / ☐ 있음 | |
| **출력이 운영적으로 의미 있는가?** | ☐ Y / ☐ N / ☐ N/A | |
| **Threshold 적절한가?** | ☐ Y / ☐ N / ☐ 조정함 | (N 값: __ ) |
| **종합** | ☐ Pass / ☐ Fail | |

**사용한 PromQL (환경 반영 시):**  
```
(기록)
```

---

### Panel 4: Critical service endpoint empty

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ☐ Y / ☐ N / ☐ N/A | (endpoint metric 이름 확인) |
| **필요 metric이 존재하는가?** | ☐ Y / ☐ N / ☐ N/A | 실제 metric 이름: _____________ |
| **Metric/라벨 조정 필요?** | ☐ 없음 / ☐ 있음 | (kube_endpoint_* vs kube_endpoints_* 등) |
| **출력이 운영적으로 의미 있는가?** | ☐ Y / ☐ N / ☐ N/A | |
| **Threshold 적절한가?** | ☐ Y / ☐ N / ☐ 조정함 | |
| **종합** | ☐ Pass / ☐ Fail | |

**사용한 PromQL (환경 반영 시):**  
```
(기록)
```

---

## Block 2 — Early Risk

### Panel 5: Node CPU utilization

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ☐ Y / ☐ N / ☐ N/A | |
| **필요 metric이 존재하는가?** | ☐ Y / ☐ N / ☐ N/A | metric: `node_cpu_seconds_total` |
| **Metric/라벨 조정 필요?** | ☐ 없음 / ☐ 있음 | |
| **출력이 운영적으로 의미 있는가?** | ☐ Y / ☐ N / ☐ N/A | (gauge % 해석 가능 여부) |
| **Threshold 적절한가?** | ☐ Y / ☐ N / ☐ 조정함 | (80/95 또는 조정 값) |
| **종합** | ☐ Pass / ☐ Fail | |

**사용한 PromQL (환경 반영 시):**  
```
(기록)
```

---

### Panel 6: Node memory / OOM risk

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ☐ Y / ☐ N / ☐ N/A | |
| **필요 metric이 존재하는가?** | ☐ Y / ☐ N / ☐ N/A | `node_memory_MemAvailable_bytes`, `node_memory_MemTotal_bytes` |
| **Metric/라벨 조정 필요?** | ☐ 없음 / ☐ 있음 | |
| **출력이 운영적으로 의미 있는가?** | ☐ Y / ☐ N / ☐ N/A | |
| **Threshold 적절한가?** | ☐ Y / ☐ N / ☐ 조정함 | |
| **종합** | ☐ Pass / ☐ Fail | |

**사용한 PromQL (환경 반영 시):**  
```
(기록)
```

---

### Panel 7: Node disk space

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ☐ Y / ☐ N / ☐ N/A | (mountpoint, fstype 필터) |
| **필요 metric이 존재하는가?** | ☐ Y / ☐ N / ☐ N/A | `node_filesystem_avail_bytes`, `node_filesystem_size_bytes` |
| **Metric/라벨 조정 필요?** | ☐ 없음 / ☐ 있음 | (mountpoint="/" 또는 다른 경로) |
| **출력이 운영적으로 의미 있는가?** | ☐ Y / ☐ N / ☐ N/A | |
| **Threshold 적절한가?** | ☐ Y / ☐ N / ☐ 조정함 | (90/95 또는 조정) |
| **종합** | ☐ Pass / ☐ Fail | |

**사용한 PromQL (환경 반영 시):**  
```
(기록)
```

---

### Panel 8: Pending pods trend

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ☐ Y / ☐ N / ☐ N/A | (time series 또는 stat) |
| **필요 metric이 존재하는가?** | ☐ Y / ☐ N / ☐ N/A | `kube_pod_status_phase` |
| **Metric/라벨 조정 필요?** | ☐ 없음 / ☐ 있음 | |
| **출력이 운영적으로 의미 있는가?** | ☐ Y / ☐ N / ☐ N/A | (추세 또는 현재값 해석 가능) |
| **Threshold 적절한가?** | ☐ Y / ☐ N / ☐ 조정함 | |
| **종합** | ☐ Pass / ☐ Fail | |

**사용한 PromQL (환경 반영 시):**  
```
(기록)
```

---

## Block 3 — Investigation / Top Offenders

### Panel 9: CPU TOP10 nodes

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ☐ Y / ☐ N / ☐ N/A | (format=table, instant) |
| **필요 metric이 존재하는가?** | ☐ Y / ☐ N / ☐ N/A | |
| **Metric/라벨 조정 필요?** | ☐ 없음 / ☐ 있음 | (instance vs node 등) |
| **출력이 운영적으로 의미 있는가?** | ☐ Y / ☐ N / ☐ N/A | (노드명 + CPU % 테이블) |
| **Threshold 적절한가?** | ☐ Y / ☐ N / ☐ N/A | (테이블 셀 색) |
| **종합** | ☐ Pass / ☐ Fail | |

**사용한 PromQL (환경 반영 시):**  
```
(기록)
```

---

### Panel 10: Memory TOP10 nodes

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ☐ Y / ☐ N / ☐ N/A | |
| **필요 metric이 존재하는가?** | ☐ Y / ☐ N / ☐ N/A | |
| **Metric/라벨 조정 필요?** | ☐ 없음 / ☐ 있음 | |
| **출력이 운영적으로 의미 있는가?** | ☐ Y / ☐ N / ☐ N/A | |
| **Threshold 적절한가?** | ☐ Y / ☐ N / ☐ N/A | |
| **종합** | ☐ Pass / ☐ Fail | |

**사용한 PromQL (환경 반영 시):**  
```
(기록)
```

---

### Panel 11: Restart TOP10 pods

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ☐ Y / ☐ N / ☐ N/A | (1h increase, retention 확인) |
| **필요 metric이 존재하는가?** | ☐ Y / ☐ N / ☐ N/A | |
| **Metric/라벨 조정 필요?** | ☐ 없음 / ☐ 있음 | |
| **출력이 운영적으로 의미 있는가?** | ☐ Y / ☐ N / ☐ N/A | (namespace, pod, restarts) |
| **Threshold 적절한가?** | ☐ Y / ☐ N / ☐ N/A | |
| **종합** | ☐ Pass / ☐ Fail | |

**사용한 PromQL (환경 반영 시):**  
```
(기록)
```

---

### Panel 12: Pending pods by workload

| 검증 항목 | 결과 | 비고 |
|-----------|------|------|
| **PromQL이 실제 환경에서 동작하는가?** | ☐ Y / ☐ N / ☐ N/A | |
| **필요 metric이 존재하는가?** | ☐ Y / ☐ N / ☐ N/A | |
| **Metric/라벨 조정 필요?** | ☐ 없음 / ☐ 있음 | |
| **출력이 운영적으로 의미 있는가?** | ☐ Y / ☐ N / ☐ N/A | (namespace + count) |
| **Threshold 적절한가?** | ☐ Y / ☐ N / ☐ N/A | |
| **종합** | ☐ Pass / ☐ Fail | |

**사용한 PromQL (환경 반영 시):**  
```
(기록)
```

---

## 임포트·레이아웃 점검

| 항목 | 결과 | 비고 |
|------|------|------|
| 대시보드 JSON 임포트 성공 | ☐ Y / ☐ N | |
| Row 1 (Operational Confidence) 패널 4개 표시 | ☐ Y / ☐ N | |
| Row 2 (Early Risk) 패널 4개 표시 | ☐ Y / ☐ N | |
| Row 3 (Investigation) 기본 접힘, 펼치면 4개 표시 | ☐ Y / ☐ N | |
| Datasource 변수로 Prometheus 연결됨 | ☐ Y / ☐ N | |
| Refresh 1m 또는 2m 설정 가능 | ☐ Y / ☐ N | |

---

## 환경별 수정 사항 요약 (검증 중 반영한 경우)

| 구분 | 내용 |
|------|------|
| **Metric 이름 변경** | (예: kube_endpoint_address_available → kube_endpoints_address_available) |
| **라벨/필터 변경** | (예: mountpoint="/var", namespace=~"default\|prod") |
| **Threshold 변경** | (예: Excessive restarts N=20, Node disk 85/95) |
| **기타** | |

---

## Production usable 판정

**판정 기준:** validation-plan.md §8 참고.

| 판정 | 선택 |
|------|------|
| **Pass** — 모든 A 패널 동작, 출력 의미 있음, threshold 반영됨. 프로덕션 사용 가능. | ☐ |
| **Pass with notes** — 일부 metric/라벨/threshold 조정했으나 문서화됨. 프로덕션 사용 가능. | ☐ |
| **Fail** — 하나 이상 A 패널 미동작 또는 출력 비의미. 수정 후 재검증 필요. | ☐ |

**판정 일자:** _____________  
**판정자:** _____________  
**비고:**  
```
(자유 기술)
```

---

*Validation results template v1. 검증 후 이 파일을 복사해 실제 결과를 채운 버전을 별도 보관(예: validation-results-YYYYMMDD.md)할 수 있다.*
