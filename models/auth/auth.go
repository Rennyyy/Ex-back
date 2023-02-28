package auth

import (
	"github.com/golang-jwt/jwt/v4"
)

type JwtClaim struct {
	UID  string                   `json:"_id"`
	Name string                   `json:"name"`
	Role string                   `json:"role"`
	Cart []map[string]interface{} `json:"cart" bson:"cart`
	jwt.StandardClaims
}
