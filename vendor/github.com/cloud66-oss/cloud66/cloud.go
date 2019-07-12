package cloud66

type Cloud struct {
	Id          string            `json:"name"`
	Name        string            `json:"display_name"`
	KeyName     string            `json:"key_name"`
	Regions     []CloudRegion     `json:"regions"`
	ServerSizes []CloudServerSize `json:"server_sizes"`
}

type CloudServerSize struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CloudRegion struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) GetCloudsInfo() ([]Cloud, error) {
	req, err := c.NewRequest("GET", "/clouds.json", nil, nil)
	if err != nil {
		return nil, err
	}
	var cloudRes []Cloud
	return cloudRes, c.DoReq(req, &cloudRes, nil)
}

func (c *Client) GetCloudInfo(cloudName string) (*Cloud, error) {
	req, err := c.NewRequest("GET", "/clouds/"+cloudName+".json", nil, nil)
	if err != nil {
		return nil, err
	}
	var cloudRes *Cloud
	return cloudRes, c.DoReq(req, &cloudRes, nil)
}
