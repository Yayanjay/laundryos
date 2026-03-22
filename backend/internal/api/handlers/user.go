package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/laundryos/backend/internal/api/middleware"
	"github.com/laundryos/backend/internal/service"
	"github.com/laundryos/backend/pkg/apiresponse"
	"github.com/laundryos/backend/pkg/jwt"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup, jwtManager *jwt.JWTManager) {
	users := r.Group("/users", middleware.Auth(jwtManager))
	{
		users.GET("", h.List)
		users.GET("/:id", h.Get)
		users.POST("", middleware.RequireRole("owner"), h.Create)
		users.PUT("/:id", middleware.RequireRole("owner"), h.Update)
		users.DELETE("/:id", middleware.RequireRole("owner"), h.Delete)
	}
}

func (h *UserHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

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

func (h *UserHandler) Get(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := c.Param("id")

	user, err := h.userService.GetUserByID(c.Request.Context(), tenantID, userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, apiresponse.NotFound("User"))
			return
		}
		c.JSON(http.StatusInternalServerError, apiresponse.InternalError(
			"Gagal mengambil data user",
			"Failed to get user",
		))
		return
	}

	c.JSON(http.StatusOK, apiresponse.Success(user))
}

func (h *UserHandler) Create(c *gin.Context) {
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

func (h *UserHandler) Update(c *gin.Context) {
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
			"Gagal update user",
			"Failed to update user",
		))
		return
	}

	c.JSON(http.StatusOK, apiresponse.Success(user))
}

func (h *UserHandler) Delete(c *gin.Context) {
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
