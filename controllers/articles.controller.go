package controllers

import (
	"article-api/models"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Articles struct {
}

type createArticleRequest struct {
	Title string                `form:"title" binding:"required"`
	Body  string                `form:"body" binding:"required"`
	Image *multipart.FileHeader `form:"image" binding:"required"`
}

var articles []models.Article = []models.Article{
	{ID: 1, Title: "Title #1", Body: "Body #1"},
	{ID: 2, Title: "Title #2", Body: "Body #2"},
	{ID: 3, Title: "Title #3", Body: "Body #3"},
	{ID: 4, Title: "Title #4", Body: "Body #4"},
	{ID: 5, Title: "Title #5", Body: "Body #5"},
}

func (a *Articles) FineAll(ctx *gin.Context) {
	result := articles

	if limit := ctx.Query("limit"); limit != "" {
		n, _ := strconv.Atoi(limit)
		result = result[:n]
	}

	ctx.JSON(http.StatusOK, gin.H{"articles": result})
}

func (a *Articles) FindById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	for _, article := range articles {
		if article.ID == uint(id) {
			ctx.JSON(http.StatusOK, gin.H{"article": article})
			return
		}
	}
	errorMessage := fmt.Sprintf("Article with id '%d' not found", id)
	ctx.JSON(http.StatusNotFound, gin.H{"error": errorMessage})

}

func (a *Articles) Create(ctx *gin.Context) {
	var data createArticleRequest
	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article := models.Article{
		ID:    uint(len(articles) + 1),
		Title: data.Title,
		Body:  data.Body,
	}

	file, _ := ctx.FormFile("image")

	path := "uploads/articles/" + strconv.Itoa(int(article.ID))

	os.MkdirAll(path, 0755)

	filename := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong with upload file"})
		return
	}

	article.Image = os.Getenv("HOST") + "/" + filename

	articles = append(articles, article)

	ctx.JSON(http.StatusCreated, gin.H{"article": article})

}
