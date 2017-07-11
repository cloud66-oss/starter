package service_yml

import (
	"fmt"
	"os"
	"strings"
)

func handleEnvVarsFormat(file []byte) string {
	finalFormat := ""

	lines := strings.Split(string(file), "\n")

	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "_env") {
			for j := 0; j < len(lines[i])-4; j++ {
				if lines[i][j] == '_' && lines[i][j+1] == 'e' && lines[i][j+2] == 'n' && lines[i][j+3] == 'v' {
					lines[i] = lines[i][:j] + "$" + lines[i][j+4:]
				}
			}
		}
		finalFormat = finalFormat + lines[i] + "\n"
	}

	return finalFormat
}

func handleVolumes(serviceVolumes []string) []VolumeMounts {
	var kubeVolumes []VolumeMounts

	for _, volume := range serviceVolumes {
		name := ""
		mountPath := ""
		var i int
		var readOnly bool
		if volume[0] == '"' {
			i = 1
		} else {
			i = 0
		}
		for ; i < len(volume); i++ {
			if volume[i] == ':' {
				break
			} else {
				name = string(append([]byte(name), volume[i]))
			}
		}

		for i = i + 1; i < len(volume); i++ {
			if volume[i] == ':' || volume[i] == '"' || volume[i] == '\n' {
				break
			} else {
				mountPath = string(append([]byte(mountPath), volume[i]))
			}
		}
		if i < len(volume)-2 {
			if volume[i] == ':' && volume[i+1] == 'r' && volume[i+2] == 'o' {
				readOnly = true
			}
		}
		kubeVolume := VolumeMounts{
			Name:      name,
			MountPath: mountPath,
			ReadOnly:  readOnly,
		}
		kubeVolumes = append(kubeVolumes, kubeVolume)
	}

	return kubeVolumes
}

func getKeysValues(env_vars map[string]string) ([]interface{}, []interface{}) {
	keys := []interface{}{}
	values := []interface{}{}
	for k, v := range env_vars {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
