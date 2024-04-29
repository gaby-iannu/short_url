package handler

import (
	"fmt"
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
			Name: "shorturl_count_status",
			// Subsystem: "short",
			Help: "Count of status returned to user",
		}, 
		[]string{"count_by_status"},
	)
	prometheus.MustRegister(shorturlStatus)
	
	m = &metric{}

	router := gin.Default()
	router.GET("/metrics", prometheusHandler)
	router.POST("/tiny", createTinyUrl)
	router.GET("/long/:tiny", getUrl)
	return router
}

func prometheusHandler(c *gin.Context) {
	 promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

var count_metric = "count_"
var m *metric

type metric struct {
	lastMetric string
}

func (m *metric) buildCountMetric(status int) {
	stSatus := strconv.Itoa(status)
	m.lastMetric = fmt.Sprintf("%s%s", count_metric, stSatus)
}

func createTinyUrl(c *gin.Context) {
	var url model.Url

	defer func(){
		shorturlStatus.WithLabelValues(m.lastMetric).Inc()
	}()

	err := c.BindJSON(&url)
	if err != nil {
		m.buildCountMetric(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("user and long url are required").Error()})
		return
	}

	if len(url.Long) <= 1 {
		m.buildCountMetric(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("Url too short!").Error()})
		return
	}

	if len(url.User) <= 1 {
		m.buildCountMetric(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("User is empty").Error()})
		return
	}

	tiny := s.Tiny(url)
	m.buildCountMetric(http.StatusCreated)
	c.JSON(http.StatusCreated, gin.H{"tiny_url": tiny})
}

func getUrl(c *gin.Context) {
  	tiny := c.Param("tiny")
	if len(tiny) <= 1 {		
		m.buildCountMetric(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"msg":fmt.Errorf("Tiny url is too short!").Error()})
		return
	}

	longUrl, err := s.Get(tiny)
	if err != nil {
		m.buildCountMetric(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	m.buildCountMetric(http.StatusAccepted)
	c.JSON(http.StatusAccepted, gin.H{"long_url":longUrl})
}

