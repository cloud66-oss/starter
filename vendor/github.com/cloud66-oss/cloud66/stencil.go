package cloud66

import (
	"encoding/base64"
	"time"
)

type Stencil struct {
	Uid              string    `json:"uid"`
	Filename         string    `json:"filename"`
	TemplateFilename string    `json:"template_filename"`
	ContextID        string    `json:"context_id"`
	Status           int       `json:"status"`
	Tags             []string  `json:"tags"`
	Inline           bool      `json:"inline"`
	GitfilePath      string    `json:"gitfile_path"`
	Body             string    `json:"body"`
	BtrRepo          string    `json:"btr_repo"`
	BtrBranch        string    `json:"btr_branch"`
	Sequence         int       `json:"sequence"`
	CreatedAt        time.Time `json:"created_at_iso"`
	UpdatedAt        time.Time `json:"updated_at_iso"`
}

func (s Stencil) String() string {
	return s.Filename
}

func (c *Client) AddStencils(stackUid string, formationUid string, baseTemplateUid string, stencils []*Stencil, message string) (*AsyncResult, error) {
	params := struct {
		Message      string     `json:"message"`
		BaseTemplate string     `json:"btr_uuid"`
		Stencils     []*Stencil `json:"stencils"`
	}{
		Message:      message,
		BaseTemplate: baseTemplateUid,
		Stencils:     stencils,
	}

	if len(stencils) > 0 {
		req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/formations/"+formationUid+"/stencils.json", params, nil)
		if err != nil {
			return nil, err
		}

		var asyncRes *AsyncResult
		return asyncRes, c.DoReq(req, &asyncRes, nil)
	}

	return nil, nil
}

func (c *Client) RenderStencil(stackUID, snapshotUID, formationUID, stencilUID string, body []byte) (*Renders, error) {
	encoded := base64.StdEncoding.EncodeToString(body)
	params := struct {
		Body        string `json:"body"`
		SnapshotUID string `json:"snapshot_id"`
	}{
		SnapshotUID: snapshotUID,
		Body:        encoded,
	}

	var result *Renders
	req, err := c.NewRequest("POST", "/stacks/"+stackUID+"/formations/"+formationUID+"/stencils/"+stencilUID+"/render.json", params, nil)

	if err != nil {
		return nil, err
	}

	result = nil
	err = c.DoReq(req, &result, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) UpdateStencil(stackUID, formationUID, stencilUID, message string, body []byte) (*Stencil, error) {
	encoded := base64.StdEncoding.EncodeToString(body)
	params := struct {
		Body    string `json:"body"`
		Message string `json:"message"`
	}{
		Message: message,
		Body:    encoded,
	}

	var result *Stencil
	req, err := c.NewRequest("PUT", "/stacks/"+stackUID+"/formations/"+formationUID+"/stencils/"+stencilUID+".json", params, nil)

	if err != nil {
		return nil, err
	}

	result = nil
	err = c.DoReq(req, &result, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}
