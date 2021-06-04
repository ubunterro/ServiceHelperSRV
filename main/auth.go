package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

func basicAuth() gin.HandlerFunc {

	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			respondWithError(401, "Unauthorized", c)
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		authResult, userStatus, userId := authenticateUser(pair[0], pair[1])

		if len(pair) != 2 || !authResult {
			respondWithError(401, "Unauthorized", c)
			return
		}

		c.Set("userStatus", userStatus)
		c.Set("userId", userId)

		c.Next()
	}
}

func authenticateUser(username, password string) (result bool, status int, userId int) {
	db := DBConn()
	var DBpasswordHash string

	err := db.QueryRow("SELECT pass, status, user_id FROM users WHERE login = ?;", username).Scan(&DBpasswordHash, &status, &userId)
	log.Println("STATUS", status)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		println("HASH " + DBpasswordHash)
	}

	var hash Hash

	err = hash.Compare(DBpasswordHash, password)

	if err != nil {
		log.Println(err)
		return false, -1, userId
	} else {
		log.Println("success")
		return true, status, userId
	}

}

func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}

	c.JSON(code, resp)
	c.Abort()
}

//Hash implements root.Hash
type Hash struct{}

//Generate a salted hash for the input string
func (c *Hash) Generate(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

//Compare string to generated hash
func (c *Hash) Compare(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}
