package cloud66

import (
	"encoding/json"
)

type WorkflowWrapper struct {
	Workflow json.RawMessage `json:"pipeline"`
}

func (c *Client) GetWorkflow(stackUid, formationUid, snapshotUID string, useLatest bool, workflowName string) (*WorkflowWrapper, error) {
	params := struct {
		SnapshotUID string `json:"snapshot_uid"`
		UseLatest   bool   `json:"use_latest"`
		Workflow    string `json:"workflow"`
	}{
		SnapshotUID: snapshotUID,
		UseLatest:   useLatest,
		Workflow:    workflowName,
	}
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/formations/"+formationUid+"/pipeline.json", params, nil)
	if err != nil {
		return nil, err
	}
	var workflowWrapper *WorkflowWrapper
	return workflowWrapper, c.DoReq(req, &workflowWrapper, nil)
}
