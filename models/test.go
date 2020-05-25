package models

import (
	"github.com/jinzhu/gorm"
	"fmt"
	"errors"
)

type Test struct {
	QTModel
	Name string `json:"name"`
	UserID uint `json:"userId"`
	User User `json:"-" sql:"-"`
	UrlSlug string `json:"urlSlug"`
	Sections []Section `json:"sections" sql:"-"`
}

func (test *Test) Valid() (string, bool) {
	var errorsString = ""
	if len(test.Name) < 1 {
		errorsString += "Name is required\n"
	}
	user := &User{}

	err := GetDB().Table("users").Where("id = ?", test.UserID).First(user).Error

	if err != nil && err == gorm.ErrRecordNotFound {
		errorsString += fmt.Sprintf("User with id %d does not exist\n", user.ID)
	}

	if errorsString != "" {
		return errorsString, false
	} else {
		return "Success", true
	}
}

func (test *Test) Create() (*Test, error) {
	if resp, ok := test.Valid(); !ok {
		return nil, errors.New(resp)
	}

	GetDB().Create(test)

	if test.ID <= 0 {
		err := errors.New("Test could not be created")
		return nil, err
	}

	return test, nil

}

func FindTestsFromUser(userId uint) ([]Test, error) {
	tests := []Test{}

	err := GetDB().Table("tests").Where("user_id = ?", userId).Find(&tests).Error

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return tests, nil
}

func (test *Test) FindByID(id int) (*Test, error) {
	err := GetDB().Table("tests").Where("id = ?", id).First(test).Error

	test.Sections = []Section{}

	if err != nil {
		return nil, err
	}

	return test, nil

}