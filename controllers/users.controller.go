package controllers

import (
	"article-api/config"
	"article-api/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Users struct {
	DB *gorm.DB
}

type createUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type userUpdateRequest struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=8"`
	Name     string `json:"name"`
}

type userResponse struct {
	ID     uint   `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Role   string `json:"role"`
}

type usersPaging struct {
	Items  []userResponse `json:"items"`
	Paging *pagingResult  `json:"paging"`
}

func (u *Users) FineAll(ctx *gin.Context) {
	var users []models.User

	term := ctx.Query("term")

	filter := map[string]string{}

	if term != "" {
		filter["name ILIKE ?"] = "%" + term + "%"
	}

	pagination := pagination{
		ctx:     ctx,
		query:   u.DB,
		records: &users,
		filter:  &filter,
	}

	paging := pagination.paginate()

	serializedUsers := []userResponse{}
	copier.Copy(&serializedUsers, &users)

	ctx.JSON(http.StatusOK, gin.H{"users": usersPaging{Items: serializedUsers, Paging: paging}})
}

func (u *Users) FindById(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

func (u *Users) Create(ctx *gin.Context) {
	var data createUserRequest
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	var newUser models.User
	copier.Copy(&newUser, &data)
	newUser.Password = newUser.GenerateEncryptedPassword()
	if err := u.DB.Create(&newUser).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializeUser userResponse
	copier.Copy(&serializeUser, &newUser)
	ctx.JSON(http.StatusCreated, gin.H{"user": serializeUser})
}

func (u *Users) Update(ctx *gin.Context) {
	var data userUpdateRequest
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if data.Password != "" {
		user.Password = user.GenerateEncryptedPassword()
	}

	var updateUser models.User
	copier.Copy(&updateUser, &data)

	if err := u.DB.Model(&user).Updates(&updateUser).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializeUser userResponse
	copier.Copy(&serializeUser, &user)
	ctx.JSON(http.StatusCreated, gin.H{"user": serializeUser})

}

func (u *Users) Promote(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	user.Promote()
	u.DB.Save(&user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

func (u *Users) Demote(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	user.Demote()
	u.DB.Save(&user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

func (u *Users) Delete(ctx *gin.Context) {
	article, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	u.DB.Delete(&article)

	ctx.Status(http.StatusNoContent)

}

func (u *Users) findUserByID(ctx *gin.Context) (*models.User, error) {
	var user models.User
	id := ctx.Param("id")
	if err := u.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func setUserImage(ctx *gin.Context, user *models.User) error {
	file, err := ctx.FormFile("avatar")
	if err != nil || file == nil {
		return err
	}

	if user.Avatar != "" {
		user.Avatar = strings.Replace(user.Avatar, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + user.Avatar)
	}

	path := "uploads/users/" + strconv.Itoa(int(user.ID))
	os.Mkdir(path, 0755)

	filename := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		log.Fatal("Fail to upload avatar")
		return err
	}

	user.Avatar = os.Getenv("HOST") + "/" + filename
	db := config.GetDB()
	db.Model(&user).Save(&user)

	return nil

}
