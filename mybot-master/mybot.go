package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	token, err := ioutil.ReadFile("token.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	strToken := string(token)

	// start a websocket-based Real Time API session
	ws, id := slackConnect(strToken)
	fmt.Println("mybot ready, ^C exits")

	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		// see if we're mentioned
		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			// if so try to parse if
			parts := strings.Fields(m.Text)
			if len(parts) == 2 && parts[1] == "help" {
				go func(m Message) {
					m.Text = fmt.Sprintf("@botname yah\n@botname help")
					postMessage(ws, m)
				}(m)
			} else if len(parts) == 2 && parts[1] == "yah" {
				go func(m Message) {
					m.Text = "yeet"
					postMessage(ws, m)
				}(m)
			} else {
				// huh?
				m.Text = fmt.Sprintf("sorry, that does not compute. Try \"@botname help\"\n")
				postMessage(ws, m)
			}
		}
	}
}
