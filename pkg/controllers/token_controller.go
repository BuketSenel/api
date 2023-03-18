package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func getSigningKey() []byte {
	return []byte("selfservice")
}

func CreateJWTToken(email string) (string, gin.H) {
	role, uid, resid, header := getUser(email)
	if header["status"] != 200 {
		return "", header
	}
	SigningKey := getSigningKey()
	user_id := strconv.FormatInt(int64(uid), 10)
	rest_id := strconv.FormatInt(int64(resid), 10)
	type selfClaims struct {
		restId string `json:"resId"`
		role   string `json:"role"`
		uid    string `json:"uid"`
		jwt.RegisteredClaims
	}

	jwt_claims := selfClaims{
		role,
		user_id,
		rest_id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt_claims)
	ss, err := token.SignedString(SigningKey)
	return ss, gin.H{"status": 200, "message": "Token Created", "token": ss, "error": err}
}

func TokenValidation(c *gin.Context) error {
	SigningKey := getSigningKey()
	tokenStr := TokenExtraction(c)
	_, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SigningKey), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func TokenExtraction(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func TokenIDExtraction(c *gin.Context) (uint, error) {

	SigningKey := getSigningKey()
	tokenString := TokenExtraction(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SigningKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(uid), nil
	}
	return 0, nil
}
