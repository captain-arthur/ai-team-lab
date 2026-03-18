package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type JobSpec struct {
	TestName    string             `yaml:"test_name"`
	Tool        string             `yaml:"tool"`
	Entry       string             `yaml:"entry"`
	Scenario    ScenarioSpec       `yaml:"scenario"`
	SLO         map[string]float64 `yaml:"slo"`
	Output      OutputSpec         `yaml:"output"`
	SelectedSLI []string           `yaml:"selected_sli"`
}

type ScenarioSpec struct {
	Type      string         `yaml:"type"`
	Target    string         `yaml:"target,omitempty"`
	LoadModel string         `yaml:"load_model,omitempty"`
	RPS       int            `yaml:"rps,omitempty"`
	VUS       int            `yaml:"vus,omitempty"`
	Duration  string         `yaml:"duration,omitempty"`
	Config    map[string]any `yaml:"config,omitempty"`
}

type OutputSpec struct {
	Dir string `yaml:"dir"`
}

type CatResult struct {
	TestName      string             `json:"test_name"`
	Tool          string             `json:"tool"`
	ScenarioType  string             `json:"scenario_type"`
	SelectedSLI   map[string]float64 `json:"selected_sli"`
	SloResult     map[string]any     `json:"slo_result"`
	FinalPassFail string             `json:"final_pass_fail"`
	ExitCode      int                `json:"exit_code"`
	RawResult     RawResult          `json:"raw_result"`
	Timestamp     string             `json:"timestamp"`
}

type RawResult struct {
	Format string `json:"format"`
	Path   string `json:"path"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: cat-runner <job.yaml>")
		os.Exit(2)
	}
	jobPath := os.Args[1]

	job, err := loadJob(jobPath)
	must(err)

	// 상대 경로는 job.yaml 기준으로 해석한다.
	baseDir := filepath.Dir(jobPath)
	if job.Entry != "" && !filepath.IsAbs(job.Entry) {
		job.Entry = filepath.Join(baseDir, job.Entry)
	}
	if job.Output.Dir != "" && !filepath.IsAbs(job.Output.Dir) {
		job.Output.Dir = filepath.Join(baseDir, job.Output.Dir)
	}

	tool := job.Tool
	if tool == "" {
		must(errors.New("job.tool is required"))
	}
	if job.Output.Dir == "" {
		must(errors.New("job.output.dir is required"))
	}
	if err := os.MkdirAll(job.Output.Dir, 0o755); err != nil {
		must(err)
	}

	var adapter ToolAdapter
	switch tool {
	case "k6":
		adapter = K6Adapter{}
	case "ginkgo":
		adapter = GinkgoAdapter{}
	case "cl2":
		adapter = CL2Adapter{}
	default:
		must(fmt.Errorf("unknown tool: %s", tool))
	}

	exitCode, runErr := adapter.Run(job)
	must(runErr)

	raw, locErr := adapter.LocateRawResult(job)
	must(locErr)

	parsed, parseErr := adapter.ParseRawResult(job, raw)
	// parse 실패해도 cat-result 구조는 남길 수 있게, minimal runner에서는 Build에서 증거를 기록한다.
	if parseErr != nil {
		parsed = ParsedSLI{
			SelectedSLI: map[string]float64{},
			SloResult: map[string]any{
				"parse_error": parseErr.Error(),
			},
		}
	}

	cat := adapter.BuildCatResult(job, parsed, raw, exitCode)
	cat.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)

	outPath := filepath.Join(job.Output.Dir, "cat-result.json")
	must(writeCatResult(outPath, cat))

	// runner exit code 정책: tool exit code를 그대로 반영(권위 단일화)
	if cat.ExitCode == 0 {
		os.Exit(0)
	}
	os.Exit(1)
}
