package jwtservice

import (
	"fmt"
	"funny-login/config"
	"funny-login/model"
	modelutil "funny-login/utils/model_util"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var Conf config.TokenConfig

func CreateToken(user model.User) (string, error) {
	signatureKey := Conf.JwtSignatureKey
	claims := modelutil.JwtPayloadClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Conf.ApplicationName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(Conf.AccessTokenLifeTime)),
		},
		UserId: user.Id,
		Role:   user.Role,
	}
	token, err := jwt.NewWithClaims(Conf.JwtSigningMethod, claims).SignedString(signatureKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyToken(token string) (modelutil.JwtPayloadClaim, error) {
	tokenParse, err := jwt.ParseWithClaims(
		token,
		&modelutil.JwtPayloadClaim{},
		func(t *jwt.Token) (any, error) {
			return Conf.JwtSignatureKey, nil
		})
	if err != nil {
		return modelutil.JwtPayloadClaim{}, fmt.Errorf("error token parsing : %s", err.Error())
	}
	claim, ok := tokenParse.Claims.(*modelutil.JwtPayloadClaim)
	if !ok {
		return modelutil.JwtPayloadClaim{}, fmt.Errorf("error claim")
	}
	return *claim, nil
}
