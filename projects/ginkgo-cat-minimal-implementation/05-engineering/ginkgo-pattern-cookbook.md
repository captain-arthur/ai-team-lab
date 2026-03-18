# Ginkgo 패턴 쿡북(CAT/실전)

각 패턴은 “문제 상황 → 추천 패턴 → 예제 코드 → 안티패턴” 순서로 적는다.

## 1) HTTP endpoint 테스트 패턴
### 문제 상황
- endpoint 응답/코드/바디를 검증해야 한다.
### 추천 패턴
- 테스트에서 서버를 직접 띄우고(httptest), `It`에서 요청→검증을 짧게 끝낸다.
### 예제 코드
```go
srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("ok"))
}))
defer srv.Close()

resp, err := http.Get(srv.URL)
Expect(err).ToNot(HaveOccurred())
Expect(resp.StatusCode).To(Equal(http.StatusOK))
```
### 안티패턴
- `BeforeAll`에 서버를 띄우고, 상태를 테스트별로 바꾸며 공유하는 것(간헐 실패).

## 2) 상태 변화 대기 패턴(Eventual consistency)
### 문제 상황
- 어떤 조건이 “나중에” 참이 되는지를 기다려야 한다.
### 추천 패턴
- 조건만 함수로 만들고 `Eventually`에 넣는다.
### 예제 코드
```go
Eventually(func() bool {
  return atomic.LoadBool(&ready)
}, "2s", "20ms").Should(BeTrue())
```
### 안티패턴
- 조건 함수 안에서 매번 요청을 보내며 side-effect를 만드는 것(대기 중 중복 트래픽/상태 오염).

## 3) 반복 측정(측정 가능한 SLI 만들기)
### 문제 상황
- 평균/에러율/복구시간 같은 값을 “여러 번” 관측해야 한다.
### 추천 패턴
- 요청 loop를 helper로 분리하고, samples/latencies를 수집한다.
### 예제 코드
```go
samples := make([]time.Duration, 0, N)
errors := 0
for i := 0; i < N; i++ {
  start := time.Now()
  resp, err := http.Get(url)
  lat := time.Since(start)
  samples = append(samples, lat)
  if err != nil || resp.StatusCode != 200 {
    errors++
  }
}
errorRate := float64(errors) / float64(N)
Expect(errorRate).To(BeNumerically("<=", 0.01), "errorRate=%f", errorRate)
```
### 안티패턴
- 측정 로직과 assertion을 한 덩어리로 섞어 읽기 어려워지는 것.

## 4) 리소스 준비/정리 패턴
### 문제 상황
- 서버, 파일, 임시 디렉터리 등을 준비하고 반드시 정리해야 한다.
### 추천 패턴
- `BeforeEach + AfterEach` 또는 “증거 저장”은 `DeferCleanup`을 사용한다.
### 예제 코드
```go
BeforeEach(func() { srv = httptest.NewServer(...) })
AfterEach(func() { srv.Close() })

DeferCleanup(func() {
  _ = os.RemoveAll(tmpDir)
})
```
### 안티패턴
- 전역 변수로 srv를 공유하고, 테스트 실행 순서에 의존하는 것.

## 5) 표 기반(table-driven) 테스트
### 문제 상황
- 입력과 기대값 쌍이 많다.
### 추천 패턴
- `DescribeTable` + `Entry`로 목록을 선언한다.
### 예제 코드
```go
DescribeTable("latency bucket",
  func(ms int, expected string) {
    Expect(bucket(ms)).To(Equal(expected))
  },
  Entry("10ms", 10, "fast"),
  Entry("500ms", 500, "slow"),
)
```
### 안티패턴
- if/else가 길어지고 메시지가 케이스별로 분리되지 않는 것.

## 6) 실패 메시지 잘 남기는 패턴
### 문제 상황
- 실패했을 때 “왜”를 빠르게 알아야 한다.
### 추천 패턴
- `Expect(...).To(..., "context=%s value=%v", ...)`처럼 message를 넣는다.
### 예제 코드
```go
Expect(avgLatMs).To(BeNumerically("<=", sloMax),
  "avgLatMs=%f sloMax=%f", avgLatMs, sloMax)
```
### 안티패턴
- 실패 메시지가 전혀 없는 단순 `Expect(x).To(Equal(y))`만 반복.

## 7) 결과 파일 저장(증거 보존) 패턴
### 문제 상황
- assertion 실패여도 결과 파일은 남겨야 한다.
### 추천 패턴(권장 2단계)
- (1) 계산이 끝난 뒤 즉시 `WriteCatResult`로 먼저 저장
- (2) 마지막으로 `DeferCleanup`으로 “최종 안전장치”를 둔다
### 예제 코드
```go
_ = catutil.WriteCatResult(path, cat) // 먼저 저장
defer func() { _ = catutil.WriteCatResult(path, cat) }() // 안전장치
Expect(cat.FinalPassFail).To(Equal("PASS"))
```
### 안티패턴
- defer cleanup에서만 파일을 쓰고, 파일 존재를 assertion 전에 바로 검증하는 것.

## 8) helper 함수 분리 패턴
### 문제 상황
- scenario(주입) / measurement(계산) / evaluation(단언)이 섞이면 읽기 어렵다.
### 추천 패턴
- `buildScenarioParams`, `measureSLI`, `evaluateSLO`, `writeResult` 순으로 함수로 분리
### 예제 코드(스케치)
```go
params := buildParamsFromEnv()
sli := measureSLI(params)
ok := evaluateSLO(sli, params)
writeResult(params, sli, ok)
Expect(ok).To(BeTrue())
```
### 안티패턴
- 모든 것을 It 내부 200줄로 몰아넣는 것.
