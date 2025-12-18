package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/mferdian/Go-GraphQL/graphql/generated"
	"github.com/mferdian/Go-GraphQL/graphql/resolver"
	"github.com/mferdian/Go-GraphQL/domain/product"
	"github.com/mferdian/Go-GraphQL/config/jwt"
	"github.com/mferdian/Go-GraphQL/middleware"
)

func GraphQLRoutes(
	r *gin.Engine,
	productService product.IProductService,
	jwtService jwt.InterfaceJWTService,
) {
	graphqlHandler := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &resolver.Resolver{
					ProductService: productService,
				},
			},
		),
	)

	group := r.Group("/graphql")
	group.Use(middleware.CORSMiddleware())
	// ⬇️ optional kalau mau auth
	// group.Use(middleware.Authentication(jwtService))

	group.POST("", func(c *gin.Context) {
		graphqlHandler.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/playground", func(c *gin.Context) {
		playground.Handler("GraphQL Playground", "/graphql").
			ServeHTTP(c.Writer, c.Request)
	})
}
