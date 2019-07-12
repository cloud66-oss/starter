/*
TODO: factor out the code and clean things up a little, commit - works with firefox/chrome/safari on mac with long-polling/eventsource support
TODO: pure long-polling
TODO: callback-polling
TODO: cross-origin-polling
*/
/*
Created by Paul Crawford
Copyright (c) 2013. All rights reserved.
*/
package fayeserver

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

// =====
// WebSocket handling

type Connection struct {
	ws *websocket.Conn
	//es          eventsource.EventSource
	send        chan []byte
	isWebsocket bool
}

/*
Initial constants based on websocket example code from github.com/garyburd/go-websocket
Reader & Writer functions also implemented based on
*/
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

func (c *Connection) esWriter(f *FayeServer) {
	fmt.Println("Writer started.")
	ticker := time.NewTicker(pingPeriod)
	//	defer func() {
	//		c.es.Close()
	//	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.esWrite([]byte{})
				return
			}
			if err := c.esWrite([]byte(message)); err != nil {
				return
			}
		case <-ticker.C:
			fmt.Println("tick.")
			if err := c.esWrite([]byte{}); err != nil {
				return
			}
		}
	}
}

/*
reader - reads messages from the websocket connection and passes them through to the fayeserver message handler
*/
func (c *Connection) reader(f *FayeServer) {
	fmt.Println("reading...")
	defer func() {
		fmt.Println("reader disconnect")
		f.DisconnectChannel(c.send)
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		// ask faye client to handle faye message
		response, ferr := f.HandleMessage(message, c.send)
		if ferr != nil {
			fmt.Println("Faye Error: ", ferr)
			c.send <- []byte(fmt.Sprintf("Error: ", ferr))
		} else {
			c.send <- response
		}
	}

	fmt.Println("reader exited.")
}

func (c *Connection) esWrite(payload []byte) error {
	fmt.Println("Writing to eventsource: ", string(payload))
	//c.es.SendMessage(string(payload), "", "")
	return nil
}

/*
write - writes messages to the websocket connection
*/
func (c *Connection) wsWrite(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

/*
writer - is the write loop that reads messages off the send channel and writes them out over the websocket connection
*/
func (c *Connection) writer(f *FayeServer) {
	fmt.Println("Writer started.")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.wsWrite(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.wsWrite(websocket.TextMessage, []byte(message)); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.wsWrite(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

/*
from faye js:

handle: function(request, response) {
    var requestUrl    = url.parse(request.url, true),
        requestMethod = request.method,
        origin        = request.headers.origin,
        self          = this;

    request.originalUrl = request.url;

    request.on('error', function(error) { self._returnError(response, error) });
    response.on('error', function(error) { self._returnError(null, error) });

    if (this._static.test(requestUrl.pathname))
      return this._static.call(request, response);

    // http://groups.google.com/group/faye-users/browse_thread/thread/4a01bb7d25d3636a
    if (requestMethod === 'OPTIONS' || request.headers['access-control-request-method'] === 'POST')
      return this._handleOptions(response);

    if (Faye.EventSource.isEventSource(request))
      return this.handleEventSource(request, response);

    if (requestMethod === 'GET')
      return this._callWithParams(request, response, requestUrl.query);

    if (requestMethod === 'POST')
      return Faye.withDataFor(request, function(data) {
        var type   = (request.headers['content-type'] || '').split(';')[0],
            params = (type === 'application/json')
                   ? {message: data}
                   : querystring.parse(data);

        request.body = data;
        self._callWithParams(request, response, params);
      });

    this._returnError(response, {message: 'Unrecognized request type'});
  },


 _handleOptions: function(response) {
    var headers = {
      'Access-Control-Allow-Credentials': 'false',
      'Access-Control-Allow-Headers':     'Accept, Content-Type, Pragma, X-Requested-With',
      'Access-Control-Allow-Methods':     'POST, GET, PUT, DELETE, OPTIONS',
      'Access-Control-Allow-Origin':      '*',
      'Access-Control-Max-Age':           '86400'
    };
    response.writeHead(200, headers);
    response.end('');
  },


func serveOther(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serve other: ", r.URL)

	fmt.Println("REQUEST URL: ", r.URL.Path)
	fmt.Println("REQUEST RAW QUERY: ", r.URL.RawQuery)
	fmt.Println("REQUEST HEADER ", r.Header)

	if isEventSource(r) {
		handleEventSource(w, r)
	} else {
		serveLongPolling(f, w, r)
	}
	w.WriteHeader(http.StatusOK)
	return
}
*/
/*
func handleEventSource(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handle event source: ", r.URL.Path)
	// create a new connection for the event source action
	clientId := strings.Split(r.URL.Path, "/")[2]
	fmt.Println("clientID: ", clientId)
	es := eventsource.New(nil, nil)
	c := &Connection{send: make(chan []byte, 256), es: es, isWebsocket: false}
	// TODO: NEED TO ASSOCIATED THE EXISTING FAYE CLIENT INFO/SUBSCRIPTIONS WITH THE CONNECTION CHANNEL
	// USE CLIENT ID TO UPDATE FAYE INFO WITH ES CONNETION CHANNEL
	f.UpdateClientChannel(clientId, c.send)
	go c.esWriter(f)
	c.es.ServeHTTP(w, r)
	return
}

handleEventSource: function(request, response) {
    var es       = new Faye.EventSource(request, response, {ping: this._options.ping}),
        clientId = es.url.split('/').pop(),
        self     = this;

    this.debug('Opened EventSource connection for ?', clientId);
    this._server.openSocket(clientId, es, request);

    es.onclose = function(event) {
      self._server.closeSocket(clientId);
      es = null;
    };
  },

/*
serverWs - provides an http handler for upgrading a connection to a websocket connection
*/
func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("METHOD: ", r.Method)
	fmt.Println("REQUEST URL: ", r.URL.Path)
	fmt.Println("REQUEST RAW QUERY: ", r.URL.RawQuery)
	fmt.Println("REQUEST HEADER ", r.Header)

	// server static assets
	// TODO - detect if it's a req that matches a static asset and if so serve it

	// handle options
	if r.Method == "OPTIONS" || r.Header.Get("Access-Control-Request-Method") == "POST" {
		handleOptions(w, r)
	}

	if isEventSource(r) {
		fmt.Println("Is event source")
	}

	if r.Method != "GET" {
		//http.Error(w, "Method not allowed", 405)
		serveLongPolling(f, w, r)
		return
	}

	/*
	   if r.Header.Get("Origin") != "http://"+r.Host {
	           http.Error(w, "Origin not allowed", 403)
	           return
	   }
	*/

	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		//http.Error(w, "Not a websocket handshake", 400)
		fmt.Println("NOT A WEBSOCKET HANDSHAKE")
		serveLongPolling(f, w, r)
		return
	} else if err != nil {
		fmt.Println(err)
		return
	}
	c := &Connection{send: make(chan []byte, 256), ws: ws, isWebsocket: true}
	go c.writer(f)
	c.reader(f)
}

/*
handleOptions allows for access control awesomeness
*/
func handleOptions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handle options!")
	w.Header().Set("Access-Control-Allow-Credentials", "false")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Pragma, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.WriteHeader(http.StatusOK)
	return
}

func isEventSource(r *http.Request) bool {
	fmt.Println("isEventSource? ", r.Method)
	if r.Method != "GET" {
		return false
	}

	accept := r.Header.Get("Accept")
	fmt.Println("Accept: ", accept)
	return accept == "text/event-stream"
}

var f *FayeServer

func Start(addr string) {
	f = NewFayeServer()
	//http.HandleFunc("/faye", serveWs)
	//http.HandleFunc("/", serveOther)

	// serve static assets workaround
	http.Handle("/file/", http.StripPrefix("/file", http.FileServer(http.Dir("/Users/paul/go/src/github.com/pcrawfor/fayego/runner"))))

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
