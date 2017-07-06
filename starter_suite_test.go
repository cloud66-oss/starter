package main_test

import (
	"github.com/cloud66/starter/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
	"testing"
	"time"
)

var binPath string = "./starter"

func TestStarter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Starter Suite")
}

var _ = BeforeSuite(func() {
	current_date := time.Now().Format("2006-01-02")

	err := exec.Command("go", "build", "-ldflags", "-X \"main.VERSION=test\" -X \"main.BUILDDATE="+current_date+"\"").Run()
	Expect(err).NotTo(HaveOccurred())
	Expect(common.FileExists(binPath)).To(BeTrue())

	err = exec.Command(binPath, "-daemon", "-templates", "templates").Start()

	if true {
		timer := time.NewTimer(time.Second * 4)
		<-timer.C
	}

	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := exec.Command("rm", binPath).Run()
	Expect(err).NotTo(HaveOccurred())
})
