package lab

import (
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/events"
)

const cardReadKey = "card-read"

// Handle handles incoming events and responds to `card-read` events
func (l Lab) Handle(e events.Event) {
	log.L.Debugf("lab/handle: Caught event: %v+", e)
	if e.Key == cardReadKey {
		l.LogAttendanceForCard(e.Value)
	}
}
