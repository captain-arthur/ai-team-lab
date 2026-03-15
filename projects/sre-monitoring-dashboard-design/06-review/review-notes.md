# Review notes: Central Kubernetes Operational Dashboard Design

**Project:** sre-monitoring-dashboard-design  
**Phase:** Review  
**Input:** Architecture (`03-architecture/architecture.md`), Panel design (`04-engineering/panel-design.md`), Kubernetes Cluster Health Monitoring Model (final-report)

이 리뷰는 **기술적 정확성**보다 **Kubernetes 클러스터 운영에서의 실질적 유용성**을 기준으로 제안된 대시보드 설계를 평가한다. 운영자가 “지금 안전한가?”, “곧 불안전해질 징후는?”, “어디를 조사할 것인가?”에 5–10분 안에 답할 수 있는지, 조사 흐름이 자연스러운지, 메인 뷰가 과부하 없이 최소인지, metric이 일부 없는 환경에서도 동작하는지 검토한다.

---

## 1. 핵심 평가 질문에 대한 답

### 1.1 운영자가 5–10분 안에 “클러스터가 안전한가?”를 판단할 수 있는가?

**평가: 예. 설계대로 구현하면 가능하다.**

- **메인 뷰**는 Block 1(Operational Confidence) + Block 2(Early Risk)만 포함하며, **패널 수는 8~12개(최대 14개 이하)** 로 제한되어 있다. Block 1만 봐도 “전부 정상 = 안전, 하나라도 비정상 = 조사 필요”라는 **단일 규칙**으로 결론을 내릴 수 있다.
- Block 1 패널은 **stat** 타입으로, **현재 값 + 정상/비정상(threshold·색)** 이 한눈에 보이도록 되어 있다. 숫자를 해석할 필요 없이 **상태(녹색/빨강)** 만 보면 되므로, 5–10분 일상 점검에 적합하다.
- Block 1이 **한 행(또는 두 행)** 에 모여 스크롤 없이 보이도록 되어 있어, “한눈에 전부 정상인지” 확인하는 데 걸리는 시간이 짧다.
- **결론:** 설계는 **5–10분 내 “클러스터가 안전한가?”** 판단을 지원한다. 구현 시 Block 1이 실제로 **최상단 한 화면**에 들어오고, threshold·value mapping이 적용되면 목표가 달성된다.

### 1.2 Early Risk 패널이 “클러스터가 곧 불안전해질 수 있다”는 신호를 명확히 보여 주는가?

**평가: 예. 조기 리스크가 한 블록에 묶여 가시성이 높다.**

- **Block 2 (Early Risk)** 에는 Node CPU utilization, Node memory/OOM risk, Node disk space, Pending pods trend가 **항상** 포함되고, (선택) CPU throttling risk, Ingress stress가 **환경 가능 시** 추가된다. 이들은 모델에서 정의한 **조기 리스크**(노드 압박, OOM 위험, 디스크 부족, pending 증가, ingress 부하)와 1:1로 대응한다.
- **gauge** 또는 **stat** + **threshold**(예: 80% 초과 = 경고)로 “지금은 괜찮지만 곧 나빠질 수 있는” 구간이 **색(노랑/빨강)** 으로 드러나도록 되어 있다. Block 2 제목(“Early Risk” 또는 “조기 리스크”)으로 **이 블록의 목적**이 분명해, 운영자가 “Summary는 정상인데 여기서 경고가 나오면 사전 대응이 필요하다”고 해석하기 쉽다.
- **결론:** Early Risk 패널은 **“클러스터가 곧 불안전해질 징후”** 를 명확히 드러내도록 설계되어 있다. 구현 시 Block 2의 threshold(80% 등)와 색 구간이 일관되게 적용되면 조기 리스크 가시성은 충분하다.

### 1.3 이상 신호 시 자연스러운 조사 흐름이 있는가? (Confidence → Early Risk → Investigation)

**평가: 예. 아키텍처와 패널 설계에 흐름이 명시되어 있다.**

- **기대 흐름:** Operational Confidence → Early Risk → Investigation / Top Offenders.  
  **설계 반영:**  
  - **Block 1** 에서 비정상 발견 → “조사 필요” → **어느 신호(NotReady / Pending / Restarts / Endpoint 등)** 인지 확인 → **카테고리**(Node / Workload / Network 등) 결정 → **Block 3** 에서 해당 Investigation 패널(CPU TOP10, Memory TOP10, Restart TOP10, Pending by workload 등)로 진입.  
  - **Block 2** 에서만 경고(Summary는 정상) → “곧 나빠질 수 있음” 인지 → **유형**(노드 압박 / 디스크 / pending 추세 등) 확인 → 필요 시 Block 3 해당 TOP10으로 **어느 노드·워크로드**가 부담을 주는지 좁힌다.
- Architecture에 **판단 흐름(decision flow)** 이 단계별로 기술되어 있고, Panel design에서 Block 3 패널이 “Block 1/2의 어떤 신호일 때 드릴다운하는지”(예: “Block 2 Node CPU 높음일 때 CPU TOP10”)가 **Why the panel exists** 로 연결되어 있다. 따라서 **Confidence → Early Risk → Investigation** 순서가 설계에 반영되어 있으며, 구현만 맞추면 자연스러운 조사 흐름이 나온다.
- **결론:** 이상 시 **Operational Confidence → Early Risk → Investigation / Top Offenders** 로 이어지는 조사 흐름이 설계에 포함되어 있고, 운영 워크플로와 부합한다.

### 1.4 대시보드가 최소로 유지되어 과부하가 없는가? (메인 뷰 소규모, 신호 중복 없음, 높은 signal-to-noise)

**평가: 예. 메인 뷰가 작고, 중복이 없으며, 신호 가치가 높다.**

- **메인 뷰 규모:** 메인 뷰에는 **Block 1 + Block 2만** 들어가고, **총 10~14개 이하**(표준 8~12개)로 제한되어 있다. Block 3(Investigation)은 **메인 뷰에 포함하지 않고** 탭·접기로 “이상 시”만 진입하도록 되어 있어, **메인 뷰는 작게 유지**된다.
- **신호 중복:**  
  - Block 1의 “Workload Pending pod count”는 **현재 Pending 수**(지금 안전한가?). Block 2의 “Pending pods trend”는 **Pending 추세 또는 현재값의 트렌드**(곧 나빠질 징후인가?). 질문이 다르므로 **의도적 이중 배치**이며, 동일 정보를 두 번 나열하는 중복은 아니다.  
  - Node CPU·메모리는 Block 2에만(조기 리스크), Block 3에는 TOP10(조사 대상)으로만 나오므로, **계층이 나뉘어** 중복이 아니다.  
  **결론:** 설계상 **불필요한 신호 중복**은 없다.
- **Signal-to-noise:** 메인 뷰에는 **모니터링 모델의 Summary + Trend-Risk** 에 해당하는 소수 신호만 들어가고, 노이즈가 많은 신호(readiness probe 실패, 일시적 재시작 1~2회 등)는 모델에서 이미 Summary에 넣지 않기로 되어 있다. 따라서 **높은 signal-to-noise** 를 유지하는 설계다.
- **결론:** 메인 뷰는 **최소**이고, 신호는 **중복되지 않으며**, 패널은 **높은 signal-to-noise** 를 제공하도록 구성되어 있다. 과부하를 피하기에 적절하다.

### 1.5 metric이 일부 없는 환경(예: Managed K8s, control plane 미사용)에서도 설계가 동작하는가?

**평가: 예. 환경별 변형이 명시되어 있다.**

- **Panel design** 과 **Architecture** 에서 **Managed Kubernetes (control plane metric 없음)** 시:  
  - **Block 1:** API server health(P2), Scheduler pending pods(P3) **제외** → **4개 패널**(NotReady, Workload Pending, Excessive restarts, Critical endpoint)만 사용.  
  - **Block 2:** API server inflight·scheduler pending trend 대신 **S4 Pending count 추세** 또는 T1, T4만 사용 → **4개** 유지.  
  - **Block 3:** Control plane component health(O6) **제외** → **4개** 뷰.  
  “제어면 건강은 managed 서비스 콘솔·알림으로 확인”하는 안내를 두도록 되어 있어, **metric이 없어도 운영 판단은 가능**하다.
- **결론:** control plane metric이 없는 환경에서도 **메인 뷰 4+4=8개**, **Investigation 4개**로 **“지금 안전한가?” / “곧 위험한가?” / “어디를 조사할 것인가?”** 에 답하는 설계가 유지된다. **일부 metric이 없어도 설계는 동작한다.**

---

## 2. 실제 운영 워크플로 관점 평가

### 2.1 일상 점검(5–10분)

- **워크플로:** 아침 또는 정해진 시간에 대시보드만 보고 “클러스터 안전한가?” / “곧 위험한가?” 결론.
- **설계 부합:** Block 1만 보면 “전부 정상 = 안전”, “하나라도 비정상 = 조사 필요”. Block 2만 훑으면 “조기 리스크 유무”. **스크롤·탭 전환 없이** 메인 뷰만으로 5–10분 내 결론이 나오도록 되어 있어 **일상 점검 워크플로와 잘 맞는다.**

### 2.2 이상 시 조사

- **워크플로:** Block 1 또는 Block 2에서 비정상/경고 발견 → 원인 후보 좁히기 → kubectl/로그/이벤트 조사.
- **설계 부합:** 비정상 **신호 → 카테고리** 식별 후 **Block 3** 해당 TOP10(table)으로 진입하도록 되어 있다. “NotReady 있음 → Node 쪽 → CPU/Memory TOP10 노드”, “Pending 있음 → Pending by workload”, “Excessive restarts → Restart TOP10” 등 **매핑이 명시**되어 있어 **이상 시 조사 워크플로와 일치**한다.

### 2.3 사전 대응(용량·정리 검토)

- **워크플로:** Summary는 정상이지만 “곧 나빠질 수 있다”는 인지 → 용량 확장·워크로드 정리 검토.
- **설계 부합:** Block 2(Early Risk)가 **메인 뷰에 항상** 포함되어 있어, Block 1이 전부 녹색이어도 Block 2에서 노드 사용률 80% 근접·디스크 여유 부족·Pending 추세 등을 보고 **사전 대응**을 트리거할 수 있다. **사전 대응 워크플로와 부합**한다.

---

## 3. 강점 (Strengths)

1. **정보 계층이 명확함:** Operational Confidence / Early Risk / Investigation 세 계층이 **한 가지 운영 질문씩** 담당하고, **블록 순서(상단 → 중간 → 드릴다운)** 와 일치해, “무엇을 어디서 보는지”가 분명하다.
2. **메인 뷰 최소화:** Block 1 + Block 2만 메인 뷰, 10~14개 이하, Block 3은 탭·접기로 분리해 **과부하 방지**가 설계에 반영되어 있다.
3. **판단 규칙이 단순함:** “Block 1 전부 정상 = 안전”, “하나라도 비정상 = 조사 필요”로 **운영 확신**을 높이기 좋다. 추가 해석 규칙이 거의 필요 없다.
4. **조기 리스크가 한 블록에 집약됨:** Node CPU/메모리/OOM, 디스크, Pending 추세, (선택) CPU throttling·Ingress stress가 **Block 2 한 곳**에 있어, “곧 나빠질 징후”를 놓치기 어렵다.
5. **조사 흐름이 설계에 포함됨:** Confidence → Early Risk → Investigation 순서와 “어떤 신호일 때 어떤 TOP10으로 들어가는지”가 Architecture·Panel design에 적혀 있어, **구현·runbook 작성**에 그대로 쓸 수 있다.
6. **환경별 변형이 정의됨:** Managed K8s 등 **metric 부재** 시 Block 1·2·3의 패널 구성(제외·대체)이 문서화되어 있어, **다양한 환경에서 동일 설계 원칙**을 적용할 수 있다.

---

## 4. 약점 및 리스크 (Weaknesses / Risks)

1. **Block 3 기본 숨김 구현 누락 가능성:** Block 3이 **탭 또는 접기로 기본 숨김**이어야 메인 뷰가 최소로 유지되는데, 구현 단계에서 Row를 그냥 펼쳐 두면 **메인 뷰가 길어져** 5–10분 목표가 약해질 수 있다. **구현 가이드에 “Block 3은 기본 접기 또는 별도 탭”** 을 명시해 두는 것이 좋다.
2. **Pending이 Block 1과 Block 2 양쪽에 등장:** Block 1 “Workload Pending pod count”와 Block 2 “Pending pods trend”가 **의도적으로 다른 질문**(현재 안전한가 vs 추세)에 쓰이지만, 운영자가 “Pending이 두 번 나온다”고 느낄 수 있다. 패널 제목·설명으로 **“현재 수” vs “추세/조기 리스크”** 를 구분해 두면 혼동을 줄일 수 있다.
3. **Node CPU와 Node memory가 Block 2에서 별도 패널:** 공간을 아끼려면 하나의 “Node resource pressure” 패널에 CPU·메모리 두 stat을 넣을 수 있으나, 현재 설계는 **각각 gauge/stat** 으로 두 패널. **가독성·조기 리스크 강조** 측면에서는 현재가 유리하고, 패널 수도 10~14 이하를 유지하므로 **변경 필수는 아님.** 다만 “패널 수를 더 줄이고 싶을 때”의 옵션으로만 문서에 적어 둘 수 있다.
4. **Threshold·값 정의는 팀 책임:** “Excessive restarts”의 N, “80% 초과 = 경고” 등 **구체적 숫자**는 Panel design에서 예시 수준이고, **실제 threshold** 는 팀·환경에 따라 정해야 한다. 설계 문서에 **“threshold는 팀 정의, 구현 시 반드시 설정”** 한 줄을 넣으면 누락을 줄일 수 있다.

---

## 5. 소규모 개선 제안 (신호·복잡도 증가 없이)

1. **Block 3 기본 상태 명시:** Panel design 또는 구현 가이드에 **“Block 3(Investigation) Row는 Grafana에서 기본 접기(collapsed) 또는 별도 탭으로 두어, 메인 뷰에는 Block 1·2만 노출”** 를 한 문장으로 넣는다. 구현 시 실수로 펼쳐 두는 것을 방지한다.
2. **Pending 두 패널의 역할 구분:** Block 1 패널 제목을 **“Workload Pending pod count (현재)”**, Block 2를 **“Pending pods trend (조기 리스크)”** 처럼 괄호로 구분하거나, 패널 설명에 “현재 수 = 지금 안전한가?” / “추세 = 곧 나빠질 징후?” 한 줄을 넣어 **의도적 이중 배치**임을 분명히 한다.
3. **Threshold 팀 정의 안내:** Panel design 요약 또는 구현 시 참고 사항에 **“Block 1·2의 threshold(0/비정상, 80% 경고 등)는 팀·환경에 맞게 반드시 정의 후 적용”** 을 한 줄 추가한다.
4. **(선택) Node CPU + memory 단일 패널 옵션:** “메인 뷰 패널을 더 줄이고 싶을 때”만, Block 2에서 Node CPU와 Node memory를 **한 stat 패널에 두 개 값**으로 넣는 옵션을 **문서에 선택 사항**으로 적어 둔다. 기본 설계는 현재처럼 두 패널 유지.

위 개선은 **패널 수·신호 수·구조를 늘리지 않고** 문서·구현 가이드만 보완하는 수준이다.

---

## 6. 최종 결론

**평가 요약**

- 운영자가 **5–10분 안에** “클러스터가 안전한가?”를 판단할 수 있는 구조(Block 1 중심, stat + threshold)가 갖춰져 있다.
- **Early Risk** 패널(Block 2)이 “클러스터가 곧 불안전해질 징후”를 한 블록에서 명확히 보여 주도록 설계되어 있다.
- **조사 흐름** Operational Confidence → Early Risk → Investigation / Top Offenders가 아키텍처·패널 설계에 반영되어 있고, 이상 시 자연스럽게 드릴다운할 수 있다.
- **메인 뷰**는 10~14개 이하로 최소화되어 있고, 신호 중복 없이 **높은 signal-to-noise** 를 유지한다.
- **일부 metric이 없는 환경**(Managed K8s 등)에서도 **패널 제외·대체**로 동일 설계 원칙을 적용할 수 있다.

**결론: 대시보드 설계는 v1으로 수용 가능하다. 소규모 조정을 권장한다.**

- **The dashboard design is acceptable for v1.**
- **Minor adjustments** 로 다음을 권장한다:  
  - Block 3을 **기본 접기 또는 별도 탭**으로 두는 것을 구현 가이드에 명시.  
  - Block 1/2의 **Pending 두 패널** 역할 구분(제목 또는 설명) 보강.  
  - **Threshold는 팀 정의** 라는 안내 한 줄 추가.  
  - (선택) Node CPU+memory 단일 패널 옵션을 문서에 선택 사항으로 기록.

**Significant redesign은 필요하지 않다.** 정보 계층, 메인 뷰 최소화, 조기 리스크 가시성, 조사 흐름, 환경별 변형이 운영 유용성과 맞아 떨어지며, 제안된 소규모 개선만 반영하면 **Documentation 단계에서 최종 설계 문서를 정리**하고 구현으로 넘어가기에 충분하다.

---

*Review phase output. 다음 단계: Documentation — 최종 설계 문서 갱신(07-documentation 초안을 아키텍처·패널 설계·Review 반영으로 정리).*
