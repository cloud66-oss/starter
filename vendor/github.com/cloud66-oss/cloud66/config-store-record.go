package cloud66

import (
	"net/url"
	"strconv"
	"strings"
)

const (
	BundledConfigStoreAccountScope = "account"
	BundledConfigStoreStackScope   = "stack"
)

type BundledConfigStoreRecords struct {
	Records []BundledConfigStoreRecord `json:"records" yaml:"records"`
}

type BundledConfigStoreRecord struct {
	ConfigStoreRecord `yaml:",inline"`
	Scope             string `json:"scope" yaml:"scope"`
}

type ConfigStoreRecord struct {
	Key      string            `json:"key" yaml:"key"`
	RawValue string            `json:"raw_value" yaml:"raw_value"`
	Metadata map[string]string `json:"metadata" yaml:"metadata"`
	Ttl      int               `json:"ttl" yaml:"ttl"`
}

type configStoreRequestWrapper struct {
	Record *ConfigStoreRecord `json:"record" yaml:"record"`
}

func (c *Client) GetConfigStoreRecords(namespace string) ([]ConfigStoreRecord, error) {
	var p Pagination
	var result []ConfigStoreRecord

	query_strings := make(map[string]string)
	query_strings["page"] = "1"

	for {
		req, err := c.NewRequest("GET", "/configstore/namespaces/"+namespace+"/records.json", nil, query_strings)
		if err != nil {
			return nil, err
		}

		var intermediateResult []ConfigStoreRecord
		err = c.DoReq(req, &intermediateResult, &p)
		if err != nil {
			return nil, err
		}
		result = append(result, intermediateResult...)

		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}

	}

	return result, nil
}

func (c *Client) GetConfigStoreRecord(namespace, key string) (*ConfigStoreRecord, error) {
	req, err := c.NewRequest("GET", "/configstore/namespaces/"+namespace+"/records/"+urlEncodedKey(key)+".json", nil, nil)
	if err != nil {
		return nil, err
	}

	var result *ConfigStoreRecord
	err = c.DoReq(req, &result, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) CreateConfigStoreRecord(namespace string, record *ConfigStoreRecord) (*ConfigStoreRecord, error) {
	req, err := c.NewRequest("POST", "/configstore/namespaces/"+namespace+"/records.json", &configStoreRequestWrapper{Record: record}, nil)
	if err != nil {
		return nil, err
	}

	var result *ConfigStoreRecord
	err = c.DoReq(req, &result, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) UpdateConfigStoreRecord(namespace, key string, record *ConfigStoreRecord) (*ConfigStoreRecord, error) {
	req, err := c.NewRequest("PUT", "/configstore/namespaces/"+namespace+"/records/"+urlEncodedKey(key)+".json", &configStoreRequestWrapper{Record: record}, nil)
	if err != nil {
		return nil, err
	}

	var result *ConfigStoreRecord
	err = c.DoReq(req, &result, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) DeleteConfigStoreRecord(namespace, key string) (*ConfigStoreRecord, error) {
	req, err := c.NewRequest("DELETE", "/configstore/namespaces/"+namespace+"/records/"+urlEncodedKey(key)+".json", nil, nil)
	if err != nil {
		return nil, err
	}

	var result *ConfigStoreRecord
	err = c.DoReq(req, &result, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func urlEncodedKey(key string) string {
	result := url.QueryEscape(key)
	result = strings.ReplaceAll(result, ".", "%2E")
	return result
}
