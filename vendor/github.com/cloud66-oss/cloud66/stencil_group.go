package cloud66

import (
	"time"
)

type StencilGroup struct {
	Uid       string    `json:"uid"`
	Name      string    `json:"name"`
	Rules     string    `json:"rules"`
	CreatedAt time.Time `json:"created_at_iso"`
	UpdatedAt time.Time `json:"updated_at_iso"`
	Tags      []string  `json:"tags"`
}

func (s StencilGroup) String() string {
	return s.Name
}

func (c *Client) AddStencilGroups(stackUid string, formationUid string, groups []*StencilGroup, message string) ([]StencilGroup, error) {
	var groupRes = make([]StencilGroup, 0)
	var singleRes *StencilGroup
	for _, stencilGroup := range groups {
		params := struct {
			Message      string        `json:"message"`
			StencilGroup *StencilGroup `json:"stencil_group"`
		}{
			Message:      message,
			StencilGroup: stencilGroup,
		}
		singleRes = nil

		req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/formations/"+formationUid+"/stencil_groups.json", params, nil)
		if err != nil {
			return nil, err
		}

		err = c.DoReq(req, &singleRes, nil)
		if err != nil {
			return nil, err
		}
		groupRes = append(groupRes, *singleRes)
	}

	return groupRes, nil
}
