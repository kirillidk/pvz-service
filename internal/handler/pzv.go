package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	service "github.com/kirillidk/pvz-service/internal/service/pvz"
)

type PVZHandler struct {
	pvzService *service.PVZService
}

func NewPVZHandler(pvzService *service.PVZService) *PVZHandler {
	return &PVZHandler{
		pvzService: pvzService,
	}
}

func (h *PVZHandler) CreatePVZ(c *gin.Context) {
	var pvzReq dto.PVZCreateRequest
	if err := c.ShouldBindJSON(&pvzReq); err != nil {
		c.JSON(http.StatusBadRequest, model.Error{Message: "Invalid request data"})
		return
	}

	if _, ok := model.ValidCities[pvzReq.City]; !ok {
		c.JSON(http.StatusBadRequest, model.Error{
			Message: "City must be one of: Москва, Санкт-Петербург, Казань",
		})
		return
	}

	createdPVZ, err := h.pvzService.CreatePVZ(c.Request.Context(), pvzReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdPVZ)
}

func (h *PVZHandler) GetPVZList(c *gin.Context) {
	var filter dto.PVZFilterQuery
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, model.Error{Message: "Invalid query parameters"})
		return
	}

	result, err := h.pvzService.GetPVZList(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
