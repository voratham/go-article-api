package routes

import (
	"article-api/config"
	"article-api/controllers"

	"github.com/gin-gonic/gin"
)

func Serve(r *gin.Engine) {
	db := config.GetDB()
	articlesGroup := r.Group("/api/v1/articles")

	// dependency inject with db to articles controller
	articlesControllers := controllers.Articles{DB: db}

	{
		articlesGroup.GET("", articlesControllers.FineAll)
		articlesGroup.GET("/:id", articlesControllers.FindById)
		articlesGroup.PATCH("/:id", articlesControllers.Update)
		articlesGroup.POST("", articlesControllers.Create)
	}

}
