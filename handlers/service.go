package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"kong.com/catalog/models"
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

func NewServiceHandler(model *models.Service) *ServiceHandler {
	return &ServiceHandler{serviceModel: model}
}
