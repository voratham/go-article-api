package controllers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type pagingResult struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	PrevPage  int `json:"prevPage"`
	NextPage  int `json:"nextPage"`
	Count     int `json:"count"`
	TotalPage int `json:"totalPage"`
}

func pagingResource(ctx *gin.Context, query *gorm.DB, records interface{}) *pagingResult {

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "12"))

	var count int64
	query.Model(records).Count(&count)

	// In case page = 1  and limit = 10
	// page => 1 , 1-10 , offset => 10
	// page => 2 , 11-10 , offset => 20
	offset := (page - 1) * limit
	query.Limit(limit).Offset(offset).Find(records)

	totalPage := int(math.Ceil(float64(count) / float64(limit)))

	var nextPage int
	if page == totalPage {
		nextPage = totalPage
	} else {
		nextPage = page + 1
	}

	return &pagingResult{
		Page:      page,
		Limit:     limit,
		Count:     int(count),
		PrevPage:  page - 1,
		NextPage:  nextPage,
		TotalPage: totalPage,
	}

}
