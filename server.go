package main

import (
	"os"

	"github.com/byuoitav/central-event-system/hub/base"
	"github.com/byuoitav/central-event-system/messenger"
	"github.com/byuoitav/common"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/events"
	"github.com/byuoitav/lab-attendance/eventforwarder"
	"github.com/byuoitav/lab-attendance/handlers"
	"github.com/labstack/echo/middleware"
)

func main() {

	log.SetLevel("debug")
	deviceInfo := events.GenerateBasicDeviceInfo(os.Getenv("SYSTEM_ID"))
	msgr, err := messenger.BuildMessenger(os.Getenv("HUB_ADDRESS"), base.Messenger, 1000)
	if err != nil {
		log.L.Errorf("Error while building messenger: %s", err)
	}
	msgr.SubscribeToRooms(deviceInfo.RoomID)

	ef, _ := eventforwarder.NewService()

	router := common.NewRouter()

	port := ":8243"

	router.POST("/api/v1/login/:byuID", handlers.Login(msgr, deviceInfo))
	router.GET("/websocket", ef.HandleWebsocket)

	router.Group("/", middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "autoclave-dist",
		Index:  "index.html",
		HTML5:  true,
		Browse: true,
	}))

	router.Start(port)
}
