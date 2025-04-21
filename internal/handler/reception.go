package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	service "github.com/kirillidk/pvz-service/internal/service/reception"
)

type ReceptionHandler struct {
	receptionService service.ReceptionServiceInterface
}

func NewReceptionHandler(receptionService service.ReceptionServiceInterface) *ReceptionHandler {
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

func (h *ReceptionHandler) CloseLastReception(c *gin.Context) {
	pvzID := c.Param("pvzId")
	if pvzID == "" {
		c.JSON(http.StatusBadRequest, model.Error{Message: "PVZ ID is required"})
		return
	}

	closedReception, err := h.receptionService.CloseLastReception(c.Request.Context(), pvzID)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, closedReception)
}
