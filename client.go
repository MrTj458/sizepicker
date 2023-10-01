package main

import (
	"io"
	"log/slog"
	"time"

	"golang.org/x/net/websocket"
)

type Cmd struct {
	Cmd    string `json:"cmd"`
	Name   string `json:"name"`
	Choice int    `json:"choice"`
}

type Client struct {
	room   *Room
	conn   *websocket.Conn
	send   chan any
	Name   string `json:"name"`
	Choice int    `json:"choice"`
}

func (c *Client) read() {
	defer func() {
		c.conn.Close()
		c.room.unregister <- c
		c.room.broadcast <- c.room
	}()

	for {
		var cmd Cmd
		err := websocket.JSON.Receive(c.conn, &cmd)
		if err != nil {
			if err != io.EOF {
				slog.Error("error reading message", "err", err)
			}
			break
		}
		slog.Debug("command received", "cmd", cmd)

		switch cmd.Cmd {
		case "name":
			c.Name = cmd.Name
		case "pick":
			c.Choice = cmd.Choice
		case "show":
			c.room.Show = true
		case "reset":
			c.room.reset <- true
		}

		c.room.broadcast <- c.room
	}
}

func (c *Client) write() {
	ticker := time.NewTicker(15 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, open := <-c.send:
			if !open {
				return
			}

			err := websocket.JSON.Send(c.conn, message)
			if err != nil {
				slog.Error("error sending message", "err", err)
				return
			}
		case <-ticker.C:
			err := websocket.JSON.Send(c.conn, c.room)
			if err != nil {
				return
			}
		}
	}
}
