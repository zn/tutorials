package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
)

type ID int

const (
	REG ID = iota
	JOIN
	LEAVE
	MSG
	CHNS
	USRS
)

var (
	DELIMITER = []byte("~")
)

type command struct {
	id        ID
	recipient string
	sender    string
	body      []byte
}

type client struct {
	conn       net.Conn
	outbound   chan<- command
	register   chan<- *client
	deregister chan<- *client
	username   string
}

func (c *client) read() error {
	for {
		msg, err := bufio.NewReader(c.conn).ReadBytes('\n')
		if err == io.EOF {
			c.deregister <- c
			return nil
		}
		if err != nil {
			return err
		}
		c.handle(msg)
	}
}

func (c *client) handle(message []byte) {
	cmd := bytes.ToUpper(bytes.TrimSpace(bytes.Split(message, []byte(" "))[0]))
	args := bytes.TrimSpace(bytes.TrimPrefix(message, cmd))
	switch string(cmd) {
	case "REG":
		if err := c.reg(args); err != nil {
			c.err(err)
		}
	case "JOIN":
		if err := c.join(args); err != nil {
			c.err(err)
		}
	case "LEAVE":
		if err := c.leave(args); err != nil {
			c.err(err)
		}
	case "MSG":
		if err := c.msg(args); err != nil {
			c.err(err)
		}
	case "CHNS":
		c.chns()
	case "USRS":
		c.usrs()
	default:
		c.err(fmt.Errorf("unknown command: %s", cmd))
	}
}

func (c *client) reg(args []byte) error {
	// u := bytes.TrimSpace(args)
	if args[0] != '@' {
		return fmt.Errorf("username must begin with @")
	}
	if len(args) == 0 {
		return fmt.Errorf("username cannot be blank")
	}
	c.username = string(args)
	c.register <- c
	return nil
}

func (c *client) msg(args []byte) error {
	// args = bytes.TrimSpace(args)
	if args[0] != '#' && args[0] != '@' {
		return fmt.Errorf("recipient must be a channel ('#name') or user ('@user')")
	}
	recipient := bytes.Split(args, []byte(" "))[0]

	args = bytes.TrimSpace(bytes.TrimPrefix(args, recipient))
	l := bytes.Split(args, DELIMITER)[0]
	length, err := strconv.Atoi(string(l))
	if err != nil {
		return fmt.Errorf("body length must be present")
	}
	if length == 0 {
		return fmt.Errorf("body length must be at least 1")
	}

	padding := len(l) + len(DELIMITER)
	body := args[padding : padding+length]

	c.outbound <- command{
		id:        MSG,
		recipient: string(recipient),
		sender:    c.username,
		body:      body,
	}
	return nil
}

func (c *client) join(args []byte) error {
	if args[0] != byte('#') {
		return fmt.Errorf("channel must begin with #")
	}
	c.outbound <- command{
		id:        JOIN,
		recipient: string(bytes.TrimSpace(args)),
		sender:    c.username,
		body:      nil,
	}
	return nil
}

func (c *client) leave(args []byte) error {
	if args[0] != byte('#') {
		return fmt.Errorf("channel must begin with #")
	}
	c.outbound <- command{
		id:        LEAVE,
		recipient: string(bytes.TrimSpace(args)),
		sender:    c.username,
		body:      nil,
	}
	return nil
}

func (c *client) chns() {
	c.outbound <- command{
		id:        CHNS,
		recipient: c.username,
		sender:    c.username,
		body:      nil,
	}
}

func (c *client) usrs() {
	c.outbound <- command{
		id:        USRS,
		recipient: c.username,
		sender:    c.username,
		body:      nil,
	}
}

func (c *client) err(e error) {
	c.conn.Write([]byte("ERR " + e.Error() + "\n"))
}

func (c *client) ok() {
	c.conn.Write([]byte("OK\n"))
}
