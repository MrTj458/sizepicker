package main

import (
	"log/slog"
)

type Room struct {
	Clients []*Client `json:"clients"`
	Show    bool      `json:"show"`

	broadcast  chan any
	register   chan *Client
	unregister chan *Client
	reset      chan bool
}

func newRoom() *Room {
	return &Room{
		Clients: make([]*Client, 0),

		broadcast:  make(chan any),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		reset:      make(chan bool),
	}
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			r.Clients = append(r.Clients, client)
			slog.Debug("client registered", "clients", r.Clients)
		case client := <-r.unregister:
			close(client.send)
			r.deleteClient(client)
			slog.Debug("client unregistered", "clients", r.Clients)
		case message := <-r.broadcast:
			for _, client := range r.Clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					r.deleteClient(client)
				}
			}
		case <-r.reset:
			for _, c := range r.Clients {
				c.Choice = 0
			}
			r.Show = false
		}
	}
}

func (r *Room) deleteClient(c *Client) {
	newClients := make([]*Client, 0)
	for _, client := range r.Clients {
		if c == client {
			continue
		}
		newClients = append(newClients, client)
	}
	r.Clients = newClients
}
