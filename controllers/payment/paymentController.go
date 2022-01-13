package payment

import(
	"fmt"
	"time"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"payment/models/User"
	"payment/models/Payment"
	"payment/controllers/log"
)

func PaymentProcess(c *gin.Context) {
	// validasi jwt
	jwtKey, err := c.Cookie("jwtKey")
	if err != nil {
		log.CreateLogHistory(true, "Forbidden", "{}")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Can't access this page",
			"data":    "{}",
		})
		return
	}

	// cek token from header
	tokenAuth := c.Request.Header["Authorization"]
	if len(tokenAuth) < 1 {
		log.CreateLogHistory(true, "Token Doesn't exists", "{}")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Token Doesn't exists",
			"data":    "{}",
		})
		return
	}

	// split token to get token without bearer
	tokenSplit := strings.Split(tokenAuth[0], " ")

	token, err := jwt.ParseWithClaims(tokenSplit[1], &User.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		log.CreateLogHistory(true, err.Error(), "{}")
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   true,
			"message": err.Error(),
			"data":    "{}",
		})
		return
	}

	// get data from Token
	claims, ok := token.Claims.(*User.UserClaims)
	if !ok || !token.Valid {
		log.CreateLogHistory(true, err.Error(), "{}")
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   true,
			"message": err.Error(),
			"data":    "{}",
		})
		return
	}

	// create payment process
	var jsonPayment Payment.PaymentTf
	if err := c.ShouldBindJSON(&jsonPayment); err != nil {
		if err.Error() == "EOF" {
			log.CreateLogHistory(true, "Please insert nominal and the destination account number", "{}")
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

	if nominal < 1 || len(to) < 1{
		log.CreateLogHistory(true, "Please insert nominal or the destination account number", "{}")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Please insert nominal or the destination account number",
				"data":    "{}",
			})
			return
	}

	// open file user.json
	userFile := "json/user.json"

	userFileData, err := ioutil.ReadFile(userFile)
	if err != nil {
		return
	}

	// preapre struct from models user
	userData := User.User{}
	userDataInput := []User.Data{}

	json.Unmarshal(userFileData, &userData)

	// listing norek and balance
	norek_list := map[string]string{}
	balance_list := map[int]int{}

	for _, val := range userData.Data {
		norek_list[val.Norek] = val.Norek
		balance_list[val.Id] = val.Balance
	}

	// check number account and balance
	if _, ok := norek_list[to]; !ok {
		log.CreateLogHistory(true, "account number not available", "{}")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "account number not available",
			"data":    "{}",
		})
		return
	}

	if to == from {
		log.CreateLogHistory(true, "can't transfer to own account", "{}")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "can't transfer to own account",
			"data":    "{}",
		})
		return
	}

	balance, _ := balance_list[claims.Id]
	if balance < nominal {
		log.CreateLogHistory(true, "your funds are not enough", "{}")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "your funds are not enough",
			"data":    "{}",
		})
		return
	}

	for _, val := range userData.Data {
		// calculate payment process
		balance = val.Balance
		if from == norek_list[val.Norek] {
			fmt.Println("from")
			balance -= nominal
		} else if to == norek_list[val.Norek] {
			fmt.Println("to")
			balance += nominal
		}
		// assign data to struct Data
		newStruct := &User.Data{
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

		// assign data to struct User
		UserDataObj := User.User{
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
	data := []Payment.PaymentTf{}

	// Here the magic happens!
	json.Unmarshal(file, &data)

	newStruct := &Payment.PaymentTf{
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

	log.CreateLogHistory(false, "success transfer payment", "{}")
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "success transfer payment",
		"data":    newStruct,
	})
}