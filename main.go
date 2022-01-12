package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"fmt"
	
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

type Log struct {
	Error bool	`json:"error"`
	Message string	`json:"message"`
	Data string `json:"data"`
}

type LogHistory struct {
	Id int `json:"id"`
	Log string `json:"log"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

func CreateLogHistory (error bool, message, data string) {
	fileName := "json/log.json"
	file, err := ioutil.ReadFile(fileName)
		if err != nil {
			return
		}
	object := Log{error, message, data}
	jsonString, _ := json.Marshal(object)

	LogData := []LogHistory{}
	json.Unmarshal(file, &LogData)
	newStruct := &LogHistory{
		Id:         len(LogData) + 1,
		Log: string(jsonString),
		Created_at: time.Now().String(),
		Updated_at: time.Now().String(),
	}
	LogData = append(LogData, *newStruct)

	// Preparing the data to be marshalled and written.
	dataBytes, err := json.MarshalIndent(LogData, "", " ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(fileName, dataBytes, 0644)
	if err != nil {
		return
	}

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
		expirationTime := time.Now().Add(60 * time.Minute)

		userJson, _ := ioutil.ReadFile("json/user.json")
		var userData User
		json.Unmarshal(userJson, &userData)

		fmt.Println(userData)
		for _, val := range userData.Data {
			hashedPassword := []byte(val.Password)
			matchPass := bcrypt.CompareHashAndPassword(hashedPassword, password)
			if jsonLogin.Username != val.Username || matchPass != nil {
				CreateLogHistory(true, "unauthorized", "{}")
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

			msg := "User "+val.Fullname+" with user_id "+strconv.Itoa(val.Id)+" Login"
			CreateLogHistory(false, msg, "{}")
			c.JSON(http.StatusOK, gin.H{
				"status":  "you are logged in",
				"Token":   ss,
				"expires": expirationTime,
			})
			return
		}
	})

	// payment route
	r.POST("/payment", func(c *gin.Context) {
		// validasi jwt
		jwtKey, err := c.Cookie("jwtKey")
		if err != nil {
			CreateLogHistory(true, "unauthorized", "{}")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "Unauthorized",
			})
			return
		}

		tokenAuth := c.Request.Header["Authorization"]
		if len(tokenAuth) < 1 {
			CreateLogHistory(true, "Token Doesn't exists", "{}")
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
			CreateLogHistory(true, err.Error(), "{}")
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		claims, ok := token.Claims.(*UserClaims)
		if !ok || !token.Valid {
			CreateLogHistory(true, err.Error(), "{}")
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
		from := claims.Norek

		// pengurangan dan panambahan balance => ubah logic update data
		userFile := "json/user.json"

		userFileData, err := ioutil.ReadFile(userFile)
		if err != nil {
			return
		}

		userData := User{}
		userDataInput := []Data{}

		json.Unmarshal(userFileData, &userData)
		fmt.Println(userData.Data)

		norek_list := map[string]string{}

		for _, val := range userData.Data {
			norek_list[val.Norek] = val.Norek
		}

		fmt.Println(norek_list)

		for _, val := range userData.Data {
			if _, ok := norek_list[to]; !ok {
				CreateLogHistory(true, "nomor rekening tidak tersedia", "{}")
				c.JSON(http.StatusOK, gin.H{
					"message": "nomor rekening tidak tersedia",
				})
				return
		}
			
			balance := val.Balance
			if balance < nominal {
				CreateLogHistory(true, "dana anda tidak cukup", "{}")
				c.JSON(http.StatusOK, gin.H{
					"message": "dana anda tidak cukup",
				})
				return
			}

			fmt.Println(jsonPayment)
			if from == val.Norek{
				fmt.Println("from")
				balance -= nominal
			} else if to == val.Norek {
				fmt.Println("to")
				balance += nominal
			}
			newStruct := &Data{
				Id:         val.Id,
				Username:   val.Username,
				Password:   val.Password,
				Fullname:   val.Fullname,
				Norek:      val.Norek,
				Balance:    balance,
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

		CreateLogHistory(false, "success transfer payment", "{}")
		c.JSON(http.StatusOK, gin.H{
			"status": "success transfer payment",
			"data":   claims,
		})
	})

	r.GET("/logout", func(c *gin.Context) {
		c.SetCookie("jwtKey", "", -1, "/", "localhost", false, true)
		CreateLogHistory(false, "User Logout", "{}")
		c.JSON(http.StatusOK, gin.H{
			"status": "Cookies Delete",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
