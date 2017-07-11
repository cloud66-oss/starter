package service_yml

import (
	"fmt"
	"os"
)


func handleEnvVars(file []byte) string{
	finalFormat := ""

	return finalFormat
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
