# Program: SRE Monitoring

**목적:** SRE Monitoring은 시스템 및 클러스터에 대한 모니터링, 알림, 관찰 가능성(observability)을 정의하고 운영하기 위한 프로그램이다. 무엇을 측정할지, metrics를 어떻게 저장·조회할지, 이를 실행 가능한 신호(dashboards, alerts, SLO)로 바꾸는 방법에 초점을 둔다.

---

## 해결하는 문제

- **가시성 공백:** 명확한 metrics 없이 시스템이 동작하여 성능 저하 감지나 용량 계획이 어려움.
- **알림 피로 또는 침묵:** 알림이 너무 시끄럽거나, 중요한 장애를 놓침.
- **신뢰성 목표 불명확:** 서비스·클러스터에 대한 공통 SLO나 error budget이 없음.
- **도구·파이프라인 산재:** 일관된 observability 전략 없이 여러 임시 dashboard·스크립트가 공존.

SRE Monitoring은 metrics, retention, dashboard, alerting rule, (해당 시) SLO/error-budget 관행을 정의하여 팀 운영 방식에 맞추는 방식으로 위 문제를 다룬다.

---

## 이 프로그램에 속하는 프로젝트 유형

SRE Monitoring 아래에 포함되는 대표적인 프로젝트는 다음과 같다.

- **Metrics pipeline 설계** — metrics 수집(Prometheus, exporters 등) 선택·설정, retention, federation 구성.
- **Dashboard 및 시각화** — 클러스터, 노드, 워크로드, API용 핵심 dashboard(Grafana 등) 정의.
- **Alerting rule 및 runbook** — 알림 조건, 심각도, runbook을 정의하여 알림이 실행 가능하도록 함.
- **SLO 및 error budget** — 서비스·클러스터용 SLI/SLO를 정의하고 alerting 또는 error-budget 추적과 연결.
- **Observability 통합** — 로그, trace 등 다른 신호를 metrics 스택과 통합하고 동작 방식을 문서화.

주요 산출물이 “더 나은 모니터링, 알림, observability”인 프로젝트는 SRE Monitoring에 속한다.

---

## 프로그램과 프로젝트의 관계

- **Program:** 장기 도메인(SRE Monitoring). 관련 프로젝트를 묶고, 공유 패턴(예: dashboard 레이아웃, 네이밍 규칙, alert taxonomy)의 근거지를 제공한다.
- **Projects:** AI 팀 워크플로(manager → research → architecture → engineering → experiment → review → documentation)로 실행되는 구체적·유한한 작업. 각 프로젝트는 config, dashboard, 문서 등 산출물을 만들며, 이는 `projects/<project-name>/`에 두고 README에서 **SRE Monitoring** 프로그램에 속함을 **선언**한다.

프로젝트는 `projects/` 아래에 둔다. 프로그램이 프로젝트 폴더를 보유하지 않는다. 이 README는 프로그램의 목적과 포함하는 프로젝트 유형을 설명한다.
