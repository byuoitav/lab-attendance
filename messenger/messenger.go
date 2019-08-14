package messenger

import (
	"fmt"
	"sync"
	"time"

	"github.com/byuoitav/central-event-system/hub/base"
	"github.com/byuoitav/central-event-system/messenger"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/events"
	"github.com/byuoitav/lab-attendance/cache"
)

// Messenger represents an instance of a messenger and contains all the configuration needed to run
type Messenger struct {
	m        *messenger.Messenger
	info     events.BasicDeviceInfo
	handlers []Handler
	hMu      sync.Mutex
}

// Handler represents a function which handles an event sent to it
type Handler func(events.Event)

// New sets up and returns a new Messenger
func New(addr string, info events.BasicDeviceInfo) (*Messenger, error) {

	m, err := messenger.BuildMessenger(addr, base.Messenger, 1000)
	if err != nil {
		e := fmt.Errorf("Error while trying to create new messenger: %s", err)
		log.L.Error(e)
		return nil, e
	}
	m.SubscribeToRooms(info.RoomID)

	msgr := &Messenger{
		m:    m,
		info: info,
	}

	go msgr.handleEvents()

	return msgr, nil
}

// Register returns a channel to the caller through which events will be sent
func (m *Messenger) Register(h Handler) {
	m.hMu.Lock()
	m.handlers = append(m.handlers, h)
	m.hMu.Unlock()
}

// SendLoginEvent sends a login event for the given person
func (m *Messenger) SendLoginEvent(p cache.Person) {
	type loginData struct {
		FirstName string
		Name      string
	}
	e := events.Event{
		Key:   "login",
		Value: "true",
		User:  p.NetID,
		Data: loginData{
			Name:      p.Name,
			FirstName: p.FirstName,
		},
	}

	m.SendEvent(e)
}

// SendEvent will send the given event to the message bus ensuring that the device information and timestamp is correct
func (m *Messenger) SendEvent(e events.Event) {
	e.GeneratingSystem = m.info.DeviceID
	e.TargetDevice = m.info
	e.AffectedRoom = m.info.BasicRoomInfo
	e.Timestamp = time.Now()
	m.m.SendEvent(e)
}

func (m *Messenger) handleEvents() {
	var e events.Event

	for {
		e = m.m.ReceiveEvent()

		m.hMu.Lock()
		for _, h := range m.handlers {
			h(e)
		}

		m.hMu.Unlock()
	}
}
