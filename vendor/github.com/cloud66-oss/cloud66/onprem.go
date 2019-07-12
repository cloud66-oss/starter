package cloud66

import (
	"fmt"
	"time"
)

type Onprem struct {
	Uid       string      `json:"uid"`
	Name      string      `json:"name"`
	Config    interface{} `json:"config"`
	CreatedAt time.Time   `json:"created_at_iso"`
	UpdatedAt time.Time   `json:"updated_at_iso"`
}

func (c *Client) ListOnprems() ([]Onprem, error) {
	req, err := c.NewRequest("GET", "/onprems.json", nil, nil)
	if err != nil {
		return nil, err
	}

	var result []Onprem
	return result, c.DoReq(req, &result, nil)
}

func (c *Client) SaveOnprem(onprem Onprem) (*Onprem, error) {
	req, err := c.NewRequest("PUT", fmt.Sprintf("/onprems/%s.json", onprem.Uid), onprem, nil)
	if err != nil {
		return nil, err
	}
	var onpremRes *Onprem
	return onpremRes, c.DoReq(req, &onpremRes, nil)
}

func (c *Client) GetOnprem(uid string) (*Onprem, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("/onprems/%s.json", uid), nil, nil)
	if err != nil {
		return nil, err
	}

	var onpremRes *Onprem
	return onpremRes, c.DoReq(req, &onpremRes, nil)
}
