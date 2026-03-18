# 컴포넌트별 “핵심 지표 3개”만 확정

이 문서는 각 컴포넌트에 대해 **판정 가능한 최소 지표 3개**만 고정한다.
추가 지표는 금지(“이것도 있으면 좋지 않을까” 금지).

## istio-ingressgateway
- error_rate
- p95_latency_ms
- health_ready

왜 이 3개만 충분한가
- ingress는 사용자 영향이 곧 실패/지연/준비성으로 나타나며, 이 3개가 fail/warning/safe의 판정에 직접 기여한다.

## nodelocaldns
- dns_timeout_rate
- dns_latency_p95_ms
- upstream_reachability

왜 이 3개만 충분한가
- 로컬 DNS에서 실제 사용자는 “timeout(실패)”과 “지연(품질)”을 체감하고, upstream reachability는 원인이 로컬인지 upstream인지 첫 분리를 가능하게 한다.

## coredns
- dns_error_rate
- dns_latency_p95_ms
- health_ready

왜 이 3개만 충분한가
- coredns는 클러스터 DNS의 기준점이므로, 에러/지연/준비성만으로 현재 영향과 내일 위험(누적)을 가장 짧게 판정할 수 있다.

