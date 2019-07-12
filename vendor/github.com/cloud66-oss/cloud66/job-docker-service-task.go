package cloud66

import "encoding/json"

type dockerServiceTaskJob struct {
	Task        string `json:"task"`
	ServiceName string `json:"service_name"`
	PrivateIp   string `json:"private_ip"`
}

type DockerServiceTaskJob struct {
	*BasicJob
	*dockerServiceTaskJob
}

func (job *DockerServiceTaskJob) UnmarshalJSON(b []byte) error {
	var bj BasicJob
	if err := json.Unmarshal(b, &bj); err == nil {
		var j dockerServiceTaskJob
		if err := json.Unmarshal(bj.basicJob.ParamsRaw, &j); err != nil {
			return err
		}
		*job = DockerServiceTaskJob{BasicJob: &bj, dockerServiceTaskJob: &j}
		return nil
	} else {
		return err
	}
}

func (job DockerServiceTaskJob) GetBasicJob() BasicJob {
	return *job.BasicJob
}
