# Ginkgo 코어 가이드(실전)

## 1) Ginkgo는 무엇인가
- Go 기본 `testing`을 BDD 스타일(명세 기반)로 확장한 테스트 프레임워크다.
- 테스트를 “기능/상태/조건” 단위로 계층화해, 실패 지점이 드러나게 만든다.

## 2) Go 기본 testing만으로 부족해지는 “조건”
아래 중 하나면 Ginkgo가 즉시 유리하다.
- 테스트가 “시나리오(상태 전환) + 여러 검증”으로 커졌고, 실패 원인을 찾기 어려워졌다.
- lifecycle(준비/정리/전후 조건)이 중요하고, 각 테스트가 같은 패턴으로 반복된다.
- 비동기 검증(상태 변화 대기)이 섞이면서, `testing`의 단순 assert들이 가독성을 잃었다.

즉, “코드가 더 짧아져서”가 아니라 “테스트 구조가 더 명확해져서” 쓴다.

## 3) Ginkgo의 핵심 철학(왜 이런 구조인가)
- **테스트는 문장**이다: `Describe/Context/It`가 상황-행동-기대의 순서로 읽히게 설계된다.
- **상태 계층화**: 같은 검증이라도 “어떤 상태에서”인지가 코드 구조로 드러나야 한다.
- **실패 메시지/증거가 중요**: 단언(assert)이 실패하면, 어느 상황에서 무엇을 기대했는지 먼저 보이도록 한다.

이 철학이 custom CAT에 특히 중요한 이유:
- CAT Job은 “Scenario 실행 → SLI 수집 → SLO 단언”이 한 덩어리이고,
- 그 덩어리를 여러 상태/조건으로 나눠야 유지보수가 쉬워진다.

## 4) Describe / Context / When / It 역할 차이
아래는 “의미”만 기억하면 된다.
- `Describe("기능/시스템 영역")`  
  - 테스트 주제의 최상위 컨테이너(예: `CAT Job`, `Recovery`, `Ingress`).
- `Context("상태/전제 조건")`  
  - Describe 아래에서 “전제 조건이 바뀌는 지점”을 드러낸다(예: `when latency is high`).
- `When("조건부 상황(선택적 가독성)")`  
  - Context와 유사하지만, “특정 조건이 있을 때만”이라는 뉘앙스가 더 강하다.
  - 팀 컨벤션으로 정하면 된다(둘 다 계층화를 위한 도구).
- `It("행동 + 기대")`  
  - 실제 테스트 실행과 SLO/단언이 들어가는 최하위 단계.

실전 규칙:
- lifecycle은 `It` 밖(Context/Describe)에 걸면 재사용되지만, 상태가 섞이기 쉬우니 신중히 둔다.

## 5) 테스트를 읽고 구조화하는 법(패턴)
Ginkgo CAT 스타일 테스트는 보통 아래 순서로 읽힌다.
1. `Describe`: CAT Job의 범위(어떤 Job인지)
2. `Context`: scenario input/전제(SLO/환경 가정 포함)
3. `It`: scenario 실행(요청/상태 전환) + SLI 계산 + SLO 단언 + 결과 파일 저장

### 짧은 예제(구조만)
```go
var _ = Describe("custom CAT job", func() {
  Context("when error rate is low", func() {
    It("should PASS and write cat-result.json", func() {
      // Scenario 실행
      // SLI 계산
      // SLO 단언(PASS/FAIL)
      // cat-result.json 저장
    })
  })
})
```

안티패턴(구조 붕괴):
- `It` 안에 setup도 잔뜩, teardown도 잔뜩 넣어서 “문장이 안 읽히는 테스트”로 만드는 것.
- Context 아래에 global 변수를 마구 공유해서, 상황 전환이 일어나면 이전 상태가 남는 것.
