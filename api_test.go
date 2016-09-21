package main

import (
	"encoding/json"
	"github.com/go-resty/resty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
			Expect(string(resp.Body())).To(Equal(`"test"`))
		})
	})
	Context("get the list of files starter is using to analyse the codebase", func() {
		It("should respond with all the supported files", func() {
			resp, err := resty.R().Get("http://127.0.0.1:9090/analyze/supported")
			Expect(err).NotTo(HaveOccurred())
			Expect(string(resp.Body())).To(Equal(`{"Languages":[{"Name":"ruby","Files":["Gemfile","Procfile","config/database.yml"],"SupportedVersion":["1.9.3","1.9","1","2.0.0","2.0","2.1.1","2.1.10","2.1.2","2.1.3","2.1.4","2.1.5","2.1.6","2.1.7","2.1.8","2.1.9","2.1","2.2.0","2.2.1","2.2.2","2.2.3","2.2.4","2.2.5","2.2","2.3.0","2.3.1","2.3","2","alpine","latest","onbuild","slim","wheezy"]},{"Name":"node","Files":["package.json","Procfile"],"SupportedVersion":["0.10.28","0.10.30","0.10.31","0.10.32","0.10.33","0.10.34","0.10.35","0.10.36","0.10.37","0.10.38","0.10.39","0.10.40","0.10.41","0.10.42","0.10.43","0.10.44","0.10.45","0.10.46","0.10","0.11.13","0.11.14","0.11.15","0.11.16","0.11","0.12.0","0.12.1","0.12.10","0.12.11","0.12.12","0.12.13","0.12.14","0.12.15","0.12.2","0.12.3","0.12.4","0.12.5","0.12.6","0.12.7","0.12.8","0.12.9","0.12","0.8.28","0.8","0","4.0.0","4.0","4.1.0","4.1.1","4.1.2","4.1","4.2.0","4.2.1","4.2.2","4.2.3","4.2.4","4.2.5","4.2.6","4.2","4.3.0","4.3.1","4.3.2","4.3","4.4.0","4.4.1","4.4.2","4.4.3","4.4.4","4.4.5","4.4.6","4.4.7","4.4","4.5.0","4.5","4","5.0.0","5.0","5.1.0","5.1.1","5.1","5.10.0","5.10.1","5.10","5.11.0","5.11.1","5.11","5.12.0","5.12","5.2.0","5.2","5.3.0","5.3","5.4.0","5.4.1","5.4","5.5.0","5.5","5.6.0","5.6","5.7.0","5.7.1","5.7","5.8.0","5.8","5.9.0","5.9.1","5.9","5","6.0.0","6.0","6.1.0","6.1","6.2.0","6.2.1","6.2.2","6.2","6.3.0","6.3.1","6.3","6.4.0","6.4","6.5.0","6.5","6.6.0","6.6","6","argon","latest","onbuild","slim","wheezy"]},{"Name":"php","Files":["composer.json"],"SupportedVersion":["5.3.29","5.3","5.4.33","5.4.34","5.4.35","5.4.36","5.4.37","5.4.38","5.4.39","5.4.40","5.4.41","5.4.42","5.4.43","5.4.44","5.4.45","5.4","5.5.17","5.5.18","5.5.19","5.5.20","5.5.21","5.5.22","5.5.23","5.5.24","5.5.25","5.5.26","5.5.27","5.5.28","5.5.29","5.5.30","5.5.31","5.5.32","5.5.33","5.5.34","5.5.35","5.5.36","5.5.37","5.5.38","5.5","5.6.1","5.6.10","5.6.11","5.6.12","5.6.13","5.6.14","5.6.15","5.6.16","5.6.17","5.6.18","5.6.19","5.6.2","5.6.20","5.6.21","5.6.22","5.6.23","5.6.24","5.6.25","5.6.26","5.6.3","5.6.4","5.6.5","5.6.6","5.6.7","5.6.8","5.6.9","5.6","5","7.0.0","7.0.0RC1","7.0.0RC2","7.0.0RC3","7.0.0RC4","7.0.0RC5","7.0.0RC6","7.0.0RC7","7.0.0RC8","7.0.0beta1","7.0.0beta2","7.0.0beta3","7.0.1","7.0.10","7.0.11","7.0.2","7.0.3","7.0.4","7.0.5","7.0.6","7.0.7","7.0.8","7.0.9","7.0","7.1.0RC1","7.1.0RC2","7.1","7","alpine","apache","cli","fpm","latest","zts"]}]}`))
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
			analysis.BuildCommands = []string {}
			analysis.DeployCommands = []string {}
			analysis.Databases = []string {}


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
			analysis.BuildCommands = []string {}
			analysis.DeployCommands = []string {}
			analysis.Databases = []string {}


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
					"git_repo":  "fake.git",
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
			analysis.BuildCommands = []string {}
			analysis.DeployCommands = []string {}
			analysis.Databases = []string {}

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
			analysis.BuildCommands = []string {}
			analysis.DeployCommands = []string {}
			analysis.Databases = []string {}

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
			analysis.BuildCommands = []string {}
			analysis.DeployCommands = []string {}
			analysis.Databases = []string {}

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
			analysis.LanguageVersion = "4.2.0"
			analysis.Framework = "express"
			analysis.FrameworkVersion = "4.13.0"
			analysis.SupportedLanguageVersions = []string {"0.10.28","0.10.30","0.10.31","0.10.32","0.10.33","0.10.34","0.10.35","0.10.36","0.10.37","0.10.38","0.10.39","0.10.40","0.10.41","0.10.42","0.10.43","0.10.44","0.10.45","0.10.46","0.10","0.11.13","0.11.14","0.11.15","0.11.16","0.11","0.12.0","0.12.1","0.12.10","0.12.11","0.12.12","0.12.13","0.12.14","0.12.15","0.12.2","0.12.3","0.12.4","0.12.5","0.12.6","0.12.7","0.12.8","0.12.9","0.12","0.8.28","0.8","0","4.0.0","4.0","4.1.0","4.1.1","4.1.2","4.1","4.2.0","4.2.1","4.2.2","4.2.3","4.2.4","4.2.5","4.2.6","4.2","4.3.0","4.3.1","4.3.2","4.3","4.4.0","4.4.1","4.4.2","4.4.3","4.4.4","4.4.5","4.4.6","4.4.7","4.4","4.5.0","4.5","4","5.0.0","5.0","5.1.0","5.1.1","5.1","5.10.0","5.10.1","5.10","5.11.0","5.11.1","5.11","5.12.0","5.12","5.2.0","5.2","5.3.0","5.3","5.4.0","5.4.1","5.4","5.5.0","5.5","5.6.0","5.6","5.7.0","5.7.1","5.7","5.8.0","5.8","5.9.0","5.9.1","5.9","5","6.0.0","6.0","6.1.0","6.1","6.2.0","6.2.1","6.2.2","6.2","6.3.0","6.3.1","6.3","6.4.0","6.4","6.5.0","6.5","6.6.0","6.6","6","argon","latest","onbuild","slim","wheezy"}
			analysis.Warnings = nil
			analysis.Dockerfile = string(file)
			analysis.StartCommands = []string {"node server.js"}
			analysis.BuildCommands = []string {}
			analysis.DeployCommands = []string {}
			analysis.Databases = []string {}


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
			analysis.LanguageVersion = "0.10.28"
			analysis.Framework = "express"
			analysis.FrameworkVersion = "4.14.0"
			analysis.SupportedLanguageVersions = []string {"0.10.28","0.10.30","0.10.31","0.10.32","0.10.33","0.10.34","0.10.35","0.10.36","0.10.37","0.10.38","0.10.39","0.10.40","0.10.41","0.10.42","0.10.43","0.10.44","0.10.45","0.10.46","0.10","0.11.13","0.11.14","0.11.15","0.11.16","0.11","0.12.0","0.12.1","0.12.10","0.12.11","0.12.12","0.12.13","0.12.14","0.12.15","0.12.2","0.12.3","0.12.4","0.12.5","0.12.6","0.12.7","0.12.8","0.12.9","0.12","0.8.28","0.8","0","4.0.0","4.0","4.1.0","4.1.1","4.1.2","4.1","4.2.0","4.2.1","4.2.2","4.2.3","4.2.4","4.2.5","4.2.6","4.2","4.3.0","4.3.1","4.3.2","4.3","4.4.0","4.4.1","4.4.2","4.4.3","4.4.4","4.4.5","4.4.6","4.4.7","4.4","4.5.0","4.5","4","5.0.0","5.0","5.1.0","5.1.1","5.1","5.10.0","5.10.1","5.10","5.11.0","5.11.1","5.11","5.12.0","5.12","5.2.0","5.2","5.3.0","5.3","5.4.0","5.4.1","5.4","5.5.0","5.5","5.6.0","5.6","5.7.0","5.7.1","5.7","5.8.0","5.8","5.9.0","5.9.1","5.9","5","6.0.0","6.0","6.1.0","6.1","6.2.0","6.2.1","6.2.2","6.2","6.3.0","6.3.1","6.3","6.4.0","6.4","6.5.0","6.5","6.6.0","6.6","6","argon","latest","onbuild","slim","wheezy"}
			analysis.Warnings = []string {"No Procfile was detected. It is strongly advised to add one in order to specify the commands to run your services."}
			analysis.Dockerfile = string(file)
			analysis.StartCommands = []string {"npm start"}
			analysis.BuildCommands = []string {}
			analysis.DeployCommands = []string {}
			analysis.Databases = []string {"mysql", "postgresql"}


			b, err := json.Marshal(analysis)
			Expect(string(resp.Body())).To(Equal(string(b)))
			os.Remove(path + "/Dockerfile")
		})
	})

})
