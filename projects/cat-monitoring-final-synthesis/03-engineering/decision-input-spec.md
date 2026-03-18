# decision-input-spec.md (입력 JSON 스펙)

## 1) CAT 입력 예시
```json
{
  "final_pass_fail": "PASS"
}
```

## 2) Monitoring 입력 예시
```json
{
  "state": "safe"
}
```

## 3) 왜 이 정도 입력이면 충분한가
Decision Layer는 “지금 안전한가?”의 최종 상태만 내리면 된다.

- CAT는 기대 동작 만족 여부를 단 하나(PASS/FAIL)로 준다.
- Monitoring은 현재 리스크 수준을 단 하나(safe/warning/fail)로 준다.
- 두 값만 결합하면 최종 상태 3개(SAFE/DEGRADED/FAIL)를 완전히 결정할 수 있다.

추가 메타데이터(예: 테스트 이름, 패널 상태, metric 값)는 “이유(reason)”를 풍부하게 만들 수는 있지만, 최종 상태 결정에는 필요 없다.

