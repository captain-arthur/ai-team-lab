package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"ginkgo-cat-minimal/internal/catutil"
)

func TestCustomCAT(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Custom CAT Scenario 02: Recovery After Failure")
}

func envInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}

func envFloat(key string, def float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return def
	}
	return f
}

var _ = Describe("recovery-after-failure scenario", func() {
	It("일정 시간 실패 후 회복하는지(에러율/복구 시간) SLO로 단언하고 JSON 저장한다", func() {
		failForMS := envInt("FAIL_FOR_MS", 300)
		totalMS := envInt("TOTAL_MS", 1200)
		intervalMS := envInt("INTERVAL_MS", 50)
		successStreak := envInt("SUCCESS_STREAK", 3)

		sloMaxErrRate := envFloat("SLO_MAX_ERROR_RATE", 0.35)
		sloMaxRecoveryMs := envFloat("SLO_MAX_RECOVERY_MS", 800)

		start := time.Now()
		failUntil := start.Add(time.Duration(failForMS) * time.Millisecond)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if time.Now().Before(failUntil) {
				http.Error(w, "failing", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		}))
		defer srv.Close()

		client := &http.Client{Timeout: 2 * time.Second}

		var total int
		var errors int

		seenFailure := false
		consecutiveSuccess := 0
		var recoveryTime time.Duration // first moment of N consecutive successes after first failure

		// 요청 루프(간격 고정)
		next := start
		end := start.Add(time.Duration(totalMS) * time.Millisecond)
		for time.Now().Before(end) {
			next = next.Add(time.Duration(intervalMS) * time.Millisecond)

			reqStart := time.Now()
			resp, err := client.Get(srv.URL)
			ok := err == nil && resp != nil && resp.StatusCode == http.StatusOK
			if resp != nil {
				_ = resp.Body.Close()
			}

			total++
			if !ok {
				errors++
				consecutiveSuccess = 0
				if !seenFailure {
					seenFailure = true
				}
			} else {
				consecutiveSuccess++
				if seenFailure && recoveryTime == 0 && consecutiveSuccess >= successStreak {
					recoveryTime = time.Since(start)
				}
			}

			// 너무 빠르면 다음 tick까지 대기
			sleepFor := next.Sub(reqStart)
			if sleepFor > 0 {
				time.Sleep(sleepFor)
			}
		}

		errorRate := catutil.ErrorRate(total, errors)
		recoveryMs := float64(recoveryTime.Milliseconds())
		if !seenFailure {
			recoveryMs = 0
		}

		finalPassFail := "FAIL"
		exitCode := 1
		if errorRate <= sloMaxErrRate && recoveryMs <= sloMaxRecoveryMs {
			finalPassFail = "PASS"
			exitCode = 0
		}

		resultsDir := filepath.Join("results", "recovery-after-failure")
		resultPath := filepath.Join(resultsDir, "cat-result.json")

		cat := catutil.CatResult{
			TestName:     "custom-cat-recovery-after-failure",
			Tool:         "ginkgo",
			ScenarioType: "custom-recovery-after-failure",
			SelectedSLI: map[string]float64{
				"error_rate":        errorRate,
				"recovery_time_ms": recoveryMs,
			},
			SloResult: map[string]any{
				"error_ok":          errorRate <= sloMaxErrRate,
				"recovery_ok":       recoveryMs <= sloMaxRecoveryMs,
				"error_slo_rate_max": sloMaxErrRate,
				"recovery_slo_max_ms": sloMaxRecoveryMs,
			},
			FinalPassFail: finalPassFail,
			ExitCode:      exitCode,
			Timestamp:     catutil.NowISO(),
			ScenarioParams: map[string]any{
				"fail_for_ms":      failForMS,
				"total_ms":         totalMS,
				"interval_ms":      intervalMS,
				"success_streak":  successStreak,
			},
		}

		defer func() {
			_ = catutil.WriteCatResult(resultPath, cat)
		}()

		// Expect 실패 전에 먼저 증거 파일을 만든다.
		_ = catutil.WriteCatResult(resultPath, cat)

		Expect(finalPassFail).To(Equal("PASS"))
	})
})

