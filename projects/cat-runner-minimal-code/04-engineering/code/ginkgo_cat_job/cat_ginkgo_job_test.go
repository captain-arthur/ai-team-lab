package ginkgo_cat_job

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCATGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CAT Ginkgo Job")
}

func envInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	iv, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return iv
}

func envFloat(key string, def float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	fv, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return def
	}
	return fv
}

func percentileMs(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	if p <= 0 {
		return sorted[0]
	}
	if p >= 100 {
		return sorted[len(sorted)-1]
	}
	// nearest-rank style.
	rank := int(math.Ceil(p/100.0*float64(len(sorted)))) - 1
	if rank < 0 {
		rank = 0
	}
	if rank >= len(sorted) {
		rank = len(sorted) - 1
	}
	return sorted[rank]
}

var _ = Describe("custom-cat-ginkgo-job", func() {
	It("SLO 단언(PASS/FAIL authority)을 exit code로 남기고 raw JSON을 생성한다", func() {
		requests := envInt("CAT_REQUESTS", 50)
		delayMS := envInt("CAT_DELAY_MS", 20)
		failEvery := envInt("CAT_FAIL_EVERY", 999999)

		sloLatencyP95 := envFloat("SLO_LATENCY_MAX_MS", 450)
		sloErrorRate := envFloat("SLO_ERROR_RATE_MAX", 0.01)

		outDir := os.Getenv("CAT_OUTPUT_DIR")
		Expect(outDir).ToNot(BeEmpty(), "CAT_OUTPUT_DIR is required")

		rawPath := filepath.Join(outDir, "ginkgo-raw.json")
		Expect(os.MkdirAll(outDir, 0o755)).To(Succeed())

		var latencyP95Ms float64
		var errorRate float64
		var throughputRps float64

		DeferCleanup(func() {
			_ = writeJSON(rawPath, map[string]any{
				"latency_p95_ms": latencyP95Ms,
				"error_rate":     errorRate,
				"throughput_rps": throughputRps,
				"requests":       requests,
				"delay_ms":       delayMS,
				"fail_every":    failEvery,
			})
		})

		var reqIndex int
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqIndex++
			time.Sleep(time.Duration(delayMS) * time.Millisecond)
			if failEvery > 0 && reqIndex%failEvery == 0 {
				w.WriteHeader(500)
				_, _ = w.Write([]byte("fail"))
				return
			}
			w.WriteHeader(200)
			_, _ = w.Write([]byte("ok"))
		}))
		defer srv.Close()

		client := &http.Client{Timeout: time.Duration(requests)*time.Second + 5*time.Second}

		latencies := make([]float64, 0, requests)
		errCount := 0
		start := time.Now()

		for i := 0; i < requests; i++ {
			reqStart := time.Now()
			resp, err := client.Get(srv.URL)
			latMs := float64(time.Since(reqStart).Microseconds()) / 1000.0
			latencies = append(latencies, latMs)
			if err != nil {
				errCount++
				continue
			}
			if resp.StatusCode >= 400 {
				errCount++
			}
			_ = resp.Body.Close()
		}
		total := time.Since(start)
		throughputRps = float64(requests) / total.Seconds()

		sort.Float64s(latencies)
		latencyP95Ms = percentileMs(latencies, 95)
		errorRate = float64(errCount) / float64(requests)

		finalPassFail := "FAIL"
		if latencyP95Ms <= sloLatencyP95 && errorRate <= sloErrorRate {
			finalPassFail = "PASS"
		}

		Expect(finalPassFail).To(Equal("PASS"))
	})
})

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

