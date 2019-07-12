package cloud66

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	DefaultCheckFrequency = 15 * time.Second // every 10 seconds
	DefaultTimeout        = 10 * time.Minute // 10 minutes
)

type AsyncResult struct {
	Id              int        `json:"id"`
	User            string     `json:"user"`
	ResourceType    string     `json:"resource_type"`
	ResourceId      string     `json:"resource_id"`
	Action          string     `json:"action"`
	StartedVia      string     `json:"started_via"`
	StartedAt       time.Time  `json:"started_at"`
	FinishedAt      *time.Time `json:"finished_at"`
	FinishedSuccess *bool      `json:"finished_success"`
	FinishedMessage string     `json:"finished_message"`
}

func (c *Client) WaitStackAsyncAction(asyncId int, stackUid string, checkFrequency time.Duration, timeout time.Duration, showWorkingIndicator bool) (*GenericResponse, error) {
	var timeoutTime = time.Now().Add(timeout)

	// declare vars
	var (
		asyncRes *AsyncResult
		err      error
	)

	for {
		// fetch the current status of the async action
		asyncRes, err = c.getStackAsyncAction(asyncId, stackUid)
		if err != nil {
			return nil, err
		}
		// check for a result!
		if asyncRes.FinishedAt != nil {
			break
		}
		// check for client-side time-out
		if time.Now().After(timeoutTime) {
			return nil, errors.New("timed-out after " + strconv.FormatInt(int64(timeout)/int64(time.Second), 10) + " second(s)")
		}
		// sleep for checkFrequency secs between lookup requests
		time.Sleep(checkFrequency)
		if showWorkingIndicator {
			fmt.Printf(".")
		}
	}
	// response
	genericRes := GenericResponse{
		Status:  *asyncRes.FinishedSuccess,
		Message: asyncRes.FinishedMessage,
	}
	return &genericRes, err
}

func (c *Client) getStackAsyncAction(asyncId int, stackUid string) (*AsyncResult, error) {
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/actions/"+strconv.Itoa(asyncId)+".json", nil, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}
