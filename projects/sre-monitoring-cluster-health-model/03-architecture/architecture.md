# Cluster Health Monitoring Model — Architecture

**Project:** sre-monitoring-cluster-health-model  
**Phase:** Architecture  
**Input:** Manager brief, Research notes (`02-research/research-notes.md`)

---

## 1. 목적과 운영 질문

### 1.1 목적

이 아키텍처는 **일상 운영에 바로 쓸 수 있는 Cluster Health Monitoring Model**을 정의한다. 목표는 “많은 metric 나열”이 아니라, **최소이면서 의미 있는 신호 집합**으로 운영자가 세 가지 질문에 빠르게 답할 수 있게 하는 것이다.

### 1.2 세 가지 운영 질문

| # | 운영 질문 | 모델에서 답하는 방식 |
|---|-----------|----------------------|
| 1 | **지금 클러스터가 건강한가?** | Cluster Health Summary — 소수의 상위 신호만으로 “안전한가?” 판단 |
| 2 | **곧 불건강해질 징후가 있는가?** | Trend / Risk Indicators — 임계치 근접·추세로 조기 리스크 파악 |
| 3 | **문제가 있다면 어디를 먼저 조사할 것인가?** | Top Offenders / Drill-down — 노드·워크로드별 순위(TOP10)로 원인 후보 좁히기 |

---

## 2. 모델 구조: 세 가지 레이어

Cluster Health Monitoring Model은 신호를 **용도**에 따라 세 레이어로 구분한다.

```
┌─────────────────────────────────────────────────────────────────┐
│  Cluster Health Summary                                         │
│  "지금 클러스터가 안전한가?" — 소수 상위 신호만                    │
├─────────────────────────────────────────────────────────────────┤
│  Trend / Risk Indicators                                        │
│  "곧 불안전해질 징후는?" — 임계치 근접·추세                       │
├─────────────────────────────────────────────────────────────────┤
│  Top Offenders / Drill-down                                      │
│  "어디를 먼저 볼 것인가?" — 노드·워크로드 TOP10                   │
└─────────────────────────────────────────────────────────────────┘
```

### 2.1 Cluster Health Summary

- **역할:** 운영자가 **한눈에** “지금 클러스터가 안전한가?”에 답하게 함.
- **원칙:** 신호 개수를 **매우 적게** 유지. 카테고리당 0~1개, 전체 5~7개 수준을 목표로 함.
- **판정:** 이 레이어의 신호만 보고도 “전부 정상이면 안전, 하나라도 비정상이면 조사 필요”라고 해석할 수 있어야 함.

### 2.2 Trend / Risk Indicators

- **역할:** “지금은 괜찮지만 **곧** 불건강해질 수 있다”는 징후를 보여 줌.
- **판정 방식:** threshold 근접(예: 노드 CPU 80% 이상), 또는 시계열 **추세**(pending 파드 수 증가, API server inflight 상승 등).
- **Summary와의 관계:** Summary가 모두 정상이어도 Trend에서 경고가 나올 수 있음. 운영자가 용량·부하 대비를 미리 할 수 있게 함.

### 2.3 Top Offenders / Drill-down

- **역할:** “뭔가 나쁘다”고 Summary/Trend가 알려 주었을 때, **어떤 노드·워크로드**가 부담을 주는지 파악.
- **판정 방식:** ranking 기반(예: CPU 사용률 TOP10 노드, 재시작 수 TOP10 워크로드). 상세 목록은 필요 시 별도 뷰로 확장.
- **Summary와의 관계:** Summary·Trend에서 이상이 감지된 뒤, 이 레이어로 원인 후보를 좁힌다.

---

## 3. 건강 카테고리 (Health Categories)

Research 결과를 바탕으로, 모델은 다음 **다섯 가지 카테고리**를 사용한다. 각 카테고리는 “해당 영역이 건강한가?”를 판단하는 단위가 된다.

| 카테고리 | 담당 영역 | 건강의 의미 |
|----------|-----------|-------------|
| **Node / Infrastructure** | 노드 가용성, kubelet, 기본 리소스 | 노드가 Ready이고 스케줄 가능하며, 에이전트가 정상 동작함 |
| **Control Plane** | API server, scheduler, controller-manager, etcd | 제어면이 응답하고, 스케줄·조정이 정상 동작함 |
| **Workload / Pod** | 파드 상태, 재시작, Pending, eviction | 워크로드가 기동·스케줄되고, 과도한 재시작·eviction이 없음 |
| **Capacity / Resource pressure** | 클러스터 리소스 여유, 노드 사용률, Pending 추세 | 클러스터가 용량 한계에 있지 않고, 스케줄 여유가 있음 |
| **Network / Ingress** | 서비스·ingress 가용성, 엔드포인트 | 트래픽이 정상 전달되고, 백엔드가 비어 있지 않음 |

---

## 4. 카테고리별 대표 건강 신호

아래는 **카테고리별 대표 health signal**이다. 각 신호에 대해 **예시 Prometheus metric**, **운영적 의미**, **가능한 건강 상태**를 적고, **어느 레이어(Summary / Trend / Top Offenders)**에 두는지 명시한다.  
신호 수를 최소로 유지하기 위해, **Health Summary에는 카테고리당 최대 1개**만 배치한다.

### 4.1 Node / Infrastructure

| Health signal | Example Prometheus metric / PromQL | Operational meaning | Possible health condition | Layer |
|---------------|-------------------------------------|----------------------|---------------------------|--------|
| 노드 가용성(NotReady 개수) | `count(kube_node_status_condition{condition="Ready",status="false"} == 1)` | 스케줄 가능한 노드가 충분한가? | 0 = 정상, >0 = 노드 손실·점검 필요 | **Health Summary** |
| 노드 CPU/메모리 사용률 | node_exporter: `100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)` 등, 메모리: `node_memory_*` | 노드가 포화 직전인가? | &lt;80% = 여유, ≥80% = 리스크 | **Trend / Risk** |
| 노드별 리소스 사용 순위 | 위 metric을 노드별로 정렬, 상위 N개 | 어떤 노드가 가장 부하가 높은가? | TOP10 노드 목록 | **Top Offenders** |
| 디스크 가용 공간 | `node_filesystem_avail_bytes{mountpoint="/"}`, 또는 사용률 | 디스크 부족·I/O 리스크가 있는가? | 여유 있음 / 부족 임박 | **Trend / Risk** (선택) |

### 4.2 Control Plane

| Health signal | Example Prometheus metric / PromQL | Operational meaning | Possible health condition | Layer |
|---------------|-------------------------------------|----------------------|---------------------------|--------|
| API server 응답·에러 | `histogram_quantile(0.99, apiserver_request_duration_seconds_bucket)` (P99), `apiserver_request_total{code=~"5.."}` | 제어면이 지금 응답하는가? 에러는? | P99 &lt;1s, 5xx 극소 = 정상 | **Health Summary** |
| API server 부하(inflight) | `apiserver_current_inflight_requests` | API가 포화 직전인가? | 낮음 = 정상, 상승 추세 = 리스크 | **Trend / Risk** |
| Scheduler pending pods | `scheduler_pending_pods` (unschedulable 등) | 스케줄되지 못한 파드가 있는가? | 0 = 정상, >0 = 스케줄 불가 | **Health Summary** |
| Pending 수 추세 | `scheduler_pending_pods` 시계열 | 대기 파드가 늘어나는가? | 증가 = 용량·이슈 징후 | **Trend / Risk** |
| Control plane 구성요소별 상태 | 각 컴포넌트 `/healthz` 또는 해당 metric | scheduler, controller-manager 등 개별 건강 | 전부 정상 / 일부 비정상 | **Top Offenders** (드릴다운 시) |

### 4.3 Workload / Pod

| Health signal | Example Prometheus metric / PromQL | Operational meaning | Possible health condition | Layer |
|---------------|-------------------------------------|----------------------|---------------------------|--------|
| 장시간 Pending 파드 수 | `count(kube_pod_status_phase{phase="Pending"} == 1)` (필요 시 조건 추가) | 스케줄 대기 중인 파드가 있는가? | 0 = 정상, >0 = 스케줄·용량 이슈 | **Health Summary** |
| 과도한 재시작(예: 10분 내 N회 초과) | `increase(kube_pod_container_status_restarts_total[10m]) > threshold` 또는 CrashLoopBackOff 수 | 워크로드가 반복적으로 죽는가? | 임계치 이하 = 정상, 초과 = 문제 | **Health Summary** 또는 **Trend / Risk** |
| Eviction 발생 | kubelet eviction 관련 metric 또는 이벤트 | 노드가 리소스 부족으로 파드를 쫓아냈는가? | 없음 = 정상, 있음 = 용량·리소스 이슈 | **Health Summary** 또는 **Trend / Risk** |
| 워크로드별 재시작·에러 순위 | `kube_pod_container_status_restarts_total` by namespace, deployment 등 | 어떤 워크로드가 가장 문제인가? | TOP10 워크로드 목록 | **Top Offenders** |

### 4.4 Capacity / Resource pressure

| Health signal | Example Prometheus metric / PromQL | Operational meaning | Possible health condition | Layer |
|---------------|-------------------------------------|----------------------|---------------------------|--------|
| 클러스터 리소스 여유(또는 Pending과 연계) | 노드 allocatable vs request/usage 집계, 또는 `scheduler_pending_pods` | CPU/메모리 여유가 있는가? | 여유 있음 / 거의 없음 = 리스크 | **Trend / Risk** (Summary와 중복 피하기 위해 Trend로 두는 것을 권장) |
| 노드별 CPU/메모리 사용률 순위 | node_exporter 기반, 노드별 사용률 정렬 | 어떤 노드가 가장 포화에 가까운가? | TOP10 노드 목록 | **Top Offenders** |

*Capacity는 “지금 안전한가?”를 직접 나타내는 단일 신호보다, **Pending 수**(Control Plane·Workload에서 이미 다룸)와 **노드 사용률 추세·TOP10**으로 보는 것이 신호 수를 줄이면서도 의미를 유지하는 방식이다.*

### 4.5 Network / Ingress

| Health signal | Example Prometheus metric / PromQL | Operational meaning | Possible health condition | Layer |
|---------------|-------------------------------------|----------------------|---------------------------|--------|
| 서비스 엔드포인트 비어 있음 | `kube_endpoint_address_available == 0` (핵심 서비스만 필터) | 백엔드가 비어 있어 트래픽을 받지 못하는가? | 모두 있음 = 정상, 0개 = 서비스 불가 | **Health Summary** |
| Ingress/서비스 에러율·지연 | ingress controller 또는 서비스 메트릭(환경에 따라 다름) | 트래픽이 정상 전달되는가? | 에러율·지연 허용 범위 내 = 정상 | **Health Summary** (엔드포인트와 통합 가능) |
| 네트워크 지연·패킷 손실 | node/네트워크 metric (환경에 따라 다름) | 네트워크 병목·손실이 있는가? | 상승 시 리스크 | **Trend / Risk** (선택) |

---

## 5. Health Summary 신호 최소 집합 (요약)

운영자가 **“지금 클러스터가 안전한가?”**만 답하려면, 아래 **소수 신호만** 보면 된다. 수를 최소로 유지했다.

| # | Category | Health signal | 판정 기준(예시) |
|---|----------|---------------|-----------------|
| 1 | Node / Infrastructure | NotReady 노드 수 | 0 = 정상 |
| 2 | Control Plane | API server P99 지연·5xx | P99 &lt;1s, 5xx 없음(또는 극소) |
| 3 | Control Plane | Scheduler pending pods | 0 = 정상 |
| 4 | Workload / Pod | 장시간 Pending 파드 수 | 0(또는 허용 범위) = 정상 |
| 5 | Workload / Pod | 과도한 재시작 또는 Eviction | 없음 = 정상 |
| 6 | Network / Ingress | 핵심 서비스 엔드포인트 비어 있음 | 비어 있지 않음 = 정상 |

*실제 구현 시 Prometheus에 따라 “과도한 재시작”은 10분 내 N회 초과 등으로, “핵심 서비스”는 라벨/서비스 목록으로 정의한다. Network에서 Ingress 에러율을 Summary에 넣을 경우 6번과 통합하거나 7번째로 하나 더 둘 수 있다.*

---

## 6. 건강 판정 방식 (Health Judgement)

운영자가 신호를 어떻게 해석하는지 모델에 명시한다. 세 가지 방식이 레이어별로 대응한다.

### 6.1 Threshold 기반 (Summary, 일부 Trend)

- **적용:** Cluster Health Summary 대부분, Trend/Risk의 “임계치 근접” 판단.
- **방식:** 신호 값이 **기준값(threshold)**을 넘으면 비정상(또는 경고). 예: NotReady 노드 >0 → 비정상; API server P99 >1s → 비정상; 노드 CPU >80% → 리스크.
- **특징:** 단순하고, “지금 이 순간” 건강 여부를 빠르게 판단하기에 적합함.

### 6.2 Trend 기반 (Trend / Risk)

- **적용:** Trend / Risk Indicators.
- **방식:** 같은 metric의 **시계열 추세**로 판단. 예: `scheduler_pending_pods`가 지속적으로 증가 → “곧 스케줄 불가·용량 부족” 징후; `apiserver_current_inflight_requests` 상승 → API 포화 징후.
- **특징:** “지금은 괜찮지만 곧 나빠질 수 있다”는 조기 경고에 적합함. 구현 시 time window(예: 5m, 15m)와 증가율/절대 증가량 정의가 필요함.

### 6.3 Ranking 기반 (Top Offenders / Drill-down)

- **적용:** Top Offenders / Drill-down.
- **방식:** metric을 **노드별·워크로드별**로 집계한 뒤 **상위 N개(TOP10)**로 정렬. “어떤 노드가 CPU 사용률 1위인가?”, “어떤 deployment가 재시작 수 1위인가?” 등.
- **특징:** Summary·Trend에서 “뭔가 나쁘다”고 나온 뒤, **어디를 먼저 조사할지** 좁히는 데 사용. 목록 전체가 아니라 TOP10만 보여도 5–10분 내 판단에 충분하다는 전제.

---

## 7. Health Summary vs Trend vs Top Offender 구분 요약

| 레이어 | 신호 수 | 판정 방식 | 운영자가 묻는 질문 |
|--------|---------|------------|---------------------|
| **Health Summary** | **매우 적음 (5~7개)** | Threshold | 지금 클러스터가 안전한가? |
| **Trend / Risk** | 카테고리당 0~2개, 전체 소수 | Threshold 근접 + Trend | 곧 불안전해질 징후는? |
| **Top Offenders** | 카테고리·뷰당 1개(TOP10 목록) | Ranking | 문제 시 어디를 먼저 볼 것인가? |

Summary 신호는 **의도적으로 적게** 두어, 운영자가 **한 화면·한눈에** “전부 정상이면 안전”이라고 판단할 수 있게 한다. Trend와 Top Offenders는 “더 보려면 여기”로 이어지는 구조다.

---

## 8. 노이즈·비실행 가능 신호 처리 (Research 반영)

Research에서 정리한 **노이즈가 많거나 실행 가능하지 않은 신호**는 Health Summary에 넣지 않는다.

- **Readiness probe 실패**(Calico/CNI 등 인프라 probe 포함): 알림은 많을 수 있으나 실제 서비스 영향과 1:1이 아님 → Summary 제외, drill-down·참고용만.
- **일시적 파드 재시작(1~2회):** 배포·드레인 등으로 1~2회는 metric에 찍혀도 영향 없을 수 있음 → “과도한 재시작”을 **임계치**(예: 10분 내 N회 초과)로 정의해 Summary에 두고, 1~2회는 제외.
- **배경 재시도·cascade 알림:** Summary에는 “사용자 관점 실패”에 해당하는 신호만 두고, 재시도·cascade는 별도 처리.

이렇게 하면 **signal-to-noise**를 유지하면서 “지금 안전한가?”에만 집중할 수 있다.

---

## 9. 향후 Central Dashboard 설계에의 활용

이 Cluster Health Monitoring Model은 **향후 central dashboard**(운영자가 매일 5–10분 안에 클러스터 건강을 파악하는 대시보드) 설계의 기반이 된다.

### 9.1 Dashboard 구조 제안

1. **상단 또는 첫 화면: Cluster Health Summary**  
   - 5~7개 Summary 신호만 한 블록으로 표시.  
   - 각 신호: 현재 값 + 정상/비정상(색·아이콘).  
   - “전부 정상 = 클러스터 안전” 한 문장으로 해석 가능하게.

2. **중간: Trend / Risk Indicators**  
   - 카테고리별로 0~2개씩, “임계치 근접” 또는 “추세 경고”만 표시.  
   - Summary가 정상이어도 여기서 경고가 나오면 “용량·부하 대비 필요” 등으로 해석.

3. **하단 또는 별도 탭: Top Offenders / Drill-down**  
   - 노드 TOP10(CPU/메모리), 워크로드 TOP10(재시작·에러 등).  
   - Summary/Trend에서 이상이 있을 때 여기서 원인 후보를 찾음.

### 9.2 후속 프로젝트와의 연결

- **sre-monitoring-dashboard-design:** 이 모델의 Summary·Trend·Top Offenders를 Grafana 패널·레이아웃으로 구체화.
- **sre-monitoring-alert-policy:** Summary·Trend 신호 중 어떤 것을 알림으로 쓸지, 임계치·심각도 정의.
- **sre-monitoring-operational-runbooks:** “Summary에서 이 신호가 비정상일 때” 조사 순서·runbook 연결.

---

## 10. 설계 결정 요약

| 결정 | 내용 |
|------|------|
| **세 레이어** | Cluster Health Summary / Trend-Risk / Top Offenders — 운영 질문 3개에 1:1 대응. |
| **다섯 카테고리** | Node, Control Plane, Workload/Pod, Capacity, Network — Research와 동일. |
| **Summary 최소화** | 카테고리당 최대 1개, 전체 5~7개. “지금 안전한가?”만 답할 수 있는 수준. |
| **판정 방식** | Summary·일부 Trend = threshold; Trend = trend; Top Offenders = ranking(TOP10). |
| **노이즈 제외** | readiness probe 실패, 일시적 재시작 1~2회 등은 Summary에서 제외. |
| **Capacity** | 별도 Summary 1개보다 Pending·노드 사용률 추세·TOP10으로 처리해 신호 수 절약. |

---

*Architecture phase output. 다음 단계: Engineering — Core signal list 및 Prometheus/PromQL 예시 구체화.*
