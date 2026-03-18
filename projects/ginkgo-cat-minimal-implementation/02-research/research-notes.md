# Ginkgo 기능 최소 정리(이 구현에서 사용한 것만)

**Date:** 2026-03-18

## 사용할 Ginkgo 기능(최소)
- `Describe/It` : 테스트(Job) 1개를 단위로 구성
- `RunSpecs` + `RegisterFailHandler(Fail)` : suite 실행과 실패 전파
- `DeferCleanup` : assertion 실패 여부와 무관하게 `cat-result.json`을 “항상” 저장

## assertion 방식(Gomega)
- `Expect(cat.FinalPassFail).To(Equal("PASS"), ...)`
  - SLO 조건이 하나라도 깨지면 최종 FAIL로 설정하고, Expect에서 테스트가 실패한다.

## exit code 해석 방식
- 이 구현에서는 `go test ./...`를 실행한다.
- Ginkgo suite의 FAIL은 go test의 **비0 종료 코드**로 전달된다.
- CAT의 최종 PASS/FAIL은 이 종료 코드(=테스트 성공 여부)를 권위로 사용한다.

## 결과 파일 생성 방법
- `results/cat-result.json` 경로로 JSON 저장
- 저장은 `DeferCleanup`에서 수행한다(Expect 실패여도 파일 생성 보장).
