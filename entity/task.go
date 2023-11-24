package entity

import "github.com/asaskevich/govalidator"

type Task struct {
	GormModel
	Title       string `gorm:"not null" json:"title" valid:"required~Title is required"`
	Description string `gorm:"not null" json:"description" valid:"required~Description is required"`
	Status      bool   `gorm:"not null" json:"status"`
	UserID      uint
	CategoryID  uint
	User        User     `valid:"-"`
	Category    Category `valid:"-"`
}

func (p *Task) Validate() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
