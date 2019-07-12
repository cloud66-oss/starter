package cloud66

func (c *Client) UnauthenticatedPing() error {
	req, err := c.NewRequest("GET", "/ping.json", nil, nil)
	if err != nil {
		return err
	}
	var res GenericResponse
	return c.DoReq(req, &res, nil)
}

func (c *Client) AuthenticatedPing() error {
	req, err := c.NewRequest("GET", "/ping/auth.json", nil, nil)
	if err != nil {
		return err
	}
	var res GenericResponse
	return c.DoReq(req, &res, nil)
}
