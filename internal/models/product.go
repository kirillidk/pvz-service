package model

import "time"

type Product struct {
	ID          string    `json:"id,omitempty" format:"uuid"`
	DateTime    time.Time `json:"dateTime" binding:"required" format:"date-time"`
	Type        string    `json:"type" binding:"required,oneof=электроника одежда обувь"`
	ReceptionID string    `json:"receptionId" binding:"required,uuid"`
}
