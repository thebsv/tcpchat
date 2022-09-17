package main

import (
	"fmt"
	"log"
	"strings"
)

type ChatServer struct {
	serveClient map[*Client]bool
}

func newServer() *ChatServer {
	return &ChatServer{
		serveClient: make(map[*Client]bool),
	}
}

func (s *ChatServer) run() {
	// The client has the server object and uses
	// RPC to operate the server. Cleanup is done here for clients
	// That have exited without closing their connection.

	for {
	}

}

func (s *ChatServer) changeName(client *Client, given string) {
	if _, exist := s.serveClient[client]; exist {
		delete(s.serveClient, client)
		s.serveClient[client] = true
	}
}

func (s *ChatServer) join(client *Client) {
	log.Printf("New client %s joined the server!", client.connection.RemoteAddr().String())
	s.serveClient[client] = true
}

func (s *ChatServer) broadcast(client *Client, msg string) {
	for c := range s.serveClient {
		message := fmt.Sprintf("SERVER %s: %s", client.name, msg)
		c.channel.push(message)
	}
}

func (s *ChatServer) list(client *Client) {
	var clientList []string
	for c := range s.serveClient {
		clientList = append(clientList, c.name)
	}
	client.channel.push(strings.Join(clientList, " "))
}

func (s *ChatServer) quit(client *Client) {
	if _, ok := s.serveClient[client]; ok {
		delete(s.serveClient, client)
	}
}
