package handlers

import (
	"monetix-be-api/internal/models"
	"monetix-be-api/internal/repositories"
	"monetix-be-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func SetupUserRoutes(router *gin.Engine, db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := NewUserHandler(userService)

	userGroup := router.Group("/api/v1/users")
	{
		userGroup.POST("/", userHandler.CreateUser)
		// userGroup.GET("/:id", userHandler.GetUser)
		// userGroup.POST("/login", userHandler.Login)
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// func (h *UserHandler) GetUser(c *gin.Context) {
// 	id, err := uuid.Parse(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
// 		return
// 	}

// 	user, err := h.userService.GetUser(id)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, user)
// }
