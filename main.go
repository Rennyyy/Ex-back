package main

import (
	// "bufio"

	"fmt"
	"net/http"

	// "os"
	// "strconv"

	// "github.com/fxtlabs/primes"

	"example.com/m/v2/controller/Payment"
	Products "example.com/m/v2/controller/Products"
	Upload "example.com/m/v2/controller/Upload"
	UserM "example.com/m/v2/controller/Users"

	// char "example.com/m/v2/controller/Users"
	"example.com/m/v2/keys"
	"example.com/m/v2/libs"
	"example.com/m/v2/models/auth"

	// "github.com/labstack/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	// "github.com/gin-gonic/gin"
	// "gopkg.in/robfig/cron.v2"
)

func Authorize(ec echo.Context) error {
	return ec.JSON(http.StatusBadRequest, echo.Map{
		"message": "eeee.",
	})
}
func Token(ec echo.Context) error {

	return ec.JSON(http.StatusBadRequest, echo.Map{
		"status":  true,
		"message": "eeee.",
	})
}

func main() {

	fmt.Println("Now")

	libs.StartMongoService(keys.DB, keys.Database, "", true)
	libs.StartMongoServiceRENNY()

	r := echo.New()

	r.HideBanner = true
	r.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} ${uri} ${remote_ip} ${user_agent} status=${status} error:${error}\n",
	}))
	r.Use(middleware.Recover())
	r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderAccessControlAllowOrigin},
	}))

	var IsLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &auth.JwtClaim{},
		SigningKey: []byte("RENNYYY//"),
	})

	_privateAPI := r.Group("/private/api/v1", IsLoggedIn)
	_publicAPI := r.Group("/public/api/v1")

	r.POST("/Authorize", Authorize)
	r.POST("/Token", Token)
	// UserM.UserPublic(*gin.Engine)
	///////////////////// UserM
	UserM.UserPublic(_publicAPI)
	UserM.UserPrivate(_privateAPI)

	//////////////////// Product
	Products.ProductsPv(_privateAPI)

	//////////////////// Upload
	Upload.UploadPv(_privateAPI)

	//////////////////// Upload
	Payment.PaymentPrivate(_privateAPI)
	///////////////////
	// Character.CharPrivate(_privateAPI)
	r.Logger.Fatal(r.Start(keys.Port)) //can change port from ./keys/keys.go
	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"status": true,
	// 		// "ss":     databases,
	// 	})
	// })

	// r.Run(keys.Port)

}
