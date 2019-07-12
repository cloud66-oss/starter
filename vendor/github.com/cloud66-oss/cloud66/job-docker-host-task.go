package cloud66

import "encoding/json"

type dockerHostTaskJob struct {
	Command string `json:"command"`
}

type DockerHostTaskJob struct {
	*BasicJob
	*dockerHostTaskJob
}

func (job *DockerHostTaskJob) UnmarshalJSON(b []byte) error {
	var bj BasicJob
	if err := json.Unmarshal(b, &bj); err == nil {
		var j dockerHostTaskJob
		if err := json.Unmarshal(bj.basicJob.ParamsRaw, &j); err != nil {
			return err
		}
		*job = DockerHostTaskJob{BasicJob: &bj, dockerHostTaskJob: &j}
		return nil
	} else {
		return err
	}
}

func (job DockerHostTaskJob) GetBasicJob() BasicJob {
	return *job.BasicJob
}
