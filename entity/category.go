package entity

import "github.com/asaskevich/govalidator"

type Category struct {
	GormModel
	Type string `gorm:"not null" json:"type" valid:"required~Type is required"`
	Task []Task `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"task"`
}

func (p *Category) Validate() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
