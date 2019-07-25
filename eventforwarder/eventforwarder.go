package eventforwarder

import (
	"os"
	"sync"
	"time"

	"github.com/byuoitav/central-event-system/hub/base"
	"github.com/byuoitav/central-event-system/messenger"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/v2/events"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

const (
	writeWait = 10 * time.Second
)

var upgrader = websocket.Upgrader{}

// Service contains the running config and dependencies for an instantiation of an eventforwarder service
type Service struct {
	m         *messenger.Messenger
	wsClients map[*websocket.Conn]bool
	clientMux sync.Mutex
}

// NewService initializes a new eventforwarder service, connects to the hub, subscribes to the room events,
// starts forwarding events, and then returns the service
func NewService() (*Service, error) {
	s := Service{}
	m, err := messenger.BuildMessenger(os.Getenv("HUB_ADDRESS"), base.Messenger, 1000)
	if err != nil {
		err = nerr.Createf("Internal", "Error while attempting to build the messenger: %s", err)
		log.L.Error(err)
		return nil, err
	}

	s.m = m
	s.m.SubscribeToRooms(events.GenerateBasicDeviceInfo(os.Getenv("SYSTEM_ID")).RoomID)
	s.wsClients = make(map[*websocket.Conn]bool, 1)

	go s.forwardEvents()
	return &s, nil
}

// HandleWebsocket upgrades the connection to a websocket connection and then sends
// messages to the client as they are recieved
func (s *Service) HandleWebsocket(ctx echo.Context) error {

	c, err := upgrader.Upgrade(ctx.Response().Writer, ctx.Request(), nil)
	if err != nil {
		log.L.Errorf("Error while attempting to upgrade connection to websocket: %v", err)
	}

	s.clientMux.Lock()
	s.wsClients[c] = true
	s.clientMux.Unlock()

	return nil
}

func (s *Service) forwardEvents() {
	var e events.Event

	for {
		e = s.m.ReceiveEvent()

		log.L.Debugf("Got event: %+v\n", e)

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
