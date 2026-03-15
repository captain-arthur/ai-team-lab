# Architecture: SLI/SLO 모델 및 devcat 실험 구조

**Project:** sample-cl2-sli-slo-analysis  
**Role:** Architect  
**Input:** 01-manager/project-brief.md, 02-research/research-notes.md, knowledge/principles (cat-vision, cat-design-principles, sli-slo-philosophy, devcat-program-brief)

---

## 1. SLI 모델 (Service Level Indicator)

Research에서 도출한 ClusterLoader2 메트릭을 바탕으로 **수용 테스트용 SLI 후보**를 정의한다. 각 SLI는 “무엇을 측정하는지”, “사용자·워크로드 관점에서 무엇을 의미하는지”를 명시한다.

### 1.1 SLI 후보 정의

| SLI (후보) | ClusterLoader2 메트릭 소스 | 측정 대상 | 단위/형태 | 사용자·워크로드 관점에서의 의미 |
|------------|----------------------------|-----------|-----------|----------------------------------|
| **Pod startup latency** | CreatePhasePodStartupLatency, PodStartupLatency 모듈 | 파드 생성부터 Ready까지 소요 시간 | P50/P90/P99 (초) | 배포·스케일링 시 “파드가 언제 쓸 수 있게 되는가”. 이 값이 크면 롤아웃·스케일이 느리게 체감됨. |
| **API call latency** | APIResponsivenessPrometheus, APIResponsivenessPrometheusSimple | API 호출별 레이턴시 | verb/resource별 P99, 슬로우 콜 수 | kubectl·컨트롤러·오퍼레이터가 API를 쓸 때 응답 지연. 슬로우 콜이 많으면 작업이 멈춘 것처럼 보일 수 있음. |
| **In-cluster network latency** | InClusterNetworkLatency | 파드 간 네트워크 RTT | P99 (ms 또는 초) | 마이크로서비스·파드 간 통신 지연. 수용 시 “클러스터 내 네트워크가 정상 범위인가” 판단에 사용. |
| **Network programming latency** | NetworkProgrammingLatency | 서비스/엔드포인트 변경이 데이터플레인에 반영되는 시간 | P99 (초) | 트래픽이 새 파드로 가기까지 걸리는 시간. “서비스가 금방 반영되는가”와 연결. |
| **API availability** | APIAvailability | API 서버 응답 성공 비율 | % (측정 구간) | 클러스터가 “사용 가능”인지의 기본 지표. 일정 비율 미만이면 수용 불가. |
| **Container restarts** | ContainerRestarts, TestMetrics (OOM 등) | 테스트 구간 중 컨테이너 재시작·OOM | 횟수, OOM 발생 여부 | 워크로드 안정성. 비정상 재시작·OOM이 있으면 “이 클러스터에서 안정적으로 돌지 않는다”로 해석. |
| **Kube-proxy iptables restore failures** | KubeProxyIptablesRestoreFailures | partial restore 실패 횟수 | 횟수 | 0이어야 정상. 실패 시 트래픽 라우팅 오류 가능. 결함 지표. |
| **Phase duration (선택)** | TestMetrics (단계 타이머) | create/scale/delete 단계 소요 시간 | 초 | “N 파드 생성이 T 초 이내에 끝나는가” 같은 시나리오 완료 시간. 실용적 수용 의미 부여 가능. |

- 소규모 클러스터에서는 공식 load config에서 **pod-startup-latency**, **scheduler-throughput** 모듈이 비활성화되므로, 실제로 수집되는 SLI는 사용하는 시나리오·오버라이드에 따라 제한될 수 있다. devcat 실험 시 “어떤 SLI를 켜고 해석할지”를 시나리오별로 명시하는 것이 좋다.

### 1.2 SLI 정의 시 원칙 (Kubernetes scalability 철학)

- **Precise and well-defined:** SLI마다 “무엇을, 어떤 시나리오·파라미터(노드 수, 파드 수, QPS 등), 어떤 단위로 측정하는지”를 문서에 고정한다.
- **Consistent:** 동일 config·동일 클러스터에서 반복 실행 시 같은 방식으로 측정된다. devcat에서는 config.yaml·ol-test.yaml 등 오버라이드를 버전 관리해 재현 가능하게 한다.
- **User-oriented:** 위 표의 “사용자·워크로드 관점”을 각 SLI에 연결해, “이 지표가 나쁘면 사용자가 무엇을 체감하는가”를 설명할 수 있게 한다.
- **Testable:** ClusterLoader2(및 perfdash)가 실제로 산출하는 포맷·경로로 측정 가능해야 한다. 산출이 없거나 불안정한 SLI는 후보에서 제외하거나, 측정 방법을 먼저 확립한 뒤 SLO로 승격한다.

---

## 2. SLO 정당화 모델 (소규모 클러스터용)

### 2.1 Kubernetes scalability 철학 및 “you promise, we promise”

- **Precise and well-defined:** SLO 값(예: P99 ≤ 5s)과 “어떤 조건에서, 어떤 단위로” 측정하는지를 명시. 모호한 “빠르다/느리다”가 아니라 숫자·조건으로 고정.
- **Consistent:** 동일 환경·시나리오에서 반복 측정 가능. 측정 방법·도구 버전(ClusterLoader2, report-dir 규약)을 문서화.
- **User-oriented:** “이 SLO를 만족하면 사용자·워크로드가 기대하는 동작을 한다”는 설명을 둠. 예: “P99 pod startup ≤ 5s → 일반적인 배포/스케일링이 체감상 부담 없이 완료된다.”
- **Testable:** 수용 테스트에서 “SLO 충족 여부”를 자동 또는 수동으로 판단할 수 있어야 함. ClusterLoader2 결과에서 해당 메트릭을 추출·비교할 수 있는 경로가 있어야 한다.
- **“You promise, we promise”:**  
  - **You promise:** 팀/플랫폼이 “이 노드 수·이 워크로드 패턴·이 시나리오에서, 이 수준을 보장한다”고 **명시적으로 약속**하는 값이 SLO다.  
  - **We promise:** 수용 테스트(ClusterLoader2 실행 + 결과 해석)가 그 약속이 지켜졌는지 **검증**한다.  
  SLO 값은 “도구 기본값”이 아니라 **정당한 근거**로 채워져야 한다.

### 2.2 SLO 값 정당화의 네 가지 기준 (sli-slo-philosophy)

| 기준 | 내용 |
|------|------|
| **User expectations** | 사용자(워크로드 소유자, 플랫폼 이용자)가 기대하는 응답 시간·가용성·처리량. “5초 안에 파드가 준비되면 만족” 등. |
| **Environment assumptions** | 노드 수, 리소스, 네트워크, 시나리오(파드 수, QPS, config 이름 등). “N 노드, M 파드 create 시나리오에서”를 명시. |
| **Repeatable measurement** | 동일 조건에서 반복 측정 가능, 측정 방법·도구 버전이 문서화됨. ClusterLoader2 report-dir, results/ 규약 등. |
| **Practical acceptance meaning** | “이 SLO를 만족하면 이 클러스터를 수용한다”가 팀·운영에서 말이 되도록 정의. 수용/거부를 구분할 수 있는 수준이어야 함. |

- SLO를 정할 때 “왜 이 숫자인가?”에 대한 답을 위 네 가지로 줄 수 있어야 한다. Engineer 단계에서 **SLO 정당화 템플릿**으로 이 네 항목을 채우는 형태를 둘 수 있다.

### 2.3 벤치마크 스타일 임계값이 수용 테스트에 부적합한 이유

- **예: 1시간(1h) PodStartupLatency.**  
  ClusterLoader2 공식 measurements.yaml에는 CreatePhasePodStartupLatency에 `threshold: 1h`가 있다. 대규모 saturation(수천 파드)에서 “전부 기동될 때까지” 기다리는 용도로 쓰이지만, **수용 테스트의 “이 클러스터를 수용한다”는 판단 기준**으로는 의미가 없다. 거의 모든 클러스터가 1h 안에 통과해 수용/거부를 구분하지 못한다. 즉 **practical acceptance meaning**이 없음.

- **일반화:**  
  - 대규모·고부하 벤치마크(예: 5000노드, 수만 파드)에서 쓰는 latency/throughput 임계값은 **소규모 클러스터의 “일상 부하에서 기대하는 수준”**과 다르다. 그대로 재사용하면 과도하게 느슨(항상 통과)하거나, 반대로 환경에 맞지 않아 과도하게 엄격할 수 있다.  
  - 도구 기본값이 “어떤 시나리오·어떤 환경”에서 정의되었는지 문서화되어 있지 않으면 **consistent·testable** 원칙에도 어긋난다.  
  - 따라서 **수용 테스트용 SLO는 “우리 환경·시나리오”를 전제로**, user expectations, environment assumptions, repeatable measurement, practical acceptance meaning으로 **별도 정당화**한다. 벤치마크 스타일의 매우 큰 임계값(1h 등)은 수용 테스트에 쓰지 않는다.

---

## 3. devcat 실험 아키텍처

devcat-program-brief의 현재 현실(ClusterLoader2, perfdash, config.yaml, ol-test.yaml, results/)을 전제로, 실험이 “어디서 실행되고, 결과가 어디에 쌓이며, 메트릭을 어떻게 추출하고, SLO 검사를 어떻게 할 수 있는지”를 개념 수준에서 정의한다.

### 3.1 ClusterLoader2 실행

- **실행 주체:** 플랫폼/SRE 또는 CI. devcat 저장소에서 ClusterLoader2 바이너리(또는 래퍼 스크립트)를 사용.
- **설정:** 부하 시나리오는 **config.yaml**을 기준으로 하고, 클러스터별 차이는 **오버라이드 파일**(예: **ol-test.yaml**)로 반영. (devcat-program-brief: “copy load scenarios and execute config.yaml with cluster-specific override files”.)
- **실행 방식:** 로컬 또는 CI에서 `clusterloader2 --testconfig=... --report-dir=...` 형태로 실행. report-dir은 **한 run당 하나의 디렉터리**로 두어, 결과가 run별로 구분되게 한다.
- **소규모 클러스터:** 공식 load config는 100+ 노드를 가정하므로, 소규모용으로는 노드 수·파드 수·스킵되는 모듈을 오버라이드로 조정한다. “소규모용 최소 시나리오”는 Engineer 단계에서 config 조합·오버라이드 예시로 구체화할 수 있다.

### 3.2 결과 저장 위치

- **저장 위치:** devcat의 **results/** 디렉터리. (devcat-program-brief: “Results are stored under results/”.)
- **규약:** 한 번의 ClusterLoader2 run마다 **results/** 아래에 run을 구분할 수 있는 서브디렉터리(예: run-id, 타임스탬프, 브랜치명 등)를 두고, 그 안에 ClusterLoader2가 산출한 요약 파일·로그·(선택) perfdash 입력 파일을 둔다. 경로 규약은 devcat 현재 구조를 따르며, 필요 시 “results/<run-id>/” 형태로 통일하는 정도만 제안한다.

### 3.3 메트릭 추출

- **입력:** 위 results/ 내의 ClusterLoader2 산출물. 구체 포맷은 ClusterLoader2 버전·report-dir 내용에 따름(요약 JSON, JUnit, 로그 등). perfdash가 읽는 형식이 있다면 동일한 파일을 “메트릭 추출”의 입력으로 사용할 수 있다.
- **추출 내용:** SLI 후보에 해당하는 메트릭 값(P50/P90/P99, 슬로우 콜 수, 재시작 횟수, 단계 소요 시간 등). 추출은 (1) 수동으로 요약 파일을 열어 확인하거나, (2) 작은 스크립트로 결과 디렉터리를 파싱해 구조화된 값(JSON/표)으로 만드는 방식이 있다. Engineer 단계에서 “어떤 파일에서 어떤 키를 읽을지”를 runbook 또는 스크립트 초안으로 정리할 수 있다.
- **출력:** **SLI measurements** — SLI별로 한 run에 대응하는 측정값 목록. 예: `{ "pod_startup_latency_p99_sec": 4.2, "api_slow_calls_count": 0, ... }`. 이 출력은 SLO 평가와 해석 노트의 입력이 된다.

### 3.4 SLO 검사 (평가)

- **입력:** (1) **SLO 정의** — SLI별 목표값·조건(예: pod_startup_latency_p99 ≤ 5s, api_slow_calls = 0). (2) **SLI measurements** — 위에서 추출한 측정값.
- **로직:** 각 SLO에 대해 “측정값이 목표를 만족하는지” yes/no 판단. 전 SLO가 만족되면 “SLO evaluation: pass”, 하나라도 불만족이면 “fail”로 기록.
- **구현 수준:** 수동(스프레드시트·문서)으로 할 수도 있고, 스크립트로 measurements와 SLO 정의를 읽어 비교한 뒤 **SLO evaluation** 결과( pass/fail, 불만족 SLO 목록)를 출력할 수도 있다. Architect 단계에서는 “입력·출력·판단 기준”만 정의하고, 구현은 Engineer에 맡긴다.
- **출력:** **SLO evaluation** — run별로 “pass/fail”, (선택) 불만족 SLO 목록·초과분. 이 결과는 interpretation 단계에서 “왜 fail이었는지”, “SLO 값을 조정할지, 환경을 개선할지” 논의하는 데 쓰인다.

---

## 4. 실험 워크플로 (Experiment workflow)

다음 흐름을 한 사이클로 둔다. research → devcat experiment → interpretation → SLO refinement → devcat improvement.

```
research
    → devcat experiment (ClusterLoader2 실행, results/ 저장)
    → interpretation (메트릭 추출, SLO 평가, “왜 이 수준인가” 해석)
    → SLO refinement (필요 시 SLO 값·정당화 문서 조정)
    → devcat improvement (시나리오·오버라이드·추출/평가 스크립트·perfdash 외 시각화 등 점진적 개선)
```

- **Research:** 메트릭·SLI 후보·정당화 기준을 문서로 정리(본 프로젝트의 02-research, 03-architecture가 해당). “어떤 SLI를 볼 것인가”, “SLO 후보 값과 정당화 초안”이 나온다.
- **devcat experiment:** devcat에서 ClusterLoader2를 실행(config + 오버라이드), results/에 결과 저장. 동일 config 반복 실행으로 일관성(consistent)을 확인할 수 있다.
- **Interpretation:** results/에서 메트릭을 추출해 SLI measurements를 만들고, SLO 정의와 비교해 SLO evaluation을 수행. “실제 값이 SLO를 만족하는지”, “만족하지 않으면 어떤 지표가 나쁜지”를 기록하고, **interpretation notes**(해석 노트)에 “왜 이 수준이 나왔는지”, “환경·시나리오 가정이 맞는지”를 짧게 적는다.
- **SLO refinement:** 해석 결과에 따라 SLO 후보 값을 유지·상향·하향하거나, 정당화 문서(user expectations, environment assumptions 등)를 수정. “실제로 10s가 나오는데 5s로 두었으면” 환경 개선 또는 SLO 완화(및 정당화 업데이트) 중 하나를 선택.
- **devcat improvement:** 시나리오 추가·오버라이드 정리, 메트릭 추출·SLO 평가를 반복 가능하게 하는 runbook/스크립트, perfdash 외 시각화(예: assertion/PASS·FAIL 표현) 등 devcat을 점진적으로 개선. 한 번에 완전히 바꾸지 않고, “다음에 할 수 있는 한 단계”를 정해 진행한다.

---

## 5. 입력·출력 정의

### 5.1 입력 (Inputs)

| 입력 | 설명 | 소스 |
|------|------|------|
| **ClusterLoader2 run results** | 한 번의 ClusterLoader2 실행이 산출한 파일 전체. | devcat **results/** 아래 run별 디렉터리. report-dir에 해당하는 경로. |
| **results/ 내 메트릭** | 요약 JSON, JUnit, 로그, perfdash용 데이터 등 ClusterLoader2·perfdash가 생성한 메트릭 파일. | 동일 run 디렉터리 내. |
| **SLO 정의** | SLI별 목표값·조건(예: P99 ≤ 5s, slow_calls = 0). 정당화 문서(네 가지 기준)와 함께 버전 관리. | 팀 문서 또는 devcat 내 config/문서. |
| **시나리오·환경 정보** | 사용한 config.yaml, 오버라이드 파일 이름, 노드 수, 파드 수 등. | 실험 run 메타데이터로 기록. |

### 5.2 출력 (Outputs)

| 출력 | 설명 | 사용처 |
|------|------|--------|
| **SLI measurements** | 한 run에 대한 SLI별 측정값(숫자·단위). | SLO 평가의 입력, 해석·트렌드 분석. |
| **SLO evaluation** | 해당 run에 대한 SLO 충족 여부(pass/fail), (선택) 불만족 SLO 목록·초과분. | 수용/거부 판단, interpretation·refinement 입력. |
| **Interpretation notes** | “왜 이 수준이 나왔는지”, “환경 가정이 맞는지”, “SLO를 바꿀지 환경을 바꿀지”에 대한 짧은 해석. | SLO refinement, devcat improvement 방향 결정. |

- SLI measurements와 SLO evaluation은 **일관된 형식**(예: JSON, Markdown 표)으로 남기면, 이후 자동 시각화(Evidence 등)나 대시보드에서 “숫자 + PASS/FAIL” 표현을 붙이는 데 재사용할 수 있다.

---

## 6. 설계 결정 요약

| ID | 결정 | 선택 | 이유 |
|----|------|------|------|
| 1 | SLI 후보 범위 | Research에서 도출한 7개 + Phase duration(선택) | 수용 테스트에 의미 있고 사용자·워크로드 관점과 연결 가능한 것만 포함. |
| 2 | SLO 값 출처 | 도구 기본값이 아닌 네 가지 기준으로 정당화 | sli-slo-philosophy·벤치마크 임계값 한계 반영. |
| 3 | 실험 결과 위치 | devcat results/ | devcat-program-brief 현재 현실 준수. |
| 4 | 메트릭 추출·SLO 평가 | 수동 + (선택) 스크립트 | 점진적 개선; 먼저 수동으로 경로 확립 후 자동화. |
| 5 | 워크플로 | research → experiment → interpretation → SLO refinement → devcat improvement | devcat-program-brief의 작업 모델과 일치. |

---

## 7. 트레이드오프·리스크

- **트레이드오프:** SLI 개수를 적게 유지하면 해석이 단순하지만, 놓치는 지표가 있을 수 있음. 반대로 많으면 해석 부담과 측정 불안정 가능성이 커짐. 소규모 클러스터·초기에는 핵심 SLI(예: pod startup, API responsiveness, restarts/OOM)만 먼저 SLO로 두고 확장하는 것을 권장.
- **리스크:** (1) devcat의 results/ 구조·ClusterLoader2 산출 포맷이 버전에 따라 바뀔 수 있음. 추출 로직·runbook은 “어떤 파일·어떤 키”를 참조하는지 문서화해 두어 변경 시 수정하기 쉽게 한다. (2) 소규모 클러스터에서 공식 load config의 일부 모듈이 스킵되면, 해당 SLI는 “수집 불가”가 될 수 있음. “이 시나리오에서는 이 SLI만 평가한다”처럼 시나리오별로 평가 대상 SLI를 명시하는 것이 좋다.

---

*Architect 역할 산출물. 구현 상세(정당화 템플릿, 추출 스크립트, runbook)는 Engineer 단계에서 진행.*
