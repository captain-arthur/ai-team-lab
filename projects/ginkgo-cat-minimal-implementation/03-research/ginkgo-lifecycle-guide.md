# Ginkgo 라이프사이클 가이드(실전)

아래 항목은 “무엇을 언제 쓰고, 잘못 쓰면 어떤 문제가 생기는지” 중심이다.

## 1) BeforeEach
### 무엇인지
- 각 `It` 실행 직전에 실행되는 준비 로직.
### 언제 쓰는지
- 각 테스트 케이스마다 “같은 패턴의 준비”가 필요할 때(예: test server 생성, 임시 폴더 준비).
### 잘못 쓰면 생기는 문제
- 준비가 무거운데 매번 반복하면 테스트 시간이 증가한다.
- 공유 상태를 BeforeEach에 숨기면 다음 It로 누출될 수 있다(전역/패키지 변수 주의).
### 예제
```go
var srv *httptest.Server

BeforeEach(func() {
  srv = httptest.NewServer(...)
})
```

## 2) JustBeforeEach
### 무엇인지
- `BeforeEach`가 끝난 뒤, `It` 바로 직전에 실행된다.
- “실행 직전 상태”를 만들고 싶을 때 두 번째 준비 단계로 쓴다.
### 언제 쓰는지
- Context/BeforeEach에서 준비해 둔 값을, It 직전에 최종 조립해야 하는 경우.
### 잘못 쓰면 생기는 문제
- 코드가 복잡해지면 “어떤 단계에서 값이 확정되는지”가 흐려진다.
### 예제
```go
BeforeEach(func() { input = buildBaseInput() })
JustBeforeEach(func() { input.Finalize() })
It("runs cat job", func() { runScenario(input) })
```

## 3) AfterEach
### 무엇인지
- 각 `It`이 끝난 뒤 실행되는 정리 로직.
### 언제 쓰는지
- 서버/파일/임시 리소스 종료가 필요할 때.
### 잘못 쓰면 생기는 문제
- AfterEach가 없으면 누적 리소스가 다음 테스트에 영향을 줄 수 있다(포트/파일 핸들 누수).
### 예제
```go
AfterEach(func() {
  if srv != nil { srv.Close() }
})
```

## 4) BeforeAll / AfterAll
### 무엇인지
- `Describe`(스코프) 내부 테스트 전체에 대해 1회만 실행되는 준비/정리.
### 언제 쓰는지
- 매우 무거운 준비(예: 모델 로딩, 대용량 공용 fixture)가 필요하지만, 테스트들이 독립적일 때만.
### 잘못 쓰면 생기는 문제
- 테스트 케이스 간 상태 누출 위험이 커진다.
- 병렬/순서 실행과 섞이면 재현성 문제가 생긴다.
### 예제
```go
BeforeAll(func() { shared = loadBigFixture() })
AfterAll(func() { shared.Close() })
```

## 5) DeferCleanup
### 무엇인지
- “스코프 종료 시점”에 무조건(cleanup) 실행되는 정리.
- 결과 파일 저장처럼 “Expect 실패해도 증거를 남기는” 패턴에 특히 좋다.
### 언제 쓰는지
- cat-result.json 저장, 로그 파일 덤프, 임시 디렉터리 정리 등.
### 잘못 쓰면 생기는 문제
- cleanup에서 다시 Expect을 걸면 실패 원인이 뒤섞인다.
- cleanup은 실패해도 테스트 실패 원인을 바꾸지 않도록, 가능한 한 안전하게 작성한다.
### 예제(개념)
```go
DeferCleanup(func() {
  _ = WriteCatResult(path, cat)
})
```

## 6) Ordered / Serial / Parallel(실행 제어)
### Ordered
- 특정 Describe/It 묶음을 정의된 순서대로 실행한다.
### Serial
- 병렬로 실행되지 않게 “직렬 실행”으로 제한한다.
### Parallel
- 테스트를 병렬 실행할 수 있게 한다(단, 전역 상태/포트 충돌을 피해야 한다).
### 언제 쓰는지
- 상태가 연쇄적으로 의존하거나, 특정 리소스를 공유해야 할 때.
### 잘못 쓰면 생기는 문제
- 전역/공유 상태를 섞으면 데이터 레이스 또는 간헐 실패가 생긴다.
### 예제(스케치)
```go
var _ = OrderedDescribe("ordered", func() {
  It("step1", func(){ ... })
  It("step2", func(){ ... })
})
```

실전 결론:
- custom CAT Job에서는 보통 `BeforeEach + DeferCleanup + Expect 단언`이 핵심이고,
- 실행 제어(Ordered/Parallel)는 “리소스 공유/의존”이 있을 때만 도입한다.
