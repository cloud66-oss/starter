package main_test

import (
	"os/exec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


var helpText string = "Starter (1.0.2) Help\n"

var versionText string = "Starter version: 1.0.2 (2016-03-14)\n"


var _ = Describe("Running Starter", func() {
	Context("using the h flag", func() {
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
})
