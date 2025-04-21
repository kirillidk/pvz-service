package dto

import (
	"time"

	"github.com/kirillidk/pvz-service/internal/model"
)

type PVZCreateRequest struct {
	RegistrationDate time.Time `json:"registrationDate" format:"date-time"`
	City             string    `json:"city" binding:"required,oneof=Москва Санкт-Петербург Казань"`
}

type PVZFilterQuery struct {
	StartDate *time.Time `form:"startDate"`
	EndDate   *time.Time `form:"endDate"`
	Page      int32      `form:"page,default=1" binding:"min=1"`
	Limit     int32      `form:"limit,default=10" binding:"min=1,max=30"`
}

type PVZWithReceptionsResponse struct {
	PVZ        model.PVZ                       `json:"pvz"`
	Receptions []ReceptionWithProductsResponse `json:"receptions"`
}

type ReceptionWithProductsResponse struct {
	Reception model.Reception `json:"reception"`
	Products  []model.Product `json:"products"`
}

type PaginatedResponse struct {
	Data []PVZWithReceptionsResponse `json:"data"`
}

type Pagination struct {
	Total      int32 `json:"total"`
	Page       int32 `json:"page"`
	Limit      int32 `json:"limit"`
	TotalPages int32 `json:"totalPages"`
}
