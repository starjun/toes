package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/starjun/jobrunner"
	"log"
	"toes/internal/middleware"
)

// ------
func InitJob() {
	jobrunner.Start() // optional: jobrunner.Start(pool int, concurrent int) (10, 1)
	jobrunner.Schedule("@every 10s", ReminderEmails{Test: "xxxxx1"}, "xxxxx")
}

// Job Specific Functions
type ReminderEmails struct {
	// filtered
	Test string
	Cnt  int
}

var (
	Cnt int
)

// ReminderEmails.Run() will get triggered automatically.
func (e ReminderEmails) Run() {
	// Queries the DB
	// Sends some email
	log.Println("Every 10 sec send reminder emails \n", e.Test, e.Cnt)
	e.Cnt++
}

func TestRun() error {

	r := gin.New()
	mws := []gin.HandlerFunc{
		middleware.Logger(),
		gin.Recovery(),
		middleware.NoCache,
		middleware.Cors,
		middleware.Secure,
		middleware.RequestID(),
	}

	r.Use(mws...)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	InitJob()

	return r.Run(":8080")
}

//---------
