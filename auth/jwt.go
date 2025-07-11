package auth

import (
	"awesomeProject/models"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
	"time"
)

var jwtKey = []byte("supersecretkey")

func GenerateJWT(user models.User) (tokenString string, err error) {
	expirationTime := time.Now().Add(24 * 90 * time.Hour)
	claims := &models.JWTClaim{
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Phone:        user.Phone,
		Password:     user.Password,
		Role:         user.Role,
		ProfilePhoto: user.ProfilePhoto,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (context.Context, error) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&models.JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.JWTClaim)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("your authentication token expired. Login again to continue")
	}
	ctx := context.WithValue(context.Background(), "userClaims", claims)
	return ctx, nil

}

func GetUserDetailsFromToken(token string) (*models.JWTClaim, error) {
	ctx, err := ValidateToken(token)
	if err != nil {
		return nil, err
	}

	claims, ok := ctx.Value("userClaims").(*models.JWTClaim)
	if !ok {
		return nil, errors.New("couldn't retrieve user claims from context")
	}

	return claims, nil
}
