# Kubernetes Cluster Health Monitoring Model

**문서 제목:** Kubernetes Cluster Health Monitoring Model  
**Project:** sre-monitoring-cluster-health-model  
**Program:** SRE Monitoring  
**Version:** v1

이 문서는 Kubernetes 클러스터의 **운영 건강(Cluster Health)** 을 일상적으로 파악하기 위한 **모니터링 모델**을 정의한다. 이후 **sre-monitoring-dashboard-design**, **sre-monitoring-alert-policy**, **sre-monitoring-operational-runbooks** 프로젝트의 기반이 된다.

---

## 1. 운영 문제 (Operational problem)

### 1.1 클러스터 건강 모니터링이 필요한 이유

운영자는 매일 **“지금 클러스터가 안전한가?”, “곧 불안전해질 징후는?”, “문제라면 어디를 먼저 볼 것인가?”** 에 답해야 한다. 클러스터는 노드, 제어면(control plane), 워크로드, 네트워크가 얽혀 있어, 한 곳의 이상이 용량 부족·스케줄 불가·서비스 장애로 이어진다. **공통된 건강 정의와 최소한의 신호**가 없으면, 팀마다 다른 기준으로 판단하게 되고, 이상 징후를 놓치거나 불필요한 알림에 시달리게 된다. 따라서 **클러스터 건강을 판단하는 공통 모델**이 필요하다.

### 1.2 대형 대시보드가 실패하는 이유

이미 Prometheus와 Grafana로 수많은 metric을 수집하고 대시보드를 두었더라도, **패널이 너무 많으면** 운영자가 5–10분 안에 “지금 안전한가?”를 결론 내리기 어렵다. 수십 개 패널을 훑는 동안 실제로 중요한 신호가 묻히고, **어디가 정상이고 어디가 비정상인지** 한눈에 보이지 않는다. 결과적으로 “매일 보는 central dashboard”가 **실제로는 쓰이지 않거나**, 중요한 판단은 여전히 **경험·스크립트·다른 채널**에 의존하게 된다. 즉, **많은 metric을 나열하는 것만으로는 일상 운영 인지가 해결되지 않는다.**

### 1.3 최소 신호 모델이 필요한 이유

**모든 것을 모니터링하지 않고**, “지금 안전한가?”에 직접 답하는 **소수의 상위 신호**만 두면, 운영자가 한 화면에서 **전부 정상 = 안전, 하나라도 비정상 = 조사 필요**로 해석할 수 있다. 또한 **곧 불안정해질 징후**(트렌드/리스크)와 **조사 대상 좁히기**(TOP10 등)를 레이어로 나누면, “요약 → 이상 시 드릴다운” 흐름이 명확해진다. 따라서 **최소이면서 의미 있는 신호 모델**이 필요하다. 이 모델은 **signal-to-noise** 를 높이고, 대시보드 과부하를 피하며, 5–10분 내 운영 판단을 가능하게 한다.

---

## 2. 세 가지 운영 질문 (Three operational questions)

모델은 다음 **세 가지 질문**에 1:1로 대응한다.

| # | 운영 질문 | 모델에서 답하는 방식 |
|---|-----------|----------------------|
| 1 | **지금 클러스터가 건강한가?** | Cluster Health Summary — 소수의 상위 신호만으로 “안전한가?” 판단 |
| 2 | **클러스터가 곧 불건강해질 징후가 있는가?** | Trend / Risk Indicators — 임계치 근접·추세로 조기 리스크 파악 |
| 3 | **문제가 있다면 어디를 먼저 조사할 것인가?** | Top Offenders / Drill-down — 노드·워크로드별 TOP10으로 원인 후보 좁히기 |

운영자는 **Summary → (이상 시) Trend 확인 → Top Offenders로 조사 대상 좁히기** 순서로 대시보드를 사용할 수 있다.

---

## 3. 모니터링 모델 구조 (Monitoring model structure)

신호는 **용도**에 따라 세 레이어로 구분한다.

```
┌─────────────────────────────────────────────────────────────────┐
│  Cluster Health Summary                                         │
│  "지금 클러스터가 안전한가?" — 5~7개 상위 신호                    │
├─────────────────────────────────────────────────────────────────┤
│  Trend / Risk Indicators                                        │
│  "곧 불안전해질 징후는?" — 임계치 근접·추세                       │
├─────────────────────────────────────────────────────────────────┤
│  Top Offenders / Drill-down                                      │
│  "어디를 먼저 볼 것인가?" — 노드·워크로드 TOP10                   │
└─────────────────────────────────────────────────────────────────┘
```

### 3.1 Cluster Health Summary

- **역할:** 운영자가 **한눈에** “지금 클러스터가 안전한가?”에 답하게 함.
- **원칙:** 신호 개수를 **매우 적게** 유지(전체 5~7개). 카테고리당 0~1개.
- **판정:** 이 레이어의 신호가 **전부 정상이면 안전**, **하나라도 비정상이면 조사 필요**로 해석.

### 3.2 Trend / Risk Indicators

- **역할:** “지금은 괜찮지만 **곧** 불건강해질 수 있다”는 징후를 보여 줌.
- **판정:** threshold 근접(예: 노드 CPU 80% 이상), 또는 시계열 **추세**(pending 증가, API server inflight 상승 등).
- **Summary와의 관계:** Summary가 모두 정상이어도 Trend에서 경고가 나올 수 있음. 용량·부하 사전 대응에 사용.

### 3.3 Top Offenders / Drill-down

- **역할:** Summary/Trend에서 “이상 있음”이 나왔을 때, **어떤 노드·워크로드**가 부담을 주는지 파악.
- **판정:** ranking 기반(예: CPU 사용률 TOP10 노드, 재시작 수 TOP10 워크로드). 각 뷰당 TOP10 수준만 보여도 5–10분 내 조사 대상 좁히기에 충분함.

---

## 4. 건강 카테고리 (Health categories)

모델은 **다섯 가지 카테고리**를 사용한다. 각 카테고리는 “해당 영역이 건강한가?”를 판단하는 단위다.

| 카테고리 | 담당 영역 | 건강의 의미 |
|----------|-----------|-------------|
| **Node / Infrastructure** | 노드 가용성, kubelet, 기본 리소스 | 노드가 Ready이고 스케줄 가능하며, 에이전트가 정상 동작함 |
| **Control Plane** | API server, scheduler, controller-manager, etcd | 제어면이 응답하고, 스케줄·조정이 정상 동작함 |
| **Workload / Pod** | 파드 상태, 재시작, Pending, eviction | 워크로드가 기동·스케줄되고, 과도한 재시작·eviction이 없음 |
| **Capacity / Resource pressure** | 클러스터 리소스 여유, 노드 사용률, Pending 추세 | 클러스터가 용량 한계에 있지 않고, 스케줄 여유가 있음 |
| **Network / Ingress** | 서비스·ingress 가용성, 엔드포인트 | 트래픽이 정상 전달되고, 백엔드가 비어 있지 않음 |

---

## 5. Core signals (핵심 신호)

### 5.1 Cluster Health Summary — 최소 집합

운영자가 **“지금 클러스터가 안전한가?”**만 답하려면 아래 **5~7개** 신호만 보면 된다. (환경에 따라 control plane metric이 없으면 4~5개로 축소.)

| # | Signal name | Category | 예시 Prometheus / PromQL | 판정(예시) |
|---|--------------|----------|---------------------------|------------|
| S1 | NotReady node count | Node / Infrastructure | `count(kube_node_status_condition{condition="Ready",status="false"} == 1)` | 0 = 정상, >0 = 비정상 |
| S2 | API server health (latency + 5xx) | Control Plane | P99: `histogram_quantile(0.99, sum(rate(apiserver_request_duration_seconds_bucket[5m])) by (le))`, 5xx: `sum(rate(apiserver_request_total{code=~"5.."}[5m]))` | P99 <1s, 5xx 극소 = 정상 |
| S3 | Scheduler pending pods | Control Plane | `scheduler_pending_pods` | 0 = 정상, >0 = 비정상 |
| S4 | Workload Pending pod count | Workload / Pod | `count(kube_pod_status_phase{phase="Pending"} == 1)` | 0 = 정상 |
| S5 | Excessive restarts (또는 Eviction) | Workload / Pod | `sum(increase(kube_pod_container_status_restarts_total[10m]))` (임계치 N 초과 시 비정상) | 임계치 이하 = 정상 |
| S6 | Critical service endpoint empty | Network / Ingress | `kube_endpoint_address_available` 등으로 핵심 서비스 필터 후 비어 있음 여부 | 비어 있지 않음 = 정상 |

- **Managed Kubernetes** 등에서 control plane(API server, scheduler)을 스크래핑하지 못하면 **S2, S3는 사용 불가.** 이 경우 **S1, S4, S5, S6** 만으로 Summary를 구성한다. 제어면 건강은 managed 서비스 콘솔·알림 등 **별도 채널**로 확인한다.

### 5.2 Trend / Risk — 최소 집합

**“클러스터가 곧 불건강해질 징후는?”** 에 답하는 소수 신호.

| # | Signal name | Category | 예시 Prometheus / PromQL | 판정(예시) |
|---|--------------|----------|---------------------------|------------|
| T1 | Node CPU / memory utilization | Node, Capacity | node_exporter: CPU·메모리 사용률 | >80% = 리스크 |
| T2 | API server inflight requests | Control Plane | `apiserver_current_inflight_requests` | 상승 추세 = 리스크 |
| T3 | Scheduler pending trend | Control Plane, Capacity | `scheduler_pending_pods` 시계열 또는 `increase(...[15m])` | 증가 = 리스크 |
| T4 | Node disk space | Node / Infrastructure | `node_filesystem_*`, root 등 mountpoint 필터 | 여유 부족(예: <10%) = 리스크 |

- Control plane metric이 없으면 **T1, T4** 만 사용. 노드·디스크 리스크는 이 둘로 커버 가능.

### 5.3 Top Offenders / Drill-down — 실무 뷰

**“어디를 먼저 조사할 것인가?”** 에 답하는 TOP10 스타일 뷰.

| # | View name | Category | 용도 |
|---|------------|----------|------|
| O1 | CPU TOP10 nodes | Node, Capacity | 노드 포화·eviction 원인 노드 |
| O2 | Memory TOP10 nodes | Node, Capacity | 메모리 압박·OOM 원인 노드 |
| O3 | Restart TOP10 workloads | Workload / Pod | 과도한 재시작이 나는 워크로드 |
| O4 | Pending pods by workload | Workload / Pod | Pending이 있는 네임스페이스/워크로드 |
| O5 | Error-heavy ingress / services | Network / Ingress | 에러율이 높은 ingress·서비스 (환경에 따라 사용 가능 시) |
| O6 | Control plane component health | Control Plane | 제어면 구성요소별 상태 (Self-managed 등에서만) |

- **즉시 사용 가능한 4개 뷰(O1–O4)** 만으로도 노드·워크로드 관점 조사 좁히기는 실용적이다. O5, O6은 Ingress controller·제어면 스크래핑 가능 여부에 따라 추가.

---

## 6. 실무적 제약 (Practical constraints)

### 6.1 Prometheus 쿼리 비용

- **histogram_quantile + rate(...[5m])** (API server P99 등)는 range 쿼리로 부하가 있을 수 있음. 대시보드 새로고침 주기를 **1–2분 이상**으로 두는 것을 권장.
- **increase(...[1h])** (재시작 TOP10 등)는 retention 구간 안에서만 유의미. retention이 짧으면 구간을 줄이거나 “현재 값만” threshold로 사용.
- **topk(10, ...)** 자체는 상대적으로 가벼우나, 내부 식이 무거우면 스캔 부하 가능. 필요 시 range·집계 구간을 조정.

### 6.2 Metric 가용성 차이

- **kube-state-metrics, node_exporter** 가 있으면 S1, S4, S5, T1, T4, O1–O4는 대부분 **즉시 사용 가능**.
- **Control plane (API server, scheduler)** metric은 **Managed K8s**에서는 사용자가 스크래핑하지 못해 **없을 수 있음.** 이 경우 Summary는 S1, S4, S5, S6으로 구성하고, 제어면은 별도 채널로 확인.
- **Eviction** 은 kubelet이 Prometheus metric으로 노출하지 않는 경우가 많아, v1에서는 **재시작만** Summary에 두고 Eviction은 후속에서 이벤트 등 다른 수단 검토.
- **Endpoint, Ingress** metric은 kube-state-metrics·Ingress controller 버전에 따라 **이름·라벨이 다를 수 있음.** 환경에서 metric을 확인한 뒤 PromQL을 조정해야 함.

### 6.3 Managed Kubernetes 환경

- EKS, GKE, AKS 등에서는 **control plane을 사용자가 스크래핑하지 못하는 경우**가 많다. 따라서 **S2, S3, T2, T3, O6** 은 사용 불가.
- **S4 Workload Pending pod count** 가 “스케줄되지 못한 워크로드가 있는가?”를 대표할 수 있어, **control plane metric 없이도** “스케줄·워크로드 관점 건강” 판단은 가능하다.
- 제어면 장애 자체는 **managed 서비스의 알림·콘솔·지원 채널**로 확인해야 하며, 이 모델의 Summary만으로는 “제어면 정상”을 보장할 수 없음 — 문서·runbook에 **“제어면은 별도 채널 확인”** 을 명시하는 것이 좋다.

---

## 7. 운영 워크플로 (Operational workflow)

### 7.1 운영자가 대시보드를 사용하는 방식

1. **일상 점검(5–10분):**  
   대시보드 **상단의 Cluster Health Summary** 만 본다. 5~7개(또는 4~5개) 신호가 **전부 정상**이면 “클러스터 안전”. **하나라도 비정상**이면 “조사 필요”.
2. **Trend 확인:**  
   **중간의 Trend / Risk** 섹션을 본다. Summary가 정상이어도 **노드 사용률 80% 근접, 디스크 여유 부족** 등이 있으면 “용량·정리·확장 검토” 등 사전 대응.
3. **이상 시 조사:**  
   Summary 또는 Trend에서 이상이 나오면, **해당 카테고리에 맞는 Top Offenders 뷰**로 들어간다. 예: 노드 문제 → O1/O2(CPU/메모리 TOP10 노드). 워크로드 문제 → O3/O4(재시작 TOP10, Pending by workload). **TOP10 목록**에서 1~2개 후보를 골라 kubectl·로그·이벤트 등으로 조사.

### 7.2 전형적인 조사 흐름

```
Summary에서 비정상 신호 발견
    → 어느 카테고리인지 확인 (Node / Workload / Network 등)
    → 해당 카테고리의 Top Offenders 뷰 열기
    → TOP10에서 가장 부담이 큰 노드/워크로드 1~2개 선택
    → kubectl describe, logs, events 등으로 원인 파악
```

- “비정상 신호 → 카테고리 → TOP10으로 후보 좁히기”가 한 번에 이루어지므로, **5–10분 안에** “어디를 먼저 볼 것인가?”에 답할 수 있다.

---

## 8. v1 한계 (v1 limitations)

### 8.1 Control plane 가시성

- **Managed K8s**에서는 API server·scheduler metric을 사용할 수 없어, **“제어면이 지금 응답·스케줄하는가?”**는 이 모델의 Summary로 보이지 않는다.  
- **대응:** Summary를 S1, S4, S5, S6으로 구성하고, 제어면 건강은 **managed 서비스 알림·콘솔** 등 별도 채널로 확인한다. 문서·runbook에 이 점을 명시한다.

### 8.2 Eviction 신호

- kubelet이 **eviction 횟수**를 Prometheus metric으로 노출하지 않는 경우가 많아, v1에서는 **Eviction을 Summary에 넣지 않는다.**  
- **대응:** 노드 메모리·CPU(T1, O2)로 “메모리 압박 가능성”은 간접적으로 본다. Eviction 자체를 숫자로 보려면 **후속에서 이벤트·다른 수단**을 검토한다.

### 8.3 환경별 쿼리 조정

- **S6 Critical service endpoint empty:** endpoint metric 이름·라벨이 kube-state-metrics 버전마다 다를 수 있음. 환경에서 실제 metric을 확인한 뒤 “비어 있는 엔드포인트” 조건을 PromQL로 정의해야 함.
- **O3 Restart TOP10:** deployment 등 **워크로드 단위**로 보려면 `kube_pod_owner` 등과 조인하거나 라벨을 추가해야 할 수 있음.
- **O4 Pending by workload:** namespace 단위만으로도 가능하나, deployment 등으로 세분화하려면 owner/라벨 조인 필요.
- **Trend “추세”:** “상승 추세”를 **어떤 time window, 어떤 증가량/비율**로 정의할지는 팀에서 정해야 함. 예: `scheduler_pending_pods` 또는 Pending count의 **최근 15분 증가량 > N** 등.

---

## 9. 향후 진화 (Future evolution)

이 모델은 **고정된 명세가 아니라**, 대시보드·알림·runbook과 함께 진화할 수 있다.

### 9.1 Dashboard 설계 (sre-monitoring-dashboard-design)

- **Cluster Health Summary** 를 대시보드 **상단 한 블록**에 5~7개 패널로 배치. 각 신호: 현재 값 + 정상/비정상(색·아이콘). “전부 정상 = 클러스터 안전” 한 문장으로 해석 가능하게.
- **Trend / Risk** 를 **중간 섹션**에 배치. 임계치 근접·추세 경고만 표시.
- **Top Offenders** 를 **하단 또는 별도 탭**에 배치. O1–O4 등 TOP10 테이블/그래프.
- 환경에 따라 **control plane 미사용 시** Summary 패널을 S1, S4, S5, S6만 보이도록 구성.

### 9.2 Alert 정책 (sre-monitoring-alert-policy)

- Summary 신호 중 **어떤 것을 알림으로 쓸지**, **임계치·심각도**(critical / warning)를 정한다. 예: S1 NotReady >0 → critical, T1 노드 사용률 >80% → warning.
- **노이즈가 많은 신호**(readiness probe 실패, 일시적 재시작 1~2회)는 Summary 알림에서 제외하고, drill-down·낮은 심각도만 사용한다는 원칙을 정책에 반영.

### 9.3 Operational runbooks (sre-monitoring-operational-runbooks)

- **“Summary에서 이 신호가 비정상일 때”** 조사 순서·체크리스트·명령어를 runbook으로 정리한다. 예: S1 비정상 → 노드 점검 runbook, S4 비정상 → Pending 원인 조사 runbook, S5 비정상 → 재시작 TOP10 확인 후 해당 워크로드 runbook.
- **Control plane 미사용 환경**에서는 “제어면 이상 의심 시 managed 서비스 콘솔·지원 채널 확인”을 runbook에 포함한다.

### 9.4 모델 확장

- **Eviction** 을 이벤트·다른 수단으로 수집할 수 있게 되면, Summary 또는 Trend에 추가 검토.
- **Ingress/서비스 에러율** (O5)을 Ingress controller metric이 지원하는 환경에서 Summary 7번째 또는 Trend로 도입 검토.
- **새 카테고리**가 필요해지면(예: 보안·비용), 기존 **세 레이어(Summary / Trend / Top Offenders)** 구조는 유지하고, 해당 카테고리당 0~1개 Summary 신호만 추가하는 방식으로 확장할 수 있다.

---

## 10. 요약

- **운영 문제:** 클러스터 건강을 빠르게 판단해야 하나, 대형 대시보드는 실효가 낮다. **최소 신호 모델**이 필요하다.
- **세 가지 질문:** 지금 건강한가? → Summary. 곧 불건강해질 징후는? → Trend. 어디를 먼저 볼 것인가? → Top Offenders.
- **모델 구조:** Cluster Health Summary(5~7개) / Trend-Risk(4~5개) / Top Offenders(5~6개 뷰). **건강 카테고리** 5개: Node, Control Plane, Workload/Pod, Capacity, Network.
- **실무 제약:** Prometheus 쿼리 비용·retention, metric 가용성 차이, Managed K8s에서 control plane 미사용을 고려해 **환경별로 Summary·Trend·Top Offenders 구성을 조정**한다.
- **v1 한계:** Control plane 가시성 부재(Managed K8s), Eviction 미반영, 환경별 쿼리 조정 필요. 문서·runbook에 **“제어면 별도 채널”, “환경별 수정 포인트”** 를 명시한다.
- **향후:** Dashboard 설계, Alert 정책, Runbooks 프로젝트에서 이 모델을 구체화하고, Eviction·에러율·새 카테고리는 필요 시 **레이어 구조를 유지한 채** 확장한다.

---

*Kubernetes Cluster Health Monitoring Model v1. SRE Monitoring program — sre-monitoring-cluster-health-model project.*
