package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	service "github.com/kirillidk/pvz-service/internal/service/product"
)

type ProductHandler struct {
	productService service.ProductServiceInterface
}

func NewProductHandler(productService service.ProductServiceInterface) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var productCreateReq dto.ProductCreateRequest
	if err := c.ShouldBindJSON(&productCreateReq); err != nil {
		c.JSON(http.StatusBadRequest, model.Error{Message: "Invalid request data"})
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), productCreateReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) DeleteLastProduct(c *gin.Context) {
	pvzID := c.Param("pvzId")
	if pvzID == "" {
		c.JSON(http.StatusBadRequest, model.Error{Message: "PVZ ID is required"})
		return
	}

	err := h.productService.DeleteLastProduct(c.Request.Context(), pvzID)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Last product deleted successfully"})
}
