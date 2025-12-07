package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mferdian/Go-GraphQL/domain/user"
)

func PublicRoutes(r *gin.Engine, userController user.IUserController) {
	public := r.Group("/api")
	public.POST("/register", userController.Register)
	public.POST("/login", userController.Login)
}
