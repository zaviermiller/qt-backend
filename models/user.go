package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	u "qt-api/utils"
	"strings"
	"os"
	"golang.org/x/crypto/bcrypt"
)

// JWT Token struct
type Token struct {
	UserId uint
	jwt.StandardClaims
}

type User struct {
	gorm.Model
	Email string `json:"email"`
	FirstName string `json:"first_name`
	LastName string `json:"last_name"`
	Password string `json:"password"`
	Token string `json:"token";sql:"-"`
}

func (user *User) Validate() (map[string] interface{}, bool) {

	if !strings.Contains(user.Email, "@") {
		return u.Message(false, "Email address required"), false
	}

	if len(user.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	temp := &User{}

	// Check errors and email uniqness
	err := GetDB().Table("users").Where("email = ?", user.Email).First(temp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email already in use"), false
	}

	return u.Message(false, "Success"), true
}

func (user *User) Create() (map[string] interface{}) {

	if resp, ok := user.Validate(); !ok {
		return resp
	}

	// Generate pwd hash
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	GetDB().Create(user)

	if user.ID <= 0 {
		return u.Message(false, "Failed to create account")
	}

	// Create JWT for user
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	user.Password = ""

	response := u.Message(true, "User created")
	response["token"] = user.Token
	return response
}

func Login(email, password string) (map[string]interface{}) {

	user := &User{}
	err := GetDB().Table("users").Where("email = ?", email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "User not found")
		}
		return u.Message(false, "Connection error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid credentials")
	}

	user.Password = ""

	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	response := u.Message(true, "Success")
	response["token"] = user.Token
	return response
}

func GetUser(userId uint) *User {
	user := &User{}
	
	GetDB().Table("user").Where("id = ?", userId).First(user)
	if user.Email == "" {
		return nil
	}

	user.Password = ""

	return user
}