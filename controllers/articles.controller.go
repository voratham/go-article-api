package controllers

import (
	"article-api/models"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Articles struct {
	DB *gorm.DB
}

type createArticleRequest struct {
	Title   string                `form:"title" binding:"required"`
	Body    string                `form:"body" binding:"required"`
	Image   *multipart.FileHeader `form:"image" binding:"required"`
	Excerpt string                `form:"excerpt" binding:"required"`
}

type createdArticleResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Excerpt string `json:"Excerpt"`
	Body    string `json:"body"`
	Image   string `json:"image"`
}

func (a *Articles) FineAll(ctx *gin.Context) {

}

func (a *Articles) FindById(ctx *gin.Context) {

}

func (a *Articles) Create(ctx *gin.Context) {
	var data createArticleRequest
	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var article models.Article
	copier.Copy(&article, &data)

	if err := a.DB.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	a.setArticleImage(ctx, &article)

	serializedArticle := createdArticleResponse{}
	copier.Copy(&serializedArticle, &article)

	ctx.JSON(http.StatusCreated, gin.H{"article": serializedArticle})

}

func (a *Articles) setArticleImage(ctx *gin.Context, article *models.Article) error {
	file, err := ctx.FormFile("image")
	if err != nil || file == nil {
		return err
	}

	if article.Image != "" {
		article.Image = strings.Replace(article.Image, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + article.Image)

	}

	path := "uploads/articles/" + strconv.Itoa(int(article.ID))
	os.Mkdir(path, 0755)

	filename := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		log.Fatal("Fail to up load image")
		return err
	}

	article.Image = os.Getenv("HOST") + "/" + filename
	a.DB.Save(article)
	return nil

}
