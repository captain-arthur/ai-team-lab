# 구현이 아닌 “최소 데이터 흐름” 설계(3 후보)

**Date:** 2026-03-19

## 공통 전제(고정)
- CAT의 최종 판정/근거는 반드시 파일로 남긴다.
  - `cat-result.json` (selected_sli + slo_result + final_pass_fail + raw_result_path)
- Evidence/시각화는 이 파일을 “소비”하는 쪽에 둔다.

## 옵션 1: 파일 기반 Evidence(가장 단순)
흐름
1. tool 실행
2. raw 결과 파일 저장
3. adapter가 `cat-result.json` 생성(정규화)
4. Evidence는 `cat-result.json`(또는 raw에서 추출한 CSV/Parquet)만 읽어 시각화

필요 구성요소
- adapter: cat-result 생성
- Evidence: 파일/추출된 데이터 sources 연결

복잡한 부분
- 파일에서 필요한 “시각화용 열”을 뽑는 쿼리 구성(대부분 SQL+CSV/Parquet에서 해결)

단순한 부분
- DB 운영 없음
- 도구별 공통화가 `cat-result`에서 끝남

실무 적용 가능성
- 기본 “판정/요약/회귀 비교”에는 충분

## 옵션 2: 혼합(권장): 파일 + (필요할 때만) k6 시계열 Influx insert
흐름
1. tool 실행
2. raw 파일 + `cat-result.json` 저장(고정)
3. k6 실행은 선택적으로 InfluxDB로 실시간 metric 적재(필요한 경우만)
4. Evidence는
   - 기본 리포트: `cat-result.json`
   - 원인 분류/추세 페이지: InfluxDB에서 시계열 추출 데이터 기반

필요 구성요소
- k6는 선택적 Influx output(옵션)
- adapter는 cat-result 정규화만
- Evidence는 cat-result + (필요 시) 시계열 추출 데이터 소스를 둘 다 처리

복잡한 부분
- “어떤 시계열을 Evidence가 쓸 수 있게 만들지”가 추가 의사결정

단순한 부분
- CL2/Ginkgo는 Influx 정합을 강제로 맞출 필요가 없다.

실무 적용 가능성
- CAT의 핵심(판정 파일)과 디버깅(시계열)을 동시에 충족하면서도 adapter 범위를 제한한다.

## 옵션 3: k6만 Influx 직적재, 나머지는 파일 후 적재(비권장/경계)
흐름
1. k6는 Influx 직적재
2. CL2/Ginkgo는 파일 저장 후 adapter가 Influx에 “동일 스키마”로 적재
3. Evidence는 Influx 기반으로 통합 시각화

복잡한 부분
- CL2/Ginkgo adapter에서 Influx 스키마/태그/필드 설계가 커짐

단순한 부분
- “사용자 관점 대시보드 하나”로 만들 수 있음

실무 적용 가능성
- adapter 복잡도 증가로 비용 대비 효과가 불확실 → 단계적 접근 권장

