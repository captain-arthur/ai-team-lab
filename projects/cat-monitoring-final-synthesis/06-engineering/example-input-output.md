# example-input-output.md (실제 예시 4개)

Decision Engine 입력은 2개 JSON이다.
- `cat_result.json`: `{"final_pass_fail": "PASS|FAIL"}`
- `monitoring_state.json`: `{"state": "safe|warning|fail"}`

출력은 `decision.json`이며 아래처럼 `system_state`가 결정된다.

## 예시 1
입력
```json
{ "final_pass_fail": "FAIL" }
```
```json
{ "state": "safe" }
```
→ 출력
```json
{ "system_state": "FAIL" }
```

## 예시 2
입력
```json
{ "final_pass_fail": "PASS" }
```
```json
{ "state": "fail" }
```
→ 출력
```json
{ "system_state": "DEGRADED" }
```

## 예시 3
입력
```json
{ "final_pass_fail": "PASS" }
```
```json
{ "state": "warning" }
```
→ 출력
```json
{ "system_state": "DEGRADED" }
```

## 예시 4
입력
```json
{ "final_pass_fail": "PASS" }
```
```json
{ "state": "safe" }
```
→ 출력
```json
{ "system_state": "SAFE" }
```

