package main

import (
	// "fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Login struct{
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Data []struct {
		Id string `json:"id"`
		Username string `json:"username"`
		Password string `json:"password"`
		Fullname string `json:"fullname"`
		Norek string `json:"norek"`
		Created_at string `json:"created_at"`
		Updated_at string `json:"updated_at"`
	} `json:"data"`
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/login", func(c *gin.Context) {
		var jsonLogin Login
		if err := c.ShouldBindJSON(&jsonLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		password := []byte(jsonLogin.Password)
		
		userJson, _ := ioutil.ReadFile("json/user.json")
		var userData User
		json.Unmarshal(userJson, &userData)
		
		for _, val := range userData.Data{
			hashedPassword := []byte(val.Password)
			matchPass := bcrypt.CompareHashAndPassword(hashedPassword, password)
			if jsonLogin.Username != val.Username || matchPass != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
				return
			}
	
			c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
			return
		}

	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
