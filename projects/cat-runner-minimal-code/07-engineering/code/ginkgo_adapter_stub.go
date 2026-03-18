package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type GinkgoAdapter struct{}

func (a GinkgoAdapter) Run(job JobSpec) (int, error) {
	if job.Output.Dir == "" {
		return 1, errors.New("job.output.dir is required to run ginkgo")
	}
	// runner 실행 디렉토리에서 ./ginkgo_cat_job 를 찾는다.
	// (문서의 end-to-end 실행 단계에서 동일 디렉토리에서 go run .을 수행한다.)

	cfg := job.Scenario.Config
	requests := cfgInt(cfg, "requests", 50)
	delayMS := cfgInt(cfg, "delay_ms", 20)
	failEvery := cfgInt(cfg, "fail_every", 999999) // 기본: 실패 없음

	sloP95, sloFailRate, err := requiredSLO(job.SLO)
	if err != nil {
		return 1, err
	}

	rawDirAbs, err := filepath.Abs(job.Output.Dir)
	if err != nil {
		return 1, err
	}
	env := os.Environ()
	env = appendEnv(env, "CAT_OUTPUT_DIR", rawDirAbs)
	env = appendEnv(env, "CAT_REQUESTS", fmt.Sprintf("%d", requests))
	env = appendEnv(env, "CAT_DELAY_MS", fmt.Sprintf("%d", delayMS))
	env = appendEnv(env, "CAT_FAIL_EVERY", fmt.Sprintf("%d", failEvery))
	env = appendEnv(env, "SLO_LATENCY_MAX_MS", fmt.Sprintf("%g", sloP95))
	env = appendEnv(env, "SLO_ERROR_RATE_MAX", fmt.Sprintf("%g", sloFailRate))

	// go test 기반으로 ginkgo suite를 실행한다.
	cmd := exec.Command("go", "test", "./ginkgo_cat_job", "-run", "TestCATGinkgo", "-count=1")
	cmd.Env = env

	runErr := cmd.Run()
	return toolExitCode(runErr), nil
}

func (a GinkgoAdapter) LocateRawResult(job JobSpec) (RawRef, error) {
	if job.Output.Dir == "" {
		return RawRef{}, errors.New("job.output.dir is required to locate ginkgo raw result")
	}
	outDirAbs, err := filepath.Abs(job.Output.Dir)
	if err != nil {
		return RawRef{}, err
	}
	rawPath := filepath.Join(outDirAbs, "ginkgo-raw.json")
	return RawRef{
		Format: "json",
		Path:   rawPath,
	}, nil
}

func (a GinkgoAdapter) ParseRawResult(job JobSpec, raw RawRef) (ParsedSLI, error) {
	sloP95, sloFailRate, err := requiredSLO(job.SLO)
	if err != nil {
		return ParsedSLI{}, err
	}

	b, err := os.ReadFile(raw.Path)
	if err != nil {
		return ParsedSLI{}, err
	}

	var r ginkgoRaw
	if err := json.Unmarshal(b, &r); err != nil {
		return ParsedSLI{}, err
	}

	selected := map[string]float64{
		"latency_p95_ms": r.LatencyP95Ms,
		"error_rate":     r.ErrorRate,
		"throughput_rps": r.ThroughputRps,
	}

	latOk := r.LatencyP95Ms <= sloP95
	errOk := r.ErrorRate <= sloFailRate

	sloResult := map[string]any{
		"latency_p95_ms": map[string]any{
			"measured": r.LatencyP95Ms,
			"slo_max":  sloP95,
			"ok":       latOk,
		},
		"error_rate": map[string]any{
			"measured": r.ErrorRate,
			"slo_max":  sloFailRate,
			"ok":       errOk,
		},
	}

	return ParsedSLI{
		SelectedSLI: selected,
		SloResult:   sloResult,
	}, nil
}

func (a GinkgoAdapter) BuildCatResult(job JobSpec, parsed ParsedSLI, raw RawRef, exitCode int) CatResult {
	final := "FAIL"
	if exitCode == 0 {
		final = "PASS"
	}

	return CatResult{
		TestName:      job.TestName,
		Tool:          "ginkgo",
		ScenarioType:  job.Scenario.Type,
		SelectedSLI:   parsed.SelectedSLI,
		SloResult:     parsed.SloResult,
		FinalPassFail: final,
		ExitCode:      exitCode,
		RawResult:     RawResult{Format: raw.Format, Path: raw.Path},
	}
}

type ginkgoRaw struct {
	LatencyP95Ms  float64 `json:"latency_p95_ms"`
	ErrorRate     float64 `json:"error_rate"`
	ThroughputRps float64 `json:"throughput_rps"`

	Requests  int `json:"requests"`
	DelayMs   int `json:"delay_ms"`
	FailEvery int `json:"fail_every"`
}

func cfgInt(cfg map[string]any, key string, def int) int {
	if cfg == nil {
		return def
	}
	v, ok := cfg[key]
	if !ok || v == nil {
		return def
	}
	switch t := v.(type) {
	case int:
		return t
	case int64:
		return int(t)
	case float64:
		return int(t)
	default:
		return def
	}
}
