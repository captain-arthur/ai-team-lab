# Practical Check(운영 최소 도입 관점)

## 1) 정말 최소 4개 골든 시그널로 충분한가?
충분하도록 설계했다.
- istio-ingressgateway: 에러율/지연(품질) + 트래픽(포화 보강) + readiness(즉시 붕괴 확증)로 구성되어 Row 1~3의 모든 상태를 만든다.
- nodelocaldns: SERVFAIL 비율(이름 해석 실패) + 지연 + forward SERVFAIL(업스트림 분리) + readiness로 Q1/Q2 및 Q3 번역을 만든다.
- coredns: SERVFAIL 비율 + 지연 + readiness + restarts(불안정/회복 실패)로 Q1/Q2 및 Q3 번역을 만든다.

## 2) PromQL이 너무 복잡하지 않은가?
완전히 “짧지는” 않지만, 복잡도는 제어했다.
- 모든 패널이 boolean 조건의 합으로 0/1/2 점수만 반환한다.
- Row 3(무엇을 해야)는 Row 1/Row 2 실패/위험 조건을 그대로 번역하는 형태다.

## 3) 대시보드가 읽기 쉬운가?
읽기 쉬운 형태로 고정했다.
- 패널은 전부 Stat이며, 표시 텍스트는 1개(안전/주의/실패, 안정/위험/높음, 예방 확인/원인 분류/즉시 완화 준비)만 보여준다.
- 3개 컴포넌트 × 3개 질문(행) 구조로 운영자가 5초 내에 상태 분류를 끝낼 수 있다.

## 4) Action Row가 과하지 않은가?
과하지 않다.
- Row 3은 새로운 분석/계산을 추가하지 않고, “fail/high → 즉시”, “warning/risk → 원인 분류”를 직접 매핑한다.

## 5) 바로 운영에 도입 가능한가?
바로 시작 가능하되, 아래 2가지 확인이 필요하다.
1. DNS plugin metrics(`coredns_dns_*`, `coredns_forward_*`)이 실제로 Prometheus에 노출되는지
2. pod 라벨/패턴(`node-local-dns.*`, `coredns.*|kube-dns.*`)과 istio workload 라벨(`istio-ingressgateway`)이 실제 환경과 맞는지

