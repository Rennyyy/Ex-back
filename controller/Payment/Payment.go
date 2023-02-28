package Payment

import (
	"context"
	"fmt"
	"time"

	"example.com/m/v2/keys"
	"example.com/m/v2/libs"
	"example.com/m/v2/models/auth"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PaymentPrivate(rn *echo.Group) {

	rn.POST("/sendOrder", sendOrder)
	rn.GET("/order/getorders", GetOrders)
	rn.GET("/order/gethistory", GetHistory)
	rn.POST("/order/update", updateStatus)

}

type OrderM struct {
	ID           primitive.ObjectID       `json:"_id" bson:"_id"`
	Order        []map[string]interface{} `json:"order" bson:"order`
	Total        float64                  `json:"total" bson:"total"`
	UserID       primitive.ObjectID       `json:"user_id" bson:"user_id`
	Status       string                   `json:"status" bson:"status"`
	CreateAtDate time.Time                `json:"create_at_date" bson:"create_at_date,omitempty"`
	UpdateAtDate time.Time                `json:"update_at_date" bson:"update_at_date,omitempty"`
}

type payloadOrder struct {
	Cart  []map[string]interface{} `json:"cart" bson:"cart`
	Total float64                  `json:"total" bson:"total"`
}

type bindOrder struct {
	ID     string `json:"_id" bson:"_id`
	Status string `json:"status" bson:"status`
}

func sendOrder(c echo.Context) error {
	var (
		ctx         = context.Background()
		_OrdersColl = libs.MongoDB.Database(keys.Database).Collection("Orders")
		payload     = payloadOrder{}
		// value       []bson.M
	)

	if err := c.Bind(&payload); err != nil {
		return c.JSON(400, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}
	fmt.Println("\n payload =>", payload)

	_userPayload := c.Get("user").(*jwt.Token)
	_claims := _userPayload.Claims.(*auth.JwtClaim)
	_uid, _ := primitive.ObjectIDFromHex(_claims.UID)

	id := primitive.NewObjectID()
	t := time.Now()
	insert := &OrderM{
		ID:           id,
		Order:        payload.Cart,
		Total:        payload.Total,
		Status:       "PROCESS",
		UserID:       _uid,
		CreateAtDate: t,
		UpdateAtDate: t,
	}

	if _, err := _OrdersColl.InsertOne(ctx, insert); err != nil {
		return c.JSON(400, echo.Map{"error": err.Error(), "message": "default_error_message"})
	}

	return c.JSON(200, echo.Map{
		"status": true,
		"result": "Add Order completed",
	})

}

func GetHistory(c echo.Context) error {

	var (
		ctx         = context.Background()
		_OrdersColl = libs.MongoDB.Database(keys.Database).Collection("Orders")
		payload     []bson.M
	)

	_userPayload := c.Get("user").(*jwt.Token)
	_claims := _userPayload.Claims.(*auth.JwtClaim)
	_role := _claims.Role
	_uid, _ := primitive.ObjectIDFromHex(_claims.UID)

	// fmt.Println("\n role ==>", _role)
	// fmt.Println("\n _uid ==>", _uid)

	pipe := []bson.M{
		bson.M{
			"$match": bson.M{
				"userid": _uid,
				"status": bson.M{
					"$in": bson.A{
						"CANCEL", "COMPLETED",
					},
				},
			},
		},
		bson.M{
			"$sort": bson.M{
				"create_at_date": 1,
			},
		},
	}

	if _role == "Admin" {
		pipe = []bson.M{
			bson.M{
				"$match": bson.M{
					"status": bson.M{
						"$in": bson.A{
							"CANCEL", "COMPLETED",
						},
					},
				},
			},
			bson.M{
				"$sort": bson.M{
					"create_at_date": 1,
				},
			},
		}
	}

	cursor, err := _OrdersColl.Aggregate(ctx, pipe)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":   err.Error(),
			"message": "default_error_message",
		})
	}

	if err := cursor.All(ctx, &payload); err != nil {
		return c.JSON(400, echo.Map{
			"error":   err.Error(),
			"message": "default_error_message",
		})
	}

	return c.JSON(200, echo.Map{
		"status": true,
		"result": payload,
	})
}

func GetOrders(c echo.Context) error {

	var (
		ctx         = context.Background()
		_OrdersColl = libs.MongoDB.Database(keys.Database).Collection("Orders")
		payload     []bson.M
	)

	_userPayload := c.Get("user").(*jwt.Token)
	_claims := _userPayload.Claims.(*auth.JwtClaim)
	_role := _claims.Role
	_uid, _ := primitive.ObjectIDFromHex(_claims.UID)

	// fmt.Println("\n role ==>", _role)
	// fmt.Println("\n _uid ==>", _uid)

	pipe := []bson.M{
		bson.M{
			"$match": bson.M{
				"userid": _uid,
				"status": "PROCESS",
			},
		},
		bson.M{
			"$sort": bson.M{
				"create_at_date": 1,
			},
		},
	}

	if _role == "Admin" {
		pipe = []bson.M{
			bson.M{
				"$match": bson.M{
					"status": "PROCESS",
				},
			},
			bson.M{
				"$sort": bson.M{
					"create_at_date": 1,
				},
			},
		}
	}

	cursor, err := _OrdersColl.Aggregate(ctx, pipe)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":   err.Error(),
			"message": "default_error_message",
		})
	}

	if err := cursor.All(ctx, &payload); err != nil {
		return c.JSON(400, echo.Map{
			"error":   err.Error(),
			"message": "default_error_message",
		})
	}

	return c.JSON(200, echo.Map{
		"status": true,
		"result": payload,
	})
}

func updateStatus(c echo.Context) error {
	var (
		ctx        = context.Background()
		_OrderColl = libs.MongoDB.Database(keys.Database).Collection("Orders")
		payload    = bindOrder{}
	)

	if err := c.Bind(&payload); err != nil {
		return c.JSON(400, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	_uid, _ := primitive.ObjectIDFromHex(payload.ID)

	fmt.Println("\n id ==>", _uid)
	fmt.Println("\n payload ==>", payload)
	filter := bson.M{
		"_id": _uid,
	}

	update := bson.M{
		"$set": bson.M{
			"status": payload.Status,
		},
	}
	_ = _OrderColl.FindOneAndUpdate(ctx, filter, update)

	return c.JSON(200, echo.Map{
		"status":  true,
		"message": "Success",
	})
}
