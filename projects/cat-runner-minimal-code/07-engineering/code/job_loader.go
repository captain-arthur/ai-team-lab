package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

func loadJob(jobPath string) (JobSpec, error) {
	b, err := os.ReadFile(jobPath)
	if err != nil {
		return JobSpec{}, err
	}
	var job JobSpec
	if err := yaml.Unmarshal(b, &job); err != nil {
		return JobSpec{}, err
	}
	return job, nil
}
