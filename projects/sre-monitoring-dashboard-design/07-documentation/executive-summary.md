# Central Kubernetes Operational Dashboard — Executive Summary

**Project:** sre-monitoring-dashboard-design  
**Audience:** Leadership / 의사결정자  
**Version:** v1 (1-page)

---

## 문제 (The Problem)

Kubernetes 클러스터를 Prometheus·Grafana로 모니터링해도, 운영자는 **“지금 클러스터가 안전한가?”, “이 수치는 정상인가?”, “누가 확인해 줄 수 있나?”** 를 계속 묻는다. metric과 패널이 많을수록 **한눈에 결론을 내리기 어렵고**, 실제 판단은 **경험·스크립트·다른 채널**에 의존하게 된다. 즉 **가시성(visibility)** 은 있지만 **운영 확신(operational confidence)** 이 없다.

---

## 핵심 아이디어 (Core Idea)

**“많은 metric 나열”이 아니라 “안전한가?” / “곧 위험한가?”에 직접 답하는 최소한의 시그널**만 한 화면에 둔다.  
운영자가 **5–10분** 안에 다음을 확신 있게 말할 수 있어야 한다.

1. **지금 클러스터가 안전하다** (또는 조사가 필요하다).  
2. **클러스터가 곧 불안전해질 조기 징후가 보인다** (또는 당분간 괜찮다).  
3. **이상일 때 어디를 먼저 조사할 것인가.**

이를 위해 **운영 안전 조건(Operational Safety Conditions)** 과 **고장 전파 경로(Failure Propagation Paths)** 를 정의하고, 각 패널이 **어떤 조건·전파 단계**를 보는지 명시했다. 즉 **사용성**뿐 아니라 **운영 논리**로 설계를 정당화한다.

---

## 대시보드 구조 (Dashboard Structure)

| 블록 | 이름 | 역할 | 패널 수 (메인 뷰) |
|------|------|------|-------------------|
| **Block 1** | Operational Confidence | “지금 안전한가?” — 전부 정상이면 안전, 하나라도 비정상이면 조사 필요 | 4~6개 |
| **Block 2** | Early Risk | “곧 불안전해질 징후인가?” — 노드 압박, OOM 위험, 디스크, Pending 추세 등 | 4~6개 |
| **Block 3** | Investigation / Top Offenders | “어디를 조사할 것인가?” — CPU/메모리/재시작/Pending TOP10 (드릴다운 전용) | 4~6개 (메인 뷰에 미포함) |

**메인 뷰**는 Block 1 + Block 2만 포함하며, **총 8~12개 패널**로 제한한다. Block 3은 **이상 시** 탭·접기로만 진입한다.

---

## 운영 확신이 개선되는 이유 (Why This Improves Operational Confidence)

- **명시적 안전 조건:** “안전하다”의 정의가 **다섯 가지 운영 안전 조건(C1~C5)** 로 정해져 있어, 팀이 **동일한 기준**으로 판단할 수 있다.  
- **선행 vs 후행 구분:** “이미 손상된 상태”를 보는 신호(NotReady, Pending, 재시작 폭증, endpoint 비어 있음)와 **“곧 나빠질 것”을 미리 알려 주는** 신호(노드 CPU·메모리·디스크, Pending 추세)를 블록으로 나누어, **사전 대응**이 가능해진다.  
- **최소·충분:** 메인 뷰에 **필수 8개**만 있어도 “안전한가?” / “조기 징후인가?”에 **결정적으로** 답할 수 있으므로, 노이즈와 과부하를 줄인다.  
- **이상 시 조사 경로:** 비정상 신호 → 카테고리 → Block 3 해당 TOP10 → 1~2개 후보 좁히기 → kubectl·로그·이벤트로 원인 파악. **5–10분 내** “어디를 볼 것인가?”에 답할 수 있다.

---

## 다음 단계 (Next Steps)

- **구현:** implementation-ready-panel-spec에 따라 Grafana + Prometheus로 v1 대시보드 구축.  
- **Alert 정책:** Block 1·Block 2 신호 중 알림으로 쓸 항목·임계치·심각도 정의(sre-monitoring-alert-policy).  
- **Runbook:** Block 1 비정상 시 조사 순서·체크리스트·명령어 정리(sre-monitoring-operational-runbooks).  
- **Managed K8s:** control plane metric이 없으면 Block 1을 4개 패널로 축소하고, 제어면 건강은 managed 서비스 콘솔·알림으로 보완.

---

*Central Kubernetes Operational Dashboard — Executive Summary v1. 상세 설계: central-kubernetes-operational-dashboard-design.md; 운영 이론: 03-architecture/operational-confidence-theory.md.*
