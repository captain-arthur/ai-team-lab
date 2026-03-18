package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type k6Summary struct {
	Metrics map[string]any `json:"metrics"`
}

type K6Adapter struct{}

func (a K6Adapter) Run(job JobSpec) (int, error) {
	// 1) Job → k6 실행에 필요한 값 구성
	sloP95, sloFailRate, err := requiredSLO(job.SLO)
	if err != nil {
		return 1, err
	}
	if job.Scenario.Target == "" {
		return 1, errors.New("scenario.target is required for k6")
	}
	if job.Entry == "" {
		return 1, errors.New("job.entry is required for k6")
	}

	outDir := job.Output.Dir
	rawPath := filepath.Join(outDir, "k6-summary.json")
	outputTxt := filepath.Join(outDir, "k6-output.txt")

	mode := job.Scenario.LoadModel
	if mode == "" {
		mode = "arrival"
	}

	env := os.Environ()
	env = appendEnv(env, "TARGET_URL", job.Scenario.Target)
	env = appendEnv(env, "MODE", mode)
	if job.Scenario.RPS > 0 {
		env = appendEnv(env, "TARGET_RPS", fmt.Sprintf("%d", job.Scenario.RPS))
	}
	if job.Scenario.VUS > 0 {
		env = appendEnv(env, "VUS", fmt.Sprintf("%d", job.Scenario.VUS))
	}
	if job.Scenario.Duration != "" {
		env = appendEnv(env, "DURATION", job.Scenario.Duration)
	}
	env = appendEnv(env, "SLO_P95_MS", fmt.Sprintf("%g", sloP95))
	env = appendEnv(env, "SLO_FAIL_RATE", fmt.Sprintf("%g", sloFailRate))

	// 2) 실행: k6 run --summary-export
	cmd := exec.Command("k6", "run", "--summary-export", rawPath, job.Entry)
	cmd.Env = env

	f, err := os.Create(outputTxt)
	if err != nil {
		return 1, err
	}
	defer f.Close()
	cmd.Stdout = f
	cmd.Stderr = f

	runErr := cmd.Run()
	return toolExitCode(runErr), nil
}

func (a K6Adapter) LocateRawResult(job JobSpec) (RawRef, error) {
	if job.Output.Dir == "" {
		return RawRef{}, errors.New("job.output.dir is required to locate k6 raw result")
	}
	rawPath := filepath.Join(job.Output.Dir, "k6-summary.json")
	return RawRef{
		Format: "json",
		Path:   rawPath,
	}, nil
}

func (a K6Adapter) ParseRawResult(job JobSpec, raw RawRef) (ParsedSLI, error) {
	sloP95, sloFailRate, err := requiredSLO(job.SLO)
	if err != nil {
		return ParsedSLI{}, err
	}

	selected, sloResult, parseErr := parseK6Summary(raw.Path, sloP95, sloFailRate)
	if parseErr != nil {
		return ParsedSLI{
			SelectedSLI: map[string]float64{},
			SloResult: map[string]any{
				"parse_error": parseErr.Error(),
			},
		}, nil
	}

	return ParsedSLI{
		SelectedSLI: selected,
		SloResult:   sloResult,
	}, nil
}

func (a K6Adapter) BuildCatResult(job JobSpec, parsed ParsedSLI, raw RawRef, exitCode int) CatResult {
	final := "FAIL"
	if exitCode == 0 {
		final = "PASS"
	}

	return CatResult{
		TestName:      job.TestName,
		Tool:          "k6",
		ScenarioType:  job.Scenario.Type,
		SelectedSLI:   parsed.SelectedSLI,
		SloResult:     parsed.SloResult,
		FinalPassFail: final,
		ExitCode:      exitCode,
		RawResult:     RawResult{Format: raw.Format, Path: raw.Path},
		Timestamp:     "",
	}
}

func appendEnv(env []string, k, v string) []string {
	prefix := k + "="
	for i := range env {
		if len(env[i]) >= len(prefix) && env[i][:len(prefix)] == prefix {
			env[i] = prefix + v
			return env
		}
	}
	return append(env, prefix+v)
}

func toolExitCode(runErr error) int {
	if runErr == nil {
		return 0
	}
	var ee *exec.ExitError
	if errors.As(runErr, &ee) {
		return ee.ExitCode()
	}
	return 1
}

func parseK6Summary(summaryPath string, sloP95 float64, sloFailRate float64) (map[string]float64, map[string]any, error) {
	b, err := os.ReadFile(summaryPath)
	if err != nil {
		return nil, nil, err
	}
	var s k6Summary
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, nil, err
	}
	metrics := s.Metrics
	if metrics == nil {
		return nil, nil, errors.New("k6 summary has no metrics")
	}

	httpReqDuration := metrics["http_req_duration"].(map[string]any)
	p95Any := httpReqDuration["p(95)"]
	p95, ok := p95Any.(float64)
	if !ok {
		// 일부 버전에서 p95 key가 다른 타입일 수 있으니 간단 처리
		return nil, nil, errors.New("missing http_req_duration.p(95)")
	}

	httpReqFailed := metrics["http_req_failed"].(map[string]any)
	valAny := httpReqFailed["value"]
	errRate, ok := valAny.(float64)
	if !ok {
		// value가 없으면 failed/(passes+failed)로 재계산이 필요하지만 minimal이므로 실패
		return nil, nil, errors.New("missing http_req_failed.value")
	}

	httpReqs := metrics["http_reqs"].(map[string]any)
	rateAny := httpReqs["rate"]
	throughput, ok := rateAny.(float64)
	if !ok {
		return nil, nil, errors.New("missing http_reqs.rate")
	}

	selected := map[string]float64{
		"latency_p95_ms": p95,
		"error_rate":     errRate,
		"throughput_rps": throughput,
	}

	latOk := p95 <= sloP95
	errOk := errRate <= sloFailRate

	sloResult := map[string]any{
		"latency_p95_ms": map[string]any{
			"measured": p95,
			"slo_max":  sloP95,
			"ok":       latOk,
		},
		"error_rate": map[string]any{
			"measured": errRate,
			"slo_max":  sloFailRate,
			"ok":       errOk,
		},
	}

	return selected, sloResult, nil
}

func requiredSLO(slo map[string]float64) (float64, float64, error) {
	if slo == nil {
		return 0, 0, errors.New("job.slo is required for k6 in this minimal runner")
	}
	sloP95, ok := slo["latency_p95_ms"]
	if !ok {
		return 0, 0, errors.New("missing slo.latency_p95_ms")
	}
	sloFailRate, ok := slo["error_rate"]
	if !ok {
		return 0, 0, errors.New("missing slo.error_rate")
	}
	return sloP95, sloFailRate, nil
}
