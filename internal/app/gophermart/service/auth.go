package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/volnistii11/accumulative-loyalty-system/internal/cerrors"
	"github.com/volnistii11/accumulative-loyalty-system/internal/model"
	"time"
)

const (
	salt       = "dasd6as76das76das7d6as76d76"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type UserAuthorize interface {
	RegisterUser(user *model.User) error
	GetUser(user *model.User) *model.User
}

type Auth struct {
	db UserAuthorize
}

func NewAuth(db UserAuthorize) *Auth {
	return &Auth{
		db: db,
	}
}

func (a *Auth) RegisterUser(user *model.User) error {
	user.Password = generatePasswordHash(user.Password)
	err := a.db.RegisterUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (a *Auth) AuthenticateUser(user *model.User) (string, error) {
	user.Password = generatePasswordHash(user.Password)

	user = a.db.GetUser(user)
	if user.ID == 0 {
		return "", cerrors.ErrHTTPStatusNoContent
	}

	jwtToken, err := BuildJWTString(user.ID)
	if err != nil {
		return "", errors.Wrap(err, "service.auth.AuthenticateUser.BuildJWT")
	}

	return jwtToken, nil
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func BuildJWTString(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserID(tokenString string) (int, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(signingKey), nil
		})
	if err != nil {
		return 0, cerrors.ErrHTTPStatusUnauthorized
	}
	if !token.Valid {
		return 0, cerrors.ErrHTTPStatusUnauthorized
	}
	return claims.UserID, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
