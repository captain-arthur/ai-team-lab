# Experiment: devcat ClusterLoader2 실제 실행 및 결과 분석

**Project:** sample-cl2-sli-slo-analysis  
**Role:** Experimenter  
**Input:** 01-manager, 02-research, 03-architecture, 04-engineering

---

## 1. 실행 요약

### 1.1 실행 경로 및 명령

- **실행 위치:** devcat 저장소 루트 (`devcat/`).
- **실제 사용 명령:**

```bash
cd devcat
PROVIDER=local KUBECONFIG="${KUBECONFIG:-$HOME/.kube/config}" \
CL2_BINARY="${PWD}/bin/cl2" \
./scripts/run-devcat.sh
```

- **바이너리:** ClusterLoader2는 `devcat/bin/cl2`를 사용. (run-devcat.sh 기본값은 `bin/clusterloader2`이므로, 이번 run에서는 환경변수 `CL2_BINARY`로 `bin/cl2` 지정.)
- **Provider:** ClusterLoader2는 `--provider`를 CLI/환경변수로 받음. 오버라이드 파일의 `provider` 필드는 시나리오 템플릿용이며, 실행 시 `PROVIDER=local`을 설정해야 “unsupported provider name” 오류를 피할 수 있음.
- **Kubeconfig:** 기본 kubeconfig(`$HOME/.kube/config`) 사용. kind 클러스터(3노드)에 연결됨.

### 1.2 Run 결과

- **Run ID:** `20260315-181750`
- **시나리오:** `scenarios/load/config.yaml`
- **오버라이드:** `overrides/ol-test.yaml`
- **테스트 결과:** ClusterLoader2가 **실제로 부하 테스트를 수행**했으나, 일부 measurement 실패로 **전체 상태는 Fail**. (TestMetrics/SchedulingMetrics: kind 환경에서 scheduler 엔드포인트 미접근으로 start/gather 실패.) 그럼에도 **report-dir**에는 메트릭 JSON·junit·generatedConfig 등이 정상 기록됨.

---

## 2. 결과 디렉터리 실제 구조

실행 후 **results/20260315-181750/** 아래 구조:

```
results/20260315-181750/
├── run.log                    # ClusterLoader2 stdout/stderr 전체 (tee)
└── clusterloader2/             # ClusterLoader2 --report-dir 출력
    ├── cl2-metadata.json       # 런 메타데이터 (본 run에서는 {})
    ├── generatedConfig_load.yaml   # 적용된 load 시나리오 설정 스냅샷
    ├── junit.xml               # JUnit 형식 테스트 결과
    ├── ClusterOOMsTracker_load_2026-03-15T18:20:42+09:00.json
    ├── MetricsForE2E_load_2026-03-15T18:20:42+09:00.json
    ├── PodStartupLatency_CreatePhasePodStartupLatency_load_2026-03-15T18:20:42+09:00.json
    ├── ResourceUsageSummary_load_2026-03-15T18:20:42+09:00.json
    ├── StatefulPodStartupLatency_CreatePhasePodStartupLatency_load_2026-03-15T18:20:42+09:00.json
    ├── StatelessPodStartupLatency_CreatePhasePodStartupLatency_load_2026-03-15T18:20:42+09:00.json
    └── SystemPodMetrics_load_2026-03-15T18:20:42+09:00.json
```

- **manifest.txt:** 이번 run에서는 ClusterLoader2가 Fatal로 종료한 뒤 스크립트가 즉시 exit하여 manifest 작성 단계까지 도달하지 않음. run-id·시나리오·오버라이드·report_dir은 위 경로와 run.log 상단 로그로 확인 가능.

---

## 3. ClusterLoader2가 실제로 생성한 메트릭 파일

| 파일명 패턴 | 설명 | SLI 후보와의 관계 |
|-------------|------|--------------------|
| `cl2-metadata.json` | 런 메타데이터 | 본 run에서는 `{}`. 메타데이터 수집 조건에 따라 비어 있을 수 있음. |
| `generatedConfig_load.yaml` | 적용된 config 스냅샷 | 환경·시나리오 재현용. |
| `junit.xml` | 테스트 스위트 결과 | pass/fail 개수·실패 스펙. |
| `PodStartupLatency_CreatePhasePodStartupLatency_load_*.json` | Create phase 파드 기동 레이턴시 (전체) | **Pod startup latency** SLI 후보에 직접 대응. |
| `StatelessPodStartupLatency_CreatePhasePodStartupLatency_load_*.json` | Stateless 파드 기동 레이턴시 | **Pod startup latency** (stateless 구간). |
| `StatefulPodStartupLatency_CreatePhasePodStartupLatency_load_*.json` | Stateful 파드 기동 레이턴시 | **Pod startup latency** (stateful 구간). |
| `ClusterOOMsTracker_load_*.json` | OOM 발생·과거 실패 목록 | **Container restarts / OOM** SLI 후보. |
| `SystemPodMetrics_load_*.json` | 시스템 파드별 restartCount 등 | **Container restarts** SLI 후보(재시작 횟수). |
| `ResourceUsageSummary_load_*.json` | 리소스 사용 요약 | 리소스 사용률 등. Phase duration·보조 해석에 활용 가능. |
| `MetricsForE2E_load_*.json` | E2E용 메트릭 (대용량) | 상세 메트릭; 필요 시 특정 SLI 보조. |

**이번 run에서 수집되지 않은 measurement (로그 기준):**

- **APIResponsivenessPrometheus / APIResponsivenessPrometheusSimple** — Prometheus 비활성화로 스킵.
- **InClusterNetworkLatency** — Prometheus 비활성화로 스킵.
- **NetworkProgrammingLatency** — Prometheus 비활성화로 스킵.
- **KubeProxyIptablesRestoreFailures** — Prometheus 비활성화로 스킵.
- **APIAvailability** — 별도 옵션 없음.
- **TestMetrics (SchedulingMetrics)** — kind에서 scheduler 엔드포인트 접근 실패로 start/gather 실패.

즉, **Prometheus가 꺼진 환경**에서는 API 레이턴시·네트워크 레이턴시·kube-proxy 실패 횟수 등은 산출되지 않으며, **Pod startup latency**, **Container restarts/OOM**, **SystemPodMetrics** 만 파일로 남는다.

---

## 4. ClusterLoader2 산출물 ↔ SLI 후보 매핑

| SLI 후보 (Architecture) | 본 run에서 산출 여부 | 실제 파일·키 (예시) |
|--------------------------|----------------------|----------------------|
| **Pod startup latency** | ✅ 수집됨 | `PodStartupLatency_*_load_*.json`, `StatelessPodStartupLatency_*_*.json`, `StatefulPodStartupLatency_*_*.json`. 내부 `dataItems[].data.Perc50|Perc90|Perc99`, `labels.Metric` (예: `pod_startup`, `schedule_to_watch`). 단위: **ms**. |
| **API call latency** | ❌ 미수집 | Prometheus 비활성화. 수집 시 `APIResponsivenessPrometheus*_*.json` 등. |
| **In-cluster network latency** | ❌ 미수집 | Prometheus 비활성화. |
| **Network programming latency** | ❌ 미수집 | Prometheus 비활성화. |
| **API availability** | ❌ 미수집 | 옵션 미사용. |
| **Container restarts** | ✅ 수집됨 | `SystemPodMetrics_*_*.json`: `pods[].containers[].restartCount`, `lastRestartReason`. `ClusterOOMsTracker_*_*.json`: `failures`, `ignored`, `past`. |
| **Kube-proxy iptables restore failures** | ❌ 미수집 | Prometheus 비활성화. |
| **Phase duration (선택)** | △ 부분 | TestMetrics/SchedulingMetrics 실패로 단계 타이머는 불완전. `ResourceUsageSummary` 등으로 구간 정보는 일부 활용 가능. |

---

## 5. 추출한 메트릭 예시 값

### 5.1 Pod startup latency (SLI 후보)

**파일:** `PodStartupLatency_CreatePhasePodStartupLatency_load_2026-03-15T18:20:42+09:00.json`

- **Metric: pod_startup** (파드 생성부터 준비까지):
  - Perc50: **8991.664 ms** (~9.0 s)
  - Perc90: **44904.985 ms** (~44.9 s)
  - Perc99: **64120.349 ms** (~64.1 s)
- **Metric: run_to_watch:**
  - Perc50: 1085.66 ms, Perc90: 2157.095 ms, Perc99: 4425.602 ms
- **Metric: schedule_to_watch:** Perc50: 2480.114 ms, Perc90: 9036.191 ms, Perc99: 9727.931 ms  
(위 값들은 `dataItems[].labels.Metric`로 구분됨.)

**파일:** `StatelessPodStartupLatency_CreatePhasePodStartupLatency_load_*.json`

- **Metric: pod_startup:**
  - Perc50: **5319.432 ms** (~5.3 s), Perc90: **10612.088 ms** (~10.6 s), Perc99: **10727.931 ms** (~10.7 s)

### 5.2 Container restarts / OOM (SLI 후보)

**파일:** `ClusterOOMsTracker_load_2026-03-15T18:20:42+09:00.json`

```json
{ "failures": [], "ignored": [], "past": [] }
```

- OOM 발생 없음.

**파일:** `SystemPodMetrics_load_2026-03-15T18:20:42+09:00.json`

- 시스템 파드(coredns, etcd, kube-apiserver, kube-proxy, kube-scheduler 등)의 `restartCount`가 모두 **0**, `lastRestartReason`은 `""`.

---

## 6. 현재 평가 가능한 것과 불명확한 것

### 6.1 이미 평가 가능한 것

- **Pod startup latency:**  
  - `PodStartupLatency_*` / `StatelessPodStartupLatency_*` / `StatefulPodStartupLatency_*` JSON에서 `dataItems`를 파싱해 `labels.Metric == "pod_startup"`인 항목의 Perc50/Perc90/Perc99(ms)를 읽으면 됨.  
  - SLO가 “P99 ≤ 5s” 형태로 정의되어 있으면, 위 예시(약 64s, 10.7s 등)와 비교해 pass/fail 판정 가능.
- **Container restarts / OOM:**  
  - `ClusterOOMsTracker_*.json`에서 `failures` 비어 있음 + `SystemPodMetrics_*.json`에서 `restartCount` 합계(또는 max)로 재시작 여부 판단 가능.  
  - “재시작 0회, OOM 0건” 같은 SLO는 이번 run 기준으로 만족.

### 6.2 아직 불명확하거나 미수집인 것

- **API call latency, In-cluster network latency, Network programming latency, Kube-proxy iptables restore failures:**  
  - Prometheus가 꺼진 현재 설정에서는 measurement가 스킵되어 **파일 자체가 생성되지 않음**.  
  - 이들 SLI를 평가하려면 Prometheus를 켜거나, 해당 measurement를 Prometheus 없이 지원하는 방식으로 설정을 바꿔야 함.
- **API availability:**  
  - 이번 run에서는 옵션 미사용으로 수집되지 않음. 수집 방법·파일 형식은 별도 확인 필요.
- **Phase duration (단계별 소요 시간):**  
  - TestMetrics(SchedulingMetrics) 실패로 스케줄링 구간 타이머는 없음. 다른 measurement에서 단계 구간 정보를 쓸 수 있는지는 추가 분석 필요.
- **cl2-metadata.json:**  
  - 본 run에서는 `{}`. 어떤 조건에서 채워지는지(버전·플래그 등)는 추가 확인 필요.

---

## 7. 해석 노트

- **실행 환경:** devcat `bin/cl2` + PROVIDER=local + 기본 kubeconfig로 **실제 ClusterLoader2 부하 테스트**가 수행됨. kind 3노드 클러스터에서 create phase 등이 실행되었고, Pod startup latency·OOM·시스템 파드 재시작 메트릭이 report-dir에 기록됨.
- **테스트 실패:** SchedulingMetrics(스케줄러 엔드포인트 접근 실패)로 전체 결과는 Fail이지만, **이미 수집된 measurement는 유효**하며 SLI 추출·SLO 평가에 사용 가능.
- **소규모 클러스터:** 노드 수가 3이므로 config의 `$IS_SMALL_CLUSTER`가 true. 공식 load config에서는 pod-startup-latency·scheduler-throughput **모듈**이 스킵되지만, **CreatePhasePodStartupLatency**는 reconcile-objects 단계에서 수집되는 것으로 보이며, 이번 run에서 PodStartupLatency 관련 JSON이 생성된 것과 일치함.
- **다음 단계 제안:**  
  (1) Pod startup latency에 대한 SLO(예: P99 ≤ 5s)를 정하고, 위 JSON에서 추출한 값으로 pass/fail 평가.  
  (2) Prometheus를 활성화하거나 대체 수단을 두어 API/네트워크/kube-proxy 관련 SLI 수집 후, 동일한 매핑 표를 보완.  
  (3) 04-engineering runbook의 “SLI measurement 추출” 절차를, 이번에 확인한 **실제 파일명·키 경로**로 구체화.

---

## 8. 안전 규칙 적용 실험 (Safe run with Prometheus)

Experiment Safety Rules(`.cursor/rules/devcat-experiment-safety.mdc`)에 따라 preflight, 안전 기본 옵션, 타임아웃을 적용한 뒤 devcat 실험을 재실행했다.

### 8.1 Preflight 결과

| 항목 | 결과 | 비고 |
|------|------|------|
| Kubeconfig 존재 | ✅ 통과 | `$HOME/.kube/config` 사용 |
| 클러스터 도달 | ✅ 통과 | `kubectl cluster-info` 성공 |
| API server 접근 | ✅ 통과 | `kubectl get ns` 성공 |
| Default StorageClass | ✅ 있음 | `standard (default)`, `ssd` (rancher.io/local-path 등) |
| ClusterLoader2 바이너리 | ✅ 있음 | `devcat/bin/cl2` 실행 가능 |

### 8.2 사용한 안전 실행 옵션

- **Prometheus 활성화:** `--enable-prometheus-server`
- **메모리 요청:** `--prometheus-memory-request=400Mi`
- **Prometheus PVC 비활성화:** CLI 플래그 없음 → **test override**로 처리. 두 번째 `--testoverrides`에 `projects/sample-cl2-sli-slo-analysis/05-experiment/cl2-override-pvc-disabled.yaml` 지정. 내용: `CL2_PROMETHEUS_PVC_ENABLED: false`
- **API server scrape 포트:** `--prometheus-apiserver-scrape-port=6443`

실제 실행 예:

```bash
cd devcat
PROVIDER=local KUBECONFIG="${KUBECONFIG:-$HOME/.kube/config}" \
perl -e 'alarm 300; exec @ARGV' -- ./bin/cl2 \
  --testconfig=scenarios/load/config.yaml \
  --testoverrides=overrides/ol-test.yaml \
  --testoverrides=/path/to/ai-team-lab/projects/sample-cl2-sli-slo-analysis/05-experiment/cl2-override-pvc-disabled.yaml \
  --report-dir=results/20260315-safe-prometheus/clusterloader2 \
  --enable-prometheus-server \
  --prometheus-memory-request=400Mi \
  --prometheus-apiserver-scrape-port=6443
```

### 8.3 타임아웃

- **적용 값:** **5분(300초)**. `perl -e 'alarm 300; exec @ARGV'`로 실행.
- **결과:** 5분 경과 시 프로세스가 종료됨. **gather 단계 전에 중단**되어, measurement JSON·junit.xml 등이 report-dir에 쓰이지 않음.
- Run ID: `20260315-safe-prometheus`. 산출된 파일: `clusterloader2/cl2-metadata.json`, `clusterloader2/generatedConfig_load.yaml`, `run.log` 뿐.

### 8.4 Prometheus 의존 SLI 메트릭 수집 여부

- **수집되지 않음.**  
  - 원인: **타임아웃으로 run이 gather 단계 전에 종료**됨. create 단계·TestMetrics 실패 후 로그 시점(18:58:24) 근처에서 5분 alarm에 의해 프로세스가 종료된 것으로 보임.
  - Prometheus 자체는 배포됨(PVC 비활성 오버라이드 적용). 로그 상 약 2분 30초 동안 "no endpoints available for service prometheus-k8s" 후 18:56:28에 다음 단계로 진행했으므로, **Prometheus Pod는 기동된 뒤** load 시나리오가 시작된 상태였다. 다만 **APIResponsivenessPrometheus 등 Prometheus 기반 measurement가 gather 단계까지 실행되지 않았기 때문에** 해당 SLI 메트릭 파일은 생성되지 않았다.

### 8.5 수집이 안 된 이유 정리

| 요인 | 설명 |
|------|------|
| **타임아웃** | 5분이면 load 시나리오 전체(create → scale/update → delete + gather)를 끝내기 부족. create 단계와 TestMetrics 실패 대기만으로도 시간이 소요되어, gather 전에 프로세스가 종료됨. |
| **Prometheus 설정** | PVC 비활성화로 Pending 없이 배포됨. 단, **kind는 API server 포트가 6443이 아님**(예: 127.0.0.1:51365). `--prometheus-apiserver-scrape-port=6443`은 kind 기본 포트와 맞지 않을 수 있어, 실제 수집 시 **apiserver scrape 실패** 가능성 있음. |
| **클러스터/환경** | TestMetrics(SchedulingMetrics)는 여전히 kind에서 scheduler 엔드포인트 접근 실패. Prometheus 기반 measurement와는 별개. |

### 8.6 Prometheus 의존 메트릭을 수집하려면 필요한 변경

1. **타임아웃 완화**  
   - 로컬에서 전체 시나리오를 돌리려면 **최소 15~20분** 정도 허용하거나, “Prometheus 기동 + load create + gather”만 돌리는 축소 시나리오를 쓰고 그에 맞춰 타임아웃 설정.  
   - 또는 조기 중단 규칙만 적용하고 전체 run 타임아웃은 더 길게 두기.

2. **Prometheus API server scrape 포트**  
   - kind는 API server가 **동적 포트**(예: 51365)를 씀. `--prometheus-apiserver-scrape-port`를 **실제 kind API server 포트**로 맞추거나, ClusterLoader2가 kind에서 해당 포트를 알 수 있도록 설정/오버라이드.  
   - 포트 불일치 시 apiserver 메트릭 스크랩 실패 → API responsiveness 등 Prometheus 기반 SLI 수집 실패 가능.

3. **실험 순서**  
   - Preflight → 안전 옵션(PVC 비활성 포함) → **충분한 타임아웃**으로 한 번 완주한 뒤, `results/<run-id>/clusterloader2/`에 `APIResponsivenessPrometheus*`, `InClusterNetworkLatency*` 등 **Prometheus 의존 measurement 파일** 생성 여부 확인.  
   - 생성되면 해당 파일과 SLI 후보 매핑을 04-engineering runbook 및 이 실험 노트에 반영.

---

## 9. Refinement 실험 (실제 API server 포트 + 긴 타임아웃)

이전 안전 실험에서 확인한 두 가지(타임아웃 부족, API server scrape 포트 불일치)를 반영해, **실제 API server 포트**를 사용하고 **Prometheus 대기·실험 타임아웃**을 늘려 한 번 더 실행했다.

### 9.1 실제 Kubernetes API server 포트 확인

- **방법:** 현재 kubeconfig에서 사용 중인 클러스터의 API server URL을 보고 포트를 추출했다.
- **명령:**
  ```bash
  kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}'
  ```
- **결과:** `https://127.0.0.1:51365` → **포트 51365**
- **추가 확인:** `kubectl cluster-info`에서도 control plane이 `https://127.0.0.1:51365`로 표시됨.
- **정리:** kind 등 로컬 클러스터는 6443이 아니라 **동적 포트**(여기서는 51365)를 쓰므로, Prometheus apiserver scrape에는 `--prometheus-apiserver-scrape-port=51365`를 사용해야 한다.

### 9.2 사용한 타임아웃

| Run | 실험 타임아웃 | Prometheus ReadyTimeout | 비고 |
|-----|----------------|--------------------------|------|
| 20260315-refinement-prometheus | 18분 (alarm 1080) | 기본 15분 | Prometheus 설정 단계에서 15분 후 "timed out waiting for the condition"으로 실패 |
| 20260315-refinement2-prometheus | 30분 (alarm 1800) | **25분** (`--prometheus-ready-timeout=25m`) | Prometheus 대기 시간 연장 후 재시도 |

### 9.3 실행 옵션 (refinement)

- **Safe local defaults:** `--enable-prometheus-server`, `--prometheus-memory-request=400Mi`, Prometheus PVC 비활성화(두 번째 `--testoverrides`: `cl2-override-pvc-disabled.yaml`).
- **실제 API server scrape 포트:** `--prometheus-apiserver-scrape-port=51365`
- **Refinement2에서만:** `--prometheus-ready-timeout=25m` 추가.

### 9.4 Run 완료 여부 및 gather/measurement 산출물

- **20260315-refinement-prometheus:**  
  - **완료되지 않음.** Prometheus 스택 설정 단계에서 **"Error while setting up prometheus stack: timed out waiting for the condition"**으로 종료.  
  - 약 2분 간격으로 "no endpoints available for service prometheus-k8s" 로그가 찍힌 뒤, 15분 경과 시 ReadyTimeout으로 실패.  
  - **gather 단계 미진입** → measurement JSON·junit.xml 등 **산출 없음**. `clusterloader2/`에는 생성된 파일 없음(또는 generatedConfig만 있을 수 있음).

- **20260315-refinement2-prometheus:**  
  - 30분 실험 타임아웃 + 25분 Prometheus 대기로 재실행.  
  - 이 run이 **gather 단계까지 진행했는지**는 run.log·`clusterloader2/` 내용으로 확인 필요.  
  - 만약 **Prometheus 의존 measurement 파일**(예: `APIResponsivenessPrometheus*_*.json`)이 **생성되었다면** → 해당 파일 목록과 예시 메트릭 값을 아래 9.5에 기록.  
  - **생성되지 않았다면** → 9.6(남은 블로커)에 반영.

### 9.5 Prometheus 의존 measurement 파일 생성 여부 및 예시 값

- **refinement 1차 run (20260315-refinement-prometheus):** Prometheus 설정 단계에서 15분 타임아웃으로 실패. **Prometheus 의존 measurement 파일 없음.** `clusterloader2/`에 산출 파일 없음.
- **refinement 2차 run (20260315-refinement2-prometheus):**  
  - `results/20260315-refinement2-prometheus/clusterloader2/` 확인 결과 **디렉터리 비어 있음**(generatedConfig, junit, measurement JSON 등 없음).  
  - run.log는 Prometheus 매니페스트 적용 직후(20:31:58)에서 끝나 있어, 2차 run도 **Prometheus 대기 중 종료되었거나 로그가 tee에 완전히 flush되지 않은 상태**로 보임.  
  - **APIResponsivenessPrometheus**, API availability, 기타 Prometheus 기반 measurement 관련 파일은 **생성되지 않음.**  
  - **예시 메트릭 값:** 수집된 Prometheus 의존 산출이 없어 기재할 예시 값 없음.

### 9.6 여전히 수집되지 않았을 때의 남은 블로커 및 최소 변경 제안

- **1차 refinement에서 확인된 블로커:**  
  - **Prometheus 스택이 15분 안에 Ready 상태가 되지 않음.**  
  - kind에서 이미지 풀·Pod 기동 등으로 **Prometheus가 15분 넘게 걸리는 것**으로 보임.  
  - 따라서 **"Prometheus-dependent metrics 수집 실패"의 직접 원인**은 **Prometheus setup 단계 타임아웃**이다.

- **최소한의 다음 변경:**  
  1. **`--prometheus-ready-timeout=25m` (또는 30m)** 로 늘리고, 실험 전체 타임아웃도 그에 맞춰 **25~30분** 이상 두어, Prometheus가 준비된 뒤 load 시나리오와 gather가 끝나도록 한다.  
  2. **동일 클러스터에서 Prometheus를 이미 한 번 기동해 둔 상태**에서, CL2가 기존 스택을 재사용할 수 있다면(또는 TearDown 비활성화 후 재실행) 대기 시간을 줄일 수 있으나, 현재 run은 매번 TearDown 후 재설치하므로 **ReadyTimeout 연장**이 가장 단순한 수정이다.  
  3. **API server scrape 포트:** 당시에는 51365를 사용했으나, §10 조사에서 **클러스터 내부**에서는 51365가 connection refused이고 **6443**을 써야 함이 확인됨. 51365는 호스트 전용 포트이므로, Prometheus(클러스터 내부)가 apiserver를 스크래핑하려면 **6443**으로 설정해야 함.

---

## 10. Prometheus 측정 준비 상태 조사 (Measurement Readiness Investigation)

전략 변경: 타임아웃만 늘리는 재시도 대신, **Prometheus 의존 메트릭 파이프라인이 실제로 준비되었는지**를 먼저 검증하는 절차를 도입하고, 현재 스크래핑 실패를 **연결/설정 문제**로 분석했다.

### 10.1 Prometheus 측정 준비 Precheck 절차

**전체 부하 실험(full load test)을 돌리기 전에** 아래를 순서대로 확인한다. 모두 통과해야 Prometheus 의존 SLI 수집이 가능하다고 판단하고, 하나라도 실패하면 부하 실험을 돌리지 않고 원인 수정을 먼저 한다.

| 단계 | 확인 항목 | 방법 | 통과 기준 |
|------|-----------|------|-----------|
| 1 | Prometheus가 실행 중인가? | `kubectl get pods -n monitoring -l app.kubernetes.io/name=prometheus` (또는 CL2가 배포한 네임스페이스/라벨) | Prometheus Pod가 Running. |
| 2 | 스크래핑 대상( scrape target )이 정의되어 있는가? | Prometheus ConfigMap 또는 `kubectl get servicemonitor` / 배포된 Prometheus config 확인 | API server용 target(또는 master IP + 포트)이 존재. |
| 3 | API server 메트릭 엔드포인트가 **Prometheus가 있는 위치(클러스터 내부)**에서 도달 가능한가? | 클러스터 내부 Pod에서 `curl -sk --connect-timeout 2 https://<MasterInternalIP>:<PORT>/metrics` 실행 (예: `172.25.0.10`) | 연결 거부(connection refused)가 아니어야 함. HTTP 403은 엔드포인트 존재 의미(인증은 별도). |
| 4 | 대상 포트·엔드포인트가 올바른가? | kind 등에서는 **호스트 포트(예: 51365)**와 **컨트롤 플레인 노드 내부 수신 포트(예: 6443)**가 다름. Precheck 시 두 포트 모두 시도해 봄. | 클러스터 **내부**에서 접근 시에는 **6443** 사용. 51365는 호스트 전용. |
| 5 | 기본 메트릭 쿼리가 성공하는가? | Prometheus Pod에서 `wget -qO- http://localhost:9090/api/v1/query?query=up` 또는 CL2가 사용하는 Prometheus API 호출 | 200 응답 및 유효한 JSON. |

**Precheck 실패 시:** full load 실험을 돌리지 않음. 실패한 단계의 원인(스크래핑 대상 미정의, 잘못된 포트, RBAC/인증 등)을 해결한 뒤 Precheck를 다시 실행.

### 10.2 현재 스크래핑 실패 분석 (Connection Refused)

**관찰:** Prometheus 의존 메트릭 수집이 되지 않고, 스크래핑 대상에 대한 **connection refused** (또는 이에 상응하는 “target connection refused”)가 발생한다. 실험에서는 API server 스크래핑 타깃이 `https://172.25.0.10:51365/metrics` 형태로 설정된 것으로 추정된다.

**실제 검증 (클러스터 내부 Pod에서):**

- `https://172.25.0.10:6443/metrics` → **HTTP 403** (연결 성공, 서버가 응답; 403은 인증/권한 이슈).
- `https://172.25.0.10:51365/metrics` → **연결 실패 (connection refused, HTTP 000)**.

**원인 분류:**  
이 실패는 **“wrong endpoint or wrong port”** 에 해당한다.

- **Wrong port:**  
  - **51365**는 kind가 **호스트**에 노출한 API server 포트이다. 클러스터 **내부** Pod(예: Prometheus)는 컨트롤 플레인 **노드**(172.25.0.10)에 접속한다. 해당 노드(컨테이너) **안**에서는 API server가 **6443**에서만 listen 한다.  
  - 따라서 Prometheus가 `MasterInternalIP:51365` 로 스크래핑하면, 노드의 51365 포트에는 아무것도 listen 하지 않아 **connection refused** 가 발생한다.
- **Wrong endpoint:**  
  - 엔드포인트 경로 `/metrics` 는 일반적으로 맞다. 포트만 6443으로 바꾸면 “target not listening” 문제는 해소된다.

**다른 카테고리와의 구분:**

- **Target not listening:** 6443에서는 listen 하고 응답(403)하므로 해당 없음. 51365는 “대상이 그 포트에서 listen 하지 않음”에 해당.
- **Cluster-specific networking:** 동일 클러스터 내 Pod에서 6443은 도달 가능하므로, 네트워크 격리 문제라기보다 **포트 설정 오류**에 가깝다.
- **Scrape configuration issue:** 스크래핑 설정에서 **포트를 호스트 포트(51365)로 지정한 것**이 문제. 클러스터 내부 스크래핑일 때는 **6443**을 써야 함.
- **RBAC/certificate:** 403은 TLS 핸셰이크 후 응답이므로, “connection refused” 와는 별개. 인증/권한은 Precheck 5 또는 Prometheus scrape 시 별도 처리(서비스 어카운트 등)로 해결.

**결론:**  
실패의 직접 원인은 **스크래핑 포트 오류**이다. `--prometheus-apiserver-scrape-port=51365` 는 **호스트 기준** 포트이며, **클러스터 내부**에서 동작하는 Prometheus가 사용할 포트로는 부적절하다. **클러스터 내부에서는 6443** 을 사용해야 한다.

### 10.3 Prometheus 의존 SLI 수집 환경 준비 여부

- **현재 상태:** **준비되지 않음.**  
  - 이유: API server 스크래핑 타깃이 **클러스터 내부 관점에서 잘못된 포트(51365)** 로 설정되어 있어, Prometheus가 apiserver 메트릭을 수집하지 못함.  
  - 그 결과 APIResponsivenessPrometheus, API availability, InClusterNetworkLatency, NetworkProgrammingLatency, KubeProxy 등 **Prometheus 기반 measurement** 가 수집되지 않음.

### 10.4 최소 필요 수정 (다음 full experiment 전)

1. **스크래핑 포트 수정**  
   - kind(및 동일 구조의 로컬 클러스터)에서 **Prometheus가 클러스터 내부**에서 API server를 스크래핑할 때는 **`--prometheus-apiserver-scrape-port=6443`** 을 사용해야 한다.  
   - 51365는 **호스트에서** API server에 접근할 때만 사용한다. CL2/배포 매니페스트가 `MasterInternalIP` + 포트로 스크래핑 타깃을 만들므로, 이 포트는 **노드(컨트롤 플레인) 내부 수신 포트 6443** 이어야 한다.
2. **Precheck 통과 후에만 full run**  
   - 위 10.1 Precheck(특히 3·4번: 클러스터 내부에서 `https://<MasterInternalIP>:6443/metrics` 도달 가능 여부)를 실행해 통과한 뒤에만, 긴 타임아웃의 full load 실험을 실행한다.
3. **타임아웃은 Precheck 통과 후에만 의미 있음**  
   - 스크래핑이 실패하는 상태에서 타임아웃만 늘려도 Prometheus 의존 메트릭은 수집되지 않는다. 먼저 포트 수정 및 Precheck 통과가 선행되어야 함.

### 10.5 권장 다음 액션 (엔지니어링 결정용)

| 옵션 | 내용 | 권장 조건 |
|------|------|-----------|
| **A. Prometheus 설정 수정 후 full load 실험** | `--prometheus-apiserver-scrape-port=6443` 으로 변경하고, 10.1 Precheck를 통과한 뒤 full load 실험 재실행. | Prometheus 기반 SLI(API 레이턴시, 네트워크 레이턴시 등)를 v1에서 수집하려는 경우. |
| **B. Prometheus 설정 수정 우선 (full run 보류)** | Precheck 절차를 문서/스크립트로 고정하고, 포트 6443 적용 후 Precheck만 실행해 “측정 준비 완료”를 확인. full load는 별도 결정 후 실행. | “측정 파이프라인 준비”를 먼저 검증하고 싶을 때. |
| **C. v1 수용은 non-Prometheus SLI 세트로** | Prometheus 의존 메트릭 없이, 이미 수집 가능한 **Pod startup latency**, **Container restarts/OOM**, **SystemPodMetrics** 등만으로 v1 수용 기준을 정의. | Prometheus 수정을 당장 하지 않고, 우선 비-Prometheus SLI로 마일스톤을 맞추려는 경우. |

**권장:**  
- **Prometheus 기반 SLI 수집이 목표라면:** **A** 또는 **B**를 선택. 반드시 **`--prometheus-apiserver-scrape-port=6443`** 적용 및 **Precheck 통과 후**에만 full load 또는 긴 타임아웃 실험을 진행한다.  
- **당장 full load 재실행은 하지 않는다.** Precheck 없이 타임아웃만 늘린 실험은 하지 않음(blind long-timeout experiment 금지).

### 10.6 요약

- **Precheck:** full load 전에 Prometheus 실행, 스크래핑 타깃 존재, **클러스터 내부에서** API server 메트릭 엔드포인트 도달 가능 여부, 올바른 포트(6443), 기본 쿼리 성공을 확인.
- **현재 스크래핑 실패:** **Wrong port.** 172.25.0.10:51365 는 클러스터 내부에서 수신하지 않음 → connection refused. 172.25.0.10:6443 은 수신 중(403 응답).
- **환경 준비:** 수정 전에는 **미준비**. 포트를 6443으로 바꾸고 Precheck를 통과하면 준비된 것으로 볼 수 있음.
- **다음 단계:** Prometheus scrape 포트를 6443으로 수정 → Precheck 실행 → 통과 시에만 full load 실험 또는 B 옵션으로 “측정 준비 완료” 문서화. Review 단계는 아직 진행하지 않음.

---

## 11. Precheck 전용 검증 실행 (Prometheus Measurement Readiness)

full load 실험을 다시 돌리기 전에, **올바른 API server scrape 포트(6443)** 로 Prometheus 측정 파이프라인이 실제로 준비되었는지 **Precheck만** 수행했다.

### 11.1 사용한 설정

- **API server scrape 포트:** `--prometheus-apiserver-scrape-port=6443` (클러스터 내부 포트).
- **CL2 실행:** Prometheus 활성화, PVC 비활성화(두 번째 `--testoverrides`: `cl2-override-pvc-disabled.yaml`), 25분 alarm. Run ID: `20260315-precheck-6443`.  
- **검증 시점:** 기존 `monitoring` 네임스페이스에 이전 run에서 배포된 Prometheus가 이미 동작 중이었고, API server 타깃만 **51365**로 설정되어 있어 down 상태였음.  
- **6443 반영 방법:** Precheck 검증을 위해 `monitoring` 네임스페이스의 **master** Endpoints 리소스에서 apiserver 포트를 51365 → **6443**으로 수동 패치 후, Prometheus가 해당 타깃을 다시 스크래핑하도록 함.  
  - `kubectl patch endpoints master -n monitoring --type='json' -p='[{"op":"replace","path":"/subsets/0/ports/0/port","value":6443}]'`

### 11.2 Precheck 결과 요약

| 단계 | 확인 항목 | 결과 | 비고 |
|------|-----------|------|------|
| 1 | Prometheus Pod 실행 | ✅ 통과 | `prometheus-k8s-0`, `prometheus-operator-*` (monitoring NS) Running. |
| 2 | API server용 스크래핑 타깃 존재 | ✅ 통과 | `serviceMonitor/monitoring/master/0` (apiserver), master/1·2 (scheduler, controller-manager). |
| 3 | 타깃 포트가 클러스터 내부 기준 올바른가 (6443) | ✅ 통과 | 패치 전: `https://172.25.0.10:51365/metrics` → down. 패치 후: `https://172.25.0.10:6443/metrics` 사용. |
| 4 | Prometheus 타깃 상태 확인 가능·건강 상태 | ✅ 통과 | `api/v1/targets` 로 확인. master/0 패치 후 **up**. |
| 5 | 기본 메트릭 쿼리 성공 | ✅ 통과 | `query?query=up` → status success, job `master` 172.25.0.10:6443 => 1. `query?query=apiserver_request_total` → 382건 수집. |

### 11.3 타깃 실패 시 남은 블로커 판별

- **패치 전:** API server 타깃(master/0)이 `https://172.25.0.10:51365/metrics` 로 설정되어 **down**, lastError에 connection refused(dial tcp ... 51365: connect: connection refused) 유사 메시지.
- **패치 후(6443):** 동일 타깃 **up**, apiserver 메트릭(`apiserver_request_total`) 정상 수집.
- **결론:** 이 환경에서의 남은 블로커는 **포트 오류(51365)** 뿐이었음. **authentication / authorization** 문제는 관찰되지 않음(타깃 up, 메트릭 수집됨). TLS/서비스 어카운트는 Prometheus Operator 구성이 그대로 유효한 것으로 보임.

### 11.4 정리 및 권장

- **Prometheus 의존 측정 수집 가능 여부:** 올바른 포트(6443)가 적용되면 **수집 가능**하다고 판단함. Precheck 1~5 모두 통과.
- **다음 full load 실험:**  
  - **진행 가능.**  
  - 조건: CL2 실행 시 **`--prometheus-apiserver-scrape-port=6443`** 사용하고, Prometheus가 **master** Endpoints(또는 동등한 스크래핑 소스)에서 apiserver 포트를 **6443**으로 받도록 해야 함.  
  - kind/CL2 매니페스트가 기본으로 호스트 포트를 넣는다면, 배포 후 Precheck로 master/0 타깃이 6443으로 up인지 한 번 확인한 뒤 full load를 돌리는 것을 권장.
- Review 단계는 수행하지 않음.

---

*Experimenter 역할 산출물. 실제 ClusterLoader2 run 기준으로 갱신함. Review 단계는 수행하지 않음.*
