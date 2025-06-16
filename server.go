package main

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/events"
	"github.com/byuoitav/lab-attendance/cache"
	"github.com/byuoitav/lab-attendance/eventforwarder"
	"github.com/byuoitav/lab-attendance/handlers"
	"github.com/byuoitav/lab-attendance/lab"
	"github.com/byuoitav/lab-attendance/messenger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//go:embed web/*
var embeddedFiles embed.FS

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

	port := ":8243"

	// Setup the Frontend
	subFS, err := fs.Sub(embeddedFiles, "web")
	if err != nil {
		log.L.Fatal("failed to create sub filesystem for web files", zap.Error(err))
	}

	router := gin.Default()

	// API routes
	router.POST("/api/v1/login/:byuID", func(c *gin.Context) {
		handlers.Login(msgr, deviceInfo, lab)(c)
	})
	router.GET("/api/v1/config", func(c *gin.Context) {
		handlers.GetConfig(config)(c)
	})
	router.GET("/websocket", func(c *gin.Context) {
		ef.HandleWebsocket(c)
	})

	// Serve static files from embedded FS under /web/
	webFS := http.FS(subFS)
	router.GET("/web/*filepath", func(c *gin.Context) {
		file := strings.TrimPrefix(c.Param("filepath"), "/")
		if file == "" || file == "/" {
			file = "index.html"
		}
		f, err := webFS.Open(file)
		if err != nil {
			// fallback to index.html for SPA routes
			f, err = webFS.Open("index.html")
			if err != nil {
				c.String(http.StatusInternalServerError, "index.html not found")
				return
			}
		}
		defer f.Close()
		stat, err := f.Stat()
		if err != nil {
			c.String(http.StatusInternalServerError, "could not stat file")
			return
		}
		content, err := io.ReadAll(f)
		if err != nil {
			c.String(http.StatusInternalServerError, "could not read file")
			return
		}
		http.ServeContent(c.Writer, c.Request, file, stat.ModTime(), bytes.NewReader(content))
	})

	// Redirect root to /web/
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/web/")
	})

	// Custom 404 handler
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/web/") {
			// Serve index.html for SPA routes
			f, err := webFS.Open("index.html")
			if err != nil {
				c.String(http.StatusInternalServerError, "index.html not found")
				return
			}
			defer f.Close()
			stat, err := f.Stat()
			if err != nil {
				c.String(http.StatusInternalServerError, "could not stat index.html")
				return
			}
			content, err := io.ReadAll(f)
			if err != nil {
				c.String(http.StatusInternalServerError, "could not read index.html")
				return
			}
			http.ServeContent(c.Writer, c.Request, "index.html", stat.ModTime(), bytes.NewReader(content))
		} else {
			c.String(http.StatusNotFound, "Not found")
			log.L.Error("404 Not Found", zap.String("path", path))
		}
	})

	router.Run(port)
}
