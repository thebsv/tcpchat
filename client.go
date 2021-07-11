package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Client struct {
	connection net.Conn
	name       string
	channel    chan string
	serverObj  *ChatServer
}

func newClient(conn net.Conn) *Client {
	return &Client{
		connection: conn,
		name:       "new",
		channel:    make(chan string, 10),
		serverObj:  nil,
	}
}

func (c *Client) loop() {
	// Works by calling a corresponding function in the server (rpc),
	// depending on the argument passed to it from the command line.

	for {
		log.Printf("%s trying to receive", c.name)
		select {
		case from := <-c.channel:
			rmsg := fmt.Sprintf("%s received from server: %s ", c.name, from)
			log.Printf(rmsg)
			if _, err := c.connection.Write([]byte(rmsg)); err != nil {
				c.quit()
				log.Fatalf("Client %s COULD NOT WRITE MESSAGE", c.name)
			}
			c.connection.Write([]byte("\r\n"))
		default:
			log.Printf("%s trying to read input", c.name)
			if msg, err := bufio.NewReader(c.connection).ReadString('\n'); err == nil {

				log.Printf("%s INPUT %s", c.name, msg)

				cmd := strings.Split(strings.Trim(msg, "\r\n"), " ")
				// switch using message and send the command to the server.
				switch cmd[0] {
				case "LIST":
					log.Printf("%s CASE LIST ", c.name)
					c.list()
				case "NAME":
					log.Printf("%s CASE NAME ", c.name)
					if len(cmd) > 1 {
						c.changeName(cmd[1])
					}
				case "QUIT":
					log.Printf("%s CASE QUIT ", c.name)
					c.quit()
					err := c.connection.Close()
					if err != nil {
						log.Fatalf("%s CASE QUIT connection close error ", c.name)
					}
					return
				default:
					log.Printf("%s CASE DEFAULT ", c.name)
					c.sendMessage(cmd[0])
				}

				log.Printf("%s DONE", c.name)
			}
		}
	}
}

func (c *Client) join(serv *ChatServer) {
	serv.join(*c)
	c.serverObj = serv
	go c.loop()
}

func (c *Client) changeName(given string) {
	c.serverObj.changeName(*c, given)
	c.name = given
}

func (c *Client) list() {
	c.serverObj.list(*c)
}

func (c *Client) sendMessage(msg string) {
	c.serverObj.broadcast(*c, msg)
}

func (c *Client) quit() {
	c.serverObj.quit(*c)
	c.connection.Close()
}
