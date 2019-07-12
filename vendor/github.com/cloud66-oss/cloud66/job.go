package cloud66

import "strconv"
import "encoding/json"
import "fmt"

var JobStatus = map[int]string{
	0: "Updated", // ST_UPDATED
	1: "Started", // ST_STARTED
	2: "Success", // ST_SUCCESS
	3: "Failed",  // ST_FAILED
}

type basicJob struct {
	Id        int             `json:"id"`
	Uid       string          `json:"uid"`
	Name      string          `json:"name"`
	Type      string          `json:"type"`
	Cron      string          `json:"cron"`
	Status    int             `json:"status"`
	ParamsRaw json.RawMessage `json:"params"`
	Params    map[string]interface{}
}

type BasicJob struct {
	*basicJob
}

func (bj *BasicJob) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &bj.basicJob); err == nil {
		if err = json.Unmarshal(bj.ParamsRaw, &bj.basicJob.Params); err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}

func (bj BasicJob) GetBasicJob() BasicJob {
	return bj
}

func (c *Client) GetJobs(stackUid string, serverUid *string) ([]Job, error) {
	fmt.Printf("")
	var params interface{}
	if serverUid == nil {
		params = nil
	} else {
		params = struct {
			ServerUid string `json:"server_uid"`
		}{
			ServerUid: *serverUid,
		}
	}

	query_strings := make(map[string]string)

	var p Pagination
	var result []Job
	var jobRes []*json.RawMessage

	for {
		req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/jobs.json", params, query_strings)
		if err != nil {
			return nil, err
		}

		jobRes = nil
		err = c.DoReq(req, &jobRes, &p)
		if err != nil {
			return nil, err
		}

		for _, j := range jobRes {
			var job *Job
			if job, err = JobFactory(*j); err != nil {
				return nil, err
			}
			result = append(result, *job)
		}

		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}

	}

	return result, nil
}

func (c *Client) GetJob(stackUid string, jobUid string) (*Job, error) {
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/jobs/"+jobUid+".json", nil, nil)
	if err != nil {
		return nil, err
	}
	var jobRes *json.RawMessage
	err = c.DoReq(req, &jobRes, nil)
	if err != nil {
		return nil, err
	}
	return JobFactory(*jobRes)
}

func (c *Client) RunJobNow(stackUid string, jobUid string, jobArgs *string) (*AsyncResult, error) {
	var params interface{}
	if jobArgs == nil {
		params = nil
	} else {
		params = struct {
			JobArgs string `json:"job_args"`
		}{
			JobArgs: *jobArgs,
		}
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/jobs/"+jobUid+"/run_now.json", params, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}
