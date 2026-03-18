# Example 03: 비동기 검증(Eventually/Consistently)

## 목적
- 상태 변화 대기(Eventually)와 안정 유지(Consistently)를 “조건 함수만”으로 분리해 쓰는 법을 익힌다.
- custom CAT에서 Ready/복구 같은 상태 검증에 그대로 연결된다.

## 실행
```bash
cd projects/ginkgo-cat-minimal-implementation/07-engineering/examples/03-async
go test -v ./...
```

## 파일
- `async_example_test.go`

## CAT 관점 연결
- “언젠가 좋아질 것”은 Eventually
- “좋은 상태가 유지될 것”은 Consistently
