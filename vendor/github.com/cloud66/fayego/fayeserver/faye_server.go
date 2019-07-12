/*
	Faye Server

*/
package fayeserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/pborman/uuid"
)

const CHANNEL_HANDSHAKE = "/meta/handshake"
const CHANNEL_CONNECT = "/meta/connect"
const CHANNEL_DISCONNECT = "/meta/disconnect"
const CHANNEL_SUBSCRIBE = "/meta/subscribe"
const CHANNEL_UNSUBSCRIBE = "/meta/unsubscribe"

type FayeServer struct {
	Connections   []Connection
	Subscriptions map[string][]Client
	SubMutex      sync.RWMutex
	Clients       map[string]Client
	ClientMutex   sync.RWMutex
	idCount       int
}

/*
Instantiate a new faye server
*/
func NewFayeServer() *FayeServer {
	return &FayeServer{Connections: []Connection{},
		Subscriptions: make(map[string][]Client),
		Clients:       make(map[string]Client)}
}

// general message handling
/*

 */
func (f *FayeServer) publishToChannel(channel, data string) {
	subs, ok := f.Subscriptions[channel]
	fmt.Println("Subs: ", f.Subscriptions, "count: ", len(f.Subscriptions[channel]))
	if ok {
		f.multiplexWrite(subs, data)
	}
}

/*

 */
func (f *FayeServer) multiplexWrite(subs []Client, data string) {
	var group sync.WaitGroup
	for i := range subs {
		fmt.Println("subs[i]: ", subs[i])
		group.Add(1)
		go func(client chan<- []byte, data string) {
			if client != nil {
				fmt.Println("WRITE FOR CLIENT")
				client <- []byte(data)
			} else {
				fmt.Println("NO CHANNEL DON'T TRY TO WRITE")
			}
			group.Done()
		}(subs[i].WriteChannel, data)
	}
	group.Wait()
}

func (f *FayeServer) findClientForChannel(c chan []byte) *Client {
	f.ClientMutex.Lock()
	defer f.ClientMutex.Unlock()

	for _, client := range f.Clients {
		if client.WriteChannel == c {
			fmt.Println("Matched Client: ", client.ClientId)
			return &client
		}
	}
	return nil
}

func (f *FayeServer) DisconnectChannel(c chan []byte) {
	client := f.findClientForChannel(c)
	if client != nil {
		fmt.Println("Disconnect Client: ", client.ClientId)
		f.removeClientFromServer(client.ClientId)
	}
}

// ========

type FayeMessage struct {
	Channel                  string      `json:"channel"`
	ClientId                 string      `json:"clientId,omitempty"`
	Subscription             string      `json:"subscription,omitempty"`
	Data                     interface{} `json:"data,omitempty"`
	Id                       string      `json:"id,omitempty"`
	SupportedConnectionTypes []string    `json:"supportedConnectionTypes,omitempty"`
}

// Message handling

func (f *FayeServer) HandleMessage(message []byte, c chan []byte) ([]byte, error) {
	// parse message JSON
	fm := FayeMessage{}
	err := json.Unmarshal(message, &fm)

	if err != nil {
		fmt.Println("Error parsing message json, try array parse:", err)

		ar := []FayeMessage{}
		jerr := json.Unmarshal(message, &ar)
		if jerr != nil {
			fmt.Println("Error parsing message json as array:", err)
		} else {
			fm = ar[0]
			fmt.Println("Parsed as: ", fm)
		}
	}

	switch fm.Channel {
	case CHANNEL_HANDSHAKE:
		fmt.Println("handshake")
		return f.handshake()
	case CHANNEL_CONNECT:
		fmt.Println("connect")
		return f.connect(fm.ClientId)
	case CHANNEL_DISCONNECT:
		fmt.Println("disconnect")
		return f.disconnect(fm.ClientId)
	case CHANNEL_SUBSCRIBE:
		fmt.Println("subscribe")
		return f.subscribe(fm.ClientId, fm.Subscription, c)
	case CHANNEL_UNSUBSCRIBE:
		fmt.Println("subscribe")
		return f.unsubscribe(fm.ClientId, fm.Subscription)
	default:
		fmt.Println("publish")
		fmt.Println("data is: ", fm.Data)
		return f.publish(fm.Channel, fm.Id, fm.Data)
	}
}

/*
FayeResponse
*/

type FayeResponse struct {
	Channel                  string                 `json:"channel,omitempty"`
	Successful               bool                   `json:"successful,omitempty"`
	Version                  string                 `json:"version,omitempty"`
	SupportedConnectionTypes []string               `json:"supportedConnectionTypes,omitempty"`
	ConnectionType           string                 `json:"connectionType,omitempty"`
	ClientId                 string                 `json:"clientId,omitempty"`
	Advice                   map[string]interface{} `json:"advice,omitempty"`
	Subscription             string                 `json:"subscription,omitempty"`
	Error                    string                 `json:"error,omitempty"`
	Id                       string                 `json:"id,omitempty"`
	Data                     interface{}            `json:"data,omitempty"`
	Ext                      interface{}            `json:"ext,omitempty"`
}

/*

Handshake:

Example response:
{
    "channel": "/meta/handshake",
    "successful": true,
    "version": "1.0",
    "supportedConnectionTypes": [
        "long-polling",
        "cross-origin-long-polling",
        "callback-polling",
        "websocket",
        "eventsource",
        "in-process"
    ],
    "clientId": "1fg1b9s10zm29e0ahpk490mzkqk3",
    "advice": {
        "reconnect": "retry",
        "interval": 0,
        "timeout": 45000
    }
}

Bayeux Handshake response

*/

func (f *FayeServer) handshake() ([]byte, error) {
	fmt.Println("handshake!")

	// build response
	resp := FayeResponse{
		Id:                       "1",
		Channel:                  "/meta/handshake",
		Successful:               true,
		Version:                  "1.0",
		SupportedConnectionTypes: []string{"websocket", "callback-polling", "long-polling", "cross-origin-long-polling", "eventsource", "in-process"},
		ClientId:                 generateClientId(),
		Advice:                   map[string]interface{}{"reconnect": "retry", "interval": 0, "timeout": 45000},
	}

	// wrap it in an array & convert to json
	return json.Marshal([]FayeResponse{resp})
}

/*

Connect:

Example response
[
  {
     "channel": "/meta/connect",
     "successful": true,
     "error": "",
     "clientId": "Un1q31d3nt1f13r",
     "timestamp": "12:00:00 1970",
     "advice": { "reconnect": "retry" }
   }
]
*/

func (f *FayeServer) connect(clientId string) ([]byte, error) {
	// TODO: setup client connection state

	resp := FayeResponse{
		Channel:    "/meta/connect",
		Successful: true,
		Error:      "",
		ClientId:   clientId,
		Advice:     map[string]interface{}{"reconnect": "retry"},
	}

	// wrap it in an array & convert to json
	return json.Marshal([]FayeResponse{resp})
}

/*
Disconnect

Example response
[
  {
     "channel": "/meta/disconnect",
     "clientId": "Un1q31d3nt1f13r"
     "successful": true
  }
]
*/

func (f *FayeServer) disconnect(clientId string) ([]byte, error) {
	// tear down client connection state
	f.removeClientFromServer(clientId)

	resp := FayeResponse{
		Channel:    "/meta/disconnect",
		Successful: true,
		ClientId:   clientId,
	}

	// wrap it in an array & convert to json
	return json.Marshal([]FayeResponse{resp})
}

/*
Subscribe

Example response
[
  {
     "channel": "/meta/subscribe",
     "clientId": "Un1q31d3nt1f13r",
     "subscription": "/foo/**",
     "successful": true,
     "error": ""
   }
]
*/

func (f *FayeServer) subscribe(clientId, subscription string, c chan []byte) ([]byte, error) {

	// subscribe the client to the given channel
	if len(subscription) == 0 {
		return []byte{}, errors.New("Subscription channel not present")
	}

	f.addClientToSubscription(clientId, subscription, c)

	// if successful send success response
	resp := FayeResponse{
		Channel:      "/meta/subscribe",
		ClientId:     clientId,
		Subscription: subscription,
		Successful:   true,
		Error:        "",
	}

	// TODO: handle failure case

	// wrap it in an array and convert to json
	return json.Marshal([]FayeResponse{resp})
}

/*
Unsubscribe

Example response
[
  {
     "channel": "/meta/unsubscribe",
     "clientId": "Un1q31d3nt1f13r",
     "subscription": "/foo/**",
     "successful": true,
     "error": ""
   }
]
*/

func (f *FayeServer) unsubscribe(clientId, subscription string) ([]byte, error) {
	// TODO: unsubscribe the client from the given channel
	if len(subscription) == 0 {
		return []byte{}, errors.New("Subscription channel not present")
	}

	// remove the client as a subscriber on the channel
	if f.removeClientFromSubscription(clientId, subscription) {
		fmt.Println("Successful unsubscribe")
	} else {
		fmt.Println("Failed to unsubscribe")
	}

	// if successful send success response
	resp := FayeResponse{
		Channel:      "/meta/unsubscribe",
		ClientId:     clientId,
		Subscription: subscription,
		Successful:   true,
		Error:        "",
	}

	// TODO: handle failure case

	// wrap it in an array and convert to json
	return json.Marshal([]FayeResponse{resp})
}

/*
Publish

Example response
[
  {
     "channel": "/some/channel",
     "successful": true,
     "id": "some unique message id"
  }
]

*/
func (f *FayeServer) publish(channel, id string, data interface{}) ([]byte, error) {

	//convert data back to json string
	message := FayeResponse{
		Channel: channel,
		Id:      id,
		Data:    data,
	}

	dataStr, err := json.Marshal([]FayeResponse{message})
	if err != nil {
		fmt.Println("Error parsing message!")
		return []byte{}, errors.New("Invalid Message Data")
	}
	fmt.Println("publish to: ", channel)
	fmt.Println("data: ", string(dataStr))

	f.publishToChannel(channel, string(dataStr))

	resp := FayeResponse{
		Channel:    channel,
		Successful: true,
		Id:         id,
	}

	return json.Marshal([]FayeResponse{resp})
}

// Helper functions:

/*
	Generate a clientId for use in the communication with the client
*/
func generateClientId() string {
	return uuid.New()
}

func (f *FayeServer) nextMessageId() string {
	f.idCount++
	return string(f.idCount)
}
