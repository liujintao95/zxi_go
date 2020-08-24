package utils

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	_ "os"
	"time"
)

func ErrCheck(g *gin.Context, err error, msg string, httpCode int) {
	if err != nil && err != sql.ErrNoRows {
		Logging.Error(msg + ":" + err.Error())
		if httpCode != 0 {
			g.JSON(httpCode, gin.H{
				"errmsg": msg,
				"data":   nil,
			})
			panic(msg + ":" + err.Error())
		}
	}
}

func CreatToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(""))
	return tokenString, err
}

func FileSha1(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	iSha1 := sha1.New()
	_, err = io.Copy(iSha1, src)
	return hex.EncodeToString(iSha1.Sum(nil)), err
}
