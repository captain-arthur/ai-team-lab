package main

import (
	"encoding/json"
	"os"
)

func writeCatResult(outPath string, cat CatResult) error {
	b, err := json.MarshalIndent(cat, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, b, 0o644)
}
