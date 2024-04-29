package handler

import (
	"fmt"
	"log"
	"net/http"
	"short_url/short"
	"short_url/short/cache"
	"short_url/short/model"
	"short_url/short/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var shorturlStatus *prometheus.CounterVec 

var s short.Short 


// InitializeAndRun Short services and run rest controller
func InitializeAndRun(repository repository.Repository, cache cache.Cache) *gin.Engine{
	s = short.New(cache, repository)

	shorturlStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "url_count",
			// Subsystem: "short",
			Help: "Count of status returned to user",
		}, 
		[]string{"url","status"},
	)
	prometheus.MustRegister(shorturlStatus)

	router := gin.Default()
	router.GET("/metrics", prometheusHandler)
	router.POST("/tiny", createTinyUrl)
	router.GET("/long/:tiny", getUrl)
	return router
}

func prometheusHandler(c *gin.Context) {
	log.Printf("prometheus\n")
	 promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

func createTinyUrl(c *gin.Context) {
	var url model.Url
	var status string
	var urlStr string = "empty"

	defer func() {
		log.Printf("inc - url:%s, status:%s\n", urlStr, status)
		shorturlStatus.WithLabelValues(urlStr, status).Inc()
	}()

	err := c.BindJSON(&url)
	if err != nil {
		status = strconv.Itoa(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("user and long url are required").Error()})
		return
	}

	if len(url.Long) <= 1 {
		urlStr = url.Long
		status = strconv.Itoa(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("Url too short!").Error()})
		return
	}

	if len(url.User) <= 1 {
		urlStr = url.Long
		status = strconv.Itoa(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("User is empty").Error()})
		return
	}

	tiny := s.Tiny(url)
	urlStr = url.Long
	status = strconv.Itoa(http.StatusCreated)
	c.JSON(http.StatusCreated, gin.H{"tiny_url": tiny})
}

func getUrl(c *gin.Context) {
  	tiny := c.Param("tiny")
	if len(tiny) <= 1 {		
		c.JSON(http.StatusBadRequest, gin.H{"msg":fmt.Errorf("Tiny url is too short!").Error()})
		return
	}

	longUrl, err := s.Get(tiny)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"long_url":longUrl})
}
