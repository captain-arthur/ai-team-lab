# Grafana 빌드 노트(최소 수정 포인트)

## 1) datasource만 바꾸면 바로 쓸 수 있는가?
대부분의 경우 가능합니다.
- 대시보드 JSON 안에서 datasource는 `${datasource}`로 placeholder 처리되어 있습니다.
- import 후 Grafana에서 `Datasource`만 실제 Prometheus 데이터소스로 매핑하면 패널이 보통 바로 렌더링됩니다.

## 2) import 후 반드시 확인할 패널 순서(1~9)
1. Panel 1: `istio-ingressgateway - 지금 안전한가?`
2. Panel 2: `nodelocaldns - 지금 안전한가?`
3. Panel 3: `coredns - 지금 안전한가?`
4. Panel 4: `istio-ingressgateway - 내일도 안전한가?`
5. Panel 5: `nodelocaldns - 내일도 안전한가?`
6. Panel 6: `coredns - 내일도 안전한가?`
7. Panel 7: `istio-ingressgateway - 무엇을 해야 하는가?`
8. Panel 8: `nodelocaldns - 무엇을 해야 하는가?`
9. Panel 9: `coredns - 무엇을 해야 하는가?`

## 3) 환경별 수정 포인트(필요 시 JSON의 expr 라벨만 교체)
아래 값이 실제 metric 라벨과 다르면 해당 패널의 `expr`에서 라벨 값을 조정해야 합니다.

### 3.1 istio-ingressgateway(Panel 1/4/7)
- Istio metric에서 `destination_workload="istio-ingressgateway"` 라벨이 실제로 다르면 교체 필요
- readiness 판단에 사용하는 kube metric:
  - `pod=~"istio-ingressgateway.*"` 패턴이 실제 pod 네이밍과 다르면 수정
- latency histogram metric:
  - `istio_request_duration_milliseconds_bucket` 가 실제 이름과 다르면 해당 metric으로 교체
- response_code label:
  - `response_code=~"5.."` 사용(5xx 계열 기준). timeout이 다른 방식으로 노출되면(예: 0/504 등) 기준 조정 필요

### 3.2 nodelocaldns(Panel 2/5/8)
- CoreDNS plugin metric에서 `k8s_app="node-local-dns"` 라벨이 실제로 다르면 교체 필요
- histogram:
  - `coredns_dns_request_duration_seconds_bucket` 사용
- forward/upstream 분리:
  - `coredns_forward_response_rcode_count_total{...,rcode="SERVFAIL"}` 사용
  - forward 관련 metric 이름/라벨이 다르면 해당 expr만 교체
- readiness/restart:
  - readiness: `kube_pod_container_status_ready{pod=~"node-local-dns.*",condition="true"}`
  - 이 metric이 없거나 라벨이 다르면 `kube-state-metrics` 스키마에 맞춰 expr 수정

### 3.3 coredns(Panel 3/6/9)
- “node-local-dns 제외”를 위해 `k8s_app!="node-local-dns"` 로 분리
  - 환경에서 CoreDNS가 다른 라벨(`k8s_app="kube-dns"` 등)로만 잡히면 여전히 동작하지만, 다른 DNS 컴포넌트가 섞일 경우 더 구체화 필요
- restarts:
  - `sum(increase(kube_pod_container_status_restarts_total{pod=~"coredns.*|kube-dns.*"}[10m]))`
  - pod 패턴이 다르면 교체
- readiness:
  - `kube_pod_container_status_ready{pod=~"coredns.*|kube-dns.*",condition="true"}`

## 4) 어떤 metric 이름은 환경 따라 조정이 필요한가(우선순위)
1. `istio_request_duration_milliseconds_bucket` 계열(istio 버전/telemetry 설정에 따라 이름/라벨 변경 가능)
2. CoreDNS plugin metric:
   - `coredns_dns_response_rcode_count_total`
   - `coredns_dns_request_duration_seconds_bucket`
   - `coredns_forward_response_rcode_count_total`
3. kube-state-metrics readiness/restarts:
   - `kube_pod_container_status_ready`
   - `kube_pod_container_status_restarts_total`

## 5) import 후 바로 확인할 quick check
- Panel 2/3/5/6/8/9에서 DNS 관련 2~3개 패널이 정상 렌더링되는지(데이터 없음이면 해당 plugin metrics 미구성 가능성)
- Panel 1/4/7에서 istio latency와 error_rate가 함께 나오는지(데이터 없음이면 Istio telemetry/metric export 미구성 가능성)

