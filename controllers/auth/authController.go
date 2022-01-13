package auth

import(
	"time"
	"strconv"
	"reflect"
	"net/http"
	"io/ioutil"
	"encoding/json"
	
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	
	"payment/models/Auth"
	"payment/models/User"
	"payment/controllers/log"
)

// handler login
func Login(c *gin.Context) {
	// check authenticaction
	var jsonLogin Auth.Login
	if err := c.ShouldBindJSON(&jsonLogin); err != nil {
		if err.Error() == "EOF" {
			log.CreateLogHistory(true, "blank username and password", "{}")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "blank username and password",
				"data":    "{}",
			})
		}
		return
	}

	// assign the key JWT
	mySigningKey := []byte(jsonLogin.Username)

	// Duration of JWT in minute
	expirationTime := time.Now().Add(60 * time.Minute)

	// open file user.json
	fileName := "json/user.json"
	userJson, _ := ioutil.ReadFile(fileName)
	var userData User.User
	json.Unmarshal(userJson, &userData)

	// listing auth user
	authList := map[string]interface{}{}
	auth := map[string]string{}
	for _, val := range userData.Data {
		authList[val.Username] = val
		auth[val.Username] = val.Username
		auth[val.Password] = val.Password
	}

	for _, val := range userData.Data {
		usernameData, usernameFound := auth[jsonLogin.Username]
		hashedPassword, _ := auth[val.Password]
		bcryptPass := []byte(hashedPassword)
		pass := []byte(jsonLogin.Password)
		matchPass := bcrypt.CompareHashAndPassword(bcryptPass, pass)

		// check auth user is exists
		if !usernameFound || matchPass != nil {
			log.CreateLogHistory(true, "username or password wrong", "{}")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "username or password wrong",
				"data":    "{}",
			})
			return
		}

		// get data user from interface authList
		reflectValue := reflect.ValueOf(authList[usernameData])

		Id := reflectValue.Field(0).Interface().(int)
		Username := reflectValue.Field(1).Interface().(string)
		Fullname := reflectValue.Field(3).Interface().(string)
		Norek := reflectValue.Field(4).Interface().(string)
		Balance := reflectValue.Field(5).Interface().(int)
		Created_at := reflectValue.Field(6).Interface().(string)
		Updated_at := reflectValue.Field(7).Interface().(string)

		// set value claims to JWT
		claims := User.UserClaims{
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

		// generate Token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, _ := token.SignedString(mySigningKey)

		// save key JWT in cookies
		c.SetCookie("jwtKey", jsonLogin.Username, 3600, "/", "localhost", false, true)

		msg := "User " + Fullname + " with user_id " + strconv.Itoa(Id) + " Login"
		log.CreateLogHistory(false, msg, "{}")
		c.JSON(http.StatusOK, gin.H{
			"Token":   ss,
			"expires": expirationTime,
			"error":   false,
			"message": "Success Login",
			"data":    claims,
		})
		return
	}
}

func Logout(c *gin.Context) {
	// delete cookies key JWT
	c.SetCookie("jwtKey", "", -1, "/", "localhost", false, true)
	log.CreateLogHistory(false, "User Logout", "{}")
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "User Logout",
	})
}