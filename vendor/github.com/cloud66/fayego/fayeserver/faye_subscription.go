package fayeserver

import (
	"errors"
	"fmt"
)

/*
Client

Clients represent connected Faye Clients, each has an Id negotiated during handshake, a write channel tied to their network connection
and a list of subscriptions(faye channels) that have been subscribed to by the client.

*/
type Client struct {
	ClientId     string
	WriteChannel chan []byte
	ClientSubs   []string
}

func (c *Client) isSubscribed(sub string) bool {
	for _, clientSub := range c.ClientSubs {
		if clientSub == sub {
			return true
		}
	}
	return false
}

/*
subscription management
*/
func (f *FayeServer) subscriptionClientIndex(subscriptions []Client, clientId string) int {
	for i, c := range subscriptions {
		if c.ClientId == clientId {
			return i
		}
	}
	return -1
}

func (f *FayeServer) removeSubFromClient(client Client, sub string) Client {
	for i, clientSub := range client.ClientSubs {
		if clientSub == sub {
			client.ClientSubs = append(client.ClientSubs[:i], client.ClientSubs[i+1:]...)
			return client
		}
	}
	return client
}

func (f *FayeServer) removeClientFromSubscription(clientId, subscription string) bool {
	fmt.Println("Remove Client From Subscription: ", subscription)

	// grab the client subscriptions array for the channel
	f.SubMutex.Lock()
	defer f.SubMutex.Unlock()

	subs, ok := f.Subscriptions[subscription]

	if !ok {
		return false
	}

	index := f.subscriptionClientIndex(subs, clientId)

	if index >= 0 {
		f.Subscriptions[subscription] = append(subs[:index], subs[index+1:]...)
	} else {
		return false
	}

	// remove sub from client subs list
	f.Clients[clientId] = f.removeSubFromClient(f.Clients[clientId], subscription)

	return true
}

func (f *FayeServer) addClientToSubscription(clientId, subscription string, c chan []byte) bool {
	fmt.Println("Add Client to Subscription: ", subscription)

	// Add client to server list if it is not present
	client := f.addClientToServer(clientId, subscription, c)

	// add the client as a subscriber to the channel if it is not already one
	f.SubMutex.Lock()
	defer f.SubMutex.Unlock()
	subs, cok := f.Subscriptions[subscription]
	if !cok {
		f.Subscriptions[subscription] = []Client{}
	}

	index := f.subscriptionClientIndex(subs, clientId)

	fmt.Println("Subs: ", f.Subscriptions, "count: ", len(f.Subscriptions[subscription]))

	if index < 0 {
		f.Subscriptions[subscription] = append(subs, *client)
		return true
	}

	return false
}

// client management

/*
updateClientChannel
*/
func (f *FayeServer) UpdateClientChannel(clientId string, c chan []byte) bool {
	fmt.Println("update client for channel: clientId: ", clientId)
	f.ClientMutex.Lock()
	defer f.ClientMutex.Unlock()
	client, ok := f.Clients[clientId]
	if !ok {
		client = Client{clientId, c, []string{}}
		f.Clients[clientId] = client
		return true
	}

	client.WriteChannel = c
	f.Clients[clientId] = client
	fmt.Println("Worked")

	return true
}

/*
Add Client to server only if the client is not already present
*/
func (f *FayeServer) addClientToServer(clientId, subscription string, c chan []byte) *Client {
	fmt.Println("Add client: ", clientId)

	f.ClientMutex.Lock()
	defer f.ClientMutex.Unlock()
	client, ok := f.Clients[clientId]
	if !ok {
		client = Client{clientId, c, []string{}}
		f.Clients[clientId] = client
	}

	fmt.Println("Client subs: ", len(client.ClientSubs), " | ", client.ClientSubs)

	// add the subscription to the client subs list
	if !client.isSubscribed(subscription) {
		fmt.Println("Client not subscribed")
		client.ClientSubs = append(client.ClientSubs, subscription)
		f.Clients[clientId] = client
		fmt.Println("Client sub count: ", len(client.ClientSubs))
	} else {
		fmt.Println("Client already subscribed")
	}

	return &client
}

/*
Remove the Client from the server and unsubscribe from any subscriptions
*/
func (f *FayeServer) removeClientFromServer(clientId string) error {
	fmt.Println("Remove client: ", clientId)

	f.ClientMutex.Lock()
	defer f.ClientMutex.Unlock()

	client, ok := f.Clients[clientId]
	if !ok {
		return errors.New("Error removing client")
	}

	// clear any subscriptions
	for _, sub := range client.ClientSubs {
		fmt.Println("Remove sub: ", sub)
		if f.removeClientFromSubscription(client.ClientId, sub) {
			fmt.Println("Removed sub!")
		} else {
			fmt.Println("Failed to remove sub.")
		}
	}

	// remove the client from the server
	delete(f.Clients, clientId)

	return nil
}
