package main_test

import (
	"os/exec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


var helpText string = 
`Cloud 66 Starter (dev) Help
Copyright 2016 Cloud66 Inc.
`

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
})
