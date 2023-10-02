package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

//go:embed client/dist/*
var webFS embed.FS

var (
	room *Room
)

func main() {
	portFlag := flag.Int("port", 3000, "the port for the server to run on")
	debugFlag := flag.Bool("debug", false, "print debug statements")
	flag.Parse()

	var slogOpts *slog.HandlerOptions
	if *debugFlag {
		slogOpts = &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, slogOpts)))

	room = newRoom()
	go room.run()

	http.Handle("/ws", websocket.Handler(handleWebSocket))

	staticFS, err := fs.Sub(webFS, "client/dist")
	if err != nil {
		slog.Error("error getting client directory", "err", err)
		os.Exit(1)
	}
	http.Handle("/", http.FileServer(http.FS(staticFS)))

	slog.Info("starting server", "port", *portFlag, "debug", *debugFlag)
	err = http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), nil)
	if err != nil {
		slog.Error("error running server", "err", err)
	}
}

func handleWebSocket(ws *websocket.Conn) {
	client := &Client{
		room: room,
		conn: ws,
		send: make(chan any, 256),
	}
	client.room.register <- client
	client.room.broadcast <- client.room

	go client.write()
	client.read()
}
