# SLI 후보 정의: 소규모 클러스터 수용 테스트용

**Project:** sample-cl2-sli-slo-analysis  
**Role:** Experimenter  
**Input:** 02-research, 03-architecture, 04-engineering, 05-experiment/experiment-notes.md, CAT 원칙·devcat-program-brief

이 문서는 **첫 CAT SLI 후보 세트**를 소규모 클러스터 수용 테스트 관점에서 정의한다. Research·Architecture·Engineering 산출물과 실제 ClusterLoader2 실험에서 확인된 산출물을 바탕으로, **측정 가능한 것**과 **Prometheus 의존·불확실** 항목을 구분해 Review 단계에서 “첫 CAT SLI 세트” 및 “선택/실험용” 결정에 활용할 수 있도록 한다.

---

## 1. 정당화 원칙 (CAT·Research 반영)

- **Precise and well-defined:** 각 SLI는 “무엇을, 어떤 시나리오·조건에서, 어떤 단위로” 측정하는지 명시한다.
- **User-oriented:** 수용 테스트에서 “이 지표가 나쁘면 사용자·워크로드가 무엇을 체감하는가”를 연결할 수 있어야 한다.
- **Testable:** ClusterLoader2(및 Prometheus)가 실제로 산출하는 파일·메트릭으로 측정 가능해야 한다.
- **Practical acceptance meaning:** “이 SLO를 만족하면 이 클러스터를 수용한다”가 팀·운영에서 말이 되도록, 소규모 클러스터·일반 부하에 맞는 임계값을 정당화한다. 벤치마크 스타일의 매우 큰 임계값(예: 1h)은 수용 판단에 사용하지 않는다.

---

## 2. SLI 후보별 정의

### 2.1 Pod startup latency

| 항목 | 내용 |
|------|------|
| **정의** | 파드 생성부터 Ready까지 소요 시간. |
| **metric source** | ClusterLoader2 measurement: **CreatePhasePodStartupLatency** (PodStartupLatency, StatelessPodStartupLatency, StatefulPodStartupLatency). |
| **measurement location** | `PodStartupLatency_CreatePhasePodStartupLatency_load_*.json`, `StatelessPodStartupLatency_*_*.json`, `StatefulPodStartupLatency_*_*.json`. 내부 `dataItems[].data.Perc50|Perc90|Perc99`, `labels.Metric` (예: `pod_startup`, `schedule_to_watch`). |
| **unit** | ms (또는 SLO 정의 시 초로 환산). |
| **수용 테스트에서의 의미** | 배포·스케일링 시 “파드가 언제 쓸 수 있게 되는가”를 반영. 값이 크면 롤아웃·스케일이 느리게 체감됨. |
| **user/workload relevance** | **워크로드 직접 체감.** 파드 기동이 느리면 배포·스케일 경험이 나쁨. |

### 2.2 API server request latency (API call latency)

| 항목 | 내용 |
|------|------|
| **정의** | API 서버에 대한 요청별 레이턴시(verb/resource별). |
| **metric source** | ClusterLoader2 measurement: **APIResponsivenessPrometheus**, **APIResponsivenessPrometheusSimple**. (Prometheus 기반.) |
| **measurement location** | `APIResponsivenessPrometheus_*_load_*.json`, `APIResponsivenessPrometheus_simple_*_*.json` (Prometheus 활성화·gather 완료 시 생성). Prometheus 쿼리 결과에서 P99·슬로우 콜 수 등. |
| **unit** | ms 또는 초; 슬로우 콜은 횟수. |
| **수용 테스트에서의 의미** | kubectl·컨트롤러·오퍼레이터가 API를 쓸 때 응답 지연. 슬로우 콜이 많으면 작업이 멈춘 것처럼 보일 수 있음. |
| **user/workload relevance** | **사용자·워크로드 모두.** API 호출 주체(kubectl, 컨트롤러)의 체감 반응성. |

### 2.3 API availability

| 항목 | 내용 |
|------|------|
| **정의** | 측정 구간 동안 API 서버 응답 성공 비율. |
| **metric source** | ClusterLoader2 measurement: **APIAvailability** (옵션 활성화 시). |
| **measurement location** | APIAvailability 관련 산출 파일(옵션·시나리오에 따라 파일명·경로 확인 필요). 비율(%) 형태. |
| **unit** | % (예: 99.5). |
| **수용 테스트에서의 의미** | 클러스터가 “사용 가능”인지의 기본 지표. 일정 비율 미만이면 수용 불가. |
| **user/workload relevance** | **사용자·워크로드 공통.** 가용성은 기본 기대. |

### 2.4 Container restarts / OOM events

| 항목 | 내용 |
|------|------|
| **정의** | 테스트 구간 중 컨테이너 재시작 횟수 및 OOM 발생 여부. |
| **metric source** | ClusterLoader2: **ClusterOOMsTracker**, **TestMetrics**(시스템 파드 메트릭)·**SystemPodMetrics**. |
| **measurement location** | `ClusterOOMsTracker_load_*.json` (`failures`, `ignored`, `past`), `SystemPodMetrics_load_*.json` (`pods[].containers[].restartCount`, `lastRestartReason`). |
| **unit** | 횟수; OOM은 발생 여부(목록). |
| **수용 테스트에서의 의미** | 비정상 재시작·OOM은 워크로드 안정성과 직결. “이 클러스터에서 워크로드가 안정적으로 돌아가는가” 수용 판단에 필요. |
| **user/workload relevance** | **워크로드 직접.** 재시작·OOM은 서비스 품질·가용성에 직결. |

### 2.5 System pod health

| 항목 | 내용 |
|------|------|
| **정의** | 시스템 컴포넌트(coredns, kube-apiserver, kube-proxy, kube-scheduler 등) 파드의 재시작 횟수·이유. |
| **metric source** | ClusterLoader2: **SystemPodMetrics** (TestMetrics 계열). |
| **measurement location** | `SystemPodMetrics_load_*.json`. `pods[].containers[].restartCount`, `lastRestartReason`. |
| **unit** | 재시작 횟수(정수). |
| **수용 테스트에서의 의미** | 컨트롤 플레인·시스템 컴포넌트가 안정적인지 판단. 재시작이 많으면 플랫폼 불안정 신호. |
| **user/workload relevance** | **플랫폼 건강도.** 사용자·워크로드는 간접적으로 영향(API·네트워크 품질). |

### 2.6 In-cluster network latency

| 항목 | 내용 |
|------|------|
| **정의** | 클러스터 내 파드 간 네트워크 RTT(또는 유사 지표). |
| **metric source** | ClusterLoader2 measurement: **InClusterNetworkLatency** (Prometheus 기반). |
| **measurement location** | InClusterNetworkLatency 관련 JSON (Prometheus 활성화·gather 완료 시). |
| **unit** | ms 또는 초 (P99 등). |
| **수용 테스트에서의 의미** | “클러스터 내 네트워크가 정상 범위인가” 판단. 마이크로서비스·파드 간 통신 지연 반영. |
| **user/workload relevance** | **워크로드.** 파드 간 통신 체감 지연. |

### 2.7 Network programming latency

| 항목 | 내용 |
|------|------|
| **정의** | 서비스/엔드포인트 변경이 데이터플레인에 반영되는 시간. |
| **metric source** | ClusterLoader2 measurement: **NetworkProgrammingLatency** (Prometheus·kube-proxy 메트릭). |
| **measurement location** | NetworkProgrammingLatency 관련 산출 (Prometheus 활성화 시). |
| **unit** | 초 (P99 등). |
| **수용 테스트에서의 의미** | 트래픽이 새 파드로 가기까지 걸리는 시간. “서비스가 금방 반영되는가”와 연결. |
| **user/workload relevance** | **워크로드.** 스케일/롤아웃 시 트래픽 전환 체감. |

### 2.8 Kube-proxy iptables restore failures

| 항목 | 내용 |
|------|------|
| **정의** | kube-proxy의 iptables partial restore 실패 횟수. |
| **metric source** | ClusterLoader2 measurement: **GenericPrometheusQuery** (KubeProxyIptablesRestoreFailures). Prometheus scrape. |
| **measurement location** | KubeProxy 관련 measurement 산출 파일 (Prometheus 활성화 시). |
| **unit** | 횟수. |
| **수용 테스트에서의 의미** | 0이어야 정상. 실패 시 트래픽 라우팅 오류 가능. 결함 지표. |
| **user/workload relevance** | **워크로드·플랫폼.** 서비스 디스커버리·라우팅 안정성. |

### 2.9 Phase duration (선택)

| 항목 | 내용 |
|------|------|
| **정의** | create/scale/delete 단계별 소요 시간. |
| **metric source** | ClusterLoader2: **TestMetrics** (단계 타이머). |
| **measurement location** | TestMetrics/SchedulingMetrics 산출 또는 `ResourceUsageSummary_load_*.json` 등에서 구간 정보. kind 등에서는 SchedulingMetrics가 실패할 수 있어 불완전. |
| **unit** | 초. |
| **수용 테스트에서의 의미** | “N 파드 생성이 T 초 이내에 끝나는가” 같은 시나리오 완료 시간. 실용적 수용 의미 부여 가능. |
| **user/workload relevance** | **시나리오 완료 체감.** 간접적. |

---

## 3. SLI 후보 요약 테이블

| SLI name | metric source | measurement location | unit | user/workload relevance | notes |
|----------|---------------|----------------------|------|--------------------------|--------|
| **Pod startup latency** | ClusterLoader2 CreatePhasePodStartupLatency | `PodStartupLatency_*_load_*.json`, `StatelessPodStartupLatency_*_*.json`, `StatefulPodStartupLatency_*_*.json` → `dataItems[].data.Perc50|90|99`, `labels.Metric` | ms (또는 s) | 워크로드: 배포·스케일 체감 | 공식 threshold 1h는 수용 부적합; 소규모는 예: P99 ≤ 5s 정당화 필요. |
| **API server request latency** | APIResponsivenessPrometheus, APIResponsivenessPrometheusSimple | `APIResponsivenessPrometheus*_*_load_*.json` (Prometheus 활성화·gather 완료 시) | ms, 슬로우 콜 횟수 | 사용자·워크로드: API 반응성 | **Prometheus 의존.** |
| **API availability** | APIAvailability | APIAvailability 산출 파일 (옵션에 따라) | % | 사용자·워크로드: 가용성 | 옵션 활성화·파일 경로 확인 필요. |
| **Container restarts / OOM** | ClusterOOMsTracker, SystemPodMetrics | `ClusterOOMsTracker_load_*.json`, `SystemPodMetrics_load_*.json` | 횟수, OOM 목록 | 워크로드: 안정성 | |
| **System pod health** | SystemPodMetrics | `SystemPodMetrics_load_*.json` (시스템 파드 restartCount 등) | 재시작 횟수 | 플랫폼: 컨트롤 플레인·시스템 안정성 | Container restarts와 동일 파일, 해석만 구분. |
| **In-cluster network latency** | InClusterNetworkLatency | InClusterNetworkLatency 산출 (Prometheus 활성화 시) | ms 또는 s (P99) | 워크로드: 파드 간 통신 | **Prometheus 의존.** |
| **Network programming latency** | NetworkProgrammingLatency | NetworkProgrammingLatency 산출 (Prometheus 활성화 시) | s (P99) | 워크로드: 서비스 반영 속도 | **Prometheus 의존.** |
| **Kube-proxy iptables restore failures** | GenericPrometheusQuery (KubeProxyIptablesRestoreFailures) | KubeProxy measurement 산출 (Prometheus 활성화 시) | 횟수 | 워크로드·플랫폼: 라우팅 결함 | **Prometheus 의존.** 0 기대. |
| **Phase duration** | TestMetrics (단계 타이머) | TestMetrics/SchedulingMetrics 또는 ResourceUsageSummary | s | 시나리오 완료 체감 (간접) | kind에서 SchedulingMetrics 실패 시 **불확실.** |

---

## 4. 측정 가능성 분류 (실험 기준)

이전 실험(experiment-notes.md) 및 Prometheus measurement readiness precheck(§10–§11)를 반영한 분류다.

### 4.1 이미 실험으로 측정 가능 확인됨 (non-Prometheus)

| SLI | 비고 |
|-----|------|
| **Pod startup latency** | Run `20260315-181750` 등에서 JSON 산출·Perc50/90/99 추출 가능. |
| **Container restarts / OOM** | `ClusterOOMsTracker_*.json`, `SystemPodMetrics_*.json`에서 확인됨. |
| **System pod health** | `SystemPodMetrics_*.json`에서 시스템 파드 restartCount 확인됨. |

→ **첫 CAT SLI 세트에 포함 후보.** SLO 값만 정당화하면 즉시 평가 가능.

### 4.2 Prometheus 측정 필요 (파이프라인 준비 완료 시 수집 가능)

Precheck(§11)에서 API server scrape 포트 6443 적용 시 타깃 up·기본 쿼리 성공이 확인되었으므로, **full load 실험 완료 시** 아래 SLI는 산출될 수 있다.

| SLI | 비고 |
|-----|------|
| **API server request latency** | APIResponsivenessPrometheus* 산출. |
| **In-cluster network latency** | InClusterNetworkLatency. |
| **Network programming latency** | NetworkProgrammingLatency. |
| **Kube-proxy iptables restore failures** | GenericPrometheusQuery (KubeProxy). |

→ **Prometheus 활성화 full run 완료 후** 산출 파일 유무로 “측정 가능” 여부 확정. Review 단계에서 “첫 CAT SLI 세트” 또는 “실험용”으로 구분 가능.

### 4.3 옵션·불확실

| SLI | 비고 |
|-----|------|
| **API availability** | APIAvailability 옵션 활성화·산출 경로 확인 필요. **불확실.** |
| **Phase duration** | TestMetrics/SchedulingMetrics가 kind에서 실패한 이력 있음. **불확실.** 필요 시 ResourceUsageSummary 등으로 부분 활용. |

→ **선택 SLI** 또는 “측정 방법 확립 후 추가” 후보.

---

## 5. Review 단계 결정을 위한 요약

- **첫 CAT SLI 세트 후보 (즉시 평가 가능):**  
  Pod startup latency, Container restarts/OOM, System pod health.  
  → 이미 실험에서 파일·키 경로·예시 값 확인됨.

- **Prometheus 의존 SLI (full run 완료 후 포함 여부 결정):**  
  API server request latency, In-cluster network latency, Network programming latency, Kube-proxy iptables restore failures.  
  → Precheck 통과로 수집 가능성 높음; 실제 full run 산출물로 확인 후 “첫 세트” 또는 “실험용”으로 표시.

- **선택/실험용 또는 후순위:**  
  API availability (옵션·경로 확인 필요), Phase duration (kind 환경에서 불안정).  
  → “첫 CAT SLI 세트”가 아닌 선택 SLI 또는 후속 실험으로 측정 방법 확립 후 추가하는 것을 권장.

---

*Experimenter 역할 산출물. Review 단계에서 첫 CAT SLI 세트·선택 SLI·실험용 구분에 활용.*
