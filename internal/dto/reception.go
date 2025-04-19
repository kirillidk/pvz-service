package dto

type ReceptionCreateRequest struct {
	PVZID string `json:"pvzId" binding:"required,uuid"`
}
