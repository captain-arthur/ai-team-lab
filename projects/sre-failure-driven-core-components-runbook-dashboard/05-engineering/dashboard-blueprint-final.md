# 최종 대시보드 청사진(바로 Grafana로 옮길 수 있는 수준)

이 대시보드는 “관측”이 아니라 “판정 → 경고/위험 → 즉시 행동(runbook 링크)”을 목표로 한다.

## 1) 최종 대시보드 구조
총 3 Row, Row별 3개 패널(컴포넌트당 1개 패널씩)

1. Row 1: `지금 안전한가(Failure / Today)`
   - 목적: 오늘 사용자 영향(실패 증거) 유무를 즉시 판정
   - Row 목적: safe / watch / fail 분류(컴포넌트별)
   - 패널 수: 3
2. Row 2: `내일도 안전한가(Confirmation / Tomorrow)`
   - 목적: 오늘의 상태가 내일 운영 창에서 악화될 가능성인지 승격
   - Row 목적: stable / risk / high 분류(컴포넌트별)
   - 패널 수: 3
3. Row 3: `지금 무엇을 해야 하는가(Action / Next runbook)`
   - 목적: Row1/Row2 상태에서 바로 runbook_id로 연결되는 “다음 행동” 제공
   - Row 목적: 상태 기반 Action class + runbook 링크 + 첫 5분 체크 요약
   - 패널 수: 3

전체 패널 수: 9

## 2) 대상별 공통 패널 상태 정의(모든 컴포넌트에 동일 규칙 적용)

### Row 1(오늘): safe / watch / fail
- safe: 실패 증거가 없음(성공 품질이 유지)
- watch: 실패 증거는 약하거나 특정 패턴(간헐/일부 경로)으로만 나타남
- fail: 실패 증거가 지속/확대(사용자 영향이 현실로 관측)

### Row 2(내일): stable / risk / high
- stable: 오늘 품질이 유지되고 악화 추세가 관측되지 않음
- risk: 품질이 흔들리거나 조기 악화 추세가 “완만하게” 존재
- high: 악화 추세가 명확하거나(회복 실패 포함) 내일 운영 창에 영향이 예상

### Row 3(Action): 상태 → next action → runbook
- Action class는 컴포넌트별 runbook 트리의 루트 상태로 매핑한다.

## 3) 컴포넌트별 Row 1 패널(지금 안전한가)

### Panel 1: `istio-ingressgateway - 지금 안전한가`
- panel 목적: 외부 사용자 관측 실패를 가장 먼저 드러냄
- 어떤 질문에 답하는가: “지금 안전한가?”
- panel 유형: `State timeline + Stat 카드(3상태)`
- 입력 지표(표현 원칙): “사용자 관측 실패(성공률/타임아웃/지연)”와 “게이트웨이 준비성(ready/health)”을 묶어 사용
- 판정 기준 초안
  - safe: 성공이 유지되고 타임아웃/오류가 낮으며 게이트웨이 준비성이 안정
  - watch: 성공률이 소폭 하락하거나 지연 p95가 기준치 근접/간헐 초과
  - fail: 타임아웃/오류율이 지속 상승하거나 준비성이 흔들려 트래픽 영향이 명확
- 이상 시 첫 행동(첫 5분)
  - upstream 연결/라우팅 실패 유형부터 분류(게이트웨이 vs upstream)
  - 게이트웨이 readiness/리소스 압박을 확인
- 연결 runbook
  - `runbook://istio-ingressgateway/failure/fail`

### Panel 2: `nodelocaldns - 지금 안전한가`
- panel 목적: 이름 해석 실패가 사용자 서비스 통신 실패로 연쇄되는 지점 포착
- 어떤 질문에 답하는가: “지금 안전한가?”
- panel 유형: `State timeline + Stat 카드(3상태)`
- 입력 지표(표현 원칙): DNS 질의 성공/실패, 타임아웃/응답 지연(캐시/포워딩 관점), 로컬 DNS 준비성
- 판정 기준 초안
  - safe: 질의 실패/타임아웃이 낮고 로컬 DNS가 ready 상태를 유지
  - watch: 실패가 증가하거나 노드 편차가 나타남(일부 노드 영향)
  - fail: 질의 실패/타임아웃이 확산되거나 응답 지연이 지속 상승
- 이상 시 첫 행동
  - upstream(coredns) 도달성 확인
  - nodelocaldns 포드/데몬셋 상태와 재시작 징후 확인
- 연결 runbook
  - `runbook://nodelocaldns/failure/fail`

### Panel 3: `coredns - 지금 안전한가`
- panel 목적: 클러스터 DNS 품질의 “원천 실패”를 조기 드러냄
- 어떤 질문에 답하는가: “지금 안전한가?”
- panel 유형: `State timeline + Stat 카드(3상태)`
- 입력 지표(표현 원칙): DNS 에러 계열(타임아웃/SERVFAIL 유사)/응답 지연, coredns 헬스/준비성, 최근 재시작/헬스 플랩
- 판정 기준 초안
  - safe: 헬스/ready 안정, 에러율 낮고 응답 지연이 기준 이내
  - watch: 에러 유형 중 특정 계열이 늘거나 지연이 완만히 상승
  - fail: 에러율 급상승/지연 지속 또는 준비성 불안정과 함께 사용자 영향 관측
- 이상 시 첫 행동
  - 로그/헬스 이벤트로 “원인 분류(업스트림/설정/리소스)” 시작
  - forward/upstream 품질을 먼저 확인
- 연결 runbook
  - `runbook://coredns/failure/fail`

## 4) 컴포넌트별 Row 2 패널(내일도 안전한가)

### Panel 4: `istio-ingressgateway - 내일도 안전한가`
- panel 목적: 오늘의 이상이 내일 운영 창으로 전이될 조짐을 승격
- 어떤 질문에 답하는가: “내일도 안전한가?”
- panel 유형: `Time series + Risk state Stat(stable/risk/high)`
- 입력 지표(표현 원칙): 실패/타임아웃/지연의 최근 추세와 회복 실패 여부, 게이트웨이 준비성 안정성
- 판정 기준 초안
  - stable: 최근 창에서 회복이 이루어지고 실패 증거가 감소
  - risk: 실패/지연이 완만히 상승하고 회복이 지연
  - high: 회복 실패 + 반복 패턴(특정 구간/경로/도메인 편차)이 확실
- 이상 시 첫 행동
  - 완화 준비(타임아웃/라우팅/용량 대응) 우선 검토
- 연결 runbook
  - `runbook://istio-ingressgateway/risk/high` (risk/high 공통 첫 행동)

### Panel 5: `nodelocaldns - 내일도 안전한가`
- panel 목적: 로컬 DNS가 누적적으로 악화되는 조기 경보 제공
- 어떤 질문에 답하는가: “내일도 안전한가?”
- panel 유형: `Time series(노드 편차) + Risk state Stat`
- 입력 지표(표현 원칙): 노드별 성공률/지연 분산 확대, upstream 의존성 악화 징후, 캐시 효율 저하 경향
- 판정 기준 초안
  - stable: 노드 편차가 작고 성공/지연이 안정
  - risk: 편차가 늘고 upstream 의존성이 악화(회복이 느림)
  - high: 실패가 확산되며 내일 운영 창 영향이 예상
- 이상 시 첫 행동
  - upstream reachability/coredns 헬스를 우선 확인
- 연결 runbook
  - `runbook://nodelocaldns/risk/high`

### Panel 6: `coredns - 내일도 안전한가`
- panel 목적: coredns의 포화/불안정이 내일로 전이되는 위험을 포착
- 어떤 질문에 답하는가: “내일도 안전한가?”
- panel 유형: `Time series + Risk state Stat`
- 입력 지표(표현 원칙): 에러/지연의 누적 추세, 헬스 플랩, 리소스 압박 징후
- 판정 기준 초안
  - stable: 에러/지연이 감소 또는 안정
  - risk: 에러/지연이 회복되지 않고 완만히 증가
  - high: 헬스 불안정/재시작 반복 + 지연 누적
- 이상 시 첫 행동
  - 내일 전이 위험이면 즉시 완화(스케일/설정 롤백) 계획 시작
- 연결 runbook
  - `runbook://coredns/risk/high`

## 5) 컴포넌트별 Row 3 패널(지금 무엇을 해야 하는가)

Row 3는 “상태 → next action → runbook 링크”를 한 화면에 끝내기 위한 패널이다.

### Panel 7: `istio-ingressgateway - Action & Runbook`
- panel 목적: 첫 5분 체크를 고정하고 runbook을 즉시 열기
- 어떤 질문에 답하는가: “지금 무엇을 해야 하는가?”
- panel 유형: `Table(상태별 next action + runbook 링크)`
- 입력 지표: Row1/Row2 상태(컴포넌트별 safe/watch/fail + stable/risk/high)
- 판정 기준 초안(매핑 규칙)
  - fail면: `failure/fail` runbook
  - watch 또는 risk면: `risk/high` 또는 `risk/investigate` 중 선택(공통 first action: 게이트웨이 vs upstream 분류)
- 상태값 예시
  - safe + stable: runbook 링크는 “예방 확인”만 제공(없어도 됨)
  - watch + risk: 조사 runbook 링크 제공
  - fail + high: fail runbook 링크 제공
- 이상 시 첫 행동
  - upstream 분류(엔드포인트/라우팅/인증 오류 계열) 1차 체크
- 연결 runbook
  - `runbook://istio-ingressgateway/*`

### Panel 8: `nodelocaldns - Action & Runbook`
- panel 목적: DNS 실패의 “원천”을 upstream(coredns) vs 로컬로 분리
- 어떤 질문에 답하는가: “지금 무엇을 해야 하는가?”
- panel 유형: `Table`
- 판정 기준 초안(매핑)
  - fail면 로컬 안정성(runbook: nodelocaldns/failure/fail) 먼저
  - risk-high면 upstream 확인을 먼저 포함(nodelocaldns/risk/high)
- 이상 시 첫 행동
  - upstream 도달성 확인 → 로컬 포드 상태 확인(순서 고정)
- 연결 runbook
  - `runbook://nodelocaldns/*`

### Panel 9: `coredns - Action & Runbook`
- panel 목적: coredns 실패 원인 분류(업스트림/설정/리소스) → 즉시 완화로 연결
- 어떤 질문에 답하는가: “지금 무엇을 해야 하는가?”
- panel 유형: `Table`
- 판정 기준 초안(매핑)
  - fail/high: fail runbook
  - watch/risk: investigate runbook(로그/구성 변경/업스트림 품질 1차 체크)
- 이상 시 첫 행동
  - 로그에서 원인 분류 후, 원인 계열에 맞는 완화 후보로 이동
- 연결 runbook
  - `runbook://coredns/*`

