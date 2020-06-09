package main

import (
	"errors"
	"strings"
)

type hub struct {
	channels        map[string]*channel
	clients         map[string]*client
	commands        chan command
	deregistrations chan *client
	registrations   chan *client
}

func newHub() *hub {
	return &hub{
		channels:        make(map[string]*channel),
		clients:         make(map[string]*client),
		commands:        make(chan command),
		deregistrations: make(chan *client),
		registrations:   make(chan *client),
	}
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.registrations:
			h.register(client)
		case client := <-h.deregistrations:
			h.deregister(client)
		case cmd := <-h.commands:
			switch cmd.id {
			case JOIN:
				h.joinChannel(cmd.sender, cmd.recipient)
			case LEAVE:
				h.leaveChannel(cmd.sender, cmd.recipient)
			case MSG:
				h.message(cmd.sender, cmd.recipient, cmd.body)
			case USRS:
				h.listUsers(cmd.sender)
			case CHNS:
				h.listChannels(cmd.sender)
			}
		}
	}
}

func (h *hub) register(c *client) {
	if _, exists := h.clients[c.username]; exists {
		c.username = ""
		c.err(errors.New("username taken"))
	} else {
		h.clients[c.username] = c
		c.ok()
	}
}

func (h *hub) deregister(c *client) {
	if _, exists := h.clients[c.username]; exists {
		delete(h.clients, c.username)
		for _, channel := range h.channels {
			delete(channel.clients, c)
		}
	}
}

func (h *hub) joinChannel(u, c string) {
	if client, ok := h.clients[u]; ok {
		if channel, ok := h.channels[c]; ok {
			channel.clients[client] = true
		} else {
			h.channels[c] = newChannel(c)
			h.channels[c].clients[client] = true
		}
	}
}

func (h *hub) leaveChannel(u, c string) {
	if client, ok := h.clients[u]; ok {
		if channel, ok := h.channels[c]; ok {
			delete(channel.clients, client)
		}
	}
}

func (h *hub) message(u, r string, m []byte) {
	if sender, ok := h.clients[u]; ok {
		switch r[0] {
		case '#':
			if channel, ok := h.channels[r]; ok {
				if _, ok := channel.clients[sender]; ok {
					channel.broadcast(sender.username, m)
				}
			}
		case '@':
			if user, ok := h.clients[r]; ok {
				user.conn.Write(append(m, '\n'))
			}
		}
	}
}

func (h *hub) listUsers(u string) {
	if sender, ok := h.clients[u]; ok {
		keys := make([]string, 0, len(h.clients))
		for k, _ := range h.clients {
			keys = append(keys, k)
		}
		response := []byte(strings.Join(keys, ", "))
		response = append(response, byte('\n'))
		h.clients[sender.username].conn.Write(response)
	}
}

func (h *hub) listChannels(u string) {
	if sender, ok := h.clients[u]; ok {
		keys := make([]string, 0, len(h.channels))
		for k, _ := range h.channels {
			keys = append(keys, k)
		}
		response := []byte(strings.Join(keys, ", "))
		response = append(response, byte('\n'))
		h.clients[sender.username].conn.Write(response)
	}
}
