# decision-model.md (통합 Decision Layer 최소 모델)

## 1) 왜 Decision Layer가 필요한가
현재는 `CAT(PASS/FAIL)`과 `Monitoring(safe/warning/fail)`이 각각 따로 존재하고, 최종 판단은 사람이 두 결과를 머릿속에서 합쳐야 한다.

Decision Layer는 이 “합치기”를 시스템 규칙으로 고정해서, 매번 같은 입력에 항상 같은 최종 상태로 수렴시키는 역할을 한다.

## 2) CAT와 Monitoring의 역할 차이
- CAT = 검증
  - 기대 동작이 실제로 만족되는지 확인한다.
  - 출력은 `PASS / FAIL`로 “기대 미충족”이 있는지의 권위 있는 판정이다.
- Monitoring = 관측
  - 현재 운영 상태가 안전한지(품질/에러/지연 등) 관측한다.
  - 출력은 `safe / warning / fail`로 “현재 리스크 수준”을 의미한다.

## 3) 최종 상태 정의
Decision Layer가 출력하는 최종 상태는 3개만 사용한다.

## SAFE
- 정의: 시스템이 기대 동작을 만족하고( CAT PASS ) 동시에 운영 리스크가 안전 수준이다( Monitoring safe ).

## DEGRADED
- 정의: 기대 동작은 만족하지만( CAT PASS ) 운영 상태가 경고/불안정 수준이다( Monitoring warning 또는 fail ).
- DEGRADED는 “즉시 실패”가 아니라 “품질 저하로 위험이 커졌음”을 뜻한다.

## FAIL
- 정의: 기대 동작이 만족되지 않는다( CAT FAIL ).
- CAT FAIL은 Monitoring 상태와 무관하게 최종 FAIL을 만든다.

