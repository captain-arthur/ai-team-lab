# Final Report: k6 Tool POC (CAT 적합성 평가)

**Date:** 2026-03-18  
**Program:** CAT

## 1) k6 한줄 정의
k6는 “정의된 클라이언트 시나리오로 트래픽을 발생시키고, SLI를 집계해 threshold로 PASS/FAIL을 자동 결정하는” 실행형 테스트 도구다.

## 2) k6 사용 흐름(요약)
- 스크립트에 **scenario(부하 모델)** + **요청/체크** + **threshold(SLO)** 정의
- `k6 run ...` 실행
- 실행 요약/`--summary-export`로 **SLI(p95/실패율/요청률)** 확인
- threshold 위반 여부로 **PASS/FAIL(종료 코드 포함)** 확정

## 3) CAT에서 k6의 역할(무엇을 맡기고, 어디서 멈추는가)
- **역할(권장)**: “클러스터 외부 관측 기반” CAT 테스트에서 **부하 발생기 + 외부 SLI 계산기 + SLO 게이트(PASS/FAIL)**.
- **멈추는 지점(경계)**: 내부 SLI(제어면/노드/리소스/재시작/큐 길이 등)를 “측정/판정”하는 주체가 아님.

## 4) 핵심 기능을 CAT 관점으로 정리(Scenario / SLI / SLO)
- **Scenario**: executor로 “동시성(VU)” 또는 “도착률(RPS)” 축에서 부하 형태를 정의
- **SLI**: `http_req_duration p(95)`, `http_req_failed rate`, (참고) `http_reqs rate`
- **SLO**: threshold로 조건 위반 시 **FAIL + 종료 코드**(자동화 가능한 판정)

## 5) 실험 결과 요약(실제 실행 기반, 의미 포함)
산출물: `04-engineering/engineering-notes.md`, `04-engineering/results/`

- **Test 1(단일 요청 반복)**: PASS  
  - 무엇을 보여줌: “요청→SLI 집계→threshold 판정” 기본 루프가 성립
- **Test 2(3-step 흐름 + think time)**: PASS  
  - 무엇을 보여줌: 다단계 사용자 흐름을 표현 가능, think time은 latency가 아니라 **달성 처리량**을 바꿈
- **Test 3(부하 모델/threshold/sleep)**:
  - VU vs arrival-rate는 “부하 정의 축”이 다름(throughput 해석이 달라짐)
  - threshold를 바꾸면 실제로 FAIL이 나고 종료 코드가 바뀜(=SLO 게이트로 사용 가능)

## 6) CAT 적합성 평가(필수 질문, YES/NO + 이유)
- **k6로 Scenario 표현 가능한가?** **YES**
  - 이유: 단일 반복/다단계 흐름/부하 패턴(VU 램프, arrival-rate 램프)을 스크립트로 재현.
- **k6로 SLI 측정 가능한가?** **YES**
  - 이유: p95/실패율/요청률이 요약에 출력되고 `--summary-export`로 저장됨.
- **k6로 SLO 기반 PASS/FAIL 가능한가?** **YES**
  - 이유: threshold 위반 시 FAIL 표시 + 비0 종료 코드(실행에서 확인).

## 7) Prometheus는 언제/왜 필요한가(모호함 제거)
- **k6 단독으로 충분한 경우**: 합격 기준이 “외부 SLI”(예: ingress 경유 p95/실패율)로 완결될 때.
- **결합이 필요한 경우**: (a) 외부 SLI가 나빠졌을 때 원인을 내부로 좁혀야 하거나, (b) 합격 기준 자체가 내부 SLI(자원 포화/재시작/제어면 지연 등)를 포함할 때.  
  - 이때 Prometheus는 “PASS/FAIL을 대체”가 아니라, **(1) 원인 분류** 또는 **(2) 내부 SLO 게이트** 역할을 한다.

## 8) CL2(clusterloader2)와의 관계(필수 질문)
- **k6가 CL2를 대체하는가?** **NO**
- **보완하는가?** **YES**
- **핵심 차이(짧게)**:
  - k6: **클러스터 밖**에서 사용자 트래픽을 만들고 **외부 SLI**로 판정
  - CL2: **클러스터 안**에서 리소스/오브젝트/제어면 성격의 부하를 만들고 **내부 SLI**를 다루는 쪽에 강함

## 9) 최종 결론(판단형)
1. **최종 판단**: **조건부 YES**
2. **판단 이유**
   - Scenario/SLI/SLO(PASS/FAIL) 3요소를 k6 단독으로 실행해 확인했고, threshold 위반 시 종료 코드까지 바뀌어 “게이트”로 동작함.
3. **적합한 사용 사례**
   - ingress/HTTP 앞단처럼 “클라이언트 관측 SLI”로 합격 판정하는 CAT 테스트
4. **부적합한 사용 사례**
   - 내부 SLI가 합격 기준의 핵심인 테스트(제어면/노드/리소스/재시작 기반), 또는 브라우저 수준 현실성이 필요한 경우
5. **CAT 아키텍처 내 권장 위치**
   - 외부 트래픽 시나리오 실행기(Scenario) + 외부 SLI 계산기 + SLO 게이트(판정).  
   - 실패 시: Prometheus는 “원인 분류(진단)” 또는 “내부 SLO 게이트(추가 조건)”로 결합.
