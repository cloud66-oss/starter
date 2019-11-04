package cloud66

import (
	"fmt"
	"strconv"
	"time"
)

type StackProfileType struct {
	StackUid string `json:"stack_uid"`
	Role     string `json:"role"`
}

type AccountProfileType struct {
	CanCreateStack              bool     `json:"can_create_stack"`
	CanAdminUsers               bool     `json:"can_admin_users"`
	CanAdminPayments            bool     `json:"can_payment"`
	CanAddCloudKey              bool     `json:"can_add_cloud_key"`
	CanDelCloudKey              bool     `json:"can_del_cloud_key"`
	CanViewAccountNotifications bool     `json:"can_view_acc_notifications"`
	CanEditAccountNotifications bool     `json:"can_edit_acc_notifications"`
	CanViewAudit                bool     `json:"can_view_audit"`
	CanViewDockerImageKey       bool     `json:"can_view_docker_img_key"`
	CanDelSshKey                bool     `json:"can_del_ssh_key"`
	CanEditPersonalToken        bool     `json:"can_edit_personal_token"`
	CanDelAuthorizedApp         bool     `json:"can_del_authorized_app"`
	CanViewCustomEnv            bool     `json:"can_view_custom_env"`
	CanEditCustomEnv            bool     `json:"can_edit_custom_env"`
	CanAddDevelopersApp         bool     `json:"can_add_developers_app"`
	CanDelDevelopersAdd         bool     `json:"can_del_developers_app"`
	CanEditGitKey               bool     `json:"can_edit_git_key"`
	CanEditGateway              bool     `json:"can_edit_gateway"`
	DefaultRoles                []string `json:"default_roles"`
}

type AclsType struct {
	AccountId int    `json:"account_id"`
	EntityUrl string `json:"entity_uri"`
	Action    string `json:"action"`
}

type AccessProfileType struct {
	AccountProfile AccountProfileType `json:"account_profile"`
	StackProfiles  []StackProfileType `json:"stack_profiles"`
	AclsProfile    []AclsType         `json:"acls_profile"`
	Override       bool               `json:"override"`
}

type User struct {
	Id               int               `json:"id"`
	Email            string            `json:"email"`
	PrimaryAccountId int               `json:"primary_account_id"`
	Locked           bool              `json:"locked"`
	AccessProfile    AccessProfileType `json:"access_profile"`
	UsesTfa          bool              `json:"uses_tfa"`
	Timezone         string            `json:"timezone"`
	HasValidPhone    bool              `json:"has_valid_phone"`
	DeveloperProgram bool              `json:"developer_program"`
	GithubLogin      bool              `json:"github_login"`
	LastLogin        time.Time         `json:"last_login"`
	Devices          interface{}       `json:"devices"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	CloudStatus      string            `json:"cloud_status"`
}

func (c *Client) ListUsers() ([]User, error) {
	query_strings := make(map[string]string)
	query_strings["page"] = "1"

	var p Pagination
	var result []User
	var userRes []User

	for {
		req, err := c.NewRequest("GET", "/users.json", nil, query_strings)
		if err != nil {
			return nil, err
		}

		userRes = nil
		err = c.DoReq(req, &userRes, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, userRes...)
		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}
	}

	return result, nil
}

func (c *Client) GetUser(userId int) (*User, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("/users/%d.json", userId), nil, nil)
	if err != nil {
		return nil, err
	}

	var userRes *User
	return userRes, c.DoReq(req, &userRes, nil)
}

func (c *Client) UpdateUser(userId int, user User) (*User, error) {
	req, err := c.NewRequest("PUT", fmt.Sprintf("/users/%d.json", userId), user, nil)
	if err != nil {
		return nil, err
	}
	var userRes *User
	return userRes, c.DoReq(req, &userRes, nil)
}
