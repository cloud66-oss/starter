package fayeclient

/*

TODO:

* handle extensions
* implement other protocol comm functions

*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloud66/fayego/fayeserver"
	"github.com/gorilla/websocket"
	"net"
	"net/url"
	"time"
)

const DEFAULT_HOST = "localhost:4001/faye"
const DEFAULT_KEEP_ALIVE_SECS = 30

const ( // iota is reset to 0
	StateWSDisconnected   = iota // == 0
	StateWSConnected      = iota
	StateFayeDisconnected = iota
	StateFayeConnected    = iota
)

// ==========================================
/*
A subscription represents a subscription to a channel by the client
Each sub has a path representing the channel on faye, a messageChan which is recieve any messages sent to it and a connected indicator to indicate the state of the sub on the faye server
*/
type ClientSubscription struct {
	channel   string
	connected bool
}

func (f *FayeClient) addSubscription(channel string) {
	c := ClientSubscription{channel: channel, connected: false}
	f.subscriptions = append(f.subscriptions, &c)
}

func (f *FayeClient) removeSubscription(channel string) {
	for i, sub := range f.subscriptions {
		if channel == sub.channel {
			f.subscriptions = append(f.subscriptions[:i], f.subscriptions[i+1:]...)
		}
	}
}

func (f *FayeClient) updateSubscription(channel string, connected bool) {
	s := f.getSubscription(channel)
	s.connected = connected
}

func (f *FayeClient) getSubscription(channel string) *ClientSubscription {
	for i, sub := range f.subscriptions {
		if channel == sub.channel {
			return f.subscriptions[i]
		}
	}
	return nil
}

func (f *FayeClient) resubscribeSubscriptions() {
	for _, sub := range f.subscriptions {
		fmt.Println("resubscribe: ", sub.channel)
		f.subscribe(sub.channel)
	}
}

// ==========================================
type FayeClient struct {
	Host          string
	MessageChan   chan ClientMessage // any messages recv'd by the client will be sent to the message channel - TODO: remap this to a set of subscription message channels one per active subscription
	conn          *Connection
	fayeState     int
	readyChan     chan bool
	clientId      string
	messageNumber int
	subscriptions []*ClientSubscription
	keepAliveSecs int
	keepAliveChan chan bool
}

type ClientMessage struct {
	Channel string
	Data    map[string]interface{}
	Ext     map[string]interface{}
}

func NewFayeClient(host string) *FayeClient {
	if len(host) == 0 {
		host = DEFAULT_HOST
	}
	// instantiate a FayeClient and return
	return &FayeClient{Host: host, fayeState: StateWSDisconnected, MessageChan: make(chan ClientMessage, 100), messageNumber: 0, keepAliveSecs: DEFAULT_KEEP_ALIVE_SECS, keepAliveChan: make(chan bool)}
}

func (f *FayeClient) SetKeepAliveIntervalSeconds(secs int) {
	f.keepAliveSecs = secs
}

func (f *FayeClient) Start(ready chan bool) error {
	fmt.Println("Starting...")
	err := f.connectToServer()
	if err != nil {
		return err
	}

	// kick off the connection handshake
	f.readyChan = ready
	f.handshake()
	return nil
}

/*
Open the websocket connection to the faye server and initialize the client state
*/
func (f *FayeClient) connectToServer() error {
	fmt.Println("start client")
	fmt.Println("connectToServer")

	url, _ := url.Parse("ws://" + f.Host)
	c, err := net.Dial("tcp", url.Host)

	if err != nil {
		fmt.Println("Error connecting to server: ", err)
		return err
	}

	ws, resp, err := websocket.NewClient(c, url, nil, 1024, 1024)

	if err != nil {
		return err
	}

	f.fayeState = StateWSConnected

	if resp != nil {
		fmt.Println("Resp: ", resp)
	}

	conn := NewConnection(ws)
	f.conn = conn
	f.conn.writerConnected = true
	f.conn.readerConnected = true
	go conn.writer()
	go conn.reader(f)

	// close keep alive channel to stop any running keep alive
	close(f.keepAliveChan)
	f.keepAliveChan = make(chan bool)
	go f.keepAlive()
	return nil
}

func (f *FayeClient) keepAlive() {
	fmt.Println("START KEEP ALIVE")
	c := time.Tick(time.Duration(f.keepAliveSecs) * time.Second)
	for {
		select {
		case _, ok := <-f.keepAliveChan:
			if !ok {
				fmt.Println("exit keep alive")
				return
			}
		case <-c:
			fmt.Println("Send keep-alive: ", time.Now())
			f.connect()
		}

	}
	fmt.Println("exiting keepalive func")
}

/*
Close the websocket connection and set the faye client state
*/
func (f *FayeClient) disconnectFromServer() {
	fmt.Println("DISCONNECT FROM SERVER")
	f.fayeState = StateWSDisconnected
	f.conn.exit <- true
	f.conn.ws.Close()
}

/*
ReaderDisconnect - called by the connection handler if the reader connection is dropped by the loss of a server connection
*/
func (f *FayeClient) ReaderDisconnect() {
	f.readyChan <- false
}

/*
Write a message to the faye server over the websocket connection
*/
func (f *FayeClient) Write(msg string) error {
	f.conn.send <- []byte(msg)
	return nil
}

/*
Parse and interpret a faye message response
*/
func (f *FayeClient) HandleMessage(message []byte) error {
	// parse the faye message and interpret the logic to set client state appropriately
	resp := []fayeserver.FayeResponse{}
	err := json.Unmarshal(message, &resp)
	var fm fayeserver.FayeResponse

	if err != nil {
		fmt.Println("Error parsing json. ", err)
	}

	for i := range resp {
		fm = resp[i]
		switch fm.Channel {
		case fayeserver.CHANNEL_HANDSHAKE:
			f.clientId = fm.ClientId
			f.connect() // send faye connect message
			f.fayeState = StateFayeConnected
			f.readyChan <- true

		case fayeserver.CHANNEL_CONNECT:
			//fmt.Println("Recv'd connect response")

		case fayeserver.CHANNEL_DISCONNECT:
			f.fayeState = StateFayeDisconnected
			f.disconnectFromServer()

		case fayeserver.CHANNEL_SUBSCRIBE:
			f.updateSubscription(fm.Subscription, fm.Successful)

		case fayeserver.CHANNEL_UNSUBSCRIBE:
			if fm.Successful {
				f.removeSubscription(fm.Subscription)
			}
		default:
			if fm.Data != nil {
				if fm.ClientId == f.clientId {
					return nil
				}
				var data map[string]interface{}
				var ext map[string]interface{}

				if fm.Data != nil {
					data = fm.Data.(map[string]interface{})
				}

				if fm.Ext != nil {
					ext = fm.Ext.(map[string]interface{})
				}

				// tell the client we got a message on a channel
				go func(d, e map[string]interface{}) {
					select {
					case f.MessageChan <- ClientMessage{Channel: fm.Channel, Data: d, Ext: e}:
						return
					case <-time.After(100 * time.Millisecond):
						return
					}
				}(data, ext)
			}
		}
	}

	return nil
}

func (f *FayeClient) Subscribe(channel string) error {
	if len(channel) == 0 {
		return errors.New("Channel must have a value.")
	}
	//fmt.Println("Subscribe to channel: ", channel)
	f.addSubscription(channel)
	return f.subscribe(channel)
}

func (f *FayeClient) Unsubscribe(channel string) error {
	if len(channel) == 0 {
		return errors.New("Channel must have a value.")
	}
	//fmt.Println("Unsubscribe from channel: ", channel)
	return f.unsubscribe(channel)
}

func (f *FayeClient) Publish(channel string, data map[string]interface{}) error {
	return f.publish(channel, data)
}

func (f *FayeClient) Disconnect() {
	f.disconnect()
}

/*
Faye protocol messages
*/

/*
type FayeResponse struct {
	Channel                  string            `json:"channel,omitempty"`
	Successful               bool              `json:"successful,omitempty"`
	Version                  string            `json:"version,omitempty"`
	SupportedConnectionTypes []string          `json:"supportedConnectionTypes,omitempty"`
	ClientId                 string            `json:"clientId,omitempty"`
	Advice                   map[string]string `json:"advice,omitempty"`
	Subscription             string            `json:"subscription,omitempty"`
	Error                    string            `json:"error,omitempty"`
	Id                       string            `json:"id,omitempty"`
	Data                     interface{}       `json:"data,omitempty"`
}
*/

// Faye message functions

/*
 */
func (f *FayeClient) handshake() {
	message := fayeserver.FayeResponse{Channel: fayeserver.CHANNEL_HANDSHAKE, Version: "1.0", SupportedConnectionTypes: []string{"websocket"}}
	err := f.writeMessage(message)
	if err != nil {
		fmt.Println("Error generating handshake message")
	}
}

/*
Connect to Faye
*/
func (f *FayeClient) connect() {
	message := fayeserver.FayeResponse{Channel: fayeserver.CHANNEL_CONNECT, ClientId: f.clientId, ConnectionType: "websocket"}
	//fmt.Println("Connect message: ", message)
	err := f.writeMessage(message)
	if err != nil {
		fmt.Println("Error generating connect message")
	}
}

/*
Disconnect from Faye
*/
func (f *FayeClient) disconnect() {
	message := fayeserver.FayeResponse{Channel: fayeserver.CHANNEL_DISCONNECT, ClientId: f.clientId}
	//fmt.Println("Connect message: ", message)
	err := f.writeMessage(message)
	if err != nil {
		fmt.Println("Error generating connect message")
	}
}

/*
Subscribe the client to a channel
*/
func (f *FayeClient) subscribe(channel string) error {
	message := fayeserver.FayeResponse{Channel: fayeserver.CHANNEL_SUBSCRIBE, ClientId: f.clientId, Subscription: channel}
	//fmt.Println("Subscribe message: ", message)
	err := f.writeMessage(message)
	if err != nil {
		fmt.Println("Error generating subscribe message")
		return err
	}
	return nil
}

/*
Unsubscribe from a channel.
*/
func (f *FayeClient) unsubscribe(channel string) error {
	message := fayeserver.FayeResponse{Channel: fayeserver.CHANNEL_UNSUBSCRIBE, ClientId: f.clientId, Subscription: channel}
	//fmt.Println("Unsubscribe message: ", message)
	err := f.writeMessage(message)
	if err != nil {
		fmt.Println("Error generating unsubscribe message")
		return err
	}
	return nil
}

/*
  Publish a message to a channel.
*/
func (f *FayeClient) publish(channel string, data map[string]interface{}) error {
	message := fayeserver.FayeResponse{Channel: channel, ClientId: f.clientId, Id: f.messageId(), Data: data}
	//fmt.Println("publish message: ", message)
	err := f.writeMessage(message)
	if err != nil {
		fmt.Println("Error generating unsubscribe message")
		return err
	}
	return nil
}

// ------------------

/*
Encode the json and send the message over the wire.
*/
func (f *FayeClient) writeMessage(message fayeserver.FayeResponse) error {
	if !f.conn.Connected() {
		// reconnect
		fmt.Println("RECONNECT")
		cerr := f.connectToServer()
		if cerr != nil {
			return cerr
		}
		if !f.conn.Connected() {
			errors.New("Not Connected, Reconnect Failed.")
		}

	}

	json, err := json.Marshal(message)
	if err != nil {
		return err
	}
	f.Write(string(json))
	return nil
}

// Message Id
func (f *FayeClient) messageId() string {
	return "1"
}
