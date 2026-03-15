# Research: ClusterLoader2 부하 테스트와 SLI/SLO

**Project:** sample-cl2-sli-slo-analysis  
**Role:** Researcher  
**Input:** 01-manager/project-brief.md, knowledge/principles (cat-vision, cat-design-principles, sli-slo-philosophy, devcat-program-brief)

**참고:** devcat 저장소는 본 AI 팀 워크스페이스와 별도이므로 `devcat/scenarios/load/` 및 `modules/` 디렉터리를 직접 읽지 못했다. 아래는 **Kubernetes 공식 perf-tests 리포의 ClusterLoader2**(testing/load, modules/measurements.yaml 등)를 기준으로 한 분석이다. devcat을 사용할 때는 동일한 모듈·메트릭 구조가 복사되어 있을 가능성이 높으므로, devcat 실험 시 이 목록과 매핑해 검증하는 것을 권장한다.

---

## 1. ClusterLoader2가 산출하는 메트릭 (Q1)

ClusterLoader2는 **measurements** 모듈을 통해 다음 유형의 메트릭을 수집·산출한다. (출처: [kubernetes/perf-tests clusterloader2/testing/load/modules/measurements.yaml](https://github.com/kubernetes/perf-tests/blob/master/clusterloader2/testing/load/modules/measurements.yaml), [pkg/measurement](https://github.com/kubernetes/perf-tests/tree/master/clusterloader2/pkg/measurement).)

### 1.1 메트릭 목록

| 메트릭(Identifier) | 의미 | 비고 |
|--------------------|------|------|
| **APIResponsivenessPrometheus** | API 호출별 레이턴시 (Prometheus 기반) | 슬로우 콜·커스텀 임계값 가능 |
| **APIResponsivenessPrometheusSimple** | 단순화된 API 레이턴시 쿼리 | 위와 동일 계열, simple 쿼리 |
| **CreatePhasePodStartupLatency** | Pod 기동 레이턴시 (create phase) | labelSelector=load, **기본 threshold 1h** (대규모 saturation용) |
| **InClusterNetworkLatency** | 클러스터 내 네트워크 레이턴시 (ping 등) | probe 기반 |
| **NodeLocalDNSLatency** | Node-local DNS 레이턴시 | (옵션) |
| **SLOMeasurement** | SIG-Scalability SLO 측정 통합 | probe 기반 |
| **NetworkProgrammingLatency** | 네트워크 프로그래밍 레이턴시 (e.g. iptables 반영) | 기본 threshold 30s 등 |
| **KubeProxyIptablesRestoreFailures** | kube-proxy iptables partial restore 실패 횟수 | threshold 0 (실패 없어야 함) |
| **APIAvailability** | API 가용성 비율 | (옵션) threshold 예: 99.5% |
| **ContainerRestarts** | 컨테이너 재시작 횟수 | 허용 재시작 수 설정 가능 |
| **ContainerCPU / ContainerMemory** | 컨테이너 CPU/메모리 사용 (P50/P90/P99) | (옵션) Prometheus |
| **TestMetrics** | 테스트 단계 타이머, OOM, 시스템 Pod 메트릭 등 | 단계별 소요 시간·리소스 |
| **ResourceSize** | API 서버 리소스 크기 추정 (MiB) | apiserver 메트릭 |
| **DNS 성능** | modules/dns-performance-metrics.yaml | DNS 관련 지표 |
| **Pod startup (전용 모듈)** | modules/pod-startup-latency.yaml | **소규모 클러스터(<100노드)에서는 기본 config에서 스킵됨** |
| **Scheduler throughput** | 스케줄러 처리량 | **소규모 클러스터에서는 스킵됨** |

- 공식 load config는 **100+ 노드**를 가정하며, `$IS_SMALL_CLUSTER := lt .Nodes 100` 일 때 scheduler-throughput, pod-startup-latency 모듈이 비활성화된다. 소규모 클러스터 수용 테스트에서는 “어떤 메트릭을 실제로 켜고 해석할지”를 선택해야 한다.

---

## 2. 수용 테스트에 의미 있는 메트릭 (Q2, Q3)

### 2.1 사용자·워크로드 관점에서 “의미 있다”의 기준

Kubernetes scalability 철학을 참고한다([kubernetes/community sig-scalability/slos](https://github.com/kubernetes/community/tree/master/contributors/devel/sig-scalability/slos)):

- **Precise and well-defined:** 지표와 측정 방법이 모호하지 않다.
- **Consistent:** 동일 조건에서 재현 가능한 방식으로 측정된다.
- **User-oriented:** 워크로드·개발자가 체감하는 동작(파드 기동, API 응답, 네트워크)을 반영한다.
- **Testable:** 수용 테스트에서 반복 측정·판단 가능하다.

“You promise, we promise”에 따르면, **SLO는 플랫폼이 “이 조건에서 이 수준을 보장한다”고 약속하는 것**이고, **수용 테스트는 그 약속이 지켜지는지 검증**한다. 따라서 수용 테스트에 의미 있는 메트릭은 “사용자·워크로드가 기대하는 동작”과 연결될 수 있어야 한다.

### 2.2 수용 테스트 관점에서 의미 있는 메트릭과 이유

| 메트릭 | 수용 테스트에서 의미 있는 이유 (사용자·워크로드 관점) |
|--------|--------------------------------------------------------|
| **PodStartupLatency** | 파드가 “언제 준비되는가”는 배포·스케일링 체감 품질과 직결. 너무 길면 배포/롤아웃이 불안정하게 느껴짐. |
| **APIResponsiveness (Prometheus)** | API 호출 지연은 kubectl·컨트롤러·오퍼레이터 동작에 직접 영향. 슬로우 호출이 많으면 사용자 작업이 멈춘 것처럼 보일 수 있음. |
| **InClusterNetworkLatency** | 파드 간 통신 지연은 마이크로서비스·데이터 플로우에 영향. 수용 시 “클러스터 내 네트워크가 정상 범위인가”를 볼 수 있음. |
| **NetworkProgrammingLatency** | 서비스/엔드포인트 반영 지연은 트래픽이 올바른 파드로 가는 시점에 영향. 사용자 관점에서는 “서비스가 금방 반영되는가”와 연결. |
| **APIAvailability** | API가 일정 비율 이상 응답해야 클러스터를 “사용 가능”이라 할 수 있음. 가용성은 사용자·워크로드의 기본 기대. |
| **ContainerRestarts / TestMetrics (OOM 등)** | 비정상 재시작·OOM은 워크로드 안정성과 직결. “이 클러스터에서 워크로드가 안정적으로 돌아가는가” 수용 판단에 필요. |
| **KubeProxyIptablesRestoreFailures** | 0이어야 정상. 실패가 있으면 트래픽 라우팅이 깨질 수 있어 사용자 관점에서 결함 지표. |

- **ResourceSize, ContainerCPU/Memory** 등은 “리소스 사용이 정상 범위인가”를 보는 보조 지표로 의미 있으나, 수용 테스트의 **주요** 통과 기준으로 쓰려면 “어떤 워크로드·어떤 한도”인지 환경 가정과 함께 정당화하는 것이 좋다.

### 2.3 SLI 후보 목록

위 “의미 있는 메트릭”을 **SLI(서비스 수준 지표) 후보**로 정리하면 다음과 같다. 실제 SLO로 쓸 때는 “정의·측정 방법·단위·조건”을 precise and well-defined 하게 고정해야 한다.

| SLI 후보 | 측정 대상 | 단위/형태 예 | 비고 |
|----------|-----------|--------------|------|
| **Pod startup latency** | 파드 생성부터 준비까지 시간 | P50/P90/P99 (초) | SIG-Scalability pod_startup_latency SLO(예: 5s)와 구분해 소규모용 정당화 필요 |
| **API call latency** | API 호출별 레이턴시 | verb/resource별 P99 등 | 슬로우 콜 개수·임계값과 함께 정의 |
| **In-cluster network latency** | 파드 간 RTT 등 | P99 (ms 또는 초) | probe 기반, 조건(노드 수·프로브 수) 명시 |
| **Network programming latency** | 서비스/엔드포인트 반영 시간 | P99 (초) | 30s 등 대규모용 임계값과 구분 |
| **API availability** | API 성공 비율 | % (시간 구간) | 폴링 주기·대상 명시 |
| **Container restarts** | 테스트 구간 중 재시작 횟수 | 횟수 (허용 상한) | 0 또는 팀 정책에 따른 허용치 |
| **Kube-proxy iptables restore failures** | partial restore 실패 | 횟수 | 0 기대 |

- **TestMetrics**의 단계 타이머(create/scale/delete 소요 시간)는 “시나리오 완료 시간” SLI 후보로 둘 수 있다. “N 파드 생성이 T 초 이내에 끝나야 수용”처럼 **실용적 수용 의미**(practical acceptance meaning)를 부여하면 된다.

---

## 3. 벤치마크 스타일 임계값의 한계 (Q5 관련)

### 3.1 예: 1시간(1h) PodStartupLatency threshold

공식 measurements.yaml에는 **CreatePhasePodStartupLatency**에 `threshold: 1h`가 설정되어 있다. 주석에 “Ideally this should be 5s”라고 되어 있으나, 대규모 saturation 테스트(수천 파드)에서는 “전부 기동될 때까지” 기다리는 용도로 1h를 둔 것으로 보인다.

- **한계:** 1h는 **수용 테스트의 “이 클러스터를 수용한다”는 판단 기준**으로는 의미가 거의 없다. 거의 모든 클러스터가 1h 안에 파드를 기동할 수 있어, 수용/거부를 구분하지 못한다.
- **소규모 클러스터:** 노드·파드 수가 적을 때는 “P99 pod startup ≤ 5s” 같은 **짧은 구간**이 사용자 기대와 맞다. 대규모 벤치마크용 1h를 그대로 쓰면 안 된다.

### 3.2 일반적 한계

- **대규모·고부하용 임계값:** 5000노드·수만 파드 규모의 벤치마크에서 쓰는 latency/throughput 임계값은, 소규모 클러스터의 “일상 부하에서 기대하는 수준”과 다르다. 그대로 재사용하면 과도하게 느슨하거나(통과만 함) 또는 오히려 과도하게 엄격할 수 있다.
- **일관성·재현성:** 도구 기본값이 “어떤 시나리오·어떤 환경”에서 정의되었는지 문서화되어 있지 않으면, consistent·testable 원칙에 어긋난다. 수용 테스트용 SLO는 **우리 환경·시나리오**를 전제로 새로 정당화하는 것이 안전하다.

---

## 4. 소규모 클러스터에서 SLO를 정할 때 고려 사항 (Q4, Q6)

### 4.1 Kubernetes scalability 철학 반영

- **Precise and well-defined:** SLO마다 “무엇을, 어떤 조건에서, 어떤 단위로 측정하는지”를 명시. ClusterLoader2의 경우 시나리오 이름·파라미터(노드 수, 파드 수, QPS 등)·measurement identifier를 함께 기록.
- **Consistent:** 동일 config·동일 클러스터 스펙에서 반복 실행 시 같은 방식으로 측정. devcat에서는 config.yaml·ol-test.yaml 등 오버라이드를 버전 관리해 재현 가능하게.
- **User-oriented:** “이 SLO를 만족하면 사용자·워크로드가 기대하는 동작을 한다”는 설명을 둠. 예: “P99 pod startup ≤ 5s → 일반적인 배포/스케일링이 체감상 부담 없이 완료된다.”
- **Testable:** ClusterLoader2(및 perfdash)가 실제로 산출하는 메트릭·포맷으로 측정 가능해야 함. 산출이 없거나 불안정한 지표는 SLO 후보에서 제외하거나, 측정 방법을 먼저 확립.

### 4.2 “You promise, we promise”

- **You promise:** “이 노드 수·이 워크로드 패턴·이 시나리오에서, P99 pod startup은 5초 이하로 보장한다”처럼 **팀/플랫폼이 명시적으로 약속**하는 값이 SLO다.
- **We promise:** 수용 테스트(ClusterLoader2 실행 + 결과 해석)가 그 약속이 지켜졌는지 **검증**한다. 따라서 SLO 값은 “도구가 주는 기본값”이 아니라 **user expectations, environment assumptions, repeatable measurement, practical acceptance meaning**으로 정당화되어야 한다.

### 4.3 SLO 정당화 시 체크 (sli-slo-philosophy 정리)

| 기준 | 내용 |
|------|------|
| **User expectations** | 사용자(워크로드 소유자, 플랫폼 이용자)가 기대하는 응답 시간·가용성·처리량. “5초 안에 파드가 준비되면 만족” 등. |
| **Environment assumptions** | 노드 수, 리소스, 네트워크, 시나리오(파드 수, QPS 등). “N 노드, M 파드 create 시나리오에서”를 명시. |
| **Repeatable measurement** | 동일 조건에서 반복 측정 가능, 측정 방법·도구 버전이 문서화됨. ClusterLoader2 report-dir, perfdash 포맷 등. |
| **Practical acceptance meaning** | “이 SLO를 만족하면 이 클러스터를 수용한다”가 팀·운영에서 말이 되도록 정의. 1h 같은 값은 practical meaning이 없음. |

- 소규모 클러스터에서는 **“우리 클러스터는 N 노드, 이런 워크로드를 받을 때 이 수준을 만족해야 수용”**으로 구체화하고, 대규모 벤치마크의 임계값을 그대로 가져오지 않는다.

---

## 5. devcat 실험으로 SLO 후보 검증·정제 (Q5)

- **실험 흐름:** devcat에서 ClusterLoader2를 실행(config.yaml, 클러스터별 ol-test.yaml) → results/에 저장된 결과(및 perfdash 입력)에서 위 SLI 후보 메트릭을 추출 → 실제 숫자 분포(P50/P90/P99 등)를 확인 → “현재 클러스터가 어떤 수준인가”를 파악한 뒤, 그에 맞춰 SLO 후보 값을 설정하거나 조정.
- **검증:** 동일 config로 여러 번 돌려서 **일관성(consistent)** 확인. 환경(노드 수, 리소스)을 바꾼 실험으로 **environment assumptions**이 SLO에 미치는 영향을 관찰.
- **정제:** “P99 5s로 두었는데 실제로는 항상 2s 근처다”면 5s는 유지해도 되고, “실제로는 10s가 나온다”면 5s SLO는 환경 개선 또는 SLO 완화(정당화 문서 업데이트) 중 하나가 필요. 실험 결과가 **user expectations, practical acceptance meaning**과 맞는지 팀에서 해석하는 단계가 필요하다.
- **한계:** 본 워크스페이스에서는 devcat 저장소의 scenarios/load·modules를 직접 읽지 못했다. devcat이 perf-tests의 load 시나리오·modules를 복사해 쓴다면 위 메트릭 목록이 그대로 적용 가능하다. devcat 전용 시나리오가 있다면, 해당 시나리오에서 어떤 measurement가 수집되는지 확인한 뒤 위 SLI 후보와 매핑하는 작업을 실험 전에 수행하는 것이 좋다.

---

## 6. 요약 (Architect·Engineer용)

- **ClusterLoader2 메트릭:** APIResponsivenessPrometheus, PodStartupLatency, InClusterNetworkLatency, NetworkProgrammingLatency, APIAvailability, ContainerRestarts, TestMetrics, ResourceSize, DNS·kube-proxy 등. 공식 load config는 100+ 노드 가정, 소규모 클러스터에서는 scheduler-throughput·pod-startup-latency 모듈이 비활성화됨.
- **수용 테스트에 의미 있는 것:** Pod startup, API 레이턴시, in-cluster network, network programming, API 가용성, 재시작/OOM, iptables 실패. 사용자·워크로드 관점(배포 체감, API 반응, 네트워크·가용성·안정성)과 연결해 정당화.
- **SLI 후보:** Pod startup latency, API call latency, In-cluster network latency, Network programming latency, API availability, Container restarts, Kube-proxy iptables failures, (선택) 단계 타이머. 각각 precise·consistent·user-oriented·testable 하게 정의 필요.
- **벤치마크 임계값 한계:** 1h 같은 큰 threshold는 수용 판단에 무의미. 소규모 클러스터는 “짧은 구간(예: 5s)”과 “우리 환경·시나리오”에 맞는 SLO를 별도 정당화.
- **소규모 클러스터 SLO 고려 사항:** scalability 철학 4원칙 + you promise we promise, user expectations·environment assumptions·repeatable measurement·practical acceptance meaning으로 정당화. 대규모 벤치마크 임계값 무비판 재사용 금지.
- **devcat 실험:** 동일 config 반복 실행·결과 해석으로 SLO 후보 값 검증·정제. devcat의 scenarios/load·modules는 별도 리포이므로 실험 시 위 메트릭 목록과 매핑해 사용할 것.

---

## 7. 오픈 질문·후속

- devcat 저장소의 `scenarios/load/` 및 `modules/`를 실제로 열어, 위 메트릭 중 어떤 것이 수집되는지·threshold 설정이 어떻게 되어 있는지 확인하면 좋다.
- PodStartupLatency를 소규모 클러스터에서 쓰려면, 공식 load config에서 스킵되는 pod-startup-latency 모듈을 대체하는 시나리오 또는 오버라이드가 필요할 수 있다. Architect 단계에서 “소규모용 최소 시나리오·measurement 세트”를 설계할 때 반영할 수 있다.
- SLO 값(예: 5s, 99%)을 **숫자로 고정**하는 것은 Engineer·실험 단계에서 “정당화 템플릿”을 채운 뒤, devcat 실험으로 검증·정제하는 흐름을 권장한다.
