# Runbook 연결 모델(패널 상태 → 다음 행동)

이 문서는 “대시보드와 runbook을 어떻게 이어 붙일지”를 실무적으로 고정한다.

## 1) 대시보드와 runbook 연결 방식(실용 판단)

후보:
- 패널 단위 연결
- Row 단위 연결
- 컴포넌트 단위 runbook tree 연결

최종 추천: **패널 단위 연결 + 컴포넌트 단위 runbook tree(보조)**

- 패널 단위 연결
  - “이 패널이 safe/watch/fail 중 무엇인지”에 따라 runbook_id가 바뀐다.
  - 패널에서 상태가 바뀌면 곧바로 행동이 바뀌어 운영자가 탐색하지 않는다.
- 컴포넌트 단위 runbook tree(보조)
  - runbook은 컴포넌트별로 묶고, 상태에 따라 first-step checklist가 달라지도록 한다.

이 혼합은 과도한 runbook 분기를 만들지 않으면서도 “다음 행동”을 정확히 고정한다.

## 2) runbook 설계 원칙

### 언제 사람이 읽어야 하는가
- 화면에서 상태가 `fail` 또는 `risk-high`로 바뀌는 순간
- Action(Row 3)에서 runbook 링크를 누를 때

### 얼마나 짧아야 하는가
- runbook은 “5분 체크 + 다음 10분 체크 + 흔한 원인 후보 + 완화 후보 + escalation”으로 고정
- 길이의 상한을 두고, “자세한 설명”은 제외(링크/문서로만 둔다)

### 어떤 형태가 적절한가
- first action / investigation checklist / mitigation flow
  - first action: 지금 즉시 볼 것 1~3개
  - investigation: 원인 분기(2~4개)
  - mitigation: 완화 후보(2~3개)
  - escalation: 자동/사람 중 어디로 넘길지

## 3) runbook 구조 제안(최소 필드)
각 runbook은 최소 아래를 포함한다.

- `runbook_id`
- `대상 컴포넌트`
- `트리거 상태`(예: fail, watch, risk-high)
- `첫 번째 확인 항목`(1~3개)
- `두 번째 확인 항목`(2~4개)
- `완화/조치 후보`(2~3개)
- `escalation 기준`(언제 상위/다른 팀으로 넘길지)

## 4) 최종 추천(대시보드에서 어디까지 보여줄지)
- 대시보드(Action 패널)는 runbook을 “열기 위한 버튼/링크”와 “첫 행동 요약”까지만 보여준다.
- 세부 체크리스트는 runbook starter를 참조하도록 분리한다.

즉:
- Row 3: 요약 + 링크
- runbook starter(본 문서의 07-engineering): 상세 체크리스트

이 구조가 “별도의 문서를 찾아 읽는 부담”을 줄이면서도, runbook이 자연스러운 다음 단계가 되게 한다.

