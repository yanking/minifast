package middlewares

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type CustomClaims struct {
	UserID      uint64 `json:"userid"`
	Nickname    string
	AuthorityId uint64
	jwt.RegisteredClaims
}

type JWT struct {
	SigningKey []byte
}

func NewJWT(signKey string) *JWT {
	return &JWT{
		SigningKey: []byte(signKey), // 可以设置过期时间
	}
}

func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string, leeway *time.Duration) (*CustomClaims, error) {
	var opts []jwt.ParserOption
	if leeway != nil {
		opts = append(opts, jwt.WithLeeway(*leeway))
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	}, opts...)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func (j *JWT) RefreshToken(tokenString string) (string, error) {
	leeway := time.Minute * 5
	claims, err := j.ParseToken(tokenString, &leeway)
	if err != nil {
		return "", err
	}

	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
	claims.RegisteredClaims.IssuedAt = jwt.NewNumericDate(time.Now())

	return j.CreateToken(*claims)
}
