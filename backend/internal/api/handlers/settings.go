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
	userService   *service.UserService
}

func NewSettingsHandler(tenantService *service.TenantService, userService *service.UserService) *SettingsHandler {
	return &SettingsHandler{
		tenantService: tenantService,
		userService:   userService,
	}
}

func (h *SettingsHandler) RegisterRoutes(r *gin.RouterGroup, jwtManager *jwt.JWTManager) {
	settings := r.Group("/settings", middleware.Auth(jwtManager))
	{
		settings.GET("", h.GetSettings)
		settings.PUT("", h.UpdateSettings)
	}

	users := r.Group("/users", middleware.Auth(jwtManager))
	{
		users.GET("", h.ListUsers)
		users.GET("/:id", h.GetUser)
		users.POST("", middleware.RequireRole("owner"), h.CreateUser)
		users.PUT("/:id", middleware.RequireRole("owner"), h.UpdateUser)
		users.DELETE("/:id", middleware.RequireRole("owner"), h.DeleteUser)
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

func (h *SettingsHandler) ListUsers(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	page := 1
	limit := 100

	users, total, err := h.userService.GetUsers(c.Request.Context(), tenantID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiresponse.InternalError(
			"Gagal mengambil data user",
			"Failed to get users",
		))
		return
	}

	pagination := apiresponse.CalculatePagination(page, limit, total)
	c.JSON(http.StatusOK, apiresponse.SuccessWithPagination(users, pagination))
}

func (h *SettingsHandler) GetUser(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := c.Param("id")

	user, err := h.userService.GetUserByID(c.Request.Context(), tenantID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, apiresponse.NotFound("User"))
		return
	}

	c.JSON(http.StatusOK, apiresponse.Success(user))
}

func (h *SettingsHandler) CreateUser(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiresponse.ValidationError([]apiresponse.ValidationErrorDetail{
			{Field: "body", Message: err.Error()},
		}))
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), tenantID, &req)
	if err != nil {
		if err == service.ErrEmailAlreadyExists {
			c.JSON(http.StatusConflict, apiresponse.Error(
				"EMAIL_EXISTS",
				"Email Sudah Terdaftar",
				"Email Already Exists",
				"",
				"",
			))
			return
		}
		c.JSON(http.StatusInternalServerError, apiresponse.InternalError(
			"Gagal membuat user",
			"Failed to create user",
		))
		return
	}

	c.JSON(http.StatusCreated, apiresponse.Success(user))
}

func (h *SettingsHandler) UpdateUser(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := c.Param("id")

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiresponse.ValidationError([]apiresponse.ValidationErrorDetail{
			{Field: "body", Message: err.Error()},
		}))
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), tenantID, userID, &req)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, apiresponse.NotFound("User"))
			return
		}
		c.JSON(http.StatusInternalServerError, apiresponse.InternalError(
			"Gagal update user",
			"Failed to update user",
		))
		return
	}

	c.JSON(http.StatusOK, apiresponse.Success(user))
}

func (h *SettingsHandler) DeleteUser(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := c.Param("id")

	err := h.userService.DeleteUser(c.Request.Context(), tenantID, userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, apiresponse.NotFound("User"))
			return
		}
		c.JSON(http.StatusInternalServerError, apiresponse.InternalError(
			"Gagal hapus user",
			"Failed to delete user",
		))
		return
	}

	c.JSON(http.StatusOK, apiresponse.Success(gin.H{"message": "User deleted successfully"}))
}
