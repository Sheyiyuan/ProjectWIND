package main

import (
	"ProjectWIND/protocol"
	"log"
)

func init() {

}

func main() {
	_, err := protocol.WebSocketHandler()
	if err != nil {
		log.Fatal(err)
	}
}
