package main

import (
	"log"
	"net"
)

func main() {
	s := newServer()

	go s.run()

	listener, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("unable to start the server!")
	}

	defer listener.Close()
	log.Printf("Started server on port 9001")

	for {
		client, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection from %s %s", client, err)
			continue
		}

		clientObj := newClient(client)
		clientObj.join(s)
	}

}
