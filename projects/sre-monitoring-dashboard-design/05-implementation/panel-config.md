# Panel Configuration — Central Kubernetes Operational Dashboard v1

**Project:** sre-monitoring-dashboard-design  
**Phase:** 05-implementation  
**Input:** implementation-ready-panel-spec, promql-spec

이 문서는 **A 우선순위 패널 12개**에 대한 **Grafana 패널 설정**을 정의한다. 패널 타입, 단위, threshold, value mapping, 색상 규칙, 새로고침 주기, 레이아웃 힌트를 정리하며, Block 1 / Block 2 / Block 3 소속을 명시한다.

---

## 1. 블록별 패널 소속

| Block | 블록 이름 | 패널 (A 우선순위) |
|-------|-----------|-------------------|
| **Block 1** | Operational Confidence | NotReady node count, Workload Pending pod count, Excessive restarts, Critical service endpoint empty |
| **Block 2** | Early Risk | Node CPU utilization, Node memory / OOM risk, Node disk space, Pending pods trend |
| **Block 3** | Investigation / Top Offenders | CPU TOP10 nodes, Memory TOP10 nodes, Restart TOP10 pods, Pending pods by workload |

---

## 2. Block 1 — Operational Confidence

### 2.1 NotReady node count

| 설정 항목 | 값 |
|-----------|-----|
| **Block** | Block 1 — Operational Confidence |
| **Panel type** | stat |
| **Unit** | none (short) |
| **Thresholds** | Base: 0. Green: 0. Red: 0.01 (즉 >0 이면 red). 또는 Custom: 0 = green, 1 = red. |
| **Value mappings** | 0 → "OK" (green 또는 text). 1+ → "Critical" (red). 또는 Null → "N/A". |
| **Color rules** | Value &lt;= 0 → green. Value &gt; 0 → red. |
| **Refresh interval** | 대시보드 기본(권장 1m~2m). |
| **Panel width / layout** | Grid: 너비 4~5 (24 그리드 기준). Block 1 한 행에 4개 패널이면 각 6. |

---

### 2.2 Workload Pending pod count

| 설정 항목 | 값 |
|-----------|-----|
| **Block** | Block 1 — Operational Confidence |
| **Panel type** | stat |
| **Unit** | none (short) |
| **Thresholds** | 0 = green, 0.01 = red (즉 >0 red). |
| **Value mappings** | 0 → "OK". &gt;0 → "Critical". |
| **Color rules** | Value &lt;= 0 → green. Value &gt; 0 → red. |
| **Refresh interval** | 대시보드 기본. |
| **Panel width / layout** | Grid 4~6. Block 1 내 일관. |

---

### 2.3 Excessive restarts

| 설정 항목 | 값 |
|-----------|-----|
| **Block** | Block 1 — Operational Confidence |
| **Panel type** | stat |
| **Unit** | none (short) |
| **Thresholds** | 팀 정의 N. 예: 0–10 green, 10–100 yellow, 100+ red. (N=10 기준: Base 10, green 0–10, red 10.01+) |
| **Value mappings** | (선택) 0 → "OK". N 미만 → "OK", N 이상 → "Critical". |
| **Color rules** | Value &lt; N → green. Value &gt;= N → red (또는 yellow 구간 추가). |
| **Refresh interval** | 대시보드 기본. |
| **Panel width / layout** | Grid 4~6. |

---

### 2.4 Critical service endpoint empty

| 설정 항목 | 값 |
|-----------|-----|
| **Block** | Block 1 — Operational Confidence |
| **Panel type** | stat |
| **Unit** | none (short) |
| **Thresholds** | 0 = green. 0.01 = red (&gt;0 이면 비어 있는 endpoint 있음). |
| **Value mappings** | 0 → "OK". &gt;0 → "Critical". |
| **Color rules** | Value &lt;= 0 → green. Value &gt; 0 → red. |
| **Refresh interval** | 대시보드 기본. |
| **Panel width / layout** | Grid 4~6. |

---

## 3. Block 2 — Early Risk

### 3.1 Node CPU utilization

| 설정 항목 | 값 |
|-----------|-----|
| **Block** | Block 2 — Early Risk |
| **Panel type** | gauge (권장) 또는 stat |
| **Unit** | percent (0–100) |
| **Thresholds** | 0–80 green. 80–95 yellow (또는 orange). 95–100 red. |
| **Value mappings** | (gauge에서는 보통 사용 안 함. Stat이면 80 미만 "OK", 80+ "Warning" 등) |
| **Color rules** | Value 0–80 green, 80–95 yellow, 95–100 red. |
| **Refresh interval** | 대시보드 기본 (1m~2m). |
| **Panel width / layout** | Grid 4~6. Block 2 한 행에 4개면 각 6. |

---

### 3.2 Node memory / OOM risk

| 설정 항목 | 값 |
|-----------|-----|
| **Block** | Block 2 — Early Risk |
| **Panel type** | gauge (권장) 또는 stat |
| **Unit** | percent (0–100) |
| **Thresholds** | 0–80 green. 80–95 yellow. 95–100 red. |
| **Value mappings** | (선택) |
| **Color rules** | 0–80 green, 80–95 yellow, 95–100 red. |
| **Refresh interval** | 대시보드 기본. |
| **Panel width / layout** | Grid 4~6. |

---

### 3.3 Node disk space

| 설정 항목 | 값 |
|-----------|-----|
| **Block** | Block 2 — Early Risk |
| **Panel type** | gauge (권장) 또는 stat |
| **Unit** | percent (0–100). (사용률이므로 90 = 90% 사용 = 10% 여유) |
| **Thresholds** | 0–90 green. 90–95 yellow. 95–100 red. |
| **Value mappings** | (선택) |
| **Color rules** | 0–90 green, 90–95 yellow, 95–100 red. |
| **Refresh interval** | 대시보드 기본. |
| **Panel width / layout** | Grid 4~6. |

---

### 3.4 Pending pods trend

| 설정 항목 | 값 |
|-----------|-----|
| **Block** | Block 2 — Early Risk |
| **Panel type** | time series (권장, 추세 가시화) 또는 stat |
| **Unit** | none (short). Time series면 Y축 단위 short. |
| **Thresholds** | Stat일 때: 0 = green, 0.01 = yellow 또는 red. Time series일 때 Y축에 0, 5, 10 등 참조선 가능. |
| **Value mappings** | Stat일 때만. 0 → "OK", &gt;0 → "Warning". |
| **Color rules** | Time series: 단일 시리즈면 기본 색. 영역 채우기 권장. |
| **Refresh interval** | 대시보드 기본. |
| **Panel width / layout** | Grid 8~12 (가로로 넓게 해서 추세가 보이도록). 또는 Stat이면 4~6. |

---

## 4. Block 3 — Investigation / Top Offenders

### 4.1 CPU TOP10 nodes

| 설정 항목 | 값 |
|-----------|-----|
| **Block** | Block 3 — Investigation |
| **Panel type** | table |
| **Unit** | 컬럼 Value: percent (0–100). |
| **Thresholds** | (테이블 셀 색) Value 컬럼: 0–80 green, 80–95 yellow, 95–100 red. |
| **Value mappings** | (테이블에서는 보통 생략) |
| **Color rules** | Table cell display: Value 기준 threshold 적용. |
| **Refresh interval** | 대시보드 기본. Block 3은 드릴다운이므로 1m~2m. |
| **Panel width / layout** | Grid 12. Row 3에서 다른 Block 3 패널과 나란히 또는 세로 배치. |

**Table 컬럼:**  
- `instance` (또는 Node) — 라벨 from query.  
- `Value` — 쿼리 결과 값(CPU %). 정렬: 내림차순.  
- (선택) `Value`를 "CPU %" 등으로 표시 이름 변경.

---

### 4.2 Memory TOP10 nodes

| 설정 항목 | 값 |
|-----------|-----|
| **Block** | Block 3 — Investigation |
| **Panel type** | table |
| **Unit** | Value 컬럼: percent (0–100). |
| **Thresholds** | Value: 0–80 green, 80–95 yellow, 95–100 red. |
| **Value mappings** | — |
| **Color rules** | Table cell by value. |
| **Refresh interval** | 대시보드 기본. |
| **Panel width / layout** | Grid 12. |

**Table 컬럼:** instance, Value (메모리 %). 내림차순.

---

### 4.3 Restart TOP10 pods

| 설정 항목 | 값 |
|-----------|-----|
| **Block** | Block 3 — Investigation |
| **Panel type** | table |
| **Unit** | Value 컬럼: none (short). |
| **Thresholds** | (선택) 0–5 green, 5–20 yellow, 20+ red. |
| **Value mappings** | — |
| **Color rules** | (선택) Value 기준. |
| **Refresh interval** | 대시보드 기본. |
| **Panel width / layout** | Grid 12. |

**Table 컬럼:** namespace, pod, Value (restarts). 내림차순.

---

### 4.4 Pending pods by workload

| 설정 항목 | 값 |
|-----------|-----|
| **Block** | Block 3 — Investigation |
| **Panel type** | table |
| **Unit** | Value 컬럼: none (short). |
| **Thresholds** | (선택) 0 green, &gt;0 yellow/red. |
| **Value mappings** | — |
| **Color rules** | (선택) |
| **Refresh interval** | 대시보드 기본. |
| **Panel width / layout** | Grid 12. |

**Table 컬럼:** namespace (또는 workload), Value (pending count). 내림차순.

---

## 5. 공통 설정 요약

| 항목 | 권장값 |
|------|--------|
| **Dashboard refresh** | 1m 또는 2m. Prometheus 부하 고려. |
| **Block 1 row** | 제목: "Operational Confidence" (또는 "운영 확신"). 패널 4개, 동일 행. |
| **Block 2 row** | 제목: "Early Risk" (또는 "조기 리스크"). 패널 4개, 동일 행. |
| **Block 3 row** | 제목: "Investigation / Top Offenders". **collapsed by default**. 패널 4개. |
| **Datasource** | `${datasource}` 또는 "Prometheus" (환경 변수 권장). |
| **Panel grid width** | Block 1·2: 각 패널 6 (24/4). Block 3: 각 12 또는 6. |

---

*Panel config v1. A-priority panels only. Grafana 8.x/9.x 기준.*
