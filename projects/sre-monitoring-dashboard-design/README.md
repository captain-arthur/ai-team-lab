# sre-monitoring-dashboard-design

**Central Kubernetes Operational Dashboard** 설계 프로젝트다. “지금 클러스터가 안전한가?”, “곧 불안전해질 징후는?”에 **확신 있게** 답하고, 이상 시 **어디를 조사할지** 빠르게 찾을 수 있는 대시보드의 레이아웃·신호 계층·조기 리스크 배치를 정의한다.

**Program:** 이 프로젝트는 **SRE Monitoring** 프로그램에 속한다. `programs/sre-monitoring/README.md` 참고.

**Foundation:** 설계는 **Cluster Health Monitoring Model**(`projects/sre-monitoring-cluster-health-model`)을 기반으로 한다. 해당 프로젝트의 final-report, architecture, core-signal-list를 입력으로 사용한다.

---

## 컨텍스트

- **문제:** 대시보드는 많지만 “이게 정상인가?”, “클러스터가 안전한가?”에 대한 **운영 확신(operational confidence)** 이 부족함.
- **목표:** metric 양이 아닌 **운영 확신**과 **조기 리스크 감지**에 초점을 둔 Central Dashboard 설계.
- **주요 산출물:** Central Kubernetes Operational Dashboard 설계 문서(레이아웃, 신호 계층, 조기 리스크 강조, 5–10분 사용 흐름).
