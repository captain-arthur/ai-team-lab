# Research notes: Kubernetes Cluster Health Signal Model

**Project:** sre-monitoring-cluster-health-model  
**Phase:** Research  
**Source:** Manager brief, Kubernetes/Prometheus/Grafana/CNCF/SRE 참고 자료

---

## 1. 연구 목적과 초점

### 1.1 핵심 질문

**“운영자가 Kubernetes 클러스터 건강을 빠르게 판단할 수 있게 하는 신호는 무엇인가?”**

이 연구는 **metric을 많이 나열하는 것**이 아니라, 프로덕션 모니터링 시스템에서 **클러스터 건강 신호가 어떻게 모델링·구성되는지**를 이해하는 데 초점을 둔다.

### 1.2 참고한 방향

- Kubernetes 공식 observability·metrics 문서  
- Prometheus + Grafana 운영용 dashboard 구성 (Grafana Labs, Control Plane/Node/Workload 구분)  
- SRE 모니터링 관행(Google SRE Book golden signals, RED/USE)  
- CNCF observability 참고  
- 노이즈·alert fatigue 관련 사례(readiness probe, cascade 알림 등)

---

## 2. 프로덕션에서의 신호 구성 방식

### 2.1 계층·카테고리 구분

실무에서 Kubernetes 모니터링은 보통 **계층(layer)** 또는 **카테고리**로 나뉜다. 대표 구분:

| 카테고리 | 담당 영역 | “건강”의 의미 |
|----------|-----------|----------------|
| **Node / Infrastructure** | 노드 가용성, 리소스, kubelet | 노드가 스케줄 가능하고 리소스 여유가 있음 |
| **Control Plane** | API server, etcd, scheduler, controller-manager | 제어면이 응답하고 스케줄·조정이 정상 동작 |
| **Workload / Pod** | 파드 상태, 재시작, probe, throttling | 워크로드가 의도대로 기동·트래픽 처리 |
| **Capacity / Resource pressure** | CPU/메모리 사용률, eviction, pending | 클러스터가 용량 한계에 있지 않음 |
| **Network / Ingress** | 네트워크 지연, 패킷 손실, ingress/서비스 가용성 | 트래픽이 정상 전달됨 |

대량의 metric을 나열하기보다, **이 카테고리별로 대표 신호 1~2개**를 두어 “이 카테고리가 건강한가?”를 빠르게 보는 패턴이 공통된다.

### 2.2 Golden signals·RED·USE와의 대응

- **Golden signals (SRE):** Latency, Traffic, Errors, Saturation — 클러스터 전체를 보는 상위 신호로 쓰인다.  
- **RED (Rate, Errors, Duration):** 주로 **애플리케이션/서비스** 수준 SLI에 가깝다. 클러스터 건강 “요약”에는 API server latency·error rate 등으로 대표될 수 있다.  
- **USE (Utilization, Saturation, Errors):** **인프라/노드** 건강에 잘 맞는다. Node·Capacity 카테고리의 대표 신호(CPU/메모리 사용률, 디스크 I/O, eviction 등)와 직접 연결된다.

즉, “클러스터 건강을 빠르게 판단”하려면 **카테고리별로 Golden/USE/RED 중 하나에 대응하는 소수 대표 신호**를 두는 방식이 일반적이다.

---

## 3. 신호의 세 가지 역할 (모델 구분)

프로덕션 대시보드와 SRE 관행을 보면, 신호는 **용도**에 따라 다음 세 가지로 구분하는 것이 유효하다.

| 역할 | 목적 | 운영자가 묻는 질문 |
|------|------|---------------------|
| **Cluster Health Summary** | “지금 클러스터가 건강한가?”에 직접 답 | 지금 안전한가? |
| **Trend / Risk Indicators** | “곧 불건강해질 수 있다”는 조기 징후 | 불안전해질 징후는? |
| **Top Offenders / Drill-down** | “어디가 부담을 주는가?” 파악 | 문제 시 어디를 먼저 볼 것인가? |

- **Summary:** 카테고리당 1~2개, 전체를 한눈에 보는 상위 지표(예: API server 가용성, NotReady 노드 수, Pending 파드 수).  
- **Trend/Risk:** 같은 metric이라도 **임계치 근접·상승 추세**로 보는 것(예: 노드 CPU 80% 근접, pending 수 증가).  
- **Top Offenders:** **어떤 노드/워크로드가** 부하·에러·eviction을 만드는지 보는 뷰(TOP10 등). Summary가 “뭔가 나쁘다”고 하면, 여기서 원인 후보를 좁힌다.

---

## 4. 카테고리별 대표 신호

아래는 **카테고리별 대표 신호**와, 각 신호가 **현재 건강(Summary) / 조기 리스크(Trend) / 드릴다운(Top Offenders)** 중 어디에 쓰이기 적합한지, 그리고 **노이즈 가능성**을 정리한 것이다. 목표는 “긴 metric 목록”이 아니라 **신호 모델의 명확성**이다.

### 4.1 Node / Infrastructure health

| 신호 (개념) | 예시 Prometheus metric / PromQL | 답하는 운영 질문 | 역할 | 비고 |
|-------------|----------------------------------|------------------|------|------|
| 노드 가용성 | `kube_node_status_condition{condition="Ready",status="true"}` vs `false` | 지금 스케줄 가능한 노드가 충분한가? | **Health Summary** | NotReady 개수/비율이 0 또는 허용 범위여야 “건강” |
| 노드 리소스 여유 | `1 - (node_cpu_seconds_total / ignoring(mode) group_left() count without(mode,cpu)(...) )` 등, 또는 node_exporter 기반 CPU/메모리 사용률 | 노드가 포화 직전인가? | **Trend/Risk** | 80~90% 근접 시 “곧 불건강” 징후 |
| 노드별 리소스 사용 | 동일 metric을 노드별·상위 N개 | 어떤 노드가 가장 부하가 높은가? | **Top Offenders** | TOP10 노드 뷰 |
| 디스크 I/O·디스크 풀 | `node_filesystem_avail_bytes`, I/O 대기 | 디스크 부족·I/O 병목 가능성 | **Trend/Risk** | 고사용률·낮은 avail은 리스크 |
| kubelet 건강 | kubelet `/healthz` 또는 해당 metric | 노드 에이전트가 정상인가? | **Health Summary** (선택) | 요약에서 “한 개라도 실패하면 안 됨” 수준으로만 쓸 수 있음 |

### 4.2 Control plane health

| 신호 (개념) | 예시 Prometheus metric / PromQL | 답하는 운영 질문 | 역할 | 비고 |
|-------------|----------------------------------|------------------|------|------|
| API server 응답·가용성 | `apiserver_request_duration_seconds` (P99), `apiserver_request_total` (5xx 비율) | 제어면이 지금 응답하는가? 지연·에러는? | **Health Summary** | P99 &lt; 1s, 5xx 없음(또는 극소)가 “건강” |
| API server 부하 | `apiserver_current_inflight_requests` | API가 포화 직전인가? | **Trend/Risk** | 상승 추세·높은 값은 리스크 |
| Scheduler 대기 파드 | `scheduler_pending_pods` (unschedulable 등) | 스케줄되지 못한 파드가 있는가? | **Health Summary** | 0 또는 허용 범위 |
| Pending 파드 추세 | `scheduler_pending_pods` 시계열 | 대기 파드가 늘어나는가? | **Trend/Risk** | 증가 추세는 용량·이슈 징후 |
| etcd 저장소 크기 | `apiserver_storage_size_bytes` 등 | etcd가 비대해지고 있는가? | **Trend/Risk** | 8GB 등 제한 근접 시 리스크 |
| Control plane 구성요소별 상태 | 각 컴포넌트 `/healthz` 또는 해당 metric | scheduler, controller-manager 등 개별 건강 | **Health Summary** (요약 시 집계) 또는 **Drill-down** | 요약에서는 “전부 정상” 하나로 묶을 수 있음 |

### 4.3 Workload / Pod health

| 신호 (개념) | 예시 Prometheus metric / PromQL | 답하는 운영 질문 | 역할 | 비고 |
|-------------|----------------------------------|------------------|------|------|
| 크래시/재시작 빈도 | `kube_pod_container_status_restarts_total` 증가율, CrashLoopBackOff 수 | 워크로드가 반복적으로 죽는가? | **Health Summary** 또는 **Trend/Risk** | 10분 내 N회 초과 등은 “지금 문제” 또는 “악화 중” |
| Pending 파드 수 | `kube_pod_status_phase{phase="Pending"}` (장시간 Pending) | 스케줄되지 못한 파드가 있는가? | **Health Summary** | Control plane의 pending과 함께 봐도 됨 |
| CPU throttling | container CPU throttling metric (해당 exporter 있을 때) | 파드가 CPU 제한에 막혀 있는가? | **Trend/Risk** | 포화·성능 저하 징후 |
| Eviction 발생 | kubelet eviction 관련 metric 또는 이벤트 | 노드가 리소스 부족으로 파드를 쫓아내는가? | **Health Summary** 또는 **Trend/Risk** | 발생 시 “지금 불건강” 또는 “용량 리스크” |
| 네임스페이스/워크로드별 재시작·에러 | 위 metric을 namespace, deployment 등으로 그룹 | 어떤 워크로드가 가장 문제인가? | **Top Offenders** | TOP10 워크로드 뷰 |

### 4.4 Capacity / Resource pressure

| 신호 (개념) | 예시 Prometheus metric / PromQL | 답하는 운영 질문 | 역할 | 비고 |
|-------------|----------------------------------|------------------|------|------|
| 클러스터 리소스 여유 | 노드별 allocatable vs request/usage 집계 | CPU/메모리 여유가 있는가? | **Health Summary** 또는 **Trend/Risk** | 여유가 거의 없으면 “곧 불건강” |
| 노드 메모리/CPU 사용률 | node_exporter 등, 노드별 사용률 | 어떤 노드가 포화에 가까운가? | **Trend/Risk** + **Top Offenders** | 80~90%는 리스크, TOP10은 드릴다운 |
| Pending 지속 시간·수 | `scheduler_pending_pods`, Pending 파드 수 | 스케줄 대기가 길어지는가? | **Trend/Risk** | 용량 부족·스케줄 불가 징후 |

### 4.5 Network / Ingress health

| 신호 (개념) | 예시 Prometheus metric / PromQL | 답하는 운영 질문 | 역할 | 비고 |
|-------------|----------------------------------|------------------|------|------|
| 서비스/Ingress 가용성 | ingress controller 또는 서비스 엔드포인트 성공률·지연 | 트래픽이 정상 전달되는가? | **Health Summary** | 에러율·지연이 허용 범위 |
| 네트워크 지연·패킷 손실 | node/네트워크 metric (환경에 따라 다름) | 네트워크 병목·손실이 있는가? | **Trend/Risk** | 상승 시 리스크 |
| 엔드포인트 비어 있음 | `kube_endpoint_address_available` 등 | 서비스 백엔드가 비어 있는가? | **Health Summary** | 0개 available이면 “서비스 불가” |

---

## 5. 노이즈가 많거나 실행 가능하지 않은 신호 (Noisy / Non-actionable)

다음과 같은 신호는 **알림을 자주 유발하지만, 실제 서비스 영향과 1:1로 연결되지 않을 수 있어** 모델에서 “Summary가 아닌 참고용” 또는 “제외/별도 처리”로 두는 것이 좋다.

| 신호 예시 | 설명 | 권장 |
|-----------|------|------|
| **Readiness probe 실패** | 파드는 Running이지만 Service에서 제외됨. 재시작은 없어 기본 모니터링에 안 잡힐 수 있음. 일시적 실패·의존성 지연으로 알림만 많이 올 수 있음. | “실제 사용자 영향”과 연결된 신호(예: 5xx, latency)를 Summary로 두고, readiness 실패는 drill-down·조사용으로만 사용하거나, 팀별 whitelist 후 알림. |
| **Calico/CNI readiness 등 인프라 probe** | 네트워크 플러그인 probe 실패가 알림을 만들지만, 실제 트래픽 장애와 무관할 수 있음. | Health Summary에서 제외. 필요 시 별도 패널·낮은 심각도 알림. |
| **배경 재시도·retry** | 클라이언트 재시도, 일시적 오류로 인한 retry는 metric 상 에러처럼 보일 수 있음. | 에러율 정의 시 “최종 실패” vs “재시도” 구분. Summary에는 “사용자 관점 실패”만 넣는 것이 좋음. |
| **일시적 파드 재시작(1~2회)** | 배포·노드 드레인 등으로 재시작이 1~2회 있으면 metric에는 찍히지만, 사용자 영향이 없을 수 있음. | “재시작 빈도”를 10분 내 N회 초과 등 **임계치**로 두어, 일시적 1~2회는 Summary에서 제외하거나 Trend만 보는 방식. |
| **Cascade 알림** | DB 장애 등 하나의 원인으로 수백 개 서비스에서 알림이 쏟아짐. | 원인(DB) 쪽 Health Summary·Trend를 강화하고, 하위 서비스 알림은 집계·deduplication. |

정리하면, **Cluster Health Summary에는 “지금 클러스터/서비스가 안전한가?”에 직접 답하는 소수 신호만** 두고, readiness probe 실패·배경 재시도·일시적 재시작·cascade는 **별도 뷰·낮은 우선순위·또는 제외**하는 것이 signal-to-noise를 높이는 방법이다.

---

## 6. 요약 및 권장 사항

### 6.1 신호 모델 요약

- **카테고리:** Node/Infrastructure, Control Plane, Workload/Pod, Capacity/Resource pressure, Network/Ingress로 구분하는 것이 실무와 잘 맞는다.  
- **역할 구분:**  
  - **Cluster Health Summary:** 카테고리당 1~2개, “지금 건강한가?”에 직접 답.  
  - **Trend/Risk:** 같은 metric의 임계치 근접·상승 추세로 “곧 불건강” 징후.  
  - **Top Offenders/Drill-down:** “어디가 부담을 주는가?”를 보는 TOP10 스타일 뷰.  
- **노이즈:** readiness probe 실패, 배경 재시도, 일시적 재시작, cascade 알림은 Summary에서 제외하거나 별도·낮은 우선순위로 두는 것이 좋다.

### 6.2 Architecture 단계로의 전달

- 위 **카테고리 5개**와 **세 가지 역할(Summary / Trend-Risk / Top Offenders)**를 Cluster Health Monitoring Model의 뼈대로 사용할 수 있다.  
- 각 카테고리에서 **대표 신호 1~2개**만 골라 “최소 집합”을 만들면, 운영자가 5–10분 안에 “안전한가 / 곧 불안전해질 징후는 / 어디를 먼저 볼 것인가?”에 답하는 **central dashboard**의 기반이 된다.  
- Prometheus에서 **실제로 쿼리 가능한지**, retention·쿼리 복잡도·리소스 제한은 Engineering·Experiment 단계에서 검증하는 것이 좋다.

---

*Research phase output. 다음 단계: Architecture — Cluster Health Monitoring Model 설계.*
