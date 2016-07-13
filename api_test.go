package main

import (
	. "github.com/onsi/ginkgo"
	"github.com/go-resty/resty"
	. "github.com/onsi/gomega"
	"encoding/json"
	"io/ioutil"
	"os"
)

var _ = Describe("Running Starter in damon mode", func() {


	Context("ping the service", func() {
		It("should respond with ok", func() {
			resp, err := resty.R().Get("http://127.0.0.1:9090/ping")
			Expect(err).NotTo(HaveOccurred())
			Expect(string(resp.Body())).To(Equal(`"ok"`))
		})
	})
	Context("get the version", func() {
		It("should respond with version number", func() {
			resp, err := resty.R().Get("http://127.0.0.1:9090/version")
			Expect(err).NotTo(HaveOccurred())
			Expect(string(resp.Body())).To(Equal(`"1.0.2"`))
		})
	})

	FContext("analayse a ruby project", func() {

		FIt("should respond with a dockerfile", func() {
			path := "test/ruby/rails_mysql/src"
			resp, err := resty.R().SetBody(`{"path":"` + path + `", "generate":"dockerfile"}`).Post("http://127.0.0.1:9090/analyze")
			Expect(err).NotTo(HaveOccurred())
			file, _ := ioutil.ReadFile(path + "/Dockerfile")
			
			analysis := CodebaseAnalysis{}
    		analysis.Ok = true
			analysis.Warnings = nil
			analysis.Dockerfile = string(file)
			b, err := json.Marshal(analysis)
		    Expect(string(resp.Body())).To(Equal(string(b)))
		    os.Remove(path + "/Dockerfile")
		})
	})
	
})
