package controllers

import (
	"article-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Categories struct {
	DB *gorm.DB
}

type createCategoryRequest struct {
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc" binding:"required"`
}

type updateCategoryRequest struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type categoryResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Articles []struct {
		ID    uint   `json:"id"`
		Title string `json:"title"`
	} `json:"articles"`
}

type allCategoriesResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type categoiresPaging struct {
	Items  []allCategoriesResponse `json:"items"`
	Paging *pagingResult           `json:"paging"`
}

func (c *Categories) FindAll(ctx *gin.Context) {
	var categories []models.Category

	pagination := pagination{
		ctx:     ctx,
		query:   c.DB,
		records: &categories,
	}

	paging := pagination.paginate()

	var serializedCategories []allCategoriesResponse
	copier.Copy(&serializedCategories, &categories)

	ctx.JSON(http.StatusOK, gin.H{"categoires": categoiresPaging{Items: serializedCategories, Paging: paging}})
}

func (c *Categories) FinById(ctx *gin.Context) {

	category, err := c.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var serializeCategory categoryResponse
	copier.Copy(&serializeCategory, &category)

	ctx.JSON(http.StatusOK, gin.H{"category": serializeCategory})
}

func (c *Categories) Create(ctx *gin.Context) {
	var data createCategoryRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	copier.Copy(&category, &data)

	if err := c.DB.Create(&category).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedCategory categoryResponse
	copier.Copy(&serializedCategory, &category)

	ctx.JSON(http.StatusCreated, gin.H{"category": serializedCategory})

}

func (c *Categories) Update(ctx *gin.Context) {
	var data updateCategoryRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	category, err := c.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var categoryUpdate models.Category
	copier.Copy(&categoryUpdate, &data)

	if err := c.DB.Model(&category).Updates(&categoryUpdate).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedCategory categoryResponse
	copier.Copy(&serializedCategory, &category)

	ctx.JSON(http.StatusCreated, gin.H{"category": serializedCategory})

}

func (c *Categories) Delete(ctx *gin.Context) {

	category, err := c.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.DB.Delete(&category)
	ctx.Status(http.StatusNoContent)

}

func (c *Categories) findCategoryByID(ctx *gin.Context) (*models.Category, error) {
	var category models.Category
	id := ctx.Param("id")

	if err := c.DB.Preload("Articles").First(&category, id).Error; err != nil {
		return nil, err
	}

	return &category, nil
}
