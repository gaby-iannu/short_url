package handler

import (
	"fmt"
	"log"
	"net/http"
	"short_url/short"
	"short_url/short/cache"
	"short_url/short/model"
	"short_url/short/repository"

	"github.com/gin-gonic/gin"
)

var s short.Short 

// InitializeAndRun Short services and run rest controller
func InitializeAndRun(repository repository.Repository, cache cache.Cache) {
	s = short.New(cache, repository)
	router := gin.Default()
	router.POST("/tiny", createTinyUrl)
	router.GET("/long/:tiny", getUrl)
	router.Run("localhost:8080")
}

func createTinyUrl(c *gin.Context) {
	var url model.Url

	err := c.BindJSON(&url)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}

	if len(url.Long) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("Url too short!")})
		return
	}

	if len(url.User) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("User is empty")})
		return
	}

	tiny := s.Tiny(url)

	c.JSON(http.StatusCreated, gin.H{"tiny_url": tiny})
}

func getUrl(c *gin.Context) {
  	tiny := c.Param("tiny")
	if len(tiny) == 0 {		
		c.JSON(http.StatusBadRequest, gin.H{"msg":fmt.Errorf("Tiny url too short!")})
		return
	}

	longUrl, err := s.Get(tiny)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"long_url":longUrl})
}
