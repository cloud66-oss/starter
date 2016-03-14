package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"github.com/cloud66/starter/common"
)

var binPath string = "./starter"

func TestStarter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Starter Suite")
}

var _ = BeforeSuite(func() {
	Expect(common.FileExists(binPath)).To(BeTrue())
})

var _ = AfterSuite(func() {
})

