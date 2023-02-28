package UserM

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"example.com/m/v2/helpers"
	"example.com/m/v2/keys"
	"example.com/m/v2/libs"
	"example.com/m/v2/model"
	"example.com/m/v2/models/auth"

	// "github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt"
	// "github.com/golang-jwt/jwt/v4"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/cloudinary/cloudinary-go/v2"
	// "github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"

	"gopkg.in/mgo.v2/bson"
)

// "github.com/gin-gonic/gin"

func UserPublic(rn *echo.Group) {

	rn.POST("/addUser", AddUserManagement)
	rn.POST("/login", loginUser)

}

func UserPrivate(rn *echo.Group) {

	rn.GET("/profile", getuserM)
	rn.PUT("/createC", createC)
	rn.PUT("/updateP", updateexp)
	rn.GET("/user/infor", checkUser)
	rn.GET("/user/Char", showChar)
	rn.GET("/user/datailChar/:id", datailChar)
	rn.PATCH("/addToCart", addToCart)

}

type ModelAddUserM struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	Email           string             `json:"email"`
	Username        string             `json:"username"`
	Password        string             `json:"password"`
	ConfirmPassword string             `json:"confirm_password"`
	Phone           string             `json:"phone"`

	Cart            []string `json:"cart,omitempty" bson:"cart"`
	Is_verify_email bool     `json:"is_verify_email"`
	Active          bool     `json:"active"`
	Block           bool     `json:"block"`

	Status Status `json:"status,omitempty" bson:"status"`
	Image  string `json:"image,omitempty" bson:"image,omitempty"`

	RoleID primitive.ObjectID `json:"role_id"`
	Role   string             `json:"role"`

	CreateAtDate time.Time `json:"create_at_date" bson:"create_at_date,omitempty"`
}

// type Cart struct {
// 	name bool `json:"name" bson:"name"`
// }

type Status struct {
	Online    bool      `json:"online,omitempty" bson:"online"`
	LastLogin LastLogin `json:"last_login,omitempty" bson:"last_login"`
}

type LastLogin struct {
	Date      time.Time `json:"date,omitempty" bson:"date"`
	IpAddr    string    `json:"ip_address,omitempty" bson:"ip_address"`
	UserAgent string    `json:"user_agent,omitempty" bson:"user_agent"`
}

type Roles struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}

type Login struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type Cart struct {
	Cart []map[string]interface{} `json:"cart" bson:"cart`
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

func GetOutboundIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}

func AddUserManagement(c echo.Context) error {
	var (
		userM           = new(ModelAddUserM)
		value           []bson.M
		dulUser         bson.M
		ctx             = context.Background()
		collection      = libs.MongoDB.Database(keys.Database).Collection("Users")
		role_collection = libs.MongoDB.Database(keys.Database).Collection("Role")
	)

	userM.Email = c.FormValue("email")
	userM.Username = c.FormValue("username")
	// userM.Password = c.FormValue("password")
	// userM.ConfirmPassword = c.FormValue("confirm_password")
	password := c.FormValue("password")
	confpassword := c.FormValue("confirm_password")
	userM.Phone = c.FormValue("phone")

	image, errImage := c.FormFile("image")
	if errImage != nil {
		fmt.Println(errImage)
	}

	collection.FindOne(context.TODO(), bson.M{
		"username": userM.Username,
	}).Decode(&dulUser)

	if len(dulUser) != 0 {
		// fmt.Println("This username is already taken.")
		return c.JSON(422, echo.Map{
			"status":  false,
			"message": "This username is already taken.",
		})
	}

	collection.FindOne(context.TODO(), bson.M{
		"email": userM.Email,
	}).Decode(&dulUser)

	if len(dulUser) != 0 {
		// fmt.Println("This email is already taken.")
		return c.JSON(422, echo.Map{
			"status":  false,
			"message": "This email is already taken.",
		})
	}

	userM.CreateAtDate = time.Now()
	userM.Active = true
	userM.Block = false
	userM.Status.LastLogin.Date = time.Now()
	userM.Cart = []string{}

	role := c.FormValue("role_id")

	if password == confpassword {
		_h := sha256.Sum256([]byte(password))
		_encodedHex := hex.EncodeToString(_h[:])
		PasswordHex, _ := bcrypt.GenerateFromPassword([]byte(_encodedHex), 10)

		userM.Password = string(PasswordHex)
		userM.ConfirmPassword = string(PasswordHex)
	} else {
		return c.JSON(400, echo.Map{"message": "passnotmatch"})
	}

	//////////////////// Hex ObjectID
	userM.ID = primitive.NewObjectID()
	roleID, err := primitive.ObjectIDFromHex(role)
	if err != nil {
		return c.JSON(422, echo.Map{"error": "role not match", "message": "default_error_message"})
	}
	userM.RoleID = roleID

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"_id": roleID,
			},
		},
	}

	cursor, err := role_collection.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(400, echo.Map{
			"status":  false,
			"message": "Cannot Aggregatet",
		})
	}

	if err := cursor.All(ctx, &value); err != nil {
		fmt.Println(err.Error())
		return c.JSON(400, echo.Map{
			"status":  false,
			"message": "cannot Aggregatet",
		})
	}
	fmt.Println(value)

	for _, v := range value {
		roles := v["name"].(bson.M)["en"].(string)
		userM.Role = roles
	}
	//////////////////// Image Upload

	// cld, _ := cloudinary.New()
	cld, err3 := cloudinary.NewFromURL(keys.CLOUDINARY_URL)
	if err3 != nil {
		return c.JSON(500, echo.Map{"error": err.Error(), "message": "default_error_message"})
	}

	if errImage == nil {
		ImageName, err := GenerateRandomString(32)
		if err != nil {
			return c.JSON(400, echo.Map{"error": err.Error(), "message": "default_error_message"})
		}
		fmt.Println(ImageName)

		src, err := image.Open()
		if err != nil {
			return c.JSON(400, echo.Map{"error": err.Error(), "message": "default_error_message"})
		}
		defer src.Close()

		resp, _ := cld.Upload.Upload(ctx, src, uploader.UploadParams{PublicID: ImageName, Folder: "Rennyyy2"})

		fmt.Println("resp ==> ", resp.SecureURL)
		userM.Image = resp.SecureURL

	}

	if _, err := collection.InsertOne(context.Background(), userM); err != nil {
		return c.JSON(400, echo.Map{"error": err.Error(), "message": "default_error_message"})
	}

	return c.JSON(200, echo.Map{
		"status": true,
		"detail": userM,
		"result": "Add user management completed",
	})
}

func loginUser(c echo.Context) error {
	clientRoleID, _ := primitive.ObjectIDFromHex("63fa0575b9702faca94cfa89")
	AdminRoleID, _ := primitive.ObjectIDFromHex("63e8bb9af5e3c520aa09379c")
	arr := []primitive.ObjectID{clientRoleID, AdminRoleID}

	return loginAuthentication(c, arr)

}

func checkUser(c echo.Context) error {

	_userPayload := c.Get("user").(*jwt.Token)
	_claims := _userPayload.Claims.(*auth.JwtClaim)
	_uid := _claims.UID

	fmt.Println(_uid)
	fmt.Println("_claims ==>", _claims)
	fmt.Println("\n\n")

	// tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI2M2U4ZmEwYzBiYmJmNTNmMGRiYWFjZjciLCJyb2xlIjoiQWRtaW4iLCJleHAiOjE2NzY1MzQ5Mjh9.0WqcCWdRZpscT4S0N1h4QMflOYsm9DEanAXzN9qsG7o"
	// claims := jwt.MapClaims{}
	// token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte("RENNYYY//>"), nil
	// })
	// // ... error handling

	// fmt.Println(token)
	// // do something with decoded claims
	// for key, val := range claims {
	// 	fmt.Printf("Key: %v, value: %v\n", key, val)
	// }

	return c.JSON(http.StatusOK, echo.Map{
		"status": true,
		"data":   _claims,
		// "_role":  _role,
	})

}

func loginAuthentication(c echo.Context, roleID []primitive.ObjectID) error {

	var payload = Login{}
	user := model.Users_R{}
	role := model.Role_R{}

	_userCollection := libs.MongoDB.Database(keys.Database).Collection("Users")
	_RoleCollection := libs.MongoDB.Database(keys.Database).Collection("Role")
	_HistoryCollection := libs.MongoDB.Database(keys.Database).Collection("Login_Histories")
	inputUser := new(model.Users_R)

	if err := c.Bind(&payload); err != nil {
		return c.JSON(400, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	// inputUser.Username = c.FormValue("username")
	// inputUser.Password = c.FormValue("password")
	inputUser.Username = payload.Username
	inputUser.Password = payload.Password

	// Email := strings.ToLower(inputUser.Email)
	errFindEmail := _userCollection.FindOne(context.TODO(), bson.M{
		// "email": bson.M{
		// 	"$regex": inputUser.Email, "$options": "i",
		// },
		// "block":  false,
		"username": inputUser.Username,
		"block":    false,
		"active":   true,
	}).Decode(&user)

	if errFindEmail != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  false,
			"message": "This username doesn't exist.",
		})
	}

	if helpers.Verify(user.Password, inputUser.Password) == true {

		RoleID, _ := primitive.ObjectIDFromHex(user.RoleID)
		_ = _RoleCollection.FindOne(context.TODO(), bson.M{
			"_id": RoleID,
		}).Decode(&role)
		fmt.Println(RoleID)
		fmt.Println(role)
		// if errFindRole == nil {
		// 	return c.JSON(http.StatusBadRequest, echo.Map{
		// 		"status":  false,
		// 		"message": "wrongrole",
		// 	})
		// }

		validateRole := false
		for _, v := range roleID {
			if RoleID == v {
				validateRole = true
			}
		}
		if validateRole == false {
			return c.JSON(200, echo.Map{
				"status":  false,
				"message": "wrongrole",
			})
		}

		Ip := GetOutboundIP()

		user_agent := c.Request().Header.Get("User-Agent")

		userid := user.UID
		////////////////////////////////////////////////////////////////////////////////
		_userid, _ := primitive.ObjectIDFromHex(userid)
		var value []bson.M
		_value := []bson.M{
			bson.M{
				"$match": bson.M{
					"user_id": _userid,
				},
			},
		}

		_valueinsert := bson.M{
			"user_id": _userid,
			"login_histories": []bson.M{
				bson.M{
					"useragent": user_agent,
					"ip":        Ip,
					"logindate": time.Now(),
				},
			},
		}

		_cursorFind, _ := _HistoryCollection.Aggregate(context.TODO(), _value)
		if err1 := _cursorFind.All(context.TODO(), &value); err1 != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  false,
				"message": "default_error_message",
			})
		} else {
			if value == nil {
				_, errInsert := _HistoryCollection.InsertOne(context.TODO(), _valueinsert)
				if errInsert != nil {
					return c.JSON(http.StatusBadRequest, echo.Map{
						"status":  false,
						"message": "default_error_message",
					})
				}
			} else {
				_, _errhis := _HistoryCollection.UpdateOne(context.TODO(), bson.M{
					"user_id": _userid,
				}, bson.M{
					"$addToSet": bson.M{
						"login_histories": bson.M{
							"useragent": user_agent,
							"ip":        Ip,
							"logindate": time.Now(),
						},
					},
				})

				if _errhis != nil {
					fmt.Println("Error : update", _errhis)
					return c.JSON(http.StatusBadRequest, echo.Map{
						"status":  false,
						"message": "default_error_message",
					})
				}
			}
		}

		/////////////////////////////////////////////////////////////////

		_, _err := _userCollection.UpdateOne(context.TODO(), bson.M{
			"username": c.FormValue("username"),
		}, bson.M{
			"$set": bson.M{
				"push_token":                   c.FormValue("push_token"),
				"status.last_login.date":       time.Now(),
				"status.last_login.ip_address": Ip,
				"status.last_login.user_agent": user_agent,
			},
		})

		_claims := &model.JwtClaim{
			user.UID,
			role.Name.En,
			user.Cart,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			},
		}

		_token := jwt.NewWithClaims(jwt.SigningMethodHS256, _claims)

		_t, _err := _token.SignedString([]byte("RENNYYY//"))

		if _err != nil {
			fmt.Println("Error : JWT<LOG>")
			return _err
		}

		if user.Is_verify_email == false {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  false,
				"message": "Email not verify",
				"error":   inputUser.Is_verify_email,
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"status":       true,
			"access_token": _t,
			"userId":       user.UID,
			"login_histories": bson.M{
				"useragent": user_agent,
				"ip":        Ip,
				"logindate": time.Now(),
			},
			// "login_histories": user.Status.LastLogin,
			"role":          role.Name.En,
			"active_acount": user.Active,
			"cart":          user.Cart,
		})
	}

	return c.JSON(http.StatusBadRequest, echo.Map{
		"status":  false,
		"message": "wrongpass",
	})
}

func addToCart(c echo.Context) error {
	var (
		// value      []bson.M
		ctx        = context.Background()
		_UsersColl = libs.MongoDB.Database(keys.Database).Collection("Users")
		payload    = Cart{}
	)

	_userPayload := c.Get("user").(*jwt.Token)
	_claims := _userPayload.Claims.(*auth.JwtClaim)
	_uid, _ := primitive.ObjectIDFromHex(_claims.UID)

	if err := c.Bind(&payload); err != nil {
		return c.JSON(400, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	filter := bson.M{
		"_id": _uid,
	}

	update := bson.M{
		"$set": payload,
	}
	_ = _UsersColl.FindOneAndUpdate(ctx, filter, update)

	return c.JSON(200, echo.Map{
		"status":  true,
		"message": "Add To Cart.",
		"payload": payload,
	})

}

func getuserM(c echo.Context) error {
	var (
		value      []bson.M
		ctx        = context.Background()
		_UsersColl = libs.MongoDB.Database(keys.Database).Collection("Users")
	)

	pipeline := []bson.M{
		// bson.M{
		// 	"$match": bson.M{
		// 		"role": 1,
		// 	},
		// },
	}

	cursor, err := _UsersColl.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(400, echo.Map{
			"status":  false,
			"message": "cannot get booking list",
		})
	}

	if err := cursor.All(ctx, &value); err != nil {
		fmt.Println(err.Error())
		return c.JSON(400, echo.Map{
			"status":  false,
			"message": "cannot get booking list",
		})
	}

	// _ = _UsersColl.FindOne(context.TODO(), bson.M{"name": "admin"}).Decode(&value)
	fmt.Println("Test")
	fmt.Println("Test==>", value)

	return c.JSON(http.StatusOK, echo.Map{
		"status": true,
		"len":    len(value),
		"result": value,
	})
}

func createC(c echo.Context) error {
	_IDS := primitive.NewObjectID()
	userId, _ := primitive.ObjectIDFromHex("60e6364e1d20443d1c8b5274")
	payload := bson.M{
		"_id":             _IDS,
		"name":            "Renny-Admin",
		"exp-min":         0,
		"exp-max":         100,
		"level":           1,
		"STR":             1,
		"INT":             1,
		"create_datetime": time.Now(),
		"update_datetime": time.Now(),
		"Coin":            500,
		"userId":          userId,
	}

	fmt.Println("DB ==> ", keys.Database)

	_, errInsert := libs.MongoDB.Database(keys.Database).Collection("Char").InsertOne(context.TODO(), payload)
	if errInsert != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  false,
			"message": "Error : " + errInsert.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status": true,
		"result": payload,
	})
}

func showChar(c echo.Context) error {
	var (
		payload []bson.M
		ColChar = libs.MongoDB.Database(keys.Database).Collection("Char")
		ctx     = context.Background()
	)

	_userPayload := c.Get("user").(*jwt.Token)
	_claims := _userPayload.Claims.(*auth.JwtClaim)
	_uid := _claims.UID

	userId_token, _ := primitive.ObjectIDFromHex(_uid)

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"userId": userId_token,
			},
		},
	}

	cursor, err := ColChar.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(400, echo.Map{
			"status":  false,
			"message": "cannot get booking list",
		})
	}

	if err := cursor.All(ctx, &payload); err != nil {
		fmt.Println(err.Error())
		return c.JSON(400, echo.Map{
			"status":  false,
			"message": "cannot get booking list",
		})
	}

	// _ = ColChar.FindOne(ctx, bson.M{"userId": userId_token}).Decode(&payload)
	fmt.Println(len(payload))

	return c.JSON(http.StatusOK, echo.Map{
		"status":     true,
		"result":     payload,
		"len_result": len(payload),
	})
}

func datailChar(c echo.Context) error {
	var (
		payload []bson.M
		ColChar = libs.MongoDB.Database(keys.Database).Collection("Char")
		ctx     = context.Background()
	)
	det := c.Param("id")
	fmt.Println(det)
	// _userPayload := c.Get("user").(*jwt.Token)
	// _claims := _userPayload.Claims.(*auth.JwtClaim)
	// _uid := _claims.UID

	userId_token, _ := primitive.ObjectIDFromHex(det)

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"_id": userId_token,
			},
		},
	}

	cursor, err := ColChar.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(400, echo.Map{
			"status":  false,
			"message": "cannot get booking list",
		})
	}

	if err := cursor.All(ctx, &payload); err != nil {
		fmt.Println(err.Error())
		return c.JSON(400, echo.Map{
			"status":  false,
			"message": "cannot get booking list",
		})
	}

	// _ = ColChar.FindOne(ctx, bson.M{"userId": userId_token}).Decode(&payload)
	fmt.Println(len(payload))

	return c.JSON(http.StatusOK, echo.Map{
		"status":     true,
		"result":     payload,
		"len_result": len(payload),
	})
}

func CronD() {

	// var payload []bson.M
	// userId_token, _ := primitive.ObjectIDFromHex("60e6364e1d20443d1c8b5274")

	// _ = libs.MongoDB.Database(keys.Database).Collection("Char").FindOne(context.Background(), bson.M{"userId": userId_token}).Decode(&payload)

	// fmt.Println(payload[0]["level"])

	// if payload[0]["exp-min"] >= payload[0]["exp-max"] {
	// 	exp_max := payload[0]["exp-max"]
	// 	new_exp_max := exp_max * 150 / 100
	// 	level := payload[0]["level"]
	// 	levelup := level + 1
	// 	upINT := 0.4
	// 	INT := payload[0]["INT"]
	// 	new_INT := INT + upINT
	// 	exp_min := payload[0]["exp-min"]
	// 	new_exp_min := exp_min - exp_max

	// 	_, _ = libs.MongoDB.Database(keys.Database).Collection("Char").UpdateOne(context.TODO(), bson.M{
	// 		"userId": _ususerId_tokenerid,
	// 	}, bson.M{
	// 		"$Set": bson.M{
	// 			"level":   levelup,
	// 			"INT":     new_INT,
	// 			"exp-max": new_exp_max,
	// 			"exp-min": new_exp_min,
	// 		},
	// 	},
	// 	)

	// 	// return

	// } else {
	// 	_INT := payload[0]["INT"]
	// 	_exp_min := payload[0]["exp-min"]
	// 	update_exp := _exp_min + _INT
	// 	_, _ = libs.MongoDB.Database(keys.Database).Collection("Char").UpdateOne(context.TODO(), bson.M{
	// 		"userId": _ususerId_tokenerid,
	// 	}, bson.M{
	// 		"$Set": bson.M{
	// 			"exp-min": update_exp,
	// 		},
	// 	},
	// 	)
	// }
	// return
	fmt.Println(time.Now())
	// return

}

func updateexp(c echo.Context) error {
	var (
		value   []bson.M
		ctx     = context.Background()
		payload bson.M

		_UsersColl = libs.MongoDB.Database(keys.Database).Collection("Users")
	)

	userId_token, _ := primitive.ObjectIDFromHex("63e8fa0c0bbbf53f0dbaacf7")

	_ = libs.MongoDB.Database(keys.Database).Collection("Char").FindOne(context.Background(), bson.M{"userId": userId_token}).Decode(&payload)

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"userId": userId_token,
			},
		},
	}

	cursor, err := _UsersColl.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(400, echo.Map{
			"status":  false,
			"message": "cannot get booking list",
		})
	}

	if err := cursor.All(ctx, &value); err != nil {
		fmt.Println(err.Error())
		return c.JSON(400, echo.Map{
			"status":  false,
			"message": "cannot get booking list",
		})
	}

	var dd int32
	var _ususerId_tokenerid primitive.ObjectID
	for _, v := range value {
		dd = v["level"].(int32)

	}

	_ususerId_tokenerid = payload["userId"].(primitive.ObjectID)

	// dd := value.level.(int32)
	// fmt.Println("level ==>", v)
	fmt.Println("level ==>", dd)

	fmt.Println("payload ==>", payload)
	fmt.Println("_ususerId_tokenerid ==>", _ususerId_tokenerid)
	fmt.Println("payload min (int32) ==>", payload["exp-min"])
	fmt.Println("payload max (int32) ==>", payload["exp-max"])

	if payload["exp-min"].(int32) >= payload["exp-max"].(int32) {
		exp_max := payload["exp-max"].(int32)
		new_exp_max := exp_max * 150 / 100
		level := payload["level"].(int32)
		levelup := level + 1
		upINT := 0.4
		INT := payload["INT"].(float64)
		new_INT := INT + upINT
		exp_min := payload["exp-min"].(int32)
		new_exp_min := exp_min - exp_max
		fmt.Println("case1  ==>")

		fmt.Println("new_exp_max  ==>", new_exp_max)
		fmt.Println("new_exp_max  ==>")
		fmt.Println("new_exp_max  ==>", new_exp_max)

		_, _ = libs.MongoDB.Database(keys.Database).Collection("Char").UpdateOne(context.TODO(), bson.M{
			"userId": _ususerId_tokenerid,
		}, bson.M{
			"$Set": bson.M{
				"level":   levelup,
				"INT":     new_INT,
				"exp-max": new_exp_max,
				"exp-min": new_exp_min,
			},
		},
		)

		return c.JSON(http.StatusOK, echo.Map{
			"status":  true,
			"message": "update level",
		})

	} else {
		fmt.Println("case2  ==>")

		_INT := payload["INT"].(int32)
		_exp_min := payload["exp-min"].(int32)
		update_exp := _exp_min + _INT

		fmt.Println("_INT  ==>", _INT)
		fmt.Println("_exp_min  ==>", _exp_min)
		fmt.Println("update_exp  ==>", update_exp)

		_, errUpdate := libs.MongoDB.Database(keys.Database).Collection("Char").UpdateOne(context.TODO(), bson.M{
			"userId": _ususerId_tokenerid,
		}, bson.M{
			"$set": bson.M{
				"exp-min": 3,
			},
		},
		)

		if errUpdate != nil {
			return c.JSON(500, echo.Map{
				"status":  false,
				"message": "Error : " + errUpdate.Error(),
			})
		}

	}

	fmt.Println("case3  ==>")

	return c.JSON(http.StatusOK, echo.Map{
		"status":  true,
		"message": "update exp",
	})

}

// _, _errhis := _HistoryCollection.UpdateOne(context.TODO(), bson.M{
// 	"user_id": _userid,
// }, bson.M{
// 	"$addToSet": bson.M{
// 		"login_histories": bson.M{
// 			"useragent": user_agent,
// 			"ip":        Ip,
// 			"logindate": time.Now(),
// 		},
// 	},
// })
