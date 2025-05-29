package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type CustomClaims struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	UserPhone string `json:"user_phone"`
	Token     string `json:"-"` // Optional field to include token in claims
	jwt.RegisteredClaims
}

const ContextClaimKey = "ctx.mw.auth.claim"

var (
	TokenExpiration = 24 * time.Hour
)

func GenerateToken(jwtSecret []byte, userID, userEmail, userName, userPhone string) (string, error) {
	claims := CustomClaims{
		UserID:    userID,
		UserName:  userName,
		UserEmail: userEmail,
		UserPhone: userPhone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	return token.SignedString(jwtSecret)
}

func ParseToken(jwtSecret []byte, tokenString string) (*CustomClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // check signing method
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	// Check for parsing errors
	if err != nil {
		return nil, err
	}

	// Validate the token and retrieve claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		claims.Token = tokenString
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func GetClaimsFromContext(c *gin.Context) *CustomClaims {
	v := c.Value(ContextClaimKey)
	token := new(CustomClaims)

	if v == nil {
		return token
	}

	out, ok := v.(*CustomClaims)
	if !ok {
		return token
	}

	zap.L().
		With(zap.String("user_id", out.UserID)).
		Sugar().
		Info("token parsed")

	return out
}
