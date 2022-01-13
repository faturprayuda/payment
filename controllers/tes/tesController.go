package tes

import(
	"github.com/gin-gonic/gin"
)

func Tes(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}