package main

import (
	"net/http"
	"rest-example/db"

	"github.com/gin-gonic/gin"
	"rest-example/models"
)

func main() {
  db.InitDB()
	router := gin.Default()

  router.GET("/events", getEvents)
  router.POST("/events", createEvents)


	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.Run(":8080") // listens on port 8080 by default
}

func getEvents(c *gin.Context){
  events, err := models.GetAllEvents()

  if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch events"})
    return
	}
  c.JSON(http.StatusOK, events)
}

func createEvents(c *gin.Context){
  var event models.Event
  err := c.ShouldBindJSON(&event)

  if err != nil{
    c.JSON(http.StatusBadRequest, gin.H{"message": "could not parse the request"})
    return
  }

  err = event.Save()

  if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to save events"})
    return
	}
  c.JSON(http.StatusCreated, gin.H{"message": "Event created", "event": event})
}