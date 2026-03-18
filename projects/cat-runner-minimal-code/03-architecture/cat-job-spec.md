# CAT Job Spec (공통)

CAT Runner는 “하나의 공통 Job Spec”만 이해한다. Job Spec은 tool을 추상화하지 않고, tool별 추가 값은 **최소한의 확장 필드(`scenario.config`)**로만 넣는다.

## 최소 예시(job.yaml)
```yaml
test_name: http-basic
tool: k6
entry: ./examples/k6/http-basic.js

scenario:
  type: http
  target: https://test.k6.io/
  load_model: arrival   # tool 실행에 필요한 최소 힌트(도구가 이해할 수 있게)
  rps: 40               # 최소 주입 파라미터
  duration: 20s

slo:
  latency_p95_ms: 450
  error_rate: 0.01

output:
  dir: ./results/http-basic
```

## 공통 필드
`test_name`, `tool`, `entry`, `scenario`, `slo`, `output.dir`

## tool별 최소 확장 방식
- `scenario` 아래의 표준 필드만 사용해 기본 의미를 고정한다.
- tool이 추가로 필요로 하는 값은 `scenario.config` 같은 **명시적 “도구 전용 확장” 맵**으로 넣는다.
- Runner는 `scenario.config`를 해석하지 않고 adapter에 그대로 전달한다.

## 불필요한 추상화 금지
- tool별 옵션을 공통 인터페이스로 “전부” 끌어오지 않는다.
- 공통 Job Spec은 CAT이 책임져야 하는 최소한의 정보만 남기고, 나머지는 adapter의 몫으로 둔다.

