package eventforwarder

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/byuoitav/central-event-system/hub/base"
	"github.com/byuoitav/central-event-system/messenger"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/events"
	"github.com/byuoitav/device-monitoring/localsystem"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{}

// Service contains the running config and dependencies for an instantiation of an eventforwarder service
type Service struct {
	wsClients map[*websocket.Conn]bool
	clientMux sync.Mutex
}

// New initializes a new eventforwarder service, connects to the hub, subscribes to the room events,
// starts forwarding events, and then returns the service
func New() *Service {
	s := Service{}

	s.wsClients = make(map[*websocket.Conn]bool, 1)
	go s.reportWebSocketCount()
	// go s.lengthCheck()
	return &s
}

// HandleWebsocket upgrades the connection to a websocket connection and then sends
// messages to the client as they are recieved
func (s *Service) HandleWebsocket(ctx echo.Context) error {

	c, err := upgrader.Upgrade(ctx.Response().Writer, ctx.Request(), nil)
	// c.SetPingHandler(func(string) error { c.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })
	if err != nil {
		log.L.Errorf("Error while attempting to upgrade connection to websocket: %v", err)
	}

	s.clientMux.Lock()
	s.wsClients[c] = true
	go s.handleClose(c)
	s.clientMux.Unlock()

	return nil
}

// ForwardEvent forwards the given event to all of the currently registered websocket clients
func (s *Service) ForwardEvent(e events.Event) {

	if e.Key == "login" || e.Key == "card-read-error" {

		s.clientMux.Lock()
		for c := range s.wsClients {
			c.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.WriteJSON(e)
			if err != nil {
				log.L.Errorf("Error while forwarding event to ws client: %s", err)
				delete(s.wsClients, c)
				c.WriteMessage(websocket.CloseMessage, []byte{})
				c.Close()
			}
		}

		s.clientMux.Unlock()
	}

}

func (s *Service) handleClose(c *websocket.Conn) {
	defer func() {
		delete(s.wsClients, c)
		c.WriteMessage(websocket.CloseMessage, []byte{})
		c.Close()
	}()
	c.SetReadLimit(maxMessageSize)
	c.SetReadDeadline(time.Now().Add(pongWait))
	c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				log.L.Infof("error: %v", err)
			}
			break
		}
		log.L.Infof("Received message from socket: %s", msg)
	}
}

func (s *Service) reportWebSocketCount() {
	id := localsystem.MustSystemID()
	deviceInfo := events.GenerateBasicDeviceInfo(id)
	roomInfo := events.GenerateBasicRoomInfo(deviceInfo.RoomID)
	messenger, err := messenger.BuildMessenger(os.Getenv("HUB_ADDRESS"), base.Messenger, 1000)
	if err != nil {
		log.L.Errorf("unable to build websocket count messenger: %s", err.Error())
	}
	for {
		log.L.Debugf("sending websocket count of : %d", len(s.wsClients))
		countEvent := events.Event{
			GeneratingSystem: id,
			Timestamp:        time.Now(),
			EventTags:        []string{events.DetailState},
			TargetDevice:     deviceInfo,
			AffectedRoom:     roomInfo,
			Key:              "websocket-count",
			Value:            fmt.Sprintf("%v", len(s.wsClients)),
		}
		if messenger != nil {
			messenger.SendEvent(countEvent)
		}
		time.Sleep(1 * time.Minute)
	}
}

func (s *Service) lengthCheck() {
	for {
		log.L.Infof("Length of the thing: %d", len(s.wsClients))
		time.Sleep(20 * time.Second)
	}
}
