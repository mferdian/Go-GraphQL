package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mferdian/Go-GraphQL/config/jwt"
	"github.com/mferdian/Go-GraphQL/constants"
	"github.com/mferdian/Go-GraphQL/domain/user"
	"github.com/mferdian/Go-GraphQL/middleware"
)

func AdminRoutes(r *gin.Engine, userController user.IUserController,
	jwtService jwt.InterfaceJWTService) {
	admin := r.Group("/api/users")
	admin.Use(middleware.Authentication(jwtService))
	admin.Use(middleware.AuthorizeRole(constants.ENUM_ROLE_ADMIN))

	// User management
	admin.POST("", userController.CreateUser)
	admin.GET("", userController.GetAllUser)
}
