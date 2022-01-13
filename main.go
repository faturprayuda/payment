package main

import (

	"github.com/gin-gonic/gin"

	"payment/controllers/tes"
	"payment/controllers/auth"
	"payment/controllers/payment"
)

func main() {
	r := gin.Default()

	// endpoint tes
	r.GET("/ping", tes.Tes)

	// endpoint login
	r.POST("/login", auth.Login)

	// endpoint payment
	r.POST("/payment", payment.PaymentProcess)

	// endpoint logout
	r.GET("/logout", auth.Logout)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
