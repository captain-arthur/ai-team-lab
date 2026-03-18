package main

import (
	"encoding/json"
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

func TestExamples(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Example 05: Result Persistence")
}

type sample struct {
	latency time.Duration
	ok      bool
}

func runCustomScenario(targetURL string, requests int) ([]sample, int) {
	client := &http.Client{Timeout: 2 * time.Second}

	var samples []sample
	var errors int
	for i := 1; i <= requests; i++ {
		start := time.Now()
		resp, err := client.Get(targetURL)
		lat := time.Since(start)

		ok := err == nil && resp != nil && resp.StatusCode == http.StatusOK
		if resp != nil {
			_ = resp.Body.Close()
		}
		if !ok {
			errors++
		}
		samples = append(samples, sample{latency: lat, ok: ok})
	}
	return samples, errors
}

func averageLatencyMs(samples []sample) float64 {
	if len(samples) == 0 {
		return 0
	}
	var sum float64
	for _, s := range samples {
		sum += float64(s.latency.Milliseconds())
	}
	return sum / float64(len(samples))
}

var _ = Describe("CAT-style result persistence", func() {
	It("cat-result.json을 파일로 남기고, SLO에 따라 PASS/FAIL을 단언한다", func() {
		// scenario injection
		delayMS := catutilInt("DELAY_MS", 20)
		failEvery := catutilInt("FAIL_EVERY", 999999) // 기본 실패 없음
		requests := catutilInt("REQUESTS", 10)

		// SLO
		sloLatencyMaxMs := catutilFloat("SLO_LATENCY_MAX_AVG_MS", 300)
		sloErrorRateMax := catutilFloat("SLO_ERROR_RATE_MAX", 0.0)

		var reqCount int
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCount++
			if failEvery > 0 && reqCount%failEvery == 0 {
				http.Error(w, "forced error", http.StatusInternalServerError)
				return
			}
			time.Sleep(time.Duration(delayMS) * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		}))
		defer srv.Close()

		var samples []sample
		var errCount int

		// 항상 결과 파일을 남긴다(예: Expect 실패해도 증거 보존)
		resultsDir := filepath.Join("results", "result-persistence")
		resultPath := filepath.Join(resultsDir, "cat-result.json")

		cat := catutil.CatResult{
			TestName:     "result-persistence-example",
			Tool:         "ginkgo",
			ScenarioType: "custom-http",
			SelectedSLI:  map[string]float64{},
			SloResult:    map[string]any{},
			FinalPassFail: "FAIL",
			ExitCode:     1,
			Timestamp:    catutil.NowISO(),
			ScenarioParams: map[string]any{
				"delay_ms":    delayMS,
				"fail_every":  failEvery,
				"requests":    requests,
			},
		}

		defer func() {
			_ = catutil.WriteCatResult(resultPath, cat)
		}()

		samples, errCount = runCustomScenario(srv.URL, requests)
		avgLatencyMs := averageLatencyMs(samples)
		errorRate := catutil.ErrorRate(requests, errCount)

		cat.SelectedSLI["avg_latency_ms"] = avgLatencyMs
		cat.SelectedSLI["error_rate"] = errorRate

		latOk := avgLatencyMs <= sloLatencyMaxMs
		errOk := errorRate <= sloErrorRateMax
		cat.SloResult["latency_ok"] = latOk
		cat.SloResult["error_ok"] = errOk
		cat.SloResult["latency_slo_max_ms"] = sloLatencyMaxMs
		cat.SloResult["error_slo_rate_max"] = sloErrorRateMax

		if latOk && errOk {
			cat.FinalPassFail = "PASS"
			cat.ExitCode = 0
		}

		// defer만으로는 Expect 실패 시점에 파일 검증이 먼저 일어날 수 있으므로,
		// 파일을 한 번 먼저 써두고(증거 생성), defer는 "안전장치"로 남긴다.
		_ = catutil.WriteCatResult(resultPath, cat)

		Expect(cat.FinalPassFail).To(Equal("PASS"))

		// 파일이 실제로 생성되었는지 최소 검증(증거가 비면 CAT 이후가 어려우므로)
		b, err := os.ReadFile(resultPath)
		Expect(err).ToNot(HaveOccurred())
		var decoded catutil.CatResult
		Expect(json.Unmarshal(b, &decoded)).To(Succeed())
	})
})

func catutilInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		i, err := strconv.Atoi(v)
		if err == nil {
			return i
		}
		return def
	}
	return def
}

func catutilFloat(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return f
		}
		return def
	}
	return def
}

