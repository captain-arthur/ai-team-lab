# 골든 시그널 최종(각 컴포넌트 4개, 최소)

## istio-ingressgateway
### 골든 시그널 4개
1. `error_rate_5m` (5xx 비율)
   - 왜 중요한가: 외부 사용자 관측 실패 품질을 가장 직접적으로 반영한다.
   - 어디에 들어가는가: Row 1(지금 안전), Row 2(내일 위험 추세)
2. `p95_latency_ms_5m`
   - 왜 중요한가: 실패가 없어도 지연이 누적되면 “내일 사용자 영향”이 시작된다.
   - 어디에 들어가는가: Row 1(지금 안전), Row 2(내일 위험 추세)
3. `request_rate_5m`
   - 왜 중요한가: 동일한 품질이라도 부하가 커지면 위험이 더 빠르게 현실화된다(포화 가능성).
   - 어디에 들어가는가: Row 2(내일 위험 high 승격 보강)
4. `ready_ok`
   - 왜 중요한가: 게이트웨이 자체가 준비되지 않으면 품질 붕괴가 반복될 가능성이 높다.
   - 어디에 들어가는가: Row 1(지금 안전 fail 확증)

## nodelocaldns
### 골든 시그널 4개
1. `dns_error_rate_5m` (SERVFAIL 비율; node-local-dns)
   - 왜 중요한가: 이름 해석 실패 품질을 직접 반영한다.
   - 어디에 들어가는가: Row 1(지금 안전), Row 2(내일 위험 추세)
2. `dns_latency_p95_ms_5m`
   - 왜 중요한가: DNS 지연은 재시도/타임아웃으로 연쇄 장애를 만든다.
   - 어디에 들어가는가: Row 1(지금 안전), Row 2(내일 위험 추세)
3. `forward_error_rate_5m` (forward SERVFAIL 비율; node-local-dns)
   - 왜 중요한가: 로컬 캐시만의 문제인지 upstream 연결/품질 문제인지 첫 분리를 돕는다.
   - 어디에 들어가는가: Row 1(지금 안전 fail/warning 판정)
4. `ready_ok`
   - 왜 중요한가: 로컬 DNS 준비성이 흔들리면 “내일도 깨질 가능성”이 커진다.
   - 어디에 들어가는가: Row 2(내일 위험 high 승격)

## coredns
### 골든 시그널 4개
1. `dns_error_rate_5m` (SERVFAIL 비율; node-local-dns 제외)
   - 왜 중요한가: 클러스터 DNS 품질 붕괴를 직접 반영한다.
   - 어디에 들어가는가: Row 1(지금 안전), Row 2(내일 위험 추세)
2. `dns_latency_p95_ms_5m`
   - 왜 중요한가: 지연 누적은 재시도 증가 → 큐잉/포화로 전이될 수 있다.
   - 어디에 들어가는가: Row 1(지금 안전), Row 2 보조(구현에서 Panel 3에 포함)
3. `ready_ok`
   - 왜 중요한가: readiness 불안정은 장애 확률을 높인다.
   - 어디에 들어가는가: Row 1(지금 안전 fail 확증)
4. `restarts_10m`
   - 왜 중요한가: 짧은 기간의 재시작은 회복 실패/불안정의 증거다.
   - 어디에 들어가는가: Row 2(내일 위험 high 승격)

