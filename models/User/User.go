package User

import (
	"github.com/golang-jwt/jwt"
)

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