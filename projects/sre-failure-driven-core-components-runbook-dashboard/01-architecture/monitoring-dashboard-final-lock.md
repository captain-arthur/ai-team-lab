# monitoring-dashboard-final-lock.md (Monitoring만으로 완결)

## Monitoring 프로젝트 최종 목표
이 대시보드는 Monitoring 관측 지표만으로 운영자가 아래 3개 질문에 즉시 답할 수 있게 하는 “최소 판정 화면”이다.

- Q1. 지금 안전한가?
- Q2. 내일도 안전한가?
- Q3. 무엇을 해야 하는가?

## CAT와의 관계(필수)
- 이 문서/대시보드는 외부 검증 결과(PASS/FAIL) 또는 기대 동작 검증 입력을 사용하지 않는다.
- 모든 판정은 Monitoring의 관측 지표(istio-ingressgateway / nodelocaldns / coredns Prometheus 지표)만으로 결정된다.

## Row 구조(고정)
- Row 1: `지금 안전한가?` (컴포넌트별 safe / warning / fail)
  - Panel 1: istio-ingressgateway
  - Panel 2: nodelocaldns
  - Panel 3: coredns
- Row 2: `내일도 안전한가?` (컴포넌트별 stable / risk / high)
  - Panel 4: istio-ingressgateway
  - Panel 5: nodelocaldns
  - Panel 6: coredns
- Row 3: `무엇을 해야 하는가?` (컴포넌트별 prevention / root-cause / immediate)
  - Panel 7: istio-ingressgateway
  - Panel 8: nodelocaldns
  - Panel 9: coredns

## Panel 구조(고정 규칙)
- 각 Panel은 “상태 1개(색/라벨)”만 반환한다.
- Row 3은 Row 1/Row 2의 상태를 번역해서 “첫 행동 1줄”만 제공한다.
- Runbook은 번역된 첫 행동에 대한 연결만 제공한다(과도한 설명/설계 금지).

## 왜 외부 검증 입력 없이도 이 대시보드가 완결되는가
- 외부 검증 입력은 “기대 동작 충족 여부 검증”이지만, 이 대시보드의 질문은 “현재/미래 운영 관측 상태” 자체다.
- 따라서 Monitoring 지표로 safe/warning/fail과 stable/risk/high를 직접 만들 수 있으며, Row 3은 그 분류를 행동으로 번역만 하면 된다.

## Threshold(오류율) 최종 수치 규칙(모든 패널 공통)
## Q1: 지금 안전한가?(오류율_5m)
- safe: error_rate(5m) <= 1%
- warning: error_rate(5m) > 1% AND error_rate(5m) <= 3%
- fail: error_rate(5m) > 3%

## Q2: 내일도 안전한가?(오류율_1h)
- stable: error_rate(1h) <= 1%
- risk: error_rate(1h) > 1% AND error_rate(1h) <= 3%
- high: error_rate(1h) > 3%

## 오류율 정의(패널 스코프 동일)
- istio-ingressgateway: `5xx 요청 비율 = 5xx / 전체 요청`
- nodelocaldns: `SERVFAIL 비율(노드 로컬 DNS) = SERVFAIL 응답 / 전체 응답`
- coredns: `SERVFAIL 비율(coredns) = SERVFAIL 응답 / 전체 응답` (단, node-local-dns 제외)

## Action(Row 3) 명령형 1줄(10~15자 수준)
- Panel 7(istio)
  - 예방 확인(0): `준비성 확인`
  - 원인 분류(1): `업스트림 확인`
  - 즉시 완화(2): `즉시 완화`
- Panel 8(nodelocaldns)
  - 예방 확인(0): `forward 확인`
  - 원인 분류(1): `업스트림 분리`
  - 즉시 완화(2): `즉시 조치`
- Panel 9(coredns)
  - 예방 확인(0): `준비성 확인`
  - 원인 분류(1): `forward 분리`
  - 즉시 완화(2): `즉시 조치`

## Runbook 연결 규칙(최소)
- action=즉시 완화(2)일 때만 runbook을 연결한다.
- Row 1(fail) 또는 Row 2(high)인 경우 action=2가 된다.
- runbook id는 컴포넌트별 기존 규칙을 그대로 사용한다.

