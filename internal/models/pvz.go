package model

import "time"

type PVZ struct {
	ID               string    `json:"id,omitempty" format:"uuid"`
	RegistrationDate time.Time `json:"registrationDate" format:"date-time"`
	City             string    `json:"city" binding:"required,oneof=Москва Санкт-Петербург Казань"`
}
