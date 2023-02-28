package Character

import (
	// "context"
	// "fmt"

	// "example.com/m/v2/keys"
	// "example.com/m/v2/libs"
	// "example.com/m/v2/model"
	"github.com/labstack/echo/v4"
	// "gopkg.in/mgo.v2/bson"
)

func CharPrivate(rn *echo.Group) {

	// rn.POST("/creatCharacter", createChar)

}

func createChar(_ch echo.Context) error {
	// var (
	// 	payload       = model.Char_Model{}
	// 	_charModelCol = libs.MongoDB.Database(keys.Database).Collection("CharModels")
	// 	_charCol      = libs.MongoDB.Database(keys.Database).Collection("Char")
	// 	ctx           = context.Background()
	// 	dulname       bson.M
	// 	charM         bson.M
	// 	ss            = new(model.Char_Model)
	// )

	// if err := _ch.Bind(&payload); err != nil {
	// 	return _ch.JSON(400, echo.Map{
	// 		"status": false,
	// 		"result": err.Error(),
	// 	})
	// }

	// filter := bson.M{
	// 	"name": payload.NAME,
	// }

	// filterM := bson.M{
	// 	"No": 1,
	// }

	// fmt.Println("pay ==> ", payload)

	// _ = _charCol.FindOne(ctx, filter).Decode(&dulname)

	// if dulname != nil {
	// 	return _ch.JSON(400, echo.Map{
	// 		"status": false,
	// 		"result": "Dul name",
	// 	})
	// }

	// _ = _charModelCol.FindOne(ctx, filterM).Decode(&charM)

	// fmt.Println(payload)
	// fmt.Println("charM ==>", charM["INT"])
	// s := charM["INT"]

	// ss.INT = s
	// payload.INT = 1

	return _ch.JSON(200, echo.Map{
		"status": true,
		"result": "ssss.",
	})

}
