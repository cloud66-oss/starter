package cloud66

import (
	"fmt"
	"strconv"
	"time"
)

var baseTemplateStatus = map[int]string{
	1: "Unknown",                          // ST_UNKNOWN
	2: "Queued to be pulled and verified", // ST_QUEUED
	3: "Pulling repository",               // ST_PULLING
	4: "Verifying repository",             // ST_VERIFYING
	5: "Failed to pull the repository",    // ST_CONNECTION_ERROR
	6: "Available",                        // ST_AVAILABLE
	7: "Failed to verify the repository",  // ST_VERIFICATION_ERROR
}

type StencilTemplate struct {
	BaseTemplate      string   `json:"base_template"`
	Filename          string   `json:"filename"`
	FilenamePattern   string   `json:"filename_pattern"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	ContextType       string   `json:"context_type"`
	Tags              []string `json:"tags"`
	PreferredSequence int      `json:"preferred_sequence"`
	Content           string   `json:"content"`
}

type BaseTemplate struct {
	Uid        string            `json:"uid"`
	Name       string            `json:"name"`
	ShortName  string            `json:"short_name"`
	GitRepo    string            `json:"git_repo"`
	GitBranch  string            `json:"git_branch"`
	StatusCode int               `json:"status"`
	Stencils   []StencilTemplate `json:"stencils"`
	LastSync   *time.Time        `json:"last_sync_iso"`
	CreatedAt  time.Time         `json:"created_at_iso"`
	UpdatedAt  time.Time         `json:"updated_at_iso"`
}

type wrappedBaseTemplate struct {
	BaseTemplate *BaseTemplate `json:"base_template"`
}

func (bt BaseTemplate) String() string {
	return fmt.Sprintf("%s:%s", bt.GitRepo, bt.GitBranch)
}

func (bt BaseTemplate) Status() string {
	return baseTemplateStatus[bt.StatusCode]
}

func (c *Client) ListBaseTemplates() ([]BaseTemplate, error) {
	queryStrings := make(map[string]string)
	queryStrings["page"] = "1"

	var p Pagination
	var result []BaseTemplate
	var pageResult []BaseTemplate

	for {
		req, err := c.NewRequest("GET", "/base_templates.json", nil, queryStrings)
		if err != nil {
			return nil, err
		}

		pageResult = nil
		err = c.DoReq(req, &pageResult, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, pageResult...)
		if p.Current < p.Next {
			queryStrings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}
	}

	return result, nil
}

func (c *Client) GetBaseTemplate(baseTemplateUID string, includeStencils bool, includeContent bool) (*BaseTemplate, error) {
	queryStrings := make(map[string]string)
	if includeStencils {
		queryStrings["include_stencils"] = "1"
	}

	if includeContent {
		queryStrings["include_content"] = "1"
	}

	req, err := c.NewRequest("GET", fmt.Sprintf("/base_templates/%s.json", baseTemplateUID), nil, queryStrings)
	if err != nil {
		return nil, err
	}

	var baseTemplateResult *BaseTemplate
	return baseTemplateResult, c.DoReq(req, &baseTemplateResult, nil)
}

func (c *Client) UpdateBaseTemplate(baseTemplateUID string, baseTemplate *BaseTemplate) (*BaseTemplate, error) {
	req, err := c.NewRequest("PUT", fmt.Sprintf("/base_templates/%s.json", baseTemplateUID), wrappedBaseTemplate{baseTemplate}, nil)
	if err != nil {
		return nil, err
	}

	var baseTemplateResult *BaseTemplate
	return baseTemplateResult, c.DoReq(req, &baseTemplateResult, nil)
}

func (c *Client) CreateBaseTemplate(baseTemplate *BaseTemplate) (*BaseTemplate, error) {
	req, err := c.NewRequest("POST", "/base_templates.json", wrappedBaseTemplate{baseTemplate}, nil)
	if err != nil {
		return nil, err
	}

	var baseTemplateResult *BaseTemplate
	return baseTemplateResult, c.DoReq(req, &baseTemplateResult, nil)
}

func (c *Client) DestroyBaseTemplate(baseTemplateUID string) (*BaseTemplate, error) {
	req, err := c.NewRequest("DELETE", fmt.Sprintf("/base_templates/%s.json", baseTemplateUID), nil, nil)
	if err != nil {
		return nil, err
	}

	var baseTemplateResult *BaseTemplate
	return baseTemplateResult, c.DoReq(req, &baseTemplateResult, nil)
}

func (c *Client) SyncBaseTemplate(baseTemplateUID string) (*BaseTemplate, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("/base_templates/%s/sync.json", baseTemplateUID), nil, nil)
	if err != nil {
		return nil, err
	}

	var baseTemplateResult *BaseTemplate
	return baseTemplateResult, c.DoReq(req, &baseTemplateResult, nil)
}
