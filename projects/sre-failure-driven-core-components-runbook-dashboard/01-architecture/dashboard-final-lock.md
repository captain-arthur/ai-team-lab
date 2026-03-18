# 최소 지표 기반 운영 대시보드 최종 확정(3개 컴포넌트)

## 최종 Row 구조 확정
- Row 1: `지금 안전한가?` (컴포넌트별 safe / warning / fail)
  - Panel 1: istio-ingressgateway
  - Panel 2: nodelocaldns
  - Panel 3: coredns
- Row 2: `내일도 안전한가?` (컴포넌트별 stable / risk / high)
  - Panel 4: istio-ingressgateway
  - Panel 5: nodelocaldns
  - Panel 6: coredns
- Row 3: `무엇을 해야 하는가?` (컴포넌트별 예방 확인 / 원인 분류 / 즉시 완화 준비)
  - Panel 7: istio-ingressgateway
  - Panel 8: nodelocaldns
  - Panel 9: coredns

## Row의 역할(설명용이 아니라 구현 규칙)
- Row 1: “현재 사용자 영향 품질”을 기준으로 safe/warning/fail을 만든다.
- Row 2: Row 1의 신호가 “악화 추세 + 회복 실패 가능성”으로 승격되는지로 stable/risk/high를 만든다.
- Row 3: Row 1/Row 2 상태를 그대로 번역해 “첫 행동 1줄”만 보여준다(새 분석/새 metric 금지).

## 왜 이 구조로 충분한가
- 질문 3개가 곧 Row 3개이며, Row마다 판정 스코프(현재/추세/행동)가 고정된다.
- 각 패널은 Stat으로 “상태 라벨 1개”만 반환하므로, 운영자가 5초 내에 색/라벨로 분류할 수 있다.
- Action(Row 3)은 Row 1/Row 2에서 만들어진 실패(high/fail) 우선순위를 그대로 사용한다.

## 컴포넌트별 4개 골든 시그널 확정(최소)
아래 4개는 각 컴포넌트의 Row 1/Row 2에서 실제 상태 계산에 사용되는 신호다.

### istio-ingressgateway 골든 시그널 4개
1. `error_rate_5m` (5xx 비율, 현재 사용자 실패 품질)
2. `p95_latency_ms_5m` (현재 지연 품질)
3. `request_rate_5m` (트래픽 압박/포화 가능성)
4. `ready_ok` (게이트웨이 준비성)

### nodelocaldns 골든 시그널 4개
1. `dns_error_rate_5m` (SERVFAIL 비율, 이름 해석 실패 품질)
2. `dns_latency_p95_ms_5m` (DNS 지연 품질)
3. `forward_error_rate_5m` (upstream forward SERVFAIL 비율, upstream 도달/품질 신호)
4. `ready_ok` (노드 로컬 DNS 준비성)

### coredns 골든 시그널 4개
1. `dns_error_rate_5m` (SERVFAIL 비율)
2. `dns_latency_p95_ms_5m` (DNS 지연 품질)
3. `ready_ok` (coredns 준비성)
4. `restarts_10m` (coredns 재시작 누적, 불안정/회복 실패 신호)

