package main

import "time"

type CatResult struct {
	TestName       string             `json:"test_name"`
	Tool           string             `json:"tool"`
	ScenarioType   string             `json:"scenario_type"`
	SelectedSLI    SelectedSLI        `json:"selected_sli"`
	SloResult      SloResult          `json:"slo_result"`
	FinalPassFail  string             `json:"final_pass_fail"`
	ExitCode       int                `json:"exit_code"`
	Timestamp      string             `json:"timestamp"`
	ScenarioParams ScenarioParams    `json:"scenario_params,omitempty"`
}

type ScenarioParams struct {
	DelayMS       int `json:"delay_ms"`
	FailEvery     int `json:"fail_every"`
	Requests      int `json:"requests"`
}

type SelectedSLI struct {
	AvgLatencyMs float64 `json:"avg_latency_ms"`
	ErrorRate    float64 `json:"error_rate"`
}

type SloResult struct {
	LatencyOk bool    `json:"latency_ok"`
	ErrorOk   bool    `json:"error_ok"`
	LatencySloMs float64 `json:"latency_slo_ms"`
	ErrorSloRate float64 `json:"error_slo_rate"`
}

func nowISO() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

