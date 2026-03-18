# decision-output-spec.md (출력 JSON 스펙)

## 출력 정의 예시
```json
{
  "system_state": "SAFE",
  "reason": "CAT PASS and monitoring safe",
  "source": {
    "cat": "PASS",
    "monitoring": "safe"
  }
}
```

## 필드 의미
- `system_state`: `SAFE | DEGRADED | FAIL`
- `reason`: 사람이 읽는 짧은 근거(룰 매칭 결과)
- `source`: 입력이 무엇이었는지 원본 값을 그대로 보관

## 왜 이 정도면 충분한가
- Decision Layer는 “최종 수렴”이 목적이므로 복잡한 evidence 구조는 필요 없다.
- 운영자가 문제를 추적할 때도 `source`만 보면 입력이 무엇이었는지 즉시 확인된다.

