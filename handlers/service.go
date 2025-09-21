package handlers

import (
	"errors"
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
	// The char `-` before a value denotes sort in descending order else ascending order.
	ValidSortOptions = map[string]bool{"name": true, "-name": true, "version": true, "-version": true}
)

type ServiceHandler struct {
	serviceModel *models.Service
}

func (h *ServiceHandler) GetServiceById(ctx *gin.Context) {
	serviceId := ctx.Param("id")
	service, err := h.serviceModel.FetchServiceById(ctx, serviceId)

	if err != nil {
		slog.Error("failed to find service with id: %s: %w", err)
		if errors.Is(err, datastore.ResourceNotFoundError) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("failed to find service with id: %s", serviceId)})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal occur occurred. Please try again."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": service.Id, "name": service.Name, "description": service.Description, "version": service.Version})
}

func (h *ServiceHandler) Search(ctx *gin.Context) {
	searchRequest := datastore.ServiceSearchRequest{}

	searchRequest.Name = strings.TrimSpace(ctx.Query("name"))
	searchRequest.Version = strings.TrimSpace(ctx.Query("version"))
	sortOrder, err := parseSortParameter(strings.TrimSpace(ctx.Query("sort")))
	searchRequest.Sort = sortOrder
	searchRequest.Keyword = strings.TrimSpace(ctx.Query("keyword"))

	if err != nil {
		if errors.Is(err, SearchIncorrectQuerySyntaxError) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to parse sort parameter, check query again.")})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("internal error occurred. Please try again in some time.")})
		return
	}

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

func parseSortParameter(sortQuery string) ([]string, error) {
	if sortQuery == "" {
		return []string{}, nil
	}
	sortOptions := strings.Split(sortQuery, ",")
	selectedSortOpts := []string{}
	for _, opt := range sortOptions {
		if _, ok := ValidSortOptions[opt]; !ok {
			return selectedSortOpts, fmt.Errorf("invalid sort option in query %v: %w", opt, SearchIncorrectQuerySyntaxError)
		}
		selectedSortOpts = append(selectedSortOpts, opt)
	}

	return selectedSortOpts, nil
}
