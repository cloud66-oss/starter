package registry

type repositoriesResponse struct {
	Repositories []string `json:"repositories"`
}

func (registry *Registry) Repositories() (repos []string, err error) {
	url := registry.url("/v2/_catalog")

	var response repositoriesResponse
	for {
		registry.Logf("registry.repositories url=%s", url)
		url, err = registry.getPaginatedJson(url, &response)
		switch err {
		case ErrNoMorePages:
			repos = append(repos, response.Repositories...)
			return
		case nil:
			repos = append(repos, response.Repositories...)
			continue
		default:
			return nil, err
		}
	}
}
