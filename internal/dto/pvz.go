package dto

import "time"

type PVZReq struct {
	RegistrationDate time.Time `json:"registrationDate" format:"date-time"`
	City             string    `json:"city" binding:"required,oneof=Москва Санкт-Петербург Казань"`
}
