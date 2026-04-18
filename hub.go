package main

import (
	"bytes"
	"html/template"
	"log"
	"sync"
)

// the actual message getting sent to the hub
type Message struct {
	ClientID string
	Text     string
}

// the message we will be receiving from the htmx
type WSMessage struct {
	Headers interface{} `json:"HEADERS"`
	Text    string      `json:"text"`
}

type Hub struct {
	sync.RWMutex
	messages   []*Message
	clients    map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	//parsing template once, instead of every message
	tmpl, err := parseTemplate()
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}

	for {
		select {
		case client := <-h.register:
			// TODO: add mutex to make it thread safe
			// h.Lock()
			h.clients[client] = true
			log.Printf("client registered %s", client.id)
			// adding logic to give client chat history
			for i := 0; i < len(h.messages); i++ {
				client.send <- getMessageTemplate(h.messages[i], tmpl)
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Printf("client unregistered %s", client.id)
				close(client.send)
				delete(h.clients, client)
			}
		case msg := <-h.broadcast:
			h.messages = append(h.messages, msg)

			for client := range h.clients {
				select {
				case client.send <- getMessageTemplate(msg, tmpl):
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func parseTemplate() (*template.Template, error) {
	tmpl, err := template.ParseFiles("templates/message.html")
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	return tmpl, err
}

func getMessageTemplate(msg *Message, tmpl *template.Template) []byte {
	var renderedMessage bytes.Buffer
	err := tmpl.Execute(&renderedMessage, msg)
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}

	return renderedMessage.Bytes()
}

func (h *Hub) UsernameTaken(un string) bool{
	for c := range h.clients{
		if (c.un == un){
			return true
		}
	}
	return false
}
