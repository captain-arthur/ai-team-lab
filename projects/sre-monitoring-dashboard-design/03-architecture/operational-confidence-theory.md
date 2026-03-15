# Operational Confidence Theory — Kubernetes Cluster Safety

**Project:** sre-monitoring-dashboard-design  
**Phase:** Architecture (theoretical foundation)  
**Input:** Cluster Health Monitoring Model final-report, Central Kubernetes Operational Dashboard Design, implementation-ready-panel-spec

이 문서는 **Central Kubernetes Operational Dashboard** 가 의존하는 **운영 이론**을 정식화한다. 목표는 “어떤 metric을 쓸지”가 아니라, **운영자가 언제 “클러스터가 안전하다”고 확신할 수 있는지**, **어떤 신호가 서비스 영향의 진짜 전조인지**, **어떤 신호는 이미 손상 상태이거나 노이즈인지**를 명확히 하는 것이다. 따라서 **일반적인 observability 프레임워크(USE, RED 등)에 의존하지 않고**, Kubernetes 클러스터의 **실제 고장 전파와 운영 판단**에 기반한 이론을 서술한다.

---

## 1. 문제의 본질: Metric이 아니라 확신

### 1.1 “어떤 metric이 있는가?”가 아닌 질문

모니터링 설계에서 자주 나오는 질문은 “어떤 metric을 볼 것인가?”다. 그러나 **실제 운영 문제**는 다른 곳에 있다.

- **운영자가 확신 있게 “클러스터가 안전하다”고 말할 수 있는 조건이 무엇인가?**
- **어떤 신호가 “곧 불안전해질 것”을 미리 알려 주는가?**
- **어떤 신호는 서비스 영향의 진짜 전조(true precursor)이고, 어떤 신호는 이미 손상이 난 뒤의 결과이거나 노이즈인가?**

metric 목록만 나열하면 **가시성(visibility)** 은 늘어나지만 **운영 확신(operational confidence)** 은 자동으로 생기지 않는다. 확신은 **“어떤 조건이 만족될 때 안전하다고 볼 수 있는가”** 와 **“어떤 신호가 그 조건의 위반 또는 위협을 나타내는가”** 에 대한 명시적 정의에서 나온다.

### 1.2 이 문서의 역할

아래에서 우리는 다음을 정의한다.

1. **Operational Safety Conditions** — “클러스터가 현재 안전하다”고 말할 수 있는 **최소 조건** (metric 이름이 아니라 운영 조건).
2. **Failure Propagation Paths** — 클러스터 문제가 **서비스 영향**으로 이어지는 전형적 경로; 대시보드 시그널이 **전파 단계의 어디에** 해당하는지 매핑.
3. **시그널의 운영적 역할 분류** — 객체 타입(node/pod/network)이 아니라 **운영 판단에서의 역할**: safety condition / early warning / instability / damage / investigation.
4. **Leading vs lagging vs already-damaged vs noisy** — 각 중요 패널이 어느 그룹에 속하는지, 그 이유.
5. **Operational Confidence** — “지금 안전하다”와 “가까운 미래에도 안전할 가능성이 높다”에 대한 확신이 **어떤 시그널 조합**에서 나오는지.

이것은 대시보드 UI 설명이 아니라 **Kubernetes 클러스터 신뢰성에 대한 운영 설계 노트**로 읽을 수 있어야 한다.

---

## 2. Operational Safety Conditions (운영 안전 조건)

“클러스터가 **현재** 안전하다”고 말하려면, 아래 **최소 조건**이 모두 만족되어야 한다. 이 조건들은 metric 이름이 아니라 **시스템이 동작하기 위해 필요한 운영 상태**를 기술한다.

### 2.1 조건 정의

| # | 조건 (영문) | 조건 (설명) | 만족 시 의미 |
|---|-------------|-------------|--------------|
| **C1** | **Scheduling path is healthy** | 파드를 노드에 배치하는 경로(control plane: API server, scheduler, 그리고 노드의 kubelet·스케줄 가능 용량)가 정상 동작하며, “대기 중인 파드가 없거나 허용 가능한 수준”이다. | 새 워크로드가 스케줄될 수 있고, 기존 워크로드의 재스케줄이 가능하다. |
| **C2** | **Critical traffic path is healthy** | 사용자 트래픽이 들어오는 경로(ingress/controller)와, 트래픽이 도달하는 백엔드(Service → Endpoints)가 유효하다. 핵심 서비스의 endpoint가 비어 있지 않다. | 요청이 백엔드까지 도달할 수 있다. |
| **C3** | **Core workload stability is intact** | 워크로드가 과도한 재시작 폭증·연쇄 실패 상태가 아니다. “일시적 재시작 1~2회”가 아니라 **임계치를 넘는 불안정**이 없어야 한다. | 파드가 계속 실행 가능하고, 백엔드 용량이 유지된다. |
| **C4** | **No evidence of active resource collapse** | 노드가 Ready를 잃거나, 노드 리소스(CPU·메모리·디스크)가 이미 붕괴 단계(예: 대량 OOM, 디스크 풀)에 있지 않다. | 스케줄 가능한 노드와 리소스가 존재한다. |
| **C5** | **No evidence of cascading instability** | 한 영역의 실패가 다른 영역으로 연쇄적으로 퍼져 나가는 징후가 없다. 예: 노드 압박 → 대량 Pending → endpoint 감소 → 트래픽 실패 같은 **연쇄**가 진행 중이 아니다. | 문제가 한정되어 있거나, 아직 확대되기 전이다. |

### 2.2 조건과 Block 1 시그널의 대응

대시보드의 **Block 1 (Operational Confidence)** 패널은 위 조건의 **위반 여부**를 직접 또는 간접적으로 나타낸다.

| Safety condition | 위반 시 보이는 대표 시그널 (Block 1) |
|-------------------|----------------------------------------|
| C1 Scheduling path healthy | NotReady node count &gt;0 (노드 손실), Scheduler pending pods &gt;0, Workload Pending pod count &gt;0 (스케줄 불가·대기) |
| C2 Critical traffic path healthy | Critical service endpoint empty (백엔드 비어 있음) |
| C3 Core workload stability | Excessive restarts (재시작 폭증) |
| C4 No resource collapse | NotReady node count (노드 상실은 리소스 붕괴의 결과일 수 있음) |
| C5 No cascading instability | 위 여러 신호가 **동시에** 비정상이면 연쇄를 의심; 단일 신호만으로는 “연쇄”라 단정하지 않음 |

**API server health**는 C1의 **상단**을 이룬다. API server가 응답하지 않거나 5xx가 많으면 스케줄링 경로 전체가 막힌다. 따라서 C1이 만족되려면 (control plane metric이 있는 환경에서) API server가 정상이어야 한다.

이렇게 **운영 안전 조건**을 먼저 두고, “Block 1의 각 패널이 어떤 조건의 위반을 나타내는가?”를 매핑하면, **“전부 정상 = 위 조건 전부 만족 = 현재 안전”** 이라는 판단 규칙이 이론적으로 정당화된다.

---

## 3. Failure Propagation Paths (고장 전파 경로)

클러스터 문제가 **서비스 영향(지연·에러·다운)** 으로 이어지는 전형적인 **전파 경로**를 정의하면, 각 시그널이 **전파의 어느 단계**에 해당하는지 이해할 수 있다. 이는 “어떤 신호가 선행 지표이고 어떤 것이 이미 손상 상태인가?”를 구분하는 기초가 된다.

### 3.1 경로 1: 노드 압박 → 스케줄 지연 → Pending 증가 → Endpoint 부족 → Ingress/서비스 저하

| 단계 | 현상 | 이 단계에서 보이는 시그널 |
|------|------|---------------------------|
| 1. 노드 압박 | 노드 CPU·메모리·디스크 사용률 상승, 일부 노드에서 eviction·성능 저하 | **Node CPU utilization**, **Node memory / OOM risk**, **Node disk space** (Early Risk) |
| 2. 스케줄 지연 | 스케줄러가 적합한 노드를 찾기 어렵거나, API server 부하로 스케줄 지연 | **Scheduler pending pods** (control plane 있으면), **API server health** 지연 |
| 3. Pending 증가 | 스케줄되지 못한 파드 수 증가 | **Workload Pending pod count**, **Pending pods trend** |
| 4. Endpoint 부족 | 파드가 줄어들거나 기동 실패로 Service의 ready endpoint 수 감소 | **Critical service endpoint empty** |
| 5. Ingress/서비스 저하 | 트래픽이 백엔드로 가지 못하거나, 백엔드 과부하로 지연·5xx | **Ingress stress**, (조사) **Error-heavy ingress / services** |

**관찰:**  
- **Node CPU/memory/disk, Pending trend** 는 **1~3단계**에 해당한다. 서비스 영향(5단계) **이전**에 보이는 **선행·조기** 신호다.  
- **Workload Pending pod count, Critical endpoint empty** 는 이미 **3~4단계** — 스케줄 실패·백엔드 부족이 **현재 상태**로 나타난다. 즉 **이미 손상이 진행 중**일 수 있다.  
- **Ingress stress** 는 트래픽 경로에서의 **조기 압박**이거나, 이미 5단계에 들어선 **결과**일 수 있으므로, 맥락(다른 시그널과 함께)으로 해석해야 한다.

### 3.2 경로 2: 메모리 압박 → OOM / 재시작 폭증 → 백엔드 불안정 → 서비스 저하

| 단계 | 현상 | 이 단계에서 보이는 시그널 |
|------|------|---------------------------|
| 1. 메모리 압박 | 노드·컨테이너 메모리 사용률 상승, OOM 가능성 증가 | **Node memory / OOM risk** (Early Risk) |
| 2. OOM / 재시작 폭증 | 컨테이너 kill, 파드 재시작 급증 | **Excessive restarts** (Block 1), (조사) **Restart TOP10 pods** |
| 3. 백엔드 불안정 | replica 수 유지 어려움, endpoint 변동·감소 | **Critical service endpoint empty**, Pending (재스케줄 실패) |
| 4. 서비스 저하 | 지연·타임아웃·5xx | **Ingress stress**, 에러율 |

**관찰:**  
- **Node memory / OOM risk** 는 **1단계** — 서비스 영향 전 **선행 지표**.  
- **Excessive restarts** 는 **2단계** — 이미 OOM·불안정이 **발생한 뒤**의 결과이므로 **후행/이미 손상** 성격이 강하다. 다만 “재시작이 계속 늘어나는가?”는 **연쇄 불안정**의 진행을 보여 주므로, **불안정 신호(instability signal)** 로도 쓸 수 있다.  
- 따라서 **안전 확신**에는 “재시작이 임계치 이하인가?”로 **C3(workload stability)** 를 검사하고, **조기 리스크**는 **Node memory / OOM risk** 로 먼저 본다.

### 3.3 경로 3: Ingress/Controller 스트레스 → 요청 지연·에러 증가 → 서비스 영향

| 단계 | 현상 | 이 단계에서 보이는 시그널 |
|------|------|---------------------------|
| 1. Ingress/controller 부하 | 요청량·동시 연결·처리 지연 상승 | **Ingress stress** (Early Risk, Optional) |
| 2. 지연·에러 증가 | P99 지연 상승, 5xx 비율 증가 | **Ingress stress** 세부, (조사) **Error-heavy ingress / services** |
| 3. 서비스 영향 | 사용자 체감 지연·실패 | — (애플리케이션·SLO 메트릭) |

**관찰:**  
- **Ingress stress** 는 트래픽 경로에서의 **선행/조기** 신호일 수 있다. 다만 “이미 백엔드가 줄어든 결과”일 수도 있어, **Pending / endpoint** 와 함께 보면 선행인지 후행인지 구분에 도움이 된다.

### 3.4 경로 4: 제어면 장애 → 스케줄·조정 정지 → 전역 영향

| 단계 | 현상 | 이 단계에서 보이는 시그널 |
|------|------|---------------------------|
| 1. API server / scheduler 장애 | 지연·5xx, pending 처리 불가 | **API server health**, **Scheduler pending pods** |
| 2. 스케줄·조정 정지 | 새 파드 스케줄 불가, 기존 조정 지연 | **Workload Pending pod count** 증가, endpoint 변동 |
| 3. 전역 영향 | 클러스터 전체가 “멈춘 것처럼” 동작 | 여러 Block 1 신호 동시 비정상 |

**관찰:**  
- **API server health**, **Scheduler pending pods** 는 **1단계** — 제어면 문제의 **직접 표현**이자, 이후 전파의 **선행**이다.

---

## 4. 시그널의 운영적 역할 분류 (Operational Role of Signals)

시그널을 **객체 타입(node/pod/network)** 이 아니라 **운영 판단에서의 역할**로 분류한다.

### 4.1 역할 정의

| 역할 (영문) | 의미 | 사용 시점 |
|-------------|------|-----------|
| **Safety condition signal** | “현재 안전한가?”를 판단하는 **필수 조건**의 만족 여부. 위반 시 “지금 안전하지 않다”고 결론. | 매일 점검 시, Block 1에서 전부 정상인지 확인. |
| **Early warning signal** | 아직 C1~C5가 크게 깨지지 않았지만, **곧 깨질 가능성**이 있음을 나타냄. 전파 경로의 **앞단**에 해당. | Block 2 확인. 사전 대응(용량·정리·확장) 트리거. |
| **Instability signal** | **불안정이 진행 중**임을 나타냄. “한 번의 비정상”이 아니라 **추세·연쇄**를 보여 줄 수 있음. | Block 1의 Excessive restarts(임계치 초과), Block 2의 Pending trend 등. |
| **Damage signal** | **이미 손상이 난 상태**를 나타냄. 서비스 영향 직전 또는 직후. 선행 대응이 아니라 **현재 상태 인지·복구**에 사용. | Critical endpoint empty, Workload Pending &gt;0 (원인에 따라), NotReady &gt;0 등. |
| **Investigation signal** | “안전한가?”/“조기 징후인가?”를 **결정**하는 것이 아니라, **이상일 때 어디를 볼 것인가**를 알려 주는 뷰. | Block 3 TOP10 전부. |

### 4.2 패널별 역할 매핑

| Panel (Block) | 주된 운영적 역할 | 보조 역할 |
|---------------|------------------|-----------|
| NotReady node count | Safety condition (C4, C1 관련) | Damage (이미 노드 상실) |
| API server health | Safety condition (C1 상단) | — |
| Scheduler pending pods | Safety condition (C1) | Early warning (증가 추세일 때) |
| Workload Pending pod count | Safety condition (C1), Damage (스케줄 실패 이미 발생) | Instability (추세와 함께 보면) |
| Excessive restarts | Safety condition (C3), Instability | Damage (이미 재시작 폭증 발생) |
| Critical service endpoint empty | Safety condition (C2), Damage | — |
| Node CPU utilization | Early warning | — |
| Node memory / OOM risk | Early warning | — |
| Node disk space | Early warning | — |
| Pending pods trend | Early warning, Instability | — |
| CPU throttling risk | Early warning | — |
| Ingress stress | Early warning (트래픽 경로 압박) | Damage (이미 저하일 때) |
| CPU TOP10 / Memory TOP10 / Restart TOP10 / Pending by workload / Error-heavy ingress / Control plane | Investigation | — |

이 분류에 따르면, **Block 1** 은 대부분 **safety condition** 이며 일부는 **damage** 또는 **instability** 를 동시에 나타낸다. **Block 2** 는 **early warning** (및 일부 **instability**). **Block 3** 은 전부 **investigation** 이다.

---

## 5. Leading vs Lagging vs Already-Damaged vs Noisy Indicators

각 중요 패널이 **선행 지표(leading)**, **후행 지표(lagging)**, **이미 손상 상태 지표(already-damaged-state)**, **노이즈/비행동 지표(noisy/non-actionable)** 중 어디에 해당하는지, 그리고 **그 이유**를 정리한다.

### 5.1 정의

| 유형 | 의미 | 운영적 사용 |
|------|------|-------------|
| **Leading indicator** | **서비스 영향이나 심각한 손상 이전**에 나타나는 신호. 전파 경로의 **앞단**. 사전 대응이 가능함. | Early Risk 블록에서 강조. “곧 나빠질 수 있음” 판단. |
| **Lagging indicator** | **원인·사건이 이미 발생한 뒤**에 변하는 신호. 원인 파악·복구 확인에는 유용하나, “미리 알림”에는 부적합. | Safety condition으로 “이미 안전하지 않다” 결론 내리는 데 사용. |
| **Already-damaged-state indicator** | **손상이 이미 난 상태**를 나타냄. “예방”이 아니라 “현재 상태 인지·복구 우선순위”에 사용. | Block 1에서 “조사 필요” 트리거. |
| **Noisy / non-actionable indicator** | **행동으로 이어지기 어려운** 신호. 일시적이거나 임계치 해석이 팀마다 다르거나, 알림만 늘고 대응은 불명확한 경우. | Summary에서 제외하거나, 낮은 심각도·드릴다운만 사용. |

### 5.2 패널별 분류 및 이유

| Panel | 유형 | 이유 |
|-------|------|------|
| **NotReady node count** | **Already-damaged-state** (및 lagging) | 노드가 이미 Ready를 잃은 **이후**의 결과. “노드가 곧 NotReady가 될 것”을 미리 보지는 않음. 다만 “현재 스케줄 용량이 줄었다”는 **안전 조건 위반**을 나타내므로 safety condition으로 사용. |
| **API server health** | **Leading** (제어면 관점) | API server 지연·5xx는 **그 이후**의 스케줄 지연·클라이언트 실패의 **원인**. 전파 경로 4의 1단계. |
| **Scheduler pending pods** | **Lagging / already-damaged** | 스케줄러가 **이미** 처리하지 못하는 파드가 쌓인 상태. “곧 쌓일 것”이 아니라 “이미 쌓여 있음”. C1 위반의 **결과** 지표. |
| **Workload Pending pod count** | **Lagging / already-damaged** | 스케줄·리소스 문제로 **이미** Pending인 파드 수. 전파 경로 1의 3단계. 선행이 아님. |
| **Excessive restarts** | **Lagging / already-damaged**, 일부 **instability** | 재시작은 OOM·실패 **이후**에 올라감. “곧 재시작이 폭증할 것”을 직접 보지는 않음. 다만 **계속 증가**하면 불안정 진행(instability)으로 해석 가능. |
| **Critical service endpoint empty** | **Already-damaged-state** | 백엔드가 **이미** 비어 있는 상태. 트래픽 실패 직전 또는 직후. 선행이 아님. |
| **Node CPU utilization** | **Leading** | 노드 CPU가 높아지는 것은 **그 다음**의 eviction·스케줄 실패·성능 저하 **이전**에 나타남. 전파 경로 1의 1단계. |
| **Node memory / OOM risk** | **Leading** | 메모리 압박은 OOM·재시작 폭증 **이전**에 나타남. 전파 경로 2의 1단계. |
| **Node disk space** | **Leading** | 디스크 여유 부족은 이미지 풀·로그 실패·노드 불안정 **이전**에 나타남. |
| **Pending pods trend** | **Leading** (추세로 볼 때) / **Instability** | “Pending 수가 **늘어나고 있다**”는 것은 스케줄·용량 문제가 **진행 중**이라는 신호. 단순 현재값만 보면 lagging에 가깝고, **추세**로 보면 선행·불안정 신호. |
| **CPU throttling risk** | **Leading** | CPU limit에 의한 throttling은 **그 다음**의 지연·불만 **이전**에 나타남. |
| **Ingress stress** | **Leading** (트래픽 경로 압박 시) / **Already-damaged** (이미 저하 시) | 부하·지연이 **증가하는 단계**에서는 선행. 이미 5xx·타임아웃이 많다면 이미 손상 상태. |
| **Block 3 (TOP10, Pending by workload 등)** | **Investigation** (leading/lagging 아님) | “안전한가?”/“조기인가?”를 **결정**하는 지표가 아니라, **이상일 때 원인 후보를 좁히는** 뷰. 선행/후행 구분은 해당 뷰가 기반하는 metric에 따름. |

### 5.3 요약

- **Leading:** Node CPU, Node memory/OOM risk, Node disk, Pending trend(추세), CPU throttling risk, (맥락에 따라) Ingress stress, API server health. → **Block 2 중심** + API server.  
- **Lagging / already-damaged:** NotReady, Scheduler pending, Workload Pending, Excessive restarts, Critical endpoint empty. → **Block 1** 의 대부분. “현재 안전하지 않다”는 **결과**를 보여 주며, 선행 대응보다는 **조사·복구** 트리거.  
- **Noisy:** 설계에서 Summary에 넣지 않은 신호 — 예: readiness probe 일시 실패, 재시작 1~2회만으로 critical로 쓰는 경우. 이런 것은 **노이즈**가 되기 쉽고, 팀에서 임계치·해석을 정하지 않으면 **non-actionable** 이다.

---

## 6. Operational Confidence (운영 확신)의 정의

**Operational Confidence** 는 “지금 클러스터가 안전하다”와 “가까운 미래에도 안전할 가능성이 높다”에 대해 **운영자가 확신을 가질 수 있는 근거**를 말한다. 이 근거는 **특정 시그널의 조합**과 **해석 규칙**으로 구체화된다.

### 6.1 “지금 안전하다”에 대한 확신 (Current safety confidence)

**정의:**  
운영자가 “클러스터가 **현재** 안전하다”고 말할 수 있는 것은, **Operational Safety Conditions C1~C5**가 모두 만족된다고 **판단할 수 있을 때**이다.

**판단 규칙:**  
- **Block 1 (Operational Confidence)** 의 **모든** 패널이 **정상(threshold 이내)** 이면, 다음을 가정할 수 있다.  
  - C1: NotReady=0, (있으면) API server·Scheduler pending 정상, Workload Pending=0 → 스케줄 경로 정상.  
  - C2: Critical endpoint empty가 “비어 있지 않음” → 트래픽 경로 유효.  
  - C3: Excessive restarts가 임계치 이하 → 워크로드 불안정이 심하지 않음.  
  - C4: NotReady=0 → 노드 붕괴·대량 상실 없음.  
  - C5: 한 개의 비정상도 없으므로 “연쇄가 진행 중”이라 볼 만한 다중 위반 없음.  
- 따라서 **“Block 1 전부 정상 = C1~C5 만족 = 현재 안전”** 이라는 **단일 규칙**이 이론적으로 정당화된다.  
- **확신의 정도:** Block 1에 control plane(API server, Scheduler pending)이 포함되면 C1을 더 직접적으로 검사하므로 **확신이 높다**. Managed K8s처럼 해당 패널이 없으면, C1의 “스케줄 경로”는 **Workload Pending=0** 등으로 간접 추론하며, 제어면 자체는 별도 채널로 보완해야 한다.

### 6.2 “가까운 미래에도 안전할 가능성이 높다”에 대한 확신 (Future risk confidence)

**정의:**  
“곧 불안전해질 수 있다”는 **조기 리스크**가 없을 때, 운영자는 “가까운 미래에도 현재 수준의 안전이 유지될 가능성이 높다”고 볼 수 있다.

**판단 규칙:**  
- **Block 2 (Early Risk)** 의 패널에서 **경고/위험(예: 80% 초과, 디스크 여유 &lt;10%, Pending 추세 상승)** 이 **없으면**, 다음을 가정할 수 있다.  
  - 노드 압박(CPU·메모리·디스크)이 심하지 않음 → 전파 경로 1·2의 **1단계**가 활성화되지 않음.  
  - Pending이 급증하는 추세가 아님 → 스케줄·용량이 **당장은** 악화되지 않는 방향.  
  - (있으면) CPU throttling·Ingress stress가 임계치 이내 → 트래픽·워크로드 경로에서 **당장** 위협이 크지 않음.  
- 따라서 **“Block 2 전부 경고/위험 없음 = 조기 리스크 없음 = 가까운 미래 안전 가능성 높음”** 이라는 해석이 가능하다.  
- **확신의 정도:** Block 2가 **leading indicator** 위주로 구성되어 있으므로, “아직 손상이 나기 전”에 리스크를 인지할 수 있어, **사전 대응**을 한 뒤 “미래 안전”을 유지할 수 있다.

### 6.3 두 가지 확신의 조합

- **현재 안전 확신:** Block 1 전부 정상.  
- **미래 안전 가능성 확신:** Block 2 전부 경고/위험 없음.  

**둘 다 만족** → “클러스터는 지금 안전하고, 가까운 미래에도 그렇게 유지될 가능성이 높다”는 **운영 확신**을 갖는 것이 타당하다.  
**Block 1 비정상** → “현재 안전하지 않다” → 조사·복구.  
**Block 1 정상, Block 2 경고** → “지금은 괜찮지만 곧 나빠질 수 있다” → 사전 대응(용량·정리·확장) 또는 Block 3으로 원인 후보 좁히기.

이렇게 **안전 조건(C1~C5)** 과 **전파 경로**와 **시그널 역할·선행/후행**을 명시해 두면, 대시보드의 “Block 1 = 안전, Block 2 = 조기 리스크” 구조가 **단순 UI 규칙**이 아니라 **운영 이론에 기반한 설계**가 된다.

---

## 7. 요약 및 설계로의 연결

- **Operational Safety Conditions:** C1 스케줄 경로, C2 트래픽 경로, C3 워크로드 안정성, C4 리소스 붕괴 없음, C5 연쇄 불안정 없음. **Block 1** 은 이 조건들의 위반 여부를 보여 준다.  
- **Failure Propagation Paths:** 노드 압박 → 스케줄 → Pending → endpoint → 서비스; 메모리 압박 → OOM/재시작 → 백엔드 불안정 → 서비스; Ingress 스트레스 → 지연/에러 → 서비스; 제어면 장애 → 전역. 시그널을 **전파 단계**에 매핑하면 **선행 vs 이미 손상** 구분이 명확해진다.  
- **시그널 역할:** safety condition / early warning / instability / damage / investigation. 객체 타입이 아니라 **운영 판단 역할**로 분류.  
- **Leading / lagging / already-damaged / noisy:** Block 2 중심의 Node CPU·memory·disk·Pending trend·(선택) CPU throttling·Ingress stress는 **leading**. Block 1의 NotReady·Pending·Restarts·Endpoint 등은 **lagging 또는 already-damaged**. Block 3은 **investigation**.  
- **Operational Confidence:** “지금 안전” = Block 1 전부 정상 = C1~C5 만족. “미래 안전 가능성” = Block 2 경고/위험 없음 = 조기 리스크 없음.  

이 문서는 **Central Kubernetes Operational Dashboard** 와 **implementation-ready-panel-spec** 의 이론적 기반이다. 대시보드 레이아웃과 패널 목록은 이 **운영 확신 이론**에 맞춰 설계되었으며, 새로운 시그널을 추가할 때도 **어느 조건·전파 단계·역할·선행/후행**에 해당하는지 먼저 정리하면, 일관된 운영 설계를 유지할 수 있다.

---

*Operational Confidence Theory v1. Kubernetes cluster safety에 대한 운영 설계 노트. SRE Monitoring program — sre-monitoring-dashboard-design project.*
