package model

import "time"

type Reception struct {
	ID       string    `json:"id,omitempty" format:"uuid"`
	DateTime time.Time `json:"dateTime" binding:"required" format:"date-time"`
	PVZID    string    `json:"pvzId" binding:"required,uuid"`
	Status   string    `json:"status" binding:"required,oneof=in_progress close"`
}
