package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

//go:embed web/*
var webFS embed.FS

var (
	indexPage *template.Template
	port      int
	room      *Room
)

func main() {
	portFlag := flag.Int("port", 3000, "the port for the server to run on")
	debugFlag := flag.Bool("debug", false, "print debug statements")
	flag.Parse()
	port = *portFlag

	var slogOpts *slog.HandlerOptions
	if *debugFlag {
		slogOpts = &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, slogOpts)))

	var err error
	indexPage, err = template.ParseFS(webFS, "web/templates/index.html")
	if err != nil {
		slog.Error("error parsing template", "err", err)
		os.Exit(1)
	}

	room = newRoom()
	go room.run()

	staticFS, err := fs.Sub(webFS, "web/static")
	if err != nil {
		slog.Error("error getting static directory", "err", err)
		os.Exit(1)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	http.Handle("/ws", websocket.Handler(handleWebSocket))
	http.HandleFunc("/", handleIndex)

	slog.Info("starting server", "port", port, "debug", *debugFlag)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		slog.Error("error running server", "err", err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	err := indexPage.Execute(w, nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
