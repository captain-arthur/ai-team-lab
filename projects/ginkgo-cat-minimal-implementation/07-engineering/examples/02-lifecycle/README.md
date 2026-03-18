# Example 02: lifecycle(BeforeEach/AfterEach)

## 목적
- 테스트 케이스마다 필요한 리소스를 `BeforeEach`로 만들고 `AfterEach`로 정리하는 패턴을 익힌다.
- 누락하면 리소스/포트/서버가 누적될 수 있다는 관점을 “구조로” 남긴다.

## 실행
```bash
cd projects/ginkgo-cat-minimal-implementation/07-engineering/examples/02-lifecycle
go test -v ./...
```

## 파일
- `lifecycle_example_test.go`

## CAT 관점 연결
- custom CAT Job에서도 서버/워크로드/임시 폴더 같은 “시나리오 구성 요소”는 lifecycle로 안전하게 묶는다.
