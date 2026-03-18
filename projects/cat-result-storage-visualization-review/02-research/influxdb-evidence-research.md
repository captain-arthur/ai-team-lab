# InfluxDB + Evidence 관점 리서치( CAT 결과 시각화용 )

**Date:** 2026-03-19

## 1) InfluxDB가 무엇인가(저장/조회 관점)
- InfluxDB는 **시계열 데이터(Time Series Data)** 를 저장하는 데이터베이스다.
- 데이터 모델 핵심 개념
  - `measurement`: 데이터의 “그룹/테이블 이름” 역할(문자열)
  - `tags`: 인덱싱되는 문자열 메타데이터(쿼리로 필터링/그룹화에 사용)
  - `fields`: 실제 값(부동소수/정수/문자열/불리언 등), 보통 인덱싱되지 않음
  - `timestamp`: 각 포인트의 시간. InfluxDB는 내부적으로 `_time`을 가진다.
- 왜 시계열에 쓰나
  - 같은 종류의 metric을 시간축으로 대량 기록하고, 구간별 집계/조회가 필요하기 때문.
- 스키마 설계가 중요한 이유
  - `tags`/`measurement` 선택이 “시리즈 수(cardinality)”와 쿼리 성능을 좌우한다.

## 2) k6와 InfluxDB 연동(흐름/데이터가 무엇으로 들어가는가)
- k6는 InfluxDB에 “실행 중 실시간으로 metric을 밀어넣는” 형태로 연동 가능하다.
- k6→InfluxDB가 보내는 것(개념)
  - k6가 실행하며 생성하는 metric(예: HTTP 요청 지연, 실패율, 요청률 등)을
  - InfluxDB의 `bucket/measurement/tags/fields` 형태로 기록해두고,
  - 이후 Grafana/Influx 도구로 쿼리/시각화한다.
- 기본 흐름
  - (1) k6 실행
  - (2) k6가 metric을 버퍼링/flush하며 InfluxDB로 전송
  - (3) Evidence/Grafana가 쿼리로 데이터를 뽑아 시각화
- 장점
  - “실행 중/직후”의 시계열 관측을 쉽게 붙일 수 있다.
  - k6 metric 집계 구조가 InfluxDB 쿼리 모델과 잘 맞는다.
- 제약/주의
  - CAT의 최종 PASS/FAIL(게이트)은 다른 레이어에서도 동일하게 결정돼야 한다.
  - InfluxDB에 저장된 값은 “측정 관측 데이터”이지, CAT 최종 판정의 권위가 될 필요는 없다(권위는 tool exit code로 고정하는 정책이 필요).

## 3) Evidence는 어떤 방식으로 데이터를 시각화하는가
- Evidence는 SQL+마크다운 기반으로 “데이터 제품(대시보드/리포트)”을 렌더링하는 프레임워크다.
- 핵심 동작
  - Evidence app이 마크다운 페이지를 렌더링
  - 페이지 안의 SQL 코드가 data source에 대해 쿼리를 실행
  - 쿼리 결과로 차트/컴포넌트를 렌더링
- 데이터 소스 처리
  - Evidence는 data source를 공통 저장소(Parquet)로 “추출(sources)”하고
  - 이후 SQL로 조회하는 방식이다.
- CAT 결과 시각화 적합성
  - Evidence가 SQL을 실행할 수 있는 “데이터 소스 커넥터”가 있어야 한다.
  - InfluxDB가 Evidence에서 바로 “SQL로” 쿼리되는지 여부는 확정이 필요하다(문서에 InfluxDB가 명시적으로 보이지 않으므로, 확실하지 않은 상태로 본다).
  - 따라서 Evidence를 CAT에 붙이려면:
    - (a) Evidence가 InfluxDB를 직접 data source로 지원하는지 확인하거나
    - (b) InfluxDB에서 필요한 데이터를 파일(Parquet/CSV/SQL DB)로 내보내 Evidence가 읽는 경로를 둬야 한다.

## 4) CAT 관점 결론(리서치로 확정 가능한 범위)
- InfluxDB는 “시계열 관측 데이터 저장소”로 CAT에 자연스럽게 붙는다.
- Evidence는 “SQL+마크다운 렌더링” 레이어로 붙는다.
- 단, Evidence가 InfluxDB를 직접 쿼리할지, 아니면 InfluxDB→중간 추출(파일/SQL DB) 단계를 둘지에 따라 복잡도가 달라진다.
