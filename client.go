package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type Queue struct {
	queue []string
	lock  sync.RWMutex
}

func (q *Queue) push(message string) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue = append(q.queue, message)
}

func (q *Queue) pop() (string, error) {
	if len(q.queue) > 0 {
		q.lock.Lock()
		defer q.lock.Unlock()
		ret := q.queue[0]
		q.queue = q.queue[1:]
		return ret, nil
	}
	return "-1", fmt.Errorf("Queue is empty")
}

func (q *Queue) top() (string, error) {
	if len(q.queue) > 0 {
		return q.queue[0], nil
	}
	return "-1", fmt.Errorf("Queue is empty")
}

func (q *Queue) empty() bool {
	return len(q.queue) == 0
}

func (q *Queue) length() int {
	return len(q.queue)
}

type Client struct {
	connection net.Conn
	name       string
	channel    *Queue
	serverObj  *ChatServer
}

func newClient(conn net.Conn) *Client {
	qu := &Queue{queue: make([]string, 0)}
	return &Client{
		connection: conn,
		name:       "new",
		channel:    qu,
		serverObj:  nil,
	}
}

func writeToConnectionOrDestroy(c *Client, msg string) bool {
	msg = msg + "\r\n"
	log.Println(msg)
	_, err := net.Conn.Write(c.connection, []byte(msg))
	if err != nil {
		log.Printf("Could not write to connection %s ", err)
		return false
	}
	return true
}

func flushQueue(c *Client) {
	if !c.channel.empty() {
		rc := true
		for c.channel.length() > 0 {
			mesg, _ := c.channel.pop()
			rc := rc && writeToConnectionOrDestroy(c, mesg)
			if !rc {
				c.quit()
				break
			}
		}
	}
}

func (c *Client) loop() {
	// Works by calling a corresponding function in the server (rpc),
	// depending on the argument passed to it from the command line.

	for {

		log.Printf("%s trying to receive", c.name)
		flushQueue(c)
		log.Println("Finished receiving, reading from stdin")
		msg, err := bufio.NewReader(c.connection).ReadString('\n')
		if err != nil {
			log.Printf("nothing to read from stdin")
			continue
		} else {
			log.Printf("%s INPUT %s", c.name, msg)
			cmd := strings.Split(strings.Trim(msg, "\n"), " ")

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
				log.Printf("%s CASE DEFAULT MESSAGE TO BE BR %s ", c.name, cmd[0])
				c.sendMessage(cmd[0])
			}
			log.Printf("%s DONE", c.name)
		}
	}
}

func (c *Client) join(serv *ChatServer) {
	serv.join(c)
	c.serverObj = serv
	go c.loop()
}

func (c *Client) changeName(given string) {
	c.name = given
	c.serverObj.changeName(c, given)
}

func (c *Client) list() {
	c.serverObj.list(c)
}

func (c *Client) sendMessage(msg string) {
	c.serverObj.broadcast(c, msg)
	flushQueue(c)
}

func (c *Client) quit() {
	c.serverObj.quit(c)
	c.connection.Close()
}
