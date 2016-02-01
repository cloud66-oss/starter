package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"testing"

	"github.com/cloud66/starter/common"
)

func AssertFilesHaveSameContent(t *testing.T, expectedFile string, generatedFile string) {
	filename := path.Base(expectedFile)
	expected, err1 := ioutil.ReadFile(expectedFile)
	if err1 != nil {
		t.Errorf("Cannot open/read %s", expectedFile)
	}
	generated, err2 := ioutil.ReadFile(generatedFile)
	if err2 != nil {
		t.Errorf("Cannot open/read %s", generatedFile)
	}

	if filename == "service.yml" {
		// NOTE: starter will be detected as the test projects git repo, so in order
		// for tests to always work we replace the current starter branch (which may
		// change) to 'master' in the generated file.
		generated = regexp.MustCompile(`git_branch: .*`).ReplaceAll(generated, []byte("git_branch: master"))
		generated = regexp.MustCompile(`git_url: .*`).ReplaceAll(generated, []byte("git_url: git@github.com:cloud66/starter.git"))
	}

	if string(expected) != string(generated) {
		t.Errorf("%s generated is wrong\nGenerated:\n%s\nExpected:\n%s\n", filename, string(generated), string(expected))
	}
}

func testApplication(t *testing.T, path string) {
	rootDir := "test/" + path
	var binPath string
	if common.FileExists("./starter") {
		binPath = "./starter"
	} else {
		binPath = "./starter-source"
	}

	command := exec.Command(binPath, "-y", "-p", rootDir+"/src", "-templates", "templates/")
	defer os.Remove(rootDir + "/src/Dockerfile")
	defer os.Remove(rootDir + "/src/service.yml")
	defer os.Remove(rootDir + "/src/docker-compose.yml")
	
	_, err := command.Output()
	if err != nil {
		t.FailNow()
	}
	AssertFilesHaveSameContent(t, rootDir+"/expected/Dockerfile", rootDir+"/src/Dockerfile")
	AssertFilesHaveSameContent(t, rootDir+"/expected/service.yml", rootDir+"/src/service.yml")
	AssertFilesHaveSameContent(t, rootDir+"/expected/docker-compose.yml", rootDir+"/src/docker-compose.yml")

}

func TestRuby13592(t *testing.T) {
	testApplication(t, "ruby/13592")
}

func TestRuby15333(t *testing.T) {
	testApplication(t, "ruby/15333")
}

func TestRuby23080(t *testing.T) {
	testApplication(t, "ruby/23080")
}

func TestRuby25528(t *testing.T) {
	testApplication(t, "ruby/25528")
}

func TestRuby25769(t *testing.T) {
	testApplication(t, "ruby/25769")
}

func init() {
	fmt.Println("Building starter..")
	err := exec.Command("go", "build").Run()
	if err != nil {
		fmt.Println("Failed to build starter")
		fmt.Println(err.Error())
		return
	}
}
