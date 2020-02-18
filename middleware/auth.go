package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"muses.service/model"
)

const contextKeyUserObj = "authedUserObj"
const bearerLength = len("Bearer ")

func ctxTokenToUser(c *gin.Context) {
	token, ok := c.GetQuery("_t")
	if !ok {
		hToken := c.GetHeader("Authorization")
		if len(hToken) < bearerLength {
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"msg": "header Authorization has not Bearer token"})
			return
		}
		token = strings.TrimSpace(hToken[bearerLength : len(hToken)-1])
	}

	println(token)
	usr, err := jwtParseUser(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"msg": err.Error()})
		return
	}

	//store the user Model in the context
	c.Set(contextKeyUserObj, *usr)
	c.Next()
	// after request
}

func MwUser(c *gin.Context) {
	ctxTokenToUser(c)
}

var AppSecret = ""

type userStdClaims struct {
	jwt.StandardClaims
	*model.User
}

func (c userStdClaims) Valid() (err error) {
	if c.VerifyExpiresAt(time.Now().Unix(), true) == false {
		return errors.New("token is expired")
	}

	if c.User.ID < 0 {
		return errors.New("invalid user in jwt")
	}
	return
}

func JwtGenerateToken(m *model.User) (string, error) {
	m.Password = ""
	expireTime := time.Now().Add(time.Hour * 24)
	stdClaims := jwt.StandardClaims{
		ExpiresAt: expireTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        fmt.Sprintf("%d", m.ID),
	}

	uClaims := userStdClaims{
		StandardClaims: stdClaims,
		User:           m,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uClaims)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(AppSecret))
	if err != nil {
		logrus.WithError(err).Fatal("config is wrong, can not generate jwt")
	}
	return tokenString, err
}

func jwtParseUser(tokenString string) (*model.User, error) {
	if tokenString == "" {
		return nil, errors.New("no token is found in Authorization Bearer")
	}

	claims := userStdClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		println(token)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(AppSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims.User, err
}
