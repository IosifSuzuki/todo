package utility

import (
	"errors"
	"github.com/IosifSuzuki/todo/internall/model"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func CheckHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateAccessToken(userId int, userName string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := model.AccessClaims{
		UserId:   userId,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(Config.SecretKey))
	return tokenString, err
}

func GenerateRefreshToken(userId int) (string, error) {
	expirationTime := time.Now().Add(3 * 24 * time.Hour)
	claims := model.RefreshClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(Config.SecretKey))
	return tokenString, err
}

func VerifyToken(tokenString string) (bool, error) {
	claims := model.AccessClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(Config.SecretKey), nil
	})
	if err != nil {
		return false, err
	}
	return token.Valid, nil
}

func GetUserIdByFromToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(Config.SecretKey), nil
	})
	if err != nil {
		return 0, err
	}
	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("dont have access to payload of token")
	}
	if userId, ok := claimsMap[UserIdKey].(float64); ok {
		return int(userId), nil
	}
	return 0, errors.New("token not contains `user-id` field")
}

func VerifyIsAccessToken(tokenString string) (bool, error) {
	claims := model.AccessClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(Config.SecretKey), nil
	})
	if err != nil {
		return false, err
	}
	return len(claims.UserName) != 0, nil
}
