package main

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestExamples(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Example 04: Table-driven")
}

func bucket(ms int) string {
	switch {
	case ms < 100:
		return "fast"
	case ms < 300:
		return "ok"
	default:
		return "slow"
	}
}

var _ = Describe("DescribeTable pattern", func() {
	DescribeTable("latency bucket 분류", func(ms int, expected string) {
		Expect(bucket(ms)).To(Equal(expected))
	},
		Entry("10ms", 10, "fast"),
		Entry("150ms", 150, "ok"),
		Entry("500ms", 500, "slow"),
	)
})

