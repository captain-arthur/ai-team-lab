# Final Report: CAT 결과 정규화 전략(최종 추천)

## 정규화 핵심 원칙
1. **Raw vs Normalized를 분리**한다.
   - raw는 증거, normalized는 비교/누적/시각화의 단일 표준.
2. Normalized는 **도구 독립적인 최소 스키마**로만 고정한다.
   - `metrics[]`는 “selected SLI만” 넣는다.
3. **PASS/FAIL 권위는 재판정하지 않는다.**
   - tool exit code(또는 테스트 성공/실패)가 최종이며, adapter는 기록만 한다.
4. Evidence 친화성은 JSON이 어렵다면 “선택적으로 CSV 파생”으로 맞춘다.

## 가장 적절한 포맷
- **JSON(옵션 A)**: `cat-result.json`을 표준 정규화 포맷으로 사용한다.

## InfluxDB 사용 범위
- InfluxDB/시계열 저장소는 normalized PASS/FAIL의 “권위”가 아니다.
- 필요할 때만 “관측/원인 분류”용 시계열 데이터로 별도 저장한다.

## Evidence 연결 위치
- Evidence는 기본적으로 `cat-result.json`(또는 runner가 생성한 `cat-result-table.csv`)를 읽어서
  - selected_sli
  - status(PASS/FAIL)
  를 기반으로 리포트를 렌더링한다.

## 최종 추천안
- 1단계: JSON 표준화(`cat-result.json`)를 단일 소스 오브 트루스로 고정
- 2단계(필요 시): CSV `cat-result-table.csv`를 추가로 생성(파생/복제 수준, ETL 금지)

## 최종 결론(요구 한 문장)
👉 **CAT 결과 정규화는 “단일 run 단위 JSON(cat-result.json) + 필요 시 CSV 파생” 구조가 가장 단순하고 실용적이다.**

## Final Question 답
👉 “CAT 결과를 일관성 있게 만드는 최적의 정규화 전략은 무엇인가?”

**결론: JSON 중심(필수) + CSV 파생(선택) 혼합 전략이 최적이다.**
