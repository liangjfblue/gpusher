package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Context struct {
	UUID string
}

var (
	_secret  = "suVyr9228ohAhg7823A"
	_expired = 3600 * 8 //默认8小时, gateway会根据心跳时间来续期
)

func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secret), nil
	}
}

func Parse(tokenString string, secret string) (*Context, error) {
	ctx := &Context{}

	token, err := jwt.Parse(tokenString, secretFunc(secret))
	if err != nil {
		return ctx, err
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx.UUID = claims["uuid"].(string)
		return ctx, nil
	} else {
		return nil, err
	}
}

func ParseRequest(token string) (*Context, error) {
	return Parse(token, _secret)
}

func SignToken(c Context) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid": c.UUID,
		"nbf":  time.Now().Unix(),
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Second * time.Duration(_expired)).Unix(),
	})

	tokenString, err = token.SignedString([]byte(_secret))
	return
}
