# PromQL 스펙 최종(9개 패널, 최소 state_score)

각 패널은 숫자만 반환한다.
- Q1(지금 안전): safe=0, warning=1, fail=2
- Q2(내일 안전): stable=0, risk=1, high=2
- Q3(무엇을 해야): 예방 확인=0, 원인 분류=1, 즉시 완화 준비=2

핵심 계산식
- Q1: `state_score = warning_cond + fail_cond`
- Q2: `risk_score = warning_cond + high_cond`
- Q3: `action_score = warning_action_cond + immediate_action_cond`

조건은 각자 boolean(0/1)으로 만들고, `fail/high`가 `warning/risk`를 포함하도록 구성해 점수가 정확히 2로 고정된다.

임계값(초안, 환경 튜닝 필요)
- istio
  - error_rate: warning 0.01, fail 0.03
  - p95 latency: warning 500ms, fail 1000ms
  - request_rate high 보강: 200 rps
- DNS(node-local-dns/coredns)
  - error_rate(SERVFAIL fraction): warning 0.01, fail 0.03
  - p95 latency: warning 50ms, fail 150ms
- coredns restarts(10m 증가): warning >2, high >5

---

## Panel 1: istio-ingressgateway - 지금 안전한가?
1. 질문(Row): Row 1
2. 입력 시그널(3개): error_rate_5m, p95_latency_ms_5m, ready_ok
3. PromQL 초안(0/1/2)
```promql
(
  (
    (
      (
        sum(rate(istio_requests_total{destination_workload="istio-ingressgateway",reporter="destination",response_code=~"5.."}[5m]))
        /
        sum(rate(istio_requests_total{destination_workload="istio-ingressgateway",reporter="destination"}[5m]))
      ) > 0.01
    ) bool
    +
    (
      histogram_quantile(0.95, sum(rate(istio_request_duration_milliseconds_bucket{destination_workload="istio-ingressgateway",reporter="destination"}[5m])) by (le))
      > 500
    ) bool
    +
    (sum(kube_pod_container_status_ready{condition="true",pod=~"istio-ingressgateway.*"}) == 0) bool
  ) > 0
) bool
+ (
  (
    (
      (
        sum(rate(istio_requests_total{destination_workload="istio-ingressgateway",reporter="destination",response_code=~"5.."}[5m]))
        /
        sum(rate(istio_requests_total{destination_workload="istio-ingressgateway",reporter="destination"}[5m]))
      ) > 0.03
    ) bool
    +
    (
      histogram_quantile(0.95, sum(rate(istio_request_duration_milliseconds_bucket{destination_workload="istio-ingressgateway",reporter="destination"}[5m])) by (le))
      > 1000
    ) bool
    +
    (sum(kube_pod_container_status_ready{condition="true",pod=~"istio-ingressgateway.*"}) == 0) bool
  ) > 0
) bool
```
4. 표시 방식: Stat(0/1/2 mapping + green/yellow/red)
5. 첫 행동 1줄
   - upstream vs 게이트웨이 준비성부터 분리

---

## Panel 2: nodelocaldns - 지금 안전한가?
1. 질문(Row): Row 1
2. 입력 시그널(3개): dns_error_rate_5m, dns_latency_p95_ms_5m, forward_error_rate_5m
3. PromQL 초안(0/1/2)
```promql
(
  (
    (
      (
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app="node-local-dns",rcode="SERVFAIL"}[5m]))
        /
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app="node-local-dns"}[5m]))
      ) > 0.01
    ) bool
    +
    (
      histogram_quantile(0.95, sum(rate(coredns_dns_request_duration_seconds_bucket{k8s_app="node-local-dns"}[5m])) without(instance,pod))
      * 1000 > 50
    ) bool
    +
    (
      sum(rate(coredns_forward_response_rcode_count_total{k8s_app="node-local-dns",rcode="SERVFAIL"}[5m]))
      /
      sum(rate(coredns_forward_response_rcode_count_total{k8s_app="node-local-dns"}[5m]))
      > 0.01
    ) bool
  ) > 0
) bool
+ (
  (
    (
      (
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app="node-local-dns",rcode="SERVFAIL"}[5m]))
        /
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app="node-local-dns"}[5m]))
      ) > 0.03
    ) bool
    +
    (
      histogram_quantile(0.95, sum(rate(coredns_dns_request_duration_seconds_bucket{k8s_app="node-local-dns"}[5m])) without(instance,pod))
      * 1000 > 150
    ) bool
    +
    (
      sum(rate(coredns_forward_response_rcode_count_total{k8s_app="node-local-dns",rcode="SERVFAIL"}[5m]))
      /
      sum(rate(coredns_forward_response_rcode_count_total{k8s_app="node-local-dns"}[5m]))
      > 0.03
    ) bool
  ) > 0
) bool
```
4. 첫 행동 1줄
   - forward/upstream 오류 원인 분리

---

## Panel 3: coredns - 지금 안전한가?
1. 질문(Row): Row 1
2. 입력 시그널(3개): dns_error_rate_5m, dns_latency_p95_ms_5m, ready_ok
3. PromQL 초안(0/1/2)
```promql
(
  (
    (
      (
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app!="node-local-dns",rcode="SERVFAIL"}[5m]))
        /
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app!="node-local-dns"}[5m]))
      ) > 0.01
    ) bool
    +
    (
      histogram_quantile(0.95, sum(rate(coredns_dns_request_duration_seconds_bucket{k8s_app!="node-local-dns"}[5m])) without(instance,pod))
      * 1000 > 50
    ) bool
    +
    (sum(kube_pod_container_status_ready{condition="true",pod=~"coredns.*|kube-dns.*"}) == 0) bool
  ) > 0
) bool
+ (
  (
    (
      (
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app!="node-local-dns",rcode="SERVFAIL"}[5m]))
        /
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app!="node-local-dns"}[5m]))
      ) > 0.03
    ) bool
    +
    (
      histogram_quantile(0.95, sum(rate(coredns_dns_request_duration_seconds_bucket{k8s_app!="node-local-dns"}[5m])) without(instance,pod))
      * 1000 > 150
    ) bool
    +
    (sum(kube_pod_container_status_ready{condition="true",pod=~"coredns.*|kube-dns.*"}) == 0) bool
  ) > 0
) bool
```
4. 첫 행동 1줄
   - forward/upstream 품질부터 확인

---

## Panel 4: istio-ingressgateway - 내일도 안전한가?
1. 질문(Row): Row 2
2. 입력 시그널(3개): error_rate_1h, p95_latency_ms_5m, request_rate_5m(보강)
3. PromQL 초안(0/1/2)
```promql
(
  (
    (
      (
        sum(rate(istio_requests_total{destination_workload="istio-ingressgateway",reporter="destination",response_code=~"5.."}[1h]))
        /
        sum(rate(istio_requests_total{destination_workload="istio-ingressgateway",reporter="destination"}[1h]))
      ) > 0.01
    ) bool
    +
    (
      histogram_quantile(0.95, sum(rate(istio_request_duration_milliseconds_bucket{destination_workload="istio-ingressgateway",reporter="destination"}[5m])) by (le))
      > 500
    ) bool
  ) > 0
) bool
+ (
  (
    (
      (
        sum(rate(istio_requests_total{destination_workload="istio-ingressgateway",reporter="destination",response_code=~"5.."}[1h]))
        /
        sum(rate(istio_requests_total{destination_workload="istio-ingressgateway",reporter="destination"}[1h]))
      ) > 0.03
    ) bool
    +
    (
      histogram_quantile(0.95, sum(rate(istio_request_duration_milliseconds_bucket{destination_workload="istio-ingressgateway",reporter="destination"}[5m])) by (le))
      > 1000
    ) bool
    +
    (
      (sum(rate(istio_requests_total{destination_workload="istio-ingressgateway",reporter="destination"}[5m])) > 200) bool
      *
      (histogram_quantile(0.95, sum(rate(istio_request_duration_milliseconds_bucket{destination_workload="istio-ingressgateway",reporter="destination"}[5m])) by (le)) > 800) bool
    ) bool
  ) > 0
) bool
```
4. 첫 행동 1줄
   - 사전 완화 준비

---

## Panel 5: nodelocaldns - 내일도 안전한가?
1. 질문(Row): Row 2
2. 입력 시그널(3개): dns_error_rate_1h, dns_latency_p95_ms_5m, ready_ok
3. PromQL 초안(0/1/2)
```promql
(
  (
    (
      (
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app="node-local-dns",rcode="SERVFAIL"}[1h]))
        /
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app="node-local-dns"}[1h]))
      ) > 0.01
    ) bool
    +
    (
      histogram_quantile(0.95, sum(rate(coredns_dns_request_duration_seconds_bucket{k8s_app="node-local-dns"}[5m])) without(instance,pod))
      * 1000 > 50
    ) bool
    +
    (sum(kube_pod_container_status_ready{condition="true",pod=~"node-local-dns.*"}) == 0) bool
  ) > 0
) bool
+ (
  (
    (
      (
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app="node-local-dns",rcode="SERVFAIL"}[1h]))
        /
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app="node-local-dns"}[1h]))
      ) > 0.03
    ) bool
    +
    (
      histogram_quantile(0.95, sum(rate(coredns_dns_request_duration_seconds_bucket{k8s_app="node-local-dns"}[5m])) without(instance,pod))
      * 1000 > 150
    ) bool
    +
    (sum(kube_pod_container_status_ready{condition="true",pod=~"node-local-dns.*"}) == 0) bool
  ) > 0
) bool
```
4. 첫 행동 1줄
   - upstream/완화 준비

---

## Panel 6: coredns - 내일도 안전한가?
1. 질문(Row): Row 2
2. 입력 시그널(3개): dns_error_rate_1h, restarts_10m, ready_ok
3. PromQL 초안(0/1/2)
```promql
(
  (
    (
      (
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app!="node-local-dns",rcode="SERVFAIL"}[1h]))
        /
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app!="node-local-dns"}[1h]))
      ) > 0.01
    ) bool
    +
    (sum(increase(kube_pod_container_status_restarts_total{pod=~"coredns.*|kube-dns.*"}[10m])) > 2) bool
    +
    (sum(kube_pod_container_status_ready{condition="true",pod=~"coredns.*|kube-dns.*"}) == 0) bool
  ) > 0
) bool
+ (
  (
    (
      (
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app!="node-local-dns",rcode="SERVFAIL"}[1h]))
        /
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app!="node-local-dns"}[1h]))
      ) > 0.03
    ) bool
    +
    (sum(increase(kube_pod_container_status_restarts_total{pod=~"coredns.*|kube-dns.*"}[10m])) > 5) bool
    +
    (sum(kube_pod_container_status_ready{condition="true",pod=~"coredns.*|kube-dns.*"}) == 0) bool
  ) > 0
) bool
```
4. 첫 행동 1줄
   - coredns 안정화 준비

---

## Panel 7: istio-ingressgateway - 무엇을 해야 하는가?
1. 질문(Row): Row 3
2. 입력 시그널(번역): Panel1 fail OR Panel4 high 우선
3. PromQL 초안(0/1/2)
```promql
(
  (
    (
      (
        sum(rate(istio_requests_total{destination_workload="istio-ingressgateway",reporter="destination",response_code=~"5.."}[5m]))
        /
        sum(rate(istio_requests_total{destination_workload="istio-ingressgateway",reporter="destination"}[5m]))
      ) > 0.01
    ) bool
    +
    (
      histogram_quantile(0.95, sum(rate(istio_request_duration_milliseconds_bucket{destination_workload="istio-ingressgateway",reporter="destination"}[5m])) by (le))
      > 500
    ) bool
    +
    (sum(kube_pod_container_status_ready{condition="true",pod=~"istio-ingressgateway.*"}) == 0) bool
    +
    (
      (
        sum(rate(istio_requests_total{destination_workload="istio-ingressgateway",reporter="destination",response_code=~"5.."}[1h]))
        /
        sum(rate(istio_requests_total{destination_workload="istio-ingressgateway",reporter="destination"}[1h]))
      ) > 0.01
    ) bool
  ) > 0
) bool
+ (
  (
    (
      (
        sum(rate(istio_requests_total{destination_workload="istio-ingressgateway",reporter="destination",response_code=~"5.."}[5m]))
        /
        sum(rate(istio_requests_total{destination_workload="istio-ingressgateway",reporter="destination"}[5m]))
      ) > 0.03
    ) bool
    +
    (
      histogram_quantile(0.95, sum(rate(istio_request_duration_milliseconds_bucket{destination_workload="istio-ingressgateway",reporter="destination"}[5m])) by (le))
      > 1000
    ) bool
    +
    (sum(kube_pod_container_status_ready{condition="true",pod=~"istio-ingressgateway.*"}) == 0) bool
    +
    (
      (
        sum(rate(istio_requests_total{destination_workload="istio-ingressgateway",reporter="destination",response_code=~"5.."}[1h]))
        /
        sum(rate(istio_requests_total{destination_workload="istio-ingressgateway",reporter="destination"}[1h]))
      ) > 0.03
    ) bool
    +
    (
      (sum(rate(istio_requests_total{destination_workload="istio-ingressgateway",reporter="destination"}[5m])) > 200) bool
      *
      (histogram_quantile(0.95, sum(rate(istio_request_duration_milliseconds_bucket{destination_workload="istio-ingressgateway",reporter="destination"}[5m])) by (le)) > 800) bool
    ) bool
  ) > 0
) bool
```
4. 첫 행동 1줄
   - fail/high면 즉시 upstream/라우팅 분리 후 대응

---

## Panel 8: nodelocaldns - 무엇을 해야 하는가?
1. 질문(Row): Row 3
2. 입력 시그널(번역): Panel2 fail OR Panel5 high 우선
3. PromQL 초안(0/1/2)
```promql
(
  (
    (
      (
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app="node-local-dns",rcode="SERVFAIL"}[5m]))
        /
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app="node-local-dns"}[5m]))
      ) > 0.01
    ) bool
    +
    (
      histogram_quantile(0.95, sum(rate(coredns_dns_request_duration_seconds_bucket{k8s_app="node-local-dns"}[5m])) without(instance,pod))
      * 1000 > 50
    ) bool
    +
    (
      sum(rate(coredns_forward_response_rcode_count_total{k8s_app="node-local-dns",rcode="SERVFAIL"}[5m]))
      /
      sum(rate(coredns_forward_response_rcode_count_total{k8s_app="node-local-dns"}[5m]))
      > 0.01
    ) bool
    +
    (
      (
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app="node-local-dns",rcode="SERVFAIL"}[1h]))
        /
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app="node-local-dns"}[1h]))
      ) > 0.01
    ) bool
    +
    (sum(kube_pod_container_status_ready{condition="true",pod=~"node-local-dns.*"}) == 0) bool
  ) > 0
) bool
+ (
  (
    (
      (
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app="node-local-dns",rcode="SERVFAIL"}[5m]))
        /
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app="node-local-dns"}[5m]))
      ) > 0.03
    ) bool
    +
    (
      histogram_quantile(0.95, sum(rate(coredns_dns_request_duration_seconds_bucket{k8s_app="node-local-dns"}[5m])) without(instance,pod))
      * 1000 > 150
    ) bool
    +
    (
      sum(rate(coredns_forward_response_rcode_count_total{k8s_app="node-local-dns",rcode="SERVFAIL"}[5m]))
      /
      sum(rate(coredns_forward_response_rcode_count_total{k8s_app="node-local-dns"}[5m]))
      > 0.03
    ) bool
    +
    (
      (
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app="node-local-dns",rcode="SERVFAIL"}[1h]))
        /
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app="node-local-dns"}[1h]))
      ) > 0.03
    ) bool
    +
    (sum(kube_pod_container_status_ready{condition="true",pod=~"node-local-dns.*"}) == 0) bool
  ) > 0
) bool
```
4. 첫 행동 1줄
   - forward 오류 원인 분리 후 즉시 조치

---

## Panel 9: coredns - 무엇을 해야 하는가?
1. 질문(Row): Row 3
2. 입력 시그널(번역): Panel3 fail OR Panel6 high 우선
3. PromQL 초안(0/1/2)
```promql
(
  (
    (
      (
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app!="node-local-dns",rcode="SERVFAIL"}[5m]))
        /
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app!="node-local-dns"}[5m]))
      ) > 0.01
    ) bool
    +
    (
      histogram_quantile(0.95, sum(rate(coredns_dns_request_duration_seconds_bucket{k8s_app!="node-local-dns"}[5m])) without(instance,pod))
      * 1000 > 50
    ) bool
    +
    (sum(increase(kube_pod_container_status_restarts_total{pod=~"coredns.*|kube-dns.*"}[10m])) > 2) bool
    +
    (
      (
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app!="node-local-dns",rcode="SERVFAIL"}[1h]))
        /
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app!="node-local-dns"}[1h]))
      ) > 0.01
    ) bool
    +
    (sum(kube_pod_container_status_ready{condition="true",pod=~"coredns.*|kube-dns.*"}) == 0) bool
  ) > 0
) bool
+ (
  (
    (
      (
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app!="node-local-dns",rcode="SERVFAIL"}[5m]))
        /
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app!="node-local-dns"}[5m]))
      ) > 0.03
    ) bool
    +
    (
      histogram_quantile(0.95, sum(rate(coredns_dns_request_duration_seconds_bucket{k8s_app!="node-local-dns"}[5m])) without(instance,pod))
      * 1000 > 150
    ) bool
    +
    (sum(increase(kube_pod_container_status_restarts_total{pod=~"coredns.*|kube-dns.*"}[10m])) > 5) bool
    +
    (
      (
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app!="node-local-dns",rcode="SERVFAIL"}[1h]))
        /
        sum(rate(coredns_dns_response_rcode_count_total{k8s_app!="node-local-dns"}[1h]))
      ) > 0.03
    ) bool
    +
    (sum(kube_pod_container_status_ready{condition="true",pod=~"coredns.*|kube-dns.*"}) == 0) bool
  ) > 0
) bool
```
4. 첫 행동 1줄
   - forward/upstream 분리 후 coredns 안정화 우선

