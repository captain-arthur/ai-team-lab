# Final Report: CAT Orchestrator 최소 구조

## 1) CAT Orchestrator의 정의
- CAT Orchestrator는 “도구 실행 + raw 수집 + adapter 정규화 + cat-result.json 저장”을 **파일 규약과 exit code 권위**로 통일하는 최소 오케스트레이션 레이어다.

## 2) Job / Adapter / Runner 관계
- **Job Spec**: 무엇을 실행할지(Scenario, SLO, selected_sli, output dir)를 적는다.
- **Runner**: job을 읽고 tool을 선택해 adapter를 실행한다.
- **Adapter**: tool raw 결과를 읽어 CAT 표준 `cat-result.json`으로 정규화한다.
- **Result Store**: 표준/원시 결과를 지정된 디렉터리에 저장한다.

## 3) 도구 역할 분리(강제)
- k6 → 외부 트래픽 기반 CAT
- CL2 → 클러스터 내부 부하/측정 기반 CAT
- Ginkgo → custom CAT(측정/단언/파일 저장이 Go 코드로 완결)

이때 CAT Orchestrator는 도구 내부 로직을 알 필요가 없다(변환/정규화는 adapter 책임).

## 4) 실제 운영 가능한가(판단)
- 이 구조는 “복잡한 플랫폼”이 아니라, 파일 기반으로 먼저 운영 가능한 최소 형태다.
- 특히 k6/Ginkgo는 이미 `cat-result.json` 생성이 가능하므로 adapter가 얇게 유지된다.
- CL2는 SLI 추출 규칙만 확정되면 동일한 패턴으로 편입 가능하다.

## 5) 다음 단계(확장 방향, 여기서 구현 금지)
- adapter별 `selected_sli` 키/metric 네이밍을 CAT 표준에 맞게 고정
- CL2 adapter의 “필요 SLI만 추출” 규칙을 확정
- optional overall aggregator 정책(PASS_WITH_WARNINGS 등)만 추가

## Final Question
👉 이 구조로 CAT 시스템을 바로 도입할 수 있는가?

**판단: YES(최소 도입 가능).**  
단, adapter들이 표준 `cat-result.json`을 생성하고, selected_sli 키/매핑이 고정되어야 한다.

