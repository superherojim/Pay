package jwt

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type JWT struct {
	key []byte
	rdb *redis.Client
}

type MyCustomClaims struct {
	UserId   string
	UserType string
	ClientIP string
	jwt.RegisteredClaims
}

func NewJwt(conf *viper.Viper, rdb *redis.Client) *JWT {
	return &JWT{
		key: []byte(conf.GetString("security.jwt.key")),
		rdb: rdb,
	}
}

func (j *JWT) GenToken(userId, userType, clientIP string, expiresAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		UserId:   userId,
		UserType: userType,
		ClientIP: clientIP,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})

	token.Header["cty"] = "cheemshappy_pay-auth"
	return token.SignedString(j.key)
}

func (j *JWT) ParseToken(tokenString string) (*MyCustomClaims, error) {
	re := regexp.MustCompile(`(?i)Bearer `)
	tokenString = re.ReplaceAllString(tokenString, "")
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	// 先检查黑名单
	if j.IsInBlacklist(context.Background(), tokenString) {
		return nil, errors.New("token revoked")
	}

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.key, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

// 加入黑名单
func (j *JWT) AddToBlacklist(ctx context.Context, tokenString string, expiresAt *jwt.NumericDate) error {
	if expiresAt == nil {
		return errors.New("invalid token expiration")
	}

	ttl := time.Until(expiresAt.Time)
	if ttl <= 0 {
		return nil // 已过期的token不需要加入
	}

	return j.rdb.SetNX(ctx, "jwt:blacklist:"+tokenString, "1", ttl).Err()
}

// 检查黑名单
func (j *JWT) IsInBlacklist(ctx context.Context, tokenString string) bool {
	exists, _ := j.rdb.Exists(ctx, "jwt:blacklist:"+tokenString).Result()
	return exists > 0
}
