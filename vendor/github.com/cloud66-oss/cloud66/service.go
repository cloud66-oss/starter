package cloud66

import "strconv"

type Service struct {
	Name          string      `json:"name"`
	Containers    []Container `json:"containers"`
	SourceType    string      `json:"source_type"`
	GitRef        string      `json:"git_ref"`
	ImageName     string      `json:"image_name"`
	ImageUid      string      `json:"image_uid"`
	ImageTag      string      `json:"image_tag"`
	Command       string      `json:"command"`
	BuildCommand  string      `json:"build_command"`
	DeployCommand string      `json:"deploy_command"`
	WrapCommand   string      `json:"wrap_command"`
	DesiredCount  int         `json:"desired_count"`
}

func (c *Client) GetServices(stackUid string, serverUid *string) ([]Service, error) {
	var params interface{}
	if serverUid == nil {
		params = nil
	} else {
		params = struct {
			ServerUid string `json:"server_uid"`
		}{
			ServerUid: *serverUid,
		}
	}

	query_strings := make(map[string]string)
	query_strings["page"] = "1"

	var p Pagination
	var result []Service
	var serviceRes []Service

	for {
		req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/services.json", params, query_strings)
		if err != nil {
			return nil, err
		}

		serviceRes = nil
		err = c.DoReq(req, &serviceRes, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, serviceRes...)
		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}

	}

	return result, nil
}

func (c *Client) GetService(stackUid string, serviceName string, serverUid *string, wrapCommand *string) (*Service, error) {
	var params interface{}
	if serverUid == nil {
		params = nil
	} else {
		if wrapCommand == nil {
			params = struct {
				ServerUid string `json:"server_uid"`
			}{
				ServerUid: *serverUid,
			}
		} else {
			params = struct {
				ServerUid   string `json:"server_uid"`
				WrapCommand string `json:"wrap_command"`
			}{
				ServerUid:   *serverUid,
				WrapCommand: *wrapCommand,
			}
		}
	}
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/services/"+serviceName+".json", params, nil)
	if err != nil {
		return nil, err
	}
	var servicesRes *Service
	return servicesRes, c.DoReq(req, &servicesRes, nil)
}

func (c *Client) StopService(stackUid string, serviceName string, serverUid *string) (*AsyncResult, error) {
	var params interface{}
	if serverUid == nil {
		params = nil
	} else {
		params = struct {
			ServerUid string `json:"server_uid"`
		}{
			ServerUid: *serverUid,
		}
	}
	req, err := c.NewRequest("DELETE", "/stacks/"+stackUid+"/services/"+serviceName+".json", params, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}

func (c *Client) ScaleService(stackUid string, serviceName string, serverCount map[string]int) (*AsyncResult, error) {
	params := struct {
		ServiceName string         `json:"service_name"`
		ServerCount map[string]int `json:"server_count"`
	}{
		ServiceName: serviceName,
		ServerCount: serverCount,
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/services.json", params, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}

func (c *Client) ScaleServiceByGroup(stackUid string, serviceName string, groupCount map[string]int) (*AsyncResult, error) {
	params := struct {
		ServiceName      string         `json:"service_name"`
		ServerGroupCount map[string]int `json:"server_group_count"`
	}{
		ServiceName:      serviceName,
		ServerGroupCount: groupCount,
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/services.json", params, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}

func (s *Service) ServerContainerCountMap() map[string]int {
	var serverMap = make(map[string]int)
	for _, container := range s.Containers {
		if _, present := serverMap[container.ServerName]; present {
			serverMap[container.ServerName] = serverMap[container.ServerName] + 1
		} else {
			serverMap[container.ServerName] = 1
		}
	}
	return serverMap
}

func (c *Client) InvokeServiceAction(stackUid string, serviceName *string, serverUid *string, action string) (*AsyncResult, error) {
	var params interface{}
	if serverUid != nil && serviceName != nil {
		params = struct {
			Command     string `json:"command"`
			ServiceName string `json:"service_name"`
			ServerUid   string `json:"server_uid"`
		}{
			Command:     action,
			ServiceName: *serviceName,
			ServerUid:   *serverUid,
		}
	} else if serverUid == nil {
		params = struct {
			Command     string `json:"command"`
			ServiceName string `json:"service_name"`
		}{
			Command:     action,
			ServiceName: *serviceName,
		}
	} else if serviceName == nil {
		params = struct {
			Command   string `json:"command"`
			ServerUid string `json:"server_uid"`
		}{
			Command:   action,
			ServerUid: *serverUid,
		}
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/actions.json", params, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}
