# Example 04: table-driven 테스트(DescribeTable)

## 목적
- 입력/기대값 테이블로 분기 폭발을 막는 패턴을 익힌다.
- CAT에서 “SLI→SLO 평가 규칙”이 케이스 기반이면 DescribeTable이 잘 맞는다.

## 실행
```bash
cd projects/ginkgo-cat-minimal-implementation/07-engineering/examples/04-table-driven
go test -v ./...
```

## 파일
- `table_driven_example_test.go`

## CAT 관점 연결
- threshold/게이트 로직이 여러 케이스로 나뉘면 table-driven이 유지보수에 유리하다.
