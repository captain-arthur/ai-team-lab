package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"ginkgo-cat-minimal/internal/catutil"
)

func TestCustomCAT(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Custom CAT Scenario 01: Stability Window")
}

var _ = Describe("stability-window scenario", func() {
	It("준비(ready) 이후 안정 구간이 유지되는지 검증하고 JSON으로 남긴다", func() {
		warmupMS := envInt("WARMUP_MS", 120)
		stabilityMS := envInt("STABILITY_MS", 400)
		requestsInWindow := envInt("REQUESTS_IN_WINDOW", 10)

		sloReadyMaxMs := envFloat("SLO_READY_MAX_MS", 250)
		sloErrorRateMax := envFloat("SLO_ERROR_RATE_MAX", 0.0)

		var ready atomic.Bool
		ready.Store(false)
		go func() {
			time.Sleep(time.Duration(warmupMS) * time.Millisecond)
			ready.Store(true)
		}()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ready.Load() {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("ready"))
				return
			}
			http.Error(w, "warming", http.StatusServiceUnavailable)
		}))
		defer srv.Close()

		start := time.Now()
		client := &http.Client{Timeout: 2 * time.Second}

		// 1) ready latency 측정: 최초 200 응답까지
		var readyLatency time.Duration
		for {
			resp, err := client.Get(srv.URL)
			if resp != nil {
				_ = resp.Body.Close()
			}
			if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
				readyLatency = time.Since(start)
				break
			}
			time.Sleep(10 * time.Millisecond)
		}

		// 2) stability window 동안 error rate 측정
		interval := time.Duration(0)
		if requestsInWindow > 1 {
			interval = time.Duration(stabilityMS/int(requestsInWindow-1)) * time.Millisecond
		}

		var errors int
		var latSum float64
		for i := 0; i < requestsInWindow; i++ {
			reqStart := time.Now()
			resp, err := client.Get(srv.URL)
			lat := time.Since(reqStart)
			latSum += float64(lat.Milliseconds())

			ok := err == nil && resp != nil && resp.StatusCode == http.StatusOK
			if resp != nil {
				_ = resp.Body.Close()
			}
			if !ok {
				errors++
			}

			if interval > 0 && i < requestsInWindow-1 {
				time.Sleep(interval)
			}
		}

		avgLatencyMs := latSum / float64(requestsInWindow)
		errorRate := catutil.ErrorRate(requestsInWindow, errors)

		finalPassFail := "FAIL"
		exitCode := 1
		if float64(readyLatency.Milliseconds()) <= sloReadyMaxMs && errorRate <= sloErrorRateMax {
			finalPassFail = "PASS"
			exitCode = 0
		}

		resultsDir := filepath.Join("results", "stability-window")
		resultPath := filepath.Join(resultsDir, "cat-result.json")

		cat := catutil.CatResult{
			TestName:     "custom-cat-stability-window",
			Tool:         "ginkgo",
			ScenarioType: "custom-stability-window",
			SelectedSLI: map[string]float64{
				"ready_latency_ms": float64(readyLatency.Milliseconds()),
				"avg_latency_ms":   avgLatencyMs,
				"error_rate":       errorRate,
			},
			SloResult: map[string]any{
				"ready_ok":        float64(readyLatency.Milliseconds()) <= sloReadyMaxMs,
				"error_ok":        errorRate <= sloErrorRateMax,
				"ready_slo_max_ms": sloReadyMaxMs,
				"error_slo_rate_max": sloErrorRateMax,
			},
			FinalPassFail: finalPassFail,
			ExitCode:      exitCode,
			Timestamp:     catutil.NowISO(),
			ScenarioParams: map[string]any{
				"warmup_ms":          warmupMS,
				"stability_ms":       stabilityMS,
				"requests_in_window": requestsInWindow,
			},
		}

		defer func() {
			_ = catutil.WriteCatResult(resultPath, cat)
		}()

		// Expect 실패 타이밍에 상관없이 증거 파일을 먼저 남긴다.
		_ = catutil.WriteCatResult(resultPath, cat)

		Expect(finalPassFail).To(Equal("PASS"))
	})
})

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

