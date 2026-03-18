# CAT 결과 저장/시각화 3옵션 비교(파일 vs InfluxDB vs 혼합)

**Date:** 2026-03-19

## 옵션 A. 파일 기반 only
- 구조
  - tool raw 결과 파일 보존
  - `cat-result.json`(최종 판정/selected_sli) 저장
  - Evidence는 파일(또는 추출된 CSV/Parquet) 기반으로 시각화

### 구조 단순성
- 높음. adapter는 “정규화/저장”만 하면 된다.

### adapter 복잡도
- 낮음. k6/CL2/Ginkgo는 각자 raw를 남기고 adapter가 cat-result만 만든다.

### k6 친화성 / 비-k6 도구 수용성
- k6는 summary-export JSON이 이미 있으므로 잘 맞는다.
- CL2/Ginkgo도 파일 기반으로 이미 수용 가능.

### 결과 비교/누적 용이성
- 좋음. run-id/테스트명 기반으로 diff 가능.

### 운영 부담
- 낮음. DB 운영이 불필요.

### 시각화 적합성
- “최종 지표 + 판정” 중심이면 충분.
- “실행 중 세밀한 시계열(초단위/분포 변화)”이 필수라면 한계.

## 옵션 B. InfluxDB 중심
- 구조
  - tool 실행 중/후에 metric을 InfluxDB에 적재
  - Evidence는 InfluxDB 쿼리 기반으로 시각화(직접 또는 추출)
  - 최종 판정은 여전히 `cat-result.json`에 저장(권위는 exit code로 고정)

### 구조 단순성
- 중간~낮음. InfluxDB 버킷/measurement/tags 스키마 설계가 필요.

### adapter 복잡도
- 높음. k6는 비교적 쉽지만, CL2/Ginkgo는 “Influx 스키마에 맞게 time series를 어떤 형태로 넣을지”가 추가 결정이다.

### k6 친화성 / 비-k6 도구 수용성
- k6는 친화적(InfluxDB output 확장/실시간 push).
- CL2/Ginkgo는 억지로 맞추게 될 위험이 있다.

### 결과 비교/누적 용이성
- 시계열 비교는 좋다.
- 하지만 CAT의 “최종 판정 파일(cat-result)”이 같이 있어야 운영자가 납득한다.

### 운영 부담
- 중간~높음. DB 운영/권한/카디널리티 관리 필요.

### 시각화 적합성
- “실행 중 관측”과 “시간 흐름”이 중요하면 강력.

## 옵션 C. 혼합 구조(권장 후보)
- 구조
  - 모든 tool: raw 파일 + `cat-result.json`은 파일에 저장(최종 판정의 권위 유지)
  - 추가로, k6(및 필요할 때만) 일부 시계열만 InfluxDB에 적재
  - Evidence는
    - “최종 결과/판정”은 `cat-result.json` 기반
    - “시간 흐름 지표”는 InfluxDB 기반(또는 InfluxDB에서 내보낸 Parquet/CSV 기반)

### 구조 단순성
- 중간. DB 운영이 있지만, adapter 범위를 제한해 단순성을 유지.

### adapter 복잡도
- 중간. Influx 적재는 “k6에만 먼저” 허용하거나, CL2/Ginkgo는 추후 단계로 둔다.

### k6 친화성 / 비-k6 도구 수용성
- k6는 자연스럽다.
- CL2/Ginkgo는 “최종 판정 파일” 중심으로 억지 정합을 피한다.

### 결과 비교/누적 용이성
- 좋음. 파일 기반으로 비교/회귀 판단 가능.

### 운영 부담
- 낮~중간. InfluxDB를 “필수 필드”만 제한적으로 사용.

### 시각화 적합성
- 가장 실무적.
  - 운영/개발자가 보는 것은 결국 PASS/FAIL과 selected_sli
  - 동시에 실행 중 추적이 필요할 때만 시계열을 붙인다.

## 결론(이번 단계에서의 비교 결과)
- “단순성과 도구 수용성”을 최우선하면 옵션 A가 가장 단순.
- “실행 중 시계열 시각화”가 필수로 올라오면 옵션 C가 균형이 좋다.
- 옵션 B는 초기 도입/유지 비용 대비 이득이 불명확해 과하다.
