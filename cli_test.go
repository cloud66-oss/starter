package main_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
	"time"
)

var helpText string
var versionText string

var _ = Describe("Running Starter", func() {

	BeforeEach(func() {
		version := "test"
		current_date := time.Now().Format("2006-01-02")

		helpText = fmt.Sprintf("Starter (%s) Help\n", version)
		versionText = fmt.Sprintf("Starter version: %s (%s)\n", version, current_date)

	})

	Context("using the help flag", func() {
		It("should show the help", func() {
			command := exec.Command(binPath, "help")
			command_out, err := command.Output()
			output_string := string(command_out)
			Expect(string(output_string)).To(Equal(helpText))
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Context("using the version flag", func() {
		It("should show the version", func() {
			command := exec.Command(binPath, "version")
			command_out, err := command.Output()
			output_string := string(command_out)
			Expect(string(output_string)).To(Equal(versionText))
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Context("using the -v flag", func() {
		It("should show the version", func() {
			command := exec.Command(binPath, "-v")
			command_out, err := command.Output()
			output_string := string(command_out)
			Expect(string(output_string)).To(Equal(versionText))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
