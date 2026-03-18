# Adapter Review: k6 표준 실행 구조

**Date:** 2026-03-18

## 장점
- CAT 최소 조건 4요소를 “k6 실행 단위”로 닫아버림(주입/측정/단언/저장).
- **PASS/FAIL 권위가 단일**(k6 exit code)이라 운영 자동화에 강함.
- raw(`k6-summary.json`)와 표준(`cat-result.json`)을 분리해, **증거 보존**과 **소비/시각화**를 동시에 만족.

## 다른 도구에도 적용 가능한가?
- **YES.**
  - 공통 패턴: `tool run` → `tool raw output` → `adapter normalize` → `cat-result.json`
  - 도구별로 달라지는 건 “raw 파서”와 “metric 매핑”뿐이다.

## 부족한 부분(최소 구현 전 확인)
- exit code 캡처는 쉘/파이프라인에 의존하면 흔들릴 수 있어, CAT Runner는 **프로세스 종료 코드**를 직접 수집해야 한다.
- k6 summary의 threshold 정보 표현이 직관적이지 않을 수 있어, “실패한 threshold 상세”를 넣고 싶다면 콘솔 THRESHOLDS 섹션 파싱 또는 별도 출력 포맷 검토가 필요.

## 구현 난이도
- **낮음(최소 기능 기준)**:
  - k6 실행 + `--summary-export` 저장
  - JSON에서 p95/실패율/처리량 추출
  - `cat-result.json` 생성
