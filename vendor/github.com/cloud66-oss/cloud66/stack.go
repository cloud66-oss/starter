package cloud66

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var stackStatus = map[int]string{
	0: "Pending analysis",      //STK_QUEUED
	1: "Deployed successfully", //STK_SUCCESS
	2: "Deployment failed",     //STK_FAILED
	3: "Analyzing",             //STK_ANALYSING
	4: "Analyzed",              //STK_ANALYSED
	5: "Queued for deployment", //STK_QUEUED_FOR_DEPLOYING
	6: "Deploying",             //STK_DEPLOYING
	7: "Unable to analyze",     //STK_TERMINAL_FAILURE
}

var skycapStatus = map[int]string{
	0: "Pending analysis",   //STK_QUEUED
	1: "Built successfully", //STK_SUCCESS
	2: "Build failed",       //STK_FAILED
	3: "Analyzing",          //STK_ANALYSING
	4: "Analyzed",           //STK_ANALYSED
	5: "Queued for build",   //STK_QUEUED_FOR_DEPLOYING
	6: "Building",           //STK_DEPLOYING
	7: "Unable to analyze",  //STK_TERMINAL_FAILURE
}

var healthStatus = map[int]string{
	0: "Unknown",  //HLT_UNKNOWN
	1: "Building", //HLT_BUILDING
	2: "Impaired", //HLT_PARTIAL
	3: "Healthy",  //HLT_OK
	4: "Failed",   //HLT_BROKEN
}

type Stack struct {
	Uid                  string     `json:"uid"`
	Name                 string     `json:"name"`
	Git                  string     `json:"git"`
	GitBranch            string     `json:"git_branch"`
	Environment          string     `json:"environment"`
	Cloud                string     `json:"cloud"`
	Fqdn                 string     `json:"fqdn"`
	Language             string     `json:"language"`
	Framework            string     `json:"framework"`
	StatusCode           int        `json:"status"`
	HealthCode           int        `json:"health"`
	MaintenanceMode      bool       `json:"maintenance_mode"`
	HasLoadBalancer      bool       `json:"has_loadbalancer"`
	RedeployHook         *string    `json:"redeploy_hook"`
	LastActivity         *time.Time `json:"last_activity_iso"`
	UpdatedAt            time.Time  `json:"updated_at_iso"`
	CreatedAt            time.Time  `json:"created_at_iso"`
	DeployDir            string     `json:"deploy_directory"`
	Backend              string     `json:"backend"`
	Version              string     `json:"version"`
	Revision             string     `json:"revision"`
	Namespaces           []string   `json:"namespaces"`
	AccountId            int        `json:"account_id"`
	AccountName          string     `json:"account_name"`
	IsCluster            bool       `json:"is_cluster"`
	IsInsideCluster      bool       `json:"is_inside_cluster"`
	ClusterName          string     `json:"cluster_name"`
	ApplicationAddress   *string    `json:"application_address"`
	ConfigStoreNamespace string     `json:"configstore_namespace"`
}

type StackSetting struct {
	Key      string      `json:"key"`
	Value    interface{} `json:"value"`
	Readonly bool        `json:"readonly"`
}

type StackEnvVarHistory struct {
	Value     interface{} `json:"value"`
	CreatedAt time.Time   `json:"created_at_iso"`
	UpdatedAt time.Time   `json:"updated_at_iso"`
}

type StackEnvVar struct {
	Key       string               `json:"key"`
	Value     interface{}          `json:"value"`
	Readonly  bool                 `json:"readonly"`
	CreatedAt time.Time            `json:"created_at_iso"`
	UpdatedAt time.Time            `json:"updated_at_iso"`
	History   []StackEnvVarHistory `json:"history"`
}

type RedeployResponse struct {
	Status  bool   `json:"ok"`
	Message string `json:"message"`
	Queued  bool   `json:"queued"`
}

func (s Stack) Status() string {
	if s.Framework == "skycap" {
		return skycapStatus[s.StatusCode]
	}
	return stackStatus[s.StatusCode]
}

func (s Stack) Namespace() string {
	return s.Namespaces[0]
}

func (s Stack) Health() string {
	return healthStatus[s.HealthCode]
}

func (c *Client) StackList() ([]Stack, error) {
	query_strings := make(map[string]string)
	query_strings["page"] = "1"

	var p Pagination
	var result []Stack
	var stacksRes []Stack

	for {
		req, err := c.NewRequest("GET", "/stacks.json", nil, query_strings)
		if err != nil {
			return nil, err
		}

		stacksRes = nil
		err = c.DoReq(req, &stacksRes, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, stacksRes...)
		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}
	}
	return result, nil
}

func (c *Client) StackListWithFilter(filter filterFunction) ([]Stack, error) {
	query_strings := make(map[string]string)
	query_strings["page"] = "1"

	var p Pagination
	var mid_result []Stack
	var stacksRes []Stack

	for {
		req, err := c.NewRequest("GET", "/stacks.json", nil, query_strings)
		if err != nil {
			return nil, err
		}

		stacksRes = nil
		err = c.DoReq(req, &stacksRes, &p)
		if err != nil {
			return nil, err
		}

		mid_result = append(mid_result, stacksRes...)
		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}

	}

	var result []Stack
	for _, item := range mid_result {
		if filter(item) {
			result = append(result, item)
		}
	}
	return result, nil
}

func (c *Client) CreateStack(name, environment, serviceYaml, manifestYaml string, targetOptions map[string]string) (*AsyncResult, error) {
	params := struct {
		Name         string `json:"name"`
		Environment  string `json:"environment"`
		ServiceYaml  string `json:"service_yaml"`
		ManifestYaml string `json:"manifest_yaml"`
		Cloud        string `json:"cloud"`
		KeyName      string `json:"key_name"`
		Region       string `json:"region"`
		Size         string `json:"size"`
		BuildType    string `json:"build_type"`
	}{
		Name:         name,
		Environment:  environment,
		ServiceYaml:  serviceYaml,
		ManifestYaml: manifestYaml,
		Cloud:        targetOptions["cloud"],
		KeyName:      targetOptions["key_name"],
		Region:       targetOptions["region"],
		Size:         targetOptions["size"],
		BuildType:    targetOptions["build_type"],
	}
	req, err := c.NewRequest("POST", "/stacks", params, nil)
	if err != nil {
		return nil, err
	}
	var asyncResult *AsyncResult
	return asyncResult, c.DoReq(req, &asyncResult, nil)
}

func (c *Client) StackInfo(stackName string) (*Stack, error) {
	stack, err := c.FindStackByName(stackName, "")
	if err != nil {
		return nil, err
	}
	return c.FindStackByUid(stack.Uid)
}

func (c *Client) StackInfoWithEnvironment(stackName, environment string) (*Stack, error) {
	stack, err := c.FindStackByName(stackName, environment)
	if err != nil {
		return nil, err
	}
	return c.FindStackByUid(stack.Uid)
}

func (c *Client) FindStackByUid(stackUid string) (*Stack, error) {
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+".json", nil, nil)
	if err != nil {
		return nil, err
	}
	var stacksRes *Stack
	return stacksRes, c.DoReq(req, &stacksRes, nil)
}

func (c *Client) StackSettings(uid string) ([]StackSetting, error) {
	query_strings := make(map[string]string)
	query_strings["page"] = "1"

	var p Pagination
	var result []StackSetting
	var settingsRes []StackSetting

	for {
		req, err := c.NewRequest("GET", "/stacks/"+uid+"/settings.json", nil, query_strings)
		if err != nil {
			return nil, err
		}

		settingsRes = nil
		err = c.DoReq(req, &settingsRes, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, settingsRes...)
		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}

	}

	return result, nil
}

func (c *Client) StackEnvVars(uid string) ([]StackEnvVar, error) {

	query_strings := make(map[string]string)
	query_strings["page"] = "1"

	var p Pagination
	var result []StackEnvVar
	var envVarsRes []StackEnvVar

	for {
		req, err := c.NewRequest("GET", "/stacks/"+uid+"/environments.json", nil, query_strings)
		if err != nil {
			return nil, err
		}

		envVarsRes = nil
		err = c.DoReq(req, &envVarsRes, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, envVarsRes...)
		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}

	}

	return result, nil
}

func (c *Client) StackEnvVarNew(stackUid string, key string, value string) (*AsyncResult, error) {
	params := struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{
		Key:   key,
		Value: value,
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/environments.json", params, nil)
	if err != nil {
		return nil, err
	}
	var asyncResult *AsyncResult
	return asyncResult, c.DoReq(req, &asyncResult, nil)
}

func (c *Client) StackEnvVarSet(stackUid string, key string, value string) (*AsyncResult, error) {
	params := struct {
		Value string `json:"value"`
	}{
		Value: value,
	}
	req, err := c.NewRequest("PUT", "/stacks/"+stackUid+"/environments/"+key+".json", params, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}

func (c *Client) FindStackByName(stackName, environment string) (*Stack, error) {
	stacks, err := c.StackList()

	for _, b := range stacks {
		if (strings.ToLower(b.Name) == strings.ToLower(stackName)) && (environment == "" || b.Environment == "" || environment == b.Environment) {
			return &b, err
		}
	}

	return nil, errors.New("Stack not found")
}

func (c *Client) ManagedBackups(uid string) ([]ManagedBackup, error) {
	query_strings := make(map[string]string)
	query_strings["page"] = "1"

	var p Pagination
	var result []ManagedBackup
	var managedBackupsRes []ManagedBackup

	for {
		req, err := c.NewRequest("GET", "/stacks/"+uid+"/backups.json", nil, query_strings)
		if err != nil {
			return nil, err
		}

		managedBackupsRes = nil
		err = c.DoReq(req, &managedBackupsRes, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, managedBackupsRes...)
		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}

	}

	return result, nil
}

func (c *Client) Set(uid string, key string, value string) (*AsyncResult, error) {
	key = strings.Replace(key, ".", "-", -1)
	params := struct {
		Value string `json:"value"`
	}{
		Value: value,
	}
	req, err := c.NewRequest("PUT", "/stacks/"+uid+"/settings/"+key+".json", params, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}

func (c *Client) Lease(uid string, ipAddress *string, timeToOpen *int, port *int, serverUid *string) (*AsyncResult, error) {
	var (
		theIpAddress  *string
		theTimeToOpen *int
		thePort       *int
		theServerUid  *string
	)
	// set defaults
	if ipAddress == nil {
		var value = "AUTO"
		theIpAddress = &value
	} else {
		theIpAddress = ipAddress
	}
	if timeToOpen == nil {
		var value = 20
		theTimeToOpen = &value
	} else {
		theTimeToOpen = timeToOpen
	}
	if port == nil {
		var value = 22
		thePort = &value
	} else {
		thePort = port
	}
	if serverUid == nil {
		var value = ""
		theServerUid = &value
	} else {
		theServerUid = serverUid
	}

	params := struct {
		TimeToOpen *int    `json:"ttl"`
		IpAddress  *string `json:"from_ip"`
		Port       *int    `json:"port"`
		ServerUid  *string `json:"server_id"`
	}{
		TimeToOpen: theTimeToOpen,
		IpAddress:  theIpAddress,
		Port:       thePort,
		ServerUid:  theServerUid,
	}
	req, err := c.NewRequest("POST", "/stacks/"+uid+"/firewalls.json", params, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}

func (c *Client) LeaseSync(stackUid string, ipAddress *string, timeToOpen *int, port *int, serverUid *string) (*GenericResponse, error) {
	asyncRes, err := c.Lease(stackUid, ipAddress, timeToOpen, port, serverUid)
	if err != nil {
		return nil, err
	}
	genericRes, err := c.WaitStackAsyncAction(asyncRes.Id, stackUid, 2*time.Second, 5*time.Minute, false)
	if err != nil {
		return nil, err
	}
	return genericRes, err
}

func (c *Client) RedeployStack(stackUid string, gitRef string, services []string) (*RedeployResponse, error) {
	params := struct {
		GitRef   string   `json:"git_ref"`
		Services []string `json:"services"`
	}{
		GitRef:   gitRef,
		Services: services,
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/deployments.json", params, nil)
	if err != nil {
		return nil, err
	}
	var redeployRes *RedeployResponse
	return redeployRes, c.DoReq(req, &redeployRes, nil)
}

func (c *Client) StackReboot(stackUid string, strategy string, group string) (*AsyncResult, error) {
	params := struct {
		Strategy string `json:"strategy"`
		Group    string `json:"group"`
	}{
		Strategy: strategy,
		Group:    group,
	}

	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/reboot_servers.json", params, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}

func (c *Client) InvokeStackAction(stackUid string, action string) (*AsyncResult, error) {
	params := struct {
		Command string `json:"command"`
	}{
		Command: action,
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/actions.json", params, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}

func (c *Client) InvokeDbStackAction(stackUid string, serverUid string, dbType *string, action string) (*AsyncResult, error) {
	var params interface{}
	if dbType == nil {
		params = struct {
			Command   string `json:"command"`
			ServerUid string `json:"server_uid"`
		}{
			Command:   action,
			ServerUid: serverUid,
		}
	} else {
		params = struct {
			Command   string `json:"command"`
			ServerUid string `json:"server_uid"`
			DbType    string `json:"db_type"`
		}{
			Command:   action,
			ServerUid: serverUid,
			DbType:    *dbType,
		}
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/actions.json", params, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}
