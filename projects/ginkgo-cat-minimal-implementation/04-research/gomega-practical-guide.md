# Gomega 실전 가이드(패턴 중심)

## 1) 기본 matcher(자주 쓰는 것)
- `Equal(x)`  
  - 정확히 같음을 기대한다(숫자/문자열/구조체).
- `BeNil()`  
  - nil 포인터/인터페이스를 기대한다.
- `HaveOccurred()`  
  - 에러가 발생했는지 기대한다.  
  - 보통 `Expect(err).ToNot(HaveOccurred())`로 “성공 경로”를 명확히 한다.
- `BeTrue()` / `BeFalse()`  
  - bool 플래그 기대.
- `ContainSubstring("...")`  
  - 문자열 일부 포함 확인(에러 메시지/응답 바디 확인).

예제
```go
Expect(2+2).To(Equal(4))
Expect(err).ToNot(HaveOccurred())
Expect(respBody).To(ContainSubstring("ok"))
Expect(got).To(BeNil())
```

## 2) 실전 핵심 matcher: Eventually / Consistently
### Eventually
- 의미: “어떤 시점 이후에 조건이 참이 될 것이다”를 기다린다.
- CAT에서 가장 흔한 위치:
  - Kubernetes Ready 대기
  - 상태 변화 이벤트가 들어올 때까지 대기
  - 장애→복구 같은 eventual consistency 검증
### Consistently
- 의미: “일정 시간 동안 조건이 계속 참이어야 한다”.
- CAT에서 흔한 위치:
  - 안정 구간(stability window) 유지 여부
  - 실패가 복구된 뒤에도 다시 깨지지 않는지 확인

예제(비동기)
```go
Eventually(func() bool {
  return atomic.LoadBool(&ready)
}, "2s", "20ms").Should(BeTrue())

Consistently(func() bool {
  return atomic.LoadBool(&ready)
}, "300ms", "30ms").Should(BeTrue())
```

## 3) 비동기/상태변화 검증 패턴(추천)
### 패턴 A: “기다리는 조건”을 함수로 분리
안티패턴:
- 조건 함수 안에서 네트워크 요청을 매번 새로 만들거나, side-effect를 넣는 것

추천:
```go
cond := func() bool { return atomic.LoadBool(&ready) }
Eventually(cond).Should(BeTrue())
```

### 패턴 B: timeout을 SLO/가정과 연결
- `Eventually`의 timeout은 “실제 수용 가능한 시간”과 연결해야 한다.
- 그렇지 않으면 “언제까지 기다렸는지”가 재현성 SLO가 되지 못한다.

## 4) custom matcher 가능 여부와 필요성
- custom matcher는 가능하지만, 보통은 “기다림 패턴/도메인 의미”를 함수로 뽑는 것으로 충분하다.
- matcher를 꼭 만들 이유:
  - 같은 assertion을 여러 테스트에서 도메인 의미 그대로 반복할 때(예: `ToBeCatPass()`).
  - 에러 메시지/출력 품질을 일정하게 유지하고 싶을 때.

결론:
- CAT에서는 matcher보다 “기다림 조건 함수 + 단언 + cat-result.json 저장” 파이프라인이 더 중요하다.
