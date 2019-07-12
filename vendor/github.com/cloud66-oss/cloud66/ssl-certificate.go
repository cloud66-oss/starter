package cloud66

import (
	"fmt"
	"strconv"
	"time"
)

var sslCertificateStatus = map[int]string{
	-1: "Not Applicable", // ST_NOT_APPLICABLE
	0:  "Not Installed",  // ST_NOT_INSTALLED
	1:  "Installing",     // ST_INSTALLING
	2:  "Failed",         // ST_FAILED
	3:  "Installed",      // ST_INSTALLED
	4:  "Removing",       // ST_REMOVING
}

type SslCertificate struct {
	Uuid                    string     `json:"uuid"`
	Name                    string     `json:"name"`
	ServerGroupID           int        `json:"server_group_id"`
	ServerNames             string     `json:"server_names"`
	SHA256Fingerprint       *string    `json:"sha256_fingerprint"`
	CAName                  *string    `json:"ca_name"`
	Type                    string     `json:"type"`
	SSLTermination          bool       `json:"ssl_termination"`
	HasIntermediateCert     bool       `json:"has_intermediate_cert"`
	Certificate             *string    `json:"certificate"`
	Key                     *string    `json:"key"`
	IntermediateCertificate *string    `json:"intermediate_certificate"`
	StatusCode              int        `json:"status"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
	ExpiresAt               *time.Time `json:"expires_at"`
}

const LetsEncryptSslCertificateType = "lets_encrypt"
const ManualSslCertificateType = "manual"

type wrappedSslCertificate struct {
	SslCertificate *SslCertificate `json:"ssl_certificate"`
}

func (ssl_certificate SslCertificate) Status() string {
	return sslCertificateStatus[ssl_certificate.StatusCode]
}

func (c *Client) ListSslCertificates(stackUID string) ([]SslCertificate, error) {
	queryStrings := make(map[string]string)
	queryStrings["page"] = "1"

	var p Pagination
	var result []SslCertificate
	var pageResult []SslCertificate

	for {
		req, err := c.NewRequest("GET", fmt.Sprintf("/stacks/%s/ssl_certificates.json", stackUID), nil, queryStrings)
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

func (c *Client) GetSslCertificate(stackUID string, sslCertificateUUID string) (*SslCertificate, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("/stacks/%s/ssl_certificates/%s.json", stackUID, sslCertificateUUID), nil, nil)
	if err != nil {
		return nil, err
	}

	var sslCertificateResult *SslCertificate
	return sslCertificateResult, c.DoReq(req, &sslCertificateResult, nil)
}

func (c *Client) CreateSslCertificate(stackUID string, sslCertificate *SslCertificate) (*SslCertificate, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("/stacks/%s/ssl_certificates.json", stackUID), wrappedSslCertificate{sslCertificate}, nil)
	if err != nil {
		return nil, err
	}

	var sslCertificateResult *SslCertificate
	return sslCertificateResult, c.DoReq(req, &sslCertificateResult, nil)
}

func (c *Client) UpdateSslCertificate(stackUID string, sslCertificateUUID string, sslCertificate *SslCertificate) (*SslCertificate, error) {
	req, err := c.NewRequest("PUT", fmt.Sprintf("/stacks/%s/ssl_certificates/%s.json", stackUID, sslCertificateUUID), wrappedSslCertificate{sslCertificate}, nil)
	if err != nil {
		return nil, err
	}

	var sslCertificateResult *SslCertificate
	return sslCertificateResult, c.DoReq(req, &sslCertificateResult, nil)
}

func (c *Client) DestroySslCertificate(stackUID string, sslCertificateUUID string) (*SslCertificate, error) {
	req, err := c.NewRequest("DELETE", fmt.Sprintf("/stacks/%s/ssl_certificates/%s.json", stackUID, sslCertificateUUID), nil, nil)
	if err != nil {
		return nil, err
	}

	var sslCertificateResult *SslCertificate
	return sslCertificateResult, c.DoReq(req, &sslCertificateResult, nil)
}
