package model

import (
	"github.com/golang-jwt/jwt"
)

type JwtClaim struct {
	UID  string                   `json:"_id" bson:"_id"`
	Role string                   `json:"role"`
	Cart []map[string]interface{} `json:"cart" bson:"cart`
	jwt.StandardClaims
}
