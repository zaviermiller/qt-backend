package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"strings"
	"os"
	"golang.org/x/crypto/bcrypt"
	"time"
	"math/rand"
	"errors"
)

var colors = []string{"red", "pink", "purple", "blue", "teal", "cyan", "green", "success", "amber", "deep-orange", "primary"}
// JWT Token struct
type Token struct {
	UserId uint
	jwt.StandardClaims
}

type User struct {
	QTModel
	Email string `json:"email"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Color string `json:"color"`
	Password string `json:"password"`
	Token string `json:"token" sql:"-"`
}

func (user *User) Validate() (string, bool) {

	if !strings.Contains(user.Email, "@") {
		return "Email address required", false
	}

	if len(user.Password) < 6 {
		return "Password is required", false
	}

	temp := &User{}

	// Check errors and email uniqness
	err := GetDB().Table("users").Where("email = ?", user.Email).First(temp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return "Connection error", false
	}
	if temp.Email != "" {
		return "Email already in use", false
	}

	return "Success", true
}

func (user *User) Create() (*User, error) {

	if resp, ok := user.Validate(); !ok {
		err := errors.New(resp)
		return nil, err
	}

	// Generate pwd hash
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	user.Color = colors[rand.Intn(len(colors) - 1)]

	GetDB().Create(user)

	if user.ID <= 0 {
		err := errors.New("Failed to create account")
		return nil, err
	}

	// Create JWT for user
	expirationTime := time.Now().Add(5 * time.Minute)
	tk := &Token{UserId: user.ID, StandardClaims: jwt.StandardClaims { ExpiresAt: expirationTime.Unix() }}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	user.Password = ""

	return user, nil
}

func Login(email, password string) (*string, error) {

	user := &User{}
	err := GetDB().Table("users").Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil { //Password does not match!
		return nil, err
	}

	user.Password = ""

	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	return &tokenString, nil
}

func GetUser(userId uint) *User {
	user := &User{}
	
	GetDB().Table("users").Where("id = ?", userId).First(user)
	if user.Email == "" {
		return nil
	}

	user.Password = ""

	return user
}