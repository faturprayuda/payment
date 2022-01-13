package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
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

type Log struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type LogHistory struct {
	Id         int    `json:"id"`
	Log        string `json:"log"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

func CreateLogHistory(error bool, message, data string) {
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
		Log:        string(jsonString),
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

	// route tes endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// route login
	r.POST("/login", func(c *gin.Context) {
		var jsonLogin Login
		if err := c.ShouldBindJSON(&jsonLogin); err != nil {
			if err.Error() == "EOF" {
				CreateLogHistory(true, "blank username and password", "{}")
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   true,
					"message": "blank username and password",
					"data":    "{}",
				})
			}
			return
		}

		mySigningKey := []byte(jsonLogin.Username)

		// Declare the expiration time of the token
		// here, we have kept it as 5 minutes
		expirationTime := time.Now().Add(5 * time.Minute)

		// open file user.json
		userJson, _ := ioutil.ReadFile("json/user.json")
		var userData User
		json.Unmarshal(userJson, &userData)

		auth := map[string]string{}
		for _, val := range userData.Data {
			auth[val.Username] = val.Username
			auth[val.Password] = val.Password
		}

		authList := map[string]interface{}{}
		for _, val := range userData.Data {
			authList[val.Username] = val
		}

		for _, val := range userData.Data {
			usernameData, usernameFound := auth[jsonLogin.Username]
			hashedPassword, _ := auth[val.Password]
			bcryptPass := []byte(hashedPassword)
			pass := []byte(jsonLogin.Password)
			matchPass := bcrypt.CompareHashAndPassword(bcryptPass, pass)
			if !usernameFound || matchPass != nil {
				CreateLogHistory(true, "username or password wrong", "{}")
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   true,
					"message": "username or password wrong",
					"data":    "{}",
				})
				return
			}

			// Create the Claims
			reflectValue := reflect.ValueOf(authList[usernameData])

			Id := reflectValue.Field(0).Interface().(int)
			Username := reflectValue.Field(1).Interface().(string)
			Fullname := reflectValue.Field(3).Interface().(string)
			Norek := reflectValue.Field(4).Interface().(string)
			Balance := reflectValue.Field(5).Interface().(int)
			Created_at := reflectValue.Field(6).Interface().(string)
			Updated_at := reflectValue.Field(7).Interface().(string)

			// set value claims
			claims := UserClaims{
				Id,
				Username,
				Fullname,
				Norek,
				Balance,
				Created_at,
				Updated_at,
				jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			ss, _ := token.SignedString(mySigningKey)

			c.SetCookie("jwtKey", jsonLogin.Username, 3600, "/", "localhost", false, true)

			msg := "User " + Fullname + " with user_id " + strconv.Itoa(Id) + " Login"
			CreateLogHistory(false, msg, "{}")
			c.JSON(http.StatusOK, gin.H{
				"Token":   ss,
				"expires": expirationTime,
				"error":   false,
				"message": "Success Login",
				"data":    claims,
			})
			return
		}
	})

	// payment route
	r.POST("/payment", func(c *gin.Context) {
		// validasi jwt
		jwtKey, err := c.Cookie("jwtKey")
		if err != nil {
			CreateLogHistory(true, "Forbidden", "{}")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Can't access this page",
				"data":    "{}",
			})
			return
		}

		tokenAuth := c.Request.Header["Authorization"]
		if len(tokenAuth) < 1 {
			CreateLogHistory(true, "Token Doesn't exists", "{}")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Token Doesn't exists",
				"data":    "{}",
			})
			return
		}
		tokenSplit := strings.Split(tokenAuth[0], " ")

		token, err := jwt.ParseWithClaims(tokenSplit[1], &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		if err != nil {
			CreateLogHistory(true, err.Error(), "{}")
			c.JSON(http.StatusBadGateway, gin.H{
				"error":   true,
				"message": err.Error(),
				"data":    "{}",
			})
			return
		}

		claims, ok := token.Claims.(*UserClaims)
		if !ok || !token.Valid {
			CreateLogHistory(true, err.Error(), "{}")
			c.JSON(http.StatusBadGateway, gin.H{
				"error":   true,
				"message": err.Error(),
				"data":    "{}",
			})
			return
		}

		// payment
		var jsonPayment PaymentTf
		if err := c.ShouldBindJSON(&jsonPayment); err != nil {
			if err.Error() == "EOF" {
				CreateLogHistory(true, "Please insert nominal and the destination account number", "{}")
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   true,
					"message": "Please insert nominal and the destination account number",
					"data":    "{}",
				})
			}
			return
		}

		nominal := jsonPayment.Nominal
		to := jsonPayment.To
		from := claims.Norek

		// pengurangan dan panambahan balance
		userFile := "json/user.json"

		userFileData, err := ioutil.ReadFile(userFile)
		if err != nil {
			return
		}

		userData := User{}
		userDataInput := []Data{}

		json.Unmarshal(userFileData, &userData)

		norek_list := map[string]string{}
		for _, val := range userData.Data {
			norek_list[val.Norek] = val.Norek
		}

		balance_list := map[int]int{}
		for _, val := range userData.Data {
			balance_list[val.Id] = val.Balance
		}

		// pengecekan rekening dan balance
		if _, ok := norek_list[to]; !ok {
			CreateLogHistory(true, "account number not available", "{}")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   true,
				"message": "account number not available",
				"data":    "{}",
			})
			return
		}

		if to == from {
			CreateLogHistory(true, "can't transfer to own account", "{}")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   true,
				"message": "can't transfer to own account",
				"data":    "{}",
			})
			return
		}

		balance, _ := balance_list[claims.Id]
		if balance < nominal {
			CreateLogHistory(true, "your funds are not enough", "{}")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "your funds are not enough",
				"data":    "{}",
			})
			return
		}

		for _, val := range userData.Data {
			// kalkulasi
			balance = val.Balance
			if from == norek_list[val.Norek] {
				fmt.Println("from")
				balance -= nominal
			} else if to == norek_list[val.Norek] {
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
			"error":   false,
			"message": "success transfer payment",
			"data":    newStruct,
		})
	})

	r.GET("/logout", func(c *gin.Context) {
		c.SetCookie("jwtKey", "", -1, "/", "localhost", false, true)
		CreateLogHistory(false, "User Logout", "{}")
		c.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "User Logout",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
