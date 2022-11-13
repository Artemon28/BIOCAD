package handlers

import (
	"BIOCAD/internal/services"
	"BIOCAD/internal/structures"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler struct {
	service *services.Service
}

func NewHandler(sr *services.Service) *Handler {
	return &Handler{service: sr}
}

func (h *Handler) GetDevices(c *gin.Context) {
	guid, pagination := GeneratePaginationFromRequest(c)
	deviceLists, err := h.service.GetAll(guid, &pagination)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"data": deviceLists,
	})
}

func GeneratePaginationFromRequest(c *gin.Context) (string, structures.Pagination) {
	limit := 2
	page := 1
	guid := ""
	query := c.Request.URL.Query()
	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
			break
		case "page":
			page, _ = strconv.Atoi(queryValue)
			break
		case "unit_guid":
			guid = queryValue
			break
		}
	}
	return guid, structures.Pagination{
		Limit: limit,
		Page:  page,
	}

}
