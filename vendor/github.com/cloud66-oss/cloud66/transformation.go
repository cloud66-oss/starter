package cloud66

import "time"

type Transformation struct {
	Uid       string    `json:"uid"`
	Name      string    `json:"name"`
	Selector  string    `json:"selector"`
	Sequence  int       `json:"sequence"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at_iso"`
	UpdatedAt time.Time `json:"updated_at_iso"`
	Tags      []string  `json:"tags"`
}

func (t Transformation) String() string {
	return t.Name
}

func (c *Client) AddTransformations(stackUid string, formationUid string, transformations []*Transformation, message string) ([]Transformation, error) {
	var transformationsRes []Transformation = make([]Transformation, 0)
	var singleRes *Transformation
	for _, transformation := range transformations {
		params := struct {
			Message        string          `json:"message"`
			Transformation *Transformation `json:"transformation"`
		}{
			Message:        message,
			Transformation: transformation,
		}
		singleRes = nil

		req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/formations/"+formationUid+"/transformations.json", params, nil)
		if err != nil {
			return nil, err
		}

		err = c.DoReq(req, &singleRes, nil)
		if err != nil {
			return nil, err
		}
		transformationsRes = append(transformationsRes, *singleRes)
	}

	return transformationsRes, nil
}
