package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type host struct {
	name string
	client client
}

type client chan<- string

var (
	entering = make(chan host)
	leaving  = make(chan host)
	messages = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening at " + listener.Addr().String())

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[host]bool) //a channel for each connected host
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli.client <- msg
			}
		case cli := <-entering:
			clients[cli] = true
			go showInRoom(clients, cli)
		case cli := <-leaving:
			delete(clients, cli)
			close(cli.client)
		}
	}
}

func showInRoom(clients map[host]bool, to host) {
	for c := range clients {
	  if c.name != to.name{
	    to.client <- c.name + " is in the room"
	  }
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- host{name: who, client:ch}

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}

	leaving <- host{name: who, client:ch}
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
