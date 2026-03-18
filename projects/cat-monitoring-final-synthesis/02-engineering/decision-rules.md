# decision-rules.md (최소 규칙 4개)

아래 규칙은 **충돌 없이 전체 경우를 커버**하도록 구성했다. 규칙 수는 4개이며, 6개 이내다.

## Rule 1
IF `cat.final_pass_fail == "FAIL"`
THEN `system_state = "FAIL"`

## Rule 2
IF `cat.final_pass_fail == "PASS"` AND `monitoring.state == "fail"`
THEN `system_state = "DEGRADED"`

## Rule 3
IF `cat.final_pass_fail == "PASS"` AND `monitoring.state == "warning"`
THEN `system_state = "DEGRADED"`

## Rule 4
IF `cat.final_pass_fail == "PASS"` AND `monitoring.state == "safe"`
THEN `system_state = "SAFE"`

## 중복/충돌 처리
- Rule 1이 가장 우선이다(FAIL은 monitoring과 무관).
- PASS 케이스는 monitoring이 safe/warning/fail 중 하나이므로 Rule 2~4 중 정확히 1개만 매칭된다.

