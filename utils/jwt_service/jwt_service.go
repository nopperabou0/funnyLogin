package jwtservice

import (
	"funny-login/config"
	"funny-login/model"
	modelutil "funny-login/utils/model_util"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var conf config.TokenConfig

func CreateToken(user model.User) (string, error) {
	signatureKey := conf.JwtSignatureKey
	claims := modelutil.JwtPayloadClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    conf.ApplicationName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(conf.AccessTokenLifeTime)),
		},
		UserId: user.Id,
		Role:   user.Role,
	}
	token, err := jwt.NewWithClaims(conf.JwtSigningMethod, claims).SignedString(signatureKey)
	if err != nil {
		return "", err
	}

	return token, nil
}
