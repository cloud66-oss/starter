package cloud66

import (
	"time"
)

type Policy struct {
	Uid       string    `json:"uid"`
	Name      string    `json:"name"`
	Selector  string    `json:"selector"`
	Sequence  int       `json:"sequence"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at_iso"`
	UpdatedAt time.Time `json:"updated_at_iso"`
	Tags      []string  `json:"tags"`
}

func (p Policy) String() string {
	return p.Name
}

func (c *Client) AddPolicies(stackUid string, formationUid string, policies []*Policy, message string) ([]Policy, error) {
	var policiesRes []Policy = make([]Policy, 0)
	var singleRes *Policy
	for _, policy := range policies {
		params := struct {
			Message string  `json:"message"`
			Policy  *Policy `json:"policy"`
		}{
			Message: message,
			Policy:  policy,
		}
		singleRes = nil

		req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/formations/"+formationUid+"/policies.json", params, nil)
		if err != nil {
			return nil, err
		}

		err = c.DoReq(req, &singleRes, nil)
		if err != nil {
			return nil, err
		}
		policiesRes = append(policiesRes, *singleRes)
	}

	return policiesRes, nil
}
