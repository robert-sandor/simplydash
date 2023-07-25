package internal

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type WebsocketServer struct {
	appService  AppService
	connections map[string]*WebsocketConnection
	logger      *log.Entry
}

func NewWebsocketServer(appService AppService) *WebsocketServer {
	return &WebsocketServer{
		appService:  appService,
		connections: make(map[string]*WebsocketConnection),
		logger:      log.WithField("name", "websocket-server"),
	}
}

func (ws *WebsocketServer) Init() {
	ws.logger.Debug("init")
	go ws.run()
}

func (ws *WebsocketServer) Connect(id string, conn *websocket.Conn) {
	log.WithField("id", id).Debug("websocket client connected")
	connection := NewWebsocketConnection(id, conn)
	connection.Init(ws.getAppsAsString())
	ws.connections[id] = connection
}

func (ws *WebsocketServer) run() {
	for {
		<-ws.appService.UpdateCh()
		ws.logger.Debug("update received")
		go ws.notifyConnections()
	}
}

func (ws *WebsocketServer) getAppsAsString() string {
	apps := ws.appService.GetApps()

	bytes, err := json.Marshal(apps)
	if err != nil {
		ws.logger.WithError(err).Error("marshalling json")
	}

	return string(bytes)
}

func (ws *WebsocketServer) notifyConnections() {
	jsonString := ws.getAppsAsString()
	for _, conn := range ws.connections {
		ws.logger.WithField("connectionId", conn.id).Debug("sending message")
		conn.updateCh <- jsonString
	}
}

type WebsocketConnection struct {
	conn     *websocket.Conn
	updateCh chan string
	stopCh   chan struct{}
	logger   *log.Entry
	id       string
}

func NewWebsocketConnection(id string, conn *websocket.Conn) *WebsocketConnection {
	return &WebsocketConnection{
		id:       id,
		conn:     conn,
		updateCh: make(chan string, 1),
		stopCh:   make(chan struct{}, 1),
		logger:   log.WithField("id", id),
	}
}

func (wc *WebsocketConnection) Init(message string) {
	go wc.sendMessage(message)
	go wc.run()
}

func (wc *WebsocketConnection) run() {
	defer func(ws *websocket.Conn) {
		_ = ws.Close()
	}(wc.conn)

	for {
		select {
		case <-wc.stopCh:
			wc.logger.Debug("closing connection")
			return
		case message := <-wc.updateCh:
			wc.logger.Debug("got update")
			go wc.sendMessage(message)
		}
	}
}

func (wc *WebsocketConnection) sendMessage(message string) {
	wc.logger.Debug("sending message to websocket")
	err := wc.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		wc.logger.WithError(err).Error("websocket message")
		wc.stopCh <- struct{}{}
	}
}
