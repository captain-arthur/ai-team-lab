# Grafana 빌드 가이드(최소 지표 기반, Stat 중심)

이 문서만 보고도 “9개 Stat 패널”을 구현할 수 있게 구성한다.
가능하면 각 패널은 단일 쿼리 결과로 `state_score`(0/1/2)를 반환하도록 만든다.

## 1) Row 구성
- Row 1(지금 안전한가?): Panel 1~3
- Row 2(내일도 안전한가?): Panel 4~6
- Row 3(지금 무엇을 해야 하는가?): Panel 7~9

## 2) 패널 배치(권장)
- 각 Row는 컴포넌트 3개를 좌→우 순서로 동일하게 고정
  - istio-ingressgateway / nodelocaldns / coredns

## 3) 패널 타입
- 모든 패널은 `Stat panel` 중심
- 그래프 금지(단, sparkline은 선택으로 1개까지 허용)

## 4) 상태 판정 결과(권장 구현 방식)
쿼리 결과를 다음처럼 통일한다.
- state_score 정의
  - safe/stable = 0
  - warning/risk = 1
  - fail/high = 2
- Stat 패널에는
  - 표시 값: state_score의 라벨(`safe`/`warning`/`fail` 등) 또는 state_score 자체
  - Threshold: 0→green, 1→yellow, 2→red

## 5) threshold 설정 방식(Stat)
- Threshold 1: `0` 경계(예: 0은 green)
- Threshold 2: `1` 경계(예: 1~1.999 yellow, 2 red)

예시(개념)
- color:
  - 0 → green
  - 1 → yellow
  - 2 → red

## 6) Stat 라벨 표시
- `Value mappings`(Grafana 기능)을 사용하면 state_score 0/1/2를 safe/warning/fail로 바로 바꿀 수 있다.

## 7) (선택) PromQL 예시(최소 수준, “패턴”)
환경마다 metric 이름/라벨이 다르므로, 아래는 “계산 패턴”을 보여준다.

### 공통: error_rate 계산 패턴
- `error_rate = 실패 요청 수 / 전체 요청 수`
- 실패 요청 수는 “에러 코드/타임아웃/에러 플래그”의 정의에 맞게 라벨만 조정

### 공통: p95_latency_ms 패턴
- `p95_latency_ms = histogram bucket 기반 quantile(0.95)` 또는 latency SLI의 p95

### 공통: health_ready 패턴
- ready/health 상태를 숫자화(예: ready면 1, 아니면 0)하거나,
- “준비성 불안정”을 나타내는 boolean/카운트 기반으로 상태 판정에 반영

### state_score 만드는 패턴(예)
1) fail_cond / warning_cond를 먼저 정의(불리언)
2) `state_score = fail_cond*2 + warning_cond*1`

예시 스켈레톤(그대로 복사하지 말고 패턴으로 사용)
```promql
state_score =
  (fail_cond)*2
  + (warning_cond)*1
```

## 8) 각 패널 쿼리 책임(정확히 무엇을 Stat 패널에 넣을지)
- Panel 1~3: 현재 state_score(0/1/2)를 반환
- Panel 4~6: 내일 추세 기반 state_score(0/1/2)를 반환
- Panel 7~9: Q1/Q2 상태를 입력으로 action class 텍스트를 반환(또는 score 기반으로 간단히 색/우선순위만 표시)

Action 패널은 “그래프 없이” 텍스트를 보여주는 것이 목표다.

## 9) 운영자가 5초 내 이해하도록 하는 UX 규칙
- Stat에는 숫자 대신 라벨을 우선 표시
  - 예: `fail`, `warning`, `safe`
- 패널 제목에 반드시 질문을 함께 적기
  - 예: “지금 안전한가”
- 색상은 무조건 동일 규칙 사용(green/yellow/red)

## 10) Q3(Action)용 최소 action_score 매핑(권장)
Q3 패널은 그래프 없이 텍스트를 Stat로 보여주는 것이 목표이므로, numeric score 하나로 만든 뒤 Value mapping으로 문자열을 표시한다.

- action_score 0 (green)
  - “예방 확인”
- action_score 1 (yellow)
  - “원인 분류 1차 확인”
- action_score 2 (red)
  - “원인 분리 후 즉시 완화 준비”

각 컴포넌트 패널 7~9의 “첫 행동 1줄”은 action_score에 대응되는 문자열로 그대로 연결한다.

