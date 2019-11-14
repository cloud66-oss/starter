package cloud66

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Formation struct {
	Uid             string           `json:"uid"`
	Name            string           `json:"name"`
	Stencils        []Stencil        `json:"stencils"`
	StencilGroups   []StencilGroup   `json:"stencil_groups"`
	BaseTemplates   []BaseTemplate   `json:"base_templates"`
	Policies        []Policy         `json:"policies"`
	Transformations []Transformation `json:"transformations"`
	Workflows       []Workflow       `json:"workflows"`
	HelmReleases    []HelmRelease    `json:"helm_releases"`
	CreatedAt       time.Time        `json:"created_at_iso"`
	UpdatedAt       time.Time        `json:"updated_at_iso"`
	Tags            []string         `json:"tags"`
}

func (f *Formation) FindStencil(stencilName string) *Stencil {
	for _, stencil := range f.Stencils {
		if stencil.Filename == stencilName {
			return &stencil
		}
	}

	return nil
}

func (c *Client) Formations(stackUid string, fullContent bool) ([]Formation, error) {
	query_strings := make(map[string]string)
	query_strings["page"] = "1"
	if fullContent {
		query_strings["full_content"] = "1"
	}

	var p Pagination
	var result []Formation
	var formationRes []Formation

	for {
		req, err := c.NewRequest("GET", fmt.Sprintf("/stacks/%s/formations.json", stackUid), nil, query_strings)
		if err != nil {
			return nil, err
		}

		formationRes = nil
		err = c.DoReq(req, &formationRes, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, formationRes...)
		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}
	}

	return result, nil
}

func (c *Client) CreateFormation(stackUid string, name string, templateRepo string, templateBranch string, tags []string) (*Formation, error) {
	type base struct {
		Repo   string `json:"repo"`
		Branch string `json:"branch"`
	}

	params := struct {
		Base base     `json:"base_template"`
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}{
		Name: name,
		Tags: tags,
	}
	params.Base = base{
		Repo:   templateRepo,
		Branch: templateBranch,
	}

	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/formations.json", params, nil)
	if err != nil {
		return nil, err
	}

	var formationRes *Formation
	err = c.DoReq(req, &formationRes, nil)
	if err != nil {
		return nil, err
	}

	return formationRes, nil
}

func (c *Client) CreateFormationMultiBtr(stackUid string, name string, baseTemplates []*BaseTemplate, tags []string) (*Formation, error) {
	type base struct {
		Repo   string `json:"repo"`
		Branch string `json:"branch"`
	}

	params := struct {
		Base []base   `json:"base_templates"`
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}{
		Name: name,
		Tags: tags,
	}
	baseList := make([]base, 0)
	for _, value := range baseTemplates {
		baseList = append(baseList, base{
			Repo:   value.GitRepo,
			Branch: value.GitBranch,
		})
	}

	params.Base = baseList
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/formations.json", params, nil)
	if err != nil {
		return nil, err
	}

	var formationRes *Formation
	err = c.DoReq(req, &formationRes, nil)
	if err != nil {
		return nil, err
	}

	return formationRes, nil
}

func (f *Formation) FindIndexByRepoAndBranch(repo string, branch string) int {
	repo = strings.TrimSpace(repo)
	branch = strings.TrimSpace(branch)
	for index, btr := range f.BaseTemplates {
		if strings.TrimSpace(btr.GitRepo) == repo && strings.TrimSpace(btr.GitBranch) == branch {
			return index
		}
	}
	return -1
}
