package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laundryos/backend/internal/api/middleware"
	"github.com/laundryos/backend/internal/service"
	"github.com/laundryos/backend/pkg/apiresponse"
	"github.com/laundryos/backend/pkg/jwt"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup, jwtManager *jwt.JWTManager) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", h.Logout)
		auth.GET("/me", middleware.Auth(jwtManager), h.Me)
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiresponse.ValidationError([]apiresponse.ValidationErrorDetail{
			{Field: "body", Message: err.Error()},
		}))
		return
	}

	result, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		if err == service.ErrEmailExists {
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
			"Registrasi Gagal",
			"Registration Failed",
		))
		return
	}

	c.JSON(http.StatusCreated, apiresponse.Success(result))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiresponse.ValidationError([]apiresponse.ValidationErrorDetail{
			{Field: "body", Message: err.Error()},
		}))
		return
	}

	result, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, apiresponse.Error(
				"INVALID_CREDENTIALS",
				"Email atau Password Salah",
				"Invalid Email or Password",
				"",
				"",
			))
			return
		}
		c.JSON(http.StatusInternalServerError, apiresponse.InternalError(
			"Login Gagal",
			"Login Failed",
		))
		return
	}

	c.JSON(http.StatusOK, apiresponse.Success(result))
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req service.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiresponse.ValidationError([]apiresponse.ValidationErrorDetail{
			{Field: "refresh_token", Message: "Refresh token is required"},
		}))
		return
	}

	result, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if err == service.ErrTokenExpired {
			c.JSON(http.StatusUnauthorized, apiresponse.Error(
				"TOKEN_EXPIRED",
				"Token Kadaluarsa",
				"Token Expired",
				"",
				"",
			))
			return
		}
		c.JSON(http.StatusUnauthorized, apiresponse.Unauthorized())
		return
	}

	c.JSON(http.StatusOK, apiresponse.Success(result))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req service.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiresponse.ValidationError([]apiresponse.ValidationErrorDetail{
			{Field: "refresh_token", Message: "Refresh token is required"},
		}))
		return
	}

	err := h.authService.Logout(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiresponse.InternalError(
			"Logout Gagal",
			"Logout Failed",
		))
		return
	}

	c.JSON(http.StatusOK, apiresponse.Success(gin.H{"message": "Logged out successfully"}))
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, apiresponse.Unauthorized())
		return
	}

	user, err := h.authService.GetCurrentUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, apiresponse.NotFound("User"))
		return
	}

	c.JSON(http.StatusOK, apiresponse.Success(user))
}
