# sre-monitoring-cluster-health-model

Kubernetes 클러스터의 **운영 건강(Cluster Health)** 을 일상적으로 파악하기 위한 **Cluster Health Monitoring Model** 설계 프로젝트다. “지금 안전한가?”, “곧 불안전해질 징후는?”, “문제 시 어디를 먼저 볼 것인가?”를 5–10분 안에 답할 수 있는 최소 신호 집합을 정의한다.

**Program:** 이 프로젝트는 **SRE Monitoring** 프로그램에 속한다. 프로그램 목적과 대표 프로젝트 유형은 `programs/sre-monitoring/README.md`를 참고한다.

---

## 컨텍스트

- **입력:** Task intake(`tasks/sre-monitoring-cluster-health-model.md`), Manager brief.
- **단계:** Manager → Research → Architecture → Engineering → Experiment → Review → Documentation.
- **주요 산출물:** Cluster Health Monitoring Model 문서, core signal list(Prometheus/PromQL 포함), 실험·제한 사항 노트, 리뷰 요약.

후속 프로젝트: `sre-monitoring-dashboard-design`, `sre-monitoring-alert-policy`, `sre-monitoring-operational-runbooks`.
