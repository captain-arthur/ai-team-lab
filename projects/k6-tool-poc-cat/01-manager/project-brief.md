# k6 Tool POC (CAT 적합성 판단)

**Date:** 2026-03-18  
**Program:** CAT

## 목적
- k6를 **직접 실행**해 보고, CAT에서 “테스트 도구”로 쓸 수 있는지 **YES/NO로 판단**한다.

## 이번 POC에서 검증할 것(실행으로만)
- **Scenario 표현력**: 단일 요청 반복 / 2~3 step 사용자 흐름 / 부하 패턴 변화
- **SLI 측정**: latency(p95), error rate, throughput(요청률)
- **SLO 게이팅**: threshold로 PASS/FAIL이 **자동 결정**되는지(종료 코드 포함)
- **결과 해석 가능성**: 출력 요약에서 “무엇을 봐야 하는지”가 명확한지

## 범위 밖
- 특정 클러스터(ingress 등) 적용/튜닝
- CAT 전체 시스템 설계/통합
- 내부 SLI(클러스터 리소스/컨트롤플레인) 기반 합격 판정 설계

## 성공 기준
- k6 테스트 3개가 **로컬에서 즉시 재현 가능**(스크립트+명령어+실행 결과 포함)
- 각 테스트에서 **SLI(p95/error/throughput)**를 읽고 해석할 수 있음
- threshold 변경으로 **PASS/FAIL이 실제로 달라짐**을 확인
- 아래 3문항에 대해 **YES/NO + 근거(실행 결과 기반)**를 제시
  - k6로 Scenario 표현 가능한가?
  - k6로 SLI 측정 가능한가?
  - k6로 SLO 기반 PASS/FAIL 가능한가?
