package handlers

import (
	"ammunition/config"

	"github.com/gin-gonic/gin"
)

func Info(c *gin.Context) {
	c.JSON(200, config.Info)
}

func Health(c *gin.Context) {
	c.JSON(200, gin.H{"Status": "UP"})
}
