package main_test

import (
	//"fmt"
	//	"bytes"
	//"time"
	"os/exec"
	. "github.com/onsi/ginkgo"
	"github.com/go-resty/resty"
	. "github.com/onsi/gomega"
)


var start_in_deamon_cmd *exec.Cmd

func runStarterInDaemonMode() () {
	start_in_deamon_cmd = exec.Command(binPath, "-daemon", "&")
	start_in_deamon_cmd.Run()

}

func stopStarterInDaamonMode() {
	start_in_deamon_cmd.Process.Kill()
}


var _ = Describe("Running Starter in damon mode", func() {

	BeforeEach(func() {
		runStarterInDaemonMode()
	})

	AfterEach(func() {
		stopStarterInDaamonMode()
	})
	Context("ping the service", func() {
		It("should respond with ok", func() {
			resp, err := resty.R().Get("http://127.0.0.1:9090/ping")
			Expect(err).NotTo(HaveOccurred())
			Expect(string(resp.Body())).To(Equal("\"ok\""))
		})
	})
	Context("get the version", func() {
		It("should respond with version number", func() {
			resp, err := resty.R().Get("http://127.0.0.1:9090/version")
			Expect(err).NotTo(HaveOccurred())
			Expect(string(resp.Body())).To(Equal("\"1.0.2\""))
		})
	})
	
})
