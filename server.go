package main

import (
	"os"

	"github.com/byuoitav/common"
	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/events"
	"github.com/byuoitav/lab-attendance/cache"
	"github.com/byuoitav/lab-attendance/eventforwarder"
	"github.com/byuoitav/lab-attendance/handlers"
	"github.com/byuoitav/lab-attendance/lab"
	"github.com/byuoitav/lab-attendance/messenger"
	"github.com/labstack/echo/middleware"
)

func main() {

	log.SetLevel("debug")

	deviceInfo := events.GenerateBasicDeviceInfo(os.Getenv("SYSTEM_ID"))

	msgr, err := messenger.New(os.Getenv("HUB_ADDRESS"), deviceInfo)
	if err != nil {
		log.L.Fatalf("Error while building messenger: %s", err)
	}

	config, err := db.GetDB().GetLabConfig(deviceInfo.RoomID)
	if err != nil {
		log.L.Fatalf("Error while trying to get Lab Config from the database: %s", err)
	}
	log.L.Debugf("Got Lab Config for room %s: %+v", deviceInfo.RoomID, config)

	cache, err := cache.New()
	if err != nil {
		log.L.Fatalf("Error while trying to create cache: %s", err)
	}

	lab := lab.Lab{
		ID:    config.LabID,
		M:     msgr,
		Cache: cache,
	}

	ef := eventforwarder.New()
	msgr.Register(lab.Handle)
	msgr.Register(ef.ForwardEvent)

	router := common.NewRouter()

	port := ":8243"

	router.POST("/api/v1/login/:byuID", handlers.Login(msgr, deviceInfo, lab))
	router.GET("/api/v1/config", handlers.GetConfig(config))
	router.GET("/websocket", ef.HandleWebsocket)

	router.Group("/", middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "autoclave-dist",
		Index:  "index.html",
		HTML5:  true,
		Browse: true,
	}))

	router.Start(port)
}
