package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"simplydash/internal"
	"time"
)

func Websocket(u *websocket.Upgrader, svc *internal.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := u.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade connetion to a websocket err = %+v", err)
		}

		defer func(ws *websocket.Conn) {
			err := ws.Close()
			if err != nil {
				log.Printf("Failed to close websocket err = %+v", err)
			}
		}(ws)

		updateChan := make(chan struct{})
		svc.AddUpdateChannel(updateChan)
		defer svc.RemoveUpdateChannel(updateChan)
		defer close(updateChan)

		stopChan := make(chan struct{})
		defer close(stopChan)

		go send(ws, updateChan, stopChan, 60*time.Second)

		listen(ws)
		stopChan <- struct{}{}
	}
}

func send(
	ws *websocket.Conn,
	updateChan chan struct{},
	stopChan chan struct{},
	pingPeriod time.Duration,
) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				log.Printf("Failed to ping the client err = %+v", err)
			}
		case <-updateChan:
			err := ws.WriteMessage(websocket.TextMessage, []byte("update"))
			if err != nil {
				log.Printf("Failed to send update message to client err = %+v", err)
			}
		case <-stopChan:
			return
		}
	}
}

func listen(ws *websocket.Conn) {
	for {
		messageType, bytes, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				log.Printf("Error when reading message from client err = %+v", err)
			}
			return
		}

		log.Printf("Got message from client type = %d message = %s", messageType, string(bytes))
	}
}
