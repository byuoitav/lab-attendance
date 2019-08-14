package eventforwarder

import (
	"sync"
	"time"

	"github.com/byuoitav/common/log"
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
	wsClients map[*websocket.Conn]bool
	clientMux sync.Mutex
}

// New initializes a new eventforwarder service, connects to the hub, subscribes to the room events,
// starts forwarding events, and then returns the service
func New() *Service {
	s := Service{}

	s.wsClients = make(map[*websocket.Conn]bool, 1)

	return &s
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
