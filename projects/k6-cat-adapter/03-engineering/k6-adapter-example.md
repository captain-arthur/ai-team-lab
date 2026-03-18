# k6 Adapter 실행 예시(전체 흐름이 실제로 동작)

**Date:** 2026-03-18  
**목표:** `k6 run → k6-summary.json → cat-result.json` 흐름을 실제로 재현

## 1) 실행할 k6 테스트
- **스크립트**: `scripts/http-basic.js`
- **타깃**: 공개 엔드포인트 `https://test.k6.io/` (예시)

## 2) k6 실행(요약 JSON 생성)
```bash
cd projects/k6-cat-adapter/03-engineering

OUTDIR="results/run-$(date +%Y%m%d-%H%M%S)"
mkdir -p "$OUTDIR"

TEST_NAME="http-basic"
TARGET="https://test.k6.io/"
DURATION="20s"
TARGET_RPS="80"
SLO_P95_MS="400"
SLO_FAIL_RATE="0.01"

TARGET_URL="$TARGET" MODE="arrival" DURATION="$DURATION" TARGET_RPS="$TARGET_RPS" \
SLO_P95_MS="$SLO_P95_MS" SLO_FAIL_RATE="$SLO_FAIL_RATE" \
k6 run --summary-export "$OUTDIR/k6-summary.json" scripts/http-basic.js | tee "$OUTDIR/k6-output.txt"
```

## 3) cat-result.json 생성(변환)
```bash
python3 scripts/k6_summary_to_cat_result.py \
  --k6-summary "$OUTDIR/k6-summary.json" \
  --out "$OUTDIR/cat-result.json" \
  --test-name "$TEST_NAME" \
  --scenario-type "http" \
  --target "$TARGET" \
  --exit-code 0 \
  --slo-latency-p95-ms "$SLO_P95_MS" \
  --slo-error-rate "$SLO_FAIL_RATE"
```

> 실제 구현에서는 `--exit-code`를 “k6 실행 종료 코드”로 자동 주입한다(쉘에 따라 파이프라인 exit code 캡처 방식이 달라질 수 있으므로, CAT Runner에서 프로세스 종료 코드를 직접 받는 방식을 권장).

## 4) 생성되는 결과 파일
- `k6-summary.json`: k6 raw 결과(증거)
- `k6-output.txt`: 콘솔 로그
- `cat-result.json`: CAT 표준 결과(판정/비교/누적용)
