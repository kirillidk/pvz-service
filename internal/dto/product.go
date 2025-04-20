package dto

type ProductCreateRequest struct {
	Type  string `json:"type" binding:"required,oneof=электроника одежда обувь"`
	PVZID string `json:"pvzId" binding:"required,uuid"`
}
