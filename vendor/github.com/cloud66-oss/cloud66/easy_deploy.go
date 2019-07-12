package cloud66

import (
	"fmt"
	"strconv"
)

type EasyDeployMaintainer struct {
	Name    *string `json:"name"`
	Email   string  `json:"email"`
	Company *string `json:"company"`
	Offical *bool   `json:"official"`
	Trusted *bool   `json:"trusted"`
}

type EasyDeploy struct {
	Name        string               `json:"name"`
	DisplayName *string              `json:"display_name"`
	Version     string               `json:"version"`
	Uid         string               `json:"uid"`
	CreatedAt   string               `json:"created_at"`
	Logo        *string              `json:"logo"`
	Maintainer  EasyDeployMaintainer `json:"maintainer"`
}

func (c *Client) EasyDeployList() ([]string, error) {
	query_strings := make(map[string]string)
	query_strings["page"] = "1"

	var p Pagination
	var result []string
	var easyDeploy []string

	for {
		req, err := c.NewRequest("GET", "/easy_deploys.json", nil, query_strings)
		if err != nil {
			return nil, err
		}

		easyDeploy = nil
		err = c.DoReq(req, &easyDeploy, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, easyDeploy...)
		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}

	}

	return result, nil

}

func (c *Client) EasyDeployInfo(name string) (*EasyDeploy, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("/easy_deploys/%s.json", name), nil, nil)
	if err != nil {
		return nil, err
	}

	var easyDeploy *EasyDeploy
	return easyDeploy, c.DoReq(req, &easyDeploy, nil)
}
