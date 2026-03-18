# decision-layer-review.md (너무 단순하지 않은가?)

## 1) 너무 단순한가?
단순하다. 하지만 단순함이 장점인 이유가 있다.
- “최종 상태 수렴”이 목적이고, 필요한 입력이 딱 두 개(CAT의 검증 결과, Monitoring의 관측 상태)로 제한되기 때문이다.
- rule engine/복잡 추론을 금지하는 요구사항과도 정확히 맞는다.

## 2) 실무에서 충분한가?
충분한 경우가 많다.
- CAT가 기대 동작 만족을 권위 있게 말해주고(PASS/FAIL),
- Monitoring이 현재 운영 리스크 수준을 말해주기 때문이다(safe/warning/fail).

다만 “CAT PASS여도 FAIL을 즉시 막을 수 있나?” 같은 더 정교한 정책(예: CAT PASS여도 일정 지표가 임계치면 FAIL)은 이 최소 버전에서 다루지 않는다.

## 3) 어떤 경우에 깨지는가(리스크)
1. CAT가 “거짓 PASS”를 낼 때
   - CAT의 기대 동작 정의가 실제 운영 위험을 충분히 커버하지 못하면, Monitoring이 safe가 아닌 한 DEGRADED로는 내려가지만 FAIL로 즉시 전환되지 않는다.
2. Monitoring 상태의 매핑 오류
   - safe/warning/fail 산정이 흔들리면 system_state도 흔들린다.
3. “CAT FAIL인데 monitoring은 safe” 같은 경우를 FAIL로 고정하는 정책이 조직의 기대와 다를 때
   - 하지만 이 정책은 안전을 우선하므로(기대 미충족이면 무조건 FAIL) 조직 정책 결정에 따라 조정 필요다.

## 4) 확장하려면 어떻게 해야 하는가(하지만 지금은 최소 유지)
확장은 rule 수를 늘리는 방향이 아니라, 입력 스키마를 풍부하게 하여 “단일 규칙”의 근거를 바꾸는 방식이 더 안전하다.
- 예: Monitoring 입력을 `state`만이 아니라 `cause`(DNS/Latency/Errors 등)까지 포함시키면, Rule 2/3에서 reason만 더 세분화 가능
- 예: CAT 입력에 `scope`나 `severity`를 추가하면 CAT FAIL을 즉시 FAIL로 내릴지 DEGRADED로 내릴지 정책을 1줄로 조정 가능

현재는 3개 상태(SAFE/DEGRADED/FAIL)와 4개 rule을 유지하는 것이 목적이므로, 확장은 나중 단계로 미룬다.

