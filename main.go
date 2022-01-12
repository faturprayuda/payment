package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Data struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Fullname   string `json:"fullname"`
	Norek      string `json:"norek"`
	Balance    int    `json:"balance"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

type User struct {
	Data []Data `json:"data"`
}

type UserClaims struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Fullname   string `json:"fullname"`
	Norek      string `json:"norek"`
	Balance    int    `json:"balance"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
	jwt.StandardClaims
}

type PaymentTf struct {
	Id         int    `json:"id"`
	User_id    int    `json:"user_id"`
	From       string `json:"from"`
	To         string `json:"to"`
	Nominal    int    `json:"Nominal"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
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

		mySigningKey := []byte(jsonLogin.Username)
		password := []byte(jsonLogin.Password)

		// Declare the expiration time of the token
		// here, we have kept it as 5 minutes
		expirationTime := time.Now().Add(5 * time.Minute)

		userJson, _ := ioutil.ReadFile("json/user.json")
		var userData User
		json.Unmarshal(userJson, &userData)

		fmt.Println(userData)
		for _, val := range userData.Data {
			hashedPassword := []byte(val.Password)
			matchPass := bcrypt.CompareHashAndPassword(hashedPassword, password)
			if jsonLogin.Username != val.Username || matchPass != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
				return
			}

			fmt.Println(val)

			// Create the Claims
			claims := UserClaims{
				val.Id,
				val.Username,
				val.Fullname,
				val.Norek,
				val.Balance,
				val.Created_at,
				val.Updated_at,
				jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			ss, _ := token.SignedString(mySigningKey)

			c.SetCookie("jwtKey", jsonLogin.Username, 3600, "/", "localhost", false, true)

			c.JSON(http.StatusOK, gin.H{
				"status":  "you are logged in",
				"Token":   ss,
				"expires": expirationTime,
			})
			return
		}
	})

	r.POST("/payment", func(c *gin.Context) {
		// validasi jwt
		jwtKey, err := c.Cookie("jwtKey")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "Unauthorized",
			})
			return
		}

		tokenAuth := c.Request.Header["Authorization"]
		if len(tokenAuth) < 1 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Token Doesn't exists",
			})
			return
		}
		tokenSplit := strings.Split(tokenAuth[0], " ")

		token, err := jwt.ParseWithClaims(tokenSplit[1], &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		claims, ok := token.Claims.(*UserClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusOK, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			fmt.Println("err = ", err.Error())
			return
		}

		// payment
		var jsonPayment PaymentTf
		if err := c.ShouldBindJSON(&jsonPayment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		nominal := jsonPayment.Nominal
		to := jsonPayment.To

		// pengurangan dan panambahan balance => ubah logic update data
		userFile := "json/user.json"

		userFileData, err := ioutil.ReadFile(userFile)
		if err != nil {
			return
		}

		userData := User{}
		userDataInput := []Data{}

		json.Unmarshal(userFileData, &userData)
		for _, val := range userData.Data {
			newStruct := &Data{
				Id:         val.Id,
				Username:   val.Username,
				Password:   val.Password,
				Fullname:   val.Fullname,
				Norek:      val.Norek,
				Balance:    val.Balance - nominal,
				Created_at: val.Created_at,
				Updated_at: time.Now().String(),
			}
			userDataInput = append(userDataInput, *newStruct)
			UserDataObj := User{
				Data: userDataInput,
			}

			// Preparing the data to be marshalled and written.
			dataBytes, err := json.MarshalIndent(UserDataObj, "", " ")
			if err != nil {
				return
			}

			err = ioutil.WriteFile(userFile, dataBytes, 0644)
			if err != nil {
				return
			}
			// if jsonPayment.From == val.Norek {
			// 	newStruct := &Data{
			// 		Id:         val.Id,
			// 		Username:   val.Username,
			// 		Password:   val.Password,
			// 		Fullname:   val.Fullname,
			// 		Norek:      val.Norek,
			// 		Balance:    val.Balance - nominal,
			// 		Created_at: val.Created_at,
			// 		Updated_at: time.Now().String(),
			// 	}
			// 	userDataInput = append(userDataInput, *newStruct)
			// 	UserDataObj := User{
			// 		Data: userDataInput,
			// 	}

			// 	// Preparing the data to be marshalled and written.
			// 	dataBytes, err := json.MarshalIndent(UserDataObj, "", " ")
			// 	if err != nil {
			// 		return
			// 	}

			// 	err = ioutil.WriteFile(userFile, dataBytes, 0644)
			// 	if err != nil {
			// 		return
			// 	}
			// } else if jsonPayment.To == val.Norek {
			// 	newStruct := &Data{
			// 		Id:         val.Id,
			// 		Username:   val.Username,
			// 		Password:   val.Password,
			// 		Fullname:   val.Fullname,
			// 		Norek:      val.Norek,
			// 		Balance:    val.Balance + nominal,
			// 		Created_at: val.Created_at,
			// 		Updated_at: time.Now().String(),
			// 	}
			// 	userDataInput = append(userDataInput, *newStruct)

			// 	// Preparing the data to be marshalled and written.
			// 	dataBytes, err := json.MarshalIndent(userDataInput, "", " ")
			// 	if err != nil {
			// 		return
			// 	}

			// 	err = ioutil.WriteFile(userFile, dataBytes, 0644)
			// 	if err != nil {
			// 		return
			// 	}
			// } else {
			// 	newStruct := &Data{
			// 		Id:         val.Id,
			// 		Username:   val.Username,
			// 		Password:   val.Password,
			// 		Fullname:   val.Fullname,
			// 		Norek:      val.Norek,
			// 		Balance:    val.Balance + nominal,
			// 		Created_at: val.Created_at,
			// 		Updated_at: time.Now().String(),
			// 	}
			// 	userDataInput = append(userDataInput, *newStruct)

			// 	// Preparing the data to be marshalled and written.
			// 	dataBytes, err := json.MarshalIndent(userDataInput, "", " ")
			// 	if err != nil {
			// 		return
			// 	}

			// 	err = ioutil.WriteFile(userFile, dataBytes, 0644)
			// 	if err != nil {
			// 		return
			// 	}
			// }
		}

		// record transfer
		filename := "json/payment.json"

		file, err := ioutil.ReadFile(filename)
		if err != nil {
			return
		}
		data := []PaymentTf{}

		// Here the magic happens!
		json.Unmarshal(file, &data)

		newStruct := &PaymentTf{
			Id:         len(data) + 1,
			User_id:    claims.Id,
			From:       claims.Norek,
			To:         to,
			Nominal:    nominal,
			Created_at: time.Now().String(),
			Updated_at: time.Now().String(),
		}

		data = append(data, *newStruct)

		// Preparing the data to be marshalled and written.
		dataBytes, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			return
		}

		err = ioutil.WriteFile(filename, dataBytes, 0644)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success transfer payment",
			"data":   claims,
		})
	})

	r.GET("/logout", func(c *gin.Context) {
		c.SetCookie("jwtKey", "", -1, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{
			"status": "Cookies Delete",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
