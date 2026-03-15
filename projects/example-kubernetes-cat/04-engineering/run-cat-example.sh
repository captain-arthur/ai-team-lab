#!/usr/bin/env bash
#
# 예제 CAT 실행 스크립트 (구현 가이드용)
# Architecture: Conformance(Sonobuoy) → Custom → 수집 → 리포팅 순서를 보여 줌.
# 실제 사용 시 KUBECONFIG, RESULTS_BASE, Sonobuoy 모드 등을 환경에 맞게 수정하세요.
#
# Usage:
#   ./run-cat-example.sh [quick|conformance]
#   기본: quick

set -euo pipefail

MODE="${1:-quick}"
RUN_ID=$(date +%Y%m%d-%H%M%S)
RESULTS_BASE="${RESULTS_BASE:-./results/cat}"
RESULTS_ROOT="${RESULTS_BASE}/${RUN_ID}"

echo "[CAT] Run ID: ${RUN_ID}"
echo "[CAT] Mode: ${MODE}"
echo "[CAT] Results: ${RESULTS_ROOT}"

mkdir -p "${RESULTS_ROOT}"/{sonobuoy,custom}

# --- Step 1: Sonobuoy (Conformance runner) ---
echo "[CAT] Step 1: Sonobuoy run (${MODE})..."
sonobuoy run --mode "${MODE}"
sonobuoy wait
OUT=$(sonobuoy retrieve)
mv "$OUT" "${RESULTS_ROOT}/sonobuoy/sonobuoy_${RUN_ID}.tar.gz"
echo "[CAT] Sonobuoy tarball: ${RESULTS_ROOT}/sonobuoy/sonobuoy_${RUN_ID}.tar.gz"

# --- Step 2: Custom tests (있을 경우) ---
echo "[CAT] Step 2: Custom tests..."
CUSTOM_EXIT=0
if [[ -x "./custom-tests/run-custom-tests.sh" ]]; then
  ./custom-tests/run-custom-tests.sh "${RESULTS_ROOT}/custom" || CUSTOM_EXIT=$?
  echo "custom_tests_exit_code=${CUSTOM_EXIT}" >> "${RESULTS_ROOT}/custom/summary.txt"
else
  echo "No custom-tests/run-custom-tests.sh found; skipping custom phase."
  echo "custom_tests_exit_code=0" > "${RESULTS_ROOT}/custom/summary.txt"
fi

# --- Step 3: Collector (이미 동일 트리 아래에 저장됨) ---
# 추가 수집 로직이 필요하면 여기에 (예: tarball 압축 해제).

# --- Step 4: Reporter (간단 요약) ---
SONOBUOY_SUMMARY=$(sonobuoy results "${RESULTS_ROOT}/sonobuoy/sonobuoy_${RUN_ID}.tar.gz" 2>/dev/null || echo " sonobuoy results failed")
{
  echo "# CAT Run Report: ${RUN_ID}"
  echo ""
  echo "## Conformance (Sonobuoy)"
  echo '```'
  echo "${SONOBUOY_SUMMARY}"
  echo '```'
  echo ""
  echo "## Custom"
  echo "- exit_code: ${CUSTOM_EXIT}"
  echo "- logs: custom/"
  echo ""
  if [[ "$CUSTOM_EXIT" -eq 0 ]]; then
    echo "## Overall: determine from Sonobuoy output above (Conformance + Custom)"
  else
    echo "## Overall: fail (Custom tests failed)"
  fi
} > "${RESULTS_ROOT}/report.md"

echo "[CAT] Report: ${RESULTS_ROOT}/report.md"
echo "[CAT] Done."
