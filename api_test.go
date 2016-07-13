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
	Context("analyse a ruby project", func() {
		It("should respond with program langange ruby", func() {
			path := "test/ruby/rails_mysql/src"
			resp, err := resty.R().SetBody(`{"path":"` + path + `", "generate":"dockerfile "}`).Post("http://127.0.0.1:9090/analyze")
			Expect(err).NotTo(HaveOccurred())
			dockerfile, err := ioutil.ReadFile(path + "/Dockerfile")
			Expect(err).NotTo(HaveOccurred())

			analysis := CodebaseAnalysis{}
			analysis.Language = "ruby"
    		analysis.Ok = true
			analysis.Warnings = nil
			analysis.Dockerfile = string(dockerfile)

			b, err := json.Marshal(analysis)
		    Expect(string(resp.Body())).To(Equal(string(b)))

			os.Remove(path + "/Dockerfile")
		})
	})




	Context("analyse a ruby project and request a dockerfile, docker-compose.yml and service.yml", func() {
		It("should respond with a dockerfile, docker-compose.yml and service.yml", func() {
			path := "test/ruby/rails_mysql/src"
			resp, err := resty.R().SetBody(`{"path":"` + path + `", "generate":"dockerfile,service,docker-compose "}`).Post("http://127.0.0.1:9090/analyze")
			Expect(err).NotTo(HaveOccurred())
			dockerfile, err := ioutil.ReadFile(path + "/Dockerfile")
			Expect(err).NotTo(HaveOccurred())
			serviceyml, err := ioutil.ReadFile(path + "/service.yml")
			Expect(err).NotTo(HaveOccurred())
			dockercomposeyml, err := ioutil.ReadFile(path + "/docker-compose.yml")
			Expect(err).NotTo(HaveOccurred())
			
			analysis := CodebaseAnalysis{}
    		analysis.Ok = true
			analysis.Warnings = nil
			analysis.Language = "ruby"
			analysis.Dockerfile = string(dockerfile)
			analysis.Service = string(serviceyml)
			analysis.DockerCompose = string(dockercomposeyml)
			b, err := json.Marshal(analysis)
		    Expect(string(resp.Body())).To(Equal(string(b)))
		    os.Remove(path + "/Dockerfile")
		    os.Remove(path + "/service.yml")
		  	os.Remove(path + "/docker-compose.yml")
		      
		})
	})


	Context("analyse a ruby project and request a dockerfile and service.yml", func() {
		It("should respond with a dockerfile and service.yml", func() {
			path := "test/ruby/rails_mysql/src"
			resp, err := resty.R().SetBody(`{"path":"` + path + `", "generate":"dockerfile,service"}`).Post("http://127.0.0.1:9090/analyze")
			Expect(err).NotTo(HaveOccurred())
			dockerfile, err := ioutil.ReadFile(path + "/Dockerfile")
			Expect(err).NotTo(HaveOccurred())
			serviceyml, err := ioutil.ReadFile(path + "/service.yml")
			Expect(err).NotTo(HaveOccurred())
			
			analysis := CodebaseAnalysis{}
    		analysis.Ok = true
    		analysis.Language = "ruby"
			analysis.Warnings = nil
			analysis.Dockerfile = string(dockerfile)
			analysis.Service = string(serviceyml)
			b, err := json.Marshal(analysis)
		    Expect(string(resp.Body())).To(Equal(string(b)))
		    os.Remove(path + "/Dockerfile")
		    os.Remove(path + "/service.yml")
		})
	})


	Context("analyse a ruby project and only request a dockerfile", func() {
		It("should respond with a dockerfile", func() {
			path := "test/ruby/rails_mysql/src"
			resp, err := resty.R().SetBody(`{"path":"` + path + `", "generate":"dockerfile"}`).Post("http://127.0.0.1:9090/analyze")
			Expect(err).NotTo(HaveOccurred())
			file, err := ioutil.ReadFile(path + "/Dockerfile")
			Expect(err).NotTo(HaveOccurred())
			analysis := CodebaseAnalysis{}
    		analysis.Ok = true
    		analysis.Language = "ruby"
			analysis.Warnings = nil
			analysis.Dockerfile = string(file)
			b, err := json.Marshal(analysis)
		    Expect(string(resp.Body())).To(Equal(string(b)))
		    os.Remove(path + "/Dockerfile")
		})
	})
	Context("analyse a node project only request a dockerfile", func() {
		It("should respond with a dockerfile", func() {
			path := "test/node/express_procfile/src"
			resp, err := resty.R().SetBody(`{"path":"` + path + `", "generate":"dockerfile"}`).Post("http://127.0.0.1:9090/analyze")
			Expect(err).NotTo(HaveOccurred())
			file, err := ioutil.ReadFile(path + "/Dockerfile")
			Expect(err).NotTo(HaveOccurred())
			analysis := CodebaseAnalysis{}
    		analysis.Ok = true
    		analysis.Language = "node"
			analysis.Warnings = nil
			analysis.Dockerfile = string(file)
			b, err := json.Marshal(analysis)
		    Expect(string(resp.Body())).To(Equal(string(b)))
		    os.Remove(path + "/Dockerfile")
		})
	})
	
})
