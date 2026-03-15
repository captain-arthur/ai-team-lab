# Central Kubernetes Operational Dashboard Design

**문서 제목:** Central Kubernetes Operational Dashboard Design  
**Project:** sre-monitoring-dashboard-design  
**Program:** SRE Monitoring  
**Version:** v1 (최종)

이 문서는 **Central Kubernetes Operational Dashboard** 의 최종 설계를 정리한 것이다. 아키텍처(03-architecture), **운영 확신 이론(03-architecture/operational-confidence-theory.md)**, 패널 설계(04-engineering), 구현 명세(implementation-ready-panel-spec), 리뷰(06-review) 및 Kubernetes Cluster Health Monitoring Model(final-report)을 통합하여, 구현·운영·후속 프로젝트(alert 정책, runbook)에서 참조할 수 있는 단일 보고서로 구성했다. 대시보드는 **사용성**뿐 아니라 **운영 논리(operational logic)** 에 의해 정당화된다.

---

## 1. Operational Problem (운영 문제)

### 1.1 기존 대시보드가 운영 확신을 주지 못하는 이유

Prometheus와 Grafana로 수많은 metric을 수집하고 대시보드가 이미 있어도, 운영자는 **“이게 정상인가?”, “클러스터가 실제로 안전한가?”, “누가 확인해 줄 수 있나?”** 를 자주 묻는다. 이는 **metric은 많지만 명확한 운영 확신(operational confidence) 모델이 부족**하기 때문이다. 패널이 많을수록 “어디가 정상이고 어디가 비정상인지” 한눈에 보이지 않고, **결론을 내리는 데 시간이 걸리거나** 결국 **경험·스크립트·다른 채널**에 의존하게 된다. 즉, **많은 metric을 나열하는 것만으로는 일상적인 운영 확신이 생기지 않는다.**

### 1.2 흔한 운영 상황

- **“이게 정상인가?”** — 수십 개 패널 중 어떤 값이 허용 범위인지, 지금 클러스터가 “괜찮은” 상태인지 판단하기 어렵다.  
- **“누가 확인해 줄 수 있나?”** — 공통된 “안전” 정의가 없어, 팀원마다 다른 기준으로 말할 수 있고, 책임 소재가 불명확해진다.

### 1.3 목표: 명확한 운영 확신 제공

이 대시보드의 목표는 **“또 하나의 큰 대시보드”**가 아니라, **운영자가 두 가지 질문에 확신 있게 답할 수 있는 단일 진입점**을 제공하는 것이다.

1. **지금 클러스터가 안전하다고 확신할 수 있는가?**  
2. **클러스터가 곧 불안전해질 조기 징후가 보이는가?**

또한 이상이 보일 때 **어디를 조사할지** 빠르게 찾을 수 있게 한다. **운영 확신**은 metric 양이 아니라 **“지금 안전한가?” / “곧 위험한가?”에 답할 수 있는지**로 측정한다.

---

## 2. Key Questions (핵심 질문)

대시보드는 다음 **세 가지 질문**에 답할 수 있어야 한다.

| # | 질문 | 대시보드에서 답하는 방식 |
|---|------|--------------------------|
| 1 | **지금 클러스터가 안전한가?** | Block 1 (Operational Confidence) — 소수 신호가 전부 정상이면 “안전”, 하나라도 비정상이면 “조사 필요”. |
| 2 | **곧 불안전해질 조기 징후가 보이는가?** | Block 2 (Early Risk) — 노드 압박, OOM 위험, 디스크, pending 추세, (가능 시) CPU throttling·ingress stress 등. |
| 3 | **이상이 보일 때 어디를 조사할 것인가?** | Block 3 (Investigation / Top Offenders) — CPU/메모리 TOP10 노드, 재시작 TOP10, Pending by workload 등으로 원인 후보 좁히기. |

---

## 3. Dashboard Model (대시보드 모델)

대시보드는 **Kubernetes Cluster Health Monitoring Model** 의 세 레이어를 그대로 화면 구조로 옮긴 **3계층 모델**을 따른다.

### 3.1 세 계층

| 계층 | 답하는 질문 | 배치 |
|------|-------------|------|
| **Operational Confidence** | 지금 클러스터가 안전한가? | 메인 뷰 최상단 (Block 1). |
| **Early Risk** | 곧 불안전해질 징후가 있는가? | 메인 뷰 Block 1 직하 (Block 2). |
| **Investigation / Top Offenders** | 어디를 먼저 조사할 것인가? | 드릴다운(탭·접기). Block 3. |

### 3.2 빠른 운영 판단을 위한 구조

- **Operational Confidence:** 4~6개 **stat** 패널만 보고 **“전부 정상 = 안전, 하나라도 비정상 = 조사 필요”** 한 규칙으로 결론. 숫자 해석 없이 **상태(색·아이콘)** 만 보면 된다.  
- **Early Risk:** 같은 메인 뷰 안에서 **“곧 나빠질 수 있다”**는 신호만 한 블록으로 묶어, Summary가 정상이어도 **용량·정리·확장** 등 사전 대응을 트리거할 수 있게 한다.  
- **Investigation:** 이상이 있을 때만 **해당 카테고리의 TOP10**으로 들어가, **어느 노드·워크로드**가 원인 후보인지 5–10분 안에 좁힌다.

이 구조로 **5–10분 일상 점검**, **높은 운영 확신**, **조기 리스크 가시성**, **낮은 대시보드 과부하**를 동시에 만족하도록 설계했다.

---

## 4. Operational Logic and Panel Justification (운영 논리와 패널 정당화)

대시보드에 나오는 패널은 **“있는 metric 나열”**이 아니라 **운영 안전 조건·고장 전파 경로**에 기반한 **의미 있는 시그널**이다. 이 절에서는 운영 확신 이론(operational-confidence-theory.md)을 요약하고, **왜 이 패널들이 최소이면서 충분한가**를 설명한다.

### 4.1 Operational Safety Conditions (운영 안전 조건)

“클러스터가 **현재** 안전하다”고 말할 수 있는 **최소 조건**은 다음 다섯 가지다. Block 1은 이 조건의 **위반 여부**를 보여 준다.

| 조건 | 설명 | Block 1에서 위반 시 보이는 시그널 |
|------|------|-----------------------------------|
| **C1. Scheduling path is healthy** | 파드를 노드에 배치하는 경로(API server, scheduler, 노드 용량)가 정상이며, 대기 파드가 없거나 허용 수준이다. | NotReady &gt;0, Scheduler pending &gt;0, Workload Pending &gt;0 |
| **C2. Critical traffic path is healthy** | 트래픽이 도달하는 백엔드(Service → Endpoints)가 유효하다. | Critical service endpoint empty |
| **C3. Core workload stability is intact** | 과도한 재시작 폭증·연쇄 실패가 없다. | Excessive restarts (임계치 초과) |
| **C4. No evidence of active resource collapse** | 노드가 Ready를 잃거나 리소스 붕괴 단계에 있지 않다. | NotReady node count &gt;0 |
| **C5. No evidence of cascading instability** | 한 영역 실패가 다른 영역으로 연쇄 퍼지는 징후가 없다. | Block 1 여러 신호 **동시** 비정상 시 연쇄 의심 |

**API server health**는 C1의 **상단**이다. API server 비정상이면 스케줄링 경로 전체가 막힌다.  
따라서 **“Block 1 전부 정상 = C1~C5 만족 = 현재 안전”** 이라는 판단 규칙은 **운영 이론**에 의해 정당화된다.

### 4.2 Failure Propagation Paths (고장 전파 경로)

클러스터 문제가 **서비스 영향**으로 이어지는 전형 경로와, 각 시그널이 **어느 단계**에 해당하는지는 다음과 같다.

| 경로 | 전파 단계 (요약) | 해당 시그널 (단계) |
|------|------------------|---------------------|
| **경로 1** | 노드 압박 → 스케줄 지연 → Pending 증가 → Endpoint 부족 → Ingress/서비스 저하 | Node CPU/memory/disk, Pending trend (1~3단계, **선행**); Workload Pending, Critical endpoint (3~4단계, **이미 손상**); Ingress stress (5단계) |
| **경로 2** | 메모리 압박 → OOM/재시작 폭증 → 백엔드 불안정 → 서비스 저하 | Node memory / OOM risk (1단계, **선행**); Excessive restarts (2단계, **후행/손상**); Critical endpoint (3단계) |
| **경로 3** | Ingress/controller 스트레스 → 지연·에러 증가 → 서비스 영향 | Ingress stress (1~2단계, **선행 또는 결과**) |
| **경로 4** | 제어면 장애 → 스케줄·조정 정지 → 전역 영향 | API server health, Scheduler pending (1단계); Workload Pending 증가 (2단계) |

**선행 지표(leading)** 는 서비스 영향 **이전**에 보이므로 Block 2에서 조기 리스크로 강조한다. **후행/이미 손상(lagging, already-damaged)** 는 Block 1에서 “현재 안전하지 않다”는 **결과**로 사용한다.

### 4.3 Why the Selected Panels Are Not Arbitrary (선택된 패널이 임의가 아닌 이유)

- **Block 1** 의 각 패널은 **특정 Operational Safety Condition의 위반**을 나타낸다. 조건(C1~C5)을 먼저 정의했기 때문에, “어떤 metric을 넣을지”가 아니라 **“어떤 조건을 검사할지”** 에서 패널이 유도된다.  
- **Block 2** 의 각 패널은 **고장 전파 경로의 앞단(1~3단계)** 에 해당하는 **선행 지표**다. “노드 CPU가 높다”는 것은 **그 다음**에 올 eviction·스케줄 실패·서비스 저하 **이전**에 보이는 신호이므로, **조기 리스크**로 두는 것이 운영적으로 타당하다.  
- **Block 3** 은 “안전한가?”/“조기인가?”를 **결정**하는 패널이 아니라, **이상일 때 원인 후보를 좁히는** investigation 뷰다. Block 1·2의 **어떤 신호가 비정상인지**에 따라 **어느 TOP10**을 볼지가 정해진다.

따라서 패널 선택은 **운영 안전 조건 + 전파 경로 + 선행/후행 구분**에 의해 정당화되며, 단순 metric 목록이 아니다.

### 4.4 A-Priority Panel Mapping (A 우선순위 패널 매핑)

구현 명세(implementation-ready-panel-spec)의 **A (Required for v1)** 패널마다, **지원하는 안전 조건**, **감지하는 고장 전파 단계**, **지표 유형(leading / lagging / already-damaged / investigation)** 을 명시한다.

#### Block 1 (Operational Confidence)

| Panel name | Supports safety condition | Failure propagation stage | Indicator type |
|------------|---------------------------|----------------------------|----------------|
| **NotReady node count** | C1, C4 | 경로 1·4 결과(노드 상실) | already-damaged (및 lagging) |
| **Workload Pending pod count** | C1 | 경로 1 단계 3, 경로 4 단계 2 | lagging / already-damaged |
| **Excessive restarts** | C3 | 경로 2 단계 2 (OOM/재시작 폭증) | lagging / already-damaged |
| **Critical service endpoint empty** | C2 | 경로 1 단계 4, 경로 2 단계 3 | already-damaged |

#### Block 2 (Early Risk)

| Panel name | Supports safety condition | Failure propagation stage | Indicator type |
|------------|---------------------------|----------------------------|----------------|
| **Node CPU utilization** | C4 위협 조기 감지 | 경로 1 단계 1 (노드 압박) | leading |
| **Node memory / OOM risk** | C4 위협 조기 감지 | 경로 1 단계 1, 경로 2 단계 1 | leading |
| **Node disk space** | C4 위협 조기 감지 | 경로 1 단계 1 (디스크 부족) | leading |
| **Pending pods trend** | C1 위협 조기 감지 | 경로 1 단계 2~3 (추세) | leading / instability |

#### Block 3 (Investigation / Top Offenders)

| Panel name | Supports safety condition | Failure propagation stage | Indicator type |
|------------|---------------------------|----------------------------|----------------|
| **CPU TOP10 nodes** | — | 경로 1 단계 1 드릴다운 (어느 노드가 압박인지) | investigation |
| **Memory TOP10 nodes** | — | 경로 2 단계 1 드릴다운 | investigation |
| **Restart TOP10 pods** | — | 경로 2 단계 2 드릴다운 | investigation |
| **Pending pods by workload** | — | 경로 1 단계 3 드릴다운 | investigation |

### 4.5 Why These Panels? (왜 이 패널들인가?)

- **최소한으로 필요한 이유:**  
  “지금 안전한가?”에 답하려면 **C1~C5**를 검사해야 하고, 각 조건당 **최소 하나의 시그널**이 필요하다. Block 1 필수 4개(NotReady, Workload Pending, Excessive restarts, Critical endpoint)는 **C1~C4**를 커버하며, control plane metric이 있으면 API server·Scheduler pending으로 C1을 더 직접 검사한다. “곧 불안전해질 징후”에 답하려면 **전파 경로의 앞단**을 봐야 하므로, Node CPU·memory·disk·Pending trend 네 가지 **선행 지표**가 최소한 필요하다.  
- **v1에서 충분한 이유:**  
  위 **8개(Block 1 필수 4 + Block 2 필수 4)** 만으로 **현재 안전 확신**(Block 1 전부 정상 = C1~C5 만족)과 **미래 안전 가능성 확신**(Block 2 경고 없음 = 조기 리스크 없음)에 **결정적으로** 답할 수 있다. Block 3은 “이상일 때 어디를 볼 것인가?”만 담당하므로, 메인 뷰에 넣지 않아도 **확신**에는 영향을 주지 않는다.  
- **다른 시그널을 메인 뷰에서 제외하는 이유:**  
  - **노이즈·비행동 지표:** readiness probe 일시 실패, 재시작 1~2회만으로 critical로 쓰는 식의 신호는 **행동으로 이어지기 어렵고** 알림 피로만 증가시키므로 Summary(Block 1)에 넣지 않는다.  
  - **상세 목록(TOP10 등):** “안전한가?”를 **결정**하는 것이 아니라 **조사 대상**을 좁히는 용도이므로 Block 3(드릴다운)으로 분리해, 메인 뷰 과부하를 피한다.  
  - **환경 의존적 시그널(API server, Scheduler, CPU throttling, Ingress stress):** control plane·ingress metric이 없는 환경에서는 **해당 패널을 숨기거나 제외**해도, 나머지 필수 8개만으로 **최소한의 운영 확신**은 유지된다. Optional(B)로 두어 “있으면 보강, 없으면 최소 집합”으로 운영한다.

---

## 5. Dashboard Layout (대시보드 레이아웃)

### 5.1 세 블록

| Block | 이름 | 내용 | 패널 수 |
|-------|------|------|---------|
| **Block 1** | Operational Confidence | NotReady node count, API server health, Scheduler pending pods, Workload Pending pod count, Excessive restarts, Critical service endpoint empty. (환경별 4~6개) | 4~6개 |
| **Block 2** | Early Risk | Node CPU utilization, Node memory / OOM risk, Node disk space, Pending pods trend, (선택) CPU throttling risk, Ingress stress. | 4~6개 |
| **Block 3** | Investigation / Top Offenders | CPU TOP10 nodes, Memory TOP10 nodes, Restart TOP10 pods, Pending pods by workload, (선택) Error-heavy ingress, Control plane component health. | 4~6개 |

### 5.2 메인 뷰 vs 드릴다운

- **메인 뷰 (Primary view):** **Block 1 + Block 2만** 포함. **총 10~14개 이하**(표준 8~12개) 패널. 운영자가 **매일 5–10분** 점검할 때 **반드시 보는 영역**. 스크롤 최소(한 화면 또는 한 번).  
- **드릴다운 (Drill-down):** **Block 3** 전부. **메인 뷰에 포함하지 않음.** 탭·접기·아래쪽 별도 섹션으로 **“이상 시”** 또는 **“더 보기”**로만 진입. Block 3은 **기본 접기(collapsed) 또는 별도 탭**으로 두어, 첫 로딩 시 Block 1·2만 보이게 하는 것을 권장한다.

### 5.3 행 구성 (Grafana)

```
Row 1: [제목] Operational Confidence
  → NotReady | API health* | Scheduler pending* | Pending pods | Restarts | Endpoint
  * control plane 미사용 시 숨김

Row 2: [제목] Early Risk
  → Node CPU | Node memory/OOM | Node disk | Pending trend | [CPU throttling]* | [Ingress stress]*
  * 선택

Row 3 (접기 또는 별도 탭): [제목] Investigation / Top Offenders
  → CPU TOP10 | Memory TOP10 | Restart TOP10 | Pending by workload | [Error-heavy ingress]* | [Control plane]*
  * 선택. 기본 숨김
```

---

## 6. Operational Workflow (운영 워크플로)

### 6.1 일상 5–10분 점검

1. **대시보드 진입** → **Block 1 (Operational Confidence)** 만 먼저 본다.  
2. **Block 1 전부 정상인가?**  
   - **예** → Block 2 (Early Risk)를 훑는다. 경고 없으면 “클러스터 안전” 결론. 경고 있으면 “곧 불안정해질 수 있음” 인지 후 용량·정리 검토 또는 Block 3 해당 뷰로 진입.  
   - **아니오** → “조사 필요”. **어느 신호가 비정상인지** 확인 → **카테고리**(Node / Workload / Network 등) 결정 → **Block 3** 해당 Investigation 뷰(TOP10)로 진입 → TOP10에서 1~2개 후보 선택 → kubectl·로그·이벤트로 조사.  
3. Block 2만 경고(Block 1은 전부 정상)인 경우: **유형**(노드 압박 / 디스크 / pending 추세 등) 확인 → 필요 시 Block 3 해당 TOP10으로 **어느 노드·워크로드**가 부담을 주는지 확인 → 사전 대응(확장·정리 등).

### 6.2 이상 신호 시 조사 흐름

```
Block 1에서 비정상 발견
    → 어느 신호인지 확인 (NotReady / Pending / Restarts / Endpoint 등)
    → 카테고리 결정 (Node / Workload / Network)
    → Block 3 해당 뷰 열기 (CPU TOP10 / Memory TOP10 / Restart TOP10 / Pending by workload 등)
    → TOP10에서 1~2개 선택 후 kubectl describe, logs, events 로 원인 파악
```

### 6.3 조기 리스크 감지

- Block 1이 전부 정상이어도 **Block 2** 에서 Node CPU·메모리 80% 근접, 디스크 여유 부족, Pending 증가, (가능 시) CPU throttling·Ingress stress 등이 보이면 **“곧 불안전해질 수 있다”** 고 인지한다.  
- 이때 Block 3의 **Node CPU TOP10 / Memory TOP10** 등으로 **어느 노드가 부담을 주는지** 먼저 확인한 뒤, 용량 확장·워크로드 정리·리소스 상향 등을 검토한다.

---

## 7. Practical Constraints (실무적 제약)

### 7.1 Prometheus 쿼리 비용

- **histogram_quantile + rate(...[5m])** (API server P99 등)는 range 쿼리로 부하가 있을 수 있다. 대시보드 **새로고침 주기 1–2분 이상** 권장.  
- **increase(...[1h])** (재시작 TOP10 등)는 **retention** 구간 안에서만 유의미. retention이 짧으면 구간을 줄이거나 “현재 값만” threshold로 사용.  
- **topk(10, ...)** 는 상대적으로 가볍지만, 내부 식이 무거우면 스캔 부하 가능. 필요 시 range·집계 구간 조정.

### 7.2 Metric 가용성 차이

- **kube-state-metrics, node_exporter** 가 있으면 Block 1의 NotReady, Pending, Restarts, Endpoint와 Block 2의 Node CPU/메모리/디스크, Block 3의 O1–O4는 대부분 **즉시 사용 가능**.  
- **Control plane (API server, scheduler)** metric은 **Managed K8s**에서는 사용자가 스크래핑하지 못해 **없을 수 있음.** 이 경우 Block 1에서 API server health, Scheduler pending 패널을 **제외**하고, Block 3에서 Control plane component를 제외. 제어면 건강은 **managed 서비스 콘솔·알림**으로 확인.  
- **Endpoint, Ingress** metric은 kube-state-metrics·Ingress controller 버전에 따라 **이름·라벨이 다를 수 있음.** 환경에서 metric을 확인한 뒤 PromQL을 조정해야 함.

### 7.3 Managed Kubernetes 환경

- EKS, GKE, AKS 등에서는 **control plane을 사용자가 스크래핑하지 못하는 경우**가 많다.  
- **Block 1:** API server health, Scheduler pending pods **패널 제거** → **4개**(NotReady, Workload Pending, Excessive restarts, Critical endpoint)만 사용.  
- **Block 2:** API server inflight·scheduler pending trend 대신 **Pending count 추세** 또는 Node CPU/메모리/디스크만 사용 → **4개** 유지.  
- **Block 3:** Control plane component health **제외** → **4개** 뷰.  
- 대시보드 상단 또는 runbook에 **“제어면 건강은 managed 서비스 콘솔·알림으로 확인”** 안내를 넣는다.

---

## 8. v1 Limitations (v1 한계)

### 8.1 Control plane 가시성

- **Managed K8s**에서는 **“제어면이 지금 응답·스케줄하는가?”** 를 이 대시보드의 Block 1으로 볼 수 없다.  
- **대응:** Block 1을 4개 패널로 구성하고, 제어면 건강은 **managed 서비스 알림·콘솔·지원 채널**로 확인한다.

### 8.2 Threshold 정의

- **Excessive restarts** 의 N(10분 내 몇 회 초과 시 비정상), **Node CPU/메모리 80%** 경고, **디스크 여유 10%** 등 **구체적 threshold** 는 **팀·환경에 맞게 반드시 정의** 후 적용해야 한다. 설계 문서는 “예시”만 제시하며, 구현 시 팀 정의가 필요하다.

### 8.3 환경별 쿼리 조정

- **Critical service endpoint empty:** endpoint metric 이름·라벨이 환경마다 다를 수 있음. “비어 있는 엔드포인트” 조건을 PromQL로 환경에 맞게 정의.  
- **Restart TOP10:** deployment 등 워크로드 단위로 보려면 `kube_pod_owner` 등과 조인하거나 라벨 추가 필요.  
- **Pending by workload:** namespace 단위만으로도 가능하나, deployment 등으로 세분화하려면 owner/라벨 조인 필요.  
- **Trend “추세”:** “Pending 추세”를 어떤 time window·어떤 증가량으로 정의할지는 팀에서 정해야 함.

---

## 9. Future Evolution (향후 진화)

대시보드는 **고정 명세가 아니라** alert·runbook·추가 신호와 함께 진화할 수 있다.

### 9.1 Alert 정책 (sre-monitoring-alert-policy)

- Block 1·Block 2 신호 중 **어떤 것을 알림으로 쓸지**, **임계치·심각도**(critical / warning)를 정한다.  
- 예: NotReady node count >0 → critical, Node memory >80% → warning.  
- 노이즈가 많은 신호(readiness probe 실패, 일시적 재시작 1~2회)는 Summary 알림에서 제외하고, drill-down·낮은 심각도만 사용한다는 원칙을 정책에 반영한다.

### 9.2 Runbooks (sre-monitoring-operational-runbooks)

- **“Block 1에서 이 신호가 비정상일 때”** 조사 순서·체크리스트·명령어를 runbook으로 정리한다.  
- 예: NotReady >0 → 노드 점검 runbook, Workload Pending >0 → Pending 원인 조사 runbook, Excessive restarts → Restart TOP10 확인 후 해당 워크로드 runbook.  
- **Control plane 미사용 환경**에서는 “제어면 이상 의심 시 managed 서비스 콘솔·지원 채널 확인”을 runbook에 포함한다.

### 9.3 추가 신호

- **Eviction:** kubelet이 eviction metric을 노출하면 Block 1 또는 Block 2에 추가 검토.  
- **Ingress/서비스 에러율:** Ingress controller metric이 안정화되면 Block 2 Early Risk 또는 Block 3에 고정 반영.  
- **새 카테고리:** 보안·비용 등이 필요해지면 기존 **3블록 구조**는 유지한 채, 해당 카테고리당 0~1개 Confidence 패널만 추가하는 방식으로 확장할 수 있다.

---

## 10. 요약

- **운영 문제:** metric은 많지만 “이게 정상인가?”, “확인해 줄 수 있나?”에 대한 **운영 확신**이 부족함. 목표는 **명확한 운영 확신** 제공.  
- **핵심 질문:** (1) 지금 클러스터가 안전한가? (2) 곧 불안전해질 징후가 있는가? (3) 이상 시 어디를 조사할 것인가?  
- **운영 논리:** **Operational Safety Conditions (C1~C5)** 와 **Failure Propagation Paths**에 따라 패널이 정당화됨. Block 1 = 조건 위반 여부, Block 2 = 전파 앞단(선행 지표), Block 3 = investigation. A 우선순위 패널은 각각 지원 안전 조건·전파 단계·지표 유형(leading/lagging/already-damaged/investigation)으로 매핑됨.  
- **대시보드 모델:** Operational Confidence / Early Risk / Investigation 세 계층. 메인 뷰는 Confidence + Early Risk만, Investigation은 드릴다운 전용.  
- **레이아웃:** Block 1(4~6개 패널), Block 2(4~6개), Block 3(4~6개). **메인 뷰 = Block 1 + Block 2**, 총 10~14개 이하. Block 3은 탭·접기로 “이상 시”만 진입.  
- **운영 워크플로:** 일상 5–10분 점검(Block 1 → Block 2), 이상 시 조사(Block 1 비정상 → 카테고리 → Block 3 TOP10), 조기 리스크 감지(Block 2 경고 → Block 3 또는 사전 대응).  
- **실무 제약:** Prometheus 쿼리 비용·retention, metric 가용성 차이, Managed K8s에서 control plane 미사용. 환경별로 패널 제외·대체로 대응.  
- **v1 한계:** Control plane 가시성 부재(Managed K8s), threshold 팀 정의 필요, 환경별 쿼리 조정 필요.  
- **향후:** Alert 정책, Runbook, Eviction·에러율·새 카테고리 등으로 **3블록 구조를 유지한 채** 확장.

**Review 결론:** 설계는 **v1으로 수용 가능**하며, Block 3 기본 접기/탭 명시, Pending 두 패널 역할 구분(제목·설명), threshold 팀 정의 안내 등 **소규모 조정**만 반영하면 구현·운영에 사용하기에 충분하다.

---

## 11. Validation — 첫 검증 패스 반영

첫 검증 패스(05-implementation/validation-results-first-pass.md) 결과를 반영한 요약이다.

### 11.1 즉시 동작한 항목 (What worked immediately)

- **대시보드:** JSON 임포트, datasource 바인딩, Row 1·2·3 렌더링, refresh 2m 정상.  
- **패널 11개:** NotReady node count, Workload Pending pod count, Excessive restarts, Node CPU utilization, Node memory / OOM risk, Node disk space(mountpoint="/" 환경), Pending pods trend, CPU TOP10 nodes, Memory TOP10 nodes, Restart TOP10 pods, Pending pods by workload — **spec PromQL 그대로** 동작.  
- **운영적 의미:** “지금 안전한가?” / “조기 징후인가?” / “조사 대상” 판단에 사용 가능.

### 11.2 소규모 수정 후 동작한 항목 (What required small fixes)

- **Critical service endpoint empty (P4):** kube-state-metrics 버전에 따라 metric 이름이 **단수** `kube_endpoint_address_available` 또는 **복수** `kube_endpoints_address_available`. 복수형을 기본으로 적용해 두었으며, No data/에러 시 단수형으로 패널 쿼리 교체. (promql-spec, implementation-notes, grafana-dashboard-v1.json 반영.)  
- **Node disk space (T3):** root 디스크가 `mountpoint="/"` 가 아닌 환경에서는 `mountpoint="/var"` 등으로 **환경별 조정** 필요. promql-spec·implementation-notes에 명시.  
- **Excessive restarts (P3):** **Threshold N은 팀에서 반드시 정의** 후 Grafana에 반영(예: N=10~20). JSON 기본값 10은 참고용.

### 11.3 운영 측 후속 조치 (What remains as operational follow-up)

- **팀 정의:** Excessive restarts N 값, Critical endpoint empty의 “critical” namespace 목록(필터 적용 시).  
- **환경별 재확인:** 다른 클러스터/프로메테우스 배포 시 endpoint metric 이름, Node disk mountpoint 재확인 후 필요 시 쿼리 수정.  
- **Alert·Runbook:** sre-monitoring-alert-policy, sre-monitoring-operational-runbooks 프로젝트에서 Block 1·2 신호 알림·조사 순서 정리.  
- **제어면·Ingress:** Managed K8s에서 제어면 건강은 별도 채널. Ingress metric 도입 시 B 패널 추가는 별도 검토.

---

*Central Kubernetes Operational Dashboard Design v1. SRE Monitoring program — sre-monitoring-dashboard-design project. Foundation: Kubernetes Cluster Health Monitoring Model; operational logic: 03-architecture/operational-confidence-theory.md.*
