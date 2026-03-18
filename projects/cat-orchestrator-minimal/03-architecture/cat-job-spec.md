# CAT Job Spec(최소, 단순, 확장 제한)

## 1) 공통 필드 정의(고정)
```yaml
test_name: <string>
tool: k6 | cl2 | ginkgo
scenario:
  type: http | custom | cluster
  target: <string>           # http인 경우(ingress 엔드포인트 등)
  config: {}                 # tool/시나리오별 구체 파라미터(자유도는 제한)
slo:
  - metric: <string>        # selected_sli 키와 1:1 매핑되는 이름
    condition: "< 300"      # 단순 문자열(어댑터가 해석)
  - metric: <string>
    condition: "< 0.01"
output:
  dir: ./results/<run-id>/
selected_sli: []            # tool이 계산해줄 SLI 키 목록(최소 키만)
```

## 2) tool별 최소 확장(불필요한 범용성 금지)
- `k6`
  - `scenario.type=http`일 때
  - `scenario.config.script_name`: 사용할 k6 스크립트 이름(예: `http-basic`)
  - `scenario.config.env`: k6가 읽을 env 키/값(최소)
- `ginkgo`
  - `scenario.type=custom`일 때
  - `scenario.config.package`: 테스트 패키지 경로(예: `./08-engineering/custom-cat-scenarios/01-stability-window`)
  - `scenario.config.env`: 테스트에 주입할 env(예: `SLO_LATENCY_MAX_MS`)
- `cl2`
  - `scenario.type=cluster`일 때
  - `scenario.config.config_yaml`: CL2 config 파일 경로(예: `config.yaml`)
  - `scenario.config.overrides_yaml`: (옵션) 오버라이드 파일 경로

## 3) scenario 주입 규칙(단순/명확)
- CAT은 `scenario`를 adapter에 그대로 전달한다.
- adapter는 `scenario.type`에 따라 필요한 입력만 사용한다.

