package controllers

import "github.com/gin-gonic/gin"

func OptionsUser(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}
