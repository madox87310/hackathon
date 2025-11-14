package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService *Service
}

func NewHandler(userService *Service) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	response, err := h.userService.SignUp(c.Request.Context(), req)
	if err != nil {
		log.Printf("h.userService.SignUp: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) SignIn(c *gin.Context) {
	var req SignInRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	response, err := h.userService.SignIn(c.Request.Context(), req)
	if err != nil {
		log.Printf("h.userService.SignIn: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	response, err := h.userService.Refresh(c.Request.Context(), req)
	if err != nil {
		log.Printf("h.userService.Refresh: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if err := h.userService.Logout(c.Request.Context(), req); err != nil {
		log.Printf("h.userService.Logout: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
	}

	c.JSON(http.StatusOK, "successfully logged out")
}
