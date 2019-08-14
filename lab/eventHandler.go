package lab

import (
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/events"
)

const cardReadKey = "card-read"

// Handle handles incoming events and responds to `card-read` events
func (l Lab) Handle(e events.Event) {
	if e.Key == cardReadKey {
		log.L.Debugf("lab/Handle: Caught Card Read Event: %+v", e)
		l.LogAttendanceForCard(e.Value)
	}
}
