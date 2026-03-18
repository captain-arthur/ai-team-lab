package main

import "errors"

type CL2Adapter struct{}

func (a CL2Adapter) Run(job JobSpec) (int, error) {
	return 1, errors.New("cl2 adapter stub: not implemented in this minimal runner")
}

func (a CL2Adapter) LocateRawResult(job JobSpec) (RawRef, error) {
	return RawRef{}, errors.New("cl2 adapter stub: locate raw not implemented")
}

func (a CL2Adapter) ParseRawResult(job JobSpec, raw RawRef) (ParsedSLI, error) {
	return ParsedSLI{}, errors.New("cl2 adapter stub: parse raw not implemented")
}

func (a CL2Adapter) BuildCatResult(job JobSpec, parsed ParsedSLI, raw RawRef, exitCode int) CatResult {
	return CatResult{
		TestName:      job.TestName,
		Tool:          "cl2",
		ScenarioType:  job.Scenario.Type,
		SelectedSLI:   parsed.SelectedSLI,
		SloResult:     parsed.SloResult,
		FinalPassFail: "FAIL",
		ExitCode:      exitCode,
		RawResult:     RawResult{Format: raw.Format, Path: raw.Path},
	}
}
