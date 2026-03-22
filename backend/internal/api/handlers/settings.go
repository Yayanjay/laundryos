package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laundryos/backend/internal/api/middleware"
	"github.com/laundryos/backend/internal/service"
	"github.com/laundryos/backend/pkg/apiresponse"
	"github.com/laundryos/backend/pkg/jwt"
)

type SettingsHandler struct {
	tenantService *service.TenantService
}

func NewSettingsHandler(tenantService *service.TenantService) *SettingsHandler {
	return &SettingsHandler{
		tenantService: tenantService,
	}
}

func (h *SettingsHandler) RegisterRoutes(r *gin.RouterGroup, jwtManager *jwt.JWTManager) {
	settings := r.Group("/settings", middleware.Auth(jwtManager))
	{
		settings.GET("", h.GetSettings)
		settings.PUT("", h.UpdateSettings)
	}
}

func (h *SettingsHandler) GetSettings(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	tenant, err := h.tenantService.GetTenant(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, apiresponse.NotFound("Settings"))
		return
	}

	c.JSON(http.StatusOK, apiresponse.Success(tenant))
}

func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req service.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiresponse.ValidationError([]apiresponse.ValidationErrorDetail{
			{Field: "body", Message: err.Error()},
		}))
		return
	}

	tenant, err := h.tenantService.UpdateTenant(c.Request.Context(), tenantID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiresponse.InternalError(
			"Gagal update settings",
			"Failed to update settings",
		))
		return
	}

	c.JSON(http.StatusOK, apiresponse.Success(tenant))
}
