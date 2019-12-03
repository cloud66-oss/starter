package cloud66

import (
	"time"
)

type Workflow struct {
	Uid       string    `json:"uid"`
	Name      string    `json:"name"`
	Body      string    `json:"body"`
	Default   bool      `json:"default"`
	CreatedAt time.Time `json:"created_at_iso"`
	UpdatedAt time.Time `json:"updated_at_iso"`
	Tags      []string  `json:"tags"`
}

func (p Workflow) String() string {
	return p.Name
}

func (c *Client) AddWorkflow(stackUid string, formationUid string, workflow *Workflow, message string) (*Workflow, error) {
	var workflowRes *Workflow = nil

	params := struct {
		Message  string    `json:"message"`
		Workflow *Workflow `json:"workflow"`
	}{
		Message:  message,
		Workflow: workflow,
	}

	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/formations/"+formationUid+"/workflows.json", params, nil)
	if err != nil {
		return nil, err
	}

	err = c.DoReq(req, &workflowRes, nil)
	if err != nil {
		return nil, err
	}

	return workflowRes, nil
}
