package cloud66

import (
	"fmt"
	"time"
)

type Gateway struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Address   string    `json:"address"`
	PrivateIp string    `json:"private_ip"`
	Ttl       string    `json:"ttl"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at_iso"`
	UpdatedAt time.Time `json:"updated_at_iso"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
}

func (c *Client) ListGateways(accountId int) ([]Gateway, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("/accounts/%d/gateways.json", accountId), nil, nil)
	if err != nil {
		return nil, err
	}

	var result []Gateway
	return result, c.DoReq(req, &result, nil)
}

func (c *Client) AddGateway(accountId int, name string, address string, username string,private_ip string) error {
	params := struct {
		Name      string `json:"name"`
		Address   string `json:"address"`
		PrivateIp string    `json:"private_ip"`
		Username  string `json:"username"`
	}{
		Name:      name,
		Address:   address,
		PrivateIp: private_ip,
		Username:  username,
	}

	req, err := c.NewRequest("POST", fmt.Sprintf("/accounts/%d/gateways.json", accountId), params, nil)
	if err != nil {
		return err
	}

	err = c.DoReq(req, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveGateway(accountId int, gatewayId int) error {
	req, err := c.NewRequest("DELETE", fmt.Sprintf("/accounts/%d/gateways/%d.json", accountId, gatewayId), nil, nil)
	if err != nil {
		return err
	}

	err = c.DoReq(req, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateGateway(accountId int, gatewayId int, keyContent string, ttl int) error {
	params := struct {
		Content string `json:"content"`
		Ttl     int    `json:"ttl"`
	}{
		Content: keyContent,
		Ttl:     ttl,
	}

	req, err := c.NewRequest("PUT", fmt.Sprintf("/accounts/%d/gateways/%d.json", accountId, gatewayId), params, nil)
	if err != nil {
		return err
	}

	err = c.DoReq(req, nil, nil)
	if err != nil {
		return err
	}

	return nil

}
