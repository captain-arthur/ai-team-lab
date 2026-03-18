# Example 01: 가장 단순한 Ginkgo/Gomega

## 목적
- `Describe/It`에서 matcher를 쓰고, 실패가 테스트 결과(exit code)로 반영되는 최소 단위를 확인한다.

## 실행
```bash
cd projects/ginkgo-cat-minimal-implementation/07-engineering/examples/01-basic
go test -v ./...
```

## 파일
- `basic_example_test.go`

## CAT 관점 연결
- CAT Job에서도 결국 마지막은 “Expect 단언 → PASS/FAIL”로 수렴한다.
