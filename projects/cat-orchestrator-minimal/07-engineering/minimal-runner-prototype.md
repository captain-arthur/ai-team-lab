# 최소 러너 프로토타입(실행 가능한 설계 예시)

## 목표
- `cat run job.yaml`이 “실제로” tool adapter를 호출하고,
- `output.dir/cat-result.json`을 생성하는 흐름을 보여준다.
- 복잡한 플러그인/스케줄러는 만들지 않는다.

## 1) Job YAML 예시
`job.yaml`
```yaml
test_name: ginkgo-stability-window
tool: ginkgo
scenario:
  type: custom
  config:
    package: ./08-engineering/custom-cat-scenarios/01-stability-window
    env:
      WARMUP_MS: "120"
      STABILITY_MS: "400"
      SLO_READY_MAX_MS: "250"
      SLO_ERROR_RATE_MAX: "0.0"
slo:
  - metric: ready_ok
    condition: "== true"
output:
  dir: ./results/run-001/
selected_sli: ["ready_latency_ms", "error_rate"]
```

## 2) Runner 동작(단순)
- 러너는 adapter를 호출한다.
- adapter는 exit code와 cat-result.json을 만든다.
- runner는 cat-result.json만 확인하고 끝낸다.

## 3) 최소 Runner 의사코드(간단 CLI)
아래 예시는 “구조”만 보여주기 위한 최소 Python 스켈레톤이다.

```python
import os, sys, subprocess, yaml, json

def run_ginkgo_adapter(job):
    outdir = job["output"]["dir"]
    pkg = job["scenario"]["config"]["package"]
    env = os.environ.copy()
    env.update(job["scenario"]["config"].get("env", {}))

    os.makedirs(outdir, exist_ok=True)
    # ginkgo는 테스트 코드 내부에서 results 파일을 output dir 규약으로 남기도록 설계 가능
    p = subprocess.run(["go", "test", pkg, "-v"], env=env)
    exit_code = p.returncode

    cat_path = os.path.join(outdir, "cat-result.json")
    if not os.path.exists(cat_path):
        raise RuntimeError(f"missing cat-result.json: {cat_path}")
    return exit_code

def main():
    job_path = sys.argv[1]
    job = yaml.safe_load(open(job_path))
    tool = job["tool"]
    if tool == "ginkgo":
        code = run_ginkgo_adapter(job)
    else:
        raise NotImplementedError(tool)
    sys.exit(code)

if __name__ == "__main__":
    main()
```

## 4) adapter 책임 예시(변환/수집만)
- k6 adapter: `k6-summary.json`을 읽어 `cat-result.json`을 만든다
- ginkgo adapter: 테스트가 만든 `cat-result.json`을 수집한다(없으면 최소 변환만)
- CL2 adapter: 측정 JSON에서 selected_sli 추출 후 `cat-result.json` 생성

## 5) 결과 생성 예시
- `output.dir/`
  - `raw/` 또는 직접 파일들(예: `k6-summary.json`, `ginkgo` artifacts, `cl2 measurement json`)
  - `cat-result.json` (표준 스키마 고정)

