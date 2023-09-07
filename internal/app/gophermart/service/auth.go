package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/volnistii11/accumulative-loyalty-system/internal/cerrors"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/gerr"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/sl"
	"github.com/volnistii11/accumulative-loyalty-system/internal/model"
	"golang.org/x/exp/slog"
	"net/http"
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
	db     UserAuthorize
	logger *slog.Logger
}

func NewAuth(db UserAuthorize, logger *slog.Logger) *Auth {
	return &Auth{
		db:     db,
		logger: logger,
	}
}

func (a *Auth) RegisterUser(w http.ResponseWriter, user *model.User) (http.ResponseWriter, error) {
	if user.Login == "" || user.Password == "" {
		a.logger.Error("wrong request format")
		w.WriteHeader(http.StatusBadRequest)
		return w, cerrors.ErrHTTPWrongRequestFormat
	}

	user.Password = generatePasswordHash(user.Password)
	err := a.db.RegisterUser(user)
	if err != nil {
		a.logger.Error("failed user register", sl.Err(err))
		if gerr.IsDuplicateKey(err) {
			w.WriteHeader(http.StatusConflict)
			return w, cerrors.ErrHTTPUserExists
		}
		w.WriteHeader(http.StatusInternalServerError)
		return w, cerrors.ErrInternalServer
	}
	a.logger.Info("user registered")
	return w, nil
}

func (a *Auth) AuthenticateUser(w http.ResponseWriter, user *model.User) (http.ResponseWriter, error) {
	user.Password = generatePasswordHash(user.Password)

	user = a.db.GetUser(user)
	if user.ID == 0 {
		a.logger.Error("failed user authentication", sl.Err(cerrors.ErrHTTPStatusNoContent))
		w.WriteHeader(http.StatusInternalServerError)
		return w, cerrors.ErrHTTPStatusNoContent
	}

	jwtToken, err := BuildJWTString(user.ID)
	if err != nil {
		a.logger.Error("failed user authentication", sl.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return w, errors.Wrap(err, "service.auth.AuthenticateUser.BuildJWT")
	}
	a.logger.Info("user authenticated")

	cookie := http.Cookie{Name: "jwtToken", Value: jwtToken}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)

	return w, nil
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
