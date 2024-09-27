package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	Participant Role = "participant"
	Maintainer  Role = "maintainer"
	Captain     Role = "captain"
)

type Order struct {
	ID          string `json:"id" gorm:"type:uuid;unique;not null;index;primary_key"`
	FIO         string `json:"fio" gorm:"type:string;not null"`
	Age         int    `json:"age" gorm:"type:int;not null"`
	Role        Role   `json:"role" gorm:"type:string;not null"`
	PhoneNumber string `json:"phone_number" gorm:"type:string;unique;not null"`
	Email       string `json:"email" gorm:"type:string;not null"`
	TeamID      string `json:"team_id" gorm:"type:uuid;not null;index"`
}

func (order *Order) BeforeCreate(tx *gorm.DB) error {
	order.ID = uuid.New().String()
	return nil
}
