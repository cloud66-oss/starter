package cloud66

import "encoding/json"

type Job interface {
	GetBasicJob() BasicJob
}

func JobFactory(jobRes json.RawMessage) (*Job, error) {
	var T = struct {
		Type string `json:"type"`
	}{}

	if err := json.Unmarshal(jobRes, &T); err != nil {
		return nil, err
	}

	var job Job

	switch T.Type {
	case "DockerHostTaskJob":
		job = new(DockerHostTaskJob)
	case "DockerServiceTaskJob":
		job = new(DockerServiceTaskJob)
	default:
		job = new(BasicJob)
	}

	if err := json.Unmarshal(jobRes, &job); err != nil {
		return nil, err
	}

	return &job, nil
}
