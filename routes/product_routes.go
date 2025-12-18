package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mferdian/Go-GraphQL/config/jwt"
	"github.com/mferdian/Go-GraphQL/domain/product"
	"github.com/mferdian/Go-GraphQL/middleware"
)

func ProductRoutes(r *gin.Engine,productController product.IProductController, jwtService jwt.InterfaceJWTService) {
	user := r.Group("/api/products")
	user.Use(middleware.Authentication(jwtService))
	
	user.POST("", productController.CreateProduct)
	user.GET("/:id", productController.UpdateProduct)
	user.DELETE("/:id", productController.DeleteProduct)

}
