package app

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 定义一个密钥，实际开发中建议放在 config.yaml 里
var jwtKey = []byte("your_secret_key")

type Claims struct {
	UserID uint
	jwt.RegisteredClaims
}

// GenerateToken 生成令牌
func GenerateToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 有效期 24 小时
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ParseToken 解析令牌
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}
