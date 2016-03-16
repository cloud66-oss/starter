package main_test

import (
	"io/ioutil"
	"os"
	"regexp"
	"os/exec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func runStarterWithProject(projectFixture string) (string, error) {
	command := exec.Command(binPath, "-y", "-p", projectFixture+"/src", "-templates", "templates/")
	command_out, err := command.Output()
	output := string(command_out)
	return output, err
}

func cleanupGeneratedFiles(projectFixture string) {
	os.Remove(projectFixture + "/src/Dockerfile")
 	os.Remove(projectFixture + "/src/service.yml")
	os.Remove(projectFixture + "/src/docker-compose.yml")
}

// NOTE: starter will be detected as the test projects git repo, so in order
// for tests to always work we replace the current starter branch (which may
// change) to 'master' in the generated file.
func convertServiceYaml(generated []byte) ([]byte) {
	generated = regexp.MustCompile(`git_branch: .*`).ReplaceAll(generated, []byte("git_branch: master"))
	generated = regexp.MustCompile(`git_url: .*`).ReplaceAll(generated, []byte("git_url: git@github.com:cloud66/starter.git"))
	return generated
}

var _ = Describe("Generating files with Starter", func() {
	Context("using a Rails project with a Mysql database", func() {
		var projectFixturePath string = "test/ruby/rails_mysql"
		
		BeforeEach(func() {
			_, err := runStarterWithProject(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})

    	AfterEach(func() {
			cleanupGeneratedFiles(projectFixturePath)
		})

		It("should generate a Dockerfile", func() {
			dockerfile_expected,_ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
			dockerfile_generated,_ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
			Expect(dockerfile_generated).To(Equal(dockerfile_expected))
		})

		It("should generate a service.yml", func() {
			service_yaml_expected,_ := ioutil.ReadFile(projectFixturePath + "/expected/service.yml")
			service_yaml_generated,_ := ioutil.ReadFile(projectFixturePath + "/src/service.yml")
			service_yaml_generated = convertServiceYaml(service_yaml_generated)
			Expect(service_yaml_generated).To(Equal(service_yaml_expected))
		})
	})
	
	Context("using a Rails project running Unicorn and using a Mysql database", func() {
		var projectFixturePath string = "test/ruby/rails_unicorn_mysql"
		
		BeforeEach(func() {
			_, err := runStarterWithProject(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})

    	AfterEach(func() {
			cleanupGeneratedFiles(projectFixturePath)
		})

		It("should generate a Dockerfile", func() {
			dockerfile_expected,_ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
			dockerfile_generated,_ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
			Expect(dockerfile_generated).To(Equal(dockerfile_expected))
		})

		It("should generate a service.yml", func() {
			service_yaml_expected,_ := ioutil.ReadFile(projectFixturePath + "/expected/service.yml")
			service_yaml_generated,_ := ioutil.ReadFile(projectFixturePath + "/src/service.yml")
			service_yaml_generated = convertServiceYaml(service_yaml_generated)
			Expect(service_yaml_generated).To(Equal(service_yaml_expected))
		})
	})
	
	Context("using a Rails project running Unicorn, some workers and using a Redis and Postgresql database", func() {
		var projectFixturePath string = "test/ruby/rails_jobs_unicorn_redis_postgresql"
		
		BeforeEach(func() {
			_, err := runStarterWithProject(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})

    	AfterEach(func() {
			cleanupGeneratedFiles(projectFixturePath)
		})

		It("should generate a Dockerfile", func() {
			dockerfile_expected,_ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
			dockerfile_generated,_ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
			Expect(dockerfile_generated).To(Equal(dockerfile_expected))
		})

		It("should generate a docker-compose.yml", func() {
			dockercompose_yaml_expected,_ := ioutil.ReadFile(projectFixturePath + "/expected/docker-compose.yml")
			dockercompose_yaml_generated,_ := ioutil.ReadFile(projectFixturePath + "/src/docker-compose.yml")
			Expect(dockercompose_yaml_generated).To(Equal(dockercompose_yaml_expected))
		})

		It("should generate a service.yml", func() {
			service_yaml_expected,_ := ioutil.ReadFile(projectFixturePath + "/expected/service.yml")
			service_yaml_generated,_ := ioutil.ReadFile(projectFixturePath + "/src/service.yml")
			service_yaml_generated = convertServiceYaml(service_yaml_generated)
			Expect(service_yaml_generated).To(Equal(service_yaml_expected))
		})
	})
})
