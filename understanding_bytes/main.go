package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}

	hub := newHub()
	go hub.run()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		c := &client{
			conn:       conn,
			outbound:   hub.commands,
			register:   hub.registrations,
			deregister: hub.deregistrations,
			username:   fmt.Sprintf("user#%d", rand.Int()),
		}
		fmt.Println("New client!")
		go c.read()
	}
}
