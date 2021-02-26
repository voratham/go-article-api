package routes

import (
	"article-api/controllers"

	"github.com/gin-gonic/gin"
)

func Serve(r *gin.Engine) {

	articlesGroup := r.Group("/api/v1/articles")

	articlesControllers := controllers.Articles{}

	{
		articlesGroup.GET("", articlesControllers.FineAll)
		articlesGroup.GET("/:id", articlesControllers.FindById)
		articlesGroup.POST("", articlesControllers.Create)
	}

}
