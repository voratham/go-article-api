package routes

import (
	"article-api/config"
	"article-api/controllers"
	"article-api/middleware"

	"github.com/gin-gonic/gin"
)

func Serve(r *gin.Engine) {

	db := config.GetDB()
	v1 := r.Group("/api/v1")

	authenticateGroup := v1.Group("auth")
	authenticateControllers := controllers.Auth{DB: db}
	{
		authenticateGroup.POST("/sign-up", authenticateControllers.Signup)
		authenticateGroup.POST("/sign-in", middleware.Authenticate().LoginHandler)
		authenticateGroup.GET("/profile", middleware.Authenticate().MiddlewareFunc(), authenticateControllers.GetProfile)
		authenticateGroup.PATCH("/profile", middleware.Authenticate().MiddlewareFunc(), authenticateControllers.UpdateProfile)
	}

	articlesGroup := v1.Group("/articles")
	// dependency inject with db to articles controller
	articlesControllers := controllers.Articles{DB: db}

	{
		articlesGroup.GET("", articlesControllers.FineAll)
		articlesGroup.GET("/:id", articlesControllers.FindById)
		articlesGroup.PATCH("/:id", articlesControllers.Update)
		articlesGroup.POST("", middleware.Authenticate().MiddlewareFunc(), articlesControllers.Create)
		articlesGroup.DELETE("/:id", articlesControllers.Delete)
	}

	categoriesGroup := v1.Group("/categories")
	categoriesControllers := controllers.Categories{DB: db}

	{
		categoriesGroup.GET("", categoriesControllers.FindAll)
		categoriesGroup.GET("/:id", categoriesControllers.FinById)
		categoriesGroup.PATCH("/:id", categoriesControllers.Update)
		categoriesGroup.POST("", categoriesControllers.Create)
		categoriesGroup.DELETE("/:id", categoriesControllers.Delete)
	}

}
