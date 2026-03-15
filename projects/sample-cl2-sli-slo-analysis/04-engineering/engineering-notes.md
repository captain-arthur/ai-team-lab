# Engineering: devcat 수용 테스트 실험 워크플로

**Project:** sample-cl2-sli-slo-analysis  
**Role:** Engineer  
**Input:** 01-manager, 02-research, 03-architecture, knowledge/principles

이 문서는 Architecture의 SLI/SLO 모델과 devcat 실험 구조를 **실행 가능한 runbook·추출·평가 모델**로 옮긴 것이다. ClusterLoader2 소스 수정 없이, 기존 산출물을 사용하고, 단순·도구 중립적인 방식을 따른다.

---

## 1. devcat 실험 Runbook

devcat 저장소에서 한 번의 수용 테스트(ClusterLoader2 실행 → 결과 수집 → 메트릭 추출 → SLO 평가 → PASS/FAIL 판정)를 수행하는 단계별 절차다.

### 사전 요건

- devcat 저장소 클론·체크아웃 완료.
- ClusterLoader2 바이너리 또는 실행 방법 준비(devcat에서 사용하는 방식 준수).
- 대상 Kubernetes 클러스터에 대한 kubeconfig 설정.
- results/ 디렉터리 생성·쓰기 권한.

### Step 1: Run ID 및 결과 디렉터리 준비

- 이번 run을 구분할 **run-id** 생성. 예: `$(date +%Y%m%d-%H%M%S)` 또는 CI job id.
- devcat **results/** 아래에 run 전용 디렉터리 생성.

```bash
export RUN_ID=$(date +%Y%m%d-%H%M%S)
export REPORT_DIR="results/${RUN_ID}"
mkdir -p "${REPORT_DIR}"
```

- 이후 ClusterLoader2의 **report-dir**을 이 경로로 지정하면, 한 run의 결과가 모두 `results/<run-id>/` 아래에 모인다.

### Step 2: 오버라이드 파일 준비

- 사용할 **config.yaml**(또는 devcat이 제공하는 부하 시나리오 config)을 결정.
- 클러스터·환경에 맞는 **오버라이드 파일**(예: **ol-test.yaml**)을 준비. 노드 수, 파드 수, 시나리오 파라미터 등이 소규모 클러스터에 맞게 조정되어 있어야 함. (devcat-program-brief: “execute config.yaml with cluster-specific override files such as ol-test.yaml”.)
- 오버라이드 파일 경로를 기록해 두고, 실험 메타데이터(run-id, config 이름, 오버라이드 이름, 노드 수 등)를 나중에 interpretation notes에 남길 수 있게 한다.

### Step 3: ClusterLoader2 실행 (devcat 방식)

- devcat에서 사용하는 실행 명령 또는 스크립트로 ClusterLoader2를 실행한다. **testconfig**에는 config.yaml(및 오버라이드 적용 결과), **report-dir**에는 위에서 만든 `results/<run-id>/`를 지정.

예시(구문만, devcat 실제 명령에 맞게 수정):

```bash
# devcat에서 사용하는 실행 방식에 맞게 조정할 것
clusterloader2 --kubeconfig="${KUBECONFIG}" \
  --testconfig=path/to/config.yaml \
  --report-dir="${REPORT_DIR}"
```

- 실행이 정상 종료되면 결과는 이미 `results/<run-id>/` 아래에 저장된 상태다.

### Step 4: 결과 수집 확인

- `results/<run-id>/` 내부에 ClusterLoader2가 생성한 파일(요약 JSON, JUnit, 로그, perfdash용 데이터 등)이 있는지 확인.
- 실제 파일 이름·구조는 ClusterLoader2 버전·옵션에 따라 다를 수 있으므로, devcat에서 한 번 실행한 뒤 디렉터리 목록을 기록해 두면 다음 run에서 “어디서 무엇을 읽을지”를 runbook에 반영할 수 있다.

### Step 5: 메트릭 추출 (SLI measurements)

- §2의 “SLI 메트릭 추출”에 따라, results/<run-id>/ 내 파일에서 SLI 후보에 해당하는 값을 읽어 **SLI measurements**를 만든다.
- 수동: 요약 파일을 열어 해당 필드를 표나 메모에 옮긴다.
- 반자동: §2에 적힌 “어떤 파일·어떤 키/필드”를 참조하는 작은 스크립트로 JSON·CSV·Markdown 표 형태로 출력한다.
- 출력은 **SLI measurements** 문서(또는 `results/<run-id>/sli-measurements.json` 등)로 남겨, 다음 단계와 interpretation notes에서 참조한다.

### Step 6: SLO 평가 (SLO evaluation)

- §3의 “SLO 평가 모델”에 따라, **SLO 정의**(SLI별 목표값·조건)와 **SLI measurements**를 비교해 각 SLO 충족 여부를 판단한다.
- 모든 평가 대상 SLO가 만족되면 **PASS**, 하나라도 불만족이면 **FAIL**로 기록.
- 불만족 SLO 목록·초과분(측정값 vs 목표값)을 함께 적어 두면 interpretation 시 유리하다.
- 출력은 **SLO evaluation** 문서(또는 `results/<run-id>/slo-evaluation.md` 등)로 남긴다. 형식 예: `Overall: PASS | FAIL`, `SLO_1: pass`, `SLO_2: fail (measured X, threshold Y)`.

### Step 7: PASS/FAIL 판정 및 해석 노트

- SLO evaluation의 **Overall** 결과가 해당 run에 대한 **수용 테스트 판정**(PASS = 수용, FAIL = 거부 또는 재검증 필요)이다.
- **Interpretation notes**를 작성한다. “왜 이 수준이 나왔는지”, “환경·시나리오 가정이 맞는지”, “SLO를 조정할지 환경을 개선할지”에 대한 짧은 메모. `results/<run-id>/interpretation-notes.md` 등으로 저장하면 된다.

---

## 2. SLI 메트릭 추출 (ClusterLoader2 결과 → SLI 후보)

### 2.1 결과 디렉터리에서 메트릭이 나오는 위치

ClusterLoader2는 **--report-dir**에 지정한 디렉터리에 실행 결과를 쓴다. 실제 파일 이름·구조는 버전과 빌드 옵션에 따라 다를 수 있으며, devcat이 사용하는 ClusterLoader2 버전과 동일한 구조를 기준으로 해야 한다. 아래는 **일반적으로 기대할 수 있는 산출물**과, Architecture에서 정의한 SLI 후보와의 **매핑**이다.

| ClusterLoader2 측정/산출 | 결과 내 위치(예상) | SLI 후보와의 매핑 |
|--------------------------|--------------------|--------------------|
| APIResponsivenessPrometheus / Simple | 요약 파일 내 API 레이턴시·슬로우 콜 | **API call latency** (P99, 슬로우 콜 수) |
| CreatePhasePodStartupLatency | PodStartupLatency 관련 요약·로그 | **Pod startup latency** (P50/P90/P99, 초) |
| InClusterNetworkLatency | 네트워크 레이턴시 요약 | **In-cluster network latency** (P99) |
| NetworkProgrammingLatency | 네트워크 프로그래밍 레이턴시 요약 | **Network programming latency** (P99) |
| APIAvailability | 가용성 비율 (옵션 수집 시) | **API availability** (%) |
| ContainerRestarts, TestMetrics (OOM) | 재시작 횟수·OOM 플래그 | **Container restarts** (횟수, OOM 여부) |
| KubeProxyIptablesRestoreFailures | kube-proxy 메트릭 요약 | **Kube-proxy iptables restore failures** (횟수) |
| TestMetrics (단계 타이머) | 단계별 소요 시간 | **Phase duration** (create/scale/delete, 초) |

- **실제 경로:** ClusterLoader2는 보통 `report-dir` 직하에 요약 JSON, JUnit XML, 또는 디렉터리(예: `measurement-*`)를 둔다. devcat에서 한 번 실행한 뒤 `results/<run-id>/` 아래에 어떤 파일이 생성되는지 확인하고, 위 매핑을 “파일 경로 + 키/필드 이름”으로 구체화하는 것을 권장한다. 예: `summary.json` → `PodStartupLatency.P99`, `APIResponsivenessPrometheus.SlowCalls` 등.
- **수집 불가 SLI:** 소규모 클러스터에서 pod-startup-latency·scheduler-throughput 모듈이 스킵되면, Pod startup latency·Phase duration의 일부는 수집되지 않을 수 있다. “이 run에서 평가하는 SLI 목록”을 run별로 명시해 두면, 수집 불가 SLI는 SLO 평가에서 제외하거나 N/A로 표시할 수 있다.

### 2.2 추출 결과 형태 (SLI measurements)

- **목적:** 한 run에 대해 “SLI 이름 → 측정값(숫자·단위)”를 일관된 형태로 남겨, SLO 평가와 해석에서 재사용한다.
- **권장 형식 예 (JSON):**

```json
{
  "run_id": "20250315-120000",
  "scenario": "load",
  "sli": {
    "pod_startup_latency_p99_sec": 4.2,
    "api_slow_calls_count": 0,
    "in_cluster_network_latency_p99_ms": 12,
    "network_programming_latency_p99_sec": 2.1,
    "api_availability_pct": 99.9,
    "container_restarts_count": 0,
    "kubeproxy_iptables_restore_failures": 0,
    "phase_create_duration_sec": 120
  }
}
```

- **수동 추출:** 위 표를 참고해 요약 파일을 열고, 해당하는 숫자를 `sli-measurements.json` 또는 Markdown 표로 옮긴다.
- **경량 스크립트:** “results/<run-id>/ 내 특정 파일을 읽어 위와 같은 키로 값을 채운 JSON을 출력”하는 스크립트를 둘 수 있다. ClusterLoader2 소스는 수정하지 않고, **기존 산출 파일만** 읽는다. 파일 경로·키 이름은 devcat 실제 구조에 맞춰 조정한다.

---

## 3. SLO 평가 모델 (단순 비교 → pass/fail)

### 3.1 모델

- **입력:** (1) **SLO 정의** — SLI별로 “목표 조건”(예: ≤ 5, = 0, ≥ 99.5). (2) **SLI measurements** — §2에서 추출한 측정값.
- **로직:**  
  - 각 SLO에 대해: **측정값**을 **SLO 임계값(또는 조건)**과 비교.  
  - “낮을수록 좋은” 지표(레이턴시, 슬로우 콜 수, 재시작 수 등): `measured ≤ threshold` 이면 pass.  
  - “높을수록 좋은” 지표(가용성 %): `measured ≥ threshold` 이면 pass.  
  - “0이어야 함” 지표(실패 횟수): `measured == 0` 이면 pass.  
- **출력:**  
  - SLO별 **pass / fail**.  
  - **Overall:** 평가 대상 SLO가 **모두 pass**이면 **PASS**, 하나라도 fail이면 **FAIL**.  
  - (선택) fail인 SLO 목록·측정값 vs 임계값을 함께 기록.

### 3.2 SLO 정의 예시 (소규모 클러스터용 초안)

- 값은 **정당화 문서**(user expectations, environment assumptions, repeatable measurement, practical acceptance meaning)와 함께 유지·갱신한다. 아래는 “예시”이며, 실제 숫자는 팀에서 정당화한 뒤 채운다.

| SLI | SLO 조건 예시 | 비고 |
|-----|----------------|------|
| Pod startup latency P99 | ≤ 5 (초) | 소규모·일반 부하 가정 |
| API slow calls | = 0 (또는 ≤ N) | N은 정당화로 결정 |
| In-cluster network latency P99 | ≤ 50 (ms) | 환경 가정에 따라 조정 |
| Network programming latency P99 | ≤ 30 (초) | 대규모용 30s와 구분해 소규모용으로 검토 |
| API availability | ≥ 99.5 (%) | |
| Container restarts | ≤ 0 (또는 팀 정책) | |
| Kube-proxy iptables restore failures | = 0 | |
| Phase create duration | ≤ T (초) | 시나리오별 T 정의 |

### 3.3 평가 실행

- **수동:** SLI measurements 표와 SLO 정의 표를 나란히 두고, 각 SLI에 대해 위 비교 규칙으로 pass/fail을 적고, 전부 pass이면 Overall PASS, 아니면 FAIL로 기록.
- **경량 스크립트:** SLO 정의를 JSON/YAML로 두고, SLI measurements와 비교한 뒤 pass/fail·Overall을 출력하는 스크립트를 둘 수 있다. ClusterLoader2나 devcat 소스를 건드리지 않고, **우리가 만든 measurements·정의 파일만** 읽어서 평가한다.

---

## 4. 최소 도구 접근 (Minimal tooling)

- **ClusterLoader2 소스 수정 금지.** as-is 사용. 설정·오버라이드로만 조정.
- **기존 산출물만 사용.** ClusterLoader2가 report-dir에 쓰는 파일만 읽어서 메트릭 추출·SLO 평가에 쓴다. 새 프로토콜이나 비공식 출력을 가정하지 않는다.
- **단순·도구 중립:**  
  - **수동:** 요약 파일을 열어 SLI 값을 옮기고, SLO 정의 표와 비교해 pass/fail·Overall·interpretation notes를 문서로 작성.  
  - **선택적 경량 스크립트:** (1) results/<run-id>/ 내 특정 파일 파싱 → SLI measurements JSON/표 출력. (2) SLO 정의 + SLI measurements 입력 → SLO evaluation( pass/fail, 불만족 목록) 출력.  
  - 스크립트는 언어·환경에 종속되지 않게, “어떤 파일·어떤 키를 읽고, 어떤 규칙으로 비교하는지”가 문서화되어 있으면 나중에 다른 도구로 대체 가능하다.

---

## 5. 실험 산출물 (Outputs)

한 번의 devcat 수용 테스트 run이 만들어 내는 산출물을 정리한다.

| 산출물 | 설명 | 권장 위치/형식 |
|--------|------|----------------|
| **ClusterLoader2 원시 결과** | report-dir에 쌓인 요약·로그·JUnit 등 | `results/<run-id>/` (ClusterLoader2가 생성) |
| **SLI measurements** | SLI별 측정값(숫자·단위). §2.2 형태. | `results/<run-id>/sli-measurements.json` 또는 동일 run 내 문서 |
| **SLO evaluation** | SLO별 pass/fail, **Overall: PASS / FAIL**. (선택) 불만족 SLO·측정값 vs 임계값 | `results/<run-id>/slo-evaluation.md` (또는 .json) |
| **Interpretation notes** | “왜 이 수준인가”, “환경 가정 적합성”, “SLO 조정 vs 환경 개선”에 대한 짧은 해석 | `results/<run-id>/interpretation-notes.md` |
| **Run 메타데이터** | run-id, config 이름, 오버라이드 이름, 노드 수, 실행 시각 등 | runbook Step 1·2에서 기록; interpretation notes 또는 별도 metadata.json에 포함 가능 |

- SLI measurements와 SLO evaluation을 **일관된 형식**(JSON, Markdown 표)으로 두면, 이후 자동 시각화(Evidence 등)나 “숫자 + PASS/FAIL” 대시보드 연동에 재사용할 수 있다.

---

## 6. 구현 노트·가정·한계

- **가정:** devcat은 Kubernetes 공식 perf-tests의 ClusterLoader2와 유사한 방식으로 config·오버라이드·report-dir을 사용한다. results/ 구조가 perf-tests의 report-dir과 동일하지 않을 수 있으므로, **실제 devcat에서 한 번 실행한 뒤 results/<run-id>/ 내용을 확인해** §2.1의 “결과 내 위치”와 §2.2의 추출 방식을 구체화하는 것이 필요하다.
- **ClusterLoader2 미수정:** 모든 단계는 ClusterLoader2 소스 변경 없이, 기존 바이너리·기존 출력만 전제로 한다.
- **TODO/후속:** (1) devcat results/ 실제 디렉터리 목록·파일 포맷을 기록한 “결과 레이아웃” 문서. (2) SLI 추출·SLO 평가를 수행하는 최소 스크립트 예시(언어 무관, “어떤 파일·키·비교 규칙”만 명시). (3) SLO 정당화 템플릿(네 가지 기준 항목을 채우는 양식) — 03-architecture §2.2에 맞춰 Engineer가 별도 파일로 둘 수 있다.
- **이탈:** Architecture와 동일한 입력·출력·워크플로를 따른다. 구체적인 파일 경로·키 이름은 devcat·ClusterLoader2 버전에 따라 달라질 수 있어, 본 문서는 “위치·매핑·형식”을 일반적으로 기술하고, devcat 측에서 한 번 검증 후 보완하는 흐름을 권장한다.
