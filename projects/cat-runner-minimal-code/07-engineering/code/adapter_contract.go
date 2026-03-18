package main

// Parsed 결과는 “tool-specific raw”를 “CAT 표준 selected_sli + slo_result”로 정규화한 값이다.
type ParsedSLI struct {
	SelectedSLI map[string]float64
	SloResult   map[string]any
}

type RawRef struct {
	Format string
	Path   string
}

type ToolAdapter interface {
	// run: tool을 실제 실행하고, raw 결과 파일을 만들도록 책임진다.
	// return: tool exit code(권위)
	Run(job JobSpec) (exitCode int, err error)

	// locate_raw_result: raw 결과가 저장된 위치를 알려준다.
	LocateRawResult(job JobSpec) (RawRef, error)

	// parse_raw_result: raw 결과에서 selected_sli를 추출한다.
	ParseRawResult(job JobSpec, raw RawRef) (ParsedSLI, error)

	// build_cat_result: parsed + exitCode로 cat-result.json에 들어갈 구조를 만든다.
	BuildCatResult(job JobSpec, parsed ParsedSLI, raw RawRef, exitCode int) CatResult
}
