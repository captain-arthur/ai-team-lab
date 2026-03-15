# Review notes: Cluster Health Monitoring Model

**Project:** sre-monitoring-cluster-health-model  
**Phase:** Review  
**Input:** Architecture (`03-architecture/architecture.md`), Core signal list (`04-engineering/core-signal-list.md`), Experiment notes (`05-experiment/experiment-notes.md`)

이 리뷰는 **기술적 정확성**보다 **실제 Kubernetes 클러스터 운영에서의 유용성**을 기준으로 제안된 모니터링 모델을 평가한다. 운영 워크플로에 맞는지, 5–10분 내 판단·조기 리스크·조사 대상 좁히기가 가능한지, 신호 수·노이즈·제한 환경이 고려되었는지 검토한다.

---

## 1. 핵심 평가 질문에 대한 답

### 1.1 “지금 클러스터가 건강한가?”를 5–10분 안에 답할 수 있는가?

**평가: 예. Summary 신호만으로 현실적으로 가능하다.**

- Summary는 **6개(또는 control plane 미사용 시 4~5개)** 로 제한되어 있어, 한 화면에 모아 두면 운영자가 “전부 정상 = 안전, 하나라도 비정상 = 조사 필요”로 해석하기에 부담이 적다.
- Experiment 결과, **kube-state-metrics + node_exporter**만 있어도 **S1(NotReady 노드), S4(Pending 파드), S5(과도한 재시작)** 는 즉시 사용 가능하다. 이 세 가지만으로도 “노드 손실 여부”, “스케줄 대기 파드 유무”, “워크로드 반복 실패 여부”를 볼 수 있어, **5–10분 안에 “대략 지금 클러스터가 안전한가?”**에는 답할 수 있다.
- S2(API server), S3(scheduler)는 control plane 스크래핑이 되는 환경에서만 사용 가능하다. Managed K8s에서는 S4 Pending count가 “스케줄·워크로드 관점 건강”을 대표할 수 있어, **control plane metric이 없어도 Summary 4~5개로 “지금 건강한가?” 판단은 가능**하다.
- **결론:** 제안된 Summary 신호 집합은 5–10분 내 운영 판단에 **실용적**이다. 다만 control plane 미사용 환경에서는 “제어면 건강”은 보이지 않으므로, 문서에 **Managed K8s용 요약 구성(예: S1, S4, S5, S6)** 을 명시해 두는 것이 좋다.

### 1.2 Trend / Risk 신호가 조기 불안정 징후 파악에 도움이 되는가?

**평가: 예. 소수 신호만으로도 조기 리스크 인지에 유효하다.**

- **T1(노드 CPU/메모리 사용률), T4(노드 디스크)** 는 대부분의 환경에서 즉시 사용 가능하다. 노드가 80% 근접·디스크 여유 부족은 eviction·성능 저하·장애로 이어질 수 있어, **“곧 불건강해질 수 있다”**는 징후로 적합하다.
- **T2(API server inflight), T3(scheduler pending 추세)** 는 control plane 스크래핑이 될 때만 사용 가능하다. 있으면 제어면 포화·스케줄 대기 증가 추세를 미리 볼 수 있어 유용하다.
- Trend 레이어는 **threshold 근접 + 추세**로만 판단하도록 되어 있어, “지금은 괜찮지만 곧 나빠질 수 있다”는 메시지를 주기에 적절하다. 신호 수도 4~5개 수준으로 적당하다.
- **결론:** Trend/Risk 신호는 **조기 불안정 징후 파악에 실질적으로 도움**이 된다. Control plane metric이 없어도 T1, T4만으로 노드·용량 리스크는 커버 가능하다.

### 1.3 Top Offender 뷰가 조사 대상 좁히기에 실용적인가?

**평가: 예. TOP10 스타일 뷰로 “어디를 먼저 볼 것인가?”에 답하기에 충분하다.**

- **O1(CPU TOP10 노드), O2(메모리 TOP10 노드)** 는 node_exporter만 있으면 사용 가능하다. Summary/Trend에서 “노드 부하·용량 리스크”가 나왔을 때 **어느 노드가 가장 부담을 주는지** 바로 좁힐 수 있다.
- **O3(재시작 TOP10), O4(Pending by workload)** 는 kube-state-metrics로 가능하다. “과도한 재시작”이나 “Pending 있음”이 Summary에서 나오면 **어떤 워크로드/네임스페이스**가 문제인지 확인하는 데 적합하다. pod/namespace 단위만 되어도 5–10분 내 조사 대상 후보를 잡는 데는 무리가 없다.
- O5(에러율 높은 ingress), O6(제어면 컴포넌트)는 환경에 따라 불가할 수 있으나, **즉시 사용 가능한 4개 뷰(O1–O4)** 만으로도 노드·워크로드 관점의 조사 좁히기는 실용적이다.
- **결론:** Top Offenders 뷰는 **조사 대상 좁히기에 실용적**이다. 신호 수를 늘리지 않고도 운영 워크플로에 맞게 동작한다.

### 1.4 신호 집합이 작아서 대시보드 과부하는 없는가?

**평가: 예. 의도적으로 소수로 유지되어 과부하 위험이 낮다.**

- Cluster Health Summary는 **5~7개(실제로는 환경에 따라 4~6개 사용)** 로, “한눈에 본다”는 목표에 맞게 적게 유지되었다.
- Trend는 4~5개, Top Offenders는 5~6개 뷰(각 뷰당 TOP10 한 테이블/그래프) 수준이라, **한 대시보드에 Summary + Trend + Top Offenders를 모두 넣어도 패널 수가 폭발하지 않는다.** 기존 “수십 개 패널을 훑어야 하는” 문제를 피하는 설계가 반영되어 있다.
- **결론:** 신호 집합은 **작게 유지되어 대시보드 과부하를 피하기에 적절**하다.

### 1.5 일반적인 모니터링 노이즈를 적절히 처리하는가?

**평가: 예. 노이즈가 많은 신호는 Summary에서 제외·구분되어 있다.**

- Architecture와 Core signal list에서 **readiness probe 실패, Calico/CNI 등 인프라 probe, 일시적 파드 재시작 1~2회** 는 Summary에 넣지 않고, drill-down·참고용으로만 두도록 명시되어 있다.
- **S5 Excessive restarts** 는 “10분 내 N회 초과” 같은 **임계치**로 정의해, 배포·드레인으로 인한 1~2회 재시작은 Summary에서 비정상으로 잡히지 않도록 설계되어 있다. 즉 **노이즈에 가까운 신호는 걸러지고, 실행 가능한 신호만 Summary에 노출**된다.
- **결론:** 모델은 **일반적인 모니터링 노이즈를 올바르게 다룬다.** Summary를 노이즈에 맡기지 않는 선택이 적절하다.

### 1.6 Control plane metric이 없는 환경에서도 모델이 동작하는가?

**평가: 예. Summary·Trend·Top Offenders 모두 축소 구성으로 동작 가능하다.**

- Experiment notes에 따르면, **Managed K8s 등에서 API server·scheduler를 스크래핑하지 못할 때** S2, S3, T2, T3, O6은 사용할 수 없다. 이 경우 **S4 Workload Pending pod count** 가 “스케줄되지 못한 워크로드가 있는가?”를 대표할 수 있어, control plane metric 없이도 “스케줄·워크로드 관점 건강”을 판단할 수 있다.
- Summary는 **S1, S4, S5(재시작), S6(수정 시)** 만으로 4개, Trend는 **T1, T4** 만으로 2개, Top Offenders는 **O1–O4** 로 4개 뷰를 유지할 수 있다. **“지금 건강한가?”, “곧 불건강해질 징후는?”, “어디를 먼저 볼 것인가?”** 세 질문에 모두 답하는 최소 집합이 유지된다.
- **결론:** **Control plane metric이 없어도 모델은 동작한다.** 제어면 가시성은 줄어들지만, 노드·워크로드·엔드포인트 관점의 운영 판단은 가능하다. 다만 “control plane 미사용 시 권장 구성”을 문서에 명시해 두면 운영자가 환경별로 적용하기 쉽다.

---

## 2. 실제 운영 워크플로 관점 평가

### 2.1 일상 점검(5–10분)

- **워크플로:** 아침 또는 정해진 시간에 “클러스터가 안전한가?”만 빠르게 확인.
- **모델 부합:** Summary 4~6개를 한 블록으로 두고, 전부 정상이면 “안전”, 하나라도 비정상이면 “조사 필요”로 해석하면 된다. 추가로 Trend에서 경고가 있으면 “용량·디스크 등 대비 필요”로 이어질 수 있어, **일상 점검 워크플로와 잘 맞는다.**

### 2.2 이상 징후 시 조사

- **워크플로:** Summary 또는 알림에서 이상이 감지된 뒤, “어디가 원인인가?”를 좁혀서 조사.
- **모델 부합:** Summary에서 **어느 카테고리**가 비정상인지(노드 / 워크로드 / 엔드포인트 등) 파악한 다음, 해당 카테고리에 대응하는 Top Offenders 뷰(O1–O4 등)로 들어가 TOP10을 보면 된다. **“비정상 신호 → 카테고리 → TOP10으로 후보 좁히기”** 흐름이 명확해 운영 워크플로와 일치한다.

### 2.3 용량·리스크 사전 대응

- **워크플로:** “아직 문제는 없지만 노드/디스크가 빡빡해지고 있다”는 것을 미리 알고 대응.
- **모델 부합:** Trend 레이어(T1 노드 사용률, T4 디스크)가 이 역할을 한다. Summary가 전부 정상이어도 Trend에서 경고가 나오면 “용량 확장·정리 검토” 등 사전 대응이 가능하다. **실제 운영에서 필요한 “조기 리스크 인지”와 부합한다.**

---

## 3. 강점 (Strengths)

1. **세 가지 운영 질문에 1:1 대응:** “지금 건강한가?”(Summary), “곧 불건강해질 징후는?”(Trend), “어디를 먼저 볼 것인가?”(Top Offenders)가 레이어별로 명확히 나뉘어 있어, 운영자가 **어디를 보면 어떤 질문에 답하는지** 바로 이해할 수 있다.
2. **Summary 신호 수 최소화:** 5~7개(환경에 따라 4~6개)로 제한되어 **한눈에 판단**하기 좋고, 대시보드 과부하를 피할 수 있다.
3. **노이즈 처리 명시:** readiness probe, 일시적 재시작 등은 Summary에서 제외하고, “과도한 재시작”만 임계치로 두어 **signal-to-noise** 를 유지했다.
4. **환경 제약 반영:** Experiment에서 control plane·Eviction·Ingress 등 **metric 부재·조건부 사용**을 정리했고, **control plane 미사용 시에도 S4 등으로 대체해 모델이 동작**하도록 되어 있다.
5. **실제 사용 가능한 신호 위주:** Core signal list와 Experiment가 **즉시 사용 가능 / 수정 후 사용 / 현재 불가**를 구분해 두어, **첫 번째 대시보드 반복을 현재 Prometheus 환경으로 구축할 수 있다**는 결론이 나와 있다.

---

## 4. 약점 및 리스크 (Weaknesses / Risks)

1. **Control plane 가시성 부재(Managed K8s):** API server·scheduler metric을 쓸 수 없으면, “제어면이 지금 응답·스케줄하는가?”는 직접 보이지 않는다. S4 Pending으로 워크로드 관점은 대체할 수 있으나, **제어면 장애 자체**는 다른 채널(managed 서비스 알림, 지원 채널 등)에 의존할 수밖에 없다. 리스크: 제어면 장애 시 Summary만 보고는 “정상”으로 보일 수 있음 — 문서에 **“control plane은 별도 채널 확인”** 같은 주의를 두는 것이 좋다.
2. **Eviction 미반영:** Eviction에 대한 Prometheus metric이 흔히 없어 v1에서는 Summary에 넣지 않았다. **메모리 압박으로 인한 eviction** 은 T1(노드 메모리)·O2(메모리 TOP10)로 간접적으로만 보인다. Eviction 자체를 숫자로 보려면 후속에서 이벤트·다른 수단이 필요하다.
3. **S6, O3, O4의 환경 의존성:** S6(엔드포인트 비어 있음)은 metric 이름·라벨에 따라 PromQL 수정이 필요하다. O3/O4는 deployment 등 상위 단위로 보려면 조인·라벨이 필요하다. **문서에 “환경별 수정 포인트”** 가 정리되어 있으면, 대시보드 설계 단계에서 빠뜨리기 쉬운 부분을 줄일 수 있다.
4. **Trend “추세” 정의 부족:** T2, T3는 “상승 추세”로 해석하도록 되어 있으나, **어떤 time window, 어떤 증가율/절대값을 쓰면 “추세 경고”로 볼지** 가 core signal list·architecture에 구체적으로 적혀 있지 않다. 구현 시 팀마다 다르게 해석할 수 있어, **가이드 한 줄(예: “최근 15분 증가량 > N” 등)** 이 있으면 좋다.

---

## 5. 소규모 개선 제안 (신호 확대 없이)

1. **Managed K8s용 “권장 Summary 구성” 명시:** Architecture 또는 Core signal list에 **“Control plane metric을 사용할 수 없는 환경에서는 Summary를 S1, S4, S5(재시작), S6(가능 시)로 구성한다”** 는 문단을 추가한다. 동일 문서 또는 Experiment에 **“이 경우 제어면 건강은 managed 서비스 알림·콘솔 등 별도 채널로 확인한다”** 는 주의를 한 줄 넣는다.
2. **S6·O3·O4 “환경별 수정 포인트” 정리:** Core signal list 또는 Experiment에 **S6(endpoint metric 이름·라벨 확인), O3(deployment 단위 시 owner/라벨 조인), O4(세분화 시 동일)** 에 대해 “대시보드 구현 시 확인할 것” 체크리스트를 짧게 추가한다.
3. **Trend “추세” 해석 가이드 한 줄:** Trend/Risk 섹션에 **“추세 경고 예: scheduler_pending_pods 또는 Pending count의 최근 15분 증가량이 N 초과”** 같은 예시를 한 줄 넣어, 구현 시 일관된 해석이 나오도록 한다.
4. **Review 결론 반영:** 아래 “최종 결론”을 Architecture 또는 Core signal list 마지막에 “Review 결과”로 요약해 두면, 이후 dashboard/alert/runbook 프로젝트에서 **v1 범위와 제한**을 참고하기 쉽다.

위 개선은 **신호 개수를 늘리지 않고** 문서만 보완하는 수준으로, 모델 구조나 레이어 설계를 바꾸지 않는다.

---

## 6. 최종 결론

**평가 요약**

- 운영자가 **5–10분 안에 “지금 클러스터가 건강한가?”**에 답하는 것은 **Summary 신호만으로 현실적으로 가능**하다.
- **Trend/Risk** 신호는 **조기 불안정 징후** 파악에 도움이 되며, control plane이 없어도 T1, T4로 노드·용량 리스크는 커버 가능하다.
- **Top Offenders** 뷰는 **조사 대상 좁히기**에 실용적이며, O1–O4만으로도 노드·워크로드 관점 조사가 가능하다.
- 신호 집합은 **작게 유지**되어 대시보드 과부하를 피하기에 적절하고, **노이즈가 많은 신호는 Summary에서 제외**되어 있다.
- **Control plane metric이 없는 환경**에서도 **S4 등으로 대체한 축소 구성**으로 모델이 동작한다.

**결론: 모니터링 모델은 v1으로 수용 가능하다. 다만 소규모 조정을 권장한다.**

- **The monitoring model is acceptable for v1.**
- **Minor adjustments** 로, 다음을 권장한다:
  - Managed K8s용 Summary 구성(S1, S4, S5, S6) 및 “제어면은 별도 채널 확인” 문구 추가.
  - S6, O3, O4에 대한 환경별 수정 포인트(체크리스트) 추가.
  - Trend “추세” 해석 예시 한 줄 추가.
  - Review 결론을 Architecture/Core signal list에 “Review 결과”로 요약 반영.

**Major redesign은 필요하지 않다.** 레이어 구조, 카테고리, 신호 수, 노이즈 처리, 제한 환경 대응이 운영 유용성과 맞아 떨어지며, Experiment에서도 “첫 번째 대시보드 반복 가능”이 확인되었다. 문서만 위와 같이 보완하면 **sre-monitoring-dashboard-design**, **sre-monitoring-alert-policy**, **sre-monitoring-operational-runbooks** 단계로 넘어가기에 충분하다.

---

*Review phase output. 다음 단계: Documentation — 최종 문서 “Kubernetes Cluster Health Monitoring Model” 작성.*
