package entity

import (
	"errors"
	"github.com/alyzsa/FinPro3/helper"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	GormModel
	Full_Name string `gorm:"not null;type:varchar(191)" json:"full_name" valid:"required~Your full name is required"`
	Email     string `gorm:"not null;type:varchar(191)" json:"email" valid:"required~Your email is required,email~Invalid email format"`
	Password  string `gorm:"not null" json:"password" valid:"required~Your password is required,minstringlength(6)~Password has to have minimum length of 6 characters"`
	Role      string `gorm:"not null" json:"role"`
}

func (u *User) Validate() error {
	if u.Role != "admin" && u.Role != "member" {
		return errors.New("Role must be admin or member")
	}
	_, err := govalidator.ValidateStruct(u)
	return err
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}

	u.Password = helper.HashPass(u.Password)
	err = nil
	return
}
