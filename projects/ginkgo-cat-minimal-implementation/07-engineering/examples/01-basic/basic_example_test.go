package main

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestExamples(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Example 01: Basic Ginkgo")
}

var _ = Describe("Basic patterns", func() {
	It("Expect + matcher(Eq)로 값을 단언한다", func() {
		Expect(2 + 2).To(Equal(4))
		Expect(true).To(BeTrue())
		Expect("ok").To(ContainSubstring("o"))
	})
})

