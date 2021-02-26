package controllers

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type Articles struct {
}

type createArticleRequest struct {
	Title string                `form:"title" binding:"required"`
	Body  string                `form:"body" binding:"required"`
	Image *multipart.FileHeader `form:"image" binding:"required"`
}

func (a *Articles) FineAll(ctx *gin.Context) {

}

func (a *Articles) FindById(ctx *gin.Context) {

}

func (a *Articles) Create(ctx *gin.Context) {

}
