package main

import (
	"encoding/json"
	"fmt"
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
)

func TestCATGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CAT Ginkgo Minimal Suite")
}

var _ = Describe("ginkgo-basic-test", func() {
	It("custom-http scenario: SLI/SLO/PASS-FAIL + cat-result.json", func() {
		// ===== CAT Job 입력(시나리오 주입) =====
		delayMS := envInt("SCENARIO_DELAY_MS", 50)
		failEvery := envInt("SCENARIO_FAIL_EVERY", 999999) // 기본은 실패 없음
		requests := envInt("SCENARIO_REQUESTS", 20)

		sloLatencyP95LikeMs := envFloat("SLO_LATENCY_MAX_MS", 300) // avg 대신 “p95 유사”로 읽되, 구현은 avg
		sloErrorRate := envFloat("SLO_ERROR_RATE_MAX", 0.0)

		// ===== SLI 수집(코드 내부 측정) =====
		var reqCount int64
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cur := atomic.AddInt64(&reqCount, 1)
			if failEvery > 0 && int(cur)%failEvery == 0 {
				http.Error(w, "forced failure", http.StatusInternalServerError)
				return
			}
			time.Sleep(time.Duration(delayMS) * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		}))
		defer srv.Close()

		type sample struct {
			latency time.Duration
			ok      bool
		}

		client := &http.Client{Timeout: 3 * time.Second}
		samples := make([]sample, 0, requests)

		// ===== PASS/FAIL과 결과 파일을 “항상” 남기기 =====
		cat := CatResult{
			TestName:     "ginkgo-basic-test",
			Tool:         "ginkgo",
			ScenarioType: "custom-http",
			SelectedSLI: SelectedSLI{
				AvgLatencyMs: 0,
				ErrorRate:    0,
			},
			SloResult: SloResult{
				LatencyOk:   false,
				ErrorOk:     false,
				LatencySloMs: sloLatencyP95LikeMs,
				ErrorSloRate: sloErrorRate,
			},
			FinalPassFail: "FAIL",
			ExitCode:      1,
			Timestamp:     nowISO(),
			ScenarioParams: ScenarioParams{
				DelayMS:   delayMS,
				FailEvery: failEvery,
				Requests:  requests,
			},
		}

		resultsDir := filepath.Join(".", "results")
		resultPath := filepath.Join(resultsDir, "cat-result.json")

		DeferCleanup(func() {
			_ = os.MkdirAll(resultsDir, 0o755)
			b, err := json.MarshalIndent(cat, "", "  ")
			if err != nil {
				// cat-result.json이 비어도 테스트 실패 원인을 남기기 위해, 여기선 콘솔 출력만 최소로 한다.
				_, _ = fmt.Fprintln(os.Stderr, "cat-result.json marshal 실패:", err)
				return
			}
			_ = os.WriteFile(resultPath, b, 0o644)
		})

		// ===== 요청 실행(시나리오 실행) =====
		url := srv.URL + "/"
		for i := 0; i < requests; i++ {
			start := time.Now()
			resp, err := client.Get(url)
			lat := time.Since(start)
			ok := err == nil && resp != nil && resp.StatusCode == http.StatusOK
			if resp != nil {
				_ = resp.Body.Close()
			}
			samples = append(samples, sample{latency: lat, ok: ok})
		}

		// ===== SLI 계산(평균 latency, error rate) =====
		var sumLatMs float64
		var errors int
		for _, s := range samples {
			sumLatMs += float64(s.latency.Milliseconds())
			if !s.ok {
				errors++
			}
		}

		avgLatMs := sumLatMs / float64(len(samples))
		errorRate := float64(errors) / float64(len(samples))

		cat.SelectedSLI.AvgLatencyMs = avgLatMs
		cat.SelectedSLI.ErrorRate = errorRate

		latencyOk := avgLatMs <= sloLatencyP95LikeMs
		errorOk := errorRate <= sloErrorRate

		cat.SloResult.LatencyOk = latencyOk
		cat.SloResult.ErrorOk = errorOk

		if latencyOk && errorOk {
			cat.FinalPassFail = "PASS"
			cat.ExitCode = 0
		}

		// ===== SLO 단언(SPASS/FAIL 권위: Ginkgo 테스트 성공 여부) =====
		Expect(cat.FinalPassFail).To(Equal("PASS"), "cat-result.json의 SLO 평가가 PASS가 아니면 FAIL로 처리한다")
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

