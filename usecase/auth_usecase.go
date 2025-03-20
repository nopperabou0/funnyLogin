package usecase

import jwtservice "funny-login/utils/jwt_service"

func Login(name string, password string) (string, error) {
	user, err := GetUserByNamePassword(name, password)
	var token string
	if err != nil {
		return "", err
	}

	token, err = jwtservice.CreateToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
}
