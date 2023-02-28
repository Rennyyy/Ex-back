package Products

import (
	// "fmt"
	"context"
	"fmt"
	"strconv"
	"time"

	"crypto/rand"

	"example.com/m/v2/keys"
	"example.com/m/v2/libs"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
	// _Upload "example.com/m/v2/controller/Upload"
)

type ModelAddProduct struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Discliption  string             `json:"discliption"`
	Price        float64            `json:"price" bson:"price"`
	Active       bool               `json:"active"`
	Image        string             `json:"image,omitempty" bson:"image,omitempty"`
	CreateAtDate time.Time          `json:"create_at_date" bson:"create_at_date,omitempty"`
	UpdateAtDate time.Time          `json:"update_at_date" bson:"update_at_date,omitempty"`
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func AddProducts(_pr echo.Context) error {
	var (
		product = new(ModelAddProduct)
		ctx     = context.Background()
		pdCol   = libs.MongoDB.Database(keys.Database).Collection("Products")
		t       = time.Now()
	)

	product.ID = primitive.NewObjectID()
	product.Name = _pr.FormValue("name")
	product.Discliption = _pr.FormValue("discliption")

	f := _pr.FormValue("price")
	fmt.Println("price ==>", f)
	fmt.Println("name ==>", product.Name)
	_price, err := strconv.ParseFloat(f, 8)
	if err != nil {
		return _pr.JSON(400, echo.Map{"error": err.Error(), "message": "default_error_message"})
	}
	product.Price = _price

	product.Active = true
	product.UpdateAtDate = t

	product.Image = _pr.FormValue("image")

	if _, err := pdCol.InsertOne(ctx, product); err != nil {
		return _pr.JSON(400, echo.Map{"error": err.Error(), "message": "default_error_message"})
	}

	return _pr.JSON(200, echo.Map{
		"status": true,
		"detail": product,
		"result": "Add Product Completed",
	})

}

func EditProduct(_pr echo.Context) error {

	var (
		ctx   = context.Background()
		pdCol = libs.MongoDB.Database(keys.Database).Collection("Products")
		t     = time.Now()
	)

	_pid, _ := primitive.ObjectIDFromHex(_pr.Param("id"))
	fmt.Println("\n _pr.QueryParam() ==> ", _pr.Param("id"))

	name := _pr.FormValue("name")
	discliption := _pr.FormValue("discliption")

	f := _pr.FormValue("price")
	price, err := strconv.ParseFloat(f, 8)
	if err != nil {
		return _pr.JSON(400, echo.Map{"error": err.Error(), "message": "default_error_message"})
	}

	image := _pr.FormValue("image")

	filter := bson.M{
		"_id": _pid,
	}

	update := bson.M{
		"$set": bson.M{
			"name":           name,
			"price":          price,
			"discliption":    discliption,
			"image":          image,
			"update_at_date": t,
		},
	}

	fmt.Println("\n update ==> ", update)
	fmt.Println("\n filter ==> ", filter)
	fmt.Println("\n _pr.QueryParam() ==> ", _pr.QueryParam("id"))
	_ = pdCol.FindOneAndUpdate(ctx, filter, update)

	return _pr.JSON(200, echo.Map{
		"status": true,
		"result": "Edit Completed",
	})
}

func testparam(e echo.Context) error {

	c := e.Param("id")
	fmt.Println("\n _pr.QueryParam() ==> ", c)

	return e.JSON(200, echo.Map{
		"status": true,
		"result": "Edit Completed",
	})
}

func Products(_pr echo.Context) error {
	var (
		payload []bson.M
		ctx     = context.Background()
		pdCol   = libs.MongoDB.Database(keys.Database).Collection("Products")
	)

	pipe := []bson.M{
		bson.M{
			"$sort": bson.M{
				"update_at_date": -1,
			},
		},
	}

	cursor, err := pdCol.Aggregate(ctx, pipe)
	if err != nil {
		fmt.Println(err.Error())
		return _pr.JSON(400, echo.Map{
			"status":  false,
			"message": "cannot get booking list",
		})
	}

	if err := cursor.All(ctx, &payload); err != nil {
		fmt.Println(err.Error())
		return _pr.JSON(400, echo.Map{
			"status":  false,
			"message": "cannot get booking list",
		})
	}

	return _pr.JSON(200, echo.Map{
		"status":     true,
		"result":     payload,
		"len_result": len(payload),
	})

}

func ProductsPv(rn *echo.Group) {

	rn.POST("/addproducts", AddProducts)
	rn.PUT("/editproducts/:id", EditProduct)

	rn.GET("/Products", Products)

	rn.POST("/testparam/:id", testparam)

}
