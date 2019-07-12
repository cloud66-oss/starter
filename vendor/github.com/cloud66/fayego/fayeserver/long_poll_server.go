package fayeserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func serveLongPolling(f *FayeServer, w http.ResponseWriter, r *http.Request) {
	jsonMessage := r.FormValue("message")
	//jsonpParam := r.FormValue("jsonp")

	if len(jsonMessage) == 0 {
		b, _ := ioutil.ReadAll(r.Body)
		fmt.Println("body: ", string(b))
		jsonMessage = string(b)
	}

	fmt.Println("THIS IS JSON MESSAGE ", jsonMessage)

	// handle the faye message
	response, error := f.HandleMessage([]byte(jsonMessage), nil)

	if error != nil {
		// We have to figure what correct error response is for HTTP for faye
		fmt.Println("HTTP SERVER ERROR: ", error)
		return
	} else {

		//finalResponse := jsonpParam + "(" + string(response) + ");"

		finalResponse := string(response)
		fmt.Println("THIS IS OUR HTTP RESPONSE: %v", finalResponse)
		fmt.Println("THIS IS THE W HEADERS: ", w.Header())
		w.Header().Set("Content-Type", "application/javascript")
		fmt.Fprint(w, finalResponse)
		return
	}
}
