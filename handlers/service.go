package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"kong.com/catalog/datastore"
	"kong.com/catalog/models"
)

var (
	defaultPage     = 1
	defaultPageSize = 10
)

type ServiceHandler struct {
	serviceModel *models.Service
}

func (h *ServiceHandler) GetServiceById(ctx *gin.Context) {
	serviceId := ctx.Param("id")
	service, err := h.serviceModel.FetchServiceById(ctx, serviceId)

	if err != nil {
		slog.Error("failed to find service with id: %s: %w", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("failed to find service with id: %s", serviceId)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": service.Id, "name": service.Name, "description": service.Description, "version": service.Version})
}

func (h *ServiceHandler) Search(ctx *gin.Context) {
	searchRequest := datastore.ServiceSearchRequest{}

	searchRequest.Name = strings.TrimSpace(ctx.Query("name"))
	searchRequest.Version = strings.TrimSpace(ctx.Query("version"))

	searchRequest.Page = parseIntOrDefault(strings.TrimSpace(ctx.Query("page")), defaultPage)
	searchRequest.PageSize = parseIntOrDefault(strings.TrimSpace(ctx.Query("page_size")), defaultPageSize)

	searchResponse, err := h.serviceModel.Search(ctx, searchRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to find services with given query: %s", err.Error())})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"items": searchResponse.Items, "page": searchResponse.Page, "page_size": searchResponse.PageSize})
}

func NewServiceHandler(model *models.Service) *ServiceHandler {
	return &ServiceHandler{serviceModel: model}
}

func parseIntOrDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
