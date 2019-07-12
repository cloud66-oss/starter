package cloud66

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type UnmanagedServer struct {
	Vendor string `json:"vendor"`
	Id     string `json:"id"`
}

type Account struct {
	Id                   int               `json:"id"`
	Owner                string            `json:"owner"`
	Name                 string            `json:"friendly_name"`
	StackCount           int               `json:"stack_count"`
	UsedClouds           []string          `json:"used_clouds"`
	CreatedAt            time.Time         `json:"created_at"`
	UpdatedAt            time.Time         `json:"updated_at"`
	CurrentAccount       bool              `json:"current_account"`
	UnmanagedServers     []UnmanagedServer `json:"unmanaged_servers"`
	ServerRegistration   string            `json:"server_registration_script"`
	ConfigStoreNamespace string            `json:"configstore_namespace"`
}

type OTP struct {
	OK   bool   `json:"ok"`
	Code string `json:"otp"`
}

func (c *Client) AccountInfo(accountId int, getUnmanaged bool) (*Account, error) {
	params := struct {
		Value bool `json:"include_servers"`
	}{
		Value: getUnmanaged,
	}

	req, err := c.NewRequest("GET", fmt.Sprintf("/accounts/%d.json", accountId), params, nil)
	if err != nil {
		return nil, err
	}

	var accountRes *Account
	return accountRes, c.DoReq(req, &accountRes, nil)
}

func (c *Client) AccountInfos() ([]Account, error) {
	query_strings := make(map[string]string)
	query_strings["page"] = "1"

	var p Pagination
	var result []Account
	var accountRes []Account

	for {
		req, err := c.NewRequest("GET", "/accounts.json", nil, query_strings)
		if err != nil {
			return nil, err
		}

		accountRes = nil
		err = c.DoReq(req, &accountRes, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, accountRes...)
		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}

	}

	return result, nil
}

// AccountOTP returns the OTP from the server. If accountId and stackId are passed
// (!= 0) then they will be passed up
func (c *Client) AccountOTP() (string, error) {
	var req *http.Request
	var err error
	req, err = c.NewRequest("GET", "/accounts/otp.json", nil, nil)

	var otp *OTP
	err = c.DoReq(req, &otp, nil)
	if err != nil {
		return "", err
	}

	return otp.Code, nil
}
