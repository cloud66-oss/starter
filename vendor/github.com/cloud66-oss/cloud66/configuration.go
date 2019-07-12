package cloud66

import (
	"strconv"
	"time"
)

type Configuration struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Body      string    `json:"body"`
	Comments  string    `json:"comments"`
	CanApply  bool      `json:"can_apply"`
	ChangedBy string    `json:"changed_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Client) ConfigurationList(stackUid string) ([]Configuration, error) {
	queryStrings := make(map[string]string)
	queryStrings["page"] = "1"

	var p Pagination
	var configurations []Configuration
	var configurationRes []Configuration

	for {
		req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/configuration.json", nil, queryStrings)
		if err != nil {
			return nil, err
		}

		configurationRes = nil
		err = c.DoReq(req, &configurationRes, &p)
		if err != nil {
			return nil, err
		}

		configurations = append(configurations, configurationRes...)
		if p.Current < p.Next {
			queryStrings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}
	}
	return configurations, nil
}

func (c *Client) ConfigurationDownload(stackUid, theType string) (*Configuration, error) {
	var configurationRes Configuration
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/configuration/"+theType+"/show.json", nil, nil)
	if err != nil {
		return nil, err
	}
	err = c.DoReq(req, &configurationRes, nil)
	if err != nil {
		return nil, err
	}
	return &configurationRes, nil
}

func (c *Client) ConfigurationUpload(stackUid, theType, commitMessage, body string, mustApply bool) (*AsyncResult, error) {
	params := struct {
		CommitMessage string `json:"commit_message"`
		Body          string `json:"body"`
		MustApply     bool   `json:"must_apply"`
	}{
		CommitMessage: commitMessage,
		Body:          body,
		MustApply:     mustApply,
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/configuration/"+theType+"/update.json", params, nil)
	if err != nil {
		return nil, err
	}
	var asyncResult *AsyncResult
	return asyncResult, c.DoReq(req, &asyncResult, nil)
}

func (c *Client) ConfigurationApply(stackUid, theType string) (*AsyncResult, error) {
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/configuration/"+theType+"/apply.json", nil, nil)
	if err != nil {
		return nil, err
	}
	var asyncResult *AsyncResult
	return asyncResult, c.DoReq(req, &asyncResult, nil)
}
