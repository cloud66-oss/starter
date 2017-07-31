package main

import (
	"encoding/json"
	"github.com/go-resty/resty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Describe("Running Starter in daemon mode", func() {

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
			Expect(string(resp.Body())).To(Equal(`"test"`))
		})
	})
	Context("get the list of files starter is using to analyse the codebase", func() {
		It("should respond with all the supported files", func() {
			resp, err := resty.R().Get("http://127.0.0.1:9090/analyze/supported")
			Expect(err).NotTo(HaveOccurred())
			Expect(string(resp.Body())).To(Equal(`{"Languages":[{"Name":"docker-compose","Files":["docker-compose.yml"],"SupportedVersion":null},{"Name":"ruby","Files":["Gemfile","Procfile","config/database.yml"],"SupportedVersion":null},{"Name":"node","Files":["package.json","Procfile",".meteor/release"],"SupportedVersion":["4.6"]},{"Name":"php","Files":["composer.json"],"SupportedVersion":null},{"Name":"service.yml","Files":["service.yml"],"SupportedVersion":null}]}`))
		})
	})

	Context("get the list of base Dockerfiles starter is supporting", func() {
		It("should respond with all the supported Dockerfiles", func() {
			resp, err := resty.R().Get("http://127.0.0.1:9090/analyze/dockerfiles")
			Expect(err).NotTo(HaveOccurred())
			dockerfiles := []Dockerfile{}

			path := "test/ruby/Dockerfile.base"
			rubyDockerFile, err := ioutil.ReadFile(path)
			dockerfile := Dockerfile{}
			dockerfile.Language = "ruby"
			dockerfile.Base = string(rubyDockerFile)
			dockerfiles = append(dockerfiles, dockerfile)

			path = "test/node/Dockerfile.base"
			nodeDockerFile, err := ioutil.ReadFile(path)
			dockerfile = Dockerfile{}
			dockerfile.Language = "node"
			dockerfile.Base = string(nodeDockerFile)
			dockerfiles = append(dockerfiles, dockerfile)

			path = "test/php/Dockerfile.base"
			phpDockerFile, err := ioutil.ReadFile(path)
			dockerfile = Dockerfile{}
			dockerfile.Language = "php"
			dockerfile.Base = string(phpDockerFile)
			dockerfiles = append(dockerfiles, dockerfile)

			b, err := json.Marshal(dockerfiles)

			Expect(err).NotTo(HaveOccurred())
			Expect(string(resp.Body())).To(Equal(string(b)))
		})
	})

	Context("analyse a ruby project", func() {
		It("should respond with program langange ruby", func() {
			path := "test/ruby/rails_mysql/src"
			resp, err := resty.R().SetBody(`{"path":"` + path + `", "generate":"dockerfile "}`).Post("http://127.0.0.1:9090/analyze")
			Expect(err).NotTo(HaveOccurred())
			dockerfile, err := ioutil.ReadFile(path + "/Dockerfile")
			Expect(err).NotTo(HaveOccurred())

			analysis := analysisResult{}
			analysis.Language = "ruby"
			analysis.Framework = "rails"
			analysis.Ok = true
			analysis.Warnings = nil
			analysis.Dockerfile = string(dockerfile)
			analysis.StartCommands = nil
			analysis.BuildCommands = []string{}
			analysis.DeployCommands = []string{}
			analysis.Databases = []string{}

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

			analysis := analysisResult{}
			analysis.Ok = true
			analysis.Warnings = nil
			analysis.Language = "ruby"
			analysis.Framework = "rails"
			analysis.Dockerfile = string(dockerfile)
			analysis.Service = string(serviceyml)
			analysis.DockerCompose = string(dockercomposeyml)
			analysis.BuildCommands = []string{}
			analysis.DeployCommands = []string{}
			analysis.Databases = []string{}

			b, err := json.Marshal(analysis)
			Expect(string(resp.Body())).To(Equal(string(b)))
			os.Remove(path + "/Dockerfile")
			os.Remove(path + "/service.yml")
			os.Remove(path + "/docker-compose.yml")

		})
	})

	Context("analyse a ruby project through upload files and request a dockerfile, docker-compose.yml and service.yml and s", func() {
		It("should respond with a dockerfile, docker-compose.yml and service.yml", func() {
			path := "test/ruby/rails_mysql/src"
			expected := "test/ruby/rails_mysql/expected/api"

			resp, err := resty.R().
				SetFile("source", path+"/source.zip").
				SetFormData(map[string]string{
					"git_repo":   "fake.git",
					"git_branch": "master",
				}).Post("http://127.0.0.1:9090/analyze/upload")
			Expect(err).NotTo(HaveOccurred())
			dockerfile, err := ioutil.ReadFile(expected + "/Dockerfile")
			Expect(err).NotTo(HaveOccurred())
			dockercomposeyml, err := ioutil.ReadFile(expected + "/docker-compose.yml")
			Expect(err).NotTo(HaveOccurred())
			serviceyml, err := ioutil.ReadFile(expected + "/service.yml")
			Expect(err).NotTo(HaveOccurred())

			analysis := analysisResult{}
			analysis.Ok = true
			analysis.Warnings = nil
			analysis.Language = "ruby"
			analysis.Framework = "rails"
			analysis.Dockerfile = string(dockerfile)
			analysis.DockerCompose = string(dockercomposeyml)
			analysis.Service = string(serviceyml)
			analysis.BuildCommands = []string{}
			analysis.DeployCommands = []string{}
			analysis.Databases = []string{}

			b, err := json.Marshal(analysis)
			Expect(string(resp.Body())).To(Equal(string(b)))

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

			analysis := analysisResult{}
			analysis.Ok = true
			analysis.Language = "ruby"
			analysis.Framework = "rails"
			analysis.Warnings = nil
			analysis.Dockerfile = string(dockerfile)
			analysis.Service = string(serviceyml)
			analysis.BuildCommands = []string{}
			analysis.DeployCommands = []string{}
			analysis.Databases = []string{}

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
			analysis := analysisResult{}
			analysis.Ok = true
			analysis.Language = "ruby"
			analysis.Framework = "rails"
			analysis.Warnings = nil
			analysis.Dockerfile = string(file)
			analysis.BuildCommands = []string{}
			analysis.DeployCommands = []string{}
			analysis.Databases = []string{}

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
			analysis := analysisResult{}
			analysis.Ok = true
			analysis.Language = "node"
			analysis.LanguageVersion = "4.6"
			analysis.Framework = "express"
			analysis.FrameworkVersion = "4.13.0"
			analysis.SupportedLanguageVersions = []string{"4.6"}
			analysis.Warnings = nil
			analysis.Dockerfile = string(file)
			analysis.StartCommands = []string{"node server.js"}
			analysis.BuildCommands = []string{}
			analysis.DeployCommands = []string{}
			analysis.Databases = []string{}

			b, err := json.Marshal(analysis)
			Expect(string(resp.Body())).To(Equal(string(b)))
			os.Remove(path + "/Dockerfile")
		})
	})
	Context("analyse a node project with databaze only request a dockerfile", func() {
		It("should respond with a dockerfile", func() {
			path := "test/node/express_mysql_pg/src"
			resp, err := resty.R().SetBody(`{"path":"` + path + `", "generate":"dockerfile"}`).Post("http://127.0.0.1:9090/analyze")
			Expect(err).NotTo(HaveOccurred())
			file, err := ioutil.ReadFile(path + "/Dockerfile")
			Expect(err).NotTo(HaveOccurred())
			analysis := analysisResult{}
			analysis.Ok = true
			analysis.Language = "node"
			analysis.LanguageVersion = "4.6"
			analysis.Framework = "express"
			analysis.FrameworkVersion = "4.14.0"
			analysis.SupportedLanguageVersions = []string{"4.6"}
			analysis.Warnings = []string{"No Procfile was detected. It is strongly advised to add one in order to specify the commands to run your services."}
			analysis.Dockerfile = string(file)
			analysis.StartCommands = []string{"npm start"}
			analysis.BuildCommands = []string{}
			analysis.DeployCommands = []string{}
			analysis.Databases = []string{"mysql", "postgresql"}

			b, err := json.Marshal(analysis)
			Expect(string(resp.Body())).To(Equal(string(b)))
			os.Remove(path + "/Dockerfile")
		})
	})

})
