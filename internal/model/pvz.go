package model

import "time"

var ValidCities = map[string]struct{}{
	"Москва":          {},
	"Санкт-Петербург": {},
	"Казань":          {},
}

type PVZ struct {
	ID               string    `json:"id,omitempty" format:"uuid"`
	RegistrationDate time.Time `json:"registrationDate" format:"date-time"`
	City             string    `json:"city" binding:"required,oneof=Москва Санкт-Петербург Казань"`
}
