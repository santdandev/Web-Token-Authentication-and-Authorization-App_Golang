package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtkey = []byte("dont_share_this")

type Creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Signup(w http.ResponseWriter, r *http.Request) {
	cred := &Creds{}
	err := json.NewDecoder(r.Body).Decode(cred)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(cred.Password), 10)
	stm, err := db.Prepare("INSERT INTO userinfo(username, password) values (?,?)")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = stm.Exec(cred.Username, string(hash))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func Signin(w http.ResponseWriter, r *http.Request) {
	cred := &Creds{}
	err := json.NewDecoder(r.Body).Decode(cred)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(cred)
	result := db.QueryRow("select password from userinfo where username=$1", cred.Username)
	var temp_cred string
	err = result.Scan(&temp_cred)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(cred.Username, cred.Password)
	if err = bcrypt.CompareHashAndPassword([]byte(temp_cred), []byte(cred.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}
	ttl := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		Username:       cred.Username,
		StandardClaims: jwt.StandardClaims{ExpiresAt: ttl.Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookie := &http.Cookie{Name: "jwt", Value: tokenString, Expires: ttl}
	http.SetCookie(w, cookie)
}

func CheckToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	claims, ok := checkToken(cookie.Value)
	if !ok {
		log.Println("wrong claims" + claims.Username)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func checkToken(tokenStr string) (*Claims, bool) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})
	if !token.Valid || err != nil {
		return claims, false
	}
	return claims, true
}
