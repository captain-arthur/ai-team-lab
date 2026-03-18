package main

import (
	"sync/atomic"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestExamples(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Example 03: Eventually/Consistently")
}

var _ = Describe("Async state verification", func() {
	It("Eventually는 '언젠가 참이 될 것'을 기다린다", func() {
		var ready atomic.Bool

		// 200ms 뒤 ready=true
		go func() {
			time.Sleep(200 * time.Millisecond)
			ready.Store(true)
		}()

		Eventually(func() bool {
			return ready.Load()
		}, 2*time.Second, 20*time.Millisecond).Should(BeTrue())
	})

	It("Consistently는 '일정 시간 동안 계속 참'을 확인한다", func() {
		var stable atomic.Bool
		stable.Store(true)

		Consistently(func() bool {
			return stable.Load()
		}, 300*time.Millisecond, 30*time.Millisecond).Should(BeTrue())
	})
})

