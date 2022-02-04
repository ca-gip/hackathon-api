package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"net/http"
	"time"
)

var queryCache = cache.New(5*time.Minute, 15*time.Minute)

func CacheContent() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, queryCache.Items())
	}

}

func CacheKeys() gin.HandlerFunc {
	return func(c *gin.Context) {
		keys := make([]string, 0)
		for key, _ := range queryCache.Items() {
			keys = append(keys, key)
		}
		c.JSON(http.StatusOK, keys)

	}

}
