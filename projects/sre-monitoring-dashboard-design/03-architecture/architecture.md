# Central Kubernetes Operational Dashboard — Architecture

**Project:** sre-monitoring-dashboard-design  
**Phase:** Architecture  
**Input:** Manager brief, Cluster Health Monitoring Model (`sre-monitoring-cluster-health-model`), 설계 초안(draft) `07-documentation/central-kubernetes-operational-dashboard-design.md`

이 문서는 **Central Kubernetes Operational Dashboard** 의 **아키텍처**를 정의한다. 기존 설계 초안은 참고하되 **검증된 최종 결과로 간주하지 않고**, 모니터링 모델에 기반한 **대시보드 구조·정보 계층·메인 뷰 vs 드릴다운 구분·판단 흐름**을 명시적으로 정리한다. Engineering·Review·Documentation 단계에서 이 아키텍처를 기준으로 구체화·검증·문서화한다.

---

## 1. 목적과 입력

### 1.1 아키텍처 목적

- **대시보드 구조를 모니터링 모델에 맞게 정식화:** Cluster Health Monitoring Model의 **Cluster Health Summary / Trend-Risk / Top Offenders** 를 대시보드의 **정보 계층**과 **화면 배치**로 1:1 매핑한다.
- **정보 계층 명확화:** Operational Confidence / Early Risk / Investigation(Top Offenders) 세 계층이 **무엇을 담당하고, 어디에 배치되며, 어떤 판단을 유도하는지**를 규정한다.
- **메인 뷰 vs 2차/드릴다운 규정:** “한 화면에 반드시 보이는 것”과 “이상 시 또는 필요 시에만 보는 것”을 구분해, **대시보드 과부하**를 막고 **5–10분 일상 점검**에 맞춘다.
- **판단 흐름 명시:** 운영자가 **어떤 순서로 무엇을 보고, 어떤 결론을 내리며, 언제 드릴다운하는지**를 결정 흐름으로 기술한다.

### 1.2 입력 (참조만, 최종 아님)

| 입력 | 용도 |
|------|------|
| **Kubernetes Cluster Health Monitoring Model** (sre-monitoring-cluster-health-model: architecture, final-report, core-signal-list) | 대시보드가 반영해야 할 **신호 집합·레이어·카테고리**의 기준. |
| **Experiment notes** (metric 가용성, Managed K8s 제한) | 어떤 신호를 **즉시/수정 후/미사용**으로 둘지, 환경별 변형 근거. |
| **설계 초안** (07-documentation/central-kubernetes-operational-dashboard-design.md) | 레이아웃·블록 이름·조기 리스크 강조 아이디어 참고. **초안이므로 이 아키텍처에서 구조·계층·판단 흐름을 정식화한 뒤, 이후 단계에서 검증·수정한다.** |

---

## 2. 설계 목표와 비목표

### 2.1 설계 목표 (최적화 대상)

| 목표 | 의미 |
|------|------|
| **5–10분 일상 점검** | 운영자가 **메인 뷰만** 보고 “지금 안전한가?” / “곧 위험한가?”에 답할 수 있어야 함. 스크롤·탭 전환·패널 해석에 시간이 많이 들면 목표 미달. |
| **높은 운영 확신** | “전부 정상 = 안전”, “하나라도 비정상 = 조사 필요”가 **한눈에** 보여, “이게 정상인가?”라는 질문이 나오지 않도록 함. |
| **조기 리스크 가시성** | “곧 불안전해질 징후”(노드 압박, OOM 위험, pending 증가, CPU throttling, ingress stress 등)가 **메인 뷰 또는 그 직하단**에 묶여 있어, Summary가 정상이어도 **조기 리스크**를 놓치지 않도록 함. |
| **낮은 대시보드 과부하** | 메인 뷰에 나열하는 **패널·신호 수를 엄격히 제한**. 2차·드릴다운은 “이상 시” 또는 “더 보기”로 진입하도록 분리. |

### 2.2 비목표

- **metric 수 최대화:** 가능한 모든 metric을 한 대시보드에 넣는 것이 아님.
- **기술적 완결성:** Prometheus/Grafana 구현 세부(쿼리 최적화, 변수 스키마)는 Architecture가 아니라 Engineering 단계에서 다룸.
- **Alert·Runbook 상세:** 이 대시보드 아키텍처는 “어떤 화면 구조로 무엇을 보여 줄지”만 규정. Alert 정책·Runbook은 별도 프로젝트.

---

## 3. 정보 계층 (Information hierarchy)

대시보드는 **세 가지 정보 계층**으로 구분된다. 각 계층은 **모니터링 모델의 한 레이어**에 대응하고, **한 가지 운영 질문**에 답한다.

### 3.1 계층 1: Operational Confidence (운영 확신)

| 항목 | 내용 |
|------|------|
| **대응하는 모델 레이어** | Cluster Health Summary |
| **답하는 질문** | **지금 클러스터가 안전하다고 확신할 수 있는가?** |
| **표시 내용** | 모델의 Summary 신호(S1–S6)에 해당하는 **최소 개수(5~7개, 환경별 4~5개)** 의 지표. 각 지표는 **현재 값 + 정상/비정상 상태**. |
| **판단 규칙** | **전부 정상** → “클러스터 안전”. **하나라도 비정상** → “조사 필요”. |
| **배치** | **메인 뷰의 최상단.** 스크롤 없이 한 화면에 전부 들어와야 함. |

**계층 1에 올 수 있는 신호(모델 기준):**  
NotReady node count (S1), API server health (S2), Scheduler pending pods (S3), Workload Pending pod count (S4), Excessive restarts (S5), Critical service endpoint empty (S6).  
환경에 따라 S2, S3가 없으면 S1, S4, S5, S6만 사용. **이 계층에 추가 신호를 넣지 않는다** — “운영 확신”은 이 소수만으로 판단하는 것이 원칙.

### 3.2 계층 2: Early Risk (조기 리스크)

| 항목 | 내용 |
|------|------|
| **대응하는 모델 레이어** | Trend / Risk Indicators |
| **답하는 질문** | **클러스터가 곧 불안전해질 조기 징후가 보이는가?** |
| **표시 내용** | 모델의 Trend-Risk 신호(T1–T4 등) + **조기 리스크 강조** 항목: node resource pressure, OOM risk, CPU throttling(가능 시), pending pods(추세), ingress controller stress(가능 시). |
| **판단 규칙** | **경고/위험** 수준이 있으면 “곧 불안정해질 수 있음” → 용량·정리·확장 검토. Summary가 전부 정상이어도 이 계층에서 경고가 나올 수 있음. |
| **배치** | **메인 뷰에서 Operational Confidence 직하단.** 일상 점검 시 “한 번만 더 훑는” 구간. 스크롤 한 번 이내 권장. |

**계층 2에 올 수 있는 신호(모델·초안 기준):**  
Node CPU/memory utilization (T1), Node disk (T4), Pending trend (T3 또는 S4 추세), API server inflight (T2, 환경 가능 시).  
**강조할 조기 리스크:** Node resource pressure (= T1), OOM risk (= T1 메모리·O2 연계), CPU throttling(metric 있을 때 별도, 없으면 T1으로 대체), Pending pods 수·추세, Ingress controller stress(metric 있을 때).  
이 계층도 **패널 수를 제한**(예: 4~6개). “모든 Trend metric”을 나열하지 않는다.

### 3.3 계층 3: Investigation / Top Offenders (조사 / 드릴다운)

| 항목 | 내용 |
|------|------|
| **대응하는 모델 레이어** | Top Offenders / Drill-down |
| **답하는 질문** | **문제가 있다면 어디를 먼저 조사할 것인가?** |
| **표시 내용** | 모델의 Top Offenders 뷰(O1–O6): CPU TOP10 nodes, Memory TOP10 nodes, Restart TOP10 workloads, Pending by workload, Error-heavy ingress(선택), Control plane component(선택). |
| **판단 규칙** | **ranking** — 상위 N개(TOP10)만 표시. Summary 또는 Early Risk에서 “어떤 카테고리가 나쁜가?”가 나오면, 해당 카테고리에 맞는 뷰를 열어 **원인 후보**를 좁힌다. |
| **배치** | **2차 뷰(secondary view)** 또는 **드릴다운.** 메인 뷰에 항상 펼쳐 두지 않고, **탭·접기·별도 행** 등으로 “이상 시” 또는 “더 보기”로 진입. |

**계층 3은 메인 뷰에 필수로 보이지 않아도 된다.** 5–10분 일상 점검에서는 계층 1·2만 보면 되고, “조사 필요”일 때만 계층 3으로 들어가는 구조가 **과부하 방지**와 **운영 확신**에 맞다.

---

## 4. 메인 뷰 vs 2차/드릴다운 뷰

### 4.1 규정

| 구분 | 정의 | 배치 원칙 |
|------|------|-----------|
| **메인 뷰 (Primary view)** | 운영자가 **매일 5–10분** 점검할 때 **반드시 보는 영역**. 스크롤 최소화(한 화면 또는 스크롤 한 번). | **계층 1 (Operational Confidence)** 전부 + **계층 2 (Early Risk)** 전부. |
| **2차 뷰 / 드릴다운 (Secondary / Drill-down)** | **이상이 있을 때** 또는 “어디가 문제인지 더 보고 싶을 때” 진입하는 영역. | **계층 3 (Investigation / Top Offenders)**. 탭·접기·아래쪽 별도 섹션 등으로 메인 뷰와 분리. |

### 4.2 메인 뷰에 반드시 나와야 하는 것

- **Operational Confidence:** Summary에 해당하는 **5~7개(또는 4~5개)** 신호. 각각 **현재 값 + 정상/비정상(색·아이콘)**. “전부 정상 = 안전”이 **한눈에** 보여야 함.
- **Early Risk:** Trend-Risk 및 조기 리스크 강조 신호 **4~6개**. “곧 나빠질 수 있다”는 메시지가 **같은 블록 안**에서 파악 가능해야 함.

**메인 뷰에 넣지 않는 것:**  
- TOP10 테이블 전체(계층 3). 단, “Early Risk” 블록 안에서 “메모리 압박 노드 수”처럼 **집계된 숫자 하나**는 허용(예: “노드 3개가 메모리 80% 초과”).  
- 상세 목록·개별 노드/파드 이름 나열은 **드릴다운**으로 미룬다.

### 4.3 2차/드릴다운에 두는 것

- **Investigation / Top Offenders** 전부: O1 CPU TOP10 nodes, O2 Memory TOP10 nodes, O3 Restart TOP10, O4 Pending by workload, O5 Error-heavy ingress(선택), O6 Control plane component(선택).  
- 각 뷰는 **TOP10 한 테이블(또는 동등한 단순 구조)** 수준. “조사 필요”일 때 해당 카테고리 탭/섹션만 열어도 되도록 구성.

### 4.4 결정 흐름 (메인 vs 드릴다운)

```
메인 뷰만 본다
    → Operational Confidence 전부 정상?
        → 예: Early Risk에 경고 있나?
            → 예: Early Risk 블록에서 “어느 카테고리(노드/용량/ pending 등)”인지 확인
                → 필요 시 Investigation 해당 뷰(TOP10)로 진입
            → 아니오: “클러스터 안전” 결론
        → 아니오: “조사 필요”. 비정상인 신호의 카테고리 확인
            → Investigation 해당 뷰(TOP10)로 진입
            → TOP10에서 1~2개 후보 선택 후 kubectl/로그/이벤트 조사
```

---

## 5. 대시보드 구조 (정식화)

### 5.1 블록 구조

대시보드는 **세 개의 논리 블록**으로 구성한다. 블록 순서는 **정보 계층 순서**와 동일하다.

| 블록 | 대응 계층 | 화면 위치 | 메인 뷰 포함 여부 |
|------|-----------|-----------|-------------------|
| **Block 1: Operational Confidence** | 계층 1 | 최상단, 1~2행 | **예** — 반드시 메인 뷰에 포함 |
| **Block 2: Early Risk** | 계층 2 | Block 1 직하단, 1~2행 | **예** — 반드시 메인 뷰에 포함 |
| **Block 3: Investigation / Top Offenders** | 계층 3 | Block 2 하단 또는 별도 탭/접기 | **아니오** — 2차/드릴다운. 이상 시 진입 |

### 5.2 블록별 내용 규정

**Block 1 (Operational Confidence)**  
- **목적:** “지금 클러스터가 안전한가?”에 대한 **확신** 제공.  
- **내용:** 모델 Summary 신호 S1–S6 중 환경에서 사용 가능한 것만(최소 4개, 최대 7개). 각 패널: **신호 이름 + 현재 값 + 상태(정상/비정상).**  
- **제한:** 이 블록에 **새 신호를 추가하지 않음**. 모델의 Summary 집합을 넘어서지 않음.

**Block 2 (Early Risk)**  
- **목적:** “곧 불안전해질 징후는?”에 대한 **가시성** 제공.  
- **내용:** 모델 Trend-Risk(T1–T4 등) + **조기 리스크 강조**: node resource pressure, OOM risk, CPU throttling(가능 시), pending pods(수·추세), ingress controller stress(가능 시).  
- **제한:** 패널 수 4~6개. “모든 Trend metric”을 나열하지 않음. **정상/경고/위험** 수준을 색·아이콘으로 통일.

**Block 3 (Investigation / Top Offenders)**  
- **목적:** “어디를 먼저 조사할 것인가?”에 대한 **진입점** 제공.  
- **내용:** O1–O6에 해당하는 TOP10(또는 동등) 뷰. 카테고리별로 탭/섹션 분리 가능.  
- **제한:** 메인 뷰에 항상 펼쳐 두지 않음. 탭·접기·링크로 “필요 시”만 노출.

### 5.3 배치 제약 (과부하 방지)

- **메인 뷰(Block 1 + Block 2) 패널 수:** 합계 **10~14개** 이하 권장. Block 1만으로 4~7개이므로 Block 2는 4~6개 수준.  
- **한 화면:** Block 1은 **스크롤 없이** 한 화면에 들어와야 함. Block 2는 스크롤 한 번 이내.  
- **Block 3:** 별도 탭이면 메인 뷰 패널 수에 포함하지 않음. 같은 페이지 하단이면 “접기”로 기본 숨김 권장.

---

## 6. 운영자 판단 흐름 (Decision flow)

운영자가 대시보드를 **어떤 순서로 보고 어떤 결론을 내리는지**를 명시한다.

### 6.1 일상 점검(5–10분) — 기본 흐름

1. **대시보드 진입** → **Block 1 (Operational Confidence)** 만 먼저 본다.
2. **Block 1 전부 정상인가?**
   - **예** → 3단계로.
   - **아니오** → “조사 필요” 결론. **어느 신호가 비정상인지** 확인 → 해당 **카테고리**(Node / Control Plane / Workload / Network) 결정 → **Block 3** 에서 해당 Investigation 뷰(O1–O6 중 해당)로 진입 → TOP10에서 1~2개 후보 선택 → kubectl/로그/이벤트로 조사. (끝)
3. **Block 2 (Early Risk)** 를 훑는다.
   - **경고/위험 없음** → “클러스터 안전” 결론. (끝)
   - **경고/위험 있음** → “곧 불안정해질 수 있음” 인지. **어느 유형인지**(노드 압박 / 디스크 / pending / ingress 등) 확인 → 필요 시 Block 3 해당 뷰로 진입해 **어느 노드·워크로드**가 부담을 주는지 확인 → 용량·정리·확장 검토. (끝)

### 6.2 “이상 시” 조사 흐름

- **Trigger:** Block 1에서 비정상 1개 이상, 또는 Block 2에서 경고/위험.
- **행동:** Block 1 비정상 → **카테고리** 식별 → Block 3 해당 TOP10 뷰 → 원인 후보 좁히기.  
  Block 2만 경고 → **유형** 식별 → Block 3 해당 뷰(예: 노드 압박이면 O1/O2) → 부담이 큰 노드/워크로드 확인.

### 6.3 계층별 “보는 목적” 요약

| 계층 | 운영자가 묻는 것 | 결론 유형 |
|------|------------------|-----------|
| Operational Confidence | 지금 안전한가? | 안전 / 조사 필요 |
| Early Risk | 곧 나빠질 징후는? | 괜찮음 / 사전 대응 필요 |
| Investigation | 어디를 먼저 볼 것인가? | 조사 대상 후보(노드·워크로드) |

---

## 7. 조기 리스크 강조 (Early Risk visibility)

설계 목표 중 **조기 리스크 가시성**을 만족하기 위해, 다음을 **Block 2 (Early Risk)** 에서 명시적으로 다룬다.

| 조기 리스크 | 모델 대응 | Block 2에서의 표현 |
|-------------|-----------|---------------------|
| **Node resource pressure** | T1 (Node CPU/memory utilization) | 노드 CPU·메모리 사용률(예: >80% 경고). 집계(예: “80% 초과 노드 수”) 또는 대표값. |
| **OOM risk** | T1 메모리, O2 (Memory TOP10) 연계 | “OOM 위험” 또는 “메모리 압박” 라벨. 노드 메모리 사용률·또는 메모리 압박 노드 수. |
| **CPU throttling** | (모델 확장 또는 T1 대체) | 컨테이너 CPU throttling metric이 있으면 별도 패널. 없으면 T1 노드 CPU로 간접 표시. |
| **Pending pods** | S4, T3 (pending 추세) | Pending pod 수·또는 “Pending 추세”(증가 시 스케줄·용량 리스크). |
| **Ingress controller stress** | Ingress metric, O5 연계 | Ingress 지연·에러율·부하(환경에서 metric 있을 때만). “Ingress stress” 라벨. |

이 항목들은 **Block 2 안에만** 두고, **상세 TOP10 목록**은 Block 3으로 미룬다. Block 2에서는 “경고가 있는지/없는지”와 “어느 유형인지”만 빠르게 파악할 수 있으면 된다.

---

## 8. 환경별 변형

### 8.1 Managed Kubernetes (control plane metric 없음)

- **Block 1:** S2(API server health), S3(Scheduler pending) **패널 제거 또는 비표시**. S1, S4, S5, S6만 사용.  
- **Block 2:** T2(API server inflight), T3(scheduler pending trend) **제거 또는 대체**. S4 Pending count 추세로 “pending 증가” 표현. T1, T4 유지.  
- **Block 3:** O6(Control plane component) **제외**.  
- **안내:** 대시보드 상단 또는 문서에 “제어면 건강은 managed 서비스 콘솔·알림으로 확인” 한 줄.

### 8.2 Metric 부재 시

- **Experiment notes** 의 “즉시 사용 가능” 신호만 Block 1·2에 반영. “수정 후 사용”은 PromQL·라벨 조정 후 Engineering 단계에서 추가. “현재 불가”는 패널을 두지 않거나 placeholder + 설명.

---

## 9. 아키텍처 결정 요약

| 결정 | 내용 |
|------|------|
| **정보 계층** | 3계층: Operational Confidence / Early Risk / Investigation. 모델의 Summary / Trend-Risk / Top Offenders와 1:1 대응. |
| **메인 뷰** | Block 1 + Block 2만. Block 3은 2차/드릴다운. |
| **메인 뷰 패널 수** | Block 1: 4~7개. Block 2: 4~6개. 합계 10~14개 이하 권장. |
| **판단 흐름** | Block 1 전부 정상? → 아니오면 Block 3 해당 뷰로 조사. 예면 Block 2 확인 → 경고 있으면 유형 파악 후 Block 3 또는 사전 대응. |
| **조기 리스크** | Node pressure, OOM risk, CPU throttling, Pending, Ingress stress를 Block 2에서 명시. 상세 TOP10은 Block 3. |
| **과부하 방지** | TOP10·상세 목록은 메인 뷰에 두지 않음. 탭·접기로 “이상 시”만 진입. |

---

## 10. 다음 단계

- **Engineering:** 이 아키텍처를 기준으로 **패널 목록·Grafana 배치 초안·PromQL 매핑**을 작성. Block 1·2·3별로 “어떤 패널에 어떤 신호를 넣을지” 구체화.
- **Review:** “5–10분 일상 점검”, “운영 확신”, “조기 리스크 가시성”, “과부하 방지” 네 가지가 설계에 반영되었는지 검증.
- **Documentation:** 검증된 설계를 **최종 설계 문서**로 정리. 기존 07-documentation 초안은 이 아키텍처·Review 반영 후 **갱신**한다.

---

*Architecture phase output. 설계 초안은 참고용이며, 이 아키텍처가 구조·계층·판단 흐름의 기준이다. 다음 단계: Engineering.*
