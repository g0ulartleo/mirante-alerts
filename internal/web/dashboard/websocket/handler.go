package websocket

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

func HandleWebSocket(broker *WebSocketBroker) echo.HandlerFunc {
	return func(c echo.Context) error {
		websocket.Handler(func(ws *websocket.Conn) {
			client := &Client{
				broker: broker,
				conn:   ws,
				send:   make(chan []byte, 256),
			}

			client.broker.register <- client

			go client.writePump()
			client.readPump()
		}).ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
