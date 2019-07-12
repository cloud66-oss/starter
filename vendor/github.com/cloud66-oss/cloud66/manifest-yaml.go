package cloud66

import (
	"strconv"
	"time"
)

type ManifestYaml struct {
	Uid       string    `json:"uid"`
	Body      string    `json:"body"`
	Comments  string    `json:"comments"`
	CreatedAt time.Time `json:"created_at_iso"`
	UpdatedAt time.Time `json:"updated_at_iso"`
}

func (c *Client) ManifestYamlList(stackUid string, include_body bool) ([]ManifestYaml, error) {
	query_strings := make(map[string]string)
	query_strings["page"] = "1"
	query_strings["include_body"] = strconv.FormatBool(include_body)

	var p Pagination
	var result []ManifestYaml
	var manifestYamlRes []ManifestYaml

	for {
		req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/manifest_yaml.json", nil, query_strings)
		if err != nil {
			return nil, err
		}

		manifestYamlRes = nil
		err = c.DoReq(req, &manifestYamlRes, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, manifestYamlRes...)
		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}

	}

	return result, nil
}

func (c *Client) ManifestYamlInfo(stackUid string, version string) (*ManifestYaml, error) {
	var manifestYamlRes ManifestYaml

	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/manifest_yaml/"+version+".json", nil, nil)
	if err != nil {
		return nil, err
	}

	err = c.DoReq(req, &manifestYamlRes, nil)
	if err != nil {
		return nil, err
	}

	return &manifestYamlRes, nil
}

func (c *Client) CreateManifestYaml(stackUid, manifestYaml, comments string) (*ManifestYaml, error) {
	params := struct {
		ManifestYaml string `json:"manifest_yaml"`
		Comments     string `json:"comments"`
	}{
		ManifestYaml: manifestYaml,
		Comments:     comments,
	}

	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/manifest_yaml", params, nil)
	if err != nil {
		return nil, err
	}

	var manifestYamlRes *ManifestYaml
	return manifestYamlRes, c.DoReq(req, &manifestYamlRes, nil)
}
