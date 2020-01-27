package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/net/websocket"
)

var acceptedSingleCommands = []string{"help", "yah", "bread"}
var acceptedTupleCommands = []string{"yah yeet"}

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
			if len(parts) == 2 {
				go func(m Message) {
					buildAndPostMessage(ws, m, processSingleCommand(parts[1]))
				}(m)
			} else if len(parts) == 3 {
				go func(m Message) {
					buildAndPostMessage(ws, m, processTupleCommand(parts[1], parts[2]))
				}(m)
			} else {
				buildAndPostMessage(ws, m, createErrorMessage())
			}
		}
	}
}

func processSingleCommand(command string) string {
	if command == acceptedSingleCommands[0] {
		return createCommandList()
	} else if command == acceptedSingleCommands[1] {
		return "yeet"
	} else if command == acceptedSingleCommands[2] {
		return "let's get this bread"
	}
	return createErrorMessage()
}

func processTupleCommand(first string, second string) string {
	commands := strings.Fields(acceptedTupleCommands[0])
	if first == commands[0] && second == commands[1] {
		return fmt.Sprintf("yah yoink")
	}
	return createErrorMessage()
}

func createErrorMessage() string {
	errorMessage := "Sorry, that does not compute. Try one of the below commands.\n"
	return errorMessage + createCommandList()
}

func createCommandList() string {
	commands := ""
	for _, command := range acceptedSingleCommands {
		commands = commands + "@botname " + command + "\n"
	}
	for _, command := range acceptedTupleCommands {
		commands = commands + "@botname " + command + "\n"
	}
	return commands
}

func buildAndPostMessage(ws *websocket.Conn, m Message, text string) {
	m.Text = text
	postMessage(ws, m)
}
