# 9개 패널 최소 지표 스펙(질문 1개 → 판정 → 1차 행동)

규칙
- 패널은 반드시 “하나의 질문”만 답한다.
  - Row 1: `지금 안전한가?` (safe / warning / fail)
  - Row 2: `내일도 안전한가?` (stable / risk / high)
  - Row 3: `지금 무엇을 해야 하는가?` (action class + 1차 행동 1줄)
- 각 패널의 “사용 지표”는 최대 3개다.
- 표시 방식은 Stat 중심(그래프 금지, sparkline은 선택).
- 첫 행동은 1줄만 제공한다(상세 runbook 금지).

색상 규칙(모든 Stat에 동일)
- green: safe / stable
- yellow: warning / risk
- red: fail / high

표기 예시
- Q1(Q now): `safe | warning | fail`
- Q2(Q tomorrow): `stable | risk | high`

---

## Panel 1: istio-ingressgateway - 지금 안전한가
1. 질문
- 지금 안전한가?
2. 사용하는 지표(최대 3)
- error_rate (현재 실패 비율)
- p95_latency_ms (현재 지연 품질)
- health_ready (현재 준비성)
3. 판정 로직
- safe: error_rate 낮고 p95_latency_ms 기준 이내이며 health_ready 정상
- warning: error_rate 또는 p95_latency_ms가 기준 근접/소폭 초과(단, health_ready 정상)
- fail: error_rate 또는 p95_latency_ms가 지속적으로 악화되고 health_ready가 흔들리거나(또는 error_rate가 기준 크게 초과)
4. 왜 이 지표만으로 충분한가
- ingress는 “사용자 관측 품질”이 곧 영향이므로 실패/지연/준비성을 최소로 묶는 것이 가장 빠른 판정 근거다.
5. 표시 방식
- Stat 1개(상태 문자열 safe/warning/fail) + 색상
- (선택) 10분 sparkline 1개는 “악화/회복” 확인용 보조만 허용
6. 첫 행동(1줄)
- upstream 연결 실패/라우팅 오류 가능성을 먼저 분리(게이트웨이 vs 타깃)

---

## Panel 2: nodelocaldns - 지금 안전한가
1. 질문
- 지금 안전한가?
2. 사용하는 지표(최대 3)
- dns_timeout_rate (현재 timeout 비율)
- dns_latency_p95_ms (현재 지연 품질)
- upstream_reachability (coredns 도달성)
3. 판정 로직
- safe: timeout_rate 낮고 p95 지연 기준 이내이며 upstream_reachability 정상
- warning: timeout_rate 또는 p95 지연이 기준 근접/소폭 상승(대부분 upstream_reachability는 유지)
- fail: timeout_rate가 기준을 지속 초과하거나, upstream_reachability가 저하되어 연쇄 실패가 나타남
4. 왜 이 지표만으로 충분한가
- nodelocaldns의 운영 목적은 “이름 해석의 성공과 짧은 지연”이며, 최소한 timeout/지연/upstream 도달성만으로 사용자 영향 전이를 빠르게 판정할 수 있다.
5. 표시 방식
- Stat 1개(상태 문자열 safe/warning/fail) + 색상
6. 첫 행동(1줄)
- coredns 도달성(upstream)이 원인인지 먼저 확인

---

## Panel 3: coredns - 지금 안전한가
1. 질문
- 지금 안전한가?
2. 사용하는 지표(최대 3)
- dns_error_rate (현재 에러 비율: SERVFAIL 유사)
- dns_latency_p95_ms (현재 지연 품질)
- health_ready (coredns 준비성)
3. 판정 로직
- safe: dns_error_rate 낮고 p95 지연 기준 이내이며 health_ready 정상
- warning: dns_error_rate 또는 p95 지연이 기준 근접/소폭 상승(health_ready 정상)
- fail: dns_error_rate 또는 p95 지연이 지속 상승하고 health_ready가 흔들림(또는 에러율이 기준 크게 초과)
4. 왜 이 지표만으로 충분한가
- coredns는 upstream 기준점이므로 에러/지연/준비성만으로 “지금 사용자 영향이 발생했는가”를 가장 직접적으로 판정한다.
5. 표시 방식
- Stat 1개(상태 문자열 safe/warning/fail) + 색상
6. 첫 행동(1줄)
- forward/upstream 품질부터 확인(에러 유형이 반복되는지)

---

## Panel 4: istio-ingressgateway - 내일도 안전한가
1. 질문
- 내일도 안전한가?
2. 사용하는 지표(최대 3)
- error_rate_trend (최근 구간 악화/회복)
- p95_latency_trend (최근 구간 악화/회복)
- recovery_signal (최근 회복이 관측되는지)
3. 판정 로직
- stable: 최근 구간에서 error_rate/p95_latency가 회복 또는 안정(회복이 관측)
- risk: 악화 징후가 있으나 회복이 지연(완만하게 지속)
- high: 회복 실패 + 반복 악화가 명확(다음 운영 창 영향 예상)
4. 왜 이 지표만으로 충분한가
- 내일 판정은 “현재 상태의 지속성”이 핵심이므로 trend/회복 신호 3개가 최소다.
5. 표시 방식
- Stat 1개(stable/risk/high) + 색상
6. 첫 행동(1줄)
- 타임아웃/라우팅/용량 대응의 “사전 완화” 준비 여부를 결정

---

## Panel 5: nodelocaldns - 내일도 안전한가
1. 질문
- 내일도 안전한가?
2. 사용하는 지표(최대 3)
- dns_timeout_rate_trend
- dns_latency_p95_trend
- upstream_reachability_trend
3. 판정 로직
- stable: timeout/지연이 회복하거나 안정, upstream 도달성이 유지
- risk: timeout/지연이 완만히 악화 + upstream 도달성이 흔들림(회복 지연)
- high: upstream reachability 저하가 지속 + 실패가 확산되는 추세
4. 왜 이 지표만으로 충분한가
- nodelocaldns는 “로컬에서 버티는 능력”이 중요하므로, 로컬 실패 추세와 upstream 도달성 추세만으로 내일 위험을 가장 짧게 추정할 수 있다.
5. 표시 방식
- Stat 1개(stable/risk/high) + 색상
6. 첫 행동(1줄)
- 내일 운영 창 기준으로 upstream 확인/완화 준비를 시작

---

## Panel 6: coredns - 내일도 안전한가
1. 질문
- 내일도 안전한가?
2. 사용하는 지표(최대 3)
- dns_error_rate_trend
- dns_latency_p95_trend
- health_ready_flap_trend
3. 판정 로직
- stable: 에러/지연이 회복 또는 안정, 준비성 플랩이 없음
- risk: 에러/지연이 완만히 증가 또는 회복 지연, 플랩 빈도는 낮음
- high: 에러/지연이 누적 악화 + 준비성 플랩이 반복(내일 영향 확률 증가)
4. 왜 이 지표만으로 충분한가
- coredns 장애는 재시도/큐잉으로 영향이 누적되기 때문에 에러/지연/준비성 불안정의 조합이 내일 예측의 최소다.
5. 표시 방식
- Stat 1개(stable/risk/high) + 색상
6. 첫 행동(1줄)
- forward/upstream 품질 또는 리소스 완화(스케일/재시작 등) 준비

---

## Panel 7: istio-ingressgateway - 지금 무엇을 해야 하는가
1. 질문
- 지금 무엇을 해야 하는가?
2. 사용하는 지표(최대 3)
- Panel 1 상태(safe/warning/fail) 1개(입력 역할)
- Panel 4 상태(stable/risk/high) 1개(참조 역할)
- (선택) health_ready의 현재 값(있다면) 1개
3. 판정 로직
- fail: 1차 확인은 “게이트웨이 준비성 → upstream/라우팅 분리” 순서
- warning: 1차 확인은 “로그에서 실패 유형 분류 → upstream/라우팅 확인”
- safe: 즉시 조치 없음(예방 확인만)
- risk/high 참조는 fail/warning의 행동을 우선순위로 가중(같은 행동 우선 수행)
4. 왜 이 지표만으로 충분한가
- Q3는 새로운 metric 해석이 아니라 Q1/Q2 상태를 행동으로 번역하는 단계이므로 입력은 상태만으로 충분하다.
5. 표시 방식
- Stat 1개(action class) + 색상(예: safe=green, warning=yellow, fail=red)
- (구현) action_score(0/1/2)를 Grafana Value mapping으로 문자열로 표시
6. 첫 행동(1줄)
- fail이면 upstream/라우팅 분리부터 시작

---

## Panel 8: nodelocaldns - 지금 무엇을 해야 하는가
1. 질문
- 지금 무엇을 해야 하는가?
2. 사용하는 지표(최대 3)
- Panel 2 상태(safe/warning/fail) 1개
- Panel 5 상태(stable/risk/high) 1개
- upstream_reachability 현재 값(가능하면) 1개
3. 판정 로직
- fail: upstream reachability 확인 → 로컬 준비성/재시작 징후 확인
- warning: upstream 도달성 재검증 → 노드 편차/집중 여부 확인
- safe: 예방 확인(추세는 Row2에서만 관찰)
- risk/high면 fail/warning 행동의 우선순위를 높임
4. 왜 이 지표만으로 충분한가
- Q3는 행동 번역 단계이므로 상태 + upstream reachability만 최소 입력으로 충분하다.
5. 표시 방식
- Stat 1개(action class) + 색상
- (구현) action_score(0/1/2)를 Grafana Value mapping으로 문자열로 표시
6. 첫 행동(1줄)
- upstream이 원인인지부터 분리

---

## Panel 9: coredns - 지금 무엇을 해야 하는가
1. 질문
- 지금 무엇을 해야 하는가?
2. 사용하는 지표(최대 3)
- Panel 3 상태(safe/warning/fail) 1개
- Panel 6 상태(stable/risk/high) 1개
- dns_error_rate 현재 값(가능하면) 1개
3. 판정 로직
- fail/high: 1차 확인은 “coredns 준비성/플랩 → forward/upstream 분리 → 포화 징후 확인” 순서
- warning/risk: 로그/에러 유형에서 upstream vs 설정/리소스 분기 후, 최소 완화 후보 선정
- safe/stable: 즉시 조치 없음
4. 왜 이 지표만으로 충분한가
- 행동은 원인 분류의 첫 단계만 필요하며, 그것을 결정하는 최소 입력이 에러/준비성 상태다.
5. 표시 방식
- Stat 1개(action class) + 색상
- (구현) action_score(0/1/2)를 Grafana Value mapping으로 문자열로 표시
6. 첫 행동(1줄)
- forward/upstream 품질부터 확인하고, 실패 유형을 분류

