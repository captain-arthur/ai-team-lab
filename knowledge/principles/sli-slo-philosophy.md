# SLI/SLO Philosophy (CAT 관점)

Cluster Acceptance Test에서 SLI/SLO를 정의하고 판단할 때 따를 철학이다. Kubernetes scalability 문서의 원칙과 “you promise, we promise” 개념을 받아들이고, **수용 테스트에 맞는 SLO**로 한정한다.

---

## Kubernetes Scalability 철학 참고

SLI/SLO 정의 시 [Kubernetes scalability](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-scalability/slos/slos.md)에서 강조하는 다음 원칙을 참고한다.

| 원칙 | 의미 |
|------|------|
| **Precise and well-defined** | 지표와 목표가 모호하지 않고, 측정 방법·조건이 명확하다. |
| **Consistent** | 환경·시나리오가 같으면 같은 방식으로 측정·판단된다. |
| **User-oriented** | 사용자(워크로드·개발자)가 체감하는 동작·품질을 반영한다. |
| **Testable** | 실제로 측정 가능하고, 수용 테스트에서 재현 가능하다. |

---

## “You promise, we promise”

- **You promise:** 팀/플랫폼이 “이 환경·이 조건에서 이 수준을 보장한다”고 명시한다 (SLO).
- **We promise:** 클러스터/시스템이 그 조건 하에서 그 수준을 만족하는지 검증한다 (CAT에서 측정·단언).
- SLO는 **약속**이므로, 막연한 벤치마크 숫자가 아니라 **정당한 근거**(사용자 기대, 환경 가정, 반복 가능한 측정, 수용 의미)가 있어야 한다.

---

## 수용 테스트용 SLO: 도구 기본값 그대로 쓰지 않기

- **벤치마크·부하 테스트 도구가 주는 넓은 임계값을 그대로 SLO로 쓰지 않는다.**  
  도구 기본값은 “대규모·일반적” 시나리오용일 수 있고, 우리 클러스터 규모·용도와 다를 수 있다.
- **수용 테스트용 SLO는 다음으로 정당화한다:**

  - **User expectations:** 사용자(워크로드 소유자, 플랫폼 이용자)가 기대하는 응답 시간·가용성·처리량 수준.
  - **Environment assumptions:** 클러스터 규모, 노드 수, 리소스, 네트워크 등 “이 조건에서”라는 가정.
  - **Repeatable measurement:** 동일 조건에서 반복 측정 가능하고, 측정 방법이 문서화되어 있음.
  - **Practical acceptance meaning:** “이 SLO를 만족하면 이 클러스터를 수용한다”는 것이 팀·운영에서 말이 되도록 정의.

- SLO를 정할 때 “왜 이 숫자인가?”에 대한 답을 위 네 가지로 줄 수 있어야 한다.

---

## 소규모 클러스터 vs 대규모 벤치마크

- **소규모 클러스터의 수용 기준은 대규모 벤치마크 임계값과 다를 수 있다.**  
  대형 벤치마크(예: 5000노드)에서 쓰는 latency/throughput SLO를 그대로 소규모 클러스터에 적용하면, 과도하게 엄격하거나 의미가 없을 수 있다.
- 수용 테스트 SLO는 **해당 클러스터의 용도·규모·환경 가정**에 맞게 정의한다.  
  “우리 클러스터는 N 노드, 이런 워크로드를 받을 때 이 수준을 만족해야 수용”으로 구체화한다.

---

## 요약

- Kubernetes scalability 원칙: precise, consistent, user-oriented, testable.
- “You promise, we promise”: SLO는 정당한 약속, CAT는 그 검증.
- 수용 테스트 SLO는 도구 기본값이 아니라 user expectations, environment assumptions, repeatable measurement, practical acceptance meaning으로 정당화.
- 소규모 수용 기준 ≠ 대규모 벤치마크 임계값; 규모·용도에 맞게 별도 정의.
