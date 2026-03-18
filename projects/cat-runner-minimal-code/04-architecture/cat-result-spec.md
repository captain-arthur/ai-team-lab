# CAT 표준 결과 스키마(cat-result.json) 최소 고정

CAT Runner의 표준 출력은 `cat-result.json` 단 하나로 고정한다.

## 최소 스키마(예시)
```json
{
  "test_name": "",
  "tool": "",
  "scenario_type": "",
  "selected_sli": {},
  "slo_result": {},
  "final_pass_fail": "",
  "exit_code": 0,
  "raw_result": {
    "format": "",
    "path": ""
  },
  "timestamp": ""
}
```

## selected_sli 구조(명명 규칙 + unit 규칙)
- key naming rule: `latency_p95_ms`, `error_rate`, `throughput_rps` 처럼 “지표 이름 + 단위/계산 형태”를 key로 고정한다.
- unit rule:
  - `latency_p95_ms`: 밀리초(ms)
  - `error_rate`: 실패 비율(0~1)
  - `throughput_rps`: 초당 처리량(RPS)

## slo_result 구조(왜 필요한가)
- `slo_result`는 CAT이 “재판정”하는 용도가 아니라, adapter가 뽑은 **measured 값과 SLO 임계값의 비교 결과**를 evidence로 남기기 위한 용도다.
- 예: `slo_result.latency_p95_ms.ok` 처럼 `ok: true/false`를 남긴다.

## raw_result.format / naming rule
- `raw_result.format`: `json`, `xml`, `text` 등 raw 포맷을 명시한다.
- `raw_result.path`: tool별 raw 산출물 파일 위치(파일 경로 문자열)를 남긴다.

## PASS/FAIL 권위(중요)
- CAT은 재판정하지 않는다.
- `final_pass_fail` 값과 `exit_code`는 **tool 실행 결과(프로세스 exit code)** 를 그대로 기록한다.
- 즉 `exit_code == 0`이면 `final_pass_fail == "PASS"`, 아니면 `"FAIL"`로 기록한다.

