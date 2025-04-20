package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	service "github.com/kirillidk/pvz-service/internal/service/reception"
)

type ReceptionHandler struct {
	receptionService *service.ReceptionService
}

func NewReceptionHandler(receptionService *service.ReceptionService) *ReceptionHandler {
	return &ReceptionHandler{
		receptionService: receptionService,
	}
}

func (h *ReceptionHandler) CreateReception(c *gin.Context) {
	var receptionCreateReq dto.ReceptionCreateRequest
	if err := c.ShouldBindJSON(&receptionCreateReq); err != nil {
		c.JSON(http.StatusBadRequest, model.Error{Message: "Invalid request data"})
		return
	}

	reception, err := h.receptionService.CreateReception(c.Request.Context(), receptionCreateReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reception)
}
