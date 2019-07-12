package cloud66

import (
	"errors"
	"fmt"

	"github.com/cloud66/fayego/fayeclient"
)

type MessageCallback func(clientMessage fayeclient.ClientMessage)

func NewFayeClient(fayeServerUrl string) (*fayeclient.FayeClient, error) {
	// instantiate client
	client := fayeclient.NewFayeClient(fayeServerUrl)

	ready := make(chan bool)
	err := client.Start(ready)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error starting client: %s", err))
	}

	// ready will recieve true when the client is connected
	<-ready

	return client, err
}

func RegisterCallback(client *fayeclient.FayeClient, channel string, callbacks ...interface{}) (*fayeclient.FayeClient, error) {

	var (
		successCallback MessageCallback
		errorCallback   *MessageCallback
	)

	// validate
	if len(callbacks) < 1 || len(callbacks) > 2 {
		return nil, errors.New("Expect at least one callback for the success messages; with an optional second callback for errors")
	} else if len(callbacks) == 2 {
		successCallback = callbacks[0].(MessageCallback)
		errorCallback = callbacks[1].(*MessageCallback)

	} else {
		successCallback = callbacks[0].(MessageCallback)
	}

	// subscribe to a channel
	client.Subscribe(channel)

	// read from stdin
	go recvMessages(client, successCallback, errorCallback)

	// no errors
	return client, nil
}

/*
Listen for messages from the client's message channel and print them to stdout
*/
func recvMessages(client *fayeclient.FayeClient, successCallback MessageCallback, errorCallback *MessageCallback) {
	for {
		select {
		case clientMessage, ok := <-client.MessageChan:
			if ok {
				successCallback(clientMessage)
			} else {
				if errorCallback != nil {
					errCbk := *errorCallback
					errCbk(clientMessage)
				}
			}
		}
	}
}
