# monitoring-dashboard-implementation-ready.md (Final Lock, Monitoring only)

이 문서는 9개 패널을 그대로 Grafana로 옮길 수 있게, 판정 기준을 **수치 범위**로 고정하고 action 문장을 **명령형 1줄**로 축약한다.

---

## 공통 Threshold(오류율 기반)

### Q1: 지금 안전한가?(error_rate_5m)
- safe: error_rate(5m) <= 1%
- warning: error_rate(5m) > 1% AND error_rate(5m) <= 3%
- fail: error_rate(5m) > 3%

### Q2: 내일도 안전한가?(error_rate_1h)
- stable: error_rate(1h) <= 1%
- risk: error_rate(1h) > 1% AND error_rate(1h) <= 3%
- high: error_rate(1h) > 3%

---

## Runbook 연결(최소)
- action=즉시(2)일 때만 runbook 연결
- Row 1 상태가 fail(2) 또는 Row 2 상태가 high(2)면 action=즉시(2)가 된다.

컴포넌트별 runbook id:
- `istio-ingressgateway`
  - `runbook://istio-ingressgateway/failure/fail`
  - `runbook://istio-ingressgateway/risk/high`
- `nodelocaldns`
  - `runbook://nodelocaldns/failure/fail`
  - `runbook://nodelocaldns/risk/high`
- `coredns`
  - `runbook://coredns/failure/fail`
  - `runbook://coredns/risk/high`

---

## Panel 1: istio-ingressgateway - 지금 안전한가?
- 질문(Row): Row 1
- 입력 지표: `error_rate_5m = 5xx / 전체 요청`
- 판정
  - safe: error_rate_5m <= 1%
  - warning: error_rate_5m > 1% AND <= 3%
  - fail: error_rate_5m > 3%
- 첫 행동 1줄: `준비성 확인`

---

## Panel 2: nodelocaldns - 지금 안전한가?
- 질문(Row): Row 1
- 입력 지표: `error_rate_5m = SERVFAIL / 전체 DNS 응답 (node-local-dns)`
- 판정
  - safe: error_rate_5m <= 1%
  - warning: error_rate_5m > 1% AND <= 3%
  - fail: error_rate_5m > 3%
- 첫 행동 1줄: `forward 확인`

---

## Panel 3: coredns - 지금 안전한가?
- 질문(Row): Row 1
- 입력 지표: `error_rate_5m = SERVFAIL / 전체 DNS 응답 (coredns, node-local-dns 제외)`
- 판정
  - safe: error_rate_5m <= 1%
  - warning: error_rate_5m > 1% AND <= 3%
  - fail: error_rate_5m > 3%
- 첫 행동 1줄: `준비성 확인`

---

## Panel 4: istio-ingressgateway - 내일도 안전한가?
- 질문(Row): Row 2
- 입력 지표: `error_rate_1h = 5xx / 전체 요청 (1시간)`
- 판정
  - stable: error_rate_1h <= 1%
  - risk: error_rate_1h > 1% AND <= 3%
  - high: error_rate_1h > 3%
- 첫 행동 1줄: `준비성 확인`

---

## Panel 5: nodelocaldns - 내일도 안전한가?
- 질문(Row): Row 2
- 입력 지표: `error_rate_1h = SERVFAIL / 전체 DNS 응답 (node-local-dns, 1시간)`
- 판정
  - stable: error_rate_1h <= 1%
  - risk: error_rate_1h > 1% AND <= 3%
  - high: error_rate_1h > 3%
- 첫 행동 1줄: `forward 확인`

---

## Panel 6: coredns - 내일도 안전한가?
- 질문(Row): Row 2
- 입력 지표: `error_rate_1h = SERVFAIL / 전체 DNS 응답 (coredns, node-local-dns 제외, 1시간)`
- 판정
  - stable: error_rate_1h <= 1%
  - risk: error_rate_1h > 1% AND <= 3%
  - high: error_rate_1h > 3%
- 첫 행동 1줄: `준비성 확인`

---

## Panel 7: istio-ingressgateway - 무엇을 해야 하는가?
- 질문(Row): Row 3
- action=0(예방): `준비성 확인`
- action=1(원인): `업스트림 확인`
- action=2(즉시): `즉시 완화`
- runbook 연결(최소)
  - action=2일 때만:
    - Row1 fail 또는 Row2 high이면:
      - `runbook://istio-ingressgateway/failure/fail` 또는 `runbook://istio-ingressgateway/risk/high`

---

## Panel 8: nodelocaldns - 무엇을 해야 하는가?
- 질문(Row): Row 3
- action=0(예방): `forward 확인`
- action=1(원인): `업스트림 분리`
- action=2(즉시): `즉시 조치`
- runbook 연결(최소)
  - action=2일 때만:
    - Row1 fail 또는 Row2 high이면:
      - `runbook://nodelocaldns/failure/fail` 또는 `runbook://nodelocaldns/risk/high`

---

## Panel 9: coredns - 무엇을 해야 하는가?
- 질문(Row): Row 3
- action=0(예방): `준비성 확인`
- action=1(원인): `forward 분리`
- action=2(즉시): `즉시 조치`
- runbook 연결(최소)
  - action=2일 때만:
    - Row1 fail 또는 Row2 high이면:
      - `runbook://coredns/failure/fail` 또는 `runbook://coredns/risk/high`

---

## Self Final Check(문서 기준)
- 모든 판정이 수치 범위로 정의됨
- warning/risk 정의가 명확함(1%~3% 포함)
- action 문장은 명령형 1줄로 축약됨
- 외부 검증 입력 없이 Monitoring만으로 성립
- 구현은 9개 패널 고정으로 시작 가능

