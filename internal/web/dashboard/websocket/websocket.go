package websocket

import (
	"bytes"
	"context"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/websocket"
)

type Client struct {
	broker *WebSocketBroker
	conn   *websocket.Conn
	send   chan []byte
}

type WebSocketBroker struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	redis      *redis.Client
	pubsub     *redis.PubSub
}

func NewWebSocketBroker(redisAddr string) (*WebSocketBroker, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	pubsub := rdb.Subscribe(context.Background(), "dashboard:updates")

	broker := &WebSocketBroker{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		redis:      rdb,
		pubsub:     pubsub,
	}

	go func() {
		for msg := range pubsub.Channel() {
			broker.broadcast <- []byte(msg.Payload)
		}
	}()
	go broker.Run()

	return broker, nil
}

func (b *WebSocketBroker) Run() {
	for {
		select {
		case client := <-b.register:
			b.mu.Lock()
			b.clients[client] = true
			b.mu.Unlock()

		case client := <-b.unregister:
			b.mu.Lock()
			if _, ok := b.clients[client]; ok {
				delete(b.clients, client)
				close(client.send)
			}
			b.mu.Unlock()

		case message := <-b.broadcast:
			b.mu.RLock()
			for client := range b.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(b.clients, client)
				}
			}
			b.mu.RUnlock()
		}
	}
}

func (b *WebSocketBroker) Close() error {
	if b.pubsub != nil {
		return b.pubsub.Close()
	}
	return nil
}

func (c *Client) readPump() {
	defer func() {
		c.broker.unregister <- c
		c.conn.Close()
	}()

	var message []byte
	for {
		if err := websocket.Message.Receive(c.conn, &message); err != nil {
			log.Printf("error receiving message: %v", err)
			break
		}
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for message := range c.send {
		currentPath := getPathFromQueryString(c.conn.Request().URL.RawQuery)
		renderedComponent, err := RenderComponent(currentPath, message)
		if err != nil {
			log.Printf("error rendering HTML: %v", err)
			continue
		}
		buf := bytes.NewBuffer(nil)
		renderedComponent.Render(context.Background(), buf)
		if err := websocket.Message.Send(c.conn, buf.Bytes()); err != nil {
			if err.Error() == "EOF" || err.Error() == "websocket: close sent" {
				return
			}
			log.Printf("error sending message: %v", err)
			continue
		}
	}
}
