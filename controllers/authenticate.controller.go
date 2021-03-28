package controllers

import (
	"article-api/models"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Auth struct {
	DB *gorm.DB
}

type authRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type authResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type updateUserRequest struct {
	Email  string                `form:"email"`
	Name   string                `form:"name"`
	Avatar *multipart.FileHeader `form:"avatar"`
}

func (a *Auth) GetProfile(ctx *gin.Context) {
	var serializedUser userResponse
	sub, _ := ctx.Get("sub")
	user := sub.(*models.User)
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})

}

func (a *Auth) UpdateProfile(ctx *gin.Context) {
	var data updateUserRequest

	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
	}

	sub, _ := ctx.Get("sub")
	user := sub.(*models.User)
	var userUpdate models.User
	copier.Copy(&userUpdate, &data)

	if err := a.DB.Model(&user).Updates(&userUpdate).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	}

	setUserImage(ctx, user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)

	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

func (a *Auth) Signup(ctx *gin.Context) {
	var data authRequest
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	copier.Copy(&user, &data)
	user.Password = user.GenerateEncryptedPassword()
	if err := a.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedUser authResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusCreated, gin.H{"user": serializedUser})

}
