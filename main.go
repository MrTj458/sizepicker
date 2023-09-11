package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

//go:embed web/*
var webFS embed.FS

var (
	indexPage *template.Template
	manager   clientManager
	port      int
)

type clientManager struct {
	Clients []*client `json:"clients"`
	Show    bool      `json:"show"`
	sync.Mutex
}

func (cm *clientManager) addClient(c *client) {
	cm.Lock()
	defer cm.Unlock()

	cm.Clients = append(cm.Clients, c)
}

func (cm *clientManager) removeClient(c *client) {
	cm.Lock()
	defer cm.Unlock()

	var newClients []*client
	for _, client := range cm.Clients {
		if client == c {
			c.conn.Close()
			continue
		}

		newClients = append(newClients, client)
	}

	cm.Clients = newClients
}

func (cm *clientManager) sendUpdate() {
	cm.Lock()
	defer cm.Unlock()

	for _, client := range cm.Clients {
		err := websocket.JSON.Send(client.conn, &manager)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (cm *clientManager) read(c *client) {
	for {
		var cmd command
		err := websocket.JSON.Receive(c.conn, &cmd)
		if err != nil {
			cm.removeClient(c)
			cm.sendUpdate()
			break
		}

		switch cmd.CMD {
		case "name":
			c.Name = cmd.Name
		case "pick":
			c.Choice = cmd.Choice
		case "show":
			cm.Show = true
		case "reset":
			cm.Show = false
			for _, client := range cm.Clients {
				client.Choice = 0
			}
		}
		cm.sendUpdate()
	}
}

type client struct {
	conn   *websocket.Conn
	Name   string `json:"name"`
	Choice int    `json:"choice"`
}

type command struct {
	CMD    string `json:"cmd"`
	Name   string `json:"name"`
	Choice int    `json:"choice"`
}

func main() {
	portFlag := flag.Int("port", 3000, "the port for the server to run on")
	flag.Parse()
	port = *portFlag

	var err error
	indexPage, err = template.ParseFS(webFS, "web/templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	staticFS, err := fs.Sub(webFS, "web/static")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	http.Handle("/ws", websocket.Handler(handleWebSocket))
	http.HandleFunc("/", handleIndex)

	fmt.Println("Starting server on port", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	err := indexPage.Execute(w, nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func handleWebSocket(ws *websocket.Conn) {
	c := &client{
		conn: ws,
	}

	manager.addClient(c)
	manager.sendUpdate()
	manager.read(c)
}
