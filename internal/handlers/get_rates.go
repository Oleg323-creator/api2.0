package handlers

import (
	"github.com/Oleg323-creator/api2.0/internal/db/rep"
	"github.com/gin-gonic/gin"
	"net/http"
)

// TO CONNECT IT WITH db.repository
type Handler struct {
	repository *rep.Repository
}

func NewHandler(repository *rep.Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) GetEndpoint(c *gin.Context) {
	var params rep.FilterParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 4
	}
	if params.OrderDir == "" {
		params.OrderDir = "asc"
	}

	data, err := h.repository.GetRatesFromDB(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
