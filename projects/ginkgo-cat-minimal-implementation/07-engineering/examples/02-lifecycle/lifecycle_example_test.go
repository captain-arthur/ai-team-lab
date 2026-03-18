package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestExamples(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Example 02: Lifecycle")
}

var _ = Describe("Lifecycle patterns", func() {
	var srv *httptest.Server

	BeforeEach(func() {
		srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintln(w, "hello")
		}))
	})

	AfterEach(func() {
		if srv != nil {
			srv.Close()
		}
	})

	It("BeforeEach로 준비하고 AfterEach로 정리한다", func() {
		// srv는 각 It 실행마다 새로 생성된다.
		resp, err := http.Get(srv.URL)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})

	It("정리(AfterEach)가 없으면 서버가 누적될 수 있다(안티패턴 경고)", func() {
		Expect(srv).ToNot(BeNil())
	})
})

