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

type pagination struct {
	ctx     *gin.Context
	query   *gorm.DB
	records interface{}
}

func (p *pagination) paginate() *pagingResult {

	page, _ := strconv.Atoi(p.ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(p.ctx.DefaultQuery("limit", "12"))
	sortBy := p.ctx.DefaultQuery("sort_by", "id")
	orderBy := p.ctx.DefaultQuery("order_by", "desc")

	ch := make(chan int)
	go p.countRecords(ch)

	offset := (page - 1) * limit
	p.query.Order(sortBy + " " + orderBy).Limit(limit).Offset(offset).Find(p.records)

	count := <-ch
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

func (p *pagination) countRecords(ch chan int) {
	var count int64
	p.query.Model(p.records).Count(&count)
	ch <- int(count)
}
