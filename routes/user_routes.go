package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mferdian/Go-GraphQL/config/jwt"
	"github.com/mferdian/Go-GraphQL/domain/user"
	"github.com/mferdian/Go-GraphQL/middleware"
)

func UserRoutes(
	r *gin.Engine,
	userController user.IUserController,
	jwtService jwt.InterfaceJWTService,
) {
	user := r.Group("/api/users")
	user.Use(middleware.Authentication(jwtService))

	// --- User Routes ---
	user.PATCH("/:id", userController.UpdateUser)
	user.GET("/:id", userController.GetUserByID)
	user.DELETE("/:id", userController.DeleteUser)
	
}
