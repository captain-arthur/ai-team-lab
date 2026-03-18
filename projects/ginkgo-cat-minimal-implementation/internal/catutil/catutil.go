package catutil

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type CatResult struct {
	TestName       string                 `json:"test_name"`
	Tool           string                 `json:"tool"`
	ScenarioType   string                 `json:"scenario_type"`
	Target         string                 `json:"target,omitempty"`
	SelectedSLI    map[string]float64    `json:"selected_sli"`
	SloResult      map[string]any        `json:"slo_result"`
	FinalPassFail  string                 `json:"final_pass_fail"`
	ExitCode       int                    `json:"exit_code"`
	Timestamp      string                 `json:"timestamp"`
	ScenarioParams map[string]any        `json:"scenario_params,omitempty"`
}

func NowISO() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0o755)
}

func WriteCatResult(path string, r CatResult) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// PercentileMs는 time.Duration slice에서 p백분위(ms)를 계산한다.
// p는 0~100 범위(예: 95)로 받는다.
func PercentileMs(durs []time.Duration, p float64) float64 {
	if len(durs) == 0 {
		return 0
	}
	cp := make([]time.Duration, len(durs))
	copy(cp, durs)
	sort.Slice(cp, func(i, j int) bool { return cp[i] < cp[j] })

	if p <= 0 {
		return float64(cp[0].Milliseconds())
	}
	if p >= 100 {
		return float64(cp[len(cp)-1].Milliseconds())
	}
	// nearest-rank 스타일(간단/재현성 우선)
	rank := int((p / 100.0) * float64(len(cp)-1))
	return float64(cp[rank].Milliseconds())
}

func ErrorRate(total int, errors int) float64 {
	if total <= 0 {
		return 0
	}
	return float64(errors) / float64(total)
}

func Must(err error) {
	if err != nil {
		panic(fmt.Errorf("catutil: %w", err))
	}
}

